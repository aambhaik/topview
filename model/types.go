package model

type Cluster struct {
	ClusterId   string `json:"clusterId"`
	Name        string `json:"name"`
	CreatedTime string `json:"createdTime"`
}

type Zone struct {
	ZoneId      string `json:"zoneId"`
	Name        string `json:"name"`
	Nodes       string `json:"nodes",omitempty`
	CreatedTime string `json:"createdTime"`
}

type Session struct {
	ClusterName string `json:"clusterName"`
	ClusterId   string `json:"clusterId"`
	ZoneName    string `json:"zoneName",omitempty`
	ZoneId      string `json:"zoneId",omitempty`
}

type DatabaseNode struct {
	NodeID    string `json:"nodeId"`
	Name      string `json:"name"`
	Host      string `json:"host"`
	AgentPort int    `json:"agentPort"`
	Status    string `json:"status"`
	Port      int    `json:"port"`
}
