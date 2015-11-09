# acbuild end

`acbuild end` will end the current build. This is accomplished by simply
deleting the directory the build context is stored in, which is `.acbuild` in
either the current directory or the directory specified via the `--work-path`
flag.

If the build was a success and an ACI is to be produced, the `write` command
must be called before `end`, otherwise the build will be lost.
