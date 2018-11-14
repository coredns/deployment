## Introduction

From Kubernetes v1.13, CoreDNS is the recommended default DNS server in Kubernetes. 

As more and more users are deploying and using CoreDNS is their Kubernetes cluster, there are a number of FAQs asked by users having trouble setting up CoreDNS, or maintaining CoreDNS. 
This document aims to maintain a collection of the most frequently asked questions by users during the support provided to users having trouble with the DNS Service after installation or upgrade of their Kubernetes cluster.

## Useful links about CoreDNS for Kubernetes:

- [Using CoreDNS for Service Discovery](https://kubernetes.io/docs/tasks/administer-cluster/coredns/)
- [Debugging DNS Resolution](https://kubernetes.io/docs/tasks/administer-cluster/dns-debugging-resolution/)
- [Troubleshooting kubeadm](https://kubernetes.io/docs/setup/independent/troubleshooting-kubeadm/)
- [DNS Specifications for Kubernetes](https://github.com/kubernetes/dns/blob/master/docs/specification.md)
- [CoreDNS Kubernetes plugin documentation](https://coredns.io/plugins/kubernetes/)

## Pre-requisites before troubleshooting

It is important to check whether the Kubernetes Network is stable before starting to troubleshoot DNS. It may appear like DNS isn't functioning when there is a network issue in the Kubernetes cluster. 
In other words, an issue like "DNS does not work properly" has often its root cause as "Network is not working properly".

Hence, iff you're confident that the network is stable and DNS still doesn't work, this FAQ document is applicable. 

**Some of the Known network issues in Kubernetes installations are:**

- Kubernetes is managing IP Tables. It may not be compatible with other firewall tool provided by the OS.
- Firewalld is running on the node
- Systemd running on the node
- Others: [Fresh deploy with CoreDNS not resolving any dns lookup](https://github.com/kubernetes/kubeadm/issues/1056)

## Known Issues

### Some Prometheus data is missing after reload

In CoreDNS v1.2.2 that is provided with K8s v1.12, after editing the ConfigMap,
CoreDNS will automatically "reload" its configuration file.

he reload causes some metrics on the Prometheus interface to be missing.

## Frequently Asked Questions

#### FAQ-1 : Why am I unable to get an answer to the DNS queries when I'm using ```busybox``` to test?

Busybox has a [known bug](https://github.com/docker-library/busybox/issues/48) and therefore does not work while using for testing DNS queries in Kubernetes.
The workaround for this is to ensure that the image of Busybox is lower or equal to 1.28.4
See: [dns can't resolve kubernetes.default and/or cluster.local](https://github.com/kubernetes/kubernetes/issues/66924#issuecomment-411804435) for more details.

#### FAQ-2 : Why isn't the `autopath` plugin in the default configuration?

Autopath requires  that the option `pods verified` is set in Kubernetes's plugin configuration. The `pods verified` takes up significantly more memory. Please look into the [documentation](https://github.com/coredns/coredns/tree/master/plugin/kubernetes#autopath) for more details.

#### FAQ-3 : How do I test the DNS service?

There are two definite ways to test this: 

- Run a pod with Busybox, as described in [Debugging DNS Resolution](https://kubernetes.io/docs/tasks/administer-cluster/dns-debugging-resolution/)
- Run a DNS tools dedication pod using the following command:

```
kubectl run -it --rm --restart=Never --image=infoblox/dnstools:latest dnstools
```

After the pod is started and you get the console propmt, you can run one of those commands:

```
# nslookup kubernetes.default
```
or
```
# dig kubernetes.default A
```

#### FAQ-4 : How do I rollback to kube-dns? I am using kubeadm tool to install the Kubernetes cluster.

The following command will reapply your cluster with existing settings but will replace CoreDNS with kube-dns.

```
kubeadm upgrade apply <current-installed-version> --feature-gates CoreDNS=false --force
```

Warning, although this command is correct, it may not work because of ongoing issues in kubeadm.

#### FAQ-5 : Why is my CoreDNS pod in `CrashLoopBackoff` state?

When installing a Kubernetes cluster on Ubuntu, when the /etc/resolv.conf file contains a localhost resolver, there are chances that CoreDNS is stuck in a loop. 
This will be detected by the `loop` plugin and will show in the logs.

If the logs show that there is a loop detected and the /etc/resolv.conf file contains, eg. 127.0.0.53 or similar, ensure to cleanup your /etc/resolv.conf. 
For more details, checkout the [documentation on loop detection](https://github.com/coredns/coredns/tree/master/plugin/loop#troubleshooting)

This issue should not happen on a Kubernetes cluster created with kubeadm with a version 1.11 or higher.
 