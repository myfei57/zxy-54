package belt

import (
	"fmt"

	"minebelt/internal/ns"
)

func (l *Line) AddZone(beltID, zoneName string, sectionIDs []string) error {
	b, ok := l.belts[beltID]
	if !ok {
		return fmt.Errorf("belt %s not found", beltID)
	}
	z := l.reg.AddZone(zoneName, sectionIDs...)
	l.zones[beltID] = z.ID
	l.zoneSec[z.ID] = append([]string(nil), sectionIDs...)
	b.Zone = zoneName
	return nil
}

func (l *Line) Extend(beltID, sectionName string) error {
	b, ok := l.belts[beltID]
	if !ok {
		return fmt.Errorf("belt %s not found", beltID)
	}
	sec := l.reg.AddSection(sectionName, ns.KindExtended)
	if err := l.reg.AssignBelt(sec.ID, beltID); err != nil {
		return err
	}
	zoneID := l.zones[beltID]
	if zoneID == "" {
		b.Section = sectionName
		return nil
	}
	sections := append(append([]string(nil), l.zoneSec[zoneID]...), sec.ID)
	if err := l.reg.RemapZoneSections(zoneID, sections); err != nil {
		return err
	}
	l.zoneSec[zoneID] = sections
	b.Section = sectionName
	return nil
}

func (l *Line) PatrolRoute(beltID string) (string, error) {
	b, ok := l.belts[beltID]
	if !ok {
		return "", fmt.Errorf("belt %s not found", beltID)
	}
	zoneID := l.zones[beltID]
	if zoneID == "" {
		return b.Section, nil
	}
	rows := l.reg.OrderedSections(zoneID)
	if len(rows) == 0 {
		return b.Section, nil
	}
	return rows[len(rows)-1].Name, nil
}
