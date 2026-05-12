package profile

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAuditLog_AppendAndFilter(t *testing.T) {
	log := &AuditLog{}
	log.Append(AuditEvent{EventType: AuditEventCreated, ProfileName: "alpha"})
	log.Append(AuditEvent{EventType: AuditEventUpdated, ProfileName: "beta"})
	log.Append(AuditEvent{EventType: AuditEventUpdated, ProfileName: "alpha"})

	if len(log.Events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(log.Events))
	}

	alpha := log.FilterByProfile("alpha")
	if len(alpha) != 2 {
		t.Errorf("expected 2 alpha events, got %d", len(alpha))
	}

	updated := log.FilterByType(AuditEventUpdated)
	if len(updated) != 2 {
		t.Errorf("expected 2 updated events, got %d", len(updated))
	}
}

func TestAuditLog_Last(t *testing.T) {
	log := &AuditLog{}
	if log.Last() != nil {
		t.Fatal("expected nil for empty log")
	}
	log.Append(AuditEvent{EventType: AuditEventCreated, ProfileName: "x"})
	log.Append(AuditEvent{EventType: AuditEventDeleted, ProfileName: "x"})
	if log.Last().EventType != AuditEventDeleted {
		t.Errorf("expected last event to be deleted")
	}
}

func TestAuditLog_TimestampAutoSet(t *testing.T) {
	before := time.Now().UTC()
	log := &AuditLog{}
	log.Append(AuditEvent{EventType: AuditEventCreated, ProfileName: "ts"})
	after := time.Now().UTC()

	ts := log.Events[0].Timestamp
	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v out of expected range [%v, %v]", ts, before, after)
	}
}

func TestAuditStore_RecordAndLoad(t *testing.T) {
	dir := t.TempDir()
	s := NewAuditStore(dir)

	if err := s.Record(AuditEvent{EventType: AuditEventCreated, ProfileName: "prod"}); err != nil {
		t.Fatalf("record: %v", err)
	}
	if err := s.Record(AuditEvent{EventType: AuditEventUpdated, ProfileName: "prod", Actor: "ci"}); err != nil {
		t.Fatalf("record: %v", err)
	}

	log, err := s.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(log.Events) != 2 {
		t.Errorf("expected 2 events, got %d", len(log.Events))
	}
	if log.Events[1].Actor != "ci" {
		t.Errorf("expected actor 'ci', got %q", log.Events[1].Actor)
	}
}

func TestAuditStore_Load_NotFound(t *testing.T) {
	dir := t.TempDir()
	s := NewAuditStore(filepath.Join(dir, "nonexistent"))
	log, err := s.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(log.Events) != 0 {
		t.Errorf("expected empty log")
	}
}

func TestAuditStore_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	s := NewAuditStore(dir)
	_ = s.Record(AuditEvent{EventType: AuditEventCreated, ProfileName: "sec"})

	info, err := os.Stat(filepath.Join(dir, auditFileName))
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("expected mode 0600, got %v", info.Mode().Perm())
	}
}
