package wls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/BurntSushi/toml"
)
// See http://docs.oracle.com/cd/E23943_01/web.1111/e24682/toc.htm#RESTS149

// http(s)://host:port/management/tenant-monitoring/path

const (
	// MonitorPath is the REST resource path from the root / that points to where the RESTful Management API endpoint is located.  As of WLS 12.1.2,
	// this is assumed to be /management/tenant-monitoring
	MonitorPath string = "/management/tenant-monitoring"
)

// AdminServer contains the configurable details necessary to request resources from a particular AdminServer.
// ServerUrl format should be similar to the following: "http(s)://[serverhost]:[adminport]"
type AdminServer struct {
	AdminURL string
	Username string
	Password string
}

// Wrapper handles all responses sent back from a WLS Rest endpoint.  These responses are wrapped by a similar body and item or items tag.
// We simply wrap that so we can get to the meat of it in the underlying Server type
//
// Wrapper is composed of 3 pieces:
// - Body contains
//   - the result of the resource request, included as either a specific Item (Datasource, Server,
//     Cluster, or Application, or
//   - an array of []Item's of the same.
// - Messages contains any error-related messages related to the query.
type Wrapper struct {
	Body struct {
		Items json.RawMessage `json:"items,omitempty"`
		Item  json.RawMessage `json:"item,omitempty"`
	} `json:"body"`
	Messages []string `json:"messages,omitempty"`
}

// requestResource is a wrapper around an http.Client{} instance assuming the following:
// - the request is a GET
// - assumes a JSON Accept header
// - set Basic Authentication based on the *AdminServer passed in
//
// returns the *http.Response or an error
func requestResource(url string, e *AdminServer) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(e.Username, e.Password)
	return client.Do(req)
}

func requestAndUnmarshal(url string, e *AdminServer) (*Wrapper, error) {
	resp, err := request(url, e)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return unmarshalWrapper(data)
}

// Wrapper function for requestResource(), handling HTTP response codes before unmarshalling responses.
func request(url string, e *AdminServer) (*http.Response, error) {
	resp, err := requestResource(url, e)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return resp, nil
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return nil, fmt.Errorf("Invalid Response Code: %v\nResponse: \n%v", resp.StatusCode, string(body))
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

// EncodeConfigFile will take an AdminServer instance and convert it to a TOML-formatted byte buffer.
// This TOML byte buffer can then be serialized to disk as a configuration file.
func (cfg AdminServer) EncodeConfigFile() *bytes.Buffer {
	var buf bytes.Buffer
	enc := toml.NewEncoder(&buf)
	err := enc.Encode(cfg)
	if err != nil {
		panic(fmt.Errorf("unable to encode wlsrest configuration: %s \n", err))
	}
	return &buf
}
