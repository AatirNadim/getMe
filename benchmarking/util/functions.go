package util

import (
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/AatirNadim/getMe/server/src"
	"github.com/AatirNadim/getMe/server/store"
)

// setupStore creates a new store in a temporary directory for isolated testing.
func SetupStoreForBenchMarking(b *testing.B) (*store.Store, func()) {
	// Create a temporary base directory for the benchmark run
	baseDir, err := os.MkdirTemp("", "getme_benchmark_*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}

	mainPath := filepath.Join(baseDir, "main")
	compactedPath := filepath.Join(baseDir, "compacted")
	loggingDisabled := true
	loggingToStdout := false

	// Initialize the store

	kvStore, err := src.InitializeStore(mainPath, compactedPath, &loggingDisabled, &loggingToStdout)

	if err != nil {
		b.Fatalf("Failed to initialize store: %v", err)
	}

	// Return the store and a cleanup function to be called deferred
	return kvStore, func() {
		kvStore.Close()
		os.RemoveAll(baseDir)
	}
}

func SetupStoreForCorrectnessCheck(t *testing.T) (*src.Controllers, func()) {
	// Create a temporary base directory for the benchmark run
	baseDir, err := os.MkdirTemp("", "getme_correctness_check_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	mainPath := filepath.Join(baseDir, "main")
	compactedPath := filepath.Join(baseDir, "compacted")
	loggingDisabled := true
	loggingToStdout := false

	// Initialize the store

	kvStore, err := src.InitializeStore(mainPath, compactedPath, &loggingDisabled, &loggingToStdout)

	if err != nil {
		t.Fatalf("Failed to initialize store: %v", err)
	}

	controllers := &src.Controllers{
		StoreInstance: kvStore,
	}

	// Return the store and a cleanup function to be called deferred
	return controllers, func() {
		kvStore.Close()
		os.RemoveAll(baseDir)
	}
}

// generateRandomString creates a random string of a given length.
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
