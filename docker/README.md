# docker based deployment

## Prerequisites

* Docker 1.12.x or later (*Docker Swarm Mode*)

## Setup

First decide which nodes you are going to run coredns on and set appropriate
labels on your nodes. I use `iface=extern` as labels on nodes with external
facing interfaces and `iface=intern` for internal facing nodes.

```#!bash
$ docker node inspect node1 | jq '.[0].Spec.Labels'
{
  "iface": "extern"
}
```

## Deploy

Connect to a "manager" node:
(*I use `docker-machine` for this*)

```#!bash
$ eval $(docker-machine env node1)
$ docker stack deploy -c dns.yml dns
```

## Verify

Verify your setup works:

```#!bash
$ dig @<node1> google.com IN A +short
```
