/*
VexCron - Drop-in replacement for the Cron daemon.
Copyright 2015 Mohit Cheppudira <mohit@muthanna.com>
*/

package lib

import (
	"reflect"
	"testing"
)

var (
	allMinutes = *RangeBitSet(0, 59)
	allHours   = *RangeBitSet(0, 23)
	allDoMs    = *RangeBitSet(1, 31)
	allMonths  = *RangeBitSet(1, 12)
	allDoWs    = *RangeBitSet(0, 6)
)

func TestParseConfig(t *testing.T) {
	cases := []struct {
		// Input
		name string
		data string

		// Output
		cfg      *Config
		hasError bool
	}{
		{
			"Empty config",
			"",
			&Config{
				"",
				Env{},
				[]Entry{},
				Stats{},
			},
			false,
		},
		{
			"Single env var.",
			"K=V",
			&Config{
				"",
				Env{
					"K": "V",
				},
				[]Entry{},
				Stats{Lines: 1, Envs: 1},
			},
			false,
		},
		{
			"Repeated env var.",
			"K=V\nK=V\nSHELL=/bin/bash",
			&Config{
				"",
				Env{
					"K":     "V",
					"SHELL": "/bin/bash",
				},
				[]Entry{},
				Stats{Lines: 3, Envs: 3},
			},
			false,
		},
		{
			"One entry",
			"* * * * * /bin/foo",
			&Config{
				"",
				Env{},
				[]Entry{
					Entry{
						[]Schedule{
							{slots: allMinutes},
							{slots: allHours},
							{slots: allDoMs},
							{slots: allMonths},
							{slots: allDoWs},
						},
						"/bin/foo",
					},
				},
				Stats{Lines: 1, Entries: 1},
			},
			false,
		},
		{
			"Mixed entries",
			"* * * * * /bin/foo\n" +
				"5 * * * * /bin/bar boo",
			&Config{
				"",
				Env{},
				[]Entry{
					Entry{
						[]Schedule{
							{slots: allMinutes},
							{slots: allHours},
							{slots: allDoMs},
							{slots: allMonths},
							{slots: allDoWs},
						},
						"/bin/foo",
					},
					Entry{
						[]Schedule{
							{slots: *NewBitSet(5)},
							{slots: allHours},
							{slots: allDoMs},
							{slots: allMonths},
							{slots: allDoWs},
						},
						"/bin/bar boo",
					},
				},
				Stats{Lines: 2, Entries: 2},
			},
			false,
		},
	}

	for i, c := range cases {
		cfg := NewConfig()
		err := cfg.Parse(c.data)

		if (err != nil) != c.hasError || !reflect.DeepEqual(c.cfg, cfg) {
			t.Errorf("Test %v: %v\nWant:\n  err: %v, cfg: %v\nGot:\n  err: %v, cfg: %v", i, c.name, c.hasError, c.cfg, err, cfg)
		}
	}
}
