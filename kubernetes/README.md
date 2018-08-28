
TEST

> **Notice**: We are conducting a [survey](https://www.surveymonkey.com/r/SKZQSLK) to evaluate the adoption of CoreDNS as the DNS for Kubernetes's cluster. 
> If you are in such configuration, **please help us by [providing your feedback](https://www.surveymonkey.com/r/SKZQSLK)**
>
> Thank you, we appreciate your collaboration here.

# Kubernetes

## Description
CoreDNS can run in place of the standard Kube-DNS in Kubernetes. Using the *kubernetes*
plugin, CoreDNS will read zone data from a Kubernetes cluster. It implements the
spec defined for Kubernetes DNS-Based service discovery:

   https://github.com/kubernetes/dns/blob/master/docs/specification.md


## deploy.sh and coredns.yaml.sed

`deploy.sh` is a convenience script to generate a manifest for running CoreDNS on a cluster
that is currently running standard kube-dns. Using the `coredns.yaml.sed` file as a template,
it creates a ConfigMap and a CoreDNS deployment, then updates the Kube-DNS service selector
to use the CoreDNS deployment. By re-using the existing service, there is no disruption in
servicing requests.

By default, the deployment script also translates the existing kube-dns configuration into the equivalent CoreDNS Corefile.
By providing the `-s` option, the deployment script will skip the translation of the ConfigMap from kube-dns to CoreDNS.

The script doesn't delete the kube-dns deployment or replication controller - you'll have to
do that manually, after deploying CoreDNS.

You should examine the manifest carefully and make sure it is correct for your particular
cluster. Depending on how you have built your cluster and the version you are running,
some modifications to the manifest may be needed.

In the best case scenario, all that's needed to replace Kube-DNS are these commands:

~~~
$ ./deploy.sh | kubectl apply -f -
$ kubectl delete --namespace=kube-system deployment kube-dns
~~~

**NOTE:** You will need to delete the kube-dns deployment (as above) since while CoreDNS and kube-dns are running at the same time, queries may randomly hit either one.

For non-RBAC deployments, you'll need to edit the resulting yaml before applying it:
1. Remove the line `serviceAccountName: coredns` from the `Deployment` section.
2. Remove the `ServiceAccount`, `ClusterRole`, and `ClusterRoleBinding` sections.


## Rollback to kube-dns

In case one wants to revert a Kubernetes cluster running CoreDNS back to kube-dns,
the `rollback.sh` script generates the kube-dns manifest to install kube-dns.
This uses the existing service, there is no disruption in servicing requests.

The script doesn't delete the CoreDNS deployment or replication controller - you'll have to
do that manually, after deploying kube-dns.

These commands will deploy kube-dns replacing CoreDNS:
~~~
$ ./rollback.sh | kubectl apply -f -
$ kubectl delete --namespace=kube-system deployment coredns
~~~

**NOTE:** You will need to delete the CoreDNS deployment (as above) since while CoreDNS and kube-dns are running at the same time, queries may randomly hit either one.
