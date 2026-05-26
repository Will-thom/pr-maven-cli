# Examples

This directory contains copyable examples for using PR Maven CLI as a command-line tool, as a Go library, and inside CI.

## CLI Against The Demo Fixture

From the repository root:

```bash
go run ./cmd/prmaven fails -project demo/multi-module-failure
go run ./cmd/prmaven why -project demo/multi-module-failure -format json
```

No-failure fixture:

```bash
go run ./cmd/prmaven fails -project demo/no-failure
```

## CLI Against A Real Maven Project

Generate reports first:

```bash
mvn -B verify
```

Then analyze the workspace:

```bash
prmaven why -project .
prmaven why -project . -format json > prmaven-report.json
```

Remember that `prmaven` exits with code `1` when findings exist. In CI, preserve the Maven exit code when Maven is the source of truth for pass/fail status.

## Go Library

Run the library example from the repository root:

```bash
go run ./examples/library demo/multi-module-failure
```

It prints a compact summary and each reproduction command returned by `pkg/prmaven`.

## GitHub Actions

The workflow example at [github-actions/triage-maven-failures.yml](github-actions/triage-maven-failures.yml) shows how to:

- run Maven;
- install PR Maven CLI;
- print text output;
- upload JSON output as an artifact;
- preserve the original Maven failure status.

