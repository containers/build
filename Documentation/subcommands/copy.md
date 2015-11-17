# acbuild copy

`acbuild copy` will copy files and directories from the local filesystem into
the ACI.

There are two modes of operation for `acbuild copy`, one for copying multiple
files/directories into a specified directory, and one for copying a single
file/directory to a specified path.

## Default Mode

When the `-T` flag isn't specified, `acbuild copy` takes any number of arguments
greater than or equal to 2. The last specified argument is the directory inside
the ACI to put the files in, and all other arguments are paths on the host to
copy. If the target directory doesn't exist in the ACI, it will be implicitly
created.

The following commands would do the same thing:

```bash
acbuild copy apache.conf sites-available/00-default sites-available/myblog /etc/apache2
```

```bash
cp apache.conf sites-available/00-default sites-available/myblog /etc/apache2
```

## Explicit Target Mode

When the `-T` flag is specified, it copies one thing from the host to a
specified path. It takes exactly two arguments, the first of which is the path
on the local system to copy from, and the second is the path inside the ACI to
copy to. If the directory the target path is in doesn't exist, it will be
implicitly created.

The following two commands would do the same thing:

```bash
acbuild copy -T ./nginx.conf /etc/nginx/nginx.conf
```

```bash
cp ./nginx.conf ./.acbuild/currentaci/rootfs/etc/nginx/nginx.conf
```
