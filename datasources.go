package remy

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

// GoString produces a GoString of a DataSource that will be more pleasant to the eyes for a command-line interface.
func (d *DataSource) GoString() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Name: %v|Type: %v\n", d.Name, d.Type))
	buffer.WriteString("Datasource Instances\n")
	for i := range d.Instances {
		inst := d.Instances[i]
		buffer.WriteString(fmt.Sprintf("Server: %v|State: %v|Enabled? %v|JDBC Driver Ver: %v\n", inst.Server, inst.State, inst.Enabled, inst.VersionJDBCDriver))
		buffer.WriteString(fmt.Sprintf("Active Connections:  Current: %v|High: %v|Average: %v\n", inst.ActiveConnectionsCurrentCount, inst.ActiveConnectionsHighCount, inst.ActiveConnectionsAverageCount))
		buffer.WriteString(fmt.Sprintf("Connection Delay Time: %v|Total Count: %v\n", inst.ConnectionDelayTime, inst.ConnectionsTotalCount))
		buffer.WriteString(fmt.Sprintf("Capacity - Current : %v|High: %v|Num Avail: %v|Num Unavailable: %v", inst.CurrCapacity, inst.CurrCapacityHighCount,
			inst.NumAvailable, inst.NumUnavailable))
		buffer.WriteString("Reserve Requests\n----------------\n")
		buffer.WriteString(fmt.Sprintf("Failures -  Reserve Reqeuest Count: %v}Reconnect Count: %v|RCLB Based Borrow : %v|Affinity Based Borrow: %v|Leaked Connection Count: %v\n",
			inst.FailedReserveRequestCount, inst.FailuresToReconnectCount, inst.FailedRCLBBasedBorrowCount, inst.FailedAffinityBasedBorrowCount, inst.LeakedConnectionCount))
		buffer.WriteString(fmt.Sprintf("Prepared Statement Cache - Access cnt: %v|Add Cnt: %v|Current Size: %v|Delete Count: %v|Hit Count: %v|Miss Count: %v\n",
			inst.PrepStmtCacheAccessCount, inst.PrepStmtCacheAddCount, inst.PrepStmtCacheCurrentSize, inst.PrepStmtCacheDeleteCount,
			inst.PrepStmtCacheHitCount, inst.PrepStmtCacheMissCount))
		buffer.WriteString("Waiting Statistics: \n------------------\n")
		buffer.WriteString(fmt.Sprintf("Wait Seconds High: %v|Currently Waiting: %v|Wait Fail Tot: %v|High Count: %v\n", inst.WaitSecondsHighCount,
			inst.WaitingForConnectionCurrentCount, inst.WaitingForConnectionFailureTotal, inst.WaitingForConnectionHighCount))
		buffer.WriteString(fmt.Sprintf("Connection Success Tot: %v|Waiting Tot: %v\n", inst.WaitingForConnectionSuccessTotal, inst.WaitingForConnectionTotal))
		buffer.WriteString("RCLB Stats\n----------\n")
		buffer.WriteString(fmt.Sprintf("Successful RCLB Borrow Count: %v|Failed RCLB Borrow Count: %v\n", inst.SuccessfulRCLBBasedBorrowCount, inst.FailedRCLBBasedBorrowCount))
		buffer.WriteString("Affinity Borrows:\n")
		buffer.WriteString(fmt.Sprintf("Success Borrows: %v|Failed Borrows: %v\n", inst.SuccessfulAffinityBasedBorrowCount, inst.FailedAffinityBasedBorrowCount))
		for j := range inst.RacInstances {
			rac := inst.RacInstances[j]
			buffer.WriteString("RAC Instances\n-------------\n")
			buffer.WriteString(fmt.Sprintf("Name: %v|State: %v|Enabled: %v|Signature: %v\n", rac.InstanceName, rac.State, rac.Enabled, rac.Signature))
			buffer.WriteString(fmt.Sprintf("Current Weight: %v|Active Conn: %v|Reserve Req. Cnt: %v\n", rac.CurrentWeight, rac.ActiveConnectionsCurrentCount, rac.ReserveRequestCount))
			buffer.WriteString(fmt.Sprintf("Total Connections: %v|Capacity: %v|Num Available: %v|Num Unavailable: %v\n", rac.ConnectionsTotalCount,
				rac.CurrCapacity, rac.NumAvailable, rac.NumUnavailable))
		}
	}
	return buffer.String()
}

// DataSources returns all generic and GridLink JDBC data sources configured in the domain, and provides run-time information for each data source.
func (a *AdminServer) DataSources(isFullFormat bool) ([]DataSource, error) {
	url := fmt.Sprintf("%v%v/datasources", a.AdminURL, MonitorPath)
	if isFullFormat {
		url = url + "?format=full"
	}
	w, err := requestAndUnmarshal(url, a)
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
func (a *AdminServer) DataSource(dataSourceName string) (*DataSource, error) {
	url := fmt.Sprintf("%v%v/datasources/%v", a.AdminURL, MonitorPath, dataSourceName)
	w, err := requestAndUnmarshal(url, a)
	if err != nil {
		return nil, err
	}
	var dataSource DataSource
	if err := json.Unmarshal(w.Body.Item, &dataSource); err != nil {
		return nil, err
	}
	return &dataSource, nil
}
