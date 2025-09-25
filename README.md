# getMe - A Simple File-Based Key-Value Store

`getMe` is a command-line interface (CLI) application that provides a persistent, file-based key-value store. It's built in Go and uses an append-only log design for durability and fast writes.

## Features

- **Simple CLI**: Easy-to-use commands for `get`, `put`, and `delete` operations.
- **Persistent Storage**: Data is saved to disk in a `.getMeStore` directory in your home folder, so it persists between application runs.
- **Append-Only Log**: All writes are appended to segment files, which is an efficient pattern for write-heavy workloads.
- **Modular Design**: The core storage logic is separated from the CLI, making the code clean and maintainable.

## Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) (version 1.23 or later) installed on your system.

### Building from Source

1. Clone the repository or download the source code.
2. Navigate to the project's root directory.
3. Build the application:

    ```sh
    go build .
    ```

    This will create an executable file named `getMe` (or `getMe.exe` on Windows).

## Usage

The CLI provides three main commands to interact with the store.

### `put`

Stores a new key-value pair. Both the key and value should be strings.

```sh
./getMe put mykey "hello world"
# Output: Successfully set value for key 'mykey'
```

### `get`

Retrieves the value associated with a given key.

```sh
./getMe get mykey
# Output: hello world
```

If the key is not found, it will return an error.

### `delete`

Removes a key-value pair from the store. This is achieved by writing a special "tombstone" entry to the log.

```sh
./getMe delete mykey
# Output: Successfully deleted key 'mykey'
```

## Project Structure

The project is organized into two main parts: the CLI and the core storage engine.

```text
.
├── go.mod
├── go.sum
├── index.go          # Main entry point for the CLI application (using Cobra)
└── store/
    ├── entry.go          # Defines the data entry structure and serialization.
    ├── hashTable.go      # In-memory hash table for quick key lookups.
    ├── segment.go        # Manages individual log segment files on disk.
    ├── segmentManager.go # Manages the collection of all segments.
    ├── store.go          # The main Store struct, orchestrating all storage operations.
    ├── logger/
    │   └── logger.go     # Custom logger.
    └── utils/
        └── errorClasses.go # Custom error types.
```

## Core Concepts

The storage engine is based on an **append-only log**. Instead of modifying files in place, all changes (`put` and `delete` operations) are written as new entries to the end of the most recent log file, called a "segment".

- **Segments**: To prevent log files from growing indefinitely, the log is broken into multiple segment files. Once a segment reaches a certain size, it is closed and a new one is created.
- **Hash Table Index**: To provide fast lookups, an in-memory hash table (`hashTable.go`) maps each key to the exact location (segment ID and file offset) of its most recent value on disk. This avoids the need to scan every segment file to find a value.
- **Deletes**: When you delete a key, the system appends a special entry with a "tombstone" marker. This indicates that the key is no longer valid. The old values are cleaned up during a future compaction process (not yet implemented).
