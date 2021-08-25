package main

import (
	"go-querynessus/querynessus"
	"log"
	"os"
)

func main() {
	creds := querynessus.TenableCredentials{
		AccessKey: os.Getenv("TENABLE_ACCESS_KEY"),
		SecretKey: os.Getenv("TENABLE_SECRET_KEY"),
	}

	params := querynessus.RequestParams{
		Size: 10000,
		Page: 1,
	}

	results, err := querynessus.FetchAllPlugins(creds, &params)

	if err != nil {
		log.Println("Failed to fetch plugins")
		os.Exit(1)
	}

	combinedPage := querynessus.PluginListPage{
		TotalCount: len(results),
		Data: querynessus.PluginDetailsList{
			PluginDetails: results,
		},
		Size: len(results),
	}

	querynessus.SavePluginsToFile("plugins.json", combinedPage)
}
