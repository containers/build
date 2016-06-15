# Development Guide

## Dependency Vendoring

acbuild vendors its dependencies. This is done mostly with the [glide][1] tool,
and [glide-vc][2] is used to strip unnecessary files from the dependencies.

These tools can be easily installed with a couple of `go get` commands.

```
go get github.com/Masterminds/glide
go get github.com/sgotti/glide-vc
```

This will fetch both tools, build them, and put the binaries in `$GOPATH/bin`.

### Adding a Dependency

```
glide get -s -u -v <package-name>
glide vc --only-code --no-tests
```

### Updating a Dependency

Edit the `glide.yaml` file in the repository, and then run:

```
glide up -s -u -v
glide vc --only-code --no-tests
```

### Removing a Dependency

```
glide rm --delete <package-name>
```

[1]: https://github.com/Masterminds/glide
[2]: https://github.com/sgotti/glide-vc
