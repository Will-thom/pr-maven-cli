# PR Maven CLI

Maven-aware PR and CI triage for Java teams.

PR Maven CLI turns local Maven test reports into actionable failure context for humans, maintainers, CI systems, and coding agents.

It answers the first production question after a Maven PR fails:

```text
What failed, in which module, through which Maven plugin, and how do I reproduce it locally?
```

## Status

Stage 1 MVP.

The current release line focuses on deterministic local analysis of Maven report artifacts. It does not require GitHub tokens, CI API access, AI providers, telemetry, or external services.

Stage 1 does not include a native GitHub or GitLab API adapter. GitHub is currently the only platform with first-party project automation and a copyable CI example, while the CLI itself remains provider-agnostic. See [Integrations](docs/integrations.md).

Target Maven baseline:

- Maven 3.9.x.
- Documented against Maven 3.9.16, the latest release currently recommended by Apache Maven for all users.
- Maven 4 support is planned later, after the Maven 4 line is production-ready.

Apache's download page currently lists Maven 3.9.16 as the recommended release and Maven 4.x as a preview line that is not safe for production use: <https://maven.apache.org/download.cgi>.

## What It Does Today

- Discovers Maven modules from `pom.xml`.
- Parses Surefire JUnit XML reports.
- Parses Failsafe JUnit XML reports.
- Parses Checkstyle XML reports.
- Maps failures back to Maven modules.
- Identifies Maven plugin and phase.
- Generates a minimal Maven reproduction command.
- Emits human-readable text.
- Emits stable JSON for CI and agent usage.
- Provides a versioned demo project with Maven report fixtures.
- Provides GitHub Actions examples, without requiring GitHub API access at runtime.

## Install From Source

```bash
git clone https://github.com/Will-thom/pr-maven-cli.git
cd pr-maven-cli
go test ./...
go install ./cmd/prmaven
```

For release artifacts, local builds, `PATH` setup, and Windows notes, read [Installation](docs/installation.md).

## Quick Start

Run against the included demo:

```bash
go run ./cmd/prmaven fails -project demo/multi-module-failure
```

Example output:

```text
PR Maven CLI - Maven failure context

Modules: 3 | Reports: 2 | Findings: 2

Module: payment-core (payment-core)
Plugin: maven-surefire-plugin
Phase: test
Test: dev.prmaven.demo.PaymentRoundingTest.shouldRejectInvalidScale
Reproduce: mvn -pl payment-core -am -Dtest=PaymentRoundingTest test
Confidence: high
```

JSON output:

```bash
go run ./cmd/prmaven why -project demo/multi-module-failure -format json
```

Run against a real Maven workspace after CI/test artifacts exist:

```bash
prmaven fails -project /path/to/maven/repo
```

For flags, exit codes, CI patterns, and real-workspace usage, read the [Usage guide](docs/usage.md).

## CLI Commands

```bash
prmaven fails -project .
prmaven fails -project . -format json
prmaven why -project .
prmaven why -project . -format json
```

Stage 1 treats `fails` and `why` as equivalent commands. The distinction is reserved for future UX where `fails` may list failures and `why` may include richer causality evidence.

Exit codes:

- `0`: analysis completed and no findings were found.
- `1`: analysis completed with Maven failure findings, or analysis failed.
- `2`: invalid CLI usage.

## Library Usage

```go
package main

import (
	"fmt"

	"github.com/Will-thom/pr-maven-cli/pkg/prmaven"
)

func main() {
	report, err := prmaven.Analyze(prmaven.Options{ProjectDir: "."})
	if err != nil {
		panic(err)
	}
	for _, finding := range report.Findings {
		fmt.Println(finding.ReproduceCommand)
	}
}
```

The public contract is centered on:

- `Report`
- `Summary`
- `Module`
- `Finding`

These structures are intentionally simple so other tools can consume the analyzer without depending on the CLI.

A runnable library example lives in [examples/library](examples/library).

## Demo Project

The demo fixture lives at:

```text
demo/multi-module-failure
```

It contains a Maven aggregator project with two modules:

- `payment-core`, with a Surefire failure.
- `payment-api`, with a Failsafe error.

The reports are intentionally committed under `target/*-reports` because they are test fixtures.

## Why This Exists

Generic PR tools can show that CI failed. Generic agents can inspect logs when prompted. Build observability platforms can offer deep build analysis.

PR Maven CLI focuses on a narrower production workflow:

```text
Turn Maven CI failures into deterministic PR context.
```

This is useful for:

- Java teams with Maven multi-module repositories.
- OSS maintainers triaging failing PRs.
- Engineering platform teams.
- Regulated teams that cannot send logs to external services.
- Internal agent platforms that need structured local evidence.

## Project Documents

- [Manifesto](MANIFESTO.md)
- [Roadmap](ROADMAP.md)
- [Installation](docs/installation.md)
- [Usage guide](docs/usage.md)
- [Examples](examples/README.md)
- [Integrations](docs/integrations.md)
- [Contributing](CONTRIBUTING.md)
- [Permission posture](docs/permissions.md)
- [Testing](docs/testing.md)
- [CI/CD](docs/ci.md)
- [Release process](docs/release.md)
- [Governance](GOVERNANCE.md)
- [Maintainers](MAINTAINERS.md)
- [Security](SECURITY.md)

## Contributing

The project is designed to accept many focused contributions.

Good contribution areas:

- More Maven report parsers.
- More fixtures.
- Better CLI ergonomics.
- JSON schema documentation.
- CI examples.
- Maven edge cases.

Read [CONTRIBUTING.md](CONTRIBUTING.md) before opening a PR.

## Founder

PR Maven CLI was founded by Will-thom.

## License

Apache-2.0.
