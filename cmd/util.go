package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/InVisionApp/tabular"
	"github.com/aambhaik/topview/model"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
)

func GetClusters(body []byte) ([]*model.Cluster, error) {
	clusters := []*model.Cluster{}
	err := json.Unmarshal(body, &clusters)
	if err != nil {
		log.Fatal(err)
	}
	return clusters, err
}

func GetZones(body []byte) ([]*model.Zone, error) {
	zones := []*model.Zone{}
	err := json.Unmarshal(body, &zones)
	if err != nil {
		log.Fatal(err)
	}
	return zones, err
}

func GetNodes(body []byte) ([]*model.Node, error) {
	nodes := []*model.Node{}
	err := json.Unmarshal(body, &nodes)
	if err != nil {
		log.Fatal(err)
	}
	return nodes, err
}

func GetSettings(body []byte) ([]*model.NodeSetting, error) {
	settings := []*model.NodeSetting{}
	err := json.Unmarshal(body, &settings)
	if err != nil {
		log.Fatal(err)
	}
	return settings, err
}

func GetSession(body []byte) (*model.Session, error) {
	session := &model.Session{}
	err := json.Unmarshal(body, &session)
	if err != nil {
		log.Fatal(err)
	}
	return session, err
}

func TruncateString(str string, num int) string {
	bnoden := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		bnoden = str[0:num] + "..."
	}
	return bnoden
}

func GetKeysOfMap(m map[string]string) ([]string, error) {
	if m == nil {
		return nil, errors.New("Nil map")
	}
	keys := []string{}
	for k := range m {
		keys = append(keys, k)
	}
	return keys, nil
}

func GetAllNodes(session *model.Session, validNodeTypeResources *map[string]string, nodeType *string) error {
	clusterId := session.ClusterId
	zoneId := session.ZoneId
	fmt.Printf("Using cluster [%v]\n", session.ClusterName)
	fmt.Printf("Using Zone [%v]\n", session.ZoneName)

	tab = tabular.New()
	tab.Col("id", "Node ID", 37)
	tab.Col("type", "Node Type", 10)
	tab.Col("name", "Node Name", 20)
	tab.Col("status", "Node Status", 20)
	tab.Col("host", "Node Host", 20)
	tab.Col("agentPort", "Node Agent Port", 15)
	tab.Col("port", "Node Service Port(s)", 16)

	var url string

	allNodes := make(map[string][]*model.Node)
	for t, nodeResource := range *validNodeTypeResources {
		if nodeType != nil && t != *nodeType {
			continue
		}
		url = baseRegistryURL + "/clusters/" + clusterId + "/zones/" + zoneId + nodeResource
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
			allNodes[t] = nodes
		}
	}

	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	format := tab.Print("*")
	for nodeType, nodes := range allNodes {
		for _, node := range nodes {
			fmt.Printf(format, node.NodeID, nodeType, node.Name, node.Status, node.Host, node.AgentPort, node.Port)
		}
	}

	return nil
}
