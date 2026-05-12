package profile

import "strings"

// SearchOptions controls how profile search is performed.
type SearchOptions struct {
	// Query is matched against profile names and variable keys/values.
	Query string
	// Tags filters results to profiles containing all specified tags.
	Tags []string
	// CaseSensitive controls whether matching is case-sensitive.
	CaseSensitive bool
}

// SearchResult holds a matched profile and the reason it matched.
type SearchResult struct {
	Profile  Profile
	MatchedOn string // "name", "var_key", or "var_value"
}

// SearchProfiles searches profiles in the store using the given options.
// It returns all profiles whose name, variable keys, or variable values
// contain the query string, optionally filtered by tags.
func SearchProfiles(s *Store, opts SearchOptions) ([]SearchResult, error) {
	names, err := s.List()
	if err != nil {
		return nil, err
	}

	var results []SearchResult
	for _, name := range names {
		p, err := s.Load(name)
		if err != nil {
			continue
		}

		if len(opts.Tags) > 0 && !profileHasAnyTag(p, opts.Tags) {
			continue
		}

		if opts.Query == "" {
			results = append(results, SearchResult{Profile: p, MatchedOn: "name"})
			continue
		}

		q := opts.Query
		pName := p.Name
		if !opts.CaseSensitive {
			q = strings.ToLower(q)
			pName = strings.ToLower(pName)
		}

		if strings.Contains(pName, q) {
			results = append(results, SearchResult{Profile: p, MatchedOn: "name"})
			continue
		}

		matched := false
		for _, v := range p.Vars {
			key := v.Key
			val := v.Value
			if !opts.CaseSensitive {
				key = strings.ToLower(key)
				val = strings.ToLower(val)
			}
			if strings.Contains(key, q) {
				results = append(results, SearchResult{Profile: p, MatchedOn: "var_key"})
				matched = true
				break
			}
			if strings.Contains(val, q) {
				results = append(results, SearchResult{Profile: p, MatchedOn: "var_value"})
				matched = true
				break
			}
		}
		_ = matched
	}
	return results, nil
}
