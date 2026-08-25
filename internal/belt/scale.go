package belt

type Scale struct {
	inst float64
	seq  int
	read int
}

func NewScale() *Scale {
	return &Scale{}
}

func (s *Scale) Read() float64 {
	s.read = s.seq
	return s.inst
}

func (s *Scale) SetInstant(v float64) {
	s.inst = v
	s.seq++
}

func (s *Scale) Pending() bool {
	return s.read != s.seq
}
