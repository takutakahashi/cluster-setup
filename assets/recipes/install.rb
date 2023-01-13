execute "install script" do
    command "bash cluster-setup/bin/install_k3s.sh"
end

execute "reboot if needed" do
  command "sysctl -A |grep fs.inotify.max_user_watches |grep 524288 || reboot"
end
