package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/InVisionApp/tabular"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"
)

var tab tabular.Table

var nodeType string
var nodeId string

var validNodeTypeResources = map[string]string{"db": "/databases", "nosql": "/cassandras", "caches": "/caches", "tm": "/gateways", "log": "/logservices"}

func init() {
	RootCmd.AddCommand(getCmd)
	getCmd.AddCommand(cmdGetClusters)
	getCmd.AddCommand(cmdGetZones)
	getCmd.AddCommand(cmdGetNodes)
	getCmd.AddCommand(cmdGetSettings)

	cmdGetNodes.Flags().StringVar(&nodeType, "type", "", "type of the node")
	cobra.MarkFlagRequired(cmdGetNodes.Flags(), "type")

	cmdGetSettings.Flags().StringVar(&nodeType, "type", "", "type of the node")
	cmdGetSettings.Flags().StringVar(&nodeId, "nodeId", "", "ID of the node")
	cobra.MarkFlagRequired(cmdGetSettings.Flags(), "type")
	cobra.MarkFlagRequired(cmdGetSettings.Flags(), "nodeId")

}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Clusterwide view of the topology",
	Long:  `Clusterwide view of the topology`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please specify one of the following subcommands: [clusters|zones|nodes|settings]")
	},
}

var cmdGetClusters = &cobra.Command{
	Use:   "clusters",
	Short: "Get available clusters",
	Long:  `Get available clusters in the topology`,
	Run: func(cmd *cobra.Command, args []string) {

		dir, err := homedir.Dir()
		if err != nil {
			log.Fatalf("Unable to detect user's HOME directory")
		}

		url := baseRegistryURL + "/clusters"
		tab = tabular.New()
		//tab.Col("id", "Cluster ID", 37)
		tab.Col("name", "Cluster Name", 20)
		tab.Col("created", "Created Time", 20)

		response, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()

		if response != nil {
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}
			clusters, err := GetClusters(body)
			if err != nil {
				log.Fatal(err)
			}

			format := tab.Print("*")
			for _, cluster := range clusters {
				fmt.Printf(format, cluster.Name, cluster.CreatedTime)
			}

			bytes, err := json.Marshal(clusters)
			if err != nil {
				log.Fatalf("Unable to marshal clusters content")
			}
			err = ioutil.WriteFile(dir+"/.topology.clusters.json", bytes, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

var cmdGetZones = &cobra.Command{
	Use:   "zones",
	Short: "Get available zones in a cluster",
	Long:  `Get available zones in a cluster in the topology`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := homedir.Dir()
		bytes, err := ioutil.ReadFile(dir + "/.topviewer-session.json")
		if err != nil {
			fmt.Println("Please select a Cluster by running \n    topview view clusters \nfollowed by \n    topview use cluster <cluster-name>")
			return
		}
		session, err := GetSession(bytes)
		if err != nil {
			log.Fatalf("Unable to un-marshall session json file from user's HOME dir")
		}

		clusterId := session.ClusterId
		url := baseRegistryURL + "/clusters/" + clusterId + "/zones"

		fmt.Printf("Using cluster [%v]\n", session.ClusterName)

		tab = tabular.New()
		//tab.Col("id", "Zone ID", 37)
		tab.Col("name", "Zone Name", 20)
		tab.Col("created", "Created Time", 20)

		response, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()

		if response != nil {
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}
			zones, err := GetZones(body)
			if err != nil {
				log.Fatal(err)
			}

			dir, err := homedir.Dir()
			if err != nil {
				log.Fatalf("Unable to detect user's HOME directory")
			}

			bytes, err := json.Marshal(zones)
			if err != nil {
				log.Fatalf("Unable to marshal zones content")
			}
			err = ioutil.WriteFile(dir+"/.topology.zones.json", bytes, 0644)

			format := tab.Print("*")
			for _, zone := range zones {
				fmt.Printf(format, zone.Name, zone.CreatedTime)
			}
		}
	},
}

var cmdGetNodes = &cobra.Command{
	Use:   "nodes",
	Short: "Get available nodes in a zone",
	Long:  `Get available nodes in a zone in the topology`,
	Run: func(cmd *cobra.Command, args []string) {
		flag.Parse()
		dir, err := homedir.Dir()
		bytes, err := ioutil.ReadFile(dir + "/.topviewer-session.json")
		if err != nil {
			fmt.Println("Please select a Cluster by running \n    topview get clusters \nfollowed by \n    topview use cluster <cluster-name>")
			fmt.Println("Please select a Zone by running \n    topview get zones \nfollowed by \n    topview use zone <zone-name>")
			return
		}
		session, err := GetSession(bytes)
		if err != nil {
			log.Fatalf("Unable to un-marshall session json file from user's HOME dir")
		}

		nodeResource, found := validNodeTypeResources[nodeType]
		if !found {
			fmt.Printf("Invalid node type [%v]\n", nodeType)
			return
		}
		clusterId := session.ClusterId
		zoneId := session.ZoneId
		url := baseRegistryURL + "/clusters/" + clusterId + "/zones/" + zoneId + nodeResource

		fmt.Printf("Using cluster [%v]\n", session.ClusterName)
		fmt.Printf("Using Zone [%v]\n", session.ZoneName)

		tab = tabular.New()
		tab.Col("id", "Node ID", 37)
		tab.Col("type", "Node Type", 10)
		tab.Col("name", "Node Name", 20)
		tab.Col("status", "Node Status", 20)
		tab.Col("host", "Node Host", 20)
		tab.Col("agentPort", "Node Agent Port", 15)
		tab.Col("port", "Node Service Port", 16)

		response, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()

		if response != nil {
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}
			nodes, err := GetNodes(body)
			if err != nil {
				log.Fatal(err)
			}
			if len(nodes) == 0 {
				fmt.Printf("No nodes found for the given node type [%v]\n", nodeType)
				return
			}
			format := tab.Print("*")
			for _, node := range nodes {
				fmt.Printf(format, node.NodeID, nodeType, node.Name, node.Status, node.Host, node.AgentPort, node.Port)
			}
		}
	},
}

var cmdGetSettings = &cobra.Command{
	Use:   "settings",
	Short: "Get settings for a specific node",
	Long:  `Get settings for a specific node in the topology`,
	Run: func(cmd *cobra.Command, args []string) {
		flag.Parse()
		dir, err := homedir.Dir()
		bytes, err := ioutil.ReadFile(dir + "/.topviewer-session.json")
		if err != nil {
			fmt.Println("Please select a Cluster by running \n    topview get clusters \nfollowed by \n    topview use cluster <cluster-name>")
			fmt.Println("Please select a Zone by running \n    topview get zones \nfollowed by \n    topview use zone <zone-name>")
			return
		}
		session, err := GetSession(bytes)
		if err != nil {
			log.Fatalf("Unable to un-marshall session json file from user's HOME dir")
		}

		clusterId := session.ClusterId
		zoneId := session.ZoneId

		nodeResource, found := validNodeTypeResources[nodeType]
		if !found {
			fmt.Printf("Invalid node type [%v]\n", nodeType)
			return
		}
		url := baseRegistryURL + "/clusters/" + clusterId + "/zones/" + zoneId + nodeResource + "/" + nodeId + "/properties"

		fmt.Printf("Using cluster [%v]\n", session.ClusterName)
		fmt.Printf("Using Zone [%v]\n", session.ZoneName)
		fmt.Printf("Using Node ID [%v] of type [%v]\n", nodeId, nodeType)

		tab = tabular.New()
		tab.Col("property", "Property", 30)
		tab.Col("value", "Value", 20)
		tab.Col("propertyHandler", "Handler", 20)
		tab.Col("encoded", "Is Encoded?", 12)
		tab.Col("fileDestination", "File Destination", 50)

		response, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()

		if response != nil {
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}
			settings, err := GetSettings(body)
			if err != nil {
				log.Fatal(err)
			}

			if len(settings) == 0 {
				fmt.Printf("No settings found for the given nodeId [%v] of type [%v]\n", nodeId, nodeType)
				return
			}
			format := tab.Print("*")
			for _, setting := range settings {
				fmt.Printf(format, truncateString(setting.Property, 30), truncateString(setting.Value, 20), truncateString(setting.Handler, 20), setting.Encoded, truncateString(setting.FileDestination, 50))
			}
		}
	},
}
