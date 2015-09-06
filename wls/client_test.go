package wls

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
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
		Environment{ServerUrl: ts.URL, Username: "user", Password: "pass"})
}

func TestAcceptJsonHeaderCall(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/json" {
			t.Fail()
		}
	}))
	defer ts.Close()
	t.Log(ts.URL)

	requestResource(ts.URL, Environment{ServerUrl: ts.URL, Username: "user", Password: "pass"})
}

func CreateTestServerResourceRouters() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc(MONITOR_PATH+"/servers", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(servers_json))
	})
	r.HandleFunc(MONITOR_PATH+"/servers/{server}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		server := vars["server"]
		if server != "adminserver" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("No servername by that name"))
		} else {
			w.Write([]byte(single_server))
		}
	})
	return r
}

//func TestServerResourceRelatedCalls(t *testing.T) {
//	ts := httptest.NewServer(CreateTestServerResourceRouters())
//	defer ts.Close()
//	service := new(ServerService)
//
//	s, err := service.Servers()
//	if err != nil {
//		t.Fail()
//	}
//	if len(s) != 2 {
//		t.Fail()
//	}
//
//}
