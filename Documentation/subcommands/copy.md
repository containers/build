# acbuild copy

`acbuild copy` will copy files and directories from the local filesystem into
the ACI.

There are two modes of operation for `acbuild copy`, one for copying multiple
files/directories into a specified directory, and one for copying a single
file/directory to a specified path.

## Default Mode

By default, `acbuild copy` will copy one thing from the host to a specified
path. It takes exactly two arguments, the first of which is the path on the
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

## Explicit Target Mode

When the `--to-dir` flag is used, `acbuild copy` takes any number of arguments
greater than or equal to 2. The last specified argument is the directory inside
the ACI to put the files in, and all other arguments are paths on the host to
copy. If the target directory doesn't exist in the ACI, it will be implicitly
created (along with any necessary parent directories).

The following commands would do the same thing:

```bash
acbuild copy --to-dir apache.conf sites-available/00-default sites-available/myblog /etc/apache2
```

```bash
cp apache.conf sites-available/00-default sites-available/myblog /etc/apache2
```

