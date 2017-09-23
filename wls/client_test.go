package wls

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestClientAuthenticatedCalls(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, p, _ := r.BasicAuth()
		if u != "user" && p != "pass" {
			t.Fail()
		}
	}))
	defer ts.Close()
	t.Log(ts.URL)
	requestResource(ts.URL,
		&AdminServer{AdminURL: ts.URL, Username: "user", Password: "pass"})
}

func TestAcceptJsonHeaderCall(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/json" {
			t.Fail()
		}
	}))
	defer ts.Close()
	t.Log(ts.URL)

	requestResource(ts.URL, &AdminServer{AdminURL: ts.URL, Username: "user", Password: "pass"})
}

func CreateTestServerResourceRouters() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc(MonitorPath+"/servers", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(serversJSON))
	})
	r.HandleFunc(MonitorPath+"/servers/{server}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		server := vars["server"]
		if server != "adminserver" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("No servername by that name"))
		} else {
			w.Write([]byte(singleServer))
		}
	})
	return r
}

func TestServerResourceRelatedCalls(t *testing.T) {
	ts := httptest.NewServer(CreateTestServerResourceRouters())
	defer ts.Close()
	service := new(AdminServer)

	service.AdminURL = ts.URL
	service.Username = "user"
	service.Password = "pass"
	s, err := service.Servers(false)
	if err != nil {
		t.Fail()
	}
	if len(s) != 2 {
		t.Fail()
	}

	server, err := service.Server("adminserver")
	if err != nil {
		t.Fail()
	}
	var serverJSONTests = []struct {
		in  string
		out string
	}{
		{server.Name, "adminserver"},
		{server.ClusterName, ""},
		{server.State, "RUNNING"},
		{server.CurrentMachine, "machine-0"},
		{server.WebLogicVersion, "WebLogic Server 12.1.1.0.0 Thu May 5 01:17:16 2011 PDT"},
		{fmt.Sprint(server.OpenSocketsCurrentCount), "2"},
		{server.Health, "HEALTH_OK"},
		{fmt.Sprint(server.HeapSizeCurrent), "536870912"},
		{fmt.Sprint(server.HeapFreeCurrent), "39651944"},
		{server.JavaVersion, "1.6.0_20"},
		{server.OsName, "Linux"},
		{server.OsVersion, "2.6.18-238.0.0.0.1.el5xen"},
		{fmt.Sprint(server.JvmProcessorLoad), "0.25"},
	}

	for _, tt := range serverJSONTests {
		if tt.in != tt.out {
			t.Errorf("want %q, got %q", tt.out, tt.in)
		}
	}
}
