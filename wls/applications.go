package wls

import (
	"encoding/json"
	"fmt"
)

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

func (s *WlsAdminServer) Applications(full_format bool) ([]Application, error) {
	url := fmt.Sprintf("%v%v/applications", s.AdminUrl, MONITOR_PATH)
	if full_format {
		url = url + "?format=full"
	}
	w, err := requestAndUnmarshal(url, s)
	if err != nil {
		return nil, err
	}
	var applications []Application
	if err := json.Unmarshal(w.Body.Items, &applications); err != nil {
		return nil, err
	}
	return applications, nil
}

func (s *WlsAdminServer) Application(app string) (*Application, error) {
	url := fmt.Sprintf("%v%v/applications/%v", s.AdminUrl, MONITOR_PATH, app)
	w, err := requestAndUnmarshal(url, s)
	if err != nil {
		return nil, err
	}
	var application Application
	if err := json.Unmarshal(w.Body.Item, &application); err != nil {
		return nil, err
	}
	return &application, nil
}
