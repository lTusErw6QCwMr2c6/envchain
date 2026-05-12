package profile

import (
	"fmt"
	"sort"
	"strings"
)

// TagFilter holds criteria for filtering profiles by tags.
type TagFilter struct {
	Include []string
	Exclude []string
}

// TaggedProfiles returns all profiles from the store that have at least one
// of the given tags. If tags is empty, all profiles are returned.
func TaggedProfiles(s *Store, tags []string) ([]*Profile, error) {
	names, err := s.List()
	if err != nil {
		return nil, fmt.Errorf("tag: list profiles: %w", err)
	}

	tagSet := make(map[string]struct{}, len(tags))
	for _, t := range tags {
		tagSet[strings.ToLower(strings.TrimSpace(t))] = struct{}{}
	}

	var result []*Profile
	for _, name := range names {
		p, err := s.Load(name)
		if err != nil {
			return nil, fmt.Errorf("tag: load profile %q: %w", name, err)
		}
		if len(tagSet) == 0 || profileHasAnyTag(p, tagSet) {
			result = append(result, p)
		}
	}
	return result, nil
}

// ProfileTags returns the sorted, deduplicated list of tags for a profile.
// Tags are stored as a special variable prefix "__tag_<name>" = "1".
func ProfileTags(p *Profile) []string {
	const prefix = "__tag_"
	seen := make(map[string]struct{})
	for _, v := range p.Vars {
		if strings.HasPrefix(v.Name, prefix) {
			tag := strings.TrimPrefix(v.Name, prefix)
			seen[tag] = struct{}{}
		}
	}
	tags := make([]string, 0, len(seen))
	for t := range seen {
		tags = append(tags, t)
	}
	sort.Strings(tags)
	return tags
}

// AddTag adds a tag to the profile by inserting the sentinel variable.
// It is a no-op if the tag already exists.
func AddTag(p *Profile, tag string) {
	tag = strings.ToLower(strings.TrimSpace(tag))
	key := "__tag_" + tag
	for _, v := range p.Vars {
		if v.Name == key {
			return
		}
	}
	p.Vars = append(p.Vars, Var{Name: key, Value: "1"})
}

// RemoveTag removes a tag from the profile. It is a no-op if not present.
func RemoveTag(p *Profile, tag string) {
	tag = strings.ToLower(strings.TrimSpace(tag))
	key := "__tag_" + tag
	filtered := p.Vars[:0]
	for _, v := range p.Vars {
		if v.Name != key {
			filtered = append(filtered, v)
		}
	}
	p.Vars = filtered
}

func profileHasAnyTag(p *Profile, tagSet map[string]struct{}) bool {
	for _, t := range ProfileTags(p) {
		if _, ok := tagSet[t]; ok {
			return true
		}
	}
	return false
}
