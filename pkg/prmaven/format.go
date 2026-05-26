package prmaven

import (
	"encoding/json"
	"fmt"
	"io"
)

func WriteJSON(w io.Writer, report Report) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}

func WriteText(w io.Writer, report Report) error {
	_, err := fmt.Fprintf(w, "PR Maven CLI - Maven failure context\n\n")
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "Project: %s\n", report.ProjectRoot); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "Modules: %d | Reports: %d | Findings: %d\n\n", report.Summary.ModuleCount, report.Summary.ReportCount, report.Summary.FindingCount); err != nil {
		return err
	}
	if len(report.Findings) == 0 {
		_, err := fmt.Fprintln(w, "No Maven test failures found in Surefire or Failsafe reports.")
		return err
	}

	for i, finding := range report.Findings {
		if _, err := fmt.Fprintf(w, "%d. %s\n", i+1, finding.ID); err != nil {
			return err
		}
		lines := []string{
			"Module: " + finding.Module + " (" + finding.ModulePath + ")",
			"Report: " + finding.ReportPath,
			"Plugin: " + finding.MavenPlugin,
			"Phase: " + finding.MavenPhase,
			"Test: " + finding.TestClass + "." + finding.TestName,
			"Kind: " + finding.FailureKind,
		}
		if finding.Message != "" {
			lines = append(lines, "Message: "+finding.Message)
		}
		lines = append(lines, "Reproduce: "+finding.ReproduceCommand)
		lines = append(lines, "Confidence: "+finding.Confidence)
		for _, line := range lines {
			if _, err := fmt.Fprintf(w, "   %s\n", line); err != nil {
				return err
			}
		}
		if len(finding.ConfidenceReasons) > 0 {
			if _, err := fmt.Fprintln(w, "   Reasons:"); err != nil {
				return err
			}
			for _, reason := range finding.ConfidenceReasons {
				if _, err := fmt.Fprintf(w, "   - %s\n", reason); err != nil {
					return err
				}
			}
		}
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
	}
	return nil
}
