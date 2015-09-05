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
	MONITOR_PATH string = "/management/tenant-monitoring"
)

type Environment struct {
	Server   WLSServer
	Username string
	Password string
}

type WLSServer struct {
	Host string
	Port int
}

type WebLogicService struct {
	Environment Environment
	Server      WLSServer
	client      *http.Client

	Servers *ServerService

	Clusters *ClusterService

	DataSources *DataSourceService

	Applications *ApplicationService
}

// All requests sent to a WLS Rest endpoint are wrapped by a similar body and item or items tag.
// We simply wrap that so we can get to the meat of it in the underlying Server type
type Wrapper struct {
	Body struct {
		Items json.RawMessage `json:"items,omitempty"`
		Item  json.RawMessage `json:"item,omitempty"`
	} `json:"body"`
	Messages []string `json:"messages,omitempty"`
}

func requestResource(url string, e Environment) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(e.Username, e.Password)
	return client.Do(req)
}

func unmarshalResponse(resp *http.Response, w *Wrapper) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		err := json.Unmarshal(bodyBytes, w)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Invalid Response Code; %v", resp.StatusCode)
}

func requestAndUnmarshal(url string, e Environment) (*Wrapper, error) {
	resp, err := requestResource(url, e)
	if err != nil {
		return nil, err
	}
	var w Wrapper
	if err = unmarshalResponse(resp, &w); err != nil {
		return nil, err
	}
	return &w, nil
}

func (w *Wrapper) getItem(r *interface{}) error {
	if err := json.Unmarshal(w.Body.Item, r); err != nil {
		return err
	}
	return nil
}

func (w *Wrapper) getItems(r *interface{}) error {
	if err := json.Unmarshal(w.Body.Items, &r); err != nil {
		return err
	}
	return nil
}
