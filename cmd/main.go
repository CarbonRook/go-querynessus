package main

import (
	"fmt"
	"go-querynessus/querynessus"
	"log"
	"os"
)

func main() {
	creds := querynessus.TenableCredentials{
		AccessKey: os.Getenv("TENABLE_ACCESS_KEY"),
		SecretKey: os.Getenv("TENABLE_SECRET_KEY"),
	}

	params := querynessus.RequestParams{}

	results, err := querynessus.FetchPlugins(creds, &params)

	if err != nil {
		log.Println("Failed to fetch plugins")
		os.Exit(1)
	}

	for _, nessusPlugin := range results {
		fmt.Println(nessusPlugin.Name)
	}
}
