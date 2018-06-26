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
	RootCmd.AddCommand(listCmd)
	listCmd.AddCommand(cmdListClusters)
	listCmd.AddCommand(cmdListZones)
	listCmd.AddCommand(cmdListNodes)
	listCmd.AddCommand(cmdListSettings)

	cmdListNodes.Flags().StringVar(&nodeType, "type", "", "type of the node")
	//cobra.MarkFlagRequired(cmdListNodes.Flags(), "type")

	cmdListSettings.Flags().StringVar(&nodeType, "type", "", "type of the node")
	cmdListSettings.Flags().StringVar(&nodeId, "nodeId", "", "ID of the node")
	cobra.MarkFlagRequired(cmdListSettings.Flags(), "type")
	cobra.MarkFlagRequired(cmdListSettings.Flags(), "nodeId")

}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List of clusters in the topology",
	Long:  `List of clusters in the topology`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please specify one of the following subcommands: [clusters|zones|nodes|settings]")
	},
}

var cmdListClusters = &cobra.Command{
	Use:   "clusters",
	Short: "List of clusters in the topology",
	Long:  `List of clusters in the topology`,
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
			err = ioutil.WriteFile(dir+"/.topview.clusters.json", bytes, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

var cmdListZones = &cobra.Command{
	Use:   "zones",
	Short: "List of zones in a cluster",
	Long:  `List of zones in a cluster in the topology`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := homedir.Dir()
		bytes, err := ioutil.ReadFile(dir + "/.topview-session.json")
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
			err = ioutil.WriteFile(dir+"/.topview.zones.json", bytes, 0644)

			format := tab.Print("*")
			for _, zone := range zones {
				fmt.Printf(format, zone.Name, zone.CreatedTime)
			}
		}
	},
}

var cmdListNodes = &cobra.Command{
	Use:   "nodes",
	Short: "List of nodes in a zone",
	Long:  `List of nodes in a zone in a cluster in the topology`,
	Run: func(cmd *cobra.Command, args []string) {
		flag.Parse()
		dir, err := homedir.Dir()
		bytes, err := ioutil.ReadFile(dir + "/.topview-session.json")
		if err != nil {
			fmt.Println("Please select a Cluster by running \n    topview get clusters \nfollowed by \n    topview use cluster <cluster-name>")
			fmt.Println("Please select a Zone by running \n    topview get zones \nfollowed by \n    topview use zone <zone-name>")
			return
		}
		session, err := GetSession(bytes)
		if err != nil {
			log.Fatalf("Unable to un-marshall session json file from user's HOME dir")
		}

		var getAllTypesOfNodes bool
		var found bool
		if len(nodeType) > 0 {
			_, found = validNodeTypeResources[nodeType]
			if !found {
				keys, err := GetKeysOfMap(validNodeTypeResources)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("Invalid node type [%v]\nPlease provide one of the following node types %v\n", nodeType, keys)
				return
			}
		} else {
			//get all type of nodes
			getAllTypesOfNodes = true
		}

		if getAllTypesOfNodes {
			err = GetAllNodes(session, &validNodeTypeResources, nil)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			err = GetAllNodes(session, &validNodeTypeResources, &nodeType)
			if err != nil {
				log.Fatalln(err)
			}
		}
	},
}

var cmdListSettings = &cobra.Command{
	Use:   "settings",
	Short: "List of settings for a specific node",
	Long:  `List of settings for a specific node in the topology`,
	Run: func(cmd *cobra.Command, args []string) {
		flag.Parse()
		dir, err := homedir.Dir()
		bytes, err := ioutil.ReadFile(dir + "/.topview-session.json")
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
				fmt.Printf(format, TruncateString(setting.Property, 30), TruncateString(setting.Value, 20), TruncateString(setting.Handler, 20), setting.Encoded, TruncateString(setting.FileDestination, 50))
			}
		}
	},
}
