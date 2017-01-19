# acbuild write

`acbuild write` will produce an ACI from the current build context. This can be
called an arbitrary number of times during a build, but in most cases it should
probably be called at least once.

## Writing the image

`acbuild write` requires one argument: the file to write the ACI to. If the
file exists, acbuild will refuse to overwrite the file unless the `--overwrite`
flag is used.
