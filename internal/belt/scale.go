package belt

type Scale struct {
	inst float64
}

func NewScale() *Scale {
	return &Scale{}
}

func (s *Scale) Read() float64 {
	return s.inst
}

func (s *Scale) SetInstant(v float64) {
	s.inst = v
}
