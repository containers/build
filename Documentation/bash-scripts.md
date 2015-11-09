# Bash scripts

Bash scripting provides a common, simple, and well-known mechanism to drive
`acbuild`. A script template, below, is used for the examples in this
repository. It makes working with acbuild a little nicer by exiting the script
when it hits an error, and calling `acbuild end` whenever the script exits.

Since the script leverages `bash` features, we open it by specifying execution
in the bash shell, rather than just inheriting the user's `SHELL` environment.

Further, we set the `-e` option to bash to ensure that the entire script exits
on the failure of any command. This gives us some atomicity, ensuring that
either a complete and valid ACI is constructed, or none at all

The `begin` line starts the build. If the build is to be started from an
existing ACI, this line will be different.

The rest of the script is concerned with cleanup and error handling. The
`acbuildend` function is interesting, as it serves as a simple "catch" mechanism
for the script, called on a trap triggered by either reaching the script's end,
or by a non-zero return value from any command in the script. `acibuildend`
stores the last command's exit code, then calls `acbuild --debug end` to
terminate the build, passing the exit code, if any, through the script exit, to
help with troubleshooting.

```bash
#!/usr/bin/env bash
set -e

acbuildend () {
    export EXIT=$?;
    acbuild --debug end && exit $EXIT;
}

acbuild --debug begin
trap acbuildend EXIT

# User entered acbuild commands go here
```
