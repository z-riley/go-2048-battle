package grid

// iter is a genetic iterator which can iterate forwards or backwards.
type iter struct {
	length  int
	reverse bool
	idx     int
}

// newIter constructs a new iterator.
func newIter(length int, reverse bool) *iter {
	if reverse {
		return &iter{length: length, reverse: true, idx: length - 1}
	}
	return &iter{length: length, reverse: false, idx: 0}
}

// hasNext returns true if another iteration is available.
func (i *iter) hasNext() bool {
	if i.reverse {
		return i.idx >= 0
	}
	return i.idx < i.length
}

// next returns the index of the next index.
func (i *iter) next() int {
	if !i.hasNext() {
		panic("no more elements")
	}

	if i.reverse {
		out := i.idx
		i.idx--
		return out
	}

	out := i.idx
	i.idx++
	return out
}
