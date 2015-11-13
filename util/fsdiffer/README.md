## FSDiffer ##

A simple filesystem differ library that will report added/modified/deleted files.

Pluggable FSdiffers can be used (they just need to implement the FSDiffer interface that is composed by only the Diff() function)

At the moment a simple fs differ is provided.
In future additional fs differs will be available (for example an overlayfs differ).



