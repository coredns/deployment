## Presentation

CoreDNS is released as DNS Discovery Service for Kubernetes:
* Kubernetes v1.9 - CoreDNS 1.0.0 is an alpha feature
* Kubernetes v1.10 - CoreDNS 1.0.6 is a beta feature
* Kubernetes v1.11 - CoreDNS 1.1.3 is GA feature, and install by default when the cluster is setup with kubeadm
* Kubernetes v1.12 - CoreDNS 1.2.2 is considered the default DNS Service.


This document intents to help people having trouble setting-up CoreDNS, or maintaining CoreDNS,
by sharing knowledge on known issues or questions already asked and resolved we collected
during the online support we provided to people having trouble with the DNS Service after installation or upgrade of their cluster.


## Useful links about CoreDNS for Kubernetes:

* [Using CoreDNS for Service Discovery](https://kubernetes.io/docs/tasks/administer-cluster/coredns/)
* [Debugging DNS Resolution](https://kubernetes.io/docs/tasks/administer-cluster/dns-debugging-resolution/)
* [Troubleshooting kubeadm](https://kubernetes.io/docs/setup/independent/troubleshooting-kubeadm/)
* [DNS Specifications for Kubernetes](https://github.com/kubernetes/dns/blob/master/docs/specification.md)
* [CoreDNS' Kubernetes plugin documentation](https://coredns.io/plugins/kubernetes/)


## Advises for troubleshooting

When starting a Kubernetes cluster with CoreDNS, if the Network is not stable, it usually happen that DNS is the first symptom.
In other words, an issue like "DNS does not work properly" has often its root cause in "Network is not working properly".

**Known network issues of some Kubernetes installations:**

Kubernetes is managing IP Tables. It may not be compatible with other firewall tool provided by the OS.
* Firewalld is running on the node
* Systemd running on the node
* Others: [Fresh deploy with CoreDNS not resolving any dns lookup](https://github.com/kubernetes/kubeadm/issues/1056)


## Known Issues

### Some Prometheus data is missing after reload

in CoreDNS v1.2.2 that is provided with K8s v1.12, after edition of the Confimap,
CoreDNS will automatically "reload" its configuration file.

Issue is raised about missing some metrics on the Prometheus interface.

## Frequently Asked Questions

### FAQ-1 : Using ```busybox``` to test DNS queries
Busybox has a known bug that make it not working properly when used for kubernetes.
workaround: You should ensure that the image of Busybox is lower or equal to 1.28.4
see [dns can't resolve kubernetes.default and/or cluster.local](https://github.com/kubernetes/kubernetes/issues/66924#issuecomment-411804435)
and initial issue on Busybox : [Nslookup does not work in latest busybox image](https://github.com/docker-library/busybox/issues/48)


### FAQ-2 : Using Autopath plugin in CoreDNS configmap
To operate properly, Autopath needs that the option "pods verified" is set in Kubernetes's plugin configuration
Please look into the [documentation for kubernetes plugin, section AutoPath](https://coredns.io/plugins/kubernetes/)

### FAQ-3 : How to test DNS service
Either you run a pod with Busybox, as describe in [Debugging DNS Resolution](https://kubernetes.io/docs/tasks/administer-cluster/dns-debugging-resolution/),
either you run ad DNS tools dedication pod using :
```
kubectl run -it --rm --restart=Never --image=infoblox/dnstools:latest dnstools
```
After the pod is started and you get the console propmt, you can run one of those commands:
```
nslookup kubernetes.default
dig kuberneted.default A
```


### FAQ-4 : How to rollback to kube-dns, when you are using kubeadm tool to install your Cluster
```
kubeadm upgrade apply <current-installed-version> --feature-gates CoreDNS=false --force
```
Warning, also this command is correct, it may not work because of ongoing issues in kubeadm.
>check with Sandeep the opened issues

### FAQ-5 : CoreDNS pod is crashing
Either by Out Of Memory, either by CrashBackoff.
On some installation of cluster on Ubuntu, there is confusion on what should be the /etc/resolv.conf file.
This file MUST not contain a localhost resolver.
If it does, eg 127.0.0.53 or similar, ensure to cleanup your /etc/resolv.conf.

* see: [CoreDNS pod dies when trying to resolve?](https://github.com/coredns/coredns/issues/1986)
* see: [Update CoreDNS to v1.12 to fix OOM & restart](https://github.com/kubernetes/kubeadm/issues/1037)

This issue should not happen on Cluster created with kubeadm with a version 1.11 or higher.
### FAQ-6 : CoreDNS pod periodically hit memory limit when stressed with large domwin queries
Check the version of CoreDNS you are using.
Before CoreDNS 1.2.2, the memory limit of CoreDNS deployent was not in sync with the amount of memory that may be needed for the cache.
if you can, upgrade to CoreDNS 1.2.2, else change the memory limit of CoreDNS deployment to 350mi.

reported issue was : [CoreDNS pods get OOM Killed](https://github.com/kubernetes/kops/issues/5652)

### FAQ-7 : Pod Security Policy
CoreDNS deployment listen, by default, on the port 53. It may not be compatible with a PSP. In that case, 2 ways:
1- extends your Pod Security Policy to include CoreDNS need to bind a system port
2- workaround by switching port in the service definition:
   * changing the listening port of the Pod to a port higher than 1024. (need to update the Configmap and deployment)
   * update kube-dns service definition by adding a switch of port from 53 to the new one defined







