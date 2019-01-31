## Introduction

From Kubernetes v1.13, CoreDNS is the recommended default DNS server in Kubernetes. 

As more and more users are deploying and using CoreDNS is their Kubernetes cluster, there are a number of FAQs from users having trouble setting up or maintaining CoreDNS.  
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

Hence, if you're confident that the network is stable and DNS still doesn't work, this FAQ document is applicable. 

## Known Issues

### Some Prometheus data is missing after reload

In CoreDNS v1.2.2 that is provided with K8s v1.12, after editing the ConfigMap, CoreDNS will automatically "reload" its configuration file.
The reload causes some metrics on the Prometheus interface to be missing.

This has been fixed in CoreDNS v1.2.3.

## Frequently Asked Questions

#### FAQ-1 : Why am I unable to get an answer to the DNS queries when I'm using ```busybox``` to test?

`nslookup` in busybox has a [known bug](https://github.com/docker-library/busybox/issues/48) that prevents DNS lookups from working unless you use FQDNs.
The workaround for this is to ensure that the image of Busybox is lower or equal to 1.28.4
See: [dns can't resolve kubernetes.default and/or cluster.local](https://github.com/kubernetes/kubernetes/issues/66924#issuecomment-411804435) for more details.

#### FAQ-2 : Why isn't the `autopath` plugin in the default configuration?

Autopath requires that the option `pods verified` is set in Kubernetes's plugin configuration. The `pods verified` takes up significantly more memory. 
There was an [analysis done for resource requirements](https://github.com/coredns/deployment/blob/master/kubernetes/Scaling_CoreDNS.md) with and without the Autopath plugin to show the memory and CPU utilization.
Also look into the [documentation](https://github.com/coredns/coredns/tree/master/plugin/kubernetes#autopath) for more details about the Autopath plugin.

#### FAQ-3 : How do I test the DNS service?

To test the DNS service, it is important that you are testing DNS resolution from a Pod, not from one of the nodes. There are two definite ways to test this: 

- Run a pod with Busybox, as described in [Debugging DNS Resolution](https://kubernetes.io/docs/tasks/administer-cluster/dns-debugging-resolution/)

- Run a DNS tools dedicated pod using the following command:

```
kubectl run -it --rm --restart=Never --image=infoblox/dnstools:latest dnstools
```

After the pod is started and you get the console prompt, run the following command:

```
# dig kubernetes.default A +search
```

#### FAQ-4 : How do I rollback to kube-dns? I am using kubeadm tool to install the Kubernetes cluster.

The following command will reapply your cluster with existing settings but will replace CoreDNS with kube-dns.

```
kubeadm upgrade apply <current-installed-version> --feature-gates CoreDNS=false --force
```

Warning, although this command is correct, it may not work because of ongoing issues in kubeadm.

#### FAQ-5 : Why is my CoreDNS pod in `CrashLoopBackoff` or `Error` state?

1. When installing a Kubernetes cluster on Ubuntu, when the /etc/resolv.conf file contains a localhost resolver, there are chances that CoreDNS is stuck in a loop. 
This will be detected by the `loop` plugin and will show in the logs.

    - If the logs show that there is a loop detected and the /etc/resolv.conf file contains, eg. 127.0.0.53 or similar, you need to cleanup your /etc/resolv.conf. 
For more details, checkout the [documentation on loop detection](https://github.com/coredns/coredns/tree/master/plugin/loop#troubleshooting)
 
    - This issue should not happen on a Kubernetes cluster created with kubeadm with a version 1.11 or higher.

2. If you have nodes that are running SELinux with an older version of Docker you might experience a scenario where the coredns pods are not starting. To solve that you can try one of the following options:
   
   - Upgrade to a newer version of Docker.
   - [Disable SELinux](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/6/html/security-enhanced_linux/sect-security-enhanced_linux-enabling_and_disabling_selinux-disabling_selinux).
   - Modify the coredns deployment to set `allowPrivilegeEscalation` to true.
 