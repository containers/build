# acbuild copy-to-dir

`acbuild copy-to-dir` will copy any number of files and directories from the
local filesystem into the ACI.

It takes at least two arguments, where all but the final argument are paths on
the host system. The final argument is a parent directory inside the ACI to
place the files and directories in.

If the target directory doesn't exist in the ACI, it will be implicitly created
(along with any necessary parent directories).

The following commands would do the same thing:

```bash
acbuild copy-to-dir apache.conf sites-available/00-default sites-available/myblog /etc/apache2
```

```bash
cp apache.conf sites-available/00-default sites-available/myblog ./.acbuild/currentaci/rootfs/etc/apache2
```

