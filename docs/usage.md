# Usage

PR Maven CLI turns Maven test, quality, and selected build log artifacts into focused failure context.

It answers:

```text
What failed, in which Maven module, through which Maven plugin, and how do I reproduce it locally?
```

## Commands

```bash
prmaven fails -project .
prmaven fails -project . -format json
prmaven why -project .
prmaven why -project . -format json
prmaven version
```

Stage 1 treats `fails` and `why` as equivalent analysis commands. The separate names leave room for future UX where `fails` lists failures and `why` adds richer causality evidence.

## Flags

- `-project`: Maven project directory to analyze. Defaults to `.`.
- `-format`: output format. Supported values are `text` and `json`. Defaults to `text`.

## Exit Codes

- `0`: analysis completed and no findings were found.
- `1`: analysis completed with Maven failure findings, or analysis failed.
- `2`: invalid CLI usage.

The non-zero finding exit code is useful in CI because a failed test suite should remain visible as a failed job while still printing structured context.

## Demo: Failure Fixture

From the repository root:

```bash
go run ./cmd/prmaven fails -project demo/multi-module-failure
```

Expected behavior:

- prints the Maven project root;
- discovers the aggregator and modules;
- reads Surefire reports from `target/surefire-reports`;
- reads Failsafe reports from `target/failsafe-reports`;
- reports two findings;
- prints minimal reproduction commands.

Example reproduction commands:

```bash
mvn -pl payment-core -am -Dtest=PaymentRoundingTest test
mvn -pl payment-api -am -Dit.test=PaymentApiIT verify
```

JSON output:

```bash
go run ./cmd/prmaven why -project demo/multi-module-failure -format json
```

The JSON contract is designed for CI systems, bots, and coding agents. It includes `summary`, `modules`, and `findings`.

For field-level details and compatibility expectations, read [JSON contract](json-contract.md).

## Provider Integrations

Stage 1 has no native GitHub or GitLab API adapter. The CLI does not need provider tokens and does not call remote PR, check-run, issue, or merge request APIs.

GitHub is the only platform with first-party project automation and a copyable CI example today. For the full integration scope and planned native adapters, read [integrations.md](integrations.md).

## Demo: No-Failure Fixture

```bash
go run ./cmd/prmaven fails -project demo/no-failure
```

Expected behavior:

- exits with code `0`;
- reports zero findings;
- confirms that no supported Maven test or quality failures were found.

## Use Against A Real Maven Workspace

First generate test report artifacts with Maven:

```bash
mvn -B test
```

For integration tests that use Failsafe:

```bash
mvn -B verify
```

Then analyze the workspace:

```bash
prmaven fails -project /path/to/maven/repo
```

The CLI scans module-level report folders and deterministic quality artifacts such as:

```text
target/surefire-reports
target/failsafe-reports
target/checkstyle-result.xml
target/spotbugsXml.xml
target/spotbugs.xml
target/maven-enforcer.log
target/jacoco.log
target/maven.log
```

## CI Pattern

This pattern keeps the Maven failure visible while still collecting PR Maven CLI output.

```bash
set +e
mvn -B verify
maven_status=$?

prmaven why -project . || true
prmaven why -project . -format json > prmaven-report.json || true

exit "$maven_status"
```

A complete GitHub Actions example is available at [../examples/github-actions/triage-maven-failures.yml](../examples/github-actions/triage-maven-failures.yml).

## Library Usage

Use `pkg/prmaven` when another Go tool needs direct access to the report model.

```go
report, err := prmaven.Analyze(prmaven.Options{ProjectDir: "."})
if err != nil {
    return err
}
for _, finding := range report.Findings {
    fmt.Println(finding.ReproduceCommand)
}
```

A runnable example is available at [../examples/library](../examples/library).

