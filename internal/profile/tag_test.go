package profile_test

import (
	"testing"

	"github.com/nicholasgasior/envchain/internal/profile"
)

func newTagStore(t *testing.T) *profile.Store {
	t.Helper()
	dir := t.TempDir()
	s, err := profile.NewStore(dir)
	if err != nil {
		t.Fatalf("newTagStore: %v", err)
	}
	return s
}

func saveTagProfile(t *testing.T, s *profile.Store, name string, tags []string) *profile.Profile {
	t.Helper()
	p := &profile.Profile{Name: name}
	for _, tag := range tags {
		profile.AddTag(p, tag)
	}
	if err := s.Save(p); err != nil {
		t.Fatalf("saveTagProfile %q: %v", name, err)
	}
	return p
}

func TestAddTag_NewTag(t *testing.T) {
	p := &profile.Profile{Name: "test"}
	profile.AddTag(p, "production")
	tags := profile.ProfileTags(p)
	if len(tags) != 1 || tags[0] != "production" {
		t.Errorf("expected [production], got %v", tags)
	}
}

func TestAddTag_Idempotent(t *testing.T) {
	p := &profile.Profile{Name: "test"}
	profile.AddTag(p, "staging")
	profile.AddTag(p, "staging")
	if got := len(profile.ProfileTags(p)); got != 1 {
		t.Errorf("expected 1 tag, got %d", got)
	}
}

func TestRemoveTag(t *testing.T) {
	p := &profile.Profile{Name: "test"}
	profile.AddTag(p, "alpha")
	profile.AddTag(p, "beta")
	profile.RemoveTag(p, "alpha")
	tags := profile.ProfileTags(p)
	if len(tags) != 1 || tags[0] != "beta" {
		t.Errorf("expected [beta], got %v", tags)
	}
}

func TestRemoveTag_NotPresent(t *testing.T) {
	p := &profile.Profile{Name: "test"}
	profile.AddTag(p, "only")
	profile.RemoveTag(p, "missing") // should not panic or error
	if got := len(profile.ProfileTags(p)); got != 1 {
		t.Errorf("expected 1 tag after removing non-existent, got %d", got)
	}
}

func TestTaggedProfiles_FilterByTag(t *testing.T) {
	s := newTagStore(t)
	saveTagProfile(t, s, "prod-api", []string{"production", "api"})
	saveTagProfile(t, s, "staging-api", []string{"staging", "api"})
	saveTagProfile(t, s, "prod-db", []string{"production"})

	result, err := profile.TaggedProfiles(s, []string{"production"})
	if err != nil {
		t.Fatalf("TaggedProfiles: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 profiles, got %d", len(result))
	}
}

func TestTaggedProfiles_EmptyTagsReturnsAll(t *testing.T) {
	s := newTagStore(t)
	saveTagProfile(t, s, "alpha", []string{"x"})
	saveTagProfile(t, s, "beta", []string{"y"})

	result, err := profile.TaggedProfiles(s, nil)
	if err != nil {
		t.Fatalf("TaggedProfiles: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 profiles, got %d", len(result))
	}
}

func TestProfileTags_Sorted(t *testing.T) {
	p := &profile.Profile{Name: "test"}
	profile.AddTag(p, "zebra")
	profile.AddTag(p, "apple")
	profile.AddTag(p, "mango")
	tags := profile.ProfileTags(p)
	expected := []string{"apple", "mango", "zebra"}
	for i, tag := range expected {
		if tags[i] != tag {
			t.Errorf("index %d: expected %q, got %q", i, tag, tags[i])
		}
	}
}
