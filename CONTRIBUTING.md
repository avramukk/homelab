# Contributing

This is a personal showcase repository, but it is written to production
standards on purpose. These conventions keep the history readable.

## Commit style

[Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <summary>

<why, not what>
```

Types: `feat` `fix` `docs` `refactor` `test` `chore` `ci`.
Scopes used here: `adr` `docs` `cluster` `gitops` `observability` `security` `terraform` `apps`.

- One logical change per commit (an ADR is one commit).
- History stays **linear** — rebase, no merge commits.

## Decision records

- Any decision that is expensive to reverse gets an **ADR** in [`docs/adr/`](docs/adr/).
- ADRs are append-only. Reversing a decision means a new ADR that supersedes it.

## Pull requests

- One concern per PR, with the [PR template](.github/pull_request_template.md).
- Every PR states **Decision / Alternatives / Validation / Rollback**.
- Docs change in the same PR as the system change.

## Security

- **No secrets in Git.** Sealed Secrets ciphertext only.
- Never commit `talosconfig`, machine secrets, `kubeconfig`, or `.env` files.
- Sign your commits (SSH signing is configured in this repo).
- **Scan before pushing:** the `gitleaks` workflow runs on every push/PR and
  weekly; to check locally, run `gitleaks git .` from the repo root.
