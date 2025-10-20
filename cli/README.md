# getMe Command-Line Interface (CLI)

This directory contains the source code for the `getMe` CLI, a command-line tool for interacting with a running `getMe` storage server.

## Purpose

The CLI provides a simple and scriptable way to perform basic key-value operations against the store from the terminal. It's useful for:

- Manual data inspection and manipulation.
- Debugging the server.
- Writing simple shell scripts for automation.

## Usage

The CLI is built and run from this directory. It supports several commands to interact with the server.

### Commands

- **`get <key>`**: Retrieves the value for a given key.

  ```bash
  go run . get mykey
  ```

- **`set <key> <value>`**: Sets a value for a given key.

  ```bash
  go run . set mykey "hello world"
  ```

- **`delete <key>`**: Deletes a key from the store.

  ```bash
  go run . delete mykey
  ```

- **`batch`**: Performs a bulk write operation. The data for the batch operation is read from a JSON file named `batch-input.json` located in the same directory.

  ```bash
  # Ensure batch-input.json exists
  go run . batch
  ```

### `batch-input.json`

This file should contain a single JSON object where keys are the database keys and values are the corresponding database values.

**Example `batch-input.json`:**

```json
{
    "key1": "value1",
    "key2": "value2",
    "key3": "another value"
}
```

## Implementation Details

- **`index.go`**: The main entry point for the CLI application. It parses the command-line arguments and flags to determine which operation to perform.
- **`core/httpClient.go`**: Contains the logic for communicating with the `getMe` server. It handles creating HTTP requests for each command (GET, POST, DELETE) and sending them to the server's API endpoints.
- **`utils/constants.go`**: Defines constants used within the CLI, such as the server's base URL.
