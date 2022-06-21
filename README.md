# Deployment

Scripts, utilities, and examples for deploying CoreDNS.


# Debian

On a debian system:

  - Run `dpkg-buildpackage -us -uc -b  --target-arch ARCH`
    Where ARCH can be any of the released architectures, like "amd64" or "arm".
  - Most users will just run: `dpkg-buildpackage -us -uc -b`
  - Note that any existing environment variables will override the default makefile variables in [debian/rules](debian/rules)
  - The above can be used, for example, to build a particular verison by setting the `VERSION` environment variable

To install:

  - Run `dpkg -i coredns_0.9.10-0~9.20_amd64.deb`.

This installs the coredns binary in /usr/bin, adds a coredns user (homedir set to /var/lib/coredns)
and a small Corefile /etc/coredns.

# Kubernetes

## Helm Chart

The repository providing the helm chart repo is available under

https://github.com/coredns/helm
