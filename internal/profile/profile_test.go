package profile_test

import (
	"testing"

	"github.com/envchain/envchain/internal/profile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProfile_Validate_ValidName(t *testing.T) {
	p := &profile.Profile{Name: "my-profile_01", Vars: []profile.Var{}}
	require.NoError(t, p.Validate())
}

func TestProfile_Validate_InvalidName(t *testing.T) {
	p := &profile.Profile{Name: "bad name!", Vars: []profile.Var{}}
	err := p.Validate()
	require.ErrorIs(t, err, profile.ErrInvalidName)
}

func TestProfile_Validate_DuplicateVar(t *testing.T) {
	p := &profile.Profile{
		Name: "dup",
		Vars: []profile.Var{
			{Key: "FOO", Value: "bar"},
			{Key: "FOO", Value: "baz"},
		},
	}
	err := p.Validate()
	require.ErrorIs(t, err, profile.ErrDuplicateVar)
}

func TestProfile_ToEnvMap(t *testing.T) {
	p := &profile.Profile{
		Name: "test",
		Vars: []profile.Var{
			{Key: "A", Value: "1"},
			{Key: "B", Value: "2"},
		},
	}
	env := p.ToEnvMap()
	assert.Equal(t, map[string]string{"A": "1", "B": "2"}, env)
}
