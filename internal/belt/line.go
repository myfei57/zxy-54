package belt

import (
	"fmt"

	"github.com/google/uuid"

	"minebelt/internal/ns"
)

type Line struct {
	reg     *ns.Registry
	belts   map[string]*Belt
	scales  map[string]*Scale
	zones   map[string]string
	zoneSec map[string][]string
}

func NewLine(reg *ns.Registry) *Line {
	return &Line{
		reg:     reg,
		belts:   make(map[string]*Belt),
		scales:  make(map[string]*Scale),
		zones:   make(map[string]string),
		zoneSec: make(map[string][]string),
	}
}

func (l *Line) AddBelt(name, kind, section string) *Belt {
	id := uuid.NewString()
	b := &Belt{ID: id, Name: name, Kind: kind, Section: section, Status: Stopped}
	l.belts[id] = b
	l.scales[id] = NewScale()
	if s, ok := l.reg.SectionByName(section); ok {
		_ = l.reg.AssignBelt(s.ID, id)
	}
	return b
}

func (l *Line) Belt(id string) (*Belt, bool) {
	b, ok := l.belts[id]
	return b, ok
}

func (l *Line) All() []*Belt {
	out := make([]*Belt, 0, len(l.belts))
	for _, b := range l.belts {
		out = append(out, b)
	}
	return out
}

func (l *Line) SetStatus(id string, s Status) error {
	b, ok := l.belts[id]
	if !ok {
		return fmt.Errorf("belt %s not found", id)
	}
	b.Status = s
	return nil
}

func (l *Line) IsRunning(id string) bool {
	b, ok := l.belts[id]
	return ok && b.Status == Running
}

func (l *Line) IsReady(id string) bool {
	return l.IsRunning(id)
}

func (l *Line) Scale(id string) (*Scale, bool) {
	s, ok := l.scales[id]
	return s, ok
}

func (l *Line) ZoneID(id string) string {
	return l.zones[id]
}
