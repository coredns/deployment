# Makefile for building packages for CoreDNS.

# ARCH can be and default to amd64 is not set.
ARCH := amd64 armhf arm64
redhat-packages-dist := $(patsubst %.centos,%,$(shell rpm --eval "%{dist}"))
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
mkfile_dir := $(dir $(mkfile_path))

ifeq ($(ARCH),)
    ARCH:=amd64
endif

.PHONY: debian
debian:
	for a in $(ARCH); do \
	    dpkg-buildpackage -us -uc -b --target-arch $$a ;\
	done

debian-clean:
	rm *.tgz

.PHONY: redhat
redhat:
	rpmbuild --undefine=_disable_source_fetch -ba \
		--verbose $(mkfile_dir)/redhat/SPECS/coredns.spec \
		--define "_topdir $(mkfile_dir)/redhat" --define "dist $(redhat-packages-dist)"

redhat-clean:
	rm -r $(mkfile_dir)/redhat/RPMS/*
	rm -r $(mkfile_dir)/redhat/SRPMS/*
	rm -r $(mkfile_dir)/redhat/BUILD/*
