# acbuild run

`acbuild run` will run the given command inside the ACI.

## --working-dir

The `--working-dir` flag can be used to specify the working directory for the
command being run inside the image.

## Options Parsing

acbuild needs to be able to differentiate between flags to acbuild and flags to
pass along to the binary being run. This is accomplished with `--`. Any flags
occurring before this are considered as being intended for acbuild, and any
flags after it are assumed to belong to the command being run.

## Dependencies

In order to be able to run the command, all dependencies of the current ACI
must be fetched. The first time `run` is called, the dependencies will be
downloaded and expanded.

## Overlayfs

acbuild utilizes overlayfs when running a command in an ACI with dependencies.
This is so that acbuild is able to separate out the files from the dependencies
and the files in your ACI after the command finishes running.

Obviously this is not necessary when there are no dependencies. If `acbuild
run` is to be used on a system without overlayfs, the ACI and its dependencies
must be flattened into a single ACI without dependencies. A command called
`acbuild squash` is being worked on to do this.

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
