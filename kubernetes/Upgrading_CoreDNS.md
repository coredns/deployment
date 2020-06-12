
# Handling CoreDNS Upgrades in Kubernetes Clusters

When upgrading a Kubernetes cluster, some care should be taken when upgrading the CoreDNS component to avoid backward incompatible configuration failures. At a high level, you'll want to review the CoreDNS release notes to see if the new version introduces backward incompatibilities or deprecations and adjust your Corefile as needed before upgrading CoreDNS.


## Identifying Backward Incompatibilities and Deprecations

To identify possible backward incompatibilities, you'll need to review the CoreDNS release notes.  CoreDNS Release Notes are located in the [CoreDNS blog](https://coredns.io/blog/). The CoreDNS deprecation policy is such that we will only introduce backward incompatibilities in x.x.0 and x.x.1 releases.  So, for example, if you are upgrading from 1.1.5, to 1.3.1, you should check release notes for 1.2.0, 1.2.1, 1.3.0, and 1.3.1 for any deprecation/backward incompatibility notices.

If you’ve discovered any backward incompatibility notices, you should review your Corefile to see if you are affected.


## Upgrading CoreDNS Manually

1. Identify and resolve any configurations in your Corefile that are backward incompatible with the new CoreDNS version before upgrading.
2. If you want to reduce the risk of CoreDNS downtime, ensure that your CoreDNS Deployment is set up to handle rolling updates. It should have > 1 replica, a RollingUpdate strategy, and a readiness probe defined that points to the `health` plugin in CoreDNS. e.g. …
```
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
```
3. Update the version of CoreDNS in the Deployment.  This step can be done at the same time as step 2.
4. Once the Deployment is updated with steps 2 and 3, monitor the Pod status
   * `kubectl -n kube-system get pods`
5. If you see a Pod Crashing, view the logs to see any errors, then adjust configs or rollback if necessary.
   * `kubectl -n kube-system logs` 


## Upgrading CoreDNS with Kubeadm

Kubeadm will update CoreDNS for you as part of a Kubernetes cluster upgrade.  In doing so, it will replace/reset any custom changes to your CoreDNS Deployment.  For example, if you have increased the number of replicas in your deployment, after an upgrade it will be reset back to the default (2). Kubeadm will not however change your CoreDNS Configmap.  If your CoreDNS Corefile contains any backward incompatible configurations, you’ll want to fix them manually before updating.


## Walkthrough - Manual Update of CoreDNS

In this walkthrough, the initial version is CoreDNS 1.0.6, which supports the `startup` plugin, and I am upgrading to 1.1.0, which no longer supports the plugin.

First let’s look at the release notes for 1.1.0 (https://coredns.io/2018/03/12/coredns-1.1.0-release/). Their notes include the following:

> * The plugins `shutdown` and `startup` where marked deprecated in 1.0.6. This release removes them. You should use `on` instead.

If we look at the coredns Configmap, we can see that I am using one of the deprecated/removed plugins, `startup`.  Here, startup executes an `echo` at startup (albeit frivolously for the sake of example).


```
CTOs-MBP:~ cohaver$ kubectl get configmap coredns -n kube-system -o yaml

apiVersion: v1
data:
  Corefile: |
    .:53 {
        startup echo
        errors
        health
        kubernetes cluster.local in-addr.arpa ip6.arpa {
           pods insecure
           upstream
           fallthrough in-addr.arpa ip6.arpa
        }
        prometheus :9153
        proxy . /etc/resolv.conf
        cache 30
        loadbalance
    }
kind: ConfigMap
metadata:
  name: coredns
  namespace: kube-system
```

If I upgrade to 1.1.0, the Corefile above will cause CoreDNS to exit with an error.  I should update the Corefile now to use the newer plugin that replaced `startup`, but for the sake of example, we will leave it unchanged, so we can see the rolling update failsafe in action.


To prepare fo the rolling update, let’s look at the current coredns Deployment, 

```
CTOs-MBP:~ cohaver$ kubectl get deployment coredns -n kube-system -o yaml
```

This produces > 100 line of output.  I’ll just go over the relevant parts…

This is the CoreDNS deployment …

```
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: coredns
  namespace: kube-system
```

It deploys 2 replicas, and uses rolling updates.

```
spec:
  replicas: 2
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 1
    type: RollingUpdate
```

It’s using CoreDNS version 1.0.6.

```
  template:
    spec:
      containers:
        name: coredns
        image: coredns/coredns:1.0.6
```

It has a liveness probe defined.

```
        livenessProbe:
          failureThreshold: 5
          httpGet:
            path: /health
            port: 8080
            scheme: HTTP
          initialDelaySeconds: 60
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 5
```

… but the Deployment doesn’t have a readiness probe defined, so we will need to define that if rolling updates are to function the way we want.
I’ll add that using the following `kubectl patch` command ...

```
CTOs-MBP:~ cohaver$ kubectl patch deployment coredns -n kube-system -p '{"spec":{"template":{"spec":{"containers":[{"name":"coredns","readinessProbe":{"httpGet":{"path":"/health","port":8080}}}]}}}}'
deployment.extensions/coredns patched
```

Now we should be ready to roll.  We can update the deployment to use coredns 1.0.6 using `kubectl patch` ...

```
CTOs-MBP:~ cohaver$ kubectl patch deployment coredns -n kube-system -p '{"spec":{"template":{"spec":{"containers":[{"name":"coredns", "image":"coredns/coredns:1.1.0"}]}}}}'
deployment.extensions/coredns patched

```

... and immediately, the rolling update begins.  If we look at the status of the coredns Pods, we can see that the old Pods are terminated, except for one, and a new set of pods are started.

```
CTOs-MBP:~ cohaver$ kubectl get pods -n kube-system -l 'k8s-app=kube-dns'

NAME                       READY   STATUS             RESTARTS   AGE
coredns-76cd54b469-jttw7   1/1     Running            0          4h30m
coredns-76cd54b469-v8gxp   0/1     Terminating        0          4h30m
coredns-7d667b54cd-4bqpf   0/1     CrashLoopBackOff   1          12s
coredns-7d667b54cd-bfk6n   0/1     CrashLoopBackOff   1          12s
```

The new Pods are crashing, but one of the original Pods is left running, so the DNS service is not down (although, it is at 1/2 capacity). To see why the new Pods are crashing we can look at the logs...

```
CTOs-MBP:~ cohaver$ kubectl -n kube-system logs coredns-7d667b54cd-4bqpf
2019/02/07 18:50:31 plugin/startup: this plugin has been deprecated
```

As expected, the startup plugin is reported as deprecated, and CoreDNS is exiting.  Now we can go back to edit the Configmap to fix it, and restart the coredns Pods...

```
CTOs-MBP:~ cohaver$ kubectl -n kube-system delete pod coredns-7d667b54cd-4bqpf coredns-7d667b54cd-bfk6n
pod "coredns-7d667b54cd-4bqpf" deleted
pod "coredns-7d667b54cd-bfk6n" deleted

CTOs-MBP:~ cohaver$ kubectl get pods -n kube-system -l 'k8s-app=kube-dns'
NAME                       READY   STATUS    RESTARTS   AGE
coredns-7d667b54cd-9d6cl   1/1     Running   0          18s
coredns-7d667b54cd-mxtrm   1/1     Running   0          18s

```

Once the crashing Pods are deleted, the Deployment spins them back up again with the updated Configmap. Seeing that they are healthy, it deletes the remaining old Pod thus completing the rolling update.
Just validate, we can see that the logs for the new Pods report they are running 1.1.0...

```
CTOs-MBP:~ cohaver$ kubectl -n kube-system logs coredns-7d667b54cd-9d6cl
.:53
2019/02/07 18:58:44 [INFO] CoreDNS-1.1.0
2019/02/07 18:58:44 [INFO] linux/amd64, go1.10, c8d91500
CoreDNS-1.1.0
linux/amd64, go1.10, c8d91500
```


