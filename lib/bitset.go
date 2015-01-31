/*
VexCron - Drop-in replacement for the Cron daemon.
Copyright 2015 Mohit Cheppudira <mohit@muthanna.com>

This file implements a simple bit set container, which allows
for O(1) lookup.
*/

package lib

import (
	"fmt"
)

// Cron fields never exceed 31, so we'll never need more
// than 64 entries. Famous last words.
type BitSet uint64

const (
	maxBitSetPos = 63
)

// Return an initialzied BitSet with all bits specified
// by poss set to true.
func NewBitSet(poss ...uint) *BitSet {
	b := new(BitSet)
	for _, pos := range poss {
		b.Set(pos, true)
	}
	return b
}

// Return a bitset with all bits between start and end set
// to true.
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
