package roller

type Filter struct {
	window int
	values map[string][]float64
}

type Verdict struct {
	Filtered float64
	Peak     float64
	Over     bool
}

func NewFilter(window int) *Filter {
	if window < 1 {
		window = 1
	}
	return &Filter{window: window, values: make(map[string][]float64)}
}

func (f *Filter) Push(beltID string, value float64) float64 {
	vals := f.values[beltID]
	vals = append(vals, value)
	if len(vals) > f.window {
		vals = vals[len(vals)-f.window:]
	}
	f.values[beltID] = vals
	return average(vals)
}

func (f *Filter) Evaluate(beltID string, value, limit float64) Verdict {
	vals := f.values[beltID]
	vals = append(vals, value)
	if len(vals) > f.window {
		vals = vals[len(vals)-f.window:]
	}
	f.values[beltID] = vals
	sum := 0.0
	peak := 0.0
	for _, v := range vals {
		sum += v
		if v > peak {
			peak = v
		}
	}
	filtered := sum / float64(len(vals))
	return Verdict{Filtered: filtered, Peak: peak, Over: filtered > limit}
}

func (f *Filter) Value(beltID string) (float64, bool) {
	vals, ok := f.values[beltID]
	if !ok || len(vals) == 0 {
		return 0, false
	}
	return average(vals), true
}

func (f *Filter) Reset(beltID string) {
	delete(f.values, beltID)
}

func average(values []float64) float64 {
	var sum float64
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}
