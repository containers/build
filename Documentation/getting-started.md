# Getting Started with acbuild

acbuild aims to support a workflow very similar to `Dockerfile`s, but with more
flexibility and adherence to the Unix tools philosophy, and native image output
in the modern ACI format. Generally, the workflow consists of a shell script
(analogous to the `Dockerfile`) driving acbuild to construct an ACI, either
from scratch or with another image as a base. This ACI is then stored in an
_image registry_, and fetched by users to execute the applications stored
within.

The following guide will walk you through building and running a simple ACI
with acbuild. The ACI will contain nginx, and will serve static files over HTTP.

Note that some of the commands shown here must be run as root.

## Making the ACI

### Starting the build

All acbuilds happen within a _context_. This build context is explicitly
created at `acbuild begin`, and deleted at build completion with `acbuild end`.
By default, the context and build are stored in the current working directory.

This means that if you change directories in the middle of a build, acbuild will
forget everything about the current build until you go back to the directory
the build was started in. If that's an issue, you can override this with the
`--work-path` flag, which will use the provided directory to store and access
build information instead of the current directory.

The command to start the build is:

```bash
acbuild begin
```

When `begin` isn't passed any arguments the build starts with an empty ACI.
acbuild can also use an existing ACI as the starting point for the build.  More
information on this feature is in the documentation for `begin`.

### Naming the ACI

All ACIs must be named. When a build is started with an empty ACI acbuild
actually gives your ACI a placeholder name, and it won't let you write out an
ACI until you change it away from that. If you're going to host this ACI
somewhere the name should match the name used during meta discovery to find the
download URL. For this example we'll just be running it off of the local
filesystem, so name it whatever you want.

```bash
acbuild set-name example.com/nginx
```

For more information on meta discovery, check out the [appc
spec](https://github.com/appc/spec/blob/master/spec/discovery.md#meta-discovery).

### Adding a dependency

Having an empty ACI isn't all that useful, so let's add a dependency.

```bash
acbuild dependency add quay.io/coreos/alpine-sh
```

Now when this ACI gets run the image at `quay.io/coreos/alpine-sh` will be
fetched, and used as a base for our image. This means that our ACI only needs
to contain the files that we add or modify on top of alpine.

### Installing nginx

Now that we've got a base image with fancy things like a shell and a package
manager, let's install nginx.

```bash
acbuild run -- apk update
acbuild run -- apk add nginx
```

The first command updates alpine's package manager, and the second command
fetches and installs nginx. The `--` in each command stops acbuild from trying
to parse flags that occur after the `--`. They're not necessary in this
instance, but they would be if one of our commands had any arguments that began
with a `-`.

### Adding a mount point

Since the goal of this is to host a static website, we're going to add a mount
point to the default location that nginx serves files out of. This can be used
when the ACI is run to make a directory on the host available inside of the
container.

```bash
acbuild mount add html /usr/share/nginx/html
```

### Setting the exec statement

The last thing we need to set in our ACI is the executable we want it to run.
Without this command, the ACI couldn't be used without specifying a path to an
executable inside of it at runtime.

```bash
acbuild set-exec -- /usr/sbin/nginx -g "daemon off;"
```

The `--` argument serves the same purpose here as it did when we were using the
`run` command, except it actually matters this time. Without it, acbuild would
think that the `-g` argument was a flag for it, instead of something to pass to
`/usr/sbin/nginx` when the ACI is run.

### Writing out the ACI

With that, our ACI is complete. Before ending the build however, we need to
write out the ACI to a file.

```bash
acbuild write nginx.aci
```

### Ending the build

This command will delete our current build context. Had we not written out an
ACI, everything we've done up to this point would be lost.

```bash
acbuild end
```

## Using the ACI

If you have a directory with some files for nginx to serve...

```bash
mkdir test
echo "Hello, world!" > test/index.html
```

... using the ACI is something as simple as the following command:

```bash
rkt run --insecure-options=image ./nginx.aci --volume html,kind=host,source=/path/to/test --net=host
```

Now point your browser at [`http://localhost`](http://localhost).
