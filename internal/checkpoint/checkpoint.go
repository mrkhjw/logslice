// Package checkpoint provides persistent read-position tracking for log files.
// It allows logslice to resume parsing from the last known offset after restart.
package checkpoint

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// State holds the persisted position for a single log file.
type State struct {
	Path    string    `json:"path"`
	Offset  int64     `json:"offset"`
	Updated time.Time `json:"updated"`
}

// Store manages checkpoint state for one or more log files.
type Store struct {
	mu       sync.Mutex
	states   map[string]*State
	filePath string
}

// New loads (or creates) a checkpoint store backed by filePath.
func New(filePath string) (*Store, error) {
	s := &Store{
		states:   make(map[string]*State),
		filePath: filePath,
	}
	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return s, nil
}

// Get returns the last saved offset for logPath, or 0 if none exists.
func (s *Store) Get(logPath string) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if st, ok := s.states[logPath]; ok {
		return st.Offset
	}
	return 0
}

// Set updates the offset for logPath and persists the store to disk.
func (s *Store) Set(logPath string, offset int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states[logPath] = &State{
		Path:    logPath,
		Offset:  offset,
		Updated: time.Now().UTC(),
	}
	return s.save()
}

// Delete removes the checkpoint for logPath and persists the change.
func (s *Store) Delete(logPath string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.states, logPath)
	return s.save()
}

func (s *Store) load() error {
	f, err := os.Open(s.filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&s.states)
}

func (s *Store) save() error {
	f, err := os.CreateTemp("", "checkpoint-*.tmp")
	if err != nil {
		return err
	}
	tmpName := f.Name()
	if err := json.NewEncoder(f).Encode(s.states); err != nil {
		f.Close()
		os.Remove(tmpName)
		return err
	}
	f.Close()
	return os.Rename(tmpName, s.filePath)
}
