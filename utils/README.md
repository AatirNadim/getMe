# getMe Utilities

This directory contains shared utility packages used across the `getMe` project.

## `logger`

A simple logging utility to provide formatted and leveled log output. It helps to:

- Standardize log messages.
- Provide different log levels (e.g., INFO, DEBUG, ERROR, SUCCESS).
- Make it easy to enable or disable logging globally.

### Usage

```go
import "getMeMod/utils/logger"

logger.Info("This is an informational message.")
logger.Error("This is an error message:", err)
```

## `constants`

Defines global constants used in multiple parts of the application, particularly in the storage engine. This package helps to:

- Centralize magic numbers and configuration values.
- Make it easy to tune the system's performance by changing values in one place.

### Key Constants

- **`DefaultMaxSegmentSize`**: The maximum size a segment file can reach before a new one is created.
- **`ThresholdForCompaction`**: The number of inactive segments that triggers the compaction process.
- **`MaxChunkSize`**: The maximum size of a batch write buffer before it is flushed to disk.
