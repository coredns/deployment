# Deployment

Scripts, utilities, and examples for deploying CoreDNS.

## MacOS

The default settings will proxy all requests to hostnames not found in your host file to Google's DNS-over-HTTPS.

To install:
  - Run `brew tap "coredns/deployment" "https://github.com/coredns/deployment"`
  - Run `brew install coredns`
  - Run `sudo brew services start coredns`
  - test with `dig google.com @127.0.0.1` and you should see  `SERVER: 127.0.0.1#53(127.0.0.1)`

Using CoreDNS as your default resolver (e.g. for your `Wi-Fi` interface):
 - Run `networksetup -setdnsservers Wi-Fi 127.0.0.1`
 
# Debian

On a debian system:

  - Run `dpkg-buildpackage -us -uc -b  --target-arch ARCH`
    Where ARCH can be any of the released architectures, like "amd64" or "arm".
  - Most users will just run: `dpkg-buildpackage -us -uc -b`

To install:

  - Run `dpkg -i coredns_0.9.10-0~9.20_amd64.deb`.

This installs the coredns binary in /usr/bin, adds a coredns user (homedir set to /var/lib/coredns)
and a small Corefile /etc/coredns.
