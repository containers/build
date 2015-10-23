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
