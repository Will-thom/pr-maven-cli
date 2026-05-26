# Permission Posture

This document records the intended repository permission model for PR Maven CLI.

## Current Stabilization Mode

The repository is public during the initial stabilization phase.

This keeps GitHub branch protection available without requiring a paid plan.

Do not switch it to private if that would disable branch protection or repository rules. A private stabilization phase should only be used when the account or organization plan supports the same branch protection controls used in public OSS mode.

Current maintainer policy:

- William Thomaz, GitHub: `@Will-thom`, is the founder and principal maintainer.
- No external collaborator should receive write, maintain, or admin access during stabilization unless intentionally promoted.
- External contributors can open issues and pull requests, but merge rights should stay limited to maintainers.

## Public OSS Mode

The public repository must keep `main` protected before accepting outside contributions.

Required public-mode controls:

- Require the `All CI checks` status before merge.
- Enforce the rule for administrators.
- Require conversation resolution before merge.
- Disable force pushes.
- Disable branch deletion.
- Keep merge rights limited to users with write, maintain, or admin access.
- Keep `@Will-thom` as the repository-wide code owner through `.github/CODEOWNERS`.

## GitHub Actions Permissions

Repository-level workflow token default:

- `contents: read`.
- Workflows cannot approve pull request reviews.

Workflow-specific policy:

- `CI` uses `contents: read`.
- `Security` uses `contents: read`, `security-events: write`, and `pull-requests: read`.
- `Release` defaults to `contents: read`; only the release publishing job receives `contents: write`.
- `Thank Contributor` uses `pull_request_target` only to read the base repository template and write a comment. It must not check out or execute contributor code.

## Secrets And External Integrations

The project should avoid repository secrets for the Stage 1 local-first MVP.

Expected empty surfaces during stabilization:

- GitHub Actions secrets.
- GitHub Actions variables.
- Webhooks.
- Environments.
- Deployments.
- GitHub Pages.

## OSS Contributor Readiness

Before actively inviting OSS contributors:

1. Validate branch protection for `main`.
2. Confirm `All CI checks` is the required status gate.
3. Confirm `CODEOWNERS` still routes ownership to `@Will-thom`.
4. Confirm no unintended collaborators, secrets, variables, webhooks, environments, deployments, releases, or pages were added.
5. Re-run CI and Security on `main`.
