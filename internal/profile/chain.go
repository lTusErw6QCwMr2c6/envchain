package profile

import "fmt"

// ResolveChain resolves a chain of profiles from the store, returning
// a merged environment map in order (later profiles override earlier ones).
func ResolveChain(store *Store, names []string) (map[string]string, error) {
	if len(names) == 0 {
		return map[string]string{}, nil
	}

	merged := map[string]string{}

	for _, name := range names {
		p, err := store.Load(name)
		if err != nil {
			return nil, fmt.Errorf("resolving chain: profile %q: %w", name, err)
		}

		env, err := p.ToEnvMap()
		if err != nil {
			return nil, fmt.Errorf("resolving chain: profile %q: %w", name, err)
		}

		for k, v := range env {
			merged[k] = v
		}
	}

	return merged, nil
}

// ChainNames returns the full ordered list of profile names in a chain,
// following each profile's Chain field recursively (depth-first, no cycles).
func ChainNames(store *Store, name string) ([]string, error) {
	return chainNamesVisit(store, name, map[string]bool{})
}

func chainNamesVisit(store *Store, name string, visited map[string]bool) ([]string, error) {
	if visited[name] {
		return nil, fmt.Errorf("cycle detected in profile chain at %q", name)
	}
	visited[name] = true

	p, err := store.Load(name)
	if err != nil {
		return nil, fmt.Errorf("chain: profile %q: %w", name, err)
	}

	var result []string
	for _, parent := range p.Chain {
		parents, err := chainNamesVisit(store, parent, visited)
		if err != nil {
			return nil, err
		}
		result = append(result, parents...)
	}
	result = append(result, name)
	return result, nil
}
