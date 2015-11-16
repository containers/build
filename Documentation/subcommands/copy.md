# acbuild copy

`acbuild copy` will copy a file or directory from the local filesystem into the
ACI. The first argument is the path on the local system to copy from, and the
second argument is the path inside the ACI to copy to.

The following two commands should do the same thing:

```bash
acbuild copy ./nginx.conf /etc/nginx/nginx.conf
```

```bash
cp ./nginx.conf ./.acbuild/currentaci/rootfs/etc/nginx/nginx.conf
```
