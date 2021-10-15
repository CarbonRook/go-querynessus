package querynessus

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"time"
)

type PluginRepository interface {
	Load() (*PluginListPage, error)
	Save(*PluginListPage) error
}

type JsonFilePluginRepository struct {
	filename string
}

func NewJsonFilePluginRepository(filename string) (*JsonFilePluginRepository, error) {
	return &JsonFilePluginRepository{
		filename: filename,
	}, nil
}

func (jfpr JsonFilePluginRepository) Load() (*PluginListPage, error) {
	jsonFile, err := os.Open(jfpr.filename)
	if err != nil {
		return &PluginListPage{}, err
	}
	results, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return &PluginListPage{}, err
	}
	var pluginPage PluginListPage
	err = json.Unmarshal(results, &pluginPage)
	if err != nil {
		return &PluginListPage{}, err
	}
	return &pluginPage, nil
}

func (jfpr JsonFilePluginRepository) Save(plugins *PluginListPage) error {
	file, err := json.Marshal(plugins)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(jfpr.filename, file, 0644)
	if err != nil {
		return err
	}
	return nil
}

var TenablePluginsServiceEndpoint = "https://cloud.tenable.com/plugins/plugin"
var TenableScannerGroupsEndpoint = "https://cloud.tenable.com/scanner-groups"
var TenableScanEndpoint = "https://cloud.tenable.com/scans"
var TenableFoldersEndpoint = "https://cloud.tenable.com/folders"
var RequestInterval = 3 * time.Second

type TenableRepository struct {
	requestInterval time.Duration
}

func (tr TenableRepository) Load() (*PluginListPage, error) {
	return &PluginListPage{}, nil
}

func (tr TenableRepository) Save(plugins *PluginListPage) error {
	return errors.New("cannot save to Tenable API")
}

func NewTenableRepository() (*TenableRepository, error) {
	return &TenableRepository{
		requestInterval: 3 * time.Second,
	}, nil
}
