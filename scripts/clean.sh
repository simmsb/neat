#!/bin/bash
# set -e


for bridge in `sudo ovs-vsctl list-br`; do
  sudo ovs-vsctl del-br $bridge
done

sudo virsh list --all | grep -o -E "mtv-\w*-\w*" | xargs -I % sh -c 'sudo virsh destroy --domain %'
sudo virsh list --all | grep -o -E "mtv-\w*-\w*" | xargs -I % sh -c 'sudo virsh undefine --domain %'
