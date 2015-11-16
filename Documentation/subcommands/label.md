# acbuild label

Labels are a part of an ACI's manifest that are used during image discovery and
dependency resolution. Each label has a name and a value, and each name must be
unique.

## Subcommands

* `acbuild label add NAME VALUE`

  Updates the ACI to contain a label with the given name and value. If the label
  already exists, its value will be changed.

* `acbuild label remove NAME`

  Removes the label with the given name from the ACI.

## Common Labels

Common labels include:

- `version`: the version of this ACI. Ideally when combined with the current
  ACI's name this will be unique for every build of an app on a given OS and
  architecture.
- `os`: the operating system the ACI is built for.
- `arch`: the architecture the ACI is built for.

## Default Labels

When an empty ACI is created with `acbuild begin`, by default the `os` and
`arch` labels are created for you. Their default values are the current
system's OS and architecture, as determined by golang's `runtime` package.

## Examples

```bash
acbuild label add version latest

acbuild label rm os
```
