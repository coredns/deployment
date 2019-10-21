#!/bin/bash

# Roll back kube-dns to the cluster which has CoreDNS installed.

show_help () {
cat << USAGE
usage: $0  [ -i DNS-IP ] [ -d CLUSTER-DOMAIN ]

    -i : Specify the cluster DNS IP address. If not specified, the IP address of
         the existing "kube-dns" service is used, if present.
    -d : Specify the Cluster Domain. Default is "cluster.local"
USAGE
exit 0
}


# Simple Defaults
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
CLUSTER_DOMAIN=cluster.local
YAML_TEMPLATE="$DIR/kube-dns.yaml.sed"

curl -L  https://raw.githubusercontent.com/kubernetes/kubernetes/master/cluster/addons/dns/kube-dns/kube-dns.yaml.base > "$YAML_TEMPLATE"

# Get Opts
while getopts "hi:d:" opt; do
    case "$opt" in
    h)  show_help
        ;;
    i)  CLUSTER_DNS_IP=$OPTARG
        ;;
    d)  CLUSTER_DOMAIN=$OPTARG
        ;;
    esac
done


if [[ -z $CLUSTER_DNS_IP ]]; then
  # Default IP to kube-dns IP
  CLUSTER_DNS_IP=$(kubectl get service --namespace kube-system kube-dns -o jsonpath="{.spec.clusterIP}")
  if [ $? -ne 0 ]; then
      >&2 echo "Error! The IP address for DNS service couldn't be determined automatically. Please specify the DNS-IP with the '-i' option."
      exit 2
  fi
fi

sed -e s/__PILLAR__DNS__SERVER__/${CLUSTER_DNS_IP}/g -e s/__PILLAR__DNS__DOMAIN__/${CLUSTER_DOMAIN}/g ${YAML_TEMPLATE}
