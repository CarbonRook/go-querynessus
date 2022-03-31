package main

import (
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/CarbonRook/go-querynessus/querynessus"
)

var TENABLE_ACCESS_KEY = "TENABLE_ACCESS_KEY"
var TENABLE_SECRET_KEY = "TENABLE_SECRET_KEY"

const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func random(length int) (string, error) {
	bytes := make([]byte, length)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = chars[b%byte(len(chars))]
	}

	return string(bytes), nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])

		flag.PrintDefaults()

		fmt.Fprintf(os.Stderr, "\nRequired environment vars:\n%s: Tenable API access key\n%s: Tenable API secret key\n", TENABLE_ACCESS_KEY, TENABLE_SECRET_KEY)

		fmt.Fprintf(os.Stderr, "\nEXAMPLES\n\nGet plugin information by name using jq:\njq '.data.plugin_details | .[] | select(.name | contains(\"QUERY\"))' plugins.json\n")
	}
	outfileArg := flag.String("out", "nessus-plugins.json", "The file to output the JSON to")
	// Plugins
	allPluginsFlag := flag.Bool("all-plugins", false, "Fetch all plugins")
	pluginsSinceFlag := flag.String("plugins-since", "", "Fetch all plugins since YYYY-MM-DD")
	singlePluginFlag := flag.Int("single-plugin", 0, "The plugin ID to fetch")
	// Scan export
	exportFlag := flag.Int("export-results", 0, "Export results from a given scan ID")
	exportFormatFlag := flag.String("format", "nessus", "Export the results in a given format \"nessus\", \"db\", \"csv\"")
	// Scans
	allScansFlag := flag.Bool("list-scans", false, "Export all scans")
	scansSinceFlag := flag.String("scans-since", "", "Fetch all scans since a given date, YYYY-MM-DD")
	singleScanFlag := flag.Int("scan", 0, "Fetch single scan details")
	// Scan Filters
	//customerFlag := flag.String("customer", "", "Customer name to filter scans for")
	//scanTypeFlag := flag.String("scan-type", "", "Type of scan to filter on")
	// Folders
	allFoldersFlag := flag.Bool("list-folders", false, "List folders in your account")
	// Update existing JSON database
	updateFileFlag := flag.String("update-plugins", "", "Add the latest plugins to a previously generated plugins file")
	flag.Parse()

	permittedFormats := map[string]bool{"nessus": true, "db": true, "csv": true}
	_, ok := permittedFormats[*exportFormatFlag]
	if !ok {
		log.Fatalf("Invalid export format provided: %s", *exportFormatFlag)
		return
	}

	tac := querynessus.NewTenableApiClient(os.Getenv(TENABLE_ACCESS_KEY), os.Getenv(TENABLE_SECRET_KEY))

	if *allPluginsFlag || *pluginsSinceFlag != "" {
		FetchAllPlugins(&tac, pluginsSinceFlag, outfileArg)
	} else if *singlePluginFlag > 0 {
		FetchSinglePlugin(&tac, singlePluginFlag)
	} else if *exportFlag != 0 {
		ExportScan(&tac, exportFlag, exportFormatFlag)
	} else if *allScansFlag || *scansSinceFlag != "" {
		FetchAllScans(&tac, scansSinceFlag)
	} else if *allFoldersFlag {
		FetchAllFolders(&tac)
	} else if *updateFileFlag != "" {
		UpdatePluginRepository(&tac, updateFileFlag)
	} else if *singleScanFlag > 0 {
		FetchSingleScan(&tac, singleScanFlag)
	}
}

func FetchSinglePlugin(tac *querynessus.TenableApiClient, pluginId *int) {
	result, err := tac.FetchSinglePluginDetails(*pluginId)
	if err != nil {
		log.Printf("Failed to fetch plugin id %d: %s\n", *pluginId, err)
		os.Exit(1)
	}
	pluginDetailsJson, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Printf("Failed to create json for plugin %d: %s\n", *pluginId, err)
		os.Exit(1)
	}
	fmt.Println(pluginDetailsJson)
}

func FetchAllPlugins(tac *querynessus.TenableApiClient, pluginsSinceFlag *string, outFilePath *string) {
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

	combinedPage.SaveToFile(*outFilePath)
}

func UpdatePluginRepository(tac *querynessus.TenableApiClient, filePath *string) {
	log.Printf("Updating file %s", *filePath)
	jfpr, err := querynessus.NewJsonFilePluginRepository(*filePath)
	if err != nil {
		log.Fatalf("Failed to create Json repository from file %s: %s\n", *filePath, err)
		os.Exit(1)
	}
	pluginPage, err := jfpr.Load()
	if err != nil {
		log.Fatalf("Failed to load plugin page from file %s: %s\n", *filePath, err)
		os.Exit(1)
	}
	log.Printf("Loaded %d plugins from %s", pluginPage.Size, *filePath)
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
	newCount, updatedCount, duplicateCount, err := pluginPage.Merge(&newPluginsPage)
	if err != nil {
		log.Fatalf("Failed to merge plugins: %s\n", err)
		os.Exit(1)
	}
	log.Printf("Merged %d new plugins, updated %d existing plugins, ignored %d duplicate plugins", newCount, updatedCount, duplicateCount)
	log.Printf("Saving new plugins to %s\n", *filePath)
	err = jfpr.Save(pluginPage)
	if err != nil {
		log.Printf("Failed to save to file %s: %s\n", *filePath, err)
		os.Exit(1)
	}
	log.Println("Complete")
}

func FetchAllFolders(tac *querynessus.TenableApiClient) {
	log.Printf("Fetching folder list")
	folderCollection, err := tac.ListFolders()
	if err != nil {
		log.Fatal("Failed to fetch list of folders\n")
	}
	for _, folder := range folderCollection.Folders {
		fmt.Printf("%s:%d\n", folder.Name, folder.Id)
	}
}

func FetchSingleScan(tac *querynessus.TenableApiClient, scanId *int) {
	result, err := tac.FetchScanDetails(*scanId)
	if err != nil {
		log.Printf("Failed to fetch scan id %d: %s\n", *scanId, err)
		os.Exit(1)
	}
	scanDetailsJson, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Printf("Failed to create json for plugin %d: %s\n", *scanId, err)
		os.Exit(1)
	}
	fmt.Println(string(scanDetailsJson))
}

func FetchAllScans(tac *querynessus.TenableApiClient, since *string) {
	params := querynessus.ScanParams{}

	if *since != "" {
		timestamp, err := time.Parse("2006-01-02", *since)
		if err != nil {
			log.Fatalf("Failed to parse provided date %s: %s", *since, err)
			return
		}
		params.EarliestStartDate = int(timestamp.Unix())
	}

	scanPage, err := tac.ListScans(&params)
	if err != nil {
		log.Fatalf("Failed to get all scans: %s", err)
		return
	}
	outFile := "scans.json"
	err = querynessus.SaveJsonToFile(outFile, scanPage)
	if err != nil {
		log.Fatalf("Failed to write scans page to file %s\n", outFile)
		return
	}
}

func ExportScan(tac *querynessus.TenableApiClient, scanId *int, format *string) {
	params := querynessus.ExportScanParams{}
	payload := querynessus.ExportScanPayload{
		Format: *format,
	}
	if *format == "db" {
		scanDetails, err := tac.FetchScanDetails(*scanId)
		if err != nil {
			log.Printf("Failed to fetch scan id %d: %s\n", *scanId, err)
			os.Exit(1)
		}
		log.Printf("Found history UUID for scan %d: %s", *scanId, scanDetails.History[0].UUID)
		params.HistoryID = scanDetails.Info.UUID

		password, err := random(12)
		if err != nil {
			log.Fatalf("Failed to generate DB password: %s", err)
			return
		}
		log.Printf("Database password set: %s", password)
		payload.Password = password
		payload.AssetID = scanDetails.Hosts[0].AssetID
	}
	log.Printf("Submitting export task to Tenable for scan %d\n", *scanId)
	fileId, _, err := tac.ExportScanResults(&params, *scanId, &payload)
	if err != nil {
		log.Fatalf("Failed to start export of scan %d: %s\n", *scanId, err)
		return
	}
	for {
		time.Sleep(time.Second * 3)
		log.Printf("Checking status for scan %d and file %s\n", *scanId, fileId)
		isReady, err := tac.ScanResultExportStatus(*scanId, fileId)
		if err != nil {
			log.Fatalf("Failed to get status for scan %d and file %s: %s\n", *scanId, fileId, err)
			return
		}
		if isReady {
			break
		}
	}
	outFile := fmt.Sprintf("%d-%s.%s", *scanId, fileId, payload.Format)
	err = tac.DownloadExportedScan(*scanId, fileId, outFile)
	if err != nil {
		log.Fatalf("Failed to write to %s: %s\n", outFile, err)
		return
	}
	log.Fatalf("Successfully downloaded %s\n", outFile)

}
