# Contributing to getMe

First off, thank you for your interest in contributing to **getMe**! We want to make contributing to this project as easy and transparent as possible, whether it's:

- Reporting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features
- Becoming a maintainer

In order to keep our issues and pull requests highly navigable, our version history clean, and to prevent repository clogging, we ask that you adhere to the following guidelines.

## 1. Issues: Reporting Bugs & Proposing Features

We use GitHub issues to track public bugs and requests. 

### Reporting Bugs
If you find a bug, please ensure an issue does not already exist. If it does not, create a new issue and include:
- **Environment:** OS, Go version, Docker version (if applicable).
- **Reproduction Steps:** Provide a minimal, reproducible example or the sequence of steps to trigger the bug.
- **Expected vs. Actual Behavior:** What did you expect to happen, and what actually happened?
- **Logs:** Any relevant logs or stack traces.

### Proposing Features
Feature requests are welcome! When proposing a new feature:
- Explain **why** this feature is needed. What use case does it solve?
- Provide a high-level overview of how you envision it being implemented.
- Discuss potential impacts on performance or existing APIs (especially since this is a high-performance key-value store).

## 2. Setting Up Your Development Environment

Please refer to the [README.md](README.md) for detailed instructions. The core server can be bootstrapped locally using:
```bash
# For local binaries
cd server && ./init-server-local.sh

# Or for Docker Compose
cd server && ./init-server-docker.sh
```
*Note: If you are contributing to one of the specific SDKs (Java, JS, Python), refer to their respective directories under `/sdks/`.*

## 3. Git Workflow & Commit Guidelines

To ensure a navigable and clean version history, we strictly enforce branch naming and commit message conventions.

### Branch Naming
Never commit directly to `main`. Create a branch from `main` using the following convention:
- `feature/description-of-feature` (e.g., `feature/batch-put-optimization`)
- `bugfix/issue-number-description` (e.g., `bugfix/12-compaction-race-condition`)
- `docs/description`

### Commit Messages
We follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification. This allows us to auto-generate changelogs and easily parse the commit history.
- `feat:` (new feature for the user, not a new feature for build script)
- `fix:` (bug fix for the user, not a fix to a build script)
- `docs:` (changes to the documentation)
- `style:` (formatting, missing semi colons, etc; no production code change)
- `refactor:` (refactoring production code, eg. renaming a variable)
- `test:` (adding missing tests, refactoring tests; no production code change)
- `chore:` (updating grunt tasks etc; no production code change)

**Example:**
`feat(server): implement asynchronous compaction thread`

## 4. Pull Requests

When you are ready to submit your code, open a Pull Request (PR) against the `main` branch. 

To prevent PR clogging and ensure rapid reviews:
1. **Scope:** Keep PRs small and focused on a single issue or feature. If you have multiple unrelated changes, break them up into multiple PRs.
2. **Tests:** All code changes must be accompanied by relevant unit and/or integration tests. Run the benchmark suite to ensure no performance regressions:
   ```bash
   go test -bench . ./...
   ```
3. **Format & Lint:** Ensure your Go code is formatted (`go fmt`) and passes standard Go linters (`go vet`, `golangci-lint` if available) before submitting.
4. **Draft PRs:** If you want feedback on a work-in-progress, open your PR as a **Draft**.
5. **Rebasing:** If your branch falls behind `main`, prefer `git rebase` over `git merge` to keep the history linear.
6. **PR Description:** Reference any relevant issue numbers (e.g., "Fixes #123"). Explain what you changed and why.

## 5. Review Process

- Maintainers will review your PR. They may ask for changes or clarifications.
- Once approved, maintainers will generally use "Squash and merge" to combine your commits into a single, clean commit on the `main` branch, utilizing your PR title as the commit message.

Thank you again for contributing!
