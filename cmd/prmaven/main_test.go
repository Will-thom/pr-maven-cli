package main

import (
	"bytes"
	"os/exec"
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
		"No Maven test or quality failures found in Surefire, Failsafe, or Checkstyle reports.",
	} {
		if !strings.Contains(text, expected) {
			t.Fatalf("CLI output missing %q\n%s", expected, text)
		}
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
