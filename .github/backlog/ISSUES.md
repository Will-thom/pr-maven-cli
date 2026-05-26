# Contributor Backlog

This backlog mirrors the planned GitHub issues for Stage 2 and Stage 3.

Each item is intentionally scoped so a contributor or automated coding agent can implement it with clear acceptance criteria.

## Stage 2 - Contributor Growth and Maven Signal Expansion

### 1. Add Checkstyle XML parser support

Labels: `stage: 2`, `area: parser`, `help wanted`, `good first issue`, `agent-friendly`

Acceptance criteria:

- Add a sanitized Checkstyle XML fixture.
- Produce a `Finding` with module path, report path, plugin name, message, and confidence reasons.
- Add unit tests for parser behavior.
- `go test ./...` passes.

### 2. Add SpotBugs XML parser support

Labels: `stage: 2`, `area: parser`, `help wanted`, `agent-friendly`

Acceptance criteria:

- Add a SpotBugs XML fixture.
- Map the bug instance to a Maven module.
- Include bug category/type and source file when available.
- `go test ./...` passes.

### 3. Extract Maven Enforcer failures from log fixtures

Labels: `stage: 2`, `area: parser`, `help wanted`, `agent-friendly`

Acceptance criteria:

- Add a sanitized Maven log fixture containing an Enforcer failure.
- Detect `maven-enforcer-plugin`.
- Emit phase/plugin/message context.
- Keep parsing deterministic and fixture-driven.

### 4. Extract JaCoCo threshold failures

Labels: `stage: 2`, `area: parser`, `help wanted`

Acceptance criteria:

- Add fixture coverage for a JaCoCo threshold failure.
- Detect `jacoco-maven-plugin`.
- Include threshold context in the finding message.
- Add focused tests.

### 5. Add parser registry abstraction

Labels: `stage: 2`, `area: architecture`, `help wanted`

Acceptance criteria:

- Introduce a small parser interface.
- Register Surefire and Failsafe through the same mechanism.
- Preserve current JSON output.
- Avoid speculative plugin framework design.

### 6. Add nested module discovery tests

Labels: `stage: 2`, `area: test`, `good first issue`, `agent-friendly`

Acceptance criteria:

- Add a nested Maven module fixture.
- Verify recursive module discovery.
- Verify report-to-module mapping for nested paths.
- `go test ./...` passes.

### 7. Add Windows and Linux path normalization tests

Labels: `stage: 2`, `area: test`, `good first issue`, `agent-friendly`

Acceptance criteria:

- Add tests for Windows-style and POSIX-style paths.
- Verify JSON output always uses slash-separated module/report paths.
- Avoid OS-specific test failures.

### 8. Document the JSON contract

Labels: `stage: 2`, `area: docs`, `good first issue`, `agent-friendly`

Acceptance criteria:

- Add `docs/json-contract.md`.
- Document all fields in `Report`, `Module`, and `Finding`.
- Include one complete JSON example from the demo.
- Explain compatibility expectations.

### 9. Add JSON schema for report output

Labels: `stage: 2`, `area: docs`, `area: test`, `help wanted`

Acceptance criteria:

- Add `schema/prmaven-report.schema.json`.
- Ensure the schema matches the current `Report` contract.
- Add documentation showing how CI tools can validate output.

### 10. Add golden tests for text output

Labels: `stage: 2`, `area: test`, `good first issue`, `agent-friendly`

Acceptance criteria:

- Add a golden fixture for the demo text output.
- Add a test that compares current output to the golden file.
- Keep output stable and intentional.

### 11. Add GitHub Actions CI for Go tests

Labels: `stage: 2`, `area: ci`, `good first issue`, `agent-friendly`

Acceptance criteria:

- Add a workflow that runs `go test ./...`.
- Use a maintained Go version.
- Run on pull requests and pushes to `main`.
- Keep the workflow minimal.

### 12. Add release workflow for tagged builds

Labels: `stage: 2`, `area: ci`, `help wanted`

Acceptance criteria:

- Add a release workflow triggered by `v*` tags.
- Build binaries for Linux, macOS, and Windows.
- Attach checksums.
- Document the release process.

### 13. Add fixture contribution guide

Labels: `stage: 2`, `area: docs`, `good first issue`, `agent-friendly`

Acceptance criteria:

- Add `docs/fixtures.md`.
- Explain how to sanitize Maven reports and logs.
- Explain where fixtures should live.
- Include examples for Surefire and Failsafe.

### 14. Add no-failure demo fixture

Labels: `stage: 2`, `area: test`, `good first issue`, `agent-friendly`

Acceptance criteria:

- Add a fixture project with passing reports.
- Verify CLI exits `0`.
- Verify text output says no failures were found.
- Verify JSON has `findingCount: 0`.

### 15. Add CLI module filter

Labels: `stage: 2`, `area: cli`, `help wanted`

Acceptance criteria:

- Add `-module` filter support.
- Limit findings to a selected Maven module path or module artifact ID.
- Add tests for matching and no-match cases.

### 16. Add output file option

Labels: `stage: 2`, `area: cli`, `good first issue`, `agent-friendly`

Acceptance criteria:

- Add `-output <path>`.
- Write text or JSON output to the selected file.
- Keep stdout behavior unchanged when `-output` is absent.
- Add tests.

### 17. Add confidence documentation

Labels: `stage: 2`, `area: docs`, `good first issue`

Acceptance criteria:

- Document what `high`, `medium`, and `low` confidence mean.
- Explain why Stage 1 currently reports high confidence for report-backed findings.
- Link the documentation from the README.

### 18. Add Maven 3.9 compatibility fixture notes

Labels: `stage: 2`, `area: docs`, `good first issue`

Acceptance criteria:

- Document Maven 3.9.x as the production baseline.
- Mention Maven 3.9.16 as the documented current baseline.
- Explain why Maven 4 is tracked separately.

### 19. Improve CLI help output

Labels: `stage: 2`, `area: cli`, `good first issue`, `agent-friendly`

Acceptance criteria:

- Add clear usage examples to `-h` output.
- Document commands and flags.
- Keep implementation dependency-free.

### 20. Add maintainer issue labeling guide

Labels: `stage: 2`, `area: docs`, `good first issue`, `agent-friendly`

Acceptance criteria:

- Add a short maintainer doc for labels.
- Explain `good first issue`, `help wanted`, and `agent-friendly`.
- Include examples of well-scoped issues.

## Stage 3 - PR and CI Context Layer

### 21. Design GitHub adapter interface

Labels: `stage: 3`, `area: architecture`, `help wanted`

Acceptance criteria:

- Add a design document for optional GitHub adapters.
- Define interfaces for changed files and check runs.
- Keep network behavior outside the core analyzer.

### 22. Add GitHub changed files adapter

Labels: `stage: 3`, `area: github`, `help wanted`

Acceptance criteria:

- Add an optional adapter for PR changed files.
- Use dependency injection so tests can run without GitHub.
- Add fixture/mocked tests.

### 23. Add GitHub check runs adapter

Labels: `stage: 3`, `area: github`, `help wanted`

Acceptance criteria:

- Add an optional adapter for check run metadata.
- Avoid requiring tokens for local analyzer tests.
- Add mock tests for success and failure states.

### 24. Add PR-to-module relevance scoring

Labels: `stage: 3`, `area: github`, `area: architecture`, `help wanted`

Acceptance criteria:

- Combine changed files with Maven module paths.
- Emit relevance score or reasons.
- Keep scoring explainable and deterministic.

### 25. Design baseline comparison model

Labels: `stage: 3`, `area: architecture`, `help wanted`

Acceptance criteria:

- Add a design document for comparing PR findings against main/baseline findings.
- Define required inputs and failure modes.
- Avoid implementation until the model is reviewed.

### 26. Implement confidence model v2

Labels: `stage: 3`, `area: architecture`, `help wanted`

Acceptance criteria:

- Add confidence levels based on multiple evidence sources.
- Preserve confidence reasons in JSON.
- Add tests for report-only and PR-context-backed findings.

### 27. Add Markdown PR summary output

Labels: `stage: 3`, `area: cli`, `help wanted`, `agent-friendly`

Acceptance criteria:

- Add `-format markdown`.
- Include findings and reproduction commands.
- Keep output suitable for PR comments.

### 28. Add `prmaven explain` command

Labels: `stage: 3`, `area: cli`, `help wanted`

Acceptance criteria:

- Add `explain` as a richer diagnostic command.
- Preserve `fails` and `why` behavior.
- Add CLI tests.

### 29. Add `prmaven ci` command

Labels: `stage: 3`, `area: cli`, `area: ci`, `help wanted`

Acceptance criteria:

- Add a command optimized for CI workspaces.
- Output JSON by default or document the chosen behavior.
- Include examples for GitHub Actions.

### 30. Investigate SARIF output

Labels: `stage: 3`, `area: ci`, `help wanted`

Acceptance criteria:

- Add a short research document.
- Map `Finding` fields to SARIF concepts.
- Recommend whether SARIF belongs in the project.

### 31. Investigate GitHub annotation output

Labels: `stage: 3`, `area: ci`, `help wanted`

Acceptance criteria:

- Document whether GitHub workflow annotations are useful for Maven findings.
- Include examples and limitations.
- Recommend next implementation steps.

### 32. Research GitLab merge request support

Labels: `stage: 3`, `area: gitlab`, `help wanted`

Acceptance criteria:

- Add a design note for GitLab MR support.
- Identify APIs needed for changed files and pipeline jobs.
- Keep GitLab support optional.

### 33. Add Maven 4 compatibility investigation

Labels: `stage: 3`, `area: maven`, `help wanted`

Acceptance criteria:

- Document Maven 4 report compatibility risks.
- Add fixtures only when Maven 4 behavior is stable enough.
- Do not declare production support prematurely.

### 34. Add agent evidence bundle output

Labels: `stage: 3`, `area: agent`, `help wanted`

Acceptance criteria:

- Design a JSON bundle for coding agents.
- Include findings, commands, confidence reasons, and relevant files when available.
- Avoid prompt-specific coupling.

### 35. Add CI artifact directory option

Labels: `stage: 3`, `area: cli`, `area: ci`, `help wanted`

Acceptance criteria:

- Add a flag for scanning a CI artifact directory separate from the source root.
- Preserve module mapping behavior where possible.
- Add tests with fixture artifact layout.

### 36. Add privacy guide for CI logs

Labels: `stage: 3`, `area: docs`, `good first issue`, `agent-friendly`

Acceptance criteria:

- Document how to sanitize logs before opening issues.
- Explain what data should never be pasted publicly.
- Link from `SECURITY.md` and `CONTRIBUTING.md`.

### 37. Add GitHub Action usage example

Labels: `stage: 3`, `area: docs`, `area: ci`, `help wanted`

Acceptance criteria:

- Add a documented example workflow using PR Maven CLI.
- Show text and JSON output modes.
- Avoid requiring a token for local report parsing.

### 38. Add internal engineering platform integration guide

Labels: `stage: 3`, `area: docs`, `help wanted`

Acceptance criteria:

- Document how platform teams can consume JSON output.
- Include examples for local/self-hosted execution.
- Explain privacy and no-telemetry behavior.

### 39. Add package manager distribution research

Labels: `stage: 3`, `area: release`, `help wanted`

Acceptance criteria:

- Research Homebrew, Scoop, npm wrapper, and direct binary releases.
- Recommend a first distribution path.
- Keep release complexity proportional to project maturity.

### 40. Add stable output compatibility policy

Labels: `stage: 3`, `area: docs`, `help wanted`

Acceptance criteria:

- Define compatibility rules for CLI output and JSON fields.
- Explain deprecation expectations before `v1.0.0`.
- Link from README and JSON contract docs.
