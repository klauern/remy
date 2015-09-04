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

type ApplicationWrapper struct {
	Body struct {
		Item  Application   `json:"item,omitempty"`
		Items []Application `json:"items,omitempty`
	}
	Messages []string
}

type Application struct {
	Name                  string
	AppType               string `json:"type"`
	State                 string
	Health                string
	TargetStates          []TargetState          `json:"targetStates,omitempty"`
	DataSources           []AppDataSource        `json:"dataSources,omitempty"`
	WorkManagers          []WorkManager          `json:"workManagers,omitempty"`
	MinThreadsConstraints []MinThreadsConstraint `json:"minThreadsConstraints,omitempty"`
	MaxThreadsConstraints []MaxThreadsConstraint `json:"maxThreadsConstraints,omitempty"`
	RequestClasses        []RequestClass         `json:"requestClasses,omitempty"`
}

type TargetState struct {
	Target string
	State  string
}

type AppDataSource struct {
	Name   string
	Server string
	State  string
}

type WorkManager struct {
	Name              string
	Server            string
	PendingRequests   int
	CompletedRequests int
}

type MinThreadsConstraint struct {
	Name                     string
	Server                   string
	PendingRequests          int
	CompletedRequests        int
	ExecutingRequests        int
	OutOfOrderExecutionCount int
	MustRunCount             int
	MaxWaitTime              int
	CurrentWaitTime          int
}

type MaxThreadsConstraint struct {
	Name              string
	Server            string
	ExecutingRequests int
	DeferredRequests  int
}

type RequestClass struct {
	Name                 string
	Server               string
	RequestClassType     string
	CompletedCount       int
	TotalThreadUse       int
	PendingRequestCount  int
	VirtualTimeIncrement int
}

func (w *ApplicationWrapper) Applications(e Environment, full_format bool) error {
	url := fmt.Sprintf("http://%v:%v%v/applications", e.Server.Host, e.Server.Port, MONITOR_PATH)
	err := WLSRestReq(e, url, &w)
	if err != nil {
		return err
	}
	return nil
}

func (w *ApplicationWrapper) Application(e Environment, app string) error {
	url := fmt.Sprintf("http://%v:%v%v/applications/%v", e.Server.Host, e.Server.Port, MONITOR_PATH, app)
	err := WLSRestReq(e, url, &w)
	if err != nil {
		return err
	}
	return nil
}

type DataSourceWrapper struct {
	Body struct {
		Item  DataSource   `json:"item,omitempty"`
		Items []DataSource `json:"items,omitempty`
	}
	Messages []string
}

type DataSource struct {
	Name      string
	Type      string
	Instances []DataSourceInstance `json:"instances,omitempty"`
}

type DataSourceInstance struct {
	Server                             string
	State                              string
	Enabled                            bool
	VersionJDBCDriver                  string        `json:",omitempty"`
	ActiveConnectionsAverageCount      int           `json:",omitempty"`
	ActiveConnectionsCurrentCount      int           `json:",omitempty"`
	ActiveConnectionsHighCount         int           `json:",omitempty"`
	ConnectionDelayTime                int           `json:",omitempty"`
	ConnectionsTotalCount              int           `json:",omitempty"`
	CurrCapacity                       int           `json:",omitempty"`
	CurrCapacityHighCount              int           `json:",omitempty"`
	FailedReserveRequestCount          int           `json:",omitempty"`
	FailuresToReconnectCount           int           `json:",omitempty"`
	HighestNumAvailable                int           `json:",omitempty"`
	LeakedConnectionCount              int           `json:",omitempty"`
	NumAvailable                       int           `json:",omitempty"`
	NumUnavailable                     int           `json:",omitempty"`
	PrepStmtCacheAccessCount           int           `json:",omitempty"`
	PrepStmtCacheAddCount              int           `json:",omitempty"`
	PrepStmtCacheCurrentSize           int           `json:",omitempty"`
	PrepStmtCacheDeleteCount           int           `json:",omitempty"`
	PrepStmtCacheHitCount              int           `json:",omitempty"`
	PrepStmtCacheMissCount             int           `json:",omitempty"`
	ReserveRequestCount                int           `json:",omitempty"`
	WaitSecondsHighCount               int           `json:",omitempty"`
	WaitingForConnectionCurrentCount   int           `json:",omitempty"`
	WaitingForConnectionFailureTotal   int           `json:",omitempty"`
	WaitingForConnectionHighCount      int           `json:",omitempty"`
	WaitingForConnectionSuccessTotal   int           `json:",omitempty"`
	WaitingForConnectionTotal          int           `json:",omitempty"`
	SuccessfulRCLBBasedBorrowCount     int           `json:",omitempty"`
	FailedRCLBBasedBorrowCount         int           `json:",omitempty"`
	SuccessfulAffinityBasedBorrowCount int           `json:",omitempty"`
	FailedAffinityBasedBorrowCount     int           `json:",omitempty"`
	RacInstances                       []RacInstance `json:",omitempty"`
}

type RacInstance struct {
	InstanceName                  string `json:",omitempty"`
	State                         string `json:",omitempty"`
	Enabled                       bool   `json:",omitempty"`
	Signature                     string `json:",omitempty"`
	CurrentWeight                 int    `json:",omitempty"`
	ActiveConnectionsCurrentCount int    `json:",omitempty"`
	ReserveRequestCount           int    `json:",omitempty"`
	ConnectionsTotalCount         int    `json:",omitempty"`
	CurrCapacity                  int    `json:",omitempty"`
	NumAvailable                  int    `json:",omitempty"`
	NumUnavailable                int    `json:",omitempty"`
}

func (w *DataSourceWrapper) DataSources(e Environment, full_format bool) error {
	url := fmt.Sprintf("http://%v:%v%v/datasources", e.Server.Host, e.Server.Port, MONITOR_PATH)
	err := WLSRestReq(e, url, &w)
	if err != nil {
		return err
	}
	return nil
}

func (w *DataSourceWrapper) DataSource(e Environment, datasource_name string) error {
	url := fmt.Sprintf("http://%v:%v%v/datasources/%v", e.Server.Host, e.Server.Port, MONITOR_PATH, datasource_name)
	err := WLSRestReq(e, url, &w)
	if err != nil {
		return err
	}
	return nil
}
