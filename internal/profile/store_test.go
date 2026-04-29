package profile_test

import (
	"os"
	"testing"

	"github.com/envchain/envchain/internal/profile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTempStore(t *testing.T) *profile.Store {
	t.Helper()
	dir := t.TempDir()
	s, err := profile.NewStore(dir)
	require.NoError(t, err)
	return s
}

func TestStore_SaveAndLoad(t *testing.T) {
	s := newTempStore(t)
	p := &profile.Profile{
		Name:    "dev",
		Extends: []string{"base"},
		Vars: []profile.Var{
			{Key: "DB_HOST", Value: "localhost"},
			{Key: "DB_PASS", Secret: true, Ref: "vault:secret/dev#DB_PASS"},
		},
	}
	require.NoError(t, s.Save(p))

	loaded, err := s.Load("dev")
	require.NoError(t, err)
	assert.Equal(t, p.Name, loaded.Name)
	assert.Equal(t, p.Extends, loaded.Extends)
	assert.Equal(t, p.Vars, loaded.Vars)
}

func TestStore_Load_NotFound(t *testing.T) {
	s := newTempStore(t)
	_, err := s.Load("nonexistent")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestStore_List(t *testing.T) {
	s := newTempStore(t)
	for _, name := range []string{"alpha", "beta", "gamma"} {
		require.NoError(t, s.Save(&profile.Profile{Name: name, Vars: []profile.Var{}}))
	}
	// add a non-toml file that should be ignored
	_ = os.WriteFile(s.Dir+"/readme.txt", []byte("ignore me"), 0o600)

	names, err := s.List()
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"alpha", "beta", "gamma"}, names)
}

func TestStore_Save_InvalidProfile(t *testing.T) {
	s := newTempStore(t)
	p := &profile.Profile{Name: "bad name", Vars: []profile.Var{}}
	err := s.Save(p)
	require.Error(t, err)
}
