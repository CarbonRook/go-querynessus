package querynessus

type ScansPage struct {
	Folders   []Folder `json:"folders"`
	Scans     []Scan   `json:"scans"`
	Timestamp int      `json:"timestamp"`
}

type Folder struct {
	UnreadCount int    `json:"unread_count"`
	Custom      int    `json:"custom"`
	DefaultTag  int    `json:"default_tag"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Id          int    `json:"id"`
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
