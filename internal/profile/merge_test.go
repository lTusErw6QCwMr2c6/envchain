package profile

import (
	"testing"
)

func TestMergeProfiles_NoOverlap(t *testing.T) {
	dst := &Profile{Name: "base", Vars: []EnvVar{{Name: "A", Value: "1"}}}
	src := &Profile{Name: "extra", Vars: []EnvVar{{Name: "B", Value: "2"}}}

	merged, result := MergeProfiles(dst, src, false)

	if len(merged.Vars) != 2 {
		t.Fatalf("expected 2 vars, got %d", len(merged.Vars))
	}
	if len(result.Added) != 1 || result.Added[0] != "B" {
		t.Errorf("expected Added=[B], got %v", result.Added)
	}
	if len(result.Overwritten) != 0 {
		t.Errorf("expected no overwritten vars, got %v", result.Overwritten)
	}
}

func TestMergeProfiles_OverwriteTrue(t *testing.T) {
	dst := &Profile{Name: "base", Vars: []EnvVar{{Name: "A", Value: "old"}}}
	src := &Profile{Name: "override", Vars: []EnvVar{{Name: "A", Value: "new"}}}

	merged, result := MergeProfiles(dst, src, true)

	if len(merged.Vars) != 1 {
		t.Fatalf("expected 1 var, got %d", len(merged.Vars))
	}
	if merged.Vars[0].Value != "new" {
		t.Errorf("expected value 'new', got %q", merged.Vars[0].Value)
	}
	if len(result.Overwritten) != 1 || result.Overwritten[0] != "A" {
		t.Errorf("expected Overwritten=[A], got %v", result.Overwritten)
	}
}

func TestMergeProfiles_OverwriteFalse_PreservesDst(t *testing.T) {
	dst := &Profile{Name: "base", Vars: []EnvVar{{Name: "A", Value: "original"}}}
	src := &Profile{Name: "override", Vars: []EnvVar{{Name: "A", Value: "ignored"}}}

	merged, result := MergeProfiles(dst, src, false)

	if merged.Vars[0].Value != "original" {
		t.Errorf("expected dst value to be preserved, got %q", merged.Vars[0].Value)
	}
	if len(result.Overwritten) != 0 {
		t.Errorf("expected no overwritten vars, got %v", result.Overwritten)
	}
	if len(result.Added) != 0 {
		t.Errorf("expected no added vars, got %v", result.Added)
	}
}

func TestMergeProfiles_DoesNotMutateDst(t *testing.T) {
	dst := &Profile{Name: "base", Vars: []EnvVar{{Name: "A", Value: "1"}}}
	src := &Profile{Name: "extra", Vars: []EnvVar{{Name: "B", Value: "2"}}}

	_, _ = MergeProfiles(dst, src, true)

	if len(dst.Vars) != 1 {
		t.Errorf("dst was mutated: expected 1 var, got %d", len(dst.Vars))
	}
}

func TestMergeProfiles_PreservesName(t *testing.T) {
	dst := &Profile{Name: "myprofile", Vars: []EnvVar{}}
	src := &Profile{Name: "other", Vars: []EnvVar{}}

	merged, _ := MergeProfiles(dst, src, false)

	if merged.Name != "myprofile" {
		t.Errorf("expected merged name 'myprofile', got %q", merged.Name)
	}
}
