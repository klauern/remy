package wls

import (
	"encoding/json"
	"fmt"
)

// Cluster is the underlying struct for a single Cluster in a domain.  Each domain may have multiple Cluster instances,
// each having their own ClusterMaster, deployed state, etc.
type Cluster struct {
	Name    string
	Servers []struct {
		Name                   string
		State                  string
		Health                 string
		IsClusterMaster        bool   `json:"clusterMaster,omitempty"`
		DropOutFrequency       string `json:",omitempty"`
		ResendRequestsCount    int    `json:",omitempty"`
		FragmentsSentCount     int    `json:",omitempty"`
		FragmentsReceivedCount int    `json:",omitempty"`
	} `json:"servers,omitempty"`
}

// Clusters returns all clusters configured in a domain and provides run-time information for each cluster and for each cluster's member servers, including all the member servers' state and health.
func (s *AdminServer) Clusters(fullFormat bool) ([]Cluster, error) {
	url := fmt.Sprintf("%v%v/clusters", s.AdminURL, MonitorPath)
	if fullFormat {
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

// Cluster returns run-time information for the specified cluster and its member servers, including the member servers' state and health.
func (s *AdminServer) Cluster(clusterName string) (*Cluster, error) {
	url := fmt.Sprintf("%v%v/clusters/%v", s.AdminURL, MonitorPath, clusterName)
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
