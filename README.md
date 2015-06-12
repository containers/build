# acbuild

acbuild is a command line utility to build and modify App Container images (ACIs).

## Rationale

Dockerfiles are powerful and feature useful concepts such as build layers, controlled build environment. At the same time, they lack flexibility (impossible to extend, re-use environment variables) and don't play nicely with the app spec and Linux toolchain (e.g. shell, makefiles)

This proposal introduces the concept of a command-line build tool, `acbuild`, that natively supports ACI builds and integrates well with shell, makefiles and other Unix tools.

## Commands

`acbuild` will support several commands:

* acbuild init in.aci
  creates an empty aci image `in.aci` with manifest filled up with auto generated stub contents

* acbuild env -var="a=b” -var=”c=d" -in=in.aci -out=out.aci   
  adds environment variables to an image `in.aci`, outputs the result to `out.aci`

* acbuild set-run /usr/local/bin/etcd -in=in.aci -out=out.aci
  sets a run command for the aci image in.aci, writing output to `out.aci`

* acbuild set-label -a=key:val -a=key:val -in=in.aci -out=out.aci  	
  sets annotation label in manifest of `in.aci` and outputs the results to `out.aci`

* acbuild add-image add.aci -in=in.aci -out=out.aci	          	
  add contents of image `add.aci` to image `in.aci` and outputs the value to `out.aci`

* acbuild add-dir /dir -in=in.aci -dir=/dir -out=out.aci                   	
  add contents of directory `/dir` to image `in.aci` and outputs the result to `out.aci`

* acbuild exec -in=in.aci -cmd=/var/run/cmd.run -out=out.aci          
  unpacks image.aci, ask systemd-nspawn (either vendored with acbuild or provided by host OS) to execute command in image.aci's environment: `/var/run/cmd.run`, and add the results to the `out.aci` as a separate layer.

* acbuild squash -in=in.aci -layers=* -out=out.aci
  squashes all layers in in.aci and outputs `out.aci` as a result

* acbuild push -user= -pass= in=in.aci url=registry-url -tag a:b
  pushes `in.aci` to the registry and add some tags to it

### acbuild exec

`acbuild exec` executes the command using systemd-nspawn with the root filesystem of the image passed as a parameter.

    acbuild exec -in=dbus.aci -out=built.aci “cd /build && ./configure && make && make install”

starts a build in the filesystem of the image `dbus.aci`

#### exec: modes of operation

The following modes of operation are possible

- un-layered build with overlayfs support
- un-layered build without overlayfs support
- layered build with overlayfs support

In un-layered mode and without overlayfs support `acbuild exec` works as follows:

- unpack `in.aci` to directory `.acbuild/run/process-id()-sha512-short-hash(in.aci)`
- start systemd-nspawn running command 
- in case of successful execution convert the contents of a build directory to `out.aci`

In un-layered mode, and with overlayfs support `acbuild exec` works as follows:

- unpack `in.aci` to `.acbuild/cas/sha512-short-hash(in.aci)`
- mount it as a lower dir using overlayfs 
- mount a new directory as an overlayfs on top of it in `.acbuild/run/process-id()-sha512-short-hash-upper(in.aci)`
- start systemd-nspawn running command, setting root directory as upper dir
- in case of successful execution, take the results of upperdir and package it into out.aci

In layered mode and with overlayfs support `acbuild exec` works as follows:

- unpack `in.aci` to `.acbuild/cas/sha512-short-hash(in.aci)`
- mount it as a lower dir using overlayfs 
- mount a new directory as an overlayfs on top of it in `.acbuild/run/process-id()-sha512-short-hash-upper(in.aci)`
- start systemd-nspawn running command passed by user setting root directory in upper dir
- in case of successful execution, take the results of the workdir and convert it to an image, add this image as a dependency to aci, thus forming a layer, this mode is explicitly activated by `acbuild exec --layer`

#### exec: caching

Caching can be available as an explicit flag for the `acbuild exec`, giving users a choice to re-use the previous execution results for a command in cases when it makes sense, e.g when command execution results are idempotent.

    acbuild exec -cache=true -in=in.aci “git clone --branch v219 --depth 1 git://anongit.freedesktop.org/systemd/systemd /tmp/out”

in case if `-cache=true` is set acbuild executes the following sequence:

- check first if there’s an image in `.acbuild/cache/hash(command line)` and if it is present, reuse it instead of executing it and consider the operation completed
- otherwise, unpack in.aci to some directory `.acbuild/cache/hash(in.aci)`
- mount cas/in directory as a lower dir using overlayfs
- mount a new directory as an overlayfs on top of it
- start systemd-nspawn running command passed by user setting root directory as a cas/dir
- in case of successful execution, take the results of the workdir and convert it to an image
- associate this command `git clone --branch v219 --depth 1 git://anongit.freedesktop.org/systemd/systemd /tmp/out` with the newly created image in `.acbuild/cache/hash(command line)`

Note that in some cases caching does not make sense, e.g. for command  `rm -rf *` would not do anything useful. We would leave the user to make this choice explicitly when writing a build script.

## Modes of operation

acbuild should support several explicit modes of operation that can be selected by user:

- Context-free: `acbuild`
- Context via file or environment variable: `acbuild -context`
- In-place updates: `acbuild --patch`

### Context-free mode

Context-free mode is useful when taking some base image used as a start of the build process, and producing a modified and customized version of it, e.g.

    acbuild add-dir /my-python-app -in=python-base.aci -out=my-app.aci

In context-free mode, input image and output image should be supplied as explicit command line flags: `-in=in.aci -out=out.aci`

### Context via file

Context-dependent build context mode too, that will deduct  `-in` and `-out` flags from the state explicitly initiated by user, e.g.:

    acbuild -c init image.aci -from=python-base.aci

This command will execute the following steps:

- create an image copying it from `python-base.aci`
- create a .acbuild/context.json file with 

```
    { 
      "type": "acbuild-context",
      "context": {
         "build-image": "image.aci",
       }
    }
```

and all the subsequent calls of the `acbuild -c` will re-use the parameters from the context, simulating a Docker-style build.

### In-place updates

In-place updates can be useful when some aci should be modified on the fly e.g.

    acbuild -p set-env HOST=$(hostname) -in image.aci

In place updates are activated by passing `-p` flag to the acbuild tool, in this case it will accept `-in` flag assuming the output to the same image


## Implementation details

acbuild can be a simpler version of [rkt](https://github.com/coreos/rkt) - it will lack systemd and will vendor stage1.aci with patched systemd-nspawn (if <220) or re-use nspawn if the host OS provides it. In fact, rkt's build system can be migrated to acbuild.

## Examples

build rkt stage1 using acbuild and buildroot

    acbuild init image.aci
    acbuild add buildroot.aci
    acbuild add systemd-buildpack.aci   
    acbuild exec "/configure && make && make strip-install" -out stage1.aci


build mongodb from official images

    acbuild init mongodb.aci
    acbuild add -dir mongodb-blabla.bin/ > out

use apt-get to install nginx


    acbuild -in=in.aci -out=int.aci add aptitude.aci 
    acbuild -in=int.aci -out=out.aci exec apt-get -y install nginx




