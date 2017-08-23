## _This project is currently unmaintained_

acbuild was originally created to be the tool used to build AppC images. Due to
the introduction of the [Open Container
Initiative](https://www.opencontainers.org/), development on AppC was
[officially suspended](https://github.com/appc/spec#-disclaimer-) in November,
2016. While acbuild has the ability to also produce OCI images, it is not the
only tool capable of doing so. In its current state, acbuild is not currently
maintained. If you wish to become a maintainer of acbuild, feel free to start
contributing and ask for direct commit access via the issue/PR tracker.

For those looking for an OCI image manipulation tool that is actively
maintained, [umoci](https://github.com/openSUSE/umoci) or
[buildah](https://github.com/projectatomic/buildah) might be able to fill the
role.

# acbuild - another container build tool

acbuild is a command line utility to build and modify container images.

It is intended to provide an image build workflow independent of specific
formats; currently, it can output the following types of container images:
* ACI, the container image format defined in the [App Container (appc) 
 spec](https://github.com/appc/spec).
* OCI, the format defined in the [Open Containers Image Format
  specification](https://github.com/opencontainers/image-spec)


[http://www.youtube.com/watch?v=WcnIDm80y68](http://www.youtube.com/watch?v=WcnIDm80y68)

[![How to use rkt acbuild to construct app container images](http://img.youtube.com/vi/WcnIDm80y68/0.jpg)](http://www.youtube.com/watch?v=WcnIDm80y68 "How to use rkt acbuild to construct app container images")

## Rationale

We needed a powerful tool for constructing and manipulating container images
that made it easy to iteratively build containers, both from scratch and atop
existing images. We wanted that tool to integrate well with Unix mechanisms
like the shell and `Makefile`s so it would fit seamlessly into well-known
administrator and developer workflows.

## Installation

### Dependencies

acbuild can only be run on a Linux system, and has only been tested on the
amd64 architecture.

For trying out acbuild on Mac OS X, it's recommended to use Vagrant.
Instructions on how to do this are a little further down in this document.

acbuild requires a handful of commands be available on the system on 
which it's run:

- `cp`
- `modprobe`

Additionally `systemd-nspawn` is required to use the [default
engine](Documentation/subcommands/run.md) for acbuild run. Thus on Ubuntu the `systemd-container` package needs to be installed.

### Prebuilt Binaries

The easiest way to get `acbuild` is to download one of the
[releases](https://github.com/containers/build/releases) from GitHub.

### Build from source

The other way to get `acbuild` is to build it from source. Building from source requires [Go 1.5+](https://golang.org/dl/).

Follow these steps to do so:

1. Grab the source code for `acbuild` by `git clone`ing the source repository:
   ```
   cd ~
   git clone https://github.com/containers/build acbuild
   ```

2. Run the `build` script from the root source repository directory:
   ```
   cd acbuild
   ./build
   ```

   Or, if you want to build in docker (assuming `$PWD` exists and contains
   `acbuild/` on your Docker host):

   ```
   cd acbuild
   ./build-docker
   ```

3. A `bin/` directory will be created that contains the `acbuild` tool. To make
   sure your shell can find this executable, append this directory to your
   environment's `$PATH` variable. You can do this in your `.bashrc` or similar
   file, for example:
   ```
   vi ~/.bashrc
   ```

and put the following lines at the end of the file:
   ```
   export ACBUILD_BIN_DIR=~/acbuild/bin
   export PATH=$PATH:$ACBUILD_BIN_DIR
   ```

### Building acbuild with rkt

If rkt is installed on the system, acbuild can also be built inside of a rkt container with the following command:

```
./build-rkt
```


### Trying out acbuild using Vagrant

For users with Vagrant 1.5.x or greater, there's a provided `Vagrantfile` that
can quickly get you set up with a Linux VM that has both acbuild and rkt. The
following steps will grab acbuild, set up the machine, and ssh into it.

```
git clone https://github.com/containers/build acbuild
cd acbuild
vagrant up
vagrant ssh
```

## Documentation

Documentation about acbuild and many of its commands is available in the
[`Documentation`
directory](https://github.com/containers/build/tree/master/Documentation) in this
repository.

## Examples

Check out the [`examples`
directory](https://github.com/containers/build/tree/master/examples) for some common
applications being packaged into ACIs with `acbuild`.

## Related work

- https://github.com/sgotti/baci
- https://github.com/appc/spec/tree/master/actool - particularly the `build` and
  `patch-manifest` subcommands. `acbuild` may subsume such functionality,
  leaving `actool` as a validator only.
- https://github.com/blablacar/dgr
