package correctnessCheck

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"

	"github.com/AatirNadim/getMe/benchmarking/util"
)

// BenchmarkGet_Correctness measures the performance of reads while also verifying data correctness.
func BenchmarkGet_Correctness(b *testing.B) {
	kv,  cleanup := util.SetupStoreForBenchMarking(b)
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

	if !noFalseFlag {
		fmt.Println("Data correctness check failed during Get benchmark")
	}

}

// BenchmarkReadWriteMixed_Correctness measures performance and correctness under a mixed workload.
func BenchmarkReadWriteMixed_Correctness(b *testing.B) {
	kv,  cleanup := util.SetupStoreForBenchMarking(b)
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
				key := util.GenerateRandomString(16)
				value := util.GenerateRandomString(128)
				expectedData.Store(key, value) // Store for potential future reads, though unlikely in this test structure
				if err := kv.Put(key, value); err != nil {
					b.Errorf("Put failed during mixed test: %v", err)
				}
			}
		}
	})

	if !noFalseFlag {
		fmt.Println("Data correctness check failed during Read/Write mixed benchmark")
	}

}

// BenchmarkBatchGet_Correctness measures the performance of bulk-reading while also verifying data correctness.
func BenchmarkBatchGet_Correctness(b *testing.B) {
	kv,  cleanup := util.SetupStoreForBenchMarking(b)
	defer cleanup()

	noFalseFlag := true

	const numKeys = 10000
	const batchSize = 100
	expectedData := make(map[string]string, numKeys)
	keys := make([]string, numKeys)

	// Pre-populate an initial batch map
	for i := 0; i < numKeys; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		expectedData[key] = value
		keys[i] = key
	}

	// Perform a batch-put to populate the store
	if _, err := kv.BatchPut(expectedData); err != nil {
		b.Fatalf("Failed to pre-populate store with BatchPut: %v", err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			batchKeys := make([]string, batchSize)
			for i := 0; i < batchSize; i++ {
				batchKeys[i] = keys[rand.Intn(numKeys)]
			}

			result, err := kv.BatchGet(batchKeys)
			if err != nil {
				b.Errorf("BatchGet failed: %v", err)
				continue
			}

			if len(result.Errors) > 0 {
				b.Errorf("BatchGet returned errors: %v", result.Errors)
			}

			for _, key := range batchKeys {
				expectedValue := expectedData[key]
				if retrievedValue, ok := result.Found[key]; ok {
					if retrievedValue != expectedValue {
						b.Errorf("Data mismatch for key %s: got %s, want %s", key, retrievedValue, expectedValue)
						noFalseFlag = false
					}
				} else {
					b.Errorf("Expected key %s to exist, but it was not found in result", key)
					noFalseFlag = false
				}
			}
		}
	})

	if !noFalseFlag {
		fmt.Println("Data correctness check failed during BatchGet benchmark")
	}
}
