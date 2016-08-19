# acbuild begin

`acbuild begin` will start a new build.

## Location of work context

By default, information about the build is stored at `.acbuild` in the current
working directory. If the current directory changes during the build, `acbuild`
will be unaware of, and unable to operate on, the build that was started until
the current directory is changed back to the location where `acbuild begin` was
run. If this is undesirable, the `--work-path` flag can be provided to specify
the location to store and access the build context.

## Starting with an empty ACI

The build will default to starting with an empty ACI. The rootfs will be empty,
and the manifest will look something like the following:

```json
{
    "acKind": "ImageManifest",
    "acVersion": "0.7.1+git",
    "name": "acbuild-unnamed",
    "labels": [
        {
            "name": "arch",
            "value": "amd64"
        },
        {
            "name": "os",
            "value": "linux"
        }
    ]
}
```

The `arch` and `os` labels are filled in with the architecture and operating
system of the machine acbuild is running on. If this is undesirable, the labels
can be modified or removed with the `acbuild label` command.

## Starting with a pre-existing ACI

The begin command can also be passed an ACI, either on the file system or an
image name to fetch via [meta
discovery](https://github.com/appc/spec/blob/master/spec/discovery.md#meta-discovery).
When an ACI is specified, it is used as the starting point for the build as
opposed to an empty image. The ACI's manifest and rootfs will both come from
the specified image.  If the image is to be fetched via meta discovery over
http (as opposed to https), the `--insecure` flag must be used.

As before, if the ACI to begin from is on the local filesystem the path to it
must start with `.`, `~`, or `/`. As an example, if the ACI is in the current
directory, then instead of passing in `alpine-latest-linux-amd64.aci`, what
would be passed in is `./alpine-latest-linux-amd64.aci`.

## Starting with a local, pre-existing rootfs

A build can be started with a rootfs of a container. This rootfs could perhaps
be produced by a tool like [buildroot](http://buildroot.org/), or downloaded
from somewhere like the [Ubuntu Core
releases](http://cdimage.ubuntu.com/ubuntu-core/xenial/daily-preinstalled/current/).

The rootfs can either be in a local directory, or a local tar file. In either
case, the rootfs is copied into the container and an empty manifest is created
that is identical to beginning with an empty ACI.

When specifying something on the local filesystem to the begin command, the
path to it _must_ start with `.`, `~`, or `/`. As an example, if the directory
is in the current directory, then instead of passing in `buildroot/output`,
what would be passed in is `./buildroot/output`.

## Examples

```bash
acbuild begin
acbuild begin ./my-app.aci
acbuild begin quay.io/coreos/alpine-sh
acbuild --work-path /tmp/mybuild begin
acbuild begin ~/projects/buildroot/output/target
acbuild begin ./ubuntu-core-14.04-core-amd64.tar.gz
```
