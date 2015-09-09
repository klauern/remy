package wls

import (
	"encoding/json"
	"fmt"
)

type Cluster struct {
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

func (s *WlsAdminServer) Clusters(full_format bool) ([]Cluster, error) {
	url := fmt.Sprintf("%v%v/clusters", s.AdminUrl, MONITOR_PATH)
	if full_format {
		url = url + "?format=full"
	}
	w, err := requestAndUnmarshal(url, s)
	if err != nil {
		return nil, err
	}
	var clusters []Cluster
	if err := json.Unmarshal(w.Body.Items, &clusters); err != nil {
		return nil, err
	}
	return clusters, nil
}

func (s *WlsAdminServer) Cluster(clustername string) (*Cluster, error) {
	url := fmt.Sprintf("%v%v/clusters/%v", s.AdminUrl, MONITOR_PATH, clustername)
	w, err := requestAndUnmarshal(url, s)
	if err != nil {
		return nil, err
	}
	var cluster Cluster
	if err := json.Unmarshal(w.Body.Item, &cluster); err != nil {
		return nil, err
	}
	return &cluster, nil
}
