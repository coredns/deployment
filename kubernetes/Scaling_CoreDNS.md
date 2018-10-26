# Scaling CoreDNS in Kubernetes Clusters


The following is a guide for tuning CoreDNS in Kubernetes clusters.


## Memory and Pods

In large scale Kubernetes clusters, CoreDNS's memory usage is predominantly affected by the number of Pods and Services in the cluster. Other factors include the size of the filled DNS answer cache, and the rate of queries received (QPS) per CoreDNS instance.

### With default CoreDNS settings

To estimate the amount of memory required for a CoreDNS instance (using default settings), you can use the following formula:

>  MB required (default settings) = (Pods + Services) * 0.000858 + 49.1

This formula has the following baked in:

* 30 MB for cache. The default cache size is 10000 entries, which uses about 30 MB when completely filled.
* 5 MB operating buffer for processing queries.  In tests, this was the amount used by a single CoreDNS replica under ~30K QPS of load.


![CoreDNS in Kubernetes Memory Use](https://docs.google.com/spreadsheets/d/e/2PACX-1vS7d2MlgN1gMrrOHXa7Zn6S3VqujST5L-4PHX7jr4IUhVcTi0guXVRCgtIYrtLm3qxZWFlMHT-Xt9n3/pubchart?oid=191775389&format=image)

### With the *autopath* plugin

The *autopath* is an optional optimization that improves performance for queries of names external to the cluster (e.g. `infoblox.com`). However, enabling the *autopath* plugin requires CoreDNS to use significantly more memory to store information about Pods. 

To estimate the amount of memory required for a CoreDNS instance (using the *autopath* plugin), you can use the following formula:

>  MB required (w/ autopath) = (Pods + Services) * 0.0025 + 50.5

This formula has the following baked in:

* 30 MB for cache. The default cache size is 10000 entries, which uses about 30 MB when completely filled.
* 5 MB operating buffer for processing queries.  In tests, this was the amount used by a single CoreDNS replica under ~30K QPS of load.

### Configuring your deployment

You can use the formulas above to estimate the required amount of memory needed by CoreDNS in your cluster, and adjust adjust the resource memory request/limit in the CoreDNS deployment accordingly.

## CPU and QPS

Max QPS was tested by using the `kubernetes/perf-tests/dns` tool, on a cluster using CoreDNS. Two types of queries used were *internal queries* (e.g. `kubernetes`), and *external queries* (e.g. `infoblox.com`).  Steps were taken to synthesize the standard Pod behavior of following the ClusterFirst domain search list (plus one local domain). The effective domain search list used in these tests was `default.svc.cluster.local svc.cluster.local cluster.local mydomain.com`).  QPS and latency here are of the client perspective.  This is of special significance for external queries, wherein a single client query actually generates 5 backend queries to the DNS server, is counted only as one query.

### With default CoreDNS settings

Single instance of CoreDNS (default settings) on a GCE n1-standard-2 node:


| Query Type  | QPS              | Avg Latency (ms)   | 
|-------------|------------------|--------------------|
| external    | 6733<sup>1</sup> | 12.02<sup>1</sup>  |
| internal    | 33669            | 2.608              |  


<sup>1</sup> From the server perspective it is processing 33667 QPS with 2.404 ms latency.

### With the *autopath* plugin

The *autopath* plugin in CoreDNS is an option that mitigates the ClusterFirst search list penalty.  When enabled, it answers Pods in one round trip rather than five. This reduces the number of DNS queries on the backend to one.  Recall that enabling the *autopath* plugin requires CoreDNS to use significantly more memory (see prior section). 

Single instance of CoreDNS (with the *autopath* plugin enabled) on a GCE n1-standard-2 node:


| Query Type  | QPS   | Avg Latency (ms) | 
|-------------|-------|------------------|
| external    | 31428 | 2.605            | 
| internal    | 33918 | 2.62             |  


Note that the resulting max QPS for external queries is much higher.  This is due to the *autopath* plugin optimization.

### Configuring your deployment

All QPS loads above pegged both vCPUs (at ~1900m), implying that the max QPS was CPU bound. You may use the results above as reference points for estimating the number of CoreDNS instances to process a given QPS load.  Note that CoreDNS is also multi threaded, meaning that it should be able to process more QPS on a node with more cores. In these tests CoreDNS ran on nodes with two cores.  

## Source Data

The formulas above were based on data collected from tests using the following setups:

**For Memory usage based on the number of Pods and Services:**

The following Kubernetes end-to-end tests (e2e):

* pull-kubernetes-e2e-gce-big-performance: A 500 node cluster with 15000 Pods, 820 Services
* pull-kubernetes-e2e-gce-large-performance: A 2000 node cluster with 60000 Pods, 3280 Services
* pull-kubernetes-e2e-gce-scale-performance: A 5000 node cluster with 150000 Pods, 8200 Services

The memory formulas above are based on a best fit linear trend of these data points.


**For QPS load testing:**

* Kubernetes Cluster: A one master plus 3 node cluster set up using `kube-up` on its default settings.
   * Master: n1-standard-1 (1 vCPU, 3.75 GB memory)
   * Nodes: n1-standard-2 (2 vCPUs, 7.5 GB memory)
   * Networking: kubenet
* Test Client: Tests were executed using the scripts in `kubernetes/perf-tests/dns`, which use `dnsperf` under the hood.
