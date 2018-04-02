# Systemd Service File

Use `coredns.service` as a systemd service file. It defaults to using a "coredns" user with
a homedir of `/var/lib/coredns` and the binary lives in `/usr/bin` and the config in
`/etc/coredns/Corefile`.

In order to work, the systemd unit needs a user named `coredns`, an handy way to provide
it is to use the `systemd-sysusers` service by installing the `coredns-sysusers.conf` file in the
`sysusers.d` folder (e.g: `/usr/lib/sysusers.d`).
