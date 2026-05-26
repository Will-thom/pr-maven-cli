package prmaven

import (
	"path/filepath"
	"regexp"
	"strings"
)

var unsafeIDCharacters = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

func relativePath(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return path
	}
	if rel == "" {
		return "."
	}
	return rel
}

func slashPath(path string) string {
	if path == "" || path == "." {
		return "."
	}
	path = strings.ReplaceAll(path, "\\", "/")
	return filepath.ToSlash(filepath.Clean(path))
}

func samePath(left, right string) bool {
	leftAbs, leftErr := filepath.Abs(left)
	rightAbs, rightErr := filepath.Abs(right)
	if leftErr != nil || rightErr != nil {
		return filepath.Clean(left) == filepath.Clean(right)
	}
	return filepath.Clean(leftAbs) == filepath.Clean(rightAbs)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}

func oneLine(value string) string {
	fields := strings.Fields(value)
	return strings.Join(fields, " ")
}

func simpleClassName(className string) string {
	className = strings.TrimSpace(className)
	if className == "" {
		return ""
	}
	parts := strings.Split(className, ".")
	return parts[len(parts)-1]
}

func classNameFromReport(reportPath string) string {
	base := filepath.Base(reportPath)
	base = strings.TrimSuffix(base, ".xml")
	return strings.TrimPrefix(base, "TEST-")
}

func findingID(modulePath, className, testName, kind string) string {
	raw := strings.Join([]string{modulePath, simpleClassName(className), testName, kind}, ":")
	raw = strings.ReplaceAll(raw, ".", "-")
	raw = strings.ReplaceAll(raw, "/", "-")
	raw = unsafeIDCharacters.ReplaceAllString(raw, "-")
	raw = strings.Trim(raw, "-")
	if raw == "" {
		return "maven-finding"
	}
	return strings.ToLower(raw)
}
