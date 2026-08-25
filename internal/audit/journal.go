package audit

type Journal struct {
	events []Event
	limit  int
	sink   Saver
}

func NewJournal(limit int, sink Saver) *Journal {
	if limit < 1 {
		limit = 1000
	}
	return &Journal{limit: limit, sink: sink}
}

func (j *Journal) Append(kind, beltID, detail string) Event {
	ev := NewEvent(kind, beltID, detail)
	j.events = append(j.events, ev)
	if len(j.events) > j.limit {
		j.events = j.events[len(j.events)-j.limit:]
	}
	if j.sink != nil {
		_ = j.sink.SaveJSON("audit-events.json", j.events)
	}
	return ev
}

func (j *Journal) Recent(kind string, n int) []Event {
	if n < 1 {
		n = 10
	}
	out := make([]Event, 0, n)
	for i := len(j.events) - 1; i >= 0 && len(out) < n; i-- {
		if kind == "" || j.events[i].Kind == kind {
			out = append(out, j.events[i])
		}
	}
	return out
}

func (j *Journal) Load(events []Event) {
	j.events = append([]Event(nil), events...)
}
