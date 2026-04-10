package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/AatirNadim/getMe/tests/util"
)

// BenchmarkPut measures the performance of concurrent single-key writes.
func runBenchmarkPut(b *testing.B) {
	kv, cleanup := util.SetupStoreForBenchMarking(b)
	defer cleanup()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Each goroutine works with a different key to avoid overwriting
			key := util.GenerateRandomString(16)
			value := util.GenerateRandomString(128)
			// fmt.Println("Putting key:", key)
			if err := kv.Put(key, value); err != nil {
				b.Errorf("Put failed: %v", err)
			}
		}
	})
}

// BenchmarkGet measures the performance of concurrent single-key reads.
func runBenchmarkGet(b *testing.B) {
	kv, cleanup := util.SetupStoreForBenchMarking(b)
	defer cleanup()

	// Pre-populate the store with a large number of keys
	const numKeys = 10000
	keys := make([]string, numKeys)
	for i := 0; i < numKeys; i++ {
		keys[i] = fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		if err := kv.Put(keys[i], value); err != nil {
			b.Fatalf("Failed to pre-populate store for Get benchmark: %v", err)
		}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		// Each goroutine reads random keys from the pre-populated set
		for pb.Next() {
			key := keys[rand.Intn(numKeys)]
			if _, _, err := kv.Get(key); err != nil {
				b.Errorf("Get failed for key %s: %v", key, err)
			}
		}
	})
}

// BenchmarkBatchPut measures the performance of bulk-writing batches of key-value pairs.
func runBenchmarkBatchPut(b *testing.B) {
	kv, cleanup := util.SetupStoreForBenchMarking(b)
	defer cleanup()

	const batchSize = 100
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Create a new batch for each benchmark iteration
		batch := make(map[string]string, batchSize)
		for j := 0; j < batchSize; j++ {
			key := fmt.Sprintf("batch%d_key%d", i, j)
			value := util.GenerateRandomString(128)
			batch[key] = value
		}

		if _, err := kv.BatchPut(batch); err != nil {
			b.Fatalf("BatchPut failed: %v", err)
		}
	}
}

// BenchmarkDelete measures the performance of concurrent single-key deletions.
func runBenchmarkDelete(b *testing.B) {
	kv, cleanup := util.SetupStoreForBenchMarking(b)
	defer cleanup()

	const numKeys = 100000
	var mu sync.Mutex
	keys := make(map[string]struct{}, numKeys)

	// Pre-populate the store
	for i := 0; i < numKeys; i++ {
		key := "key" + strconv.Itoa(i)
		keys[key] = struct{}{}
		if err := kv.Put(key, "value"); err != nil {
			b.Fatalf("Failed to pre-populate store for Delete benchmark: %v", err)
		}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Safely grab a key to delete
			mu.Lock()
			var keyToDelete string
			for k := range keys {
				keyToDelete = k
				break
			}
			delete(keys, keyToDelete)
			mu.Unlock()

			if keyToDelete == "" {
				continue // All keys have been deleted
			}

			if err := kv.Delete(keyToDelete); err != nil {
				b.Errorf("Delete failed: %v", err)
			}
		}
	})
}

// BenchmarkReadWriteMixed measures performance under a mixed workload of reads and writes.
// using closure here to avoid code duplication for different read/write ratios
func benchmarkReadWriteMixedWithRatio(readThreshold int) func(b *testing.B) {
	return func(b *testing.B) {
		kv, cleanup := util.SetupStoreForBenchMarking(b)
		defer cleanup()

		// Pre-populate with some initial data
		const numInitialKeys = 10000
		keys := make([]string, numInitialKeys)
		for i := 0; i < numInitialKeys; i++ {
			keys[i] = fmt.Sprintf("key%d", i)
			if err := kv.Put(keys[i], "initial_value"); err != nil {
				b.Fatalf("Failed to pre-populate store: %v", err)
			}
		}

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// Use the enclosed readThreshold parameter
				if rand.Intn(10) < readThreshold {
					// Perform a read
					key := keys[rand.Intn(numInitialKeys)]
					_, _, _ = kv.Get(key)
				} else {
					// Perform a write
					key := util.GenerateRandomString(16)
					value := util.GenerateRandomString(128)
					_ = kv.Put(key, value)
				}
			}
		})
	}
}

func runBatchGetBenchmark(b *testing.B) {
	kv, cleanup := util.SetupStoreForBenchMarking(b)
	defer cleanup()

	const numKeys = 100000
	const batchSize = 100

	keys := make([]string, numKeys)
	initialBatch := make(map[string]string, numKeys)

	// Pre-populate an initial batch map
	for i := 0; i < numKeys; i++ {
		key := fmt.Sprintf("key%d", i)
		keys[i] = key
		initialBatch[key] = fmt.Sprintf("value%d", i)
	}

	// Bulk write the values using BatchPut
	if _, err := kv.BatchPut(initialBatch); err != nil {
		b.Fatalf("Failed to pre-populate store with BatchPut: %v", err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			batchKeys := make([]string, batchSize)
			for i := 0; i < batchSize; i++ {
				batchKeys[i] = keys[rand.Intn(numKeys)]
			}

			if _, err := kv.BatchGet(batchKeys); err != nil {
				b.Errorf("BatchGet failed: %v", err)
			}
		}
	})
}

// runMultipleTimes wraps a benchmark function and runs it n times as sub-benchmarks.
func runMultipleTimes(b *testing.B, n int, benchFunc func(*testing.B)) {
	for i := 1; i <= n; i++ {
		b.Run(fmt.Sprintf("Iteration%d", i), benchFunc)
	}
}

func BenchmarkGet(b *testing.B) {
	runMultipleTimes(b, 3, runBenchmarkGet)
}

// BenchmarkBatchGet measures the performance of bulk-reading key-value pairs using BatchGet.
func BenchmarkBatchGet(b *testing.B) {
	runMultipleTimes(b, 3, runBatchGetBenchmark)
}

func BenchmarkPut(b *testing.B) {
	runMultipleTimes(b, 3, runBenchmarkPut)
}

func BenchmarkBatchPut(b *testing.B) {
	runMultipleTimes(b, 3, runBenchmarkBatchPut)
}

// 90% reads, 10% writes
func BenchmarkReadWriteMixed_90_10(b *testing.B) {
	runMultipleTimes(b, 3, benchmarkReadWriteMixedWithRatio(9))
}

// 80% reads, 20% writes
func BenchmarkReadWriteMixed_80_20(b *testing.B) {
	runMultipleTimes(b, 3, benchmarkReadWriteMixedWithRatio(8))
}

func BenchmarkDelete(b *testing.B) {
	runMultipleTimes(b, 3, runBenchmarkDelete)
}
