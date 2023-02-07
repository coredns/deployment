# Systemd Service File

Use `coredns.service` as a systemd service file. It defaults to using a "coredns" user with
a homedir of `/var/lib/coredns` and the binary lives in `/usr/bin` and the config in
`/etc/coredns/Corefile`.

In order to work, you need to do following jobs:

- Put `coredns` binary in `/usr/bin`
- Put `Corefile` at `/etc/coredns/Corefile`
- Put `coredns-sysusers.conf` in `/usr/lib/sysusers.d`
- Put `coredns-tmpfiles.conf` in `/usr/lib/tmpfiles.d`
- Put `coredns-log.conf.conf` in `/etc/logrotate.d`
