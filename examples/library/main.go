package main

import (
	"fmt"
	"os"

	"github.com/Will-thom/pr-maven-cli/pkg/prmaven"
)

func main() {
	projectDir := "demo/multi-module-failure"
	if len(os.Args) > 1 {
		projectDir = os.Args[1]
	}

	report, err := prmaven.Analyze(prmaven.Options{ProjectDir: projectDir})
	if err != nil {
		fmt.Fprintf(os.Stderr, "analyze Maven reports: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Modules: %d | Reports: %d | Findings: %d\n",
		report.Summary.ModuleCount,
		report.Summary.ReportCount,
		report.Summary.FindingCount,
	)

	for _, finding := range report.Findings {
		fmt.Printf("%s: %s\n", finding.ID, finding.ReproduceCommand)
	}
}
