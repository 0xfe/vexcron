package lib

import (
	"testing"
)

func testBit(t *testing.T, set *BitSet, pos uint, val bool) {
	v := set.Get(pos)
	if v != val {
		t.Errorf("Pos: %v, Want: %v, Got: %v", pos, val, v)
	}
}

func TestBitSet(t *testing.T) {
	b := NewBitSet()

	testBit(t, b, 4, false)
	b.Set(4, true)
	testBit(t, b, 4, true)
	b.Clear()
	testBit(t, b, 4, false)
	b.Set(10, true)
	b.Set(20, true)
	b.Set(30, true)
	testBit(t, b, 20, true)
	testBit(t, b, 21, false)
}
