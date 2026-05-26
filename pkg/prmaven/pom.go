package prmaven

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type pomProject struct {
	ArtifactID string   `xml:"artifactId"`
	Modules    []string `xml:"modules>module"`
}

func discoverModules(projectRoot string) ([]Module, error) {
	rootPOM := filepath.Join(projectRoot, "pom.xml")
	if _, err := os.Stat(rootPOM); err != nil {
		return nil, fmt.Errorf("read Maven project root: %w", err)
	}

	seen := map[string]bool{}
	var modules []Module

	var walk func(relPath string) error
	walk = func(relPath string) error {
		cleanRel := filepath.Clean(relPath)
		if cleanRel == "." {
			cleanRel = "."
		}
		if seen[cleanRel] {
			return nil
		}
		seen[cleanRel] = true

		pomPath := filepath.Join(projectRoot, cleanRel, "pom.xml")
		if cleanRel == "." {
			pomPath = rootPOM
		}

		project, err := readPOM(pomPath)
		if err != nil {
			return err
		}

		name := strings.TrimSpace(project.ArtifactID)
		if name == "" {
			name = moduleNameFromPath(cleanRel)
		}
		modules = append(modules, Module{
			Name: name,
			Path: slashPath(cleanRel),
			POM:  slashPath(relativePath(projectRoot, pomPath)),
		})

		for _, child := range project.Modules {
			child = strings.TrimSpace(child)
			if child == "" {
				continue
			}
			childRel := filepath.Join(cleanRel, filepath.FromSlash(child))
			if cleanRel == "." {
				childRel = filepath.FromSlash(child)
			}
			if err := walk(childRel); err != nil {
				return err
			}
		}
		return nil
	}

	if err := walk("."); err != nil {
		return nil, err
	}

	sort.Slice(modules, func(i, j int) bool {
		return modules[i].Path < modules[j].Path
	})
	return modules, nil
}

func readPOM(path string) (pomProject, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return pomProject{}, fmt.Errorf("read pom %s: %w", path, err)
	}
	var project pomProject
	if err := xml.Unmarshal(data, &project); err != nil {
		return pomProject{}, fmt.Errorf("parse pom %s: %w", path, err)
	}
	return project, nil
}

func moduleNameFromPath(path string) string {
	if path == "." {
		return "root"
	}
	return filepath.Base(path)
}
