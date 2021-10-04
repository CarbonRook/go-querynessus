package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

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
	allPluginsFlag := flag.Bool("all-plugins", false, "Fetch all plugins")
	singlePluginFlag := flag.Int("single-plugin", 0, "The plugin ID to fetch")
	exportFlag := flag.Int("export-results", 0, "Export results from a given scan ID")
	flag.Parse()

	tac := querynessus.NewTenableApiClient(os.Getenv(TENABLE_ACCESS_KEY), os.Getenv(TENABLE_SECRET_KEY))

	if *allPluginsFlag {
		params := querynessus.RequestParams{
			Size: 10000,
			Page: 1,
		}

		results, err := tac.FetchAllPlugins(&params)

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
	} else if *singlePluginFlag > 0 {

	} else if *exportFlag != 0 {
		params := querynessus.RequestParams{}
		log.Printf("Submitting export task to Tenable for scan %d\n", *exportFlag)
		fileId, _, err := tac.ExportScanResults(&params, *exportFlag)
		if err != nil {
			log.Fatalf("Failed to start export of scan %d: %s\n", *exportFlag, err)
			return
		}
		for {
			time.Sleep(time.Second * 3)
			log.Printf("Checking status for scan %d and file %s\n", *exportFlag, fileId)
			isReady, err := tac.ScanResultExportStatus(*exportFlag, fileId)
			if err != nil {
				log.Fatalf("Failed to get status for scan %d and file %s: %s\n", *exportFlag, fileId, err)
				return
			}
			if isReady {
				break
			}
		}
		outFile := fmt.Sprintf("%d-%s.nessus", *exportFlag, fileId)
		err = tac.DownloadExportedScan(*exportFlag, fileId, outFile)
		if err != nil {
			log.Fatalf("Failed to write to %s: %s\n", outFile, err)
		}
		log.Fatalf("Successfully downloaded %s\n", outFile)
	}
}
