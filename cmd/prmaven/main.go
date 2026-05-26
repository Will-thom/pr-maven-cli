package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Will-thom/pr-maven-cli/pkg/prmaven"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	command := "fails"
	if len(args) > 0 && args[0] != "" && args[0][0] != '-' {
		command = args[0]
		args = args[1:]
	}

	flags := flag.NewFlagSet("prmaven", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	projectDir := flags.String("project", ".", "Maven project directory")
	format := flags.String("format", "text", "output format: text or json")

	if err := flags.Parse(args); err != nil {
		return 2
	}

	if flags.NArg() > 0 {
		command = flags.Arg(0)
	}

	switch command {
	case "fails", "why":
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n", command)
		fmt.Fprintln(os.Stderr, "available commands: fails, why")
		return 2
	}

	report, err := prmaven.Analyze(prmaven.Options{ProjectDir: *projectDir})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	switch *format {
	case "json":
		if err := prmaven.WriteJSON(os.Stdout, report); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
	case "text":
		if err := prmaven.WriteText(os.Stdout, report); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown format %q\n", *format)
		fmt.Fprintln(os.Stderr, "available formats: text, json")
		return 2
	}

	if report.Summary.FindingCount > 0 {
		return 1
	}
	return 0
}
