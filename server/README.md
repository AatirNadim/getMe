# getMe Storage Server

This directory contains the core of the getMe key-value store: a high-performance, persistent storage engine written in Go. It is designed based on the principles of the Bitcask paper, prioritizing high write throughput and fast read access.

## Architecture

The server's storage engine is built on a log-structured hash table model. This means that all data is written sequentially to an append-only log file, and an in-memory hash table (`HashTable`) is used as an index to quickly locate data on disk.

### Key Components

- **`store.go`**: This is the main entry point and public API for the storage engine. It orchestrates the interactions between the hash table and the segment manager. It exposes methods like `Put`, `Get`, `Delete`, and `BatchPut`.

- **`core/`**: This directory contains the heart of the storage engine.
  - **`hashTable.go`**: A thread-safe, in-memory map that stores keys and their corresponding locations on disk (Segment ID, offset, size). It provides fast lookups for read operations.
  - **`segment.go`**: Represents a single log file on disk. All writes are appended to the active segment. Once a segment reaches its maximum size, it becomes immutable and a new active segment is created.
  - **`segmentManager.go`**: Manages the collection of all segment files. It handles writing data to the active segment, reading data from older segments, and orchestrating the compaction process.
  - **`entry.go`**: Defines the data structure for a single key-value record as it is serialized to disk. This includes the key, value, timestamp, and a CRC checksum for data integrity.
  - **`compactedSegmentManager.go`**: Manages the process of compaction. It reads data from older, "dirty" segments, writes only the latest value for each key into new, clean, segments, and then facilitates the atomic swap of the old segments for the new ones.

- **`src/`**: This directory contains the networking layer that exposes the storage engine.
  - **`server.go`**: Implements the primary server logic, handling incoming requests.
  - **`muxHandler.go`**: Manages the request routing and dispatches commands to the appropriate store methods.
  - **`socket.go`**: (If applicable) Contains logic for handling Unix socket connections for local, high-performance inter-process communication.

## Core Concepts

### Write Path

1. A `Put` or `BatchPut` request arrives at the `Store`.
2. The key-value pair(s) are serialized into the `Entry` format.
3. The `SegmentManager` appends this serialized data to the current active `Segment` file.
4. The `SegmentManager` returns the disk location (Segment ID and offset) of the new data.
5. The `Store` updates the in-memory `HashTable` with the key and its new location.

This append-only design makes writes extremely fast, as it avoids slow, random disk writes.

### Read Path

1. A `Get` request for a key arrives at the `Store`.
2. The `Store` performs a fast lookup in the `HashTable`.
3. If the key exists, the `HashTable` returns the disk location (Segment ID and offset).
4. The `SegmentManager` opens the corresponding `Segment` file and reads the value directly from the given offset.

### Compaction

Over time, as keys are updated or deleted, the segment files will contain stale data. The compaction process is a background task that cleans up this old data:

1. The `SegmentManager` identifies several old, inactive segments to be compacted.
2. It iterates through these segments, reading every key-value pair.
3. For each key, it checks the main `HashTable` to see if the data is still the latest version.
4. Only the latest, "live" data is written to new, clean, "compacted" segment files.
5. Once complete, the `SegmentManager` atomically updates the `HashTable` to point to the new locations and then safely deletes the old segment files.

This process reclaims disk space and keeps read performance high by reducing the amount of stale data the system needs to scan over during startup.
