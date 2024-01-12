#!/bin/bash

k3s_version=`k3s -v |grep k3s |awk '{print $3}'`
if [[ "$k3s_version" = "{{ .Version }}" ]]; then
  echo install
fi

 curl -sfL https://get.k3s.io | K3S_TOKEN={{ .Secret.Token }} \
   K3S_URL={{ .Secret.URL }} \
   INSTALL_K3S_VERSION={{ .Version }} \
{{- if eq .Node.Type "server" }}
   INSTALL_K3S_EXEC={{ .Node.Type }} \
   K3S_DATASTORE_ENDPOINT="{{ .Secret.DataStore }}" sh - 
{{- else }}
   INSTALL_K3S_EXEC={{ .Node.Type }} sh -
{{- end }} 
