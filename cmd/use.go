package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/aambhaik/topview/model"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
)

func init() {
	RootCmd.AddCommand(useCmd)
	useCmd.AddCommand(cmdUseCluster)
	useCmd.AddCommand(cmdUseZone)
}

var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Use available cluster or zone",
	Long:  `Use available cluster or zone in the topology`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please specify one of the following subcommands: [cluster|zone]")
	},
}

var cmdUseCluster = &cobra.Command{
	Use:   "cluster",
	Short: "Use available cluster ",
	Long:  `Use available cluster in the topology using the registry.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || len(args) > 1 {
			fmt.Println("Please specify only one valid cluster name")
			return
		}
		dir, err := homedir.Dir()
		if err != nil {
			log.Fatalf("Unable to detect user's HOME directory")
		}

		bytes, err := ioutil.ReadFile(dir + "/.topology.clusters.json")
		if err != nil {
			log.Fatalf("Unable to read clusters json file from user's HOME dir")
		}

		clusters, err := GetClusters(bytes)
		if err != nil {
			log.Fatal(err)
		}

		foundCluster := false

		for _, cluster := range clusters {
			if cluster.Name == args[0] {
				foundCluster = true
				//found our cluster, look up its ID
				session := model.Session{}
				session.ClusterId = cluster.ClusterId
				session.ClusterName = cluster.Name
				bytes, err = json.Marshal(session)

				err = ioutil.WriteFile(dir+"/.topviewer-session.json", bytes, 0644)
				if err != nil {
					log.Fatalf("Unable to write clusters json file from user's HOME dir, %v", err)

				}
				fmt.Println("Using cluster: ", cluster.Name)
				break
			}
		}
		if !foundCluster {
			fmt.Println("No matching cluster found: ", args[0])
		}
	},
}

var cmdUseZone = &cobra.Command{
	Use:   "zone",
	Short: "Use available zone in a cluster",
	Long:  `Use available zone in a cluster in the topology using the registry.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || len(args) > 1 {
			fmt.Println("Please specify only one valid zone name")
			return
		}
		dir, err := homedir.Dir()
		if err != nil {
			log.Fatalf("Unable to detect user's HOME directory")
		}

		bytes, err := ioutil.ReadFile(dir + "/.topology.zones.json")
		if err != nil {
			log.Fatalf("Unable to read zones json file from user's HOME dir")
		}

		zones, err := GetZones(bytes)
		if err != nil {
			log.Fatal(err)
		}
		foundZone := false
		for _, zone := range zones {
			if zone.Name == args[0] {
				//found our zone, look up its ID
				foundZone = true
				bytes, err := ioutil.ReadFile(dir + "/.topviewer-session.json")
				if err != nil {
					log.Fatalf("Unable to read session json file from user's HOME dir")
				}
				session, err := GetSession(bytes)
				if err != nil {
					log.Fatalf("Unable to un-marshall session json file from user's HOME dir")
				}

				session.ZoneId = zone.ZoneId
				session.ZoneName = zone.Name
				bytes, err = json.Marshal(session)

				err = ioutil.WriteFile(dir+"/.topviewer-session.json", bytes, 0644)
				fmt.Println("Using zone: ", zone.Name)

				break
			}
		}
		if !foundZone {
			fmt.Println("No matching zone found: ", args[0])
		}
	},
}
