package remy

import (
	"encoding/json"
	"fmt"
	"testing"
)

var applications = `{
    "body": {
        "items": [
            {
                "name": "appscopedejbs",
                "type": "ear",
                "state": "STATE_ACTIVE",
                "health": " HEALTH_OK"
            },
            {
                "name": "MyWebApp",
                "type": "war",
                "state": "STATE_NEW"
            }
        ]
    },
    "messages": []
}`

func TestUnmarshalMultipleApplications(t *testing.T) {
	wrapper, err := unmarshalWrapper([]byte(applications))
	if err != nil {
		t.Error(err)
	}
	t.Log(wrapper)
	var applications []Application
	if err := json.Unmarshal(wrapper.Body.Items, &applications); err != nil {
		t.Error(err)
	}
	if len(applications) == 0 {
		t.Errorf("Applications count should be 2, was 0")
	}
	var applicationsJSONTests = []struct {
		in  string
		out string
	}{
		{applications[0].Name, "appscopedejbs"},
		{applications[0].AppType, "ear"},
		{applications[0].State, "STATE_ACTIVE"},
		{applications[0].Health, " HEALTH_OK"},
		{applications[1].Name, "MyWebApp"},
		{applications[1].AppType, "war"},
		{applications[1].State, "STATE_NEW"},
	}

	for _, tt := range applicationsJSONTests {
		if tt.in != tt.out {
			t.Errorf("want %q, got %q", tt.out, tt.in)
		}
	}
}

var application = `{
    "body": {
        "item": {
            "name": "appscopedejbs",
            "type": "ear",
            "health": " HEALTH_OK ",
            "state": "STATE_ACTIVE",
            "targetStates": [
                {
                    "target": "ms1",
                    "state": "STATE_ACTIVE"
                },
                {
                    "target": "ms2",
                    "state": "STATE_ACTIVE"
                }
            ],
            "dataSources": [],
            "entities": [],
            "workManagers": [
                {
                    "name": "default",
                    "server": "ms1",
                    "pendingRequests": 0,
                    "completedRequests": 0
                }
            ],
            "minThreadsConstraints": [
                {
                    "name": "minThreadsConstraints-0",
                    "server": "ms1",
                    "completedRequests": 0,
                    "pendingRequests": 0,
                    "executingRequests": 0,
                    "outOfOrderExecutionCount": 0,
                    "mustRunCount": 0,
                    "maxWaitTime": 0,
                    "currentWaitTime": 0
                }
            ],
            "maxThreadsConstraints": [
                {
                    "name": "maxThreadsConstraints-0",
                    "server": "ms1",
                    "executingRequests": 0,
                    "deferredRequests": 0
                }
            ],
            "requestClasses": [
                {
                    "name": "requestClasses-0",
                    "server": "ms1",
                    "requestClassType": "fairshare",
                    "completedCount": 0,
                    "totalThreadUse": 0,
                    "pendingRequestCount": 0,
                    "virtualTimeIncrement": 0
                }
            ]
        }
    },
    "messages": []
}`

func TestUnmarshalApplication(t *testing.T) {
	wrapper, err := unmarshalWrapper([]byte(application))
	if err != nil {
		t.Error(err)
	}
	var application Application
	if err := json.Unmarshal(wrapper.Body.Item, &application); err != nil {
		t.Error(err)
	}
	//	t.Log(application)
	//	t.Log(application.TargetStates[0])
	//	t.Log(application.TargetStates[1])
	//	t.Log(application.WorkManagers[0])
	t.Log(application.MinThreadsConstraints[0])
	//	t.Log(application.MaxThreadsConstraints[0])
	t.Log(application.RequestClasses[0])
	var applicationJSONTests = []struct {
		in  string
		out string
	}{
		{application.Name, "appscopedejbs"},
		{application.AppType, "ear"},
		{application.State, "STATE_ACTIVE"},
		{application.Health, " HEALTH_OK "},
		{application.TargetStates[0].State, "STATE_ACTIVE"},
		{application.TargetStates[0].Target, "ms1"},
		{application.TargetStates[1].State, "STATE_ACTIVE"},
		{application.TargetStates[1].Target, "ms2"},
		{application.WorkManagers[0].Name, "default"},
		{application.WorkManagers[0].Server, "ms1"},
		{fmt.Sprint(application.WorkManagers[0].PendingRequests), "0"},
		{fmt.Sprint(application.WorkManagers[0].CompletedRequests), "0"},
		{application.MinThreadsConstraints[0].Name, "minThreadsConstraints-0"},
		{application.MinThreadsConstraints[0].Server, "ms1"},
		{fmt.Sprint(application.MinThreadsConstraints[0].CompletedRequests), "0"},
		{fmt.Sprint(application.MinThreadsConstraints[0].PendingRequests), "0"},
		{fmt.Sprint(application.MinThreadsConstraints[0].ExecutingRequests), "0"},
		{fmt.Sprint(application.MinThreadsConstraints[0].OutOfOrderExecutionCount), "0"},
		{fmt.Sprint(application.MinThreadsConstraints[0].MustRunCount), "0"},
		{fmt.Sprint(application.MinThreadsConstraints[0].MaxWaitTime), "0"},
		{fmt.Sprint(application.MinThreadsConstraints[0].CurrentWaitTime), "0"},
		{application.MaxThreadsConstraints[0].Name, "maxThreadsConstraints-0"},
		{application.MaxThreadsConstraints[0].Server, "ms1"},
		{fmt.Sprint(application.MaxThreadsConstraints[0].DeferredRequests), "0"},
		{fmt.Sprint(application.MaxThreadsConstraints[0].ExecutingRequests), "0"},
		{application.RequestClasses[0].Name, "requestClasses-0"},
		{application.RequestClasses[0].Server, "ms1"},
		{application.RequestClasses[0].RequestClassType, "fairshare"},
		{fmt.Sprint(application.RequestClasses[0].CompletedCount), "0"},
		{fmt.Sprint(application.RequestClasses[0].TotalThreadUse), "0"},
		{fmt.Sprint(application.RequestClasses[0].PendingRequestCount), "0"},
		{fmt.Sprint(application.RequestClasses[0].VirtualTimeIncrement), "0"},
	}

	for _, tt := range applicationJSONTests {
		if tt.in != tt.out {
			t.Errorf("want %q, got %q", tt.out, tt.in)
		}
	}
}
