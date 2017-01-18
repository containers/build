# acbuild mount

Mount points can be specified in an image's manifest. These are locations within
the image's rootfs that the app expects to have external data mounted to.

In OCI, these are referred to as volumes, and acbuild has an alias allowing any
mount commands to be referred to as volume commands (as in, `acbuild mount add`
and `acbuild volume add` do the same thing)

## Subcommands

* `acbuild mount add NAME PATH`

  Updates the image to contain a mount point with the given name and path. If
  the mount point already exists, its path will be changed.

* `acbuild mount remove NAME/PATH`

  Removes the mount point with the given name or path from the image.

## Flags

- `--read-only`: when specified, the data mounted into the image's rootfs should
  be mounted as read only. This is unsupported in the oci build mode.

## Examples

```bash
acbuild mount add source /root/source

acbuild volume add html /usr/share/nginx/html --read-only

acbuild mount remove work
```
