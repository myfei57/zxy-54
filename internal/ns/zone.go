package ns

type Zone struct {
	ID         string
	Name       string
	SectionIDs []string
}

func (z *Zone) Sections(reg *Registry) []*Section {
	out := make([]*Section, 0, len(z.SectionIDs))
	for _, id := range z.SectionIDs {
		if s, ok := reg.sections[id]; ok {
			out = append(out, s)
		}
	}
	return out
}
