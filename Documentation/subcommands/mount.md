# acbuild mount

Mount points can be specified in an ACI's manifest. These are locations withing
the ACI's rootfs that the app expects to have external data mounted to.

## Subcommands

* `acbuild mount add NAME PATH`

  Updates the ACI to contain a mount point with the given name and path. If the
  mount point already exists, its path will be changed.

* `acbuild mount remove NAME`

  Removes the mount point with the given name from the ACI.

## Flags

- `--read-only`: when specified, the data mounted into the ACI's rootfs should
  be mounted as read only.

## Examples

```bash
acbuild mount add source /root/source

acbuild mount add html /usr/share/nginx/html --read-only

acbuild mount remove work
```
