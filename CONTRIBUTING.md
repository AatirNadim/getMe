# Contributing to getMe

First off, thank you for your interest in contributing to **getMe**! The general idea is that contributing to this project be easy and transparent as possible, whether it's:

- Reporting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features
- Becoming a maintainer

In order to keep our issues and pull requests highly navigable, our version history clean, and to prevent repository clogging, adhering to the following guidelines is strongly advised.

## 1. Issues: Reporting Bugs & Proposing Features

GitHub issues are used to track public bugs and requests.

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

_Note: If you are contributing to one of the specific SDKs (Java, JS, Python), refer to their respective directories under `/sdks/`._

### MCP Server

The MCP (Model Context Protocol) server is built with Python and uses `uv` for lightning-fast dependency management. To set up:

1. Install [`uv`](https://docs.astral.sh/uv/getting-started/installation/).
2. Navigate to the directory: `cd mcp-server`
3. Install dependencies: `uv sync --all-extras --dev`

## 3. Git Workflow & Commit Guidelines

To ensure a navigable and clean version history, branch naming and commit message conventions are strictly enforced.

### Branch Naming

Never commit directly to `main`. Create a branch from `main` using the following convention:

- `feature/description-of-feature` (e.g., `feature/batch-put-optimization`)
- `bugfix/issue-number-description` (e.g., `bugfix/12-compaction-race-condition`)
- `docs/description`

### Commit Messages

[Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specifications are followed here. This allows us to auto-generate changelogs and easily parse the commit history.

- `feat:` (new feature for the user, not a new feature for build script)
- `fix:` (bug fix for the user, not a fix to a build script)
- `docs:` (changes to the documentation)
- `style:` (formatting, missing semi colons, etc; no production code change)
- `refactor:` (refactoring production code, eg. renaming a variable)
- `test:` (adding missing tests, refactoring tests; no production code change)
- `chore:` (updating grunt tasks etc; no production code change)

**Example:**
`feat(server): implement asynchronous compaction thread`

### MCP Server Versioning & Releases

The `mcp-server` package utilizes dynamic versioning and automated releases via [Python Semantic Release](https://python-semantic-release.readthedocs.io/).

When changes are merged into the `main` branch, the CI pipeline (`.github/workflows/mcp-server-ci.yml`) automatically:

1. Parses your Conventional Commits to calculate the next semantic version number (e.g., `fix:` -> PATCH, `feat:` -> MINOR, `feat!:` or `BREAKING CHANGE:` -> MAJOR).
2. Tags the release dynamically.
3. Injects the new version into the build (via `SETUPTOOLS_SCM_PRETEND_VERSION`) and updates `server.json`.
4. Publishes the new version to PyPI and the MCP Registry.

**Important:** Do NOT manually bump versions in `pyproject.toml` or `server.json`. Rely entirely on correct Conventional Commit messages.

### SDK Versioning & Releases (Ephemeral Release Structure)

The SDKs located in the `sdks/` directory are published with manual release triggers (for bump-type) via a highly sophisticated **Ephemeral Release Structure**.

To keep the `main` branch impeccably clean from meaningless "bump version to X" commits, we treat versioning as a pure deployment concern. **Just by providing a bump-type (`major`, `minor`, `patch`), the CI/CD pipeline autonomously orchestrates the entire release lifecycle.**

The gh-actions pipelines in this repo, check out the code on a **detached HEAD**, dynamically calculate the next semantic version from prior tags, mutate the necessary configuration files (like `package.json`, `build.gradle.kts`, or `pyproject.toml`), and tag that specific detached commit. From there, it automatically **generates a changelog**, **drafts a GitHub Release**, and **publishes** the build to the **public registries** (npm, PyPI, Maven Central). These version bumps are **never merged back** into `main`.

**Important:** Just like the MCP server, **do NOT manually bump versions** for any SDK in your Pull Requests. Leave the base placeholder versions exactly as they are. 

If you are curious to see how this powerful automation works, you can explore the associated workflows:
- Core Versioning & Tagging: [`.github/workflows/reusable-sdk-release.yml`](.github/workflows/reusable-sdk-release.yml)
- Changelogs & Releases: [`.github/workflows/reusable-release-drafter.yml`](.github/workflows/reusable-release-drafter.yml)
- Registry Publishing: [`.github/workflows/publish-js-sdk.yml`](.github/workflows/publish-js-sdk.yml), etc.

You can read more about this architecture in the [SDKs README](./sdks/README.md#advanced-release--versioning-architecture).

## 4. Pull Requests

When you are ready to submit your code, open a Pull Request (PR) against the `main` branch.

To prevent PR clogging and ensure rapid reviews:

1. **Scope:** Keep PRs small and focused on a single issue or feature. If you have multiple unrelated changes, break them up into multiple PRs.
2. **Tests & Checks:** All code changes must be accompanied by relevant tests.

   **For the Go Server (`/server`) & CLI:**
   Run the benchmark suite to ensure no performance regressions:

   ```bash
   go test -bench . ./...
   ```

   **For the MCP Server (`/mcp-server`):**
   Ensure your code passes formatting, linting, and tests:

   ```bash
   cd mcp-server
   uv run ruff format --check .
   uv run ruff check .
   uv run pytest
   ```

3. **Format & Lint (Go):** Ensure your Go code is formatted (`go fmt`) and passes standard Go linters (`go vet`, `golangci-lint` if available) before submitting.
4. **Draft PRs:** If you want feedback on a work-in-progress, open your PR as a **Draft**.
5. **Rebasing:** If your branch falls behind `main`, prefer `git rebase` over `git merge` to keep the history linear.
6. **PR Description:** Reference any relevant issue numbers (e.g., "Fixes #123"). Explain what you changed and why.

## 5. Review Process

- Maintainers will review your PR. They may ask for changes or clarifications.
- Once approved, maintainers will generally use "Squash and merge" to combine your commits into a single, clean commit on the `main` branch, utilizing your PR title as the commit message.

Thank you again for contributing!
