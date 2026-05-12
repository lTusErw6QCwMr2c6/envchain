package profile

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// ImportFormat represents the format of an import source.
type ImportFormat string

const (
	ImportFormatDotenv ImportFormat = "dotenv"
	ImportFormatExport ImportFormat = "export"
)

// ImportOptions controls the behaviour of ImportProfile.
type ImportOptions struct {
	Overwrite bool
}

// ImportProfile reads environment variable definitions from r in the given
// format and merges them into the named profile stored in st.
func ImportProfile(st Store, name string, r io.Reader, format ImportFormat, opts ImportOptions) error {
	if err := ValidateName(name); err != nil {
		return fmt.Errorf("import: %w", err)
	}

	vars, err := parseImport(r, format)
	if err != nil {
		return fmt.Errorf("import: %w", err)
	}

	existing, err := st.Load(name)
	if err != nil && !IsStoreNotFound(err) {
		return fmt.Errorf("import: load profile: %w", err)
	}
	if existing == nil {
		existing = &Profile{Name: name}
	}

	incoming := &Profile{Name: name, Vars: vars}
	merged := MergeProfiles(existing, incoming, opts.Overwrite)

	if err := ValidateProfile(merged); err != nil {
		return fmt.Errorf("import: %w", err)
	}
	return st.Save(merged)
}

// parseImport reads key=value pairs from r according to format.
func parseImport(r io.Reader, format ImportFormat) ([]EnvVar, error) {
	var vars []EnvVar
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if format == ImportFormatExport {
			line = strings.TrimPrefix(line, "export ")
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("malformed line: %q", line)
		}
		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), `"`)
		if err := ValidateVarName(key); err != nil {
			return nil, fmt.Errorf("invalid variable %q: %w", key, err)
		}
		vars = append(vars, EnvVar{Name: key, Value: val})
	}
	return vars, scanner.Err()
}

// IsStoreNotFound returns true when err is a store not-found error.
func IsStoreNotFound(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(interface{ NotFound() bool })
	return ok
}
