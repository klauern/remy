package wls

import (
	"encoding/json"
	"fmt"
	"testing"
)

var servers_json = `{
  "body": {
    "items": [
      {
        "name": "adminserver",
        "state": "RUNNING",
        "health": " HEALTH_OK "
      },
      {
        "name": "ms1",
        "state": "SHUTDOWN",
        "health": ""
      }
     ]
   },
   "messages": [
  ]
 }`

func TestUnmarshalServersJson(t *testing.T) {
	wrapper, err := unmarshalWrapper([]byte(servers_json))
	if err != nil {
		t.Error(err)
	}
	t.Log(wrapper)
	var servers []Server
	err = json.Unmarshal(wrapper.Body.Items, &servers)
	if err != nil {
		t.Error(err)
	}
	var servers_json_tests = []struct {
		in  string
		out string
	}{
		{servers[0].Health, " HEALTH_OK "},
		{servers[0].Name, "adminserver"},
		{servers[0].State, "RUNNING"},
		{servers[1].Name, "ms1"},
		{servers[1].State, "SHUTDOWN"},
		{servers[1].Health, ""},
	}

	for _, tt := range servers_json_tests {
		if tt.in != tt.out {
			t.Errorf("want %q, got %q", tt.out, tt.in)
		}
	}

}

var singleServer = `{
  "body": {
    "item": {
      "name": "adminserver",
      "clusterName": null,

      "state": "RUNNING",

      "currentMachine": "machine-0",
      "weblogicVersion": "WebLogic Server 12.1.1.0.0 Thu May 5 01:17:16 2011 PDT",
      "openSocketsCurrentCount": 2,
      "health": "HEALTH_OK",

      "heapSizeCurrent": 536870912,
      "heapFreeCurrent": 39651944,
      "heapSizeMax": 1073741824,
      "javaVersion": "1.6.0_20",
      "osName": "Linux",
      "osVersion": "2.6.18-238.0.0.0.1.el5xen",

      "jvmProcessorLoad": 0.25
     }
    },
     "messages": [
    ]
  }`

func TestUnmarshalSingleServer(t *testing.T) {
	wrapper, err := unmarshalWrapper([]byte(singleServer))
	if err != nil {
		t.Error(err)
	}
	var server Server
	if err := json.Unmarshal(wrapper.Body.Item, &server); err != nil {
		t.Error(err)
	}
	//	t.Log(wrapper)
	var servers_json_tests = []struct {
		in  string
		out string
	}{
		{server.Name, "adminserver"},
		{server.ClusterName, ""},
		{server.State, "RUNNING"},
		{server.CurrentMachine, "machine-0"},
		{server.WeblogicVersion, "WebLogic Server 12.1.1.0.0 Thu May 5 01:17:16 2011 PDT"},
		{fmt.Sprint(server.OpenSocketsCurrentCount), "2"},
		{server.Health, "HEALTH_OK"},
		{fmt.Sprint(server.HeapSizeCurrent), "536870912"},
		{fmt.Sprint(server.HeapFreeCurrent), "39651944"},
		{server.JavaVersion, "1.6.0_20"},
		{server.OsName, "Linux"},
		{server.OsVersion, "2.6.18-238.0.0.0.1.el5xen"},
		{fmt.Sprint(server.JvmProcessorLoad), "0.25"},
	}

	for _, tt := range servers_json_tests {
		if tt.in != tt.out {
			t.Errorf("want %q, got %q", tt.out, tt.in)
		}
	}
}
