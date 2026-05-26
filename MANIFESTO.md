# PR Maven CLI Manifesto

PR Maven CLI exists to make Maven CI failures easier to understand, reproduce, and act on.

The project starts from one practical belief:

> A failed Maven build should produce useful PR context without requiring a dashboard, a prompt, or a proprietary service.

## Founder

PR Maven CLI was founded by Will-thom.

Will-thom is the public founder identity for this project. The founder's goal is to build a useful open source tool for Java and Maven teams, while creating a contributor-friendly project where focused issues can be picked up by humans, maintainers, and automated coding agents.

## Product Principles

- Deterministic before agentic.
- Local-first before hosted.
- Maven-aware before generic.
- Small, reproducible evidence before broad claims.
- Human-readable and machine-readable output from day one.
- Useful to individual contributors, maintainers, CI systems, and internal engineering platforms.

## What We Optimize For

- Clear failure context.
- Minimal reproduction commands.
- Stable JSON contracts.
- Low setup cost.
- Contributor-friendly issues.
- Production-safe defaults.
- Respect for private source code and private CI logs.

## What We Avoid

- Sending code or logs to external services by default.
- Replacing Maven, GitHub CLI, GitLab CLI, or build observability platforms.
- Claiming root cause without evidence.
- Auto-fixing code before explaining the failure.
- Large unfocused changes that are hard to review.

## Maintainer Culture

This project should be friendly to contributors who want scoped work:

- Every issue should have clear acceptance criteria.
- Test expectations should be explicit.
- New contributors should be able to land focused improvements.
- Automated PRs are welcome when they are readable, tested, and respectful of the project scope.
- Maintainers should prefer small, reviewable changes over large speculative rewrites.
