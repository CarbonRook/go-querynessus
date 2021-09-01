package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/carbonrook/go-querynessus/querynessus"
)

var TENABLE_ACCESS_KEY = "TENABLE_ACCESS_KEY"
var TENABLE_SECRET_KEY = "TENABLE_SECRET_KEY"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])

		flag.PrintDefaults()

		fmt.Fprintf(os.Stderr, "\nRequired environment vars:\n%s: Tenable API access key\n%s: Tenable API secret key\n", TENABLE_ACCESS_KEY, TENABLE_SECRET_KEY)
	}
	outfileArg := flag.String("out", "nessus-plugins.json", "The file to output the JSON to")
	flag.Parse()

	creds := querynessus.TenableCredentials{
		AccessKey: os.Getenv(TENABLE_ACCESS_KEY),
		SecretKey: os.Getenv(TENABLE_SECRET_KEY),
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

	combinedPage.SaveToFile(*outfileArg)
}
