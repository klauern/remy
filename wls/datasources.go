package wls

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// DataSource is a specific Data Source in a domain, including all deployed DataSourceInstance's.
type DataSource struct {
	Name      string
	Type      string
	Instances []DataSourceInstance `json:"instances,omitempty"`
}

// DataSourceInstance provides run-time information for a data source instance
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

// RacInstance provide Oracle RAC instance statistics
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

// DataSources returns all generic and GridLink JDBC data sources configured in the domain, and provides run-time information for each data source.
func (s *AdminServer) DataSources(isFullFormat bool) ([]DataSource, error) {
	url := fmt.Sprintf("%v%v/datasources", s.AdminURL, MonitorPath)
	if isFullFormat {
		url = url + "?format=full"
	}
	w, err := requestAndUnmarshal(url, s)
	if err != nil {
		return nil, err
	}
	var datasources []DataSource
	if err := json.Unmarshal(w.Body.Items, &datasources); err != nil {
		return nil, err
	}
	return datasources, nil
}

// DataSource returns run-time information for the specified data source, including Oracle RAC statistics for GridLink data sources.
func (s *AdminServer) DataSource(dataSourceName string) (*DataSource, error) {
	url := fmt.Sprintf("%v%v/datasources/%v", s.AdminURL, MonitorPath, dataSourceName)
	w, err := requestAndUnmarshal(url, s)
	if err != nil {
		return nil, err
	}
	var dataSource DataSource
	if err := json.Unmarshal(w.Body.Item, &dataSource); err != nil {
		return nil, err
	}
	return &dataSource, nil
}
