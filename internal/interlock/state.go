package interlock

func (x *Interlock) Tagged(beltID string) bool {
	return x.tagged[beltID]
}

func (x *Interlock) TaggedCount() int {
	count := 0
	for _, tagged := range x.tagged {
		if tagged {
			count++
		}
	}
	return count
}
