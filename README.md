# Deployment

Scripts, utilities, and examples for deploying CoreDNS.


# Debian

On a debian system:

  - Run `dpkg-buildpackage -us -uc -b  --target-arch ARCH`
    Where ARCH can be any of the released architectures, like "amd64" or "arm".
  - Most users will just run: `dpkg-buildpackage -us -uc -b`

To install:

  - Run `dpkg -i coredns_0.9.10-0~9.20_amd64.deb`.

This installs the coredns binary in /usr/bin, adds a coredns user (homedir set to /var/lib/coredns)
and a small Corefile /etc/coredns.

# Kuebernetes

## Helm Chart

This repository provides helm chart repo. 

```
helm repo add coredns https://raw.github.com/coredns/deployment/gh-pages/
helm install coredns/coredns
```
