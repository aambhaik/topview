![topview logo](https://raw.githubusercontent.com/aambhaik/topview/master/logo.gif)

Topview is a command-line utility to access Mashery Local toplogy for read-only use.
This utility uses the following Go extension libraries

* [Cobra](https://github.com/spf13/cobra)
* [Viper](https://github.com/spf13/viper)

# Table of Contents
- [Overview](#overview)
- [How to build](#how-to-build)
- [How to use](#how-to-use)

# Overview

Topview is a command-line utility that works with Mashery Local 5.0 registry to provide:

* Get all Clusters: `topview list clusters`
* Use one of the clusters: `topview use cluster <cluster-name>`
* Get all Zones: `topview list zones`
* Use one of the zones: `topview use zone <zone-name>`
* Get all nodes of different types in a zone: `topview list nodes --type=<node-type>` or get all nodes in a zone: `topview list nodes`

   Following are valid node types
    - db (For database nodes)
    - nosql (For cassandra nodes)
    - caches (For cache nodes)
    - tm (For Traffic Manager nodes)
    - log (For log service nodes)
* Get property settings of different types of node: `topview list settings --type=<node-type> --nodeId=<node-id>`

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
    * docker run -it -e TMGC_REGISTRY_HOST=10.1.10.138 mashery/topview:1.0 /bin/bash
      - replace the env var TMGC_REGISTRY_HOST with your registry ip

You are ready to use the topview!

# How to use
 1. Get all clusters `topview list clusters`
 2. Select one of the clusters `topview use cluster <cluster-name>`
 3. Get all Zones in the selected cluster `topview list zones`
 4. Select one of the zones in the cluster `topview use zone <zone-name>`
 5. Get nodes in the selected zone `topview list nodes --type=<node-type>`
 6. Get property settings of different types of node `topview list settings --type=<node-type> --nodeId=<node-id>`

