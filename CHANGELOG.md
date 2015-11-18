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
