package cmd

import (
	"encoding/json"
	"github.com/aambhaik/topview/model"
	"log"
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

func truncateString(str string, num int) string {
	bnoden := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		bnoden = str[0:num] + "..."
	}
	return bnoden
}
