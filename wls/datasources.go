package wls

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type DataSourceService struct {
	client      *http.Client
	Environment Environment
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

func (s *DataSourceService) DataSources() ([]DataSource, error) {
	url := fmt.Sprintf("%v%v/datasources", s.Environment.ServerUrl, MONITOR_PATH)
	w, err := requestAndUnmarshal(url, s.Environment)
	if err != nil {
		return nil, err
	}
	var datasources []DataSource
	if err := json.Unmarshal(w.Body.Items, &datasources); err != nil {
		return nil, err
	}
	return datasources, nil
}

func (s *DataSourceService) DataSource(datasource_name string) (*DataSource, error) {
	url := fmt.Sprintf("%v%v/datasources/%v", s.Environment.ServerUrl, MONITOR_PATH, datasource_name)
	w, err := requestAndUnmarshal(url, s.Environment)
	if err != nil {
		return nil, err
	}
	var datasource DataSource
	if err := json.Unmarshal(w.Body.Item, &datasource); err != nil {
		return nil, err
	}
	return &datasource, nil
}
