![topview logo](https://raw.githubusercontent.com/aambhaik/topview/master/logo.gif)

Topview is a command-line utility to access Mashery Local toplogy for read-only use.
This utility uses the following Go extension libraries

* [Cobra](https://github.com/spf13/cobra)
* [Viper](https://github.com/spf13/viper)

# Table of Contents
- [How to build](#how-to-build)
- [Overview](#overview)

# How to build
Simple! Just run the following (if you already have a Golang setup)
* go get github.com/aambhaik/topview/...

OR
- if you don't want to setup Golang locally, follow the Docker way...
    * make a new directory 'topview'
      - mkdir topview
      - cd topview
    * Copy the Dockerfile
      - wget https://raw.githubusercontent.com/aambhaik/topview/master/Dockerfile

    * change the registry Host and Port as necessary in the Dockerfile
    * docker build -t mashery/topview:1.0 .
    * docker run -it mashery/topview:1.0 /bin/bash

You are ready to use the topview!


# Overview

Topview is a command-line utility that works with Mashery Local 5.0 registry to provide:

* View all Clusters: `topview get clusters`
* Use one of the clusters: `topview use cluster <cluster-name>`
* View all Zones: `topview get zones`
* Use one of the zones: `topview use zone <zone-name>`
* View all nodes of different types in a zone: `topview get nodes --type=<node-type>` or get all nodes in a zone: `topview get nodes`

   Following are valid node types
    - db (For database nodes)
    - nosql (For cassandra nodes)
    - caches (For cache nodes)
    - tm (For Traffic Manager nodes)
    - log (For log service nodes)
* View property settings of different types of node: `topview get settings --type=<node-type> --nodeId=<node-id>`

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
