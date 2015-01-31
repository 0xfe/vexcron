/*
VexCron - Drop-in replacement for the Cron daemon.
Copyright 2015 Mohit Cheppudira <mohit@muthanna.com>
*/

package lib

import (
	"testing"
)

func TestBitSet(t *testing.T) {
	b := NewBitSet()

	test := func(pos uint, val bool) {
		v := b.Get(pos)
		if v != val {
			t.Errorf("Pos: %v, Want: %v, Got: %v", pos, val, v)
		}
	}

	test(4, false)
	b.Set(4, true)
	test(4, true)
	b.Clear()
	test(4, false)
	b.Set(10, true)
	b.Set(20, true)
	b.Set(30, true)
	test(20, true)
	test(21, false)

	b = NewBitSet(5, 6, 7, 8)
	test(6, true)
	test(9, false)

	b = RangeBitSet(10, 20)
	test(15, true)
	test(9, false)
	test(21, false)
}

func TestStringer(t *testing.T) {
	b := NewBitSet(10, 11, 16)
	want := "(10 11 16)"
	got := b.String()
	if want != got {
		t.Errorf("want: %v, got %v", want, got)
	}
}
