<div align="center">
  <img src="https://raw.githubusercontent.com/AatirNadim/getMe/main/getme-landing/public/extended-logo-rounded.png" alt="getMe Logo" style="width: 400px; max-width: 100%; margin-bottom: 20px;"/>
</div>

# getMe SDKs

This directory provides Software Development Kits (SDKs) for interacting with the `getMe` storage server from different programming languages. These SDKs abstract away the underlying API calls over Unix Domain Sockets, offering a simple, idiomatic interface for developers.

## Purpose

The goal of the SDKs is to make it easy to integrate `getMe` into various applications by providing a native, language-specific client library.

## Available SDKs

This project contains SDKs for the following languages:

- **[Go](./goSdk/)**: A native Go client for `getMe`.
- **[Java](./javaSdk/)**: A Java client, built with Gradle.
- **[JavaScript/TypeScript](./jsSdk/)**: A client library for Node.js environments.
- **[Python](./pythonSdk/)**: A Python client.

Each SDK directory contains its own source code, build files, and dependencies, making them independent and easy to package and distribute.

## General Design

While implementations vary by language, all SDKs follow a similar design pattern:

1. **Client Class/Struct**: Each SDK exposes a primary `Client` or `Service` object that manages the connection to the `getMe` server over a Unix Domain Socket (`.sock` file).
2. **Core Methods**: The `Client` provides methods that map directly to the server's core operations. Depending on the language, these generally include:
    - `get(key)`
    - `getJson(key)` (or language equivalent)
    - `batchGet(keys)`
    - `put(key, value)`
    - `putJson(key, json_value)` (or language equivalent)
    - `batchPut(pairs)`
    - `delete(key)`
    - `batchDelete(keys)`
    - `clearStore()`
3. **UDS Communication**: Internally, the SDKs use an HTTP client configured to communicate over Unix Domain Sockets (UDS) instead of standard TCP to the `getMe` server's endpoints. They handle request creation, serialization of data (e.g., to JSON), and deserialization of responses.
4. **Error Handling**: Errors from the server (e.g., "key not found") or network issues are translated into idiomatic error types or exceptions for the respective language.

## 🚀 Advanced Release & Versioning Architecture

The `getMe` repository employs a highly sophisticated **Ephemeral Release Structure** for managing and publishing SDKs. We treat versioning as a deployment concern rather than a source control concern. This ensures an impeccably clean monorepo history and completely eliminates "version bump" commit clutter in the main branch.

**Just by providing a bump-type (`major`, `minor`, `patch`), the CI/CD pipeline autonomously orchestrates the entire release lifecycle.** 

### The Magic Behind the Scenes

1. **Clean `main` Branch**: Source files in the `main` branch (such as `package.json`, `pyproject.toml`, or `build.gradle.kts`) are deliberately kept at a base or placeholder version. They are *never* permanently bumped.
2. **The Detached Commit Strategy**: When a new release is triggered via GitHub Actions, the pipeline checks out the repository into a **detached HEAD** state.
3. **Dynamic Orchestration**: 
   - The workflow analyzes the git history to find the latest tag for the specific SDK.
   - It intelligently calculates the next semantic version based on your chosen bump-type.
   - It mutates the necessary configuration files *in memory* on the detached commit.
4. **End-to-End Automation**: The workflow tags the detached commit, automatically generates a comprehensive changelog, drafts a GitHub Release, and finally pushes the built artifacts to the respective public registries (npm, PyPI, Maven Central).

Because this entire process happens on a detached HEAD, the version bump commits are **never merged back** into `main`. The release tag exists strictly as an immutable snapshot for distribution.

### 🔍 Exploring the Workflows

If you want to see how this powerful automation is wired together, check out the following GitHub Actions workflows:
- **[`reusable-sdk-release.yml`](../.github/workflows/reusable-sdk-release.yml)**: The core engine that calculates versions, mutates files, and tags detached commits.
- **[`reusable-release-drafter.yml`](../.github/workflows/reusable-release-drafter.yml)**: Automatically generates changelogs and drafts the GitHub Release.
- **Language-specific pipelines**: Look at [`publish-js-sdk.yml`](../.github/workflows/publish-js-sdk.yml), [`publish-python-sdk.yml`](../.github/workflows/publish-python-sdk.yml), or [`publish-java-sdk.yml`](../.github/workflows/publish-java-sdk.yml) to see how the code is verified, built, and pushed to public registries!

### ⚠️ Important Note for Contributors

**Do not manually bump versions in your Pull Requests.** 

If you are contributing code to an SDK, leave the version configuration files exactly as they are. The CI/CD pipeline handles all version calculations, file bumping, and distribution dynamically upon release.

For detailed implementation and usage instructions, please refer to the source code within each SDK's directory.
