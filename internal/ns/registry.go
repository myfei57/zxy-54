package ns

import (
	"fmt"
	"sort"

	"github.com/google/uuid"
)

type Registry struct {
	sections  map[string]*Section
	zones     map[string]*Zone
	nextOrder int
}

func NewRegistry() *Registry {
	return &Registry{
		sections:  make(map[string]*Section),
		zones:     make(map[string]*Zone),
		nextOrder: 1,
	}
}

func (r *Registry) AddSection(name string, kind SectionKind) *Section {
	id := uuid.NewString()
	s := &Section{ID: id, Name: name, Kind: kind, Order: r.nextOrder}
	r.nextOrder++
	r.sections[id] = s
	return s
}

func (r *Registry) AddZone(name string, sectionIDs ...string) *Zone {
	z := &Zone{ID: uuid.NewString(), Name: name, SectionIDs: append([]string(nil), sectionIDs...)}
	r.zones[z.ID] = z
	return z
}

func (r *Registry) Section(id string) (*Section, bool) {
	s, ok := r.sections[id]
	return s, ok
}

func (r *Registry) SectionByName(name string) (*Section, bool) {
	for _, s := range r.sections {
		if s.Name == name {
			return s, true
		}
	}
	return nil, false
}

func (r *Registry) Zone(id string) (*Zone, bool) {
	z, ok := r.zones[id]
	return z, ok
}

func (r *Registry) AssignBelt(sectionID, beltID string) error {
	s, ok := r.sections[sectionID]
	if !ok {
		return fmt.Errorf("section %s not found", sectionID)
	}
	for _, b := range s.BeltIDs {
		if b == beltID {
			return nil
		}
	}
	s.BeltIDs = append(s.BeltIDs, beltID)
	return nil
}

func (r *Registry) RemapZoneSections(zoneID string, sectionIDs []string) error {
	z, ok := r.zones[zoneID]
	if !ok {
		return fmt.Errorf("zone %s not found", zoneID)
	}
	z.SectionIDs = append([]string(nil), sectionIDs...)
	for i, id := range z.SectionIDs {
		if s, ok := r.sections[id]; ok {
			s.Order = i + 1
		}
	}
	return nil
}

func (r *Registry) OrderedSections(zoneID string) []*Section {
	z, ok := r.zones[zoneID]
	if !ok {
		return nil
	}
	rows := z.Sections(r)
	sort.SliceStable(rows, func(i, j int) bool { return rows[i].Order < rows[j].Order })
	return rows
}
