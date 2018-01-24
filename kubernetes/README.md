# Kubernetes

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

The script doesn't delete the kube-dns deployment or replication controller - you'll have to
do that manually, after deploying CoreDNS.

You should examine the manifest carefully and make sure it is correct for your particular
cluster. Depending on how you have built your cluster and the version you are running,
some modifications to the manifest may be needed.

In the best case scenario, all that's needed to replace Kube-DNS are these two commands:

~~~
$ ./deploy.sh | kubectl apply -f -
$ kubectl delete --namespace=kube-system deployment kube-dns
~~~

For non-RBAC deployments, you'll need to edit the resulting yaml before applying it:
1. Remove the line `serviceAccountName: coredns` from the `Deployment` section.
2. Remove the `ServiceAccount`, `ClusterRole`, and `ClusterRoleBinding` sections.

TEST 6
