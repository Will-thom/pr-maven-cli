package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/Will-thom/pr-maven-cli/pkg/prmaven"
)

var version = "dev"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	command := "fails"
	if len(args) > 0 && args[0] != "" && args[0][0] != '-' {
		command = args[0]
		args = args[1:]
	}

	flags := flag.NewFlagSet("prmaven", flag.ContinueOnError)
	flags.SetOutput(stderr)
	projectDir := flags.String("project", ".", "Maven project directory")
	format := flags.String("format", "text", "output format: text or json")
	outputPath := flags.String("output", "", "write output to file instead of stdout")

	if err := flags.Parse(args); err != nil {
		return 2
	}

	if flags.NArg() > 0 {
		command = flags.Arg(0)
	}

	switch command {
	case "fails", "why", "version":
	default:
		fmt.Fprintf(stderr, "unknown command %q\n", command)
		fmt.Fprintln(stderr, "available commands: fails, why, version")
		return 2
	}

	if command == "version" {
		fmt.Fprintln(stdout, version)
		return 0
	}

	report, err := prmaven.Analyze(prmaven.Options{ProjectDir: *projectDir})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	switch *format {
	case "json", "text":
	default:
		fmt.Fprintf(stderr, "unknown format %q\n", *format)
		fmt.Fprintln(stderr, "available formats: text, json")
		return 2
	}

	output := stdout
	var outputFile *os.File
	if *outputPath != "" {
		outputFile, err = os.Create(*outputPath)
		if err != nil {
			fmt.Fprintf(stderr, "create output file %q: %v\n", *outputPath, err)
			return 1
		}
		output = outputFile
	}

	var writeErr error
	switch *format {
	case "json":
		writeErr = prmaven.WriteJSON(output, report)
	case "text":
		writeErr = prmaven.WriteText(output, report)
	}

	if outputFile != nil {
		if closeErr := outputFile.Close(); writeErr == nil && closeErr != nil {
			writeErr = closeErr
		}
	}

	if writeErr != nil {
		fmt.Fprintln(stderr, writeErr)
		return 1
	}

	if report.Summary.FindingCount > 0 {
		return 1
	}
	return 0
}
