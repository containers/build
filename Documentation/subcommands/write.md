# acbuild write

`acbuild write` will produce an ACI from the current build context. This can be
called an arbitrary number of times during a build, but in most cases it should
probably be called at least once.

## Writing the image

`acbuild write` requires one argument: the file to write the ACI to. If the
file exists, acbuild will refuse to overwrite the file unless the `--overwrite`
flag is used.

## Signing the image

`acbuild write` can exec the `gpg` command on your system to sign the ACI for you. If the `--sign` flag is used without any other arguments like so:

```bash
acbuild write mycoolapp.aci --sign
```

acbuild would run the following command after it has finished writing out the
ACI:

```bash
gpg --armor --yes --output mycoolapp.aci.asc --detach-sig mycoolapp.aci
```

If the user who is running acbuild has their gpg keyring configured correctly
this should Just Work™, but there will be cases where this isn't sufficient. To
allow the user to have control over the gpg command run, arguments to replace
the `--armor --yes` flags can be specified after the ACI path.

For example if the user wishes to generate the signature with the following command:

```bash
gpg --no-default-keyring --armor --secret-keyring ./rkt.sec --keyring ./rkt.pub --output mycoolapp.aci.asc --detach-sig mycoolapp.aci
```

the acbuild command to do it would be:

```bash
acbuild write mycoolapp.aci --sign -- --no-default-keyring --armor --secret-keyring ./rkt.sec --keyring ./rkt.pub
```
