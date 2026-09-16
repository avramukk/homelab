# Homelab

Personal home lab, planned in public. This repository is where I figure out
what I want (backlog), decide the architecture, and — soon — provision the
actual infrastructure.

## Status

- **Phase:** planning (inbox → architecture). Decisions get recorded as they converge.
- **Structure today:** a raw idea backlog in [`INBOX.md`](./INBOX.md) and nothing else — no code yet.

## How it works

1. Drop an idea / link / technology into the chat with an agent — no context needed.
2. The agent summarizes it (read-only) and records it in [`INBOX.md`](./INBOX.md):
   category, status `inbox`, sources, open questions.
3. During planning rounds we triage: `inbox → researching → adopted/rejected`.
4. Adopted decisions feed the architecture document (to appear later).

## Security policy

- **No secrets in the repo.** [`/.gitignore`](./.gitignore) blocks `*.credentials*`, `.env*`, `secrets/`, `*.tfstate`, and similar.
- No plaintext tokens or passwords. GitHub access goes through the `gh` credential helper.