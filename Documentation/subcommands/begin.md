# acbuild begin

`acbuild begin` will start a new build. A build stores information about the
image being built (its files, dependencies, manifest, etc.) in a work context on
disk. This context can be created from scratch, and thus have no files and a
skeleton manifest, or from an existing image, giving this context the files,
manifest, and so on from the specified image.

## Location of work context

By default, information about the build is stored at `.acbuild` in the current
working directory. If the current directory changes during the build, `acbuild`
will be unaware of, and unable to operate on, the build that was started until
the current directory is changed back to the location where `acbuild begin` was
run. If this is undesirable, the `--work-path` flag can be provided to specify
the location to store and access the build context.

## Picking a build mode

acbuild can produce both AppC and OCI images. Which one a build is going to
produce is specified when the build is started. Some commands are only available
in one build mode and not another, and other commands are available in both but
will have somewhat different behavior depending on the mode.

The build mode is selected with the `--build-mode` flag, and it will accept
either `appc` or `oci` as a value. If unspecified, it will default to `appc`.

## Starting with an empty image

If no additional arguments are provided to `acbuild begin`, the build will be
started with an empty image. There will be no files in the rootfs for the
container, and the manifest will include almost no information.

## Starting with a local rootfs

A build can be started with a rootfs of a container. This rootfs could perhaps
be produced by a tool like [buildroot][1], or downloaded from somewhere like the
[Ubuntu Core releases][2].

The rootfs can either be in a local directory, or a local tar file. In either
case, the rootfs is copied into the container and an empty manifest is created
that is identical to beginning with an empty image.

When specifying something on the local filesystem to the begin command, the path
to it _must_ start with `.`, `~`, or `/`. As an example, if the directory is in
the current directory, then instead of passing in `buildroot/output`, what would
be passed in is `./buildroot/output`. This is necessary for acbuild to be able
to differentiate between a local path and a remote image name that acbuild is to
fetch.

## Starting with a pre-existing image

When a build is started, an entire image can also be provided as a starting
point. The build will modify the provided image, so at the start of the build
all of the given image's files will be in the build context as will the image's
manifest and other related information.

A local image on disk can be specified with a path (again, this path _must_
start with `.`, `~`, or `/`).

A remote image can also be specified, and acbuild will download the image and
then work on it.

When in the appc build mode an image name can be specified, and acbuild will
perform [AppC discovery][3] to convert this into a URL it will download.
Additionally, if in appc build mode, a name can be prefixed with `docker://` and
acbuild will use the [docker2aci project][4] to fetch and convert a docker image
into an ACI, and then use that to begin the build.

Remote image fetching in OCI is currently unsupported.

## Examples

```bash
acbuild begin
acbuild begin ./my-app.aci
acbuild begin --build-mode oci ./my-app.oci
acbuild begin quay.io/coreos/alpine-sh
acbuild begin --build-mode appc docker://alpine
acbuild --work-path /tmp/mybuild begin
acbuild begin ~/projects/buildroot/output/target
acbuild begin --build-mode oci ./ubuntu-core-14.04-core-amd64.tar.gz
```

[1]: http://buildroot.org/
[2]: http://cdimage.ubuntu.com/ubuntu-base/xenial/daily/current/
[3]: https://github.com/appc/spec/blob/master/spec/discovery.md
[4]: https://github.com/appc/docker2aci/
