# appc Dependencies

Some of the most common questions people new to acbuild ask are related to image
dependencies. How do they work, how are they referenced, what's the difference
between `acbuild begin <image>` and `acbuild dependency add <image>`.

This document aims to provide a detailed explanation of image dependencies in
AppC and how they can be manipulated with acbuild.

## acbuild begin

One way to interact with other images is via `acbuild begin`. Specifics on the
different ways to invoke it are [here][1].

The short version of it is that when you run `acbuild begin` and provide an
image, the build will start out containing everything from the specified image.
It'll have the same manifest, and all the files in the other image will be
present. When `acbuild write` is called, the resulting image will include all of
this information, so the final image will be the same size or larger than the
image you began from.

When this method is used, the original image is copied and modified, and the
original image isn't needed at runtime when you hand the resulting container
image to a container runtime.

In AppC terminology, this isn't even using dependencies at all.

## acbuild dependency add

AppC provides a way for images to point to other images needed at runtime; these
other images are called dependencies. These pointers can be added via the
`acbuild dep add` command. Specifics on using it are [here][2].

When a container runtime goes to run an image, it also needs all of that image's
dependencies, and all of the dependencies' dependencies, and so on. Adding a
dependency to an image simply adds a pointer to a different image, and nothing
else about the dependency (other than potentially things like labels, and an
expected hash) are stored in the image being built.

These pointers take the form of AppC image names. When you give an image with
dependencies to a container runtime, the runtime should first check its local
store for any images with matching names, use those if they exist, and otherwise
perform AppC discovery to find the images on the internet and fetch them. AppC
image names and discovery are described in greater detail [here][3].

acbuild doesn't have a persistent local store, so if acbuild needs to find
dependencies (which happens if you use the `run` command after a `dep add`
command) it immediately falls back to performing AppC discovery to find the
image. This means that local dependencies are not currently supported.

[1]: Documentation/subcommands/begin.md
[2]: Documentation/subcommands/dependency.md
[3]: https://github.com/appc/spec/blob/master/spec/discovery.md
