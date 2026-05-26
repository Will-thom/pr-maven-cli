# Testing

PR Maven CLI treats tests as part of the public product contract.

The Stage 1 test system is designed for:

- contributors working locally;
- maintainers reviewing pull requests;
- automated coding agents proposing scoped changes;
- CI runs on Linux, Windows, and macOS.

## Local Commands

Run the standard suite:

```bash
go test ./...
```

On Windows PowerShell:

```powershell
.\scripts\test.ps1
.\scripts\test.ps1 -Race -Coverage
```

On Unix-like shells:

```bash
./scripts/test.sh
PRMAVEN_RACE=1 PRMAVEN_COVERAGE=1 ./scripts/test.sh
```

With `make`:

```bash
make test
make test-race
make coverage
make coverage-check
make ci
```

## Test Layers

### Library Tests

Package: `pkg/prmaven`.

Coverage includes:

- Maven module discovery from `pom.xml`;
- Surefire report parsing;
- Failsafe report parsing;
- Checkstyle report parsing;
- SpotBugs report parsing;
- Maven Enforcer log parsing;
- report-to-module mapping;
- reproduction command generation;
- JSON output contract;
- text output snapshots;
- missing project error behavior.

### CLI End-to-End Tests

Package: `cmd/prmaven`.

Coverage includes:

- `fails` text output;
- `why` JSON output;
- non-zero exit when findings exist;
- zero exit when no findings exist;
- invalid command and invalid format handling.

### Demo Fixtures

Fixtures live under `demo/`.

- `demo/multi-module-failure`: Maven aggregator with Surefire and Failsafe findings.
- `demo/no-failure`: Maven aggregator with passing Surefire report output.
- `pkg/prmaven/testdata/checkstyle-project`: Maven aggregator with a sanitized Checkstyle report fixture.
- `pkg/prmaven/testdata/spotbugs-project`: Maven aggregator with a sanitized SpotBugs report fixture.
- `pkg/prmaven/testdata/enforcer-project`: Maven aggregator with a sanitized Maven Enforcer log fixture.

The `target/*-reports` directories and selected `target/*.log` files are intentionally versioned because they are stable test fixtures, not local build output.

### Golden Files

Golden files live under:

```text
pkg/prmaven/testdata/golden
```

They protect human-readable output from accidental changes. When output changes intentionally, update the golden file in the same PR and explain the reason.

## CI

GitHub Actions runs:

- Go tests on Linux, Windows, and macOS;
- Go 1.22.x and the current stable Go release;
- race detector on Linux;
- coverage generation on Linux;
- a minimum total coverage gate of 70%.

The CI workflow is intentionally dependency-light. Core tests do not require Maven, network services, GitHub tokens, Docker, or external APIs.

## Contributor Expectations

For parser changes:

- add a sanitized fixture;
- add a focused unit test;
- update golden files if text output changes;
- keep JSON fields stable unless the issue explicitly allows a compatibility change.

For CLI changes:

- add or update end-to-end tests;
- document exit code changes;
- update README examples if user-facing commands change.

For documentation-only changes:

- tests are usually not required;
- keep examples aligned with the current CLI behavior.
