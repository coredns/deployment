#!/bin/bash

# Roll back kube-dns to the cluster which has CoreDNS installed.

show_help () {
cat << USAGE
usage: $0  [ -i DNS-IP ] [ -d CLUSTER-DOMAIN ] [-m MEMORY-LIMIT] [-v KUBERNETES-VERSION]

    -i : Specify the cluster DNS IP address. If not specified, the IP address of
         the existing "kube-dns" service is used, if present.
    -d : Specify the Cluster Domain. Default is "cluster.local"
    -m : Specify the memory limit for kube-dns. Defaults to the memory limit defined for CoreDNS deployment.
         (Example: 170Mi)
    -v : Specify the Kubernetes Version. If not specified, the latest kube-dns manifest will be applied.
         (Example: 1.18)
USAGE
exit 0
}


# Simple Defaults
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
CLUSTER_DOMAIN=cluster.local
YAML_TEMPLATE="$DIR/kube-dns.yaml.sed"

# Get Opts
while getopts "hi:d:m:v:" opt; do
    case "$opt" in
    h)  show_help
        ;;
    i)  CLUSTER_DNS_IP=$OPTARG
        ;;
    d)  CLUSTER_DOMAIN=$OPTARG
        ;;
    m)  MEMORY_LIMIT=$OPTARG
        ;;
    v)  K8S_VERSION=$OPTARG
        ;;
    esac
done

if [[ -z ${K8S_VERSION} ]]; then
  curl -L  https://raw.githubusercontent.com/kubernetes/kubernetes/master/cluster/addons/dns/kube-dns/kube-dns.yaml.base > "$YAML_TEMPLATE"
else
  curl -L  https://raw.githubusercontent.com/kubernetes/kubernetes/release-${K8S_VERSION}/cluster/addons/dns/kube-dns/kube-dns.yaml.base > "$YAML_TEMPLATE"
fi

if [[ -z ${CLUSTER_DNS_IP} ]]; then
  # Default IP to kube-dns IP
  CLUSTER_DNS_IP=$(kubectl get service --namespace kube-system kube-dns -o jsonpath="{.spec.clusterIP}")
  if [[ $? -ne 0 ]]; then
      >&2 echo "Error! The IP address for DNS service couldn't be determined automatically. Please specify the DNS-IP with the '-i' option."
      exit 2
  fi
fi

if [[ -z ${MEMORY_LIMIT} ]]; then
  # Default Memory Limit to the one specified in the CoreDNS deployment
  MEMORY_LIMIT=$(kubectl get deployment --namespace kube-system coredns -o jsonpath="{.spec.template.spec.containers[0].resources.limits.memory}")
  if [[ $? -ne 0 ]]; then
      >&2 echo "Error! The memory limit for the DNS deployment couldn't be determined automatically. Please specify the memory with the '-m' option."
      exit 2
  fi
fi

sed -e s/__PILLAR__DNS__SERVER__/${CLUSTER_DNS_IP}/g -e s/__PILLAR__DNS__DOMAIN__/${CLUSTER_DOMAIN}/g -e s/__PILLAR__DNS__MEMORY__LIMIT__/${MEMORY_LIMIT}/g ${YAML_TEMPLATE}
