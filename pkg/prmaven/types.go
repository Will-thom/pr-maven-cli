package prmaven

// Report is the stable library contract returned by the analyzer.
type Report struct {
	ProjectRoot string    `json:"projectRoot"`
	Summary     Summary   `json:"summary"`
	Modules     []Module  `json:"modules"`
	Findings    []Finding `json:"findings"`
}

type Summary struct {
	ModuleCount  int `json:"moduleCount"`
	ReportCount  int `json:"reportCount"`
	FindingCount int `json:"findingCount"`
}

type Module struct {
	Name string `json:"name"`
	Path string `json:"path"`
	POM  string `json:"pom"`
}

type Finding struct {
	ID                 string   `json:"id"`
	Module             string   `json:"module"`
	ModulePath         string   `json:"modulePath"`
	ReportPath         string   `json:"reportPath"`
	ReportKind         string   `json:"reportKind"`
	MavenPlugin        string   `json:"mavenPlugin"`
	MavenPhase         string   `json:"mavenPhase"`
	TestClass          string   `json:"testClass"`
	TestName           string   `json:"testName"`
	FailureKind        string   `json:"failureKind"`
	FailureType        string   `json:"failureType,omitempty"`
	Message            string   `json:"message,omitempty"`
	ReproduceCommand   string   `json:"reproduceCommand"`
	Confidence         string   `json:"confidence"`
	ConfidenceReasons  []string `json:"confidenceReasons"`
	SourceReportFormat string   `json:"sourceReportFormat"`
}
