# Command History

By default, acbuild will maintain annotations in the image it's modifying to
track the acbuild commands that were performed on that image. The name of the
annotation follows the pattern `coreos.com/acbuild/command-X` where X is a
number, and the value of the annotation will be the acbuild command that was
called.

For example, the following acbuild commands:

```bash
acbuild begin
acbuild set-name example.com/nginx
acbuild dep add quay.io/coreos/alpine-sh
acbuild run -- apk update
acbuild run -- apk add nginx
acbuild write nginx.aci
acbuild end
```

result in the manifest in `nginx.aci` containing the following annotations:

```json
[
    {
        "name": "appc.io/acbuild/command-1",
        "value": "acbuild set-name \"example.com/nginx\""
    },
    {
        "name": "appc.io/acbuild/command-2",
        "value": "acbuild dependency add \"quay.io/coreos/alpine-sh\""
    },
    {
        "name": "appc.io/acbuild/command-3",
        "value": "acbuild run \"apk\" \"update\""
    },
    {
        "name": "appc.io/acbuild/command-4",
        "value": "acbuild run \"apk\" \"add\" \"nginx\""
    }
]
```

This command tracking can easily be turned off, by providing the `--no-history`
flag to any command that should not generate this additional annotation.
