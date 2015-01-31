/*
VexCron - Drop-in replacement for the Cron daemon.
Copyright 2015 Mohit Cheppudira <mohit@muthanna.com>
*/

package lib

import (
	"reflect"
	"testing"
)

func TestLineType(t *testing.T) {
	cases := []struct {
		name string

		// Input
		line string

		// Output
		lineType LineType
		matches  []string
	}{
		{
			"Blank line",
			"",
			BLANK,
			nil,
		},
		{
			"Environment",
			"K=V",
			ENV,
			[]string{"K=V", "K", "V"},
		},
		{
			"Comment",
			" # foobar",
			COMMENT,
			nil,
		},
		{
			"Entry",
			"* * 1-3 6,5,4 * /bin/foo",
			ENTRY,
			[]string{"* * 1-3 6,5,4 * /bin/foo", "*", "*", "1-3", "6,5,4", "*", "/bin/foo"},
		},
		{
			"Unknown",
			"blah de blah",
			UNKNOWN,
			nil,
		},
	}

	for i, c := range cases {
		lineType, matches := getLineType(c.line)

		if c.lineType != lineType || !reflect.DeepEqual(c.matches, matches) {
			t.Errorf("Test %v: %v\nWant:\n  lineType: %v, matches: %v\nGot:\n  lineType: %v, matches: %v", i, c.name, c.lineType, c.matches, lineType, matches)
		}
	}
}

func TestSchedule(t *testing.T) {
	cases := []struct {
		name string
		// Input
		match string
		low   uint
		high  uint
		opts  scheduleOptions

		// Output
		hasError bool
		sched    Schedule
	}{
		{
			"All the things",
			"*", 0, 5, scheduleOptions{},
			false, Schedule{fields: *RangeBitSet(0, 5)},
		},
		{
			"Just one",
			"1", 0, 5, scheduleOptions{},
			false, Schedule{fields: *NewBitSet(1)},
		},
		{
			"Out of range",
			"6", 0, 5, scheduleOptions{},
			true, Schedule{},
		},
		{
			"Multiple",
			"1,4,3", 0, 5, scheduleOptions{},
			false, Schedule{fields: *NewBitSet(1, 4, 3)},
		},
		{
			"Single range",
			"3-10", 0, 20, scheduleOptions{},
			false, Schedule{fields: *RangeBitSet(3, 10)},
		},
		{
			"Multiple ranges",
			"3-5,15-16,19", 0, 20, scheduleOptions{},
			false, Schedule{fields: *NewBitSet(3, 4, 5, 15, 16, 19)},
		},
	}

	for i, c := range cases {
		sched, err := extractSchedule(c.match, c.low, c.high, c.opts)

		if (c.hasError != (err != nil)) || !reflect.DeepEqual(c.sched, sched) {
			t.Errorf("Test %v: %v\nWant:\n  err: %v, schedule: %v\nGot:\n  err: %v, schedule: %v",
				i, c.name, c.hasError, c.sched, err, sched)
		}
	}

}
