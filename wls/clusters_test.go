package wls

import (
	"encoding/json"
	"fmt"
	"testing"
)

var single_cluster = `{
    "body": {
        "item": {
            "name": "mycluster1",
            "servers": [
                {
                    "name": "ms1",
                    "state": "RUNNING",
                    "health": "OK",
                    "clusterMaster": false,
                    "dropOutFrequency": "Never",
                    "resendRequestsCount": 0,
                    "fragmentsSentCount": 3708,
                    "fragmentsReceivedCount": 3631
                },
                {
                    "name": "ms2",
                    "state": "RUNNING",
                    "health": "OK"
                }
            ]
        }
    },
    "messages": []
}`

func TestUnmarshalSingleCluster(t *testing.T) {
	wrapper, err := unmarshalWrapper([]byte(single_cluster))
	if err != nil {
		t.Error(err)
	}
	var cluster Cluster
	if err := json.Unmarshal(wrapper.Body.Item, &cluster); err != nil {
		t.Error(err)
	}
	if len(cluster.Servers) == 0 {
		t.Errorf("Servers in wrapper.Body.Item is 0, should be 2")
	}
	var servers_json_tests = []struct {
		in  string
		out string
	}{
		{cluster.Name, "mycluster1"},
		{cluster.Servers[0].Name, "ms1"},
		{cluster.Servers[0].State, "RUNNING"},
		{cluster.Servers[0].Health, "OK"},
		{fmt.Sprint(cluster.Servers[0].IsClusterMaster), "false"},
		{cluster.Servers[0].DropOutFrequency, "Never"},
		{fmt.Sprint(cluster.Servers[0].ResendRequestsCount), "0"},
		{fmt.Sprint(cluster.Servers[0].FragmentsSentCount), "3708"},
		{fmt.Sprint(cluster.Servers[0].FragmentsReceivedCount), "3631"},
		{cluster.Servers[1].Name, "ms2"},
		{cluster.Servers[1].State, "RUNNING"},
		{cluster.Servers[1].Health, "OK"},
		{fmt.Sprint(cluster.Servers[1].DropOutFrequency), ""},
	}

	for _, tt := range servers_json_tests {
		if tt.in != tt.out {
			t.Errorf("want %q, got %q", tt.out, tt.in)
		}
	}
}

var clusters = `{
    "body": {
        "items": [
            {
                "name": "mycluster1",
                "servers": [
                    {
                        "name": "ms1",
                        "state": "RUNNING",
                        "health": "HEALTH_OK"
                    },
                    {
                        "name": "ms2",
                        "state": "RUNNING",
                        "health": "HEALTH_OVERLOADED"
                    }
                ]
            }
        ]
    },
    "messages": []
}`

func TestUnmarshalMultipleClusters(t *testing.T) {
	wrapper, err := unmarshalWrapper([]byte(clusters))
	if err != nil {
		t.Error(err)
	}
	var clusters []Cluster
	if err := json.Unmarshal(wrapper.Body.Items, &clusters); err != nil {
		t.Error(err)
	}
	if len(clusters) == 0 {
		t.Errorf("Clusters count should be 1, was 0")
	}
	var servers_json_tests = []struct {
		in  string
		out string
	}{
		{clusters[0].Name, "mycluster1"},
		{clusters[0].Servers[0].Name, "ms1"},
		{clusters[0].Servers[0].State, "RUNNING"},
		{clusters[0].Servers[0].Health, "HEALTH_OK"},
		{clusters[0].Servers[1].Name, "ms2"},
		{clusters[0].Servers[1].State, "RUNNING"},
		{clusters[0].Servers[1].Health, "HEALTH_OVERLOADED"},
	}

	for _, tt := range servers_json_tests {
		if tt.in != tt.out {
			t.Errorf("want %q, got %q", tt.out, tt.in)
		}
	}
}
