%w(
    /etc/security/limits.conf
).each do |path|
  template path do
    source "../rootfs#{path}"
    mode '0644'
  end
end