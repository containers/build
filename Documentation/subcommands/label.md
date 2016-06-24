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

## The Arch Label on ARM

When `acbuild begin` is used without any starting image, `runtime.GOARCH` is
used to determine the default value for the `arch` label. This doesn't
differentiate between different arm versions, and according to the AppC spec
acbuild needs to pick one of `armv6l`, `armv7l`, and `armv7b`.

To account for this acbuild will inspect the endinanness of the machine it is
running on and pick `armv7l` or `armv7b` accordingly. As a result if acbuild is
on a platform where the label should be `armv6l`, the following command should
be used to override the default label.

```
acbuild label add arch armv6l
```

## Examples

```bash
acbuild label add version latest

acbuild label rm os
```
