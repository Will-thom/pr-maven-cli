package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLIEndToEndText(t *testing.T) {
	command := exec.Command("go", "run", ".", "fails", "-project", "../../demo/multi-module-failure")
	output, err := command.CombinedOutput()
	if err == nil {
		t.Fatal("CLI exit code = 0, want non-zero when findings are present")
	}

	text := string(output)
	for _, expected := range []string{
		"PR Maven CLI - Maven failure context",
		"Module: payment-core (payment-core)",
		"Reproduce: mvn -pl payment-core -am -Dtest=PaymentRoundingTest test",
	} {
		if !strings.Contains(text, expected) {
			t.Fatalf("CLI output missing %q\n%s", expected, text)
		}
	}
}

func TestRunReturnsFindingExitCodeAndText(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{"fails", "-project", "../../demo/multi-module-failure"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("exit code = %d, want 1", code)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	if !strings.Contains(stdout.String(), "Findings: 2") {
		t.Fatalf("stdout missing findings summary\n%s", stdout.String())
	}
}

func TestRunReturnsZeroForNoFailureProject(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{"fails", "-project", "../../demo/no-failure"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0\nstderr=%s\nstdout=%s", code, stderr.String(), stdout.String())
	}
	if !strings.Contains(stdout.String(), "Findings: 0") {
		t.Fatalf("stdout missing no-failure summary\n%s", stdout.String())
	}
}

func TestRunWritesTextOutputFile(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	outputPath := filepath.Join(t.TempDir(), "prmaven-report.txt")

	code := run([]string{"fails", "-project", "../../demo/multi-module-failure", "-output", outputPath}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("exit code = %d, want 1", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty when output file is set", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}

	output, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	text := string(output)
	for _, expected := range []string{
		"PR Maven CLI - Maven failure context",
		"Findings: 2",
		"maven-surefire-plugin",
	} {
		if !strings.Contains(text, expected) {
			t.Fatalf("output file missing %q\n%s", expected, text)
		}
	}
}

func TestRunWritesJSONOutputFile(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	outputPath := filepath.Join(t.TempDir(), "prmaven-report.json")

	code := run([]string{"why", "-project", "../../demo/no-failure", "-format", "json", "-output", outputPath}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0\nstderr=%s\nstdout=%s", code, stderr.String(), stdout.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty when output file is set", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}

	output, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}

	var report struct {
		Summary struct {
			ModuleCount  int `json:"moduleCount"`
			ReportCount  int `json:"reportCount"`
			FindingCount int `json:"findingCount"`
		} `json:"summary"`
		Findings []struct{} `json:"findings"`
	}
	if err := json.Unmarshal(output, &report); err != nil {
		t.Fatalf("output file contains invalid JSON: %v\n%s", err, string(output))
	}
	if report.Summary.ModuleCount != 2 {
		t.Fatalf("module count = %d, want 2", report.Summary.ModuleCount)
	}
	if report.Summary.ReportCount != 1 {
		t.Fatalf("report count = %d, want 1", report.Summary.ReportCount)
	}
	if report.Summary.FindingCount != 0 {
		t.Fatalf("finding count = %d, want 0", report.Summary.FindingCount)
	}
	if len(report.Findings) != 0 {
		t.Fatalf("findings = %d, want 0", len(report.Findings))
	}
}

func TestRunFiltersModuleByArtifactID(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{"fails", "-project", "../../demo/multi-module-failure", "-module", "payment-core"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("exit code = %d, want 1", code)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}

	text := stdout.String()
	for _, expected := range []string{
		"Findings: 1",
		"Module: payment-core (payment-core)",
		"maven-surefire-plugin",
	} {
		if !strings.Contains(text, expected) {
			t.Fatalf("stdout missing %q\n%s", expected, text)
		}
	}
	if strings.Contains(text, "payment-api") {
		t.Fatalf("stdout contains filtered module payment-api\n%s", text)
	}
}

func TestRunFiltersModuleByPath(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{"why", "-project", "../../pkg/prmaven/testdata/nested-module-project", "-module", "platform/service-core"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("exit code = %d, want 1", code)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}

	text := stdout.String()
	for _, expected := range []string{
		"Findings: 1",
		"Module: service-core (platform/service-core)",
		"platform/service-core/target/surefire-reports/TEST-dev.prmaven.demo.NestedPaymentTest.xml",
	} {
		if !strings.Contains(text, expected) {
			t.Fatalf("stdout missing %q\n%s", expected, text)
		}
	}
}

func TestRunFiltersModuleNoMatch(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{"fails", "-project", "../../demo/multi-module-failure", "-module", "does-not-exist"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0\nstderr=%s\nstdout=%s", code, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}

	text := stdout.String()
	if !strings.Contains(text, "Findings: 0") {
		t.Fatalf("stdout missing no-match findings summary\n%s", text)
	}
	if strings.Contains(text, "payment-core") || strings.Contains(text, "payment-api") {
		t.Fatalf("stdout contains filtered findings\n%s", text)
	}
}

func TestRunReturnsUsageExitCode(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{"unknown"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("exit code = %d, want 2", code)
	}
	if !strings.Contains(stderr.String(), `unknown command "unknown"`) {
		t.Fatalf("stderr missing unknown command\n%s", stderr.String())
	}
}

func TestRunVersion(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{"version"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if strings.TrimSpace(stdout.String()) == "" {
		t.Fatal("version output is empty")
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestCLIEndToEndJSON(t *testing.T) {
	command := exec.Command("go", "run", ".", "why", "-project", "../../demo/multi-module-failure", "-format", "json")
	output, err := command.CombinedOutput()
	if err == nil {
		t.Fatal("CLI exit code = 0, want non-zero when findings are present")
	}

	text := string(output)
	for _, expected := range []string{
		`"findingCount": 2`,
		`"mavenPlugin": "maven-surefire-plugin"`,
		`"reproduceCommand": "mvn -pl payment-api -am -Dit.test=PaymentApiIT verify"`,
	} {
		if !strings.Contains(text, expected) {
			t.Fatalf("CLI JSON output missing %q\n%s", expected, text)
		}
	}
}

func TestCLIEndToEndNoFailure(t *testing.T) {
	command := exec.Command("go", "run", ".", "fails", "-project", "../../demo/no-failure")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI exit error = %v\n%s", err, string(output))
	}

	text := string(output)
	for _, expected := range []string{
		"Modules: 2 | Reports: 1 | Findings: 0",
		"No Maven test or quality failures found in Surefire, Failsafe, Checkstyle, SpotBugs, Maven Enforcer, or JaCoCo reports.",
	} {
		if !strings.Contains(text, expected) {
			t.Fatalf("CLI output missing %q\n%s", expected, text)
		}
	}
}

func TestCLIEndToEndNoFailureJSON(t *testing.T) {
	command := exec.Command("go", "run", ".", "why", "-project", "../../demo/no-failure", "-format", "json")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI exit error = %v\n%s", err, string(output))
	}

	var report struct {
		Summary struct {
			ModuleCount  int `json:"moduleCount"`
			ReportCount  int `json:"reportCount"`
			FindingCount int `json:"findingCount"`
		} `json:"summary"`
		Findings []struct{} `json:"findings"`
	}
	if err := json.Unmarshal(output, &report); err != nil {
		t.Fatalf("CLI JSON output is invalid: %v\n%s", err, string(output))
	}

	if report.Summary.ModuleCount != 2 {
		t.Fatalf("module count = %d, want 2", report.Summary.ModuleCount)
	}
	if report.Summary.ReportCount != 1 {
		t.Fatalf("report count = %d, want 1", report.Summary.ReportCount)
	}
	if report.Summary.FindingCount != 0 {
		t.Fatalf("finding count = %d, want 0", report.Summary.FindingCount)
	}
	if len(report.Findings) != 0 {
		t.Fatalf("findings = %d, want 0", len(report.Findings))
	}
}

func TestCLIInvalidUsage(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "unknown command",
			args:     []string{"unknown"},
			expected: `unknown command "unknown"`,
		},
		{
			name:     "unknown format",
			args:     []string{"fails", "-project", "../../demo/no-failure", "-format", "xml"},
			expected: `unknown format "xml"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command := exec.Command("go", append([]string{"run", "."}, tt.args...)...)
			output, err := command.CombinedOutput()
			if err == nil {
				t.Fatalf("CLI exit code = 0, want non-zero\n%s", string(output))
			}
			if !strings.Contains(string(output), tt.expected) {
				t.Fatalf("CLI output missing %q\n%s", tt.expected, string(output))
			}
		})
	}
}
