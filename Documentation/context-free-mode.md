# Context Free Mode

If a tiny change is to be made to an image, for example if the name needs to be
changed, it can be cumbersome to call `begin`, `write`, and `end` for the
single change.

To make this use case more streamlined, the `--modify-*` flags exists. When a
command is invoked with this flag acbuild will create a directory in `/tmp` to
store the build context, and do the following with this alternate context:

- Call `acbuild begin` with the image passed in via the `--modify` flag.
- Call the provided command.
- Call `acbuild write --overwrite` with the image passed in via the `--modify`
  flag.
- Call `acbuild end`.

If more than one change needs to be made, it will be faster to avoid this flag,
as it will result in unnecessary compressing/uncompressing and copying between
the changes.

The `--modify-appc` flag is for modifying appc images, and the `--modify-oci`
flag is for modifying OCI images.
