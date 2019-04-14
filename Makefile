# Makefile for building packages for CoreDNS.

# ARCH can be and default to amd64 is not set.
ARCH := amd64 armhf arm64

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
