# Contributing to PR Maven CLI

Thank you for considering a contribution.

PR Maven CLI is intentionally designed as a contributor-friendly OSS project. The roadmap is split into focused issues so humans and automated coding agents can make useful progress without needing deep project context.

## Baseline

- Language: Go.
- Runtime dependencies: none for the stage 1 core.
- Target Maven baseline: Maven 3.9.x, documented against Maven 3.9.16.
- Maven 4 support is planned after the Maven 4 line is production-ready.
- License: Apache-2.0.

## Setup

```bash
git clone https://github.com/Will-thom/pr-maven-cli.git
cd pr-maven-cli
go test ./...
```

Run the demo:

```bash
go run ./cmd/prmaven fails -project demo/multi-module-failure
go run ./cmd/prmaven why -project demo/multi-module-failure -format json
```

The demo includes versioned Maven report fixtures under `demo/multi-module-failure/**/target/*-reports`.

## Contribution Rules

- Keep PRs focused.
- Include tests for behavior changes.
- Preserve stable JSON fields unless the issue explicitly allows a breaking change.
- Prefer deterministic parsing over guessing.
- Keep network access optional.
- Do not add telemetry.
- Do not require external services for core tests.
- Do not commit generated binaries, local caches, or private CI logs.

## Good First Contributions

Good first issues usually fit one of these areas:

- Add a parser fixture for a Maven plugin report.
- Improve output wording.
- Add edge-case coverage for Maven module discovery.
- Add documentation examples.
- Add one JSON field with a clear acceptance test.
- Improve error messages.

## Pull Request Checklist

- [ ] The PR solves one issue or one clearly scoped problem.
- [ ] `go test ./...` passes.
- [ ] New behavior is covered by tests or fixtures.
- [ ] Documentation is updated when user-facing behavior changes.
- [ ] JSON output compatibility is preserved or the breaking change is explicitly justified.

## Automated Contributions

Automated coding agents are welcome when the contribution is reviewable.

Expected standard:

- The PR title references the issue.
- The PR body explains the change in plain language.
- Tests are included.
- The diff is scoped.
- The contribution does not perform unrelated refactors.

## Maintainers

Maintainers are expected to:

- Keep issues small and actionable.
- Add acceptance criteria to issues.
- Label issues consistently.
- Review contributor PRs with clear, respectful feedback.
- Favor production reliability over feature volume.
