# acbuild port

Ports can be specified in an image's manifest. This allows for easy mapping of
ports inside the image to ports on the host, when the app is run in a separate
network namespace, among other uses.

## Subcommands

* `acbuild port add NAME PROTOCOL PORT`

  Updates the image to contain a port with the given name, protocol, and port.
  The protocol is either `udp` or `tcp`. If the port already exists, its values
  will be changed.

* `acbuild port remove NAME/PORT`

  Removes the port with the given name or number from the image.

## Flags

`acbuild port add` supports the following flags when in appc build mode:

- `--count`: when specified, represents a range of ports as opposed to a single
  one. The range starts at the port being added, and has a size of the given
  number.
- `--socket-activated`: when set, the application expects to be socket
  activated on the given ports.

## Examples:

```bash
acbuild port add http tcp 80

acbuild port add dns udp 53

acbuild port remove tftp
```
