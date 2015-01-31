package lib

import (
	"fmt"
)

// Simple implementation of a bit set for storing scheduling
// entries for O(1) lookup.

type BitSet uint64

const (
	maxBitSetPos = 63
)

func NewBitSet() *BitSet {
	return new(BitSet)
}

func SingleBitSet(pos uint) *BitSet {
	b := NewBitSet()
	b.Set(pos, true)
	return b
}

func MultiBitSet(poss ...uint) *BitSet {
	b := NewBitSet()
	for _, pos := range poss {
		b.Set(pos, true)
	}
	return b
}

func RangeBitSet(start, end uint) *BitSet {
	b := NewBitSet()
	for pos := start; pos <= end; pos++ {
		b.Set(pos, true)
	}
	return b
}

func (b *BitSet) Set(pos uint, val bool) {
	if pos > maxBitSetPos {
		panic(fmt.Sprintf("pos %v out of bounds", pos))
	}
	if val {
		*b = (*b | (1 << pos))
	} else {
		*b = (*b ^ (1 << pos))
	}
}

func (b *BitSet) Get(pos uint) bool {
	if pos > maxBitSetPos {
		panic(fmt.Sprintf("pos %v out of bounds", pos))
	}
	return (*b & (1 << pos)) > 0
}

func (b *BitSet) Clear() {
	*b = 0
}

func (b BitSet) String() string {
	str := "("
	var i uint
	first := true
	for i = 0; i <= maxBitSetPos; i++ {
		if b.Get(i) {
			if !first {
				str += " "
			}
			first = false
			str = fmt.Sprintf("%v%v", str, i)
		}
	}

	return str + ")"
}
