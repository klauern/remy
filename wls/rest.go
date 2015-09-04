/*

 */
package wls

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// See http://docs.oracle.com/cd/E23943_01/web.1111/e24682/toc.htm#RESTS149

// http(s)://host:port/management/tenant-monitoring/path

const (
	MONITOR_PATH string = "/management/tenant-monitoring"
)

type Environment struct {
	Server   WLSServer
	Username string
	Password string
}

type WLSServer struct {
	Host string
	Port int
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

// All requests sent to a WLS Rest endpoint are wrapped by a similar body and item or items tag.
// We simply wrap that so we can get to the meat of it in the underlying Server type
type ServerWrapper struct {
	Body struct {
		Items []Server `json:"items,omitempty"`
		Item  Server   `json:"item,omitempty"`
	} `json:"body"`
	Messages []string `json:"messages,omitempty"`
}

func WLSRestReq(e Environment, url string, wls interface{}) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(e.Username, e.Password)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		err := json.Unmarshal(bodyBytes, wls)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Invalid Response Code; %v", resp.StatusCode)
}

func (w *ServerWrapper) Servers(e Environment) error {
	url := fmt.Sprintf("http://%v:%v%v/servers", e.Server.Host, e.Server.Port, MONITOR_PATH)
	err := WLSRestReq(e, url, &w)
	if err != nil {
		return err
	}
	return nil
}

func (w *ServerWrapper) ServerInfo(e Environment, servername string) error {
	url := fmt.Sprintf("http://%v:%v%v/servers/%v", e.Server.Host, e.Server.Port, MONITOR_PATH, servername)
	err := WLSRestReq(e, url, w)
	if err != nil {
		return err
	}
	return nil
}

type Cluster struct {
	Name    string `json:"name"`
	Servers []Server
}

type ClusterServer struct {
	Name    string
	Servers []struct {
		Name                   string
		State                  string
		Health                 string
		IsClusterMaster        bool `json:"clusterMaster,omitempty"`
		DropOutFrequency       string
		ResendRequestsCount    int
		FragmentsSentCount     int
		FragmentsReceivedCount int
	} `json:"servers,omitempty"`
}

type ClusterWrapper struct {
	Body struct {
		Item  ClusterServer   `json:"item,omitempty"`
		Items []ClusterServer `json:"items,omitempty"`
	} `json:"body"`
	Messages []string
}

func (w *ClusterWrapper) Clusters(e Environment, full_format bool) error {
	url := fmt.Sprintf("http://%v:%v%v/clusters", e.Server.Host, e.Server.Port, MONITOR_PATH)
	err := WLSRestReq(e, url, &w)
	if err != nil {
		return err
	}
	return nil
}

func (w *ClusterWrapper) Cluster(e Environment, clustername string) error {
	url := fmt.Sprintf("http://%v:%v%v/clusters", e.Server.Host, e.Server.Port, MONITOR_PATH)
	err := WLSRestReq(e, url, w)
	if err != nil {
		return err
	}
	return nil
}
