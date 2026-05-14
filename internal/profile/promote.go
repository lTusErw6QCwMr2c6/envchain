package profile

import "fmt"

// PromoteOptions controls how a profile is promoted from one environment to another.
type PromoteOptions struct {
	// Overwrite allows overwriting the destination profile if it already exists.
	Overwrite bool
	// StripParents removes parent chain references from the promoted profile.
	StripParents bool
	// SuffixTag optionally appends a tag to the destination profile after promotion.
	SuffixTag string
}

// PromoteProfile copies a profile from src to dst, optionally transforming it
// according to opts. It is intended to "promote" a profile across environments
// (e.g. staging -> production) while allowing selective cleanup of metadata.
func PromoteProfile(st Store, src, dst string, opts PromoteOptions) error {
	if src == dst {
		return fmt.Errorf("promote: source and destination must differ")
	}

	srcProfile, err := st.Load(src)
	if err != nil {
		return fmt.Errorf("promote: load source %q: %w", src, err)
	}

	if !opts.Overwrite {
		_, err := st.Load(dst)
		if err == nil {
			return fmt.Errorf("promote: destination %q already exists (use --overwrite to replace)", dst)
		}
	}

	dstProfile := Profile{
		Name: dst,
		Vars: make([]Var, len(srcProfile.Vars)),
	}
	copy(dstProfile.Vars, srcProfile.Vars)

	if !opts.StripParents {
		dstProfile.Parents = make([]string, len(srcProfile.Parents))
		copy(dstProfile.Parents, srcProfile.Parents)
	}

	if err := st.Save(dstProfile); err != nil {
		return fmt.Errorf("promote: save destination %q: %w", dst, err)
	}

	if opts.SuffixTag != "" {
		if err := AddTag(st, dst, opts.SuffixTag); err != nil {
			return fmt.Errorf("promote: add tag to %q: %w", dst, err)
		}
	}

	return nil
}
