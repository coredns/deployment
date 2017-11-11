# Makefile for building packages for CoreDNS.

# Build the debian packages
.PHONY: debian
debian:
	dpkg-buildpackage -us -uc -b --target-arch amd64   
	dpkg-buildpackage -us -uc -b --target-arch arm
	dpkg-buildpackage -us -uc -b --target-arch arm64
	# debs are one up
	ls ../*.deb
