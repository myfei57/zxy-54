package audit

func (j *Journal) Counts() map[string]int {
	counts := make(map[string]int)
	for _, ev := range j.events {
		counts[ev.Kind]++
	}
	return counts
}

func (j *Journal) Total() int {
	return len(j.events)
}
