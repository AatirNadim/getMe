package main

import (
	"fmt"
	"getMeMod/server/store"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
)

// setupStore creates a new store in a temporary directory for isolated testing.
func setupStore(b *testing.B) (*store.Store, func()) {
	// Create a temporary base directory for the benchmark run
	baseDir, err := os.MkdirTemp("", "getme_benchmark_*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}

	mainPath := filepath.Join(baseDir, "main")
	compactedPath := filepath.Join(baseDir, "compacted")

	// Initialize the store
	kvStore := store.NewStore(mainPath, compactedPath)
	// fmt.Println("Store has been setup")
	// fmt.Println("Store main path:", mainPath)
	// fmt.Println("Store compacted path:", compactedPath)

	// Return the store and a cleanup function to be called deferred
	return kvStore, func() {
		kvStore.Close()
		os.RemoveAll(baseDir)
	}
}

// generateRandomString creates a random string of a given length.
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// BenchmarkPut measures the performance of concurrent single-key writes.
func BenchmarkPut(b *testing.B) {
	kv, cleanup := setupStore(b)
	defer cleanup()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Each goroutine works with a different key to avoid overwriting
			key := generateRandomString(16)
			value := generateRandomString(128)
			// fmt.Println("Putting key:", key)
			if err := kv.Put(key, value); err != nil {
				b.Errorf("Put failed: %v", err)
			}
		}
	})
}

// BenchmarkGet measures the performance of concurrent single-key reads.
func BenchmarkGet(b *testing.B) {
	kv, cleanup := setupStore(b)
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
func BenchmarkBatchPut(b *testing.B) {
	kv, cleanup := setupStore(b)
	defer cleanup()

	const batchSize = 100
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Create a new batch for each benchmark iteration
		batch := make(map[string]string, batchSize)
		for j := 0; j < batchSize; j++ {
			key := fmt.Sprintf("batch%d_key%d", i, j)
			value := generateRandomString(128)
			batch[key] = value
		}

		if err := kv.BatchPut(batch); err != nil {
			b.Fatalf("BatchPut failed: %v", err)
		}
	}
}

// BenchmarkDelete measures the performance of concurrent single-key deletions.
func BenchmarkDelete(b *testing.B) {
	kv, cleanup := setupStore(b)
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

// BenchmarkPut measures the performance of concurrent single-key writes.
func BenchmarkPut1(b *testing.B) {
	kv, cleanup := setupStore(b)
	defer cleanup()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Each goroutine works with a different key to avoid overwriting
			key := generateRandomString(16)
			value := generateRandomString(128)
			// fmt.Println("Putting key:", key)
			if err := kv.Put(key, value); err != nil {
				b.Errorf("Put failed: %v", err)
			}
		}
	})
}

// BenchmarkGet measures the performance of concurrent single-key reads.
func BenchmarkGet1(b *testing.B) {
	kv, cleanup := setupStore(b)
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
func BenchmarkBatchPut1(b *testing.B) {
	kv, cleanup := setupStore(b)
	defer cleanup()

	const batchSize = 100
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Create a new batch for each benchmark iteration
		batch := make(map[string]string, batchSize)
		for j := 0; j < batchSize; j++ {
			key := fmt.Sprintf("batch%d_key%d", i, j)
			value := generateRandomString(128)
			batch[key] = value
		}

		if err := kv.BatchPut(batch); err != nil {
			b.Fatalf("BatchPut failed: %v", err)
		}
	}
}

// BenchmarkDelete measures the performance of concurrent single-key deletions.
func BenchmarkDelete1(b *testing.B) {
	kv, cleanup := setupStore(b)
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

// BenchmarkPut measures the performance of concurrent single-key writes.
func BenchmarkPut2(b *testing.B) {
	kv, cleanup := setupStore(b)
	defer cleanup()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Each goroutine works with a different key to avoid overwriting
			key := generateRandomString(16)
			value := generateRandomString(128)
			// fmt.Println("Putting key:", key)
			if err := kv.Put(key, value); err != nil {
				b.Errorf("Put failed: %v", err)
			}
		}
	})
}

// BenchmarkGet measures the performance of concurrent single-key reads.
func BenchmarkGet2(b *testing.B) {
	kv, cleanup := setupStore(b)
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
func BenchmarkBatchPut2(b *testing.B) {
	kv, cleanup := setupStore(b)
	defer cleanup()

	const batchSize = 100
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Create a new batch for each benchmark iteration
		batch := make(map[string]string, batchSize)
		for j := 0; j < batchSize; j++ {
			key := fmt.Sprintf("batch%d_key%d", i, j)
			value := generateRandomString(128)
			batch[key] = value
		}

		if err := kv.BatchPut(batch); err != nil {
			b.Fatalf("BatchPut failed: %v", err)
		}
	}
}

// BenchmarkDelete measures the performance of concurrent single-key deletions.
func BenchmarkDelete2(b *testing.B) {
	kv, cleanup := setupStore(b)
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
func BenchmarkReadWriteMixed(b *testing.B) {
	kv, cleanup := setupStore(b)
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
			// 90% reads, 10% writes
			if rand.Intn(10) < 9 {
				// Perform a read
				key := keys[rand.Intn(numInitialKeys)]
				_, _, _ = kv.Get(key)
			} else {
				// Perform a write
				key := generateRandomString(16)
				value := generateRandomString(128)
				_ = kv.Put(key, value)
			}
		}
	})
}

func BenchmarkReadWriteMixed1(b *testing.B) {
	kv, cleanup := setupStore(b)
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
			// 90% reads, 10% writes
			if rand.Intn(10) < 9 {
				// Perform a read
				key := keys[rand.Intn(numInitialKeys)]
				_, _, _ = kv.Get(key)
			} else {
				// Perform a write
				key := generateRandomString(16)
				value := generateRandomString(128)
				_ = kv.Put(key, value)
			}
		}
	})
}

func BenchmarkReadWriteMixed2(b *testing.B) {
	kv, cleanup := setupStore(b)
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
			// 80% reads, 20% writes
			if rand.Intn(10) < 8 {
				// Perform a read
				key := keys[rand.Intn(numInitialKeys)]
				_, _, _ = kv.Get(key)
			} else {
				// Perform a write
				key := generateRandomString(16)
				value := generateRandomString(128)
				_ = kv.Put(key, value)
			}
		}
	})
}
