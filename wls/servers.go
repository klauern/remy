package wls

import (
	"encoding/json"
	"fmt"
)

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

func (s *WlsAdminServer) Servers(full_format bool) ([]Server, error) {
	url := fmt.Sprintf("%v%v/servers", s.ServerUrl, MONITOR_PATH)
	if full_format {
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

func (s *WlsAdminServer) Server(servername string) (*Server, error) {
	url := fmt.Sprintf("%v%v/servers/%v", s.ServerUrl, MONITOR_PATH, servername)
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
