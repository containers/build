// Copyright 2015 The appc Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package registry

import (
	"archive/tar"
	"compress/bzip2"
	"compress/gzip"
	"crypto/sha512"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/appc/spec/aci"
	"github.com/appc/spec/discovery"
	"github.com/appc/spec/pkg/acirenderer"
	"github.com/appc/spec/schema/types"
	"github.com/coreos/ioprogress"
	"xi2.org/x/xz"

	"github.com/containers/build/util"
)

func (r Registry) tmppath() string {
	return path.Join(r.DepStoreTarPath, "tmp.aci")
}

func (r Registry) tmpuncompressedpath() string {
	return path.Join(r.DepStoreTarPath, "tmp.uncompressed.aci")
}

// Fetch will download the given image, and optionally its dependencies, into
// r.DepStoreTarPath
func (r Registry) Fetch(imagename types.ACIdentifier, labels types.Labels, size uint, fetchDeps bool) error {
	_, err := r.GetACI(imagename, labels)
	if err == ErrNotFound {
		err := r.fetchACIWithSize(imagename, labels, size, fetchDeps)
		if err != nil {
			return err
		}
		return nil
	}
	return err
}

// FetchAndRender will fetch the given image and all of its dependencies if
// they have not been fetched yet, and will then render them on to the
// filesystem if they have not been rendered yet.
func (r Registry) FetchAndRender(imagename types.ACIdentifier, labels types.Labels, size uint) error {
	err := r.Fetch(imagename, labels, size, true)
	if err != nil {
		return err
	}

	filesToRender, err := acirenderer.GetRenderedACI(imagename,
		labels, r)
	if err != nil {
		return err
	}

filesloop:
	for _, fs := range filesToRender {
		_, err := os.Stat(path.Join(r.DepStoreExpandedPath, fs.Key, "rendered"))
		switch {
		case os.IsNotExist(err):
			break
		case err != nil:
			return err
		default:
			// This ACI has already been rendered
			continue filesloop
		}

		err = util.ExtractImage(path.Join(r.DepStoreTarPath, fs.Key),
			path.Join(r.DepStoreExpandedPath, fs.Key), fs.FileMap)
		if err != nil {
			return err
		}

		rfile, err := os.Create(
			path.Join(r.DepStoreExpandedPath, fs.Key, "rendered"))
		if err != nil {
			return err
		}
		rfile.Close()
	}
	return nil
}

func (r Registry) fetchACIWithSize(imagename types.ACIdentifier, labels types.Labels, size uint, fetchDeps bool) error {
	endpoint, err := r.discoverEndpoint(imagename, labels)
	if err != nil {
		return err
	}

	err = r.download(endpoint.ACI, r.tmppath(), string(imagename))
	if err != nil {
		return err
	}

	//TODO: download .asc, verify the .aci with it

	if size != 0 {
		finfo, err := os.Stat(r.tmppath())
		if err != nil {
			return err
		}
		if finfo.Size() != int64(size) {
			return fmt.Errorf(
				"dependency %s has incorrect size: expected=%d, actual=%d",
				size, finfo.Size())
		}
	}

	err = r.uncompress()
	if err != nil {
		return err
	}

	err = os.Remove(r.tmppath())
	if err != nil {
		return err
	}

	id, err := GenImageID(r.tmpuncompressedpath())
	if err != nil {
		return err
	}

	err = os.Rename(r.tmpuncompressedpath(), path.Join(r.DepStoreTarPath, id))
	if err != nil {
		return err
	}

	if !fetchDeps {
		return nil
	}

	err = os.MkdirAll(
		path.Join(r.DepStoreExpandedPath, id, aci.RootfsDir), 0755)
	if err != nil {
		return err
	}

	err = getManifestFromTar(path.Join(r.DepStoreTarPath, id),
		path.Join(r.DepStoreExpandedPath, id, aci.ManifestFile))
	if err != nil {
		return err
	}

	man, err := r.GetImageManifest(id)
	if err != nil {
		return err
	}

	if man.Name != imagename {
		return fmt.Errorf(
			"downloaded ACI name %q does not match expected image name %q",
			man.Name, imagename)
	}

	for _, dep := range man.Dependencies {
		err := r.fetchACIWithSize(dep.ImageName, dep.Labels, dep.Size, fetchDeps)
		if err != nil {
			return err
		}
		if dep.ImageID != nil {
			id, err := r.GetACI(dep.ImageName, dep.Labels)
			if err != nil {
				return err
			}
			if id != dep.ImageID.String() {
				return fmt.Errorf("dependency %s doesn't match hash",
					dep.ImageName)
			}
		}
	}
	return nil
}

// Need to uncompress the file to be able to generate the Image ID
func (r Registry) uncompress() error {
	acifile, err := os.Open(r.tmppath())
	if err != nil {
		return err
	}
	defer acifile.Close()

	typ, err := aci.DetectFileType(acifile)
	if err != nil {
		return err
	}

	// In case DetectFileType changed the cursor
	_, err = acifile.Seek(0, 0)
	if err != nil {
		return err
	}

	var in io.Reader
	switch typ {
	case aci.TypeGzip:
		in, err = gzip.NewReader(acifile)
		if err != nil {
			return err
		}
	case aci.TypeBzip2:
		in = bzip2.NewReader(acifile)
	case aci.TypeXz:
		in, err = xz.NewReader(acifile, 0)
		if err != nil {
			return err
		}
	case aci.TypeTar:
		in = acifile
	case aci.TypeText:
		return fmt.Errorf("downloaded ACI is text, not a tarball")
	case aci.TypeUnknown:
		return fmt.Errorf("downloaded ACI is of an unknown type")
	}

	out, err := os.OpenFile(r.tmpuncompressedpath(),
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return fmt.Errorf("error copying: %v", err)
	}

	err = out.Sync()
	if err != nil {
		return fmt.Errorf("error writing: %v", err)
	}

	return nil
}

func getManifestFromTar(tarpath, dst string) error {
	tarfile, err := os.Open(tarpath)
	if err != nil {
		return err
	}
	defer tarfile.Close()

	tr := tar.NewReader(tarfile)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// End of tar reached
			break
		}
		if err != nil {
			return err
		}
		switch {
		case hdr.Typeflag == tar.TypeReg:
			if hdr.Name == aci.ManifestFile {
				f, err := os.OpenFile(dst,
					os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
				if err != nil {
					return err
				}
				defer f.Close()
				_, err = io.Copy(f, tr)
				if err != nil {
					return err
				}
				return nil
			}
		default:
			continue
		}
	}
	return fmt.Errorf("manifest not found in ACI")
}

func GenImageID(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	h := sha512.New()

	_, err = io.Copy(h, file)
	if err != nil {
		return "", err
	}

	s := h.Sum(nil)

	return fmt.Sprintf("sha512-%x", s), nil
}

func (r Registry) discoverEndpoint(imageName types.ACIdentifier, labels types.Labels) (*discovery.ACIEndpoint, error) {
	labelmap := make(map[types.ACIdentifier]string)
	for _, label := range labels {
		labelmap[label.Name] = label.Value
	}

	app, err := discovery.NewApp(string(imageName), labelmap)
	if err != nil {
		return nil, err
	}
	if _, ok := app.Labels["arch"]; !ok {
		app.Labels["arch"] = runtime.GOARCH
	}
	if _, ok := app.Labels["os"]; !ok {
		app.Labels["os"] = runtime.GOOS
	}

	insecure := discovery.InsecureNone
	if r.Insecure {
		insecure = discovery.InsecureHTTP
	}

	acis, attempts, err := discovery.DiscoverACIEndpoints(*app, nil, insecure, 0)
	if err != nil {
		return nil, err
	}
	if r.Debug {
		for _, a := range attempts {
			fmt.Fprintf(os.Stderr, "meta tag not found on %s: %v\n",
				a.Prefix, a.Error)
		}
	}
	if len(acis) == 0 {
		return nil, fmt.Errorf("no endpoints discovered to download %s",
			imageName)
	}

	return &acis[0], nil
}

func (r Registry) download(url, path, label string) error {
	//TODO: auth
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	transport := http.DefaultTransport
	transport.(*http.Transport).Proxy = http.ProxyFromEnvironment
	if r.Insecure {
		transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	client := &http.Client{Transport: transport}
	//f.setHTTPHeaders(req, etag)

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) >= 10 {
			return fmt.Errorf("too many redirects")
		}
		//f.setHTTPHeaders(req, etag)
		return nil
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		break
	case http.StatusNotFound:
		return ErrNotFound
	default:
		return fmt.Errorf("bad HTTP status code: %d", res.StatusCode)
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}

	reader := newIoprogress(label, res.ContentLength, res.Body)

	_, err = io.Copy(out, reader)
	if err != nil {
		return fmt.Errorf("error copying %s: %v", label, err)
	}

	err = out.Sync()
	if err != nil {
		return fmt.Errorf("error writing %s: %v", label, err)
	}

	err = out.Close()
	if err != nil {
		return err
	}

	return nil
}

func newIoprogress(label string, size int64, rdr io.Reader) io.Reader {
	prefix := "Downloading " + label
	fmtBytesSize := 18

	// if barSize < 2, drawing the bar will panic; 3 will at least give a spinny
	// thing.
	barSize := int64(80 - len(prefix) - fmtBytesSize)
	if barSize < 2 {
		barSize = 2
	}

	bar := ioprogress.DrawTextFormatBarForW(barSize, os.Stderr)
	fmtfunc := func(progress, total int64) string {
		// Content-Length is set to -1 when unknown.
		if total == -1 {
			return fmt.Sprintf(
				"%s: %v of an unknown total size",
				prefix,
				ioprogress.ByteUnitStr(progress),
			)
		}
		return fmt.Sprintf(
			"%s: %s %s",
			prefix,
			bar(progress, total),
			ioprogress.DrawTextFormatBytes(progress, total),
		)
	}

	return &ioprogress.Reader{
		Reader:       rdr,
		Size:         size,
		DrawFunc:     ioprogress.DrawTerminalf(os.Stderr, fmtfunc),
		DrawInterval: time.Second,
	}
}
