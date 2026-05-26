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

type checkstyleReport struct {
	Files []checkstyleFile `xml:"file"`
}

type checkstyleFile struct {
	Name   string            `xml:"name,attr"`
	Errors []checkstyleError `xml:"error"`
}

type checkstyleError struct {
	Line     string `xml:"line,attr"`
	Column   string `xml:"column,attr"`
	Severity string `xml:"severity,attr"`
	Message  string `xml:"message,attr"`
	Source   string `xml:"source,attr"`
}

type spotbugsReport struct {
	Bugs []spotbugsBug `xml:"BugInstance"`
}

type spotbugsBug struct {
	Type         string             `xml:"type,attr"`
	Category     string             `xml:"category,attr"`
	ShortMessage string             `xml:"ShortMessage"`
	LongMessage  string             `xml:"LongMessage"`
	Class        spotbugsClass      `xml:"Class"`
	SourceLine   spotbugsSourceLine `xml:"SourceLine"`
}

type spotbugsClass struct {
	ClassName  string             `xml:"classname,attr"`
	SourceLine spotbugsSourceLine `xml:"SourceLine"`
}

type spotbugsSourceLine struct {
	ClassName  string `xml:"classname,attr"`
	SourceFile string `xml:"sourcefile,attr"`
	SourcePath string `xml:"sourcepath,attr"`
	Start      string `xml:"start,attr"`
	End        string `xml:"end,attr"`
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

	reportFiles, err := findReportFiles(absRoot)
	if err != nil {
		return Report{}, err
	}

	var findings []Finding
	for _, reportFile := range reportFiles {
		reportFindings, err := parseReport(absRoot, moduleByPath, reportFile)
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

func findReportFiles(projectRoot string) ([]reportFile, error) {
	junitReports, err := findJUnitReports(projectRoot)
	if err != nil {
		return nil, err
	}
	checkstyleReports, err := findCheckstyleReports(projectRoot)
	if err != nil {
		return nil, err
	}
	spotbugsReports, err := findSpotBugsReports(projectRoot)
	if err != nil {
		return nil, err
	}

	reports := append(junitReports, checkstyleReports...)
	reports = append(reports, spotbugsReports...)
	sort.Slice(reports, func(i, j int) bool {
		return reports[i].relPath < reports[j].relPath
	})
	return reports, nil
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

func findCheckstyleReports(projectRoot string) ([]reportFile, error) {
	return findNamedTargetReports(projectRoot, "checkstyle", map[string]bool{
		"checkstyle-result.xml": true,
	})
}

func findSpotBugsReports(projectRoot string) ([]reportFile, error) {
	return findNamedTargetReports(projectRoot, "spotbugs", map[string]bool{
		"spotbugs.xml":    true,
		"spotbugsXml.xml": true,
	})
}

func findNamedTargetReports(projectRoot, kind string, names map[string]bool) ([]reportFile, error) {
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
		if !names[entry.Name()] || !isInsideTarget(path) {
			return nil
		}

		rel := relativePath(projectRoot, path)
		modulePath := inferTargetModulePath(projectRoot, path)
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

func parseReport(projectRoot string, moduleByPath map[string]Module, reportFile reportFile) ([]Finding, error) {
	switch reportFile.kind {
	case "checkstyle":
		return parseCheckstyleReport(projectRoot, moduleByPath, reportFile)
	case "spotbugs":
		return parseSpotBugsReport(projectRoot, moduleByPath, reportFile)
	default:
		return parseJUnitReport(projectRoot, moduleByPath, reportFile)
	}
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

func parseCheckstyleReport(projectRoot string, moduleByPath map[string]Module, reportFile reportFile) ([]Finding, error) {
	data, err := os.ReadFile(reportFile.absPath)
	if err != nil {
		return nil, fmt.Errorf("read report %s: %w", reportFile.relPath, err)
	}

	var report checkstyleReport
	if err := xml.Unmarshal(data, &report); err != nil {
		return nil, fmt.Errorf("parse report %s: %w", reportFile.relPath, err)
	}

	module := moduleByPath[reportFile.modulePath]
	if module.Path == "" {
		module = Module{
			Name: moduleNameFromPath(filepath.FromSlash(reportFile.modulePath)),
			Path: reportFile.modulePath,
		}
	}

	var findings []Finding
	for _, file := range report.Files {
		sourcePath := checkstyleSourcePath(projectRoot, file.Name)
		for _, violation := range file.Errors {
			findings = append(findings, buildCheckstyleFinding(module, reportFile, sourcePath, violation))
		}
	}
	return findings, nil
}

func parseSpotBugsReport(projectRoot string, moduleByPath map[string]Module, reportFile reportFile) ([]Finding, error) {
	data, err := os.ReadFile(reportFile.absPath)
	if err != nil {
		return nil, fmt.Errorf("read report %s: %w", reportFile.relPath, err)
	}

	var report spotbugsReport
	if err := xml.Unmarshal(data, &report); err != nil {
		return nil, fmt.Errorf("parse report %s: %w", reportFile.relPath, err)
	}

	module := moduleByPath[reportFile.modulePath]
	if module.Path == "" {
		module = Module{
			Name: moduleNameFromPath(filepath.FromSlash(reportFile.modulePath)),
			Path: reportFile.modulePath,
		}
	}

	var findings []Finding
	for _, bug := range report.Bugs {
		findings = append(findings, buildSpotBugsFinding(module, reportFile, projectRoot, bug))
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

func buildCheckstyleFinding(module Module, reportFile reportFile, sourcePath string, violation checkstyleError) Finding {
	location := checkstyleLocation(violation.Line, violation.Column)
	message := oneLine(violation.Message)
	if violation.Severity != "" && message != "" {
		message = violation.Severity + ": " + message
	}

	return Finding{
		ID:                 findingID(module.Path, sourcePath, location, "checkstyle"),
		Module:             module.Name,
		ModulePath:         module.Path,
		ReportPath:         reportFile.relPath,
		ReportKind:         reportFile.kind,
		MavenPlugin:        pluginForReportKind(reportFile.kind),
		MavenPhase:         phaseForReportKind(reportFile.kind),
		TestClass:          sourcePath,
		TestName:           location,
		FailureKind:        "violation",
		FailureType:        firstNonEmpty(violation.Source, violation.Severity),
		Message:            message,
		ReproduceCommand:   reproduceCommand(reportFile.kind, module.Path, ""),
		Confidence:         "high",
		ConfidenceReasons:  checkstyleConfidenceReasons(module.Path, sourcePath),
		SourceReportFormat: "checkstyle-xml",
	}
}

func buildSpotBugsFinding(module Module, reportFile reportFile, projectRoot string, bug spotbugsBug) Finding {
	sourceLine := spotbugsBestSourceLine(bug)
	sourcePath := spotbugsSourcePath(projectRoot, sourceLine, bug.Class.ClassName)
	location := spotbugsLocation(sourceLine.Start, sourceLine.End)
	message := oneLine(firstNonEmpty(bug.LongMessage, bug.ShortMessage, bug.Type))

	return Finding{
		ID:                 findingID(module.Path, sourcePath, location, bug.Type),
		Module:             module.Name,
		ModulePath:         module.Path,
		ReportPath:         reportFile.relPath,
		ReportKind:         reportFile.kind,
		MavenPlugin:        pluginForReportKind(reportFile.kind),
		MavenPhase:         phaseForReportKind(reportFile.kind),
		TestClass:          sourcePath,
		TestName:           location,
		FailureKind:        "bug",
		FailureType:        spotbugsFailureType(bug),
		Message:            message,
		ReproduceCommand:   reproduceCommand(reportFile.kind, module.Path, ""),
		Confidence:         "high",
		ConfidenceReasons:  spotbugsConfidenceReasons(module.Path, sourcePath, bug.Type),
		SourceReportFormat: "spotbugs-xml",
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

func inferTargetModulePath(projectRoot, reportPath string) string {
	for dir := filepath.Dir(reportPath); dir != "." && dir != string(filepath.Separator); dir = filepath.Dir(dir) {
		if filepath.Base(dir) == "target" {
			moduleDir := filepath.Dir(dir)
			if samePath(moduleDir, projectRoot) {
				return "."
			}
			return relativePath(projectRoot, moduleDir)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
	}
	return inferModulePath(projectRoot, reportPath)
}

func isInsideTarget(path string) bool {
	for dir := filepath.Dir(path); dir != "." && dir != string(filepath.Separator); dir = filepath.Dir(dir) {
		if filepath.Base(dir) == "target" {
			return true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return false
		}
	}
	return false
}

func checkstyleSourcePath(projectRoot, sourcePath string) string {
	sourcePath = strings.TrimSpace(sourcePath)
	if sourcePath == "" {
		return "unknown"
	}
	if filepath.IsAbs(sourcePath) {
		return slashPath(relativePath(projectRoot, sourcePath))
	}
	return slashPath(sourcePath)
}

func checkstyleLocation(line, column string) string {
	line = strings.TrimSpace(line)
	column = strings.TrimSpace(column)
	switch {
	case line != "" && column != "":
		return "line " + line + ", column " + column
	case line != "":
		return "line " + line
	case column != "":
		return "column " + column
	default:
		return "location unknown"
	}
}

func spotbugsBestSourceLine(bug spotbugsBug) spotbugsSourceLine {
	if bug.SourceLine.SourcePath != "" || bug.SourceLine.SourceFile != "" || bug.SourceLine.Start != "" {
		return bug.SourceLine
	}
	return bug.Class.SourceLine
}

func spotbugsSourcePath(projectRoot string, sourceLine spotbugsSourceLine, className string) string {
	sourcePath := firstNonEmpty(sourceLine.SourcePath, sourceLine.SourceFile)
	if sourcePath == "" && className != "" {
		sourcePath = strings.ReplaceAll(className, ".", "/") + ".java"
	}
	sourcePath = strings.TrimSpace(sourcePath)
	if sourcePath == "" {
		return "unknown"
	}
	if filepath.IsAbs(sourcePath) {
		return slashPath(relativePath(projectRoot, sourcePath))
	}
	return slashPath(sourcePath)
}

func spotbugsLocation(start, end string) string {
	start = strings.TrimSpace(start)
	end = strings.TrimSpace(end)
	switch {
	case start != "" && end != "" && start != end:
		return "lines " + start + "-" + end
	case start != "":
		return "line " + start
	case end != "":
		return "line " + end
	default:
		return "location unknown"
	}
}

func spotbugsFailureType(bug spotbugsBug) string {
	if bug.Category != "" && bug.Type != "" {
		return bug.Category + "/" + bug.Type
	}
	return firstNonEmpty(bug.Type, bug.Category)
}

func pluginForReportKind(kind string) string {
	if kind == "checkstyle" {
		return "maven-checkstyle-plugin"
	}
	if kind == "spotbugs" {
		return "spotbugs-maven-plugin"
	}
	if kind == "failsafe" {
		return "maven-failsafe-plugin"
	}
	return "maven-surefire-plugin"
}

func phaseForReportKind(kind string) string {
	if kind == "checkstyle" {
		return "verify"
	}
	if kind == "spotbugs" {
		return "verify"
	}
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
	if kind == "checkstyle" {
		parts = append(parts, "checkstyle:check")
		return strings.Join(parts, " ")
	}
	if kind == "spotbugs" {
		parts = append(parts, "spotbugs:check")
		return strings.Join(parts, " ")
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

func spotbugsConfidenceReasons(modulePath, sourcePath, bugType string) []string {
	reasons := []string{
		"bug was found in a Maven SpotBugs XML report",
		"report path maps to Maven module " + modulePath,
		"SpotBugs source entry maps to " + sourcePath,
	}
	if bugType != "" {
		reasons = append(reasons, "SpotBugs bug type is "+bugType)
	}
	return reasons
}

func checkstyleConfidenceReasons(modulePath, sourcePath string) []string {
	return []string{
		"violation was found in a Maven Checkstyle XML report",
		"report path maps to Maven module " + modulePath,
		"Checkstyle file entry maps to " + sourcePath,
	}
}
