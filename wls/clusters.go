package wls

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ClusterService struct {
	client      *http.Client
	Environment Environment
}

//type Cluster struct {
//	Name    string `json:"name"`
//	Servers []Server
//}

type ClusterServer struct {
	Name    string
	Servers []struct {
		Name                   string
		State                  string
		Health                 string
		IsClusterMaster        bool   `json:"clusterMaster,omitempty"`
		DropOutFrequency       string `json",omitempty"`
		ResendRequestsCount    int    `json:",omitempty"`
		FragmentsSentCount     int    `json:",omitempty"`
		FragmentsReceivedCount int    `json:",omitempty"`
	} `json:"servers,omitempty"`
}

func (s *ClusterService) Clusters() (*[]ClusterServer, error) {
	url := fmt.Sprintf("http://%v:%v%v/clusters", s.Environment.Server.Host, s.Environment.Server.Port, MONITOR_PATH)
	w, err := requestAndUnmarshal(url, s.Environment)
	if err != nil {
		return nil, err
	}
	var clusters []ClusterServer
	if err := json.Unmarshal(w.Body.Items, &clusters); err != nil {
		return nil, err
	}
	return &clusters, nil
}

func (s *ClusterService) Cluster(clustername string) (*ClusterServer, error) {
	url := fmt.Sprintf("http://%v:%v%v/clusters/%v", s.Environment.Server.Host, s.Environment.Server.Port, MONITOR_PATH, clustername)
	w, err := requestAndUnmarshal(url, s.Environment)
	if err != nil {
		return nil, err
	}
	var cluster ClusterServer
	if err := json.Unmarshal(w.Body.Item, &cluster); err != nil {
		return nil, err
	}
	return &cluster, nil
}
