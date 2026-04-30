package secret

import (
	"os"
	"testing"
)

func TestDefaultProvider_Env(t *testing.T) {
	t.Setenv("ENVCHAIN_PROVIDER", "env")
	p, err := DefaultProvider()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := p.(*envProvider); !ok {
		t.Errorf("expected *envProvider, got %T", p)
	}
}

func TestDefaultProvider_Vault_MissingEnv(t *testing.T) {
	t.Setenv("ENVCHAIN_PROVIDER", "vault")
	os.Unsetenv("VAULT_ADDR")
	os.Unsetenv("VAULT_TOKEN")
	_, err := DefaultProvider()
	if err == nil {
		t.Fatal("expected error for missing vault env vars")
	}
}

func TestDefaultProvider_Vault_WithEnv(t *testing.T) {
	t.Setenv("ENVCHAIN_PROVIDER", "vault")
	t.Setenv("VAULT_ADDR", "http://localhost:8200")
	t.Setenv("VAULT_TOKEN", "root")
	p, err := DefaultProvider()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil provider")
	}
}

func TestDefaultProvider_Doppler_MissingEnv(t *testing.T) {
	t.Setenv("ENVCHAIN_PROVIDER", "doppler")
	os.Unsetenv("DOPPLER_TOKEN")
	os.Unsetenv("DOPPLER_PROJECT")
	os.Unsetenv("DOPPLER_CONFIG")
	_, err := DefaultProvider()
	if err == nil {
		t.Fatal("expected error for missing doppler env vars")
	}
}

func TestDefaultProvider_Doppler_WithEnv(t *testing.T) {
	t.Setenv("ENVCHAIN_PROVIDER", "doppler")
	t.Setenv("DOPPLER_TOKEN", "dp.st.abc")
	t.Setenv("DOPPLER_PROJECT", "myapp")
	t.Setenv("DOPPLER_CONFIG", "dev")
	p, err := DefaultProvider()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := p.(*dopplerProvider); !ok {
		t.Errorf("expected *dopplerProvider, got %T", p)
	}
}

func TestDefaultProvider_Unknown(t *testing.T) {
	t.Setenv("ENVCHAIN_PROVIDER", "notreal")
	_, err := DefaultProvider()
	if err == nil {
		t.Fatal("expected error for unknown provider")
	}
}

func TestDefaultProvider_FallbackKeyring(t *testing.T) {
	os.Unsetenv("ENVCHAIN_PROVIDER")
	p, err := DefaultProvider()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil provider")
	}
}
