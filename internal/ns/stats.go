package ns

func (r *Registry) SectionCount() int {
	return len(r.sections)
}

func (r *Registry) ZoneCount() int {
	return len(r.zones)
}

func (r *Registry) BeltCount() int {
	seen := make(map[string]bool)
	for _, s := range r.sections {
		for _, id := range s.BeltIDs {
			seen[id] = true
		}
	}
	return len(seen)
}

func (r *Registry) SectionsByKind(kind SectionKind) []*Section {
	out := make([]*Section, 0)
	for _, s := range r.sections {
		if s.Kind == kind {
			out = append(out, s)
		}
	}
	return out
}
