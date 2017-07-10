# -*- mode: ruby -*-
# vi: set ft=ruby :

$success_message = <<MESSAGE
You must manually enable the shared folder in the Vagrantfile.
Edit the 'config.vm.synced_folder' line and set 'disabled' to 'false',
then run 'vagrant reload --provision' to reboot and mount the shared folder.
We do apologize for the inconvenience. Upstream changes are pending
to avoid this manual bootstrapping process in the future.
MESSAGE

$script = <<SCRIPT
echo "Configuring shared folder under Ubuntu 16.04..."
sudo apt-get --no-install-recommends install -y virtualbox-guest-utils
if [[ ! -d /vagrant ]]; then
echo "#{$success_message}" 1>&2 &&
exit 1
fi
SCRIPT

Vagrant.configure('2') do |config|
    # grab Ubuntu 16.04 official image
    config.vm.box = "ubuntu/xenial64" # Ubuntu 16.04
    config.vm.synced_folder "./", "/vagrant", disabled: true

    # fix issues with slow dns http://serverfault.com/a/595010
    config.vm.provider :virtualbox do |vb, override|
        vb.customize ["modifyvm", :id, "--natdnshostresolver1", "on"]
        vb.customize ["modifyvm", :id, "--natdnsproxy1", "on"]
        # add more ram, the default isn't enough for the build
        vb.customize ["modifyvm", :id, "--memory", "768"]
    end

    # The /vagrant shared folder functionality is broken in Ubuntu 16.04, see:
    # https://bugs.launchpad.net/cloud-images/+bug/1565985
    # So let's handle the automagic ourselves.
    config.vm.provision "shell", inline: $script

    # install Build Dependencies (GOLANG)
    config.vm.provision :shell, :privileged => false, :path => "scripts/vagrant-install-go.sh"

    # Install acbuild
    config.vm.provision :shell, :privileged => false, :path => "scripts/vagrant-install-acbuild.sh"
end
