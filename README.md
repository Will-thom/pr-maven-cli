# PR Maven CLI

Maven-aware PR and CI triage for Java teams.

This repository preserves the product context, market reasoning, and initial direction for a possible OSS developer tool focused on pull request operations in Java/Maven projects.

The current goal is not to implement code yet. The goal is to keep a versioned record of the idea so it can be resumed later with clearer product and technical decisions.

## Core Thesis

Modern software teams are increasingly productive, but human review and CI failure triage remain bottlenecks.

Generic PR tools can show that a pull request failed CI. Generic agents can inspect logs when prompted. But Java/Maven teams still need a deterministic, local or self-hosted tool that answers:

```text
What exactly did this PR break in the Maven build?
Which module failed?
Which Maven plugin failed?
Which test or rule failed?
What is the smallest local command to reproduce it?
Who needs to act now?
```

PR Maven CLI should live in the practical layer between GitHub/GitLab CLI tools, CI logs, build observability platforms, and AI coding agents.

## Positioning

PR Maven CLI is not:

- An AI CI repair agent.
- A dashboard-first product.
- A generic GitHub CLI replacement.
- A simple test reporter.
- A tool for finding OSS issues for contributors.

PR Maven CLI is:

- A Maven-aware CI failure context layer.
- Terminal-first.
- Deterministic by default.
- Offline/self-hosted friendly.
- Useful for humans and agents.
- Focused on PR, Maven module, failing build signal, and reproducible command.

Short positioning:

```text
Turn Maven CI failures into actionable PR context.
```

Longer positioning:

```text
PR Maven CLI gives humans and agents the same Maven-aware build context.
```

## Why This Matters

This is especially relevant for:

- Java teams with large Maven projects.
- Multi-module Maven repositories.
- Regulated companies such as banks, insurers, healthcare, and enterprise software vendors.
- Teams that cannot send source code or CI logs to external LLM providers.
- Teams that need deterministic CI diagnostics without prompts or manual sessions.
- Internal engineering platforms that want machine-readable build context.
- OSS maintainers who need faster, lower-friction PR triage.

The strongest product angle is not "replace Codex or Claude Code". The stronger angle is:

```text
Deterministic, cheap, auditable, local/CI-first, no prompt required, no vendor lock-in.
```

For teams with internal LLMs or tools like Ollama, PR Maven CLI can become an evidence provider instead of a competitor.

## Target Users

The product should support short role concepts instead of long labels:

- `owner`: someone responsible for review, merge, or release readiness.
- `author`: someone who opened or owns a PR.
- `team`: shared queue and repository health view.
- `mine`: user-specific action view.
- `queue`: review and action queue.
- `ready`: PRs likely ready for review or merge.
- `blocked`: PRs blocked by CI, conflict, review, or missing action.
- `fails`: PRs with failing CI/build signals.

The vocabulary should be practical and memorable, closer to `mvn clean` than verbose enterprise terminology.

## CLI Direction

The preferred user experience avoids Maven goal syntax with colons for day-to-day usage.

Possible command style:

```bash
pr queue
pr fails
pr broken
pr log
pr why
pr ready
pr blocked
pr mine
pr owner
pr author
```

Alternative binary name if `pr` is too generic:

```bash
prmaven queue
prmaven fails
prmaven broken
prmaven why
```

Maven plugin syntax such as `mvn pr:queue` may still be supported later for Maven-native integration, but the preferred CLI should be space-based and terminal-friendly.

## Example Output

```text
PR #4821 likely introduced a Maven test failure

Module: payment-core
Plugin: maven-surefire-plugin
Test: PaymentRoundingTest.shouldRoundHalfEven

Reproduce:
mvn -pl payment-core -Dtest=PaymentRoundingTest test

Confidence: high

Reason:
- Test passed on main
- Test fails on the PR merge commit
- PR changed files inside payment-core
- Failure appears in the surefire report
```

For CI or agent usage:

```json
{
  "pr": 4821,
  "status": "ci_failing",
  "module": "payment-core",
  "phase": "test",
  "plugin": "maven-surefire-plugin",
  "failure": "PaymentRoundingTest.shouldRoundHalfEven",
  "reproduce": "mvn -pl payment-core -Dtest=PaymentRoundingTest test",
  "confidence": "high"
}
```

## Maven-Aware Scope

Initial Maven signals worth supporting:

- Maven reactor and multi-module projects.
- `pom.xml` module mapping.
- Surefire unit test reports.
- Failsafe integration test reports.
- Checkstyle failures.
- SpotBugs failures.
- Maven Enforcer failures.
- JaCoCo threshold failures.
- Build phase and plugin detection.
- Minimal reproduction command generation.

The core product value is understanding enough Maven structure to produce useful, reproducible PR context.

## Architecture Sketch

Potential components:

- Local CLI: terminal commands such as `pr fails`, `pr why`, `pr queue`.
- CI integration: GitHub Action, GitLab CI step, or generic CI script.
- Maven engine: parses project structure, modules, reports, plugin outputs, and build metadata.
- SCM adapter: GitHub first, GitLab later.
- PR context adapter: changed files, checks, reviews, author, labels, mergeability.
- Evidence model: deterministic structured context with confidence levels.
- Output formats: human-readable text, Markdown, JSON.

The product should work without AI. AI can be optional later.

## Competition Notes

Validated adjacent competitors:

- GitHub CLI: strong generic PR and CI status commands, but not Maven-aware.
- GitLab CLI: strong generic MR operations, but not Maven-aware.
- Develocity Build Scan for Maven: deep build observability, but not a simple terminal-first PR triage layer.
- TestGlance and other test reporters: good PR reporting, but not Maven reactor plus minimal command reproduction as the main product.
- gha-failure-analysis: open source CI failure analysis with GitHub Actions and optional local LLM/Ollama support, but not specifically Maven CLI first.
- Daxtack: enterprise/on-prem CI failure analysis, but not positioned as a lightweight Maven-first CLI.
- WarpFix: CI repair agent with Maven/Gradle support, but more agentic auto-fix than deterministic Maven evidence layer.

Current conclusion:

```text
Direct competition in the exact niche appears low.
Adjacent competition exists and validates the problem.
```

Exact niche:

```text
Maven-aware + terminal-first + deterministic + offline/self-hosted + PR + Maven module + reproduction command.
```

## License Direction

Apache-2.0 is the preferred license direction.

Reasons:

- Familiar in Java and enterprise ecosystems.
- Comfortable for corporate adoption.
- Includes explicit patent language.
- Fits Maven, Apache, Spring, and broader JVM culture.

MIT is simpler, but Apache-2.0 likely sends a better signal for this category.

## Sponsor Strategy

The short-term goal is not fame or becoming a "rockstar".

The goal is sponsorship and credibility through practical utility.

Potential sponsor audiences:

- Java OSS maintainers.
- Engineering platform teams.
- Backend teams with Maven monorepos or large multi-module repositories.
- Companies that need local/self-hosted tooling.
- Teams with internal AI agents that need deterministic build context.

Sponsor pitch:

```text
Support an OSS tool that reduces wasted time in Maven CI failure triage.
```

Realistic expectation:

- GitHub Sponsors is slow.
- First dollar is plausible only after a usable MVP, examples, docs, and distribution.
- A realistic first-sponsor window may be 3 to 6 months with consistent execution.
- Service/consulting may monetize earlier than Sponsors.

## Monetization Paths

OSS-first:

- GitHub Sponsors.
- Corporate sponsorship.
- Sponsored features.
- Support for public OSS maintainers.

Service-first:

- Consulting for Java/Spring/Maven CI bottlenecks.
- Build failure triage.
- Maven multi-module optimization.
- CI reproducibility improvements.
- Internal tooling for regulated teams.

Future commercial options:

- Hosted dashboard.
- Team policy rules.
- Historical CI failure analytics.
- Private repository support packages.
- Enterprise support.
- Integrations for internal LLM/agent workflows.

## MVP Candidate

An initial MVP could avoid network complexity and start from local CI artifacts:

1. Parse a Maven project.
2. Detect modules.
3. Parse Surefire and Failsafe reports.
4. Identify failing test class and method.
5. Map failure to Maven module.
6. Generate minimal reproduction command.
7. Output text and JSON.

Possible first command:

```bash
pr fails
```

Possible MVP output:

```text
Module: core
Failure: UserServiceTest.shouldRejectInvalidEmail
Reproduce: mvn -pl core -Dtest=UserServiceTest test
```

After that, add PR and CI context.

## Open Questions

- Should the first implementation be Java, Kotlin, Go, Rust, or Node?
- Should the primary binary be `pr`, `prmaven`, or another short name?
- Should Maven plugin support be first-class or secondary?
- Should GitHub support come before generic local artifact parsing?
- Should the project begin as a CLI only, then add CI integration?
- How much AI support should be optional in v1?
- What exact evidence is required before saying "likely caused by this PR"?

## Current Product Decision

The strongest direction is:

```text
PR Maven CLI: deterministic Maven CI failure context for humans and agents.
```

The product should avoid competing directly with GitHub CLI, Develocity, or AI repair agents.

It should focus on a narrow, painful, production-first workflow:

```text
Tell me what failed in this Maven PR, where it failed, and how to reproduce it locally.
```

