directory '/etc/rancher/k3s/' do
  mode '0755'
end

%w(
    /etc/security/limits.conf
    /etc/rancher/k3s/config.yaml
).each do |path|
  template path do
    source "../rootfs#{path}"
    mode '0644'
  end
end

execute "reboot if needed" do
  command "ulimit -n |grep 65535 || reboot"
end