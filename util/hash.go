package util

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"path"
)

func HashBlob(blob []byte) string {
	h := sha256.New()
	h.Write(blob)
	hashBytes := h.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

func MarshalHashAndWrite(ociPath string, data interface{}) (string, string, int, error) {
	blob, err := json.Marshal(data)
	if err != nil {
		return "", "", 0, err
	}
	hash := HashBlob(blob)
	err = ioutil.WriteFile(path.Join(ociPath, "blobs", "sha256", hash), blob, 0644)
	if err != nil {
		return "", "", 0, err
	}
	return "sha256", hash, len(blob), nil
}
