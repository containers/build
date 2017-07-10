# acbuild run

`acbuild run` will run the given command inside the image.

## --working-dir

The `--working-dir` flag can be used to specify the working directory for the
command being run inside the image.

## Options Parsing

acbuild needs to be able to differentiate between flags to acbuild and flags to
pass along to the binary being run. This is accomplished with `--`. Any flags
occurring before this are considered as being intended for acbuild, and any
flags after it are assumed to belong to the command being run.

## Dependencies

In order to be able to run the command, all dependencies of the image must be
available. In the appc build mode, any missing dependencies will be downloaded
the first time `run` is called. In the oci build mode, all layers should be
available if the image was pulled from a remote registry. If the build was
started from a local image and not all layers are present, `run` will be unable
to run and exit with an error.

## Overlayfs

acbuild utilizes overlayfs when running a command in an image with layers.
This is so that acbuild is able to separate out the files from lower layers
and the files belonging to the top layer after the command finishes running.

Obviously this is not necessary when there is only one layer. If `acbuild run`
is to be used on a system without overlayfs, the image and its dependencies must
be flattened into a single layer without dependencies. A command called `acbuild
squash` is being worked on to do this.

## Engines

acbuild can use different engines to perform the actual execution of the given
command. The flag `--engine` can be used to select a non-default engine.

### systemd-nspawn

The default engine in acbuild is called `systemd-nspawn`, which rather
obviously uses `systemd-nspawn` to run the given command. This means that the
machine running acbuild must have systemd installed to be able to use `acbuild
run` with the default engine.

### chroot

An alternative engine is called `chroot`, which uses the chroot syscall to
enter into the container and run the specified command. There's no namespacing
involved, so the command will be able to see and possibly interact with other
processes on the host. This engine notably has no dependency on systemd, unlike
the `systemd-nspawn` engine.

### Exiting out of systemd-nspawn

All acbuild commands can be cancelled with Ctrl+c with the exception of
`acbuild run` once it has executed systemd-nspawn. To break out of a
system-nspawn call, press Ctrl+] three times.
