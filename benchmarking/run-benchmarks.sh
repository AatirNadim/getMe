#!/bin/bash

# Runs the getMe benchmark suite end to end.
# Simply executes the Go benchmarks located under benchmarking/.
#
# This script is safe to run on developer machines as well as in CI.

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)"

echo "=== Building and running benchmarks ==="
pushd "$REPO_ROOT/benchmarking" >/dev/null
# Run only benchmarks (skip regular tests) and surface memory statistics.
go test -run=^$ -bench=. -benchmem ./...
popd >/dev/null

echo "=== Benchmark run complete ==="
