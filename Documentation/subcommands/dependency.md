# acbuild dependency

Dependencies are ACIs separate from the current ACI, that are placed down into
the rootfs before the files from the current ACI.

The ordering of an ACI's dependencies is significant. Let's say the current ACI
has two dependencies, A and B, where A comes before B. If both A and B contain
the file `/bin/sh`, and the current ACI does not have this file, then the
`/bin/sh` from B is what will appear in the container when it is run.

At the time of writing acbuild simply puts the dependency in the manifest in
the order they were added. If the ordering of an ACI's dependencies is
incorrect, the user must remove each dependency manually and then add them in
the correct order.

## Subcommands

* `acbuild dependency add IMAGE_NAME`

  Updates the ACI to contain a dependency with the given name.

* `acbuild dependency remove IMAGE_NAME`

  Removes the dependency with the given image name from the ACI.

## Flags

The `add` command also has the following optional flags:

- `--image-id`: sets the content hash of the dependency being added. When this
  ACI is run, the retrieved dependency must match this hash.

- `--label`: adds a label to the dependency being added. This is used when
  determining which image to fetch for the dependency.

- `--size`: the size of the dependency being added, in bytes. When this ACI is
  run, the retrieved dependency must have this size.

## Examples

```bash
acbuild dependency add example.com/alpine

acbuild dependency add example.com/ubuntu --image-id sha512-...

acbuild dependency add example.com/nodejs --label version=4.0.0 --label arch=noarch

acbuild dependency remove example.com/centos
```
