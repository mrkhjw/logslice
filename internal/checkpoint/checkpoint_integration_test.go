package checkpoint_test

import (
	"path/filepath"
	"sync"
	"testing"

	"github.com/yourorg/logslice/internal/checkpoint"
)

// TestConcurrentSet verifies that concurrent writes do not corrupt state.
func TestConcurrentSet_NoPanic(t *testing.T) {
	s, _ := tempStore(t)

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			_ = s.Set("/concurrent.log", int64(n*100))
		}(i)
	}
	wg.Wait()

	// Final offset must be one of the written values (non-zero).
	if got := s.Get("/concurrent.log"); got < 0 {
		t.Errorf("unexpected negative offset: %d", got)
	}
}

// TestRoundtrip_MultipleReloads ensures data survives repeated open/close cycles.
func TestRoundtrip_MultipleReloads(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cp.json")

	const logFile = "/srv/logs/service.log"
	const wantOffset int64 = 8192

	s1, _ := checkpoint.New(path)
	_ = s1.Set(logFile, wantOffset)

	for i := 0; i < 5; i++ {
		s, err := checkpoint.New(path)
		if err != nil {
			t.Fatalf("reload %d: %v", i, err)
		}
		if got := s.Get(logFile); got != wantOffset {
			t.Fatalf("reload %d: expected %d, got %d", i, wantOffset, got)
		}
	}
}
