# CoreDNS version in Kubernetes

CoreDNS has been shipping with Kubernetes since CoreDNS has been an Alpha feature in Version 1.9.

This document helps map the CoreDNS version that has been shipped with Kubernetes since v1.9.
It also includes all the changes that have been made in CoreDNS from the previous release of Kubernetes to the current.


| Kubernetes Version   |      CoreDNS version      |  Changes in CoreDNS from previous release to Kubernetes |
|:--------------------:|:-------------------------:|:----------|
| v1.14                |  v1.3.1                   | TTL is also applied to negative responses (NXDOMAIN, etc). <br> <br> k8s_external a new plugin that allows external zones to point to Kubernetes in-cluster services. <br><br>kubernetes now checks if a zone transfer is allowed. Also allow a TTL of 0 to avoid caching in the cache plugin. |
| v1.13                |  v1.2.6                   | ***Huge optimization which leads to less memory usage (~30% less)***.<br><br> Support for using a kubeconfig file, including various auth providers (Azure not supported due to a compilation issue with that code).<br><br>`loop` plugin fixes a bug when dealing with a failing upstream.  |
| v1.12                |  v1.2.2                   | Makes the default cache size smaller.<br><br> A new plugin called loop was added. When starting up it detects resolver loops and stops the process if one is detected.<br><br>The etcd plugin now supports etcd version 3 (only!). It was impossible to support v2 and v3 at the same time (even as separate plugins); so we decided to drop v2 support.<br><br>The auto plugin now works better with Kubernetes Configmaps.
| v1.11                |  v1.1.3                   | ***Deprecation notice for the reverse plugin.***<br><br>***Deprecation notice for the https_google protocol in proxy.***<br><br>kubernetes has a small fix for apex queries.<br><br>kubernetes adds option to ignore services without ready endpoints.<br><br>cache fixes the critical spoof vulnerability.<br><br>A new plugin was added: reload, which watches for changes in your Corefile and then automatically will reload the process. 
| v1.10                |  v1.0.6                   | ***The startup and shutdown plugin are deprecated (but working and included) in this release in favor of the on plugin. If you use them, this is the moment to move to on.***<br><br>A plugin called `forward` has been included in CoreDNS, this was, up until now, an external plugin.<br><br>We now support zone transfers in the kubernetes plugin.<br><br>***Fixes a vulnerability in the underlying DNS library, CVE-2017-15133.*** <br><br>kubernetes, adds a fix for pod insecure look ups for non-IP addresses.
| v1.09                |  v1.0.1                   | The v1.0.1 was the first version of CoreDNS to be shipped with Kubernetes. |