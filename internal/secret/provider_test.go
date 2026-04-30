package secret

import (
	"testing"
)

func TestProviderType_Constants(t *testing.T) {
	tests := []struct {
		pt   ProviderType
		want string
	}{
		{ProviderEnv, "env"},
		{ProviderKeyring, "keyring"},
		{ProviderVault, "vault"},
		{ProviderAWS, "aws"},
	}
	for _, tt := range tests {
		if string(tt.pt) != tt.want {
			t.Errorf("ProviderType %q: got %q, want %q", tt.pt, string(tt.pt), tt.want)
		}
	}
}

func TestErrNotFound_Error(t *testing.T) {
	err := &ErrNotFound{Profile: "myapp", Key: "DB_PASS"}
	got := err.Error()
	want := "secret not found: profile=myapp key=DB_PASS"
	if got != want {
		t.Errorf("ErrNotFound.Error() = %q, want %q", got, want)
	}
}

func TestErrNotFound_IsError(t *testing.T) {
	var err error = &ErrNotFound{Profile: "p", Key: "k"}
	if err == nil {
		t.Fatal("expected non-nil error")
	}
}
