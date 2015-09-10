package wls

import (
	"encoding/json"
	"fmt"
)

// Server is a specific Server instance deployed to the domain under the given AdminServer
type Server struct {
	Name                    string  `json:"name"`
	State                   string  `json:"state"`
	Health                  string  `json:"health"`
	ClusterName             string  `json:"clusterName,omitempty"`
	CurrentMachine          string  `json:",omitempty"`
	WeblogicVersion         string  `json:",omitempty"`
	OpenSocketsCurrentCount float64 `json:",omitempty"`
	HeapSizeCurrent         int     `json:",omitempty"`
	HeapFreeCurrent         int     `json:",omitempty"`
	JavaVersion             string  `json:",omitempty"`
	OsName                  string  `json:",omitempty"`
	OsVersion               string  `json:",omitempty"`
	JvmProcessorLoad        float64 `json:",omitempty"`
}

// Servers returns all servers configured in a domain and provides run-time information for each server, including the server state and health.
// isFullFormat determines whether to return a fully-filled out list of Servers, or only a shortened version of the Servers list.
func (s *AdminServer) Servers(isFullFormat bool) ([]Server, error) {
	url := fmt.Sprintf("%v%v/servers", s.AdminURL, MonitorPath)
	if isFullFormat {
		url = url + "?format=full"
	}
	w, err := requestAndUnmarshal(url, s)
	if err != nil {
		return nil, err
	}
	var servers []Server
	if err := json.Unmarshal(w.Body.Items, &servers); err != nil {
		return nil, err
	}
	return servers, nil
}

// Server returns information for a specified server in a domain, including the server state, health, and JVM heap availability.
func (s *AdminServer) Server(serverName string) (*Server, error) {
	url := fmt.Sprintf("%v%v/servers/%v", s.AdminURL, MonitorPath, serverName)
	w, err := requestAndUnmarshal(url, s)
	if err != nil {
		return nil, err
	}
	var server Server
	if err := json.Unmarshal(w.Body.Item, &server); err != nil {
		return nil, err
	}
	return &server, nil
}
