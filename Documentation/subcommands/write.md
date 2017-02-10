# acbuild write

`acbuild write` will produce an image from the current build context. This can
be called an arbitrary number of times during a build, but in most cases it
should probably be called once at the end of the build.

## Writing the image

`acbuild write` requires one argument: the file to write the image to. If the
file exists, acbuild will refuse to overwrite the file unless the `--overwrite`
flag is used.

The format the resulting image will be written in is dependent on what build
mode was specified when the build was started.
