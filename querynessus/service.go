package querynessus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/go-querystring/query"
)

var TenablePluginsServiceEndpoint = "https://cloud.tenable.com/plugins/plugin"

type TenableCredentials struct {
	AccessKey string
	SecretKey string
}

type RequestParams struct {
}

func FetchPlugins(creds TenableCredentials, params *RequestParams) ([]PluginDetails, error) {

	v, _ := query.Values(params)
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", TenablePluginsServiceEndpoint, nil)
	req.URL.RawQuery = v.Encode()
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-ApiKeys", "accessKey="+creds.AccessKey+";secretKey="+creds.SecretKey)
	if err != nil {
		log.Fatalf("Failed to create request")
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to submit request")
		return nil, err
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	var pluginPage PluginListPage
	err = decoder.Decode(&pluginPage)
	if err != nil {
		log.Fatalf("Failed to decode results")
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(body)
		return nil, err
	}
	return pluginPage.Data.PluginDetails, nil
}
