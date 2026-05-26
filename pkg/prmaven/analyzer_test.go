package prmaven

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestAnalyzeDemoProject(t *testing.T) {
	report, err := Analyze(Options{ProjectDir: "../../demo/multi-module-failure"})
	if err != nil {
		t.Fatal(err)
	}

	if report.Summary.ModuleCount != 3 {
		t.Fatalf("module count = %d, want 3", report.Summary.ModuleCount)
	}
	if report.Summary.ReportCount != 2 {
		t.Fatalf("report count = %d, want 2", report.Summary.ReportCount)
	}
	if report.Summary.FindingCount != 2 {
		t.Fatalf("finding count = %d, want 2", report.Summary.FindingCount)
	}

	first := report.Findings[0]
	if first.Module != "payment-api" {
		t.Fatalf("first module = %q, want payment-api", first.Module)
	}
	if first.MavenPlugin != "maven-failsafe-plugin" {
		t.Fatalf("first plugin = %q, want maven-failsafe-plugin", first.MavenPlugin)
	}
	if first.ReproduceCommand != "mvn -pl payment-api -am -Dit.test=PaymentApiIT verify" {
		t.Fatalf("first reproduce command = %q", first.ReproduceCommand)
	}

	second := report.Findings[1]
	if second.Module != "payment-core" {
		t.Fatalf("second module = %q, want payment-core", second.Module)
	}
	if second.MavenPlugin != "maven-surefire-plugin" {
		t.Fatalf("second plugin = %q, want maven-surefire-plugin", second.MavenPlugin)
	}
	if second.ReproduceCommand != "mvn -pl payment-core -am -Dtest=PaymentRoundingTest test" {
		t.Fatalf("second reproduce command = %q", second.ReproduceCommand)
	}
}

func TestWriteTextIncludesActionableContext(t *testing.T) {
	report, err := Analyze(Options{ProjectDir: "../../demo/multi-module-failure"})
	if err != nil {
		t.Fatal(err)
	}

	var output bytes.Buffer
	if err := WriteText(&output, report); err != nil {
		t.Fatal(err)
	}

	text := output.String()
	for _, expected := range []string{
		"Module: payment-core (payment-core)",
		"Plugin: maven-surefire-plugin",
		"Reproduce: mvn -pl payment-core -am -Dtest=PaymentRoundingTest test",
		"Confidence: high",
	} {
		if !strings.Contains(text, expected) {
			t.Fatalf("text output missing %q\n%s", expected, text)
		}
	}
}

func TestWriteJSONProducesStableContract(t *testing.T) {
	report, err := Analyze(Options{ProjectDir: "../../demo/multi-module-failure"})
	if err != nil {
		t.Fatal(err)
	}

	var output bytes.Buffer
	if err := WriteJSON(&output, report); err != nil {
		t.Fatal(err)
	}

	var decoded Report
	if err := json.Unmarshal(output.Bytes(), &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Summary.FindingCount != 2 {
		t.Fatalf("decoded finding count = %d, want 2", decoded.Summary.FindingCount)
	}
	if decoded.Findings[0].SourceReportFormat != "junit-xml" {
		t.Fatalf("source format = %q, want junit-xml", decoded.Findings[0].SourceReportFormat)
	}
}
