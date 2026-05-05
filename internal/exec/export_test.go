package exec

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/envchain/envchain/internal/profile"
	"github.com/envchain/envchain/internal/secret"
)

func newExportTempStore(t *testing.T) *profile.Store {
	t.Helper()
	dir := t.TempDir()
	store, err := profile.NewStore(filepath.Join(dir, "profiles.json"))
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	return store
}

func TestExporter_Export_Dotenv(t *testing.T) {
	store := newExportTempStore(t)
	p := &profile.Profile{
		Name: "web",
		Vars: []profile.Var{
			{Key: "APP_ENV", Value: "production"},
			{Key: "PORT", Value: "3000"},
		},
	}
	if err := store.Save(p); err != nil {
		t.Fatal(err)
	}

	provider := secret.NewEnvProvider()
	exporter := NewExporter(store, provider)

	var buf bytes.Buffer
	if err := exporter.Export("web", "dotenv", &buf); err != nil {
		t.Fatalf("export: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "APP_ENV=production") {
		t.Errorf("expected APP_ENV in output, got: %q", out)
	}
	if !strings.Contains(out, "PORT=3000") {
		t.Errorf("expected PORT in output, got: %q", out)
	}
}

func TestExporter_Export_ProfileNotFound(t *testing.T) {
	store := newExportTempStore(t)
	provider := secret.NewEnvProvider()
	exporter := NewExporter(store, provider)

	var buf bytes.Buffer
	err := exporter.Export("nonexistent", "dotenv", &buf)
	if err == nil {
		t.Error("expected error for missing profile")
	}
}

func TestExporter_Export_UnknownFormat(t *testing.T) {
	store := newExportTempStore(t)
	p := &profile.Profile{
		Name: "simple",
		Vars: []profile.Var{{Key: "X", Value: "1"}},
	}
	if err := store.Save(p); err != nil {
		t.Fatal(err)
	}

	provider := secret.NewEnvProvider()
	exporter := NewExporter(store, provider)

	var buf bytes.Buffer
	err := exporter.Export("simple", "toml", &buf)
	if err == nil {
		t.Error("expected error for unknown format")
	}
}
