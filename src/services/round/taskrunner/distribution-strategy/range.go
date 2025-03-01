package distributionstrategy

type Range struct {
	Start int
	End   int
}

func (r *Range) Contains(i int) bool {
	return i >= r.Start && i <= r.End
}
func (r *Range) Length() int {
	return r.End - r.Start + 1
}
func (r *Range) Add(i int) {
	if i < r.Start {
		r.Start = i
	}
	if i > r.End {
		r.End = i
	}
}
func (r *Range) Subtract(i int) {
	if i == r.Start {
		r.Start++
	}
	if i == r.End {
		r.End--
	}
}
func (r *Range) IsEmpty() bool {
	return r.Start == r.End
}
func (r *Range) IsSingle() bool {
	return r.Start == r.End
}
func (r *Range) IsEmptyOrSingle() bool {
	return r.IsEmpty() || r.IsSingle()
}
func (r *Range) Less(i int) bool {
	return r.Start < i
}
func (r *Range) Greater(i int) bool {
	return r.End > i
}
func (r *Range) ContainsRange(other *Range) bool {
	return r.Start <= other.Start && r.End >= other.End
}
func (r *Range) Overlaps(other *Range) bool {
	return r.Start <= other.End && r.End >= other.Start
}
func (r *Range) Merge(other *Range) {
	if other.Start < r.Start {
		r.Start = other.Start
	}
	if other.End > r.End {
		r.End = other.End
	}
}
