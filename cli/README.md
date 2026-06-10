# getMe Command-Line Interface (CLI)

This directory contains the source code for the `getMe` CLI, a command-line tool for interacting with a running `getMe` storage server.

## Purpose

The CLI provides a simple and scriptable way to perform basic key-value operations against the store from the terminal. It's useful for:

- Manual data inspection and manipulation.
- Debugging the server.
- Writing simple shell scripts for automation.
- Managing the getMe server lifecycle (starting/stopping).

## Usage

The CLI is built and run from this directory. It supports several commands to interact with the server.

### Commands

#### Key-Value Operations

- **`get <key>`**: Retrieves the value for a given key.
  ```bash
  go run . get mykey
  ```

- **`getJson <key> [-o, --out <path>]`**: Retrieves a JSON value by its key. Optionally writes the JSON to a file path.
  ```bash
  go run . getJson mykey --out output.json
  ```

- **`batchGet [jsonFilePath] [-d, --data <json_string>] [-o, --out <path>]`**: Batch gets values for multiple keys specified in either a JSON file or via the `--data` flag. Optionally writes the response to a file.
  ```bash
  go run . batchGet keys.json
  go run . batchGet -d '["key1", "key2"]'
  ```

- **`put <key> <value>`**: Sets a string value for a given key.
  ```bash
  go run . put mykey "hello world"
  ```

- **`putJson <key> <jsonFilePath>`**: Sets a key with a JSON value loaded from a file.
  ```bash
  go run . putJson mykey data.json
  ```

- **`batchPut [jsonFilePath] [-d, --data <json_string>]`**: Performs a bulk write operation. The data for the batch operation is read from a JSON file or passed via the `--data` flag.
  ```bash
  go run . batchPut batch-input.json
  ```

- **`delete <key>`**: Deletes a key from the store.
  ```bash
  go run . delete mykey
  ```

- **`batchDelete [jsonFilePath] [-d, --data <json_string>]`**: Batch deletes multiple keys specified in either a JSON file or via the `--data` flag.
  ```bash
  go run . batchDelete keys-to-delete.json
  ```

- **`clearStoreConfirm`**: Clears all key-value pairs from the store.
  ```bash
  go run . clearStoreConfirm
  ```

#### Server Management & REPL

- **`start-server`**: Starts the getMe server as a background daemon.
  ```bash
  go run . start-server
  ```

- **`stop-server`**: Stops the background getMe server daemon.
  ```bash
  go run . stop-server
  ```

- **`getMe_repl`**: Starts an interactive REPL session for getMe.
  ```bash
  go run . getMe_repl
  ```

### `batch-input.json` Example

For `batchPut`, the JSON file should contain a single JSON object where keys are the database keys and values are the corresponding database values:

```json
{
    "key1": "value1",
    "key2": "value2",
    "key3": "another value"
}
```

For `batchGet` and `batchDelete`, the JSON file should contain an array of keys:

```json
[
    "key1",
    "key2",
    "key3"
]
```

## Implementation Details

- **`index.go`**: The main entry point for the CLI application using Cobra. It parses the command-line arguments and flags to determine which operation to perform.
- **`core/httpClient.go`**: Contains the logic for communicating with the `getMe` server over a Unix domain socket.
- **`core/commands/`**: Contains the individual Cobra command definitions and flag parsing for each CLI operation.
- **`utils/constants.go`**: Defines context keys used within the CLI to pass the service layer.
