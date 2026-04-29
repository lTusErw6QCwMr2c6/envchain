package secret

import (
	"testing"

	"github.com/zalando/go-keyring"
)

func init() {
	// Use the mock keyring backend for tests to avoid OS keyring dependency.
	keyring.MockInit()
}

func TestKeyringProvider_SetAndGet(t *testing.T) {
	p := NewKeyringProviderWithService("envchain-test")

	err := p.Set("myproject", "DB_PASSWORD", "supersecret")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	val, err := p.Get("myproject", "DB_PASSWORD")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if val != "supersecret" {
		t.Errorf("expected %q, got %q", "supersecret", val)
	}
}

func TestKeyringProvider_Get_NotFound(t *testing.T) {
	p := NewKeyringProviderWithService("envchain-test")

	_, err := p.Get("nonexistent", "MISSING_KEY")
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestKeyringProvider_Delete(t *testing.T) {
	p := NewKeyringProviderWithService("envchain-test")

	if err := p.Set("myproject", "API_KEY", "abc123"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	if err := p.Delete("myproject", "API_KEY"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err := p.Get("myproject", "API_KEY")
	if err == nil {
		t.Fatal("expected error after delete, got nil")
	}
}

func TestKeyringProvider_Delete_NotFound(t *testing.T) {
	p := NewKeyringProviderWithService("envchain-test")

	// Deleting a non-existent key should not return an error.
	if err := p.Delete("ghost", "PHANTOM_KEY"); err != nil {
		t.Errorf("expected no error deleting missing key, got: %v", err)
	}
}

func TestKeyringProvider_Account(t *testing.T) {
	p := NewKeyringProviderWithService("envchain-test")
	account := p.account("staging", "SECRET_TOKEN")
	expected := "staging/SECRET_TOKEN"
	if account != expected {
		t.Errorf("expected account %q, got %q", expected, account)
	}
}
