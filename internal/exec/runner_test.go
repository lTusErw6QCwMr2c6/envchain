package exec

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envchain/envchain/internal/profile"
	"github.com/envchain/envchain/internal/secret"
)

func newTempStore(t *testing.T) *profile.Store {
	t.Helper()
	dir := t.TempDir()
	store, err := profile.NewStore(filepath.Join(dir, "profiles"))
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return store
}

func TestRunner_Run_Success(t *testing.T) {
	provider := secret.NewEnvProvider()
	store := newTempStore(t)

	p := &profile.Profile{
		Name: "test",
		Vars: []profile.Var{{Name: "GREETING", Secret: false, Value: "hello"}},
	}
	if err := store.Save(p); err != nil {
		t.Fatalf("Save: %v", err)
	}

	runner := NewRunner(store, provider)
	if err := runner.Run("test", []string{"env"}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunner_Run_ProfileNotFound(t *testing.T) {
	provider := secret.NewEnvProvider()
	store := newTempStore(t)
	runner := NewRunner(store, provider)

	err := runner.Run("nonexistent", []string{"env"})
	if err == nil {
		t.Error("expected error for missing profile, got nil")
	}
}

func TestRunner_Run_NoCommand(t *testing.T) {
	provider := secret.NewEnvProvider()
	store := newTempStore(t)
	runner := NewRunner(store, provider)

	err := runner.Run("any", []string{})
	if err == nil {
		t.Error("expected error for empty command, got nil")
	}
}

func TestMergeEnv(t *testing.T) {
	base := []string{"FOO=bar", "BAZ=qux", "HOME=/root"}
	overrides := map[string]string{"FOO": "overridden", "NEW": "value"}

	result := mergeEnv(base, overrides)

	resultMap := make(map[string]string)
	for _, entry := range result {
		k := envKey(entry)
		v := entry[len(k)+1:]
		resultMap[k] = v
	}

	if resultMap["FOO"] != "overridden" {
		t.Errorf("expected FOO=overridden, got %q", resultMap["FOO"])
	}
	if resultMap["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %q", resultMap["BAZ"])
	}
	if resultMap["NEW"] != "value" {
		t.Errorf("expected NEW=value, got %q", resultMap["NEW"])
	}
	if resultMap["HOME"] != "/root" {
		t.Errorf("expected HOME=/root, got %q", resultMap["HOME"])
	}
}

func TestEnvKey(t *testing.T) {
	tests := []struct{ input, want string }{
		{"FOO=bar", "FOO"},
		{"A=", "A"},
		{"NOEQUALS", "NOEQUALS"},
	}
	for _, tc := range tests {
		if got := envKey(tc.input); got != tc.want {
			t.Errorf("envKey(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
	_ = os.Getenv // suppress unused import
}
