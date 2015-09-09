package wls

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// See http://docs.oracle.com/cd/E23943_01/web.1111/e24682/toc.htm#RESTS149

// http(s)://host:port/management/tenant-monitoring/path

const (
	// path from the root / that points to where the RESTful Management API endpoint is located.  As of WLS 12.1.2,
	// this is assumed to be /management/tenant-monitoring
	MONITOR_PATH string = "/management/tenant-monitoring"
)

// Environment provides the configurable details necessary to request resources from a particular server.
// ServerUrl format should be similar to the following: "http(s)://[serverhost]:[adminport]"
type WlsAdminServer struct {
	ServerUrl string
	Username  string
	Password  string
}

// All requests sent to a WLS Rest endpoint are wrapped by a similar body and item or items tag.
// We simply wrap that so we can get to the meat of it in the underlying Server type
//
// Wrapper is composed of 3 pieces:
// - Body contains the result of the resource request, included as either a specific Item (Datasource, Server,
//   Cluster, or Application, or an array of []Item's of the same.
// - Messages contains any error-related messages related to the query.
type Wrapper struct {
	Body struct {
		Items json.RawMessage `json:"items,omitempty"`
		Item  json.RawMessage `json:"item,omitempty"`
	} `json:"body"`
	Messages []string `json:"messages,omitempty"`
}

func requestResource(url string, e *WlsAdminServer) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(e.Username, e.Password)
	return client.Do(req)
}

func requestAndUnmarshal(url string, e *WlsAdminServer) (*Wrapper, error) {
	if resp, err := request(url, e); err != nil {
		return nil, err
	} else {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return unmarshalWrapper(data)
	}
}

// Wrapper function for requestResource(), handling HTTP response codes before unmarshalling responses.
func request(url string, e *WlsAdminServer) (*http.Response, error) {
	resp, err := requestResource(url, e)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return resp, nil
	}
	return nil, fmt.Errorf("Invalid Response Code: %v\nResponse: %v", resp.StatusCode, resp.Body)
}

// Take the raw response from the server and attempt to unmarshal it into the Wrapper type.
func unmarshalWrapper(data []byte) (*Wrapper, error) {
	var w Wrapper
	err := json.Unmarshal(data, &w)
	if err != nil {
		return nil, err
	}
	return &w, nil
}
