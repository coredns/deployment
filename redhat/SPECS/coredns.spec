# Specfile for CoreDNS package
#

%define debug_package %{nil}

Name: coredns
Summary: CoreDNS is a DNS server that chains plugins
Version: 1.8.4
Release: 1%{?dist}
License: ASL 2.0
Packager: <sayf-eddine.hammemi@scality.com>
Group: System Environment/Base
URL: https://coredns.io
Source0: https://github.com/coredns/%{name}/releases/download/v%{version}/%{name}_%{version}_linux_amd64.tgz
Source1: coredns.service
Source2: coredns.default
Source3: Corefile
BuildArch: x86_64
BuildRoot: %{_tmppath}/%{name}-root

Requires(pre): shadow-utils
%{?systemd_requires}

%description

CoreDNS is a fast and flexible DNS server. The key word here is flexible: with CoreDNS you are able to do what you want with your DNS data by utilizing plugins. If some functionality is not provided out of the box you can add it by writing a plugin.

%prep
%setup -c
cp %{SOURCE1} .
cp %{SOURCE2} .
cp %{SOURCE3} .
%build
# Nothing to build

%install
rm -rf %{buildroot}
install -D -m 755 coredns %{buildroot}%{_bindir}/coredns
install -D -m 644 %{SOURCE3} %{buildroot}%{_sysconfdir}/coredns/Corefile

install -D -m 644 %{SOURCE1} %{buildroot}%{_unitdir}/coredns.service
install -D -m 644 %{SOURCE2} %{buildroot}%{_sysconfdir}/sysconfig/coredns

%clean
rm -rf %{buildroot}

%pre
getent group coredns >/dev/null || groupadd -r coredns
getent passwd coredns >/dev/null || \
  useradd -r -g coredns -s /sbin/nologin \
          -c "CoreDNS services" coredns
exit 0

%post
%systemd_post coredns.service

%preun
%systemd_preun coredns.service

%postun
%systemd_postun coredns.service

%files
%defattr(-,root,root,-)
%{_unitdir}/coredns.service
%{_bindir}/coredns
%config(noreplace) %{_sysconfdir}/coredns/Corefile
%config(noreplace) %{_sysconfdir}/sysconfig/coredns

%changelog
* Sun Jul 4 2021 SayfEddine HAMMEMI <sayf-eddine.hammemi@scality.com> 1.8.4
- Package coredns 1.8.4
