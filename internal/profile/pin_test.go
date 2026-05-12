package profile_test

import (
	"os"
	"testing"

	"github.com/nicholasgasior/envchain/internal/profile"
)

func newPinStore(t *testing.T) (*profile.Store, *profile.PinStore) {
	t.Helper()
	dir := t.TempDir()
	s, err := profile.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	ps, err := profile.NewPinStore(dir)
	if err != nil {
		t.Fatalf("NewPinStore: %v", err)
	}
	return s, ps
}

func savePinProfile(t *testing.T, s *profile.Store, name string, vars map[string]string) {
	t.Helper()
	p := &profile.Profile{Name: name}
	for k, v := range vars {
		p.Vars = append(p.Vars, profile.Var{Key: k, Value: v})
	}
	if err := s.Save(p); err != nil {
		t.Fatalf("Save: %v", err)
	}
}

func TestPinProfile_CapturesVars(t *testing.T) {
	s, ps := newPinStore(t)
	savePinProfile(t, s, "myapp", map[string]string{"DB_URL": "postgres://localhost", "PORT": "5432"})

	pin, err := profile.PinProfile(s, "myapp", "ci")
	if err != nil {
		t.Fatalf("PinProfile: %v", err)
	}
	if pin.Vars["DB_URL"] != "postgres://localhost" {
		t.Errorf("expected DB_URL captured")
	}
	if pin.PinnedBy != "ci" {
		t.Errorf("expected pinnedBy=ci, got %s", pin.PinnedBy)
	}
	_ = ps
}

func TestPinStore_SaveAndLoad(t *testing.T) {
	s, ps := newPinStore(t)
	savePinProfile(t, s, "web", map[string]string{"HOST": "localhost"})

	pin, err := profile.PinProfile(s, "web", "deploy")
	if err != nil {
		t.Fatalf("PinProfile: %v", err)
	}
	if err := ps.Save("v1", pin); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := ps.Load("web", "v1")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Vars["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost")
	}
}

func TestPinStore_Load_NotFound(t *testing.T) {
	_, ps := newPinStore(t)
	_, err := ps.Load("ghost", "v1")
	if err == nil {
		t.Fatal("expected error for missing pin")
	}
}

func TestPinStore_Delete(t *testing.T) {
	s, ps := newPinStore(t)
	savePinProfile(t, s, "svc", map[string]string{"KEY": "val"})
	pin, _ := profile.PinProfile(s, "svc", "tester")
	_ = ps.Save("stable", pin)
	if err := ps.Delete("svc", "stable"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := ps.Load("svc", "stable"); err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestDiffPin_DetectsChanges(t *testing.T) {
	s, _ := newPinStore(t)
	savePinProfile(t, s, "app", map[string]string{"A": "1", "B": "2"})
	pin, _ := profile.PinProfile(s, "app", "base")

	// Modify the profile
	savePinProfile(t, s, "app", map[string]string{"A": "99", "C": "3"})

	diff, err := profile.DiffPin(s, pin)
	if err != nil {
		t.Fatalf("DiffPin: %v", err)
	}
	if len(diff.Changed) == 0 {
		t.Error("expected changed vars")
	}
	if len(diff.Added) == 0 {
		t.Error("expected added vars")
	}
	if len(diff.Removed) == 0 {
		t.Error("expected removed vars")
	}
	_ = os.Getenv // suppress unused import
}
