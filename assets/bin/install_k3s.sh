#!/bin/bash

k3s_version=`k3s -v |grep k3s |awk '{print $3}'`
if [[ "$k3s_version" = "v1.23.12+k3s1"]]; then
  echo install
fi

# curl -sfL https://get.k3s.io | \ 
#   K3S_TOKEN=tokenfromfile \
#   K3S_URL=https://localhost:6443 \
#   INSTALL_K3S_VERSION=v1.23.12+k3s1 \
#   INSTALL_K3S_EXEC=server \
#   K3S_DATASTORE_ENDPOINT=mysql://mysqlfromfile \ 
#   sh - 