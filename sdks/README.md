# getMe SDKs

This directory provides Software Development Kits (SDKs) for interacting with the `getMe` storage server from different programming languages. These SDKs abstract away the underlying HTTP API calls, offering a simple, idiomatic interface for developers.

## Purpose

The goal of the SDKs is to make it easy to integrate `getMe` into various applications by providing a native, language-specific client library.

## Available SDKs

This project contains SDKs for the following languages:

- **[Go](./goSdk/)**: A native Go client for `getMe`.
- **[Java](./javaSdk/)**: A Java client, built with Gradle.
- **[JavaScript/TypeScript](./jsSdk/)**: A client library for Node.js and browser environments.
- **[Python](./pythonSdk/)**: A Python client.

Each SDK directory contains its own source code, build files, and dependencies, making them independent and easy to package and distribute.

## General Design

While implementations vary by language, all SDKs follow a similar design pattern:

1. **Client Class/Struct**: Each SDK exposes a primary `Client` object that manages the connection to the `getMe` server.
2. **Core Methods**: The `Client` provides methods that map directly to the server's core operations:
    - `get(key)`
    - `set(key, value)`
    - `delete(key)`
    - `batchSet(map_of_key_values)`
3. **HTTP Communication**: Internally, the SDKs use a standard HTTP client to communicate with the `getMe` server's REST API. They handle request creation, serialization of data (e.g., to JSON), and deserialization of responses.
4. **Error Handling**: Errors from the server (e.g., "key not found") or network issues are translated into idiomatic error types or exceptions for the respective language.

For detailed implementation and usage instructions, please refer to the source code within each SDK's directory.
