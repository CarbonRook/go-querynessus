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

		fmt.Fprintf(os.Stderr, "\nEXAMPLES\n\nGet plugin information by name using jq:\njq '.data.plugin_details | .[] | select(.name | contains(\"QUERY\"))' plugins.json\n")
	}
	outfileArg := flag.String("out", "nessus-plugins.json", "The file to output the JSON to")
	allPluginsFlag := flag.Bool("all-plugins", false, "Fetch all plugins")
	pluginsSinceFlag := flag.String("plugins-since", "", "Fetch all plugins since YYYY-MM-DD")
	singlePluginFlag := flag.Int("single-plugin", 0, "The plugin ID to fetch")
	exportFlag := flag.Int("export-results", 0, "Export results from a given scan ID")
	allScansFlag := flag.Bool("list-scans", false, "Export all scans")
	allFoldersFlag := flag.Bool("list-folders", false, "List folders in your account")
	updateFileFlag := flag.String("update-plugins", "", "Add the latest plugins to a previously generated plugins file")
	flag.Parse()

	tac := querynessus.NewTenableApiClient(os.Getenv(TENABLE_ACCESS_KEY), os.Getenv(TENABLE_SECRET_KEY))

	if *allPluginsFlag || *pluginsSinceFlag != "" {
		params := querynessus.RequestParams{
			Size: 10000,
			Page: 1,
		}

		if *pluginsSinceFlag != "" {
			params.LastUpdated = *pluginsSinceFlag
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
			return
		}
		log.Fatalf("Successfully downloaded %s\n", outFile)
	} else if *allScansFlag {
		params := querynessus.ScanParams{}
		scanPage, err := tac.ListScans(&params)
		if err != nil {
			log.Fatal("Failed to get all scans")
			return
		}
		outFile := "scans.json"
		err = querynessus.SaveJsonToFile(outFile, scanPage)
		if err != nil {
			log.Fatalf("Failed to write scans page to file %s\n", outFile)
			return
		}
	} else if *allFoldersFlag {
		log.Printf("Fetching folder list")
		folderCollection, err := tac.ListFolders()
		if err != nil {
			log.Fatal("Failed to fetch list of folders\n")
		}
		for _, folder := range folderCollection.Folders {
			fmt.Printf("%s:%d\n", folder.Name, folder.Id)
		}
	} else if *updateFileFlag != "" {
		log.Printf("Updating file %s", *updateFileFlag)
		jfpr, err := querynessus.NewJsonFilePluginRepository(*updateFileFlag)
		if err != nil {
			log.Fatalf("Failed to create Json repository from file %s: %s\n", *updateFileFlag, err)
			os.Exit(1)
		}
		pluginPage, err := jfpr.Load()
		if err != nil {
			log.Fatalf("Failed to load plugin page from file %s: %s\n", *updateFileFlag, err)
			os.Exit(1)
		}
		log.Printf("Loaded %d plugins from %s", pluginPage.Size, *updateFileFlag)
		lastModifiedDate, err := pluginPage.LatestModifiedDate()
		if err != nil {
			log.Fatalf("Failed to get latest modified date for plugins: %s\n", err)
			os.Exit(1)
		}
		log.Printf("Latest plugin modification date %s\n", lastModifiedDate.Format(time.RFC3339))
		params := querynessus.RequestParams{
			Size:        10000,
			Page:        1,
			LastUpdated: lastModifiedDate.Format("2006-01-02"),
		}
		log.Printf("Fetching plugins since %s", lastModifiedDate.Format(time.RFC3339))
		results, err := tac.FetchAllPlugins(&params)
		if err != nil {
			log.Println("Failed to fetch plugins")
			os.Exit(1)
		}
		newPluginsPage := querynessus.PluginListPage{
			TotalCount: len(results),
			Data: querynessus.PluginDetailsList{
				PluginDetails: results,
			},
			Size: len(results),
		}
		log.Printf("Merging in %d new plugins\n", newPluginsPage.Size)
		pluginPage.Merge(&newPluginsPage)
		log.Printf("Saving new plugins to %s\n", *updateFileFlag)
		err = jfpr.Save(pluginPage)
		if err != nil {
			log.Printf("Failed to save to file %s: %s\n", *updateFileFlag, err)
			os.Exit(1)
		}
		log.Println("Complete")
	}
}
