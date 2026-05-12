package profile

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const auditFileName = "audit.json"

// AuditStore persists an AuditLog alongside the profile store directory.
type AuditStore struct {
	path string
}

// NewAuditStore creates an AuditStore rooted at dir.
func NewAuditStore(dir string) *AuditStore {
	return &AuditStore{path: filepath.Join(dir, auditFileName)}
}

// Load reads the audit log from disk. Returns an empty log if not found.
func (s *AuditStore) Load() (*AuditLog, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &AuditLog{}, nil
		}
		return nil, err
	}
	var log AuditLog
	if err := json.Unmarshal(data, &log); err != nil {
		return nil, err
	}
	return &log, nil
}

// Save writes the audit log to disk, creating the directory if needed.
func (s *AuditStore) Save(log *AuditLog) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o600)
}

// Record appends a single event and persists the updated log.
func (s *AuditStore) Record(event AuditEvent) error {
	log, err := s.Load()
	if err != nil {
		return err
	}
	log.Append(event)
	return s.Save(log)
}
