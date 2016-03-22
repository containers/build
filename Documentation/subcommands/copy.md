# acbuild copy

`acbuild copy` will copy one file or directory from the local filesystem into
the ACI.

It takes exactly two arguments, the first of which is the path on the
local system to copy from, and the second is the path inside the ACI to copy
to. If the target path's parent directory does not exist in the filesystem, it
will be implicitly created (along with any necessary parent directories).

The following two commands would do the same thing:

```bash
acbuild copy ./nginx.conf /etc/nginx/nginx.conf
```

```bash
cp ./nginx.conf ./.acbuild/currentaci/rootfs/etc/nginx/nginx.conf
```
