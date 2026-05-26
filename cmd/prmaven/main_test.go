package main

import (
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
