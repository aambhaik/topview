Topview is a command-line utility to access Mashery Local toplogy for read-only use.
This utility uses the following Go extension libraries

* [Cobra](https://github.com/spf13/cobra)
* [Viper](https://github.com/spf13/viper)

# Table of Contents
- [How to build](#Howtobuild)
- [Overview](#overview)

# How to build
Simple! Just run the following...
* go get github.com/aambhaik/topview/...

# Overview

Topview is a command-line utility that works with Mashery Local 5.0 registry to provide:

* View all Clusters: `topview get clusters`
* Use one of the clusters: `topview use cluster <cluster-name>`
* View all Zones: `topview get zones`
* Use one of the zones: `topview use zone <zone-name>`
* View all nodes of different types in a zone: `topview get nodes --type=<node-type>`
   Following are valid node types
    - db (For database nodes)
    - nosql (For cassandra nodes)
    - caches (For cache nodes)
    - tm (For Traffic Manager nodes)
    - log (For log service nodes)

* The default registry service host is localhost and the default registry port is 1080. It can be changed
  by setting the following ENV variables
   - TMGC_REGISTRY_HOST
   - TMGC_REGISTRY_PORT
