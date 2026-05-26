# Security Policy

PR Maven CLI is local-first and does not send source code, Maven reports, CI logs, or project metadata to external services.

## Supported Versions

Until `v1.0.0`, security fixes are accepted on the default branch.

## Reporting

Please report security concerns privately to the repository owner before opening a public issue.

Do not include private CI logs, credentials, tokens, proprietary source code, or customer data in public issues.

## Data Handling Principles

- No telemetry by default.
- No external network calls in the core analyzer.
- No credentials required for local report analysis.
- GitHub or GitLab integrations must remain optional.
