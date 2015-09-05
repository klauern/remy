package wls

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ServerService struct {
	client      *http.Client
	Environment Environment
}

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

func (s *ServerService) Servers() (*[]Server, error) {
	url := fmt.Sprintf("http://%v:%v%v/servers", s.Environment.Server.Host, s.Environment.Server.Port, MONITOR_PATH)
	w, err := requestAndUnmarshal(url, s.Environment)
	if err != nil {
		return nil, err
	}
	var servers []Server
	if err := json.Unmarshal(w.Body.Items, &servers); err != nil {
		return nil, err
	}
	return &servers, nil
}

func (s *ServerService) Server(servername string) (*Server, error) {
	url := fmt.Sprintf("http://%v:%v%v/servers/%v", s.Environment.Server.Host, s.Environment.Server.Port, MONITOR_PATH, servername)
	w, err := requestAndUnmarshal(url, s.Environment)
	if err != nil {
		return nil, err
	}
	var server Server
	if err := json.Unmarshal(w.Body.Item, &server); err != nil {
		return nil, err
	}
	return &server, nil
}
