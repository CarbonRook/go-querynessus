package querynessus

type ScansPage struct {
	FolderCollection
	Scans     []Scan `json:"scans"`
	Timestamp int    `json:"timestamp"`
}

type Scan struct {
	Legacy               bool   `json:"legacy"`
	Permissions          int    `json:"permissions"`
	Type                 string `json:"type"`
	Read                 bool   `json:"read"`
	LastModificationDate int    `json:"last_modification_date"`
	CreationDate         int    `json:"creation_date"`
	Status               string `json:"status"`
	UUID                 string `json:"uuid"`
	Shared               bool   `json:"shared"`
	UserPermissions      int    `json:"user_permissions"`
	Owner                string `json:"owner"`
	ScheduleUUID         string `json:"schedule_uuid"`
	Timezone             string `json:"timezone"`
	RepetitionRules      string `json:"rrules"`
	StartTime            string `json:"starttime"`
	Enabled              bool   `json:"enabled"`
	Control              bool   `json:"control"`
	WizardUUID           string `json:"wizard_uuid"`
	PolicyId             int    `json:"policy_id"`
	Name                 string `json:"name"`
	Id                   int    `json:"id"`
}

type ScanDetails struct {
	Info             ScanInfo        `json:"info"`
	History          []History       `json:"history"`
	Hosts            []Host          `json:"hosts"`
	Vulnerabilities  []Vulnerability `json:"vulnerabilities"`
	ComplianceHosts  []Host          `json:"comphosts"`
	ComplianceChecks []Compliance    `json:"compliance"`
	Notes            []Note          `json:"notes"`
}

type ScanInfo struct {
	Owner             string   `json:"owner"`
	Name              string   `json:"name"`
	NoTarget          bool     `json:"no_target"`
	FolderId          int      `json:"folder_id"`
	Control           bool     `json:"control"`
	UserPermissions   int      `json:"user_permissions"`
	ScheduleId        string   `json:"schedule_id"`
	EditAllowed       bool     `json:"edit_allowed"`
	ScannerName       string   `json:"scanner_name"`
	Policy            string   `json:"policy"`
	Shared            bool     `json:"shared"`
	ObjectId          int      `json:"object_id"`
	TagTargets        []string `json:"tag_targets"`
	ACLs              []ACL    `json:"acls"`
	HostCount         int      `json:"hostcount"`
	UUID              string   `json:"uuid"`
	Status            string   `json:"status"`
	ScanType          string   `json:"scan_type"`
	Targets           string   `json:"targets"`
	AltTargetsUsed    bool     `json:"alt_targets_used"`
	PCICanUpload      bool     `json:"pci_can_upload"`
	ScanStart         int      `json:"scan_start"`
	ScanEnd           int      `json:"scan_end"`
	Timestamp         int      `json:"timestamp"`
	IsArchived        bool     `json:"is_archived"`
	HasKB             bool     `json:"haskb"`
	HasAuditTrail     bool     `json:"hasaudittrail"`
	ImportedScanStart int      `json:"scanner_start"`
	ImportedScanEnd   int      `json:"scanner_end"`
}

type ACL struct {
	Permissions int    `json:"permissions"`
	Owner       string `json:"owner"`
	DisplayName string `json:"display_name"`
	Name        string `json:"name"`
	ID          int    `json:"id"`
	Type        string `json:"type"`
}

type History struct {
	HistoryID            int    `json:"history_id"`
	OwnerID              int    `json:"owner_id"`
	CreationDate         int    `json:"creation_date"`
	LastModificationDate int    `json:"last_modification_date"`
	UUID                 string `json:"uuid"`
	Type                 string `json:"type"`
	Status               string `json:"status"`
	Scheduler            int    `json:"scheduler"`
	AltTargetsUsed       bool   `json:"alt_targets_used"`
	IsArchived           bool   `json:"is_archived"`
}

type Host struct {
	AssetID               int    `json:"asset_id"`
	HostID                int    `json:"host_id"`
	Hostname              string `json:"hostname"`
	Progress              string `json:"progress"`
	ScanProgressCurrent   int    `json:"scanprogresscurrent"`
	ScanProgressTotal     int    `json:"scanprogresstotal"`
	NumChecksConsidered   int    `json:"numchecksconsidered"`
	TotalChecksConsidered int    `json:"totalchecksconsidered"`
	SeverityCount         []struct {
		Count         int `json:"count"`
		SeverityLevel int `json:"severitylevel"`
	} `json:"severitycount>item"`
	Severity  int `json:"severity"`
	Score     int `json:"score"`
	Info      int `json:"info"`
	Low       int `json:"low"`
	Medium    int `json:"medium"`
	High      int `json:"high"`
	Critical  int `json:"critical"`
	HostIndex int `json:"host_index"`
}

type Vulnerability struct {
	Count              int    `json:"count"`
	PluginID           int    `json:"plugin_id"`
	PluginName         string `json:"plugin_name"`
	Severity           int    `json:"severity"`
	PluginFamily       string `json:"plugin_family"`
	VulnerabilityIndex int    `json:"vuln_index"`
}

type Compliance struct {
	Count         int    `json:"count"`
	HostID        int    `json:"host_id"`
	Hostname      string `json:"hostname"`
	PluginID      int    `json:"plugin_id"`
	PluginName    string `json:"plugin_name"`
	Severity      int    `json:"severity"`
	PluginFamily  string `json:"plugin_family"`
	SeverityIndex int    `json:"severity_index"`
}

type Note struct {
	Message string `json:"message"`
	Title   string `json:"title"`
}
