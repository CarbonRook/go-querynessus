package querynessus

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/go-querystring/query"
)

var TenablePluginsServiceEndpoint = "https://cloud.tenable.com/plugins/plugin"
var RequestInterval = 3 * time.Second

type TenableCredentials struct {
	AccessKey string
	SecretKey string
}

type RequestParams struct {
	LastUpdated string `url:"last_updated,omitempty"`
	Size        int32  `url:"size"`
	Page        int32  `url:"page"`
}

func fetchSinglePluginPage(creds TenableCredentials, params *RequestParams) (PluginListPage, error) {

	v, _ := query.Values(params)
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", TenablePluginsServiceEndpoint, nil)
	req.URL.RawQuery = v.Encode()
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-ApiKeys", "accessKey="+creds.AccessKey+";secretKey="+creds.SecretKey)
	if err != nil {
		log.Fatalf("Failed to create request")
		return PluginListPage{}, err
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("Failed to submit request")
		return PluginListPage{}, err
	}
	if resp.StatusCode != 200 {
		log.Printf("Received %s response from %s", resp.Status, TenablePluginsServiceEndpoint)
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	var pluginPage PluginListPage
	err = decoder.Decode(&pluginPage)
	if err != nil {
		log.Fatalf("Failed to decode results")
		return PluginListPage{}, err
	}
	return pluginPage, nil
}

func FetchPlugins(creds TenableCredentials, params *RequestParams) ([]PluginDetails, error) {
	pluginPage, err := fetchSinglePluginPage(creds, params)
	if err != nil {
		log.Println("Failed to fetch plugin page")
		return nil, err
	}
	return pluginPage.Data.PluginDetails, err
}

func FetchAllPlugins(creds TenableCredentials, params *RequestParams) ([]PluginDetails, error) {
	var pluginDetails []PluginDetails
	for {
		log.Printf("Requesting plugin page %d", params.Page)
		endIndex := params.Page * params.Size
		pluginPage, err := fetchSinglePluginPage(creds, params)
		if err != nil {
			log.Printf("Failed to fetch plugin page %d", params.Page)
			continue
		}
		pluginDetails = append(pluginDetails, pluginPage.Data.PluginDetails...)
		if endIndex > int32(pluginPage.TotalCount) {
			break
		} else {
			params.Page = params.Page + 1
		}
		time.Sleep(RequestInterval)
	}
	return pluginDetails, nil
}

func LoadPluginsFromFile(filename string) (PluginListPage, error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Println("Failed to open json file")
		return PluginListPage{}, nil
	}
	results, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Println("Failed to read json file")
		return PluginListPage{}, nil
	}
	var pluginPage PluginListPage
	err = json.Unmarshal(results, &pluginPage)
	if err != nil {
		log.Println("Failed to unmarshal json")
		log.Println(err)
		return PluginListPage{}, err
	}
	return pluginPage, nil
}

func SavePluginsToFile(filename string, pluginsPage PluginListPage) error {
	file, err := json.Marshal(pluginsPage)
	if err != nil {
		log.Println("Failed to marshal JSON structure")
		return err
	}
	err = ioutil.WriteFile(filename, file, 0644)
	if err != nil {
		log.Printf("Failed to write to file %s", filename)
		return err
	}
	return nil
}
