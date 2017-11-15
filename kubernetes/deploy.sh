#!/bin/bash

# Deploys CoreDNS to a cluster currently running Kube-DNS.

SERVICE_CIDR=$1
POD_CIDR=$2
CLUSTER_DOMAIN=${3:-cluster.local}
YAML_TEMPLATE=${4:-`pwd`/coredns.yaml.sed}

if [[ -z $SERVICE_CIDR ]]; then
	echo "Usage: $0 SERVICE-CIDR [ POD-CIDR ] [ CLUSTER-DOMAIN ] [ YAML-TEMPLATE ]"
	exit 1
fi

CLUSTER_DNS_IP=$(kubectl get service --namespace kube-system kube-dns -o jsonpath="{.spec.clusterIP}")

sed -e s/CLUSTER_DNS_IP/$CLUSTER_DNS_IP/g -e s/CLUSTER_DOMAIN/$CLUSTER_DOMAIN/g -e s?SERVICE_CIDR?$SERVICE_CIDR?g -e s?POD_CIDR?$POD_CIDR?g $YAML_TEMPLATE
