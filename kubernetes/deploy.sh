#!/bin/bash

# Deploys CoreDNS to a cluster currently running Kube-DNS.

show_help () {
cat << USAGE
usage: $0 [ -r REVERSE-CIDR ] [ -i DNS-IP ] [ -d CLUSTER-DOMAIN ] [ -t YAML-TEMPLATE ]

    -r : Define a reverse zone for the given CIDR. You may specify this option more
         than once to add multiple reverse zones. If no reverse CIDRs are defined,
         then the default is to handle all reverse zones (i.e. in-addr.arpa and ip6.arpa)
    -i : Specify the cluster DNS IP address. If not specified, the IP address of
         the existing "kube-dns" service is used, if present.
    -s : Skips the translation of kube-dns configmap to the corresponding CoreDNS Corefile configuration.

USAGE
exit 0
}

# Simple Defaults
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
CLUSTER_DOMAIN=cluster.local
YAML_TEMPLATE="$DIR/coredns.yaml.sed"
STUBDOMAINS=""
UPSTREAM=\\/etc\\/resolv\.conf

# Translates the kube-dns ConfigMap to equivalent CoreDNS Configuration.
function translate-kube-dns-configmap {
    kube-dns-upstreamnameserver-to-coredns
    kube-dns-stubdomains-to-coredns
}

function kube-dns-upstreamnameserver-to-coredns {
  up=$(kubectl -n kube-system get configmap kube-dns  -ojsonpath='{.data.upstreamNameservers}' 2> /dev/null | tr -d '[",]')
  if [[ ! -z ${up} ]]; then
    UPSTREAM=${up}
  fi
}

function kube-dns-stubdomains-to-coredns {
  STUBDOMAIN_TEMPLATE='
    SD_DOMAIN:53 {
      errors
      cache 30
      loop
      forward . SD_DESTINATION {
        max_concurrent 1000
      }
    }'

  function dequote {
    str=${1#\"} # delete leading quote
    str=${str%\"} # delete trailing quote
    echo ${str}
  }

  function parse_stub_domains() {
    sd=$1

  # get keys - each key is a domain
  sd_keys=$(echo -n $sd | jq keys[])

  # For each domain ...
  for dom in $sd_keys; do
    dst=$(echo -n $sd | jq '.['$dom'][0]') # get the destination

    dom=$(dequote $dom)
    dst=$(dequote $dst)

    sd_stanza=${STUBDOMAIN_TEMPLATE/SD_DOMAIN/$dom} # replace SD_DOMAIN
    sd_stanza=${sd_stanza/SD_DESTINATION/$dst} # replace SD_DESTINATION
    echo "$sd_stanza"
  done
}

  sd=$(kubectl -n kube-system get configmap kube-dns  -ojsonpath='{.data.stubDomains}' 2> /dev/null)
  STUBDOMAINS=$(parse_stub_domains "$sd")
}


# Get Opts
while getopts "hsr:i:d:t:k:" opt; do
    case "$opt" in
    h)  show_help
        ;;
    s)  SKIP=1
        ;;
    r)  REVERSE_CIDRS="$REVERSE_CIDRS $OPTARG"
        ;;
    i)  CLUSTER_DNS_IP=$OPTARG
        ;;
    d)  CLUSTER_DOMAIN=$OPTARG
        ;;
    t)  YAML_TEMPLATE=$OPTARG
        ;;
    esac
done

# Conditional Defaults
if [[ -z $REVERSE_CIDRS ]]; then
  REVERSE_CIDRS="in-addr.arpa ip6.arpa"
fi
if [[ -z $CLUSTER_DNS_IP ]]; then
  # Default IP to kube-dns IP
  CLUSTER_DNS_IP=$(kubectl get service --namespace kube-system kube-dns -o jsonpath="{.spec.clusterIP}")
  if [ $? -ne 0 ]; then
      >&2 echo "Error! The IP address for DNS service couldn't be determined automatically. Please specify the DNS-IP with the '-i' option."
      exit 2
  fi
fi

if [[ "${SKIP}" -ne 1 ]] ; then
    translate-kube-dns-configmap
fi

orig=$'\n'
replace=$'\\\n'
sed -e "s/CLUSTER_DNS_IP/$CLUSTER_DNS_IP/g" \
    -e "s/CLUSTER_DOMAIN/$CLUSTER_DOMAIN/g" \
    -e "s?REVERSE_CIDRS?$REVERSE_CIDRS?g" \
    -e "s@STUBDOMAINS@${STUBDOMAINS//$orig/$replace}@g" \
    -e "s/UPSTREAMNAMESERVER/$UPSTREAM/g" \
    "${YAML_TEMPLATE}"
