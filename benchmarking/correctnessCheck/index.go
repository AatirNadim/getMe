package correctnessCheck

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"getMeMod/server/store"
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

// BenchmarkGet_Correctness measures the performance of reads while also verifying data correctness.
func BenchmarkGet_Correctness(b *testing.B) {
	kv, cleanup := setupStore(b)
	defer cleanup()

	noFalseFlag := true

	// Pre-populate the store with a large number of keys and store them for verification
	const numKeys = 10000
	expectedData := make(map[string]string, numKeys)
	for i := 0; i < numKeys; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		expectedData[key] = value
		if err := kv.Put(key, value); err != nil {
			b.Fatalf("Failed to pre-populate store for Get benchmark: %v", err)
		}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		// Each goroutine reads random keys from the pre-populated set
		keys := make([]string, 0, len(expectedData))
		for k := range expectedData {
			keys = append(keys, k)
		}

		for pb.Next() {
			key := keys[rand.Intn(len(keys))]
			retrievedValue, exists, err := kv.Get(key)
			if err != nil {
				b.Errorf("Get failed for key %s: %v", key, err)
				continue
			}
			if !exists {
				b.Errorf("Expected key %s to exist, but it was not found", key)
				continue
			}
			if retrievedValue != expectedData[key] {
				b.Errorf("Data mismatch for key %s: got %s, want %s", key, retrievedValue, expectedData[key])
				noFalseFlag = false
			}
		}
	})

	fmt.Println("coming here after benchmark")

	if !noFalseFlag {
		fmt.Println("Data correctness check failed during Get benchmark")
	}

}

// BenchmarkReadWriteMixed_Correctness measures performance and correctness under a mixed workload.
func BenchmarkReadWriteMixed_Correctness(b *testing.B) {
	kv, cleanup := setupStore(b)
	defer cleanup()

	noFalseFlag := true

	// Pre-populate with some initial data and store for verification
	const numInitialKeys = 10000
	expectedData := new(sync.Map) // Use a concurrent map for safe access
	keys := make([]string, numInitialKeys)
	for i := 0; i < numInitialKeys; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("initial_value_%d", i)
		keys[i] = key
		expectedData.Store(key, value)
		if err := kv.Put(key, value); err != nil {
			b.Fatalf("Failed to pre-populate store: %v", err)
		}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 80% reads, 20% writes
			if rand.Intn(10) < 8 {
				// Perform a read and verify
				key := keys[rand.Intn(numInitialKeys)]
				retrievedValue, exists, err := kv.Get(key)
				if err != nil {
					b.Errorf("Get failed during mixed test for key %s: %v", key, err)
					continue
				}
				if !exists {
					b.Errorf("Expected key %s to exist, but it was not found", key)
					continue
				}
				expected, _ := expectedData.Load(key)
				if retrievedValue != expected.(string) {
					b.Errorf("Data mismatch for key %s: got %s, want %s", key, retrievedValue, expected.(string))
					noFalseFlag = false
				}
			} else {
				// Perform a write
				key := generateRandomString(16)
				value := generateRandomString(128)
				expectedData.Store(key, value) // Store for potential future reads, though unlikely in this test structure
				if err := kv.Put(key, value); err != nil {
					b.Errorf("Put failed during mixed test: %v", err)
				}
			}
		}
	})

	fmt.Println("coming here after benchmark")

	if !noFalseFlag {
		fmt.Println("Data correctness check failed during Read/Write mixed benchmark")
	}

}
