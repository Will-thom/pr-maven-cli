package prmaven

import (
	"encoding/xml"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Options struct {
	ProjectDir string
}

type junitReport struct {
	XMLName xml.Name
	Suites  []junitSuite `xml:"testsuite"`
	Cases   []junitCase  `xml:"testcase"`
}

type junitSuite struct {
	Name  string      `xml:"name,attr"`
	Cases []junitCase `xml:"testcase"`
}

type junitCase struct {
	ClassName string         `xml:"classname,attr"`
	Name      string         `xml:"name,attr"`
	Failures  []junitProblem `xml:"failure"`
	Errors    []junitProblem `xml:"error"`
}

type junitProblem struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
	Body    string `xml:",chardata"`
}

type reportFile struct {
	absPath    string
	relPath    string
	modulePath string
	kind       string
}

func Analyze(options Options) (Report, error) {
	projectRoot := options.ProjectDir
	if strings.TrimSpace(projectRoot) == "" {
		projectRoot = "."
	}
	absRoot, err := filepath.Abs(projectRoot)
	if err != nil {
		return Report{}, err
	}

	modules, err := discoverModules(absRoot)
	if err != nil {
		return Report{}, err
	}
	moduleByPath := make(map[string]Module, len(modules))
	for _, module := range modules {
		moduleByPath[module.Path] = module
	}

	reportFiles, err := findJUnitReports(absRoot)
	if err != nil {
		return Report{}, err
	}

	var findings []Finding
	for _, reportFile := range reportFiles {
		reportFindings, err := parseJUnitReport(absRoot, moduleByPath, reportFile)
		if err != nil {
			return Report{}, err
		}
		findings = append(findings, reportFindings...)
	}

	sort.Slice(findings, func(i, j int) bool {
		return findings[i].ID < findings[j].ID
	})

	return Report{
		ProjectRoot: absRoot,
		Summary: Summary{
			ModuleCount:  len(modules),
			ReportCount:  len(reportFiles),
			FindingCount: len(findings),
		},
		Modules:  modules,
		Findings: findings,
	}, nil
}

func findJUnitReports(projectRoot string) ([]reportFile, error) {
	var reports []reportFile

	err := filepath.WalkDir(projectRoot, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			switch entry.Name() {
			case ".git", ".idea", ".mvn":
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasPrefix(entry.Name(), "TEST-") || !strings.HasSuffix(entry.Name(), ".xml") {
			return nil
		}

		reportDir := filepath.Base(filepath.Dir(path))
		kind := ""
		switch reportDir {
		case "surefire-reports":
			kind = "surefire"
		case "failsafe-reports":
			kind = "failsafe"
		default:
			return nil
		}

		rel := relativePath(projectRoot, path)
		modulePath := inferModulePath(projectRoot, path)
		reports = append(reports, reportFile{
			absPath:    path,
			relPath:    slashPath(rel),
			modulePath: slashPath(modulePath),
			kind:       kind,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(reports, func(i, j int) bool {
		return reports[i].relPath < reports[j].relPath
	})
	return reports, nil
}

func parseJUnitReport(projectRoot string, moduleByPath map[string]Module, reportFile reportFile) ([]Finding, error) {
	data, err := os.ReadFile(reportFile.absPath)
	if err != nil {
		return nil, fmt.Errorf("read report %s: %w", reportFile.relPath, err)
	}

	var report junitReport
	if err := xml.Unmarshal(data, &report); err != nil {
		return nil, fmt.Errorf("parse report %s: %w", reportFile.relPath, err)
	}

	var cases []junitCase
	cases = append(cases, report.Cases...)
	for _, suite := range report.Suites {
		cases = append(cases, suite.Cases...)
	}

	module := moduleByPath[reportFile.modulePath]
	if module.Path == "" {
		module = Module{
			Name: moduleNameFromPath(filepath.FromSlash(reportFile.modulePath)),
			Path: reportFile.modulePath,
		}
	}

	var findings []Finding
	for _, testCase := range cases {
		for _, problem := range testCase.Failures {
			findings = append(findings, buildFinding(module, reportFile, testCase, problem, "failure"))
		}
		for _, problem := range testCase.Errors {
			findings = append(findings, buildFinding(module, reportFile, testCase, problem, "error"))
		}
	}
	return findings, nil
}

func buildFinding(module Module, reportFile reportFile, testCase junitCase, problem junitProblem, kind string) Finding {
	className := strings.TrimSpace(testCase.ClassName)
	if className == "" {
		className = classNameFromReport(reportFile.relPath)
	}
	testName := strings.TrimSpace(testCase.Name)
	message := firstNonEmpty(problem.Message, strings.TrimSpace(problem.Body))

	return Finding{
		ID:                 findingID(module.Path, className, testName, kind),
		Module:             module.Name,
		ModulePath:         module.Path,
		ReportPath:         reportFile.relPath,
		ReportKind:         reportFile.kind,
		MavenPlugin:        pluginForReportKind(reportFile.kind),
		MavenPhase:         phaseForReportKind(reportFile.kind),
		TestClass:          className,
		TestName:           testName,
		FailureKind:        kind,
		FailureType:        problem.Type,
		Message:            oneLine(message),
		ReproduceCommand:   reproduceCommand(reportFile.kind, module.Path, className),
		Confidence:         "high",
		ConfidenceReasons:  confidenceReasons(reportFile.kind, module.Path, className),
		SourceReportFormat: "junit-xml",
	}
}

func inferModulePath(projectRoot, reportPath string) string {
	reportDir := filepath.Dir(reportPath)
	targetDir := filepath.Dir(reportDir)
	moduleDir := filepath.Dir(targetDir)
	if samePath(moduleDir, projectRoot) {
		return "."
	}
	return relativePath(projectRoot, moduleDir)
}

func pluginForReportKind(kind string) string {
	if kind == "failsafe" {
		return "maven-failsafe-plugin"
	}
	return "maven-surefire-plugin"
}

func phaseForReportKind(kind string) string {
	if kind == "failsafe" {
		return "verify"
	}
	return "test"
}

func reproduceCommand(kind, modulePath, className string) string {
	class := simpleClassName(className)
	parts := []string{"mvn"}
	if modulePath != "." {
		parts = append(parts, "-pl", modulePath, "-am")
	}
	if kind == "failsafe" {
		parts = append(parts, "-Dit.test="+class, "verify")
	} else {
		parts = append(parts, "-Dtest="+class, "test")
	}
	return strings.Join(parts, " ")
}

func confidenceReasons(kind, modulePath, className string) []string {
	reportName := "Surefire"
	if kind == "failsafe" {
		reportName = "Failsafe"
	}
	reasons := []string{
		"failure was found in a Maven " + reportName + " JUnit XML report",
		"report path maps to Maven module " + modulePath,
	}
	if className != "" {
		reasons = append(reasons, "reproduction command targets test class "+simpleClassName(className))
	}
	return reasons
}
