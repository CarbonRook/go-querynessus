package querynessus

type PluginListPage struct {
	Size       int                  `json:"size"`
	TotalCount int                  `json:"total_count"`
	Params     PluginListPageParams `json:"params"`
	Data       PluginDetailsList    `json:"data"`
}

type PluginListPageParams struct {
	Page        int    `json:"page"`
	Size        int    `json:"size"`
	LastUpdated string `json:"last_updated"`
}

type PluginDetailsList struct {
	PluginDetails []PluginDetails `json:"plugin_details"`
}

type PluginDetails struct {
	ID         int              `json:"id"`
	Name       string           `json:"name"`
	Attributes PluginAttributes `json:"attributes"`
}

type PluginAttributes struct {
	PluginModificationDate       string                      `json:"plugin_modification_date"`
	IntelType                    string                      `json:"intel_type"`
	PluginPublicationDate        string                      `json:"plugin_publication_date"`
	VulnerabilityPublicationDate string                      `json:"vuln_publication_date,omitempty"`
	Version                      string                      `json:"plugin_version,omitempty"`
	PluginType                   string                      `json:"plugin_type,omitempty"`
	Description                  string                      `json:"description"`
	RiskFactor                   string                      `json:"risk_factor"`
	ExploitedByNessus            bool                        `json:"exploited_by_nessus,omitempty"`
	CVE                          []string                    `json:"cve"`
	DefaultAccount               string                      `json:"default_account,omitempty"`
	Solution                     string                      `json:"solution"`
	CPE                          []string                    `json:"cpe,omitempty"`
	InTheNews                    bool                        `json:"in_the_news,omitempty"`
	Synopsis                     string                      `json:"synopsis"`
	VPR                          VulnerabilityPriorityRating `json:"vpr"`
	AlwaysRun                    bool                        `json:"always_run"`
	Compliance                   bool                        `json:"compliance"`
	BID                          []int                       `json:"bid"`

	CVSSv2BaseScore     float32 `json:"cvss_base_score,omitempty"`
	CVSSv2TemporalScore float32 `json:"cvss_temporal_score,omitempty"`
	CVSSv2Vector        struct {
		VectorString          string `json:"raw"`
		AccessVector          string `json:"AccessVector"`
		AccessComplexity      string `json:"AccessComplexity"`
		Authentication        string `json:"Authentication"`
		ConfidentialityImpact string `json:"Confidentiality-Impact"`
		IntegrityImpact       string `json:"Integrity-Impact"`
		AvailabilityImpact    string `json:"Availability-Impact"`
	} `json:"cvss_vector,omitempty"`
	CVSSv2TemporalVector struct {
		VectorString     string `json:"raw"`
		Exploitability   string `json:"Exploitability"`
		RemediationLevel string `json:"RemediationLevel"`
		ReportConfidence string `json:"ReportConfidence"`
	} `json:"cvss_temporal_vector,omitempty"`
	CVSSv3BaseScore     float32 `json:"cvss3_base_score,omitempty"`
	CVSSv3TemporalScore float32 `json:"cvss3_temporal_score,omitempty"`
	CVSSv3ImpactScore   float32 `json:"cvss3_impact_score,omitempty"`
	CVSSv3Vector        struct {
		VectorString          string `json:"raw"`
		AttackVector          string `json:"AttackVector"`
		AttackComplexity      string `json:"AttackComplexity"`
		PrivilegesRequired    string `json:"PrivilegesRequired"`
		UserInteraction       string `json:"UserInteraction"`
		Scope                 string `json:"Scope"`
		ConfidentialityImpact string `json:"Confidentiality-Impact"`
		IntegrityImpact       string `json:"Integrity-Impact"`
		AvailabilityImpact    string `json:"Availability-Impact"`
	} `json:"cvss3_vector,omitempty"`
	CVSSv3TemporalVector struct {
		VectorString        string `json:"raw"`
		ExploitCodeMaturity string `json:"ExploitCodeMaturity"`
		RemediationLevel    string `json:"RemediationLevel"`
		ReportConfidence    string `json:"ReportConfidence"`
	} `json:"cvss3_temporal_vector,omitempty"`

	ExploitAvailable           bool `json:"exploit_available"`
	ExploitFrameworkCanvas     bool `json:"exploit_framework_canvas,omitempty"`
	ExploitFrameworkCore       bool `json:"exploit_framework_core,omitempty"`
	ExploitFrameworkD2Elliot   bool `json:"exploit_framework_d2_elliot,omitempty"`
	ExploitFrameworkExploitHub bool `json:"exploit_framework_exploithub,omitempty"`
	ExploitFrameworkMetasploit bool `json:"exploit_framework_metasploit,omitempty"`

	ExploitedByMalware bool `json:"exploited_by_malware,omitempty"`
	Malware            bool `json:"malware,omitempty"`

	HasPatch             bool   `json:"has_patch"`
	PatchPublicationDate string `json:"patch_publication_date,omitempty"`
	UnsupportedByVendor  bool   `json:"unsupported_by_vendor,omitempty"`

	SeeAlso []string `json:"see_also"`
	XRef    []string `json:"xref"`
	XRefs   []struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"xrefs"`
}

type VulnerabilityPriorityRating struct {
	Score   float32    `json:"score,omitempty"`
	Drivers VPRDrivers `json:"drivers,omitempty"`
	Updated string     `json:"updated,omitempty"`
}

type VPRDrivers struct {
	AgeOfVulnerability struct {
		LowerBound int `json:"lower_bound,omitempty"`
		UpperBound int `json:"upper_bound,omitempty"`
	} `json:"age_of_vulnerability"`
	ExploitCodeMaturity        string `json:"exploit_code_maturity"`
	ThreatIntensityLast28      string `json:"threat_intensity_last28"`
	IsCVSSImpactScorePredicted bool   `json:"cvss_impact_score_predicted,omitempty"`
	ThreatRecency              struct {
		LowerBound int `json:"lower_bound,omitempty"`
		UpperBound int `json:"upper_bound,omitempty"`
	} `json:"threat_recency"`
	ThreatSourcesLast28 []string `json:"threat_sources_last28"`
	ProductCoverage     string   `json:"product_coverage"`
}
