# acbuild set-exec

* `acbuild set-exec -- CMD [ARGS]`

  Sets the exec command in the image's manifest.

## Options Parsing

acbuild needs to be able to differentiate between flags to acbuild and flags to
pass along to the binary being run. This is accomplished with `--`. Any flags
occurring before this are considered as being intended for acbuild, and any
flags after it are assumed to belong to the command being run.
