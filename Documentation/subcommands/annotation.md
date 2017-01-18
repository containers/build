# acbuild annotation

Annotations are elements in an image's manifest that store extra metadata about
the image. Each annotation has two parts: a key and a value. Each annotation
key is unique for a given manifest. Annotations may be read by external tooling
(like a registry) to get additional information about an image.

## Subcommands

* `acbuild annotation add NAME VALUE`

  Updates the image to contain an annotation with the given name and value. If the
  annotation already exists, its value will be changed.

* `acbuild annotation remove NAME`

  Removes the annotation with the given name from the image.

## Common annotations

Common annotations include:

- `created`: date on which the image was built.
- `authors`: contact details of the creators responsible for the image.
- `homepage`: URL to find more information about the image.
- `documentation`: URL to get documentation on the image.

## Examples

```bash
acbuild annotation add documentation https://example.com/docs

acbuild annotation add authors "Carly Container <carly@example.com>, Nat Network <[nat@example.com](mailto:nat@example.com)>"

acbuild annotation remove homepage
```
