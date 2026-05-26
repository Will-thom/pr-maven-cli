# Integrations

PR Maven CLI is local-first in Stage 1.

The core analyzer does not call GitHub, GitLab, CI APIs, hosted services, or AI providers. It reads Maven report artifacts from the local filesystem and emits text or JSON.

## Stage 1 Integration Scope

Current first-party integration surface:

- GitHub repository automation for the project itself.
- GitHub Actions CI/CD workflows for this repository.
- A copyable GitHub Actions usage example for Maven failure triage.

Current product/runtime integration scope:

- No native GitHub API adapter yet.
- No native GitLab API adapter yet.
- No required provider token.
- No required network access for Maven report analysis.

This means GitHub is the only platform with official project automation and example coverage today, but the CLI itself remains provider-agnostic in Stage 1.

## Planned Native Adapters

Native PR and CI context adapters are planned for Stage 3.

Planned order:

1. GitHub adapter for changed files and check runs.
2. GitHub-oriented PR summary and CI workspace output.
3. GitLab merge request support investigation.

Adapters must stay optional. Local Maven report parsing should continue to work without GitHub or GitLab tokens.

## Contributor Guidance

Before adding a provider integration:

- Keep the provider client behind an interface.
- Avoid making network access part of core tests.
- Add fixtures or mocks for provider responses.
- Preserve the local analyzer contract.
- Document required tokens and permissions.
- Keep provider-specific behavior out of parser packages.

