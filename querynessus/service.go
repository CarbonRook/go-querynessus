package querynessus

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

var TenablePluginsServiceEndpoint = "https://cloud.tenable.com/plugins/plugin"
var TenableScannerGroupsEndpoint = "https://cloud.tenable.com/scanner-groups"
var TenableScanEndpoint = "https://cloud.tenable.com/scans"
var RequestInterval = 3 * time.Second

type TenableApiClient struct {
	Credentials TenableCredentials
}

type TenableCredentials struct {
	AccessKey string
	SecretKey string
}

func NewTenableApiClient(accessKey string, secretKey string) TenableApiClient {
	return TenableApiClient{
		Credentials: TenableCredentials{
			AccessKey: accessKey,
			SecretKey: secretKey,
		},
	}
}

type RequestParams struct {
	LastUpdated string `url:"last_updated,omitempty"`
	Size        int32  `url:"size"`
	Page        int32  `url:"page"`
}

func (tac TenableApiClient) sendPostRequest(tenableEndpoint string, params *RequestParams, payload string) (*http.Response, error) {
	v, _ := query.Values(params)
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("POST", tenableEndpoint, strings.NewReader(payload))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-ApiKeys", "accessKey="+tac.Credentials.AccessKey+";secretKey="+tac.Credentials.SecretKey)
	req.URL.RawQuery = v.Encode()
	if err != nil {
		log.Fatalf("Failed to create request")
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to submit request")
		return nil, err
	}
	if resp.StatusCode != 200 {
		log.Printf("Received %s response from %s", resp.Status, tenableEndpoint)
	}

	return resp, nil
}

func (tac TenableApiClient) sendGetRequest(tenableEndpoint string, params *RequestParams) (*http.Response, error) {
	v, _ := query.Values(params)
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", tenableEndpoint, nil)
	req.URL.RawQuery = v.Encode()
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-ApiKeys", "accessKey="+tac.Credentials.AccessKey+";secretKey="+tac.Credentials.SecretKey)
	if err != nil {
		log.Fatalf("Failed to create request")
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to submit request")
		return nil, err
	}
	if resp.StatusCode != 200 {
		log.Printf("Received %s response from %s", resp.Status, tenableEndpoint)
	}

	return resp, nil
}

func (tac TenableApiClient) ExportScanResults(params *RequestParams, scanId int) (fileId string, tempToken string, err error) {
	endpoint := fmt.Sprintf("%s/%d/export", TenableScanEndpoint, scanId)
	payload := "{\"format\": \"nessus\"}"
	resp, err := tac.sendPostRequest(endpoint, params, payload)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("received non 200 response %d", resp.StatusCode)
	}

	type ExportResponseBody struct {
		FileId    string `json:"file"`
		TempToken string `json:"temp_token"`
	}

	decoder := json.NewDecoder(resp.Body)
	var respBody ExportResponseBody
	err = decoder.Decode(&respBody)
	if err != nil {
		return "", "", err
	}

	return respBody.FileId, respBody.TempToken, nil
}

func (tac TenableApiClient) ScanResultExportStatus(scanId int, fileId string) (result bool, err error) {
	endpoint := fmt.Sprintf("%s/%d/export/%s/status", TenableScanEndpoint, scanId, fileId)
	resp, err := tac.sendGetRequest(endpoint, &RequestParams{})
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false, fmt.Errorf("received non 200 response %d", resp.StatusCode)
	}

	type StatusResponseBody struct {
		Status string `json:"status"`
	}

	decoder := json.NewDecoder(resp.Body)
	var respBody StatusResponseBody
	err = decoder.Decode(&respBody)
	if err != nil {
		return false, err
	}
	return strings.ToLower(respBody.Status) == "ready", nil
}

func (tac TenableApiClient) DownloadExportedScan(scanId int, fileId string, outFile string) error {
	out, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer out.Close()

	endpoint := fmt.Sprintf("%s/%d/export/%s/download", TenableScanEndpoint, scanId, fileId)
	resp, err := tac.sendGetRequest(endpoint, &RequestParams{})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("received non 200 response %d", resp.StatusCode)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (tac TenableApiClient) fetchSinglePluginPage(params *RequestParams) (*PluginListPage, error) {
	resp, err := tac.sendGetRequest(TenablePluginsServiceEndpoint, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	var pluginPage PluginListPage
	err = decoder.Decode(&pluginPage)
	if err != nil {
		log.Fatalf("Failed to decode results")
		return nil, err
	}
	return &pluginPage, nil
}

func (tac TenableApiClient) FetchPlugins(params *RequestParams) ([]PluginDetails, error) {
	pluginPage, err := tac.fetchSinglePluginPage(params)
	if err != nil {
		log.Println("Failed to fetch plugin page")
		return nil, err
	}
	return pluginPage.Data.PluginDetails, err
}

func (tac TenableApiClient) FetchAllPlugins(params *RequestParams) ([]PluginDetails, error) {
	var pluginDetails []PluginDetails
	for {
		log.Printf("Requesting plugin page %d", params.Page)
		endIndex := params.Page * params.Size
		pluginPage, err := tac.fetchSinglePluginPage(params)
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
