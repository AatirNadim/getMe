# Project Context for AI Agents (getMe)

This document provides foundational context, architectural boundaries, and guidelines for AI coding agents interacting with the `getMe` repository, structured to ensure maximum effectiveness when providing assistance.

## 1. Project Overview
`getMe` is a high-performance, local-first Key-Value Store. 

**Core Architecture:**
- **Server (`/server`):** The core daemon and storage engine handling memory/disk persistence.
- **Interfaces:** A CLI (`/cli`), SDKs (`/sdks`), and an HTTP proxy (`/http-proxy-go`) interact with the server.
- **Shared Code (`/commons`):** Routes, socket paths, and types are centralized here. If you modify an API route or limit, update it in `commons` to keep the CLI, server, and SDKs in sync.
- **Transport:** The server predominantly uses Unix Domain Sockets (default: `/tmp/getMeStore/sockDir/getMe.sock`) for fast local inter-process communication, with a REST-like interface (`GET`, `POST`, `DELETE`).

**Go Workspace Note:**
While local development uses a multi-module `go.work` workspace, **remote environments do not use `go.work`**. Be mindful of this context when generating build instructions or CI/CD steps.

## 2. Build and Test Commands
Agents should use the following commands to build and run the project:

- **Local Build & Run:** 
  Use `./server/init-server-local.sh` (or the root init scripts if symlinked) to build the Go binary and start the server with its logging stack. 
  *(Note: See Security Considerations regarding `sudo`)*
- **Docker Build & Run:** 
  Use `./server/init-server-docker.sh` (or `docker compose up --build`) for a fully containerized stack.
- **Standard Tests:** 
  `go test ./...` (Run within the specific module directory, e.g., `/server` or `/cli`).
- **Performance Benchmarks:** 
  `go test -bench . ./...`

## 3. Code Style Guidelines
- **Formatting:** All Go code must be formatted using standard `gofmt`.
- **Linting:** Ensure code passes `go vet` and `golangci-lint`. The project relies on the configuration found in `.golangci.yml`.
- **Idioms:** Mimic existing project conventions, particularly around concurrency (mutex locks) in the storage engine. 

## 4. Testing Instructions
- **Coverage:** All code changes (especially features and bug fixes) must be accompanied by relevant unit and/or integration tests.
- **Performance:** When touching the `/server/store` or hot paths, always verify there are no performance regressions by running the benchmark suite.
- **Verification:** Always execute tests proactively before declaring a task complete to ensure you have not broken the build or existing behaviors.

## 5. Security Considerations
- **Script Execution:** **NEVER run build/init scripts with `sudo`**. The scripts already call necessary setup helpers with elevated privileges. Running them as root causes local folders/files to be owned by `root`, breaking local development.
- **Data Protection:** Never introduce code that exposes, logs, or commits secrets, API keys, or raw user data. 
- **Socket Permissions:** Unix domain sockets are used for transport; ensure socket files maintain proper restrictive permissions to prevent unauthorized local access.

## 6. Extra Instructions

### Commit Message Guidelines
We strictly follow the **Conventional Commits** specification for auto-generating changelogs:
- `feat:` (new feature for the user)
- `fix:` (bug fix for the user)
- `docs:` (changes to documentation)
- `style:` (formatting, missing semicolons, etc.)
- `refactor:` (refactoring production code)
- `test:` (adding or refactoring tests)
- `chore:` (updating grunt tasks, configs, etc.)
- *Example:* `feat(server): implement asynchronous compaction thread`

### Branch & Pull Request Guidelines
- **Branch Naming:** Branch off `main` using `feature/description-of-feature`, `bugfix/issue-number-description`, or `docs/description`.
- **Rebasing:** Prefer `git rebase` over `git merge` to keep history linear.
- **PR Description:** Reference relevant issues (e.g., "Fixes #123"). Explain *what* was changed and *why*.
- **Merging:** Maintainers use "Squash and merge" upon approval.

### Deployment & CI
- Refer to `DOCKERHUB.md` for containerization details and volume mapping rules.
- Remote pipelines build the modules individually without a global `go.work` file. Ensure dependencies inside each module's `go.mod` remain consistent and self-contained.
