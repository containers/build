# acbuild script

The script command is a way to build container images written in a easy-to-understand DSL format and then ran. This allows the ability to build an image using a DSL format rather than a bash script.

Similar to building via command line, each example begins with `begin` and ends with `end`.

## Supported commands / DSL reference

Each command utilized by `acbuild` is compatible with the DSL.

```
  annotation [add|remove]                 Manage annotations
  begin                                   Start a new build
  cat-manifest                            Print the manifest from the current build
  copy-to-dir                             Copy a file or directory into a directory in an ACI
  copy                                    Copy a file or directory into an ACI
  dependency [add|remove]                 Manage dependencies
  end                                     end a current build
  environment [add|remove]                Manage environment variables
  isolator [add|remove]                   Manage isolators
  label [add|remove]                      Manage labels
  mount [add|remove]                      Manage mount points
  port [add|remove]                       Manage ports
  replace-manifest                        Replace the manifest in the current build
  run                                     Run a command in an ACI
  script                                  Runs an acbuild script
  set-event-handler [pre-start|post-stop] Manage event handlers
  set-exec                                Set the exec command
  set-group                               Set the group
  set-name                                Set the image name
  set-user                                Set the user
  set-working-directory                   Set the working directory
  version                                 Get the version of acbuild
  write                                   Write the ACI to a file
```


## Example

An HTTP server example running apache on alpine.

```
begin

set-name example.com/apache

dep add quay.io/coreos/alpine-sh

run -- apk update
run -- apk add apache2
run -- /bin/sh -c "echo 'ServerName localhost' >> /etc/apache2/httpd.conf"

port add http tcp 80

mount add html /var/www/localhost/htdocs

set-exec -- /bin/sh -c "chmod 755 / && /usr/sbin/httpd -D FOREGROUND"

write --overwrite apache-latest-linux-amd64.aci

end
```
