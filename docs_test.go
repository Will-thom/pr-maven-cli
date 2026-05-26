package prmaven_test

import (
	"os"
	"strings"
	"testing"
)

func TestDocumentationCoversInstallationUsageAndExamples(t *testing.T) {
	files := map[string]string{
		"README.md":                mustReadFile(t, "README.md"),
		"docs/installation.md":     mustReadFile(t, "docs/installation.md"),
		"docs/usage.md":            mustReadFile(t, "docs/usage.md"),
		"examples/README.md":       mustReadFile(t, "examples/README.md"),
		"examples/library/main.go": mustReadFile(t, "examples/library/main.go"),
	}

	assertContains(t, files["README.md"], "[Installation](docs/installation.md)")
	assertContains(t, files["README.md"], "[Usage guide](docs/usage.md)")
	assertContains(t, files["README.md"], "[Examples](examples/README.md)")
	assertContains(t, files["docs/installation.md"], "go install ./cmd/prmaven")
	assertContains(t, files["docs/installation.md"], "prmaven version")
	assertContains(t, files["docs/usage.md"], "prmaven fails -project .")
	assertContains(t, files["docs/usage.md"], "demo/multi-module-failure")
	assertContains(t, files["docs/usage.md"], "demo/no-failure")
	assertContains(t, files["examples/README.md"], "go run ./examples/library demo/multi-module-failure")
	assertContains(t, files["examples/library/main.go"], "prmaven.Analyze")
}

func mustReadFile(t *testing.T, path string) string {
	t.Helper()

	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(contents)
}

func assertContains(t *testing.T, text, expected string) {
	t.Helper()

	if !strings.Contains(text, expected) {
		t.Fatalf("documentation missing %q", expected)
	}
}
