package secret

import (
	"os"
	"testing"
)

func TestEnvProvider_SetAndGet(t *testing.T) {
	p := NewEnvProvider()

	if err := p.Set("my-service", "API_KEY", "supersecret"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	val, err := p.Get("my-service", "API_KEY")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if val != "supersecret" {
		t.Errorf("expected %q, got %q", "supersecret", val)
	}
}

func TestEnvProvider_Get_NotFound(t *testing.T) {
	p := NewEnvProvider()

	_, err := p.Get("nonexistent-service", "MISSING_KEY")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestEnvProvider_Delete(t *testing.T) {
	p := NewEnvProvider()
	key := envKey("svc", "TOKEN")

	os.Setenv(key, "value")

	if err := p.Delete("svc", "TOKEN"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err := p.Get("svc", "TOKEN")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestEnvKey_Normalization(t *testing.T) {
	cases := []struct {
		service, key, want string
	}{
		{"my-service", "api.key", "ENVCHAIN_SECRET_MY_SERVICE__API_KEY"},
		{"prod env", "DB PASS", "ENVCHAIN_SECRET_PROD_ENV__DB_PASS"},
	}
	for _, tc := range cases {
		got := envKey(tc.service, tc.key)
		if got != tc.want {
			t.Errorf("envKey(%q, %q) = %q, want %q", tc.service, tc.key, got, tc.want)
		}
	}
}
