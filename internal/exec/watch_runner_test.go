package exec_test

import (
	"context"
	"os"
	"testing"
	"time"

	execpkg "github.com/envchain/envchain/internal/exec"
	"github.com/envchain/envchain/internal/profile"
	"github.com/envchain/envchain/internal/secret"
)

func newWatchRunnerStore(t *testing.T) profile.Store {
	t.Helper()
	dir, err := os.MkdirTemp("", "envchain-wr-*")
	if err != nil {
		t.Fatalf("mkdirtemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	store, err := profile.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return store
}

func saveWRProfile(t *testing.T, store profile.Store, name string, vars map[string]string) {
	t.Helper()
	p := &profile.Profile{Name: name}
	for k, v := range vars {
		p.Vars = append(p.Vars, profile.Var{Name: k, Value: v})
	}
	if err := store.Save(p); err != nil {
		t.Fatalf("Save: %v", err)
	}
}

func TestWatchRunner_Run_NoCommand(t *testing.T) {
	store := newWatchRunnerStore(t)
	saveWRProfile(t, store, "dev", map[string]string{"X": "1"})
	provider := secret.NewEnvProvider()
	wr := execpkg.NewWatchRunner(store, provider, 50*time.Millisecond)
	ctx := context.Background()
	err := wr.Run(ctx, "dev", nil)
	if err == nil {
		t.Fatal("expected error for empty command")
	}
}

func TestWatchRunner_Run_ProfileNotFound(t *testing.T) {
	store := newWatchRunnerStore(t)
	provider := secret.NewEnvProvider()
	wr := execpkg.NewWatchRunner(store, provider, 50*time.Millisecond)
	ctx := context.Background()
	err := wr.Run(ctx, "missing", []string{"echo", "hi"})
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestWatchRunner_Run_CancelStops(t *testing.T) {
	store := newWatchRunnerStore(t)
	saveWRProfile(t, store, "dev", map[string]string{"HELLO": "world"})
	provider := secret.NewEnvProvider()
	wr := execpkg.NewWatchRunner(store, provider, 30*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- wr.Run(ctx, "dev", []string{"sleep", "10"})
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("WatchRunner did not stop after context cancellation")
	}
}
