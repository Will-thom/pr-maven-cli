package prmaven

import (
	"bytes"
	"encoding/json"
	"os"
	"regexp"
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

func TestAnalyzeNoFailureDemoProject(t *testing.T) {
	report, err := Analyze(Options{ProjectDir: "../../demo/no-failure"})
	if err != nil {
		t.Fatal(err)
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

func TestAnalyzeCheckstyleReport(t *testing.T) {
	report, err := Analyze(Options{ProjectDir: "testdata/checkstyle-project"})
	if err != nil {
		t.Fatal(err)
	}

	if report.Summary.ModuleCount != 2 {
		t.Fatalf("module count = %d, want 2", report.Summary.ModuleCount)
	}
	if report.Summary.ReportCount != 1 {
		t.Fatalf("report count = %d, want 1", report.Summary.ReportCount)
	}
	if report.Summary.FindingCount != 1 {
		t.Fatalf("finding count = %d, want 1", report.Summary.FindingCount)
	}

	finding := report.Findings[0]
	if finding.ModulePath != "service-core" {
		t.Fatalf("module path = %q, want service-core", finding.ModulePath)
	}
	if finding.ReportPath != "service-core/target/checkstyle-result.xml" {
		t.Fatalf("report path = %q", finding.ReportPath)
	}
	if finding.ReportKind != "checkstyle" {
		t.Fatalf("report kind = %q, want checkstyle", finding.ReportKind)
	}
	if finding.MavenPlugin != "maven-checkstyle-plugin" {
		t.Fatalf("plugin = %q, want maven-checkstyle-plugin", finding.MavenPlugin)
	}
	if finding.MavenPhase != "verify" {
		t.Fatalf("phase = %q, want verify", finding.MavenPhase)
	}
	if finding.TestClass != "service-core/src/main/java/dev/prmaven/demo/OrderService.java" {
		t.Fatalf("source file = %q", finding.TestClass)
	}
	if finding.TestName != "line 17, column 5" {
		t.Fatalf("source location = %q", finding.TestName)
	}
	if finding.FailureKind != "violation" {
		t.Fatalf("failure kind = %q, want violation", finding.FailureKind)
	}
	if finding.FailureType != "com.puppycrawl.tools.checkstyle.checks.sizes.LineLengthCheck" {
		t.Fatalf("failure type = %q", finding.FailureType)
	}
	if finding.Message != "error: Line is longer than 120 characters (found 134)." {
		t.Fatalf("message = %q", finding.Message)
	}
	if finding.ReproduceCommand != "mvn -pl service-core -am checkstyle:check" {
		t.Fatalf("reproduce command = %q", finding.ReproduceCommand)
	}
	if finding.SourceReportFormat != "checkstyle-xml" {
		t.Fatalf("source format = %q, want checkstyle-xml", finding.SourceReportFormat)
	}
	assertReasonsContain(t, finding.ConfidenceReasons, "violation was found in a Maven Checkstyle XML report")
	assertReasonsContain(t, finding.ConfidenceReasons, "report path maps to Maven module service-core")
}

func TestAnalyzeSpotBugsReport(t *testing.T) {
	report, err := Analyze(Options{ProjectDir: "testdata/spotbugs-project"})
	if err != nil {
		t.Fatal(err)
	}

	if report.Summary.ModuleCount != 2 {
		t.Fatalf("module count = %d, want 2", report.Summary.ModuleCount)
	}
	if report.Summary.ReportCount != 1 {
		t.Fatalf("report count = %d, want 1", report.Summary.ReportCount)
	}
	if report.Summary.FindingCount != 1 {
		t.Fatalf("finding count = %d, want 1", report.Summary.FindingCount)
	}

	finding := report.Findings[0]
	if finding.ModulePath != "service-core" {
		t.Fatalf("module path = %q, want service-core", finding.ModulePath)
	}
	if finding.ReportPath != "service-core/target/spotbugsXml.xml" {
		t.Fatalf("report path = %q", finding.ReportPath)
	}
	if finding.ReportKind != "spotbugs" {
		t.Fatalf("report kind = %q, want spotbugs", finding.ReportKind)
	}
	if finding.MavenPlugin != "spotbugs-maven-plugin" {
		t.Fatalf("plugin = %q, want spotbugs-maven-plugin", finding.MavenPlugin)
	}
	if finding.MavenPhase != "verify" {
		t.Fatalf("phase = %q, want verify", finding.MavenPhase)
	}
	if finding.TestClass != "service-core/src/main/java/dev/prmaven/demo/OrderAnalyzer.java" {
		t.Fatalf("source file = %q", finding.TestClass)
	}
	if finding.TestName != "line 42" {
		t.Fatalf("source location = %q", finding.TestName)
	}
	if finding.FailureKind != "bug" {
		t.Fatalf("failure kind = %q, want bug", finding.FailureKind)
	}
	if finding.FailureType != "CORRECTNESS/NP_NULL_ON_SOME_PATH" {
		t.Fatalf("failure type = %q", finding.FailureType)
	}
	if finding.Message != "Possible null pointer dereference of order in dev.prmaven.demo.OrderAnalyzer.analyze(Order)" {
		t.Fatalf("message = %q", finding.Message)
	}
	if finding.ReproduceCommand != "mvn -pl service-core -am spotbugs:check" {
		t.Fatalf("reproduce command = %q", finding.ReproduceCommand)
	}
	if finding.SourceReportFormat != "spotbugs-xml" {
		t.Fatalf("source format = %q, want spotbugs-xml", finding.SourceReportFormat)
	}
	assertReasonsContain(t, finding.ConfidenceReasons, "bug was found in a Maven SpotBugs XML report")
	assertReasonsContain(t, finding.ConfidenceReasons, "report path maps to Maven module service-core")
}

func TestAnalyzeEnforcerLog(t *testing.T) {
	report, err := Analyze(Options{ProjectDir: "testdata/enforcer-project"})
	if err != nil {
		t.Fatal(err)
	}

	if report.Summary.ModuleCount != 2 {
		t.Fatalf("module count = %d, want 2", report.Summary.ModuleCount)
	}
	if report.Summary.ReportCount != 1 {
		t.Fatalf("report count = %d, want 1", report.Summary.ReportCount)
	}
	if report.Summary.FindingCount != 1 {
		t.Fatalf("finding count = %d, want 1", report.Summary.FindingCount)
	}

	finding := report.Findings[0]
	if finding.ModulePath != "service-core" {
		t.Fatalf("module path = %q, want service-core", finding.ModulePath)
	}
	if finding.ReportPath != "service-core/target/maven-enforcer.log" {
		t.Fatalf("report path = %q", finding.ReportPath)
	}
	if finding.ReportKind != "enforcer" {
		t.Fatalf("report kind = %q, want enforcer", finding.ReportKind)
	}
	if finding.MavenPlugin != "maven-enforcer-plugin" {
		t.Fatalf("plugin = %q, want maven-enforcer-plugin", finding.MavenPlugin)
	}
	if finding.MavenPhase != "validate" {
		t.Fatalf("phase = %q, want validate", finding.MavenPhase)
	}
	if finding.TestClass != "maven-enforcer-plugin" {
		t.Fatalf("log source = %q", finding.TestClass)
	}
	if finding.TestName != "require-maven-baseline" {
		t.Fatalf("execution = %q", finding.TestName)
	}
	if finding.FailureKind != "rule" {
		t.Fatalf("failure kind = %q, want rule", finding.FailureKind)
	}
	if finding.FailureType != "org.apache.maven.enforcer.rules.version.RequireMavenVersion" {
		t.Fatalf("failure type = %q", finding.FailureType)
	}
	if finding.Message != "Detected Maven Version: 3.8.8 is not in the allowed range [3.9.0,)." {
		t.Fatalf("message = %q", finding.Message)
	}
	if finding.ReproduceCommand != "mvn -pl service-core -am enforcer:enforce" {
		t.Fatalf("reproduce command = %q", finding.ReproduceCommand)
	}
	if finding.SourceReportFormat != "maven-log" {
		t.Fatalf("source format = %q, want maven-log", finding.SourceReportFormat)
	}
	assertReasonsContain(t, finding.ConfidenceReasons, "failure was found in a Maven log containing maven-enforcer-plugin output")
	assertReasonsContain(t, finding.ConfidenceReasons, "report path maps to Maven module service-core")
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

func TestWriteTextIncludesCheckstyleSourceContext(t *testing.T) {
	report, err := Analyze(Options{ProjectDir: "testdata/checkstyle-project"})
	if err != nil {
		t.Fatal(err)
	}

	var output bytes.Buffer
	if err := WriteText(&output, report); err != nil {
		t.Fatal(err)
	}

	text := output.String()
	for _, expected := range []string{
		"Plugin: maven-checkstyle-plugin",
		"Source: service-core/src/main/java/dev/prmaven/demo/OrderService.java (line 17, column 5)",
		"Message: error: Line is longer than 120 characters (found 134).",
		"Reproduce: mvn -pl service-core -am checkstyle:check",
	} {
		if !strings.Contains(text, expected) {
			t.Fatalf("text output missing %q\n%s", expected, text)
		}
	}
}

func TestWriteTextIncludesSpotBugsSourceContext(t *testing.T) {
	report, err := Analyze(Options{ProjectDir: "testdata/spotbugs-project"})
	if err != nil {
		t.Fatal(err)
	}

	var output bytes.Buffer
	if err := WriteText(&output, report); err != nil {
		t.Fatal(err)
	}

	text := output.String()
	for _, expected := range []string{
		"Plugin: spotbugs-maven-plugin",
		"Source: service-core/src/main/java/dev/prmaven/demo/OrderAnalyzer.java (line 42)",
		"Kind: bug",
		"Message: Possible null pointer dereference of order in dev.prmaven.demo.OrderAnalyzer.analyze(Order)",
		"Reproduce: mvn -pl service-core -am spotbugs:check",
	} {
		if !strings.Contains(text, expected) {
			t.Fatalf("text output missing %q\n%s", expected, text)
		}
	}
}

func TestWriteTextIncludesEnforcerLogContext(t *testing.T) {
	report, err := Analyze(Options{ProjectDir: "testdata/enforcer-project"})
	if err != nil {
		t.Fatal(err)
	}

	var output bytes.Buffer
	if err := WriteText(&output, report); err != nil {
		t.Fatal(err)
	}

	text := output.String()
	for _, expected := range []string{
		"Plugin: maven-enforcer-plugin",
		"Phase: validate",
		"Log: maven-enforcer-plugin (require-maven-baseline)",
		"Kind: rule",
		"Message: Detected Maven Version: 3.8.8 is not in the allowed range [3.9.0,).",
		"Reproduce: mvn -pl service-core -am enforcer:enforce",
	} {
		if !strings.Contains(text, expected) {
			t.Fatalf("text output missing %q\n%s", expected, text)
		}
	}
}

func TestWriteTextMatchesGoldenFiles(t *testing.T) {
	tests := []struct {
		name       string
		projectDir string
		goldenPath string
	}{
		{
			name:       "multi module failure",
			projectDir: "../../demo/multi-module-failure",
			goldenPath: "testdata/golden/multi-module-failure.txt",
		},
		{
			name:       "no failure",
			projectDir: "../../demo/no-failure",
			goldenPath: "testdata/golden/no-failure.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report, err := Analyze(Options{ProjectDir: tt.projectDir})
			if err != nil {
				t.Fatal(err)
			}

			var output bytes.Buffer
			if err := WriteText(&output, report); err != nil {
				t.Fatal(err)
			}

			wantBytes, err := os.ReadFile(tt.goldenPath)
			if err != nil {
				t.Fatal(err)
			}

			got := normalizeTextOutput(output.String())
			want := normalizeTextOutput(string(wantBytes))
			if got != want {
				t.Fatalf("golden output mismatch\nwant:\n%s\n\ngot:\n%s", want, got)
			}
		})
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

func TestAnalyzeMissingProjectReturnsError(t *testing.T) {
	_, err := Analyze(Options{ProjectDir: "testdata/does-not-exist"})
	if err == nil {
		t.Fatal("error = nil, want missing project error")
	}
	if !strings.Contains(err.Error(), "read Maven project root") {
		t.Fatalf("error = %q, want Maven project root context", err.Error())
	}
}

var projectRootLine = regexp.MustCompile(`(?m)^Project: .+$`)

func normalizeTextOutput(value string) string {
	value = strings.ReplaceAll(value, "\r\n", "\n")
	value = projectRootLine.ReplaceAllString(value, "Project: <PROJECT_ROOT>")
	return strings.TrimRight(value, "\n")
}

func assertReasonsContain(t *testing.T, reasons []string, expected string) {
	t.Helper()

	for _, reason := range reasons {
		if reason == expected {
			return
		}
	}
	t.Fatalf("confidence reasons missing %q in %#v", expected, reasons)
}
