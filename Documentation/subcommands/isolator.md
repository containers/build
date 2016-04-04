# acbuild isolator

Isolators can be specified on an ACI. They represent a list of isolation steps
that should be performed when the ACI is run.

## Subcommands

* `acbuild isolator add NAME FILE`

  Updates the ACI to contain an isolator with the given name, and the value of
  the contents of the given file. If an isolator with the name already exists,
  its value will be changed.

  If `-` is used for FILE, the value for the isolator is read in from stdin.

* `acbuild isolator remove NAME`

  Removes the isolator with the given name from the ACI.

## Linux Capabilities

One very common usage of isolators is to grant a container a Linux capability.
This is done with an isolator named `os/linux/capabilities-retain-set`.

```
echo '{ "set": ["CAP_IPC_LOCK"] }' | acbuild isolator add "os/linux/capabilities-retain-set" -
```
