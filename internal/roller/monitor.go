package roller

import "github.com/google/uuid"

type Roller struct {
	ID     string
	BeltID string
	Index  int
	Temp   float64
}

type Monitor struct {
	rollers map[string]*Roller
	filter  *Filter
	limit   float64
}

func NewMonitor(limit float64, window int) *Monitor {
	return &Monitor{
		rollers: make(map[string]*Roller),
		filter:  NewFilter(window),
		limit:   limit,
	}
}

func (m *Monitor) AddRoller(beltID string, index int) *Roller {
	r := &Roller{ID: uuid.NewString(), BeltID: beltID, Index: index}
	m.rollers[r.ID] = r
	return r
}

func (m *Monitor) Roller(id string) (*Roller, bool) {
	r, ok := m.rollers[id]
	return r, ok
}

func (m *Monitor) RollerByBelt(beltID string, index int) (*Roller, bool) {
	for _, r := range m.rollers {
		if r.BeltID == beltID && r.Index == index {
			return r, true
		}
	}
	return nil, false
}

func (m *Monitor) Sample(id string, temp float64) (float64, error) {
	r, ok := m.rollers[id]
	if !ok {
		return 0, errRollerNotFound
	}
	r.Temp = temp
	return m.filter.Push(r.BeltID, temp), nil
}

func (m *Monitor) Filter() *Filter {
	return m.filter
}

func (m *Monitor) FilterValue(beltID string) (float64, bool) {
	return m.filter.Value(beltID)
}

func (m *Monitor) ResetFilter(beltID string) {
	m.filter.Reset(beltID)
}

func (m *Monitor) Limit() float64 {
	return m.limit
}

func (m *Monitor) Rollers(beltID string) []*Roller {
	var out []*Roller
	for _, r := range m.rollers {
		if r.BeltID == beltID {
			out = append(out, r)
		}
	}
	return out
}

var errRollerNotFound = rollerError("roller not found")

type rollerError string

func (e rollerError) Error() string {
	return string(e)
}
