#!/bin/bash

k3s_version=`k3s -v |grep k3s |awk '{print $3}'`

if [[ "$k3s_version" = "{{ .Version }}" ]]; then
  echo install
fi

export K3S_URL={{ .Secret.URL }}
export INSTALL_K3S_VERSION={{ .Version }}
export INSTALL_K3S_EXEC={{ .Node.Type }}
export K3S_TOKEN={{ .Secret.Token }}
{{ if eq .Node.Type "server" }}
export K3S_DATASTORE_ENDPOINT="{{ .Secret.DataStore }}"
{{ end }} 
curl -sfL https://get.k3s.io | sh - 
