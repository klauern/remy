package wls

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ApplicationService struct {
	client      *http.Client
	Environment Environment
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

func (s *ApplicationService) Applications(full_format bool) ([]Application, error) {
	url := fmt.Sprintf("%v%v/applications", s.Environment.ServerUrl, MONITOR_PATH)
	w, err := requestAndUnmarshal(url, s.Environment)
	if err != nil {
		return nil, err
	}
	var applications []Application
	if err := json.Unmarshal(w.Body.Items, &applications); err != nil {
		return nil, err
	}
	return applications, nil
}

func (s *ApplicationService) Application(app string) (*Application, error) {
	url := fmt.Sprintf("%v%v/applications/%v", s.Environment.ServerUrl, MONITOR_PATH, app)
	w, err := requestAndUnmarshal(url, s.Environment)
	if err != nil {
		return nil, err
	}
	var application Application
	if err := json.Unmarshal(w.Body.Item, &application); err != nil {
		return nil, err
	}
	return &application, nil
}
