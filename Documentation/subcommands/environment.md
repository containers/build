# acbuild environment

`acbuild environment` is used to set environment variables in the ACI. Each
variable name must be unique.

## Subcommands

* `acbuild environment add NAME VALUE`

  Updates the ACI to contain an environment variable with the given name and
  value. If the variable already exists, its value will be changed.

* `acbuild environment remove NAME`

  Removes the environment variable with the given name from the ACI.

## Examples

```bash
acbuild environment add REDUCE_WORKER_DEBUG true

acbuild environment remove LANG
```
