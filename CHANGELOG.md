## v0.3.0

appc spec version: v0.7.4

- Sets the PATH environment variable inside of acbuild run commands, instead of
  propagating the variable from the host.
- Copy has been split out into copy and copy-to-dir. copy is for copying a
  given file/directory to a specific path, and copy-to-dir is for copying one
  or more files/directories into a given directory.  error.
- The run and end subcommands now contain a check for if .acbuild/target is
  mounted, and will attempt to unmount it before performing their actions.
- The flag --working-dir has been added to the run subcommand, for specifying
  the working directory to run the given command in. This requires systemd >=
  229.
- Added support for passing the begin command a tarball, which when not a valid
  ACI will be used as the starting rootfs for the build.
- Now uses version 0.7.4 of the AppC spec.

## v0.2.2

appc spec version: v0.7.1+git

v0.2.2 is a minor release incorporating a handful of changes since v0.2.1.
Things this release includes:

- Don't require a starting character of `.` or `/` when specifying an image for
  the `--modify` flag.
- Adds support for specifying a local directory with `begin`, for use as the
  initial rootfs of the ACI.
- Adds support for the `run` command on systems with a version of systemd <
  209.
- Versions of dependencies can be specified using tag syntax.
- A bug related to the handling of hard links in tars was found, and fixed.
- A history of acbuild commands performed is now tracked in annotations in the
  ACI being worked on. This can be disabled via a flag.

## v0.2.1

appc spec version: v0.7.1+git5a7af19

v0.2.0 changed the default behavior of the copy command, requiring that the
user passed the `-T` flag to use the previous default behavior. This release
flips that, such that the default will be as it was before, and a new flag must
be passed to use the new behavior.

## v0.2.0

appc spec version: v0.7.1+git5a7af19

- Automatically fills in the `os` and `arch` labels on new builds.
- The `--modify` flag failed to properly use an alternate work context.
- The help page has been reformatted.
- Extracting ACIs during begin no longer requires root.
- Modifying an existing ACI that didn't have the `app` field set would cause an
  error.
- Copy now has two modes of operation, toggled with the `-T` flag. The new mode
  allows copying many files at once into the ACI, as opposed to a single one.
- When using begin with a local file, the path must now start with a `.`, `~`,
  or `/`.

The following commands have been added to acbuild since the last release:
- cat-manifest
- isolator add
- isolator remove
- replace-manifest
- set-event-handler pre-start
- set-event-handler post-stop
- set-working-directory

## v0.1.1

appc spec version: v0.7.1+git5a7af19

v0.1.1 has a couple of bug fixes since the initial release:

- When the `--debug` flag was combined with the `run` command acbuild would
  panic. (https://github.com/appc/acbuild/issues/56)
- When a remote ACI was specified with `begin`, acbuild would emit an error and
  exit after downloading the ACI. (https://github.com/appc/acbuild/pull/58)

## v0.1.0

appc spec version: v0.7.1+git5a7af19

v0.1.0 is the initial release of acbuild.


Currently functioning:
- annotation add
- annotation remove
- begin
- copy
- dependency add
- dependency remove
- end
- environment add
- environment remove
- label add
- label remove
- mount add	
- mount remove
- port add
- port remove
- run
- set-exec
- set-group
- set-name
- set-user
- write
