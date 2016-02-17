# acbuild run

`acbuild run` will run the given command inside the ACI.

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

## systemd-nspawn

acbuild currently uses `systemd-nspawn` to run commands inside the ACI. This
means that the machine running acbuild must have systemd installed to be able
to use `acbuild run`. Alternate execution tools (like `runc`) will be added in
the future.

## Exiting out of systemd-nspawn

All acbuild commands can be cancelled with Ctrl+c with the exception of
`acbuild run` once it has executed systemd-nspawn. To break out of a
system-nspawn call, press Ctrl+] three times.
