# CI/CD

PR Maven CLI uses a Stage 1 OSS-style GitHub pipeline.

The pipeline is intentionally dependency-light. Core checks do not require Maven, Docker, private credentials, hosted services, or external test infrastructure.

## CI Workflow

File: `.github/workflows/ci.yml`

Runs on:

- pull requests;
- pushes to `main`.

Jobs:

- `Quality gate`: `gofmt`, `go vet`, and unit tests.
- `Go tests`: Linux, Windows, macOS, Go 1.22.x, and current stable Go.
- `Race detector`: `go test -race ./...` on Linux.
- `Coverage gate`: coverage profile with a 70% total coverage floor.
- `Build`: cross-platform binary builds for Linux, macOS, and Windows.
- `CLI smoke test`: exercises the compiled binary against demo fixtures.
- `All CI checks`: stable aggregate job for future branch protection.

## Security Workflow

File: `.github/workflows/security.yml`

Runs on:

- pull requests;
- pushes to `main`;
- weekly schedule;
- manual dispatch.

Jobs:

- `Go vulnerability check`: runs `govulncheck`.
- `CodeQL`: static analysis for Go.
- `Dependency review`: reviews dependency changes on pull requests.

## Release Workflow

File: `.github/workflows/release.yml`

Runs on:

- tags matching `v*`;
- manual dispatch for package validation.

Release artifacts:

- Linux amd64 and arm64 tarballs.
- macOS amd64 and arm64 tarballs.
- Windows amd64 zip.
- SHA-256 checksum files.

The tag version is embedded in the CLI through:

```text
prmaven version
```

## Local Parity

Before opening a PR, contributors should run:

```bash
sh scripts/quality.sh
PRMAVEN_COVERAGE=1 sh scripts/test.sh
sh scripts/build.sh
```

On Windows PowerShell:

```powershell
.\scripts\quality.ps1
.\scripts\test.ps1 -Coverage
.\scripts\build.ps1
```

## Branch Protection Recommendation

When branch protection is enabled, use `All CI checks` as the required CI status. Keep security checks visible, but avoid making scheduled security tooling a blocker for focused contributor PRs until the project has more maintainers.
