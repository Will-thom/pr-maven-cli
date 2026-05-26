# PR Maven CLI Roadmap

This roadmap is intentionally conservative. Dates include safety margin for review, contributor onboarding, documentation, and stabilization.

Current production baseline: Maven 3.9.x, documented against Maven 3.9.16. Apache Maven currently lists Maven 3.9.16 as the recommended release and Maven 4.x as a preview line that is not safe for production use: <https://maven.apache.org/download.cgi>.

## Stage 1 - Production-Usable Local MVP

Target window: May 26, 2026 to June 7, 2026.

Target release: `v0.1.0`.

Goal: deliver a stable local CLI and Go library that can inspect Maven Surefire and Failsafe reports from a repository or CI workspace and produce actionable failure context.

Scope:

- Go library package with stable `Report`, `Module`, and `Finding` structures.
- CLI binary entrypoint: `prmaven`.
- Commands: `fails` and `why`.
- Text and JSON output.
- Maven multi-module discovery from `pom.xml`.
- Surefire report parsing.
- Failsafe report parsing.
- Module mapping from report paths.
- Minimal reproduction command generation.
- Demo Maven multi-module fixture.
- Demo no-failure fixture.
- Unit tests for parser and formatter behavior.
- End-to-end CLI tests.
- Golden tests for text output.
- GitHub Actions CI for Linux, Windows, macOS, Go 1.22.x, and stable Go.
- Race detector and coverage jobs.
- Quality, security, build, smoke, and release workflows.
- Dependabot maintenance for GitHub Actions and Go modules.
- Public README, manifesto, contribution guide, license, roadmap, issue templates.
- Repository visibility changed to public.

Acceptance criteria:

- `go test ./...` passes.
- Demo command returns findings from both Surefire and Failsafe fixtures.
- No-failure demo exits successfully and reports zero findings.
- JSON output is parseable and stable enough for CI/agent consumption.
- Text output gives module, plugin, phase, test, report path, confidence, and reproduction command.
- Text output is protected by golden tests.
- CI runs the Stage 1 suite across supported operating systems.
- Release tags produce Linux, macOS, and Windows packages with checksums.
- Security automation runs CodeQL, govulncheck, and dependency review.
- No network or external service is required for core functionality.
- GitHub is the only first-party platform with project automation and usage examples; no native provider API adapter ships in Stage 1.
- Stage 2 and Stage 3 issues exist in GitHub with contributor-friendly acceptance criteria.

Expected confidence: high.

## Stage 2 - Contributor Growth and Maven Signal Expansion

Target window: June 8, 2026 to July 19, 2026.

Target release: `v0.2.0`.

Goal: expand Maven failure coverage and make the project attractive for contributors by keeping tasks isolated, testable, and clearly documented.

Planned work:

- Add Checkstyle report support.
- Add SpotBugs report support.
- Add Maven Enforcer failure extraction from log fixtures.
- Add JaCoCo threshold failure extraction from report/log fixtures.
- Improve aggregator/root module handling.
- Add Windows/Linux path normalization tests.
- Add JSON schema documentation.
- Add golden snapshot tests for text output.
- Add GitHub Actions CI for Go tests.
- Add release workflow.
- Add fixture contribution guide.
- Add labels and issue taxonomy.
- Add maintainer guide.
- Add examples for CI artifacts.

Acceptance criteria:

- Each new signal type has fixture-based tests.
- Every issue is individually reviewable.
- CI runs on pull requests.
- Release process is documented.
- New contributors can add one fixture and one parser without understanding the entire codebase.

Expected confidence: high.

## Stage 3 - PR and CI Context Layer

Target window: July 20, 2026 to September 6, 2026.

Target release: `v0.3.0`.

Goal: connect local Maven evidence to pull request context while keeping the core local-first and deterministic.

Planned work:

- GitHub adapter for changed files and check runs.
- Optional GitLab adapter investigation.
- PR diff to module relevance scoring.
- Baseline/main comparison model.
- Confidence model v2.
- Markdown PR summary output.
- `prmaven explain` command.
- `prmaven ci` command for CI workspaces.
- SARIF or annotations investigation.
- Maven 4 compatibility investigation.
- Agent-consumable evidence bundle output.
- Documentation for internal engineering platforms.

Acceptance criteria:

- Network adapters are optional.
- Local report parsing remains usable without GitHub or GitLab tokens.
- PR context is additive evidence, not a replacement for Maven evidence.
- CI examples work with public repositories and private CI artifacts.
- Maven 4 support is not declared production-ready until validated against stable Maven 4 behavior.

Expected confidence: medium-high.

## Backlog Policy

The project should maintain many focused issues so contributors can participate quickly.

Issue quality standard:

- One problem per issue.
- Clear acceptance criteria.
- Suggested test files or fixture shape.
- Expected output impact.
- Labels that indicate area, difficulty, and type.
- No hidden dependency on private services.

Preferred labels:

- `good first issue`
- `help wanted`
- `area: parser`
- `area: cli`
- `area: docs`
- `area: test`
- `area: ci`
- `type: bug`
- `type: enhancement`
- `difficulty: focused`
- `agent-friendly`

## Release Naming

- `v0.1.x`: local Maven report context.
- `v0.2.x`: more Maven signals and contributor infrastructure.
- `v0.3.x`: PR and CI context.
- `v1.0.0`: stable CLI contract, documented JSON compatibility policy, and validated production usage.
