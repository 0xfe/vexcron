/*
VexCron - Drop-in replacement for the Cron daemon.
Copyright 2015 Mohit Cheppudira <mohit@muthanna.com>

This file implements parsing of crontab files. We currently support
the following syntax:

Field name   | Mandatory? | Allowed values  | Allowed special characters
----------   | ---------- | --------------  | --------------------------
Seconds      | Yes        | 0-59            | * / , -
Minutes      | Yes        | 0-59            | * / , -
Hours        | Yes        | 0-23            | * / , -
Day of month | Yes        | 1-31            | * / , - ? L W
Month        | Yes        | 1-12 or JAN-DEC | * / , -
Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ? L #

We also support all vixie-cron macros:

@yearly (or @annually)	Run once a year at midnight of January 1	0 0 1 1 *
@monthly	              Run once a month at midnight of the first day of the month	0 0 1 * *
@weekly	                Run once a week at midnight on Sunday morning	0 0 * * 0
@daily	                Run once a day at midnight	0 0 * * *
@hourly	                Run once an hour at the beginning of the hour	0 * * * *
@reboot	                Run at startup	@reboot
*/

package lib

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var (
	envRE, entryRE, rangeRE, valueRE, eolRE *regexp.Regexp
)

type Env map[string]string

type Schedule struct {
	fields  BitSet
	last    bool
	weekday bool
}

type Entry struct {
	pieces []Schedule
	cmd    string
}

type Stats struct {
	Lines      int
	BlankLines int
	Comments   int

	// Number of environment variables parsed (not found)
	Envs int

	// Number of entries parsed (not found)
	Entries int
}

func init() {
	// Build regexps
	compile := func(re string) *regexp.Regexp {
		r, err := regexp.Compile(re)
		if err != nil {
			log.Fatalf("Parser init: %v", err)
		}
		return r
	}

	envRE = compile("^(\\S+)\\s*=\\s*(\\S+)")
	entryRE = compile("^(\\S+)\\s+(\\S+)\\s+(\\S+)\\s+(\\S+)\\s+(\\S+)\\s+(.+)")
	rangeRE = compile("^(\\d+)\\s*-\\s*(\\d+)")
	valueRE = compile("^(\\d+)([LW#]?)")
	eolRE = compile("\\r?\\n")
}

type LineType int

const (
	COMMENT LineType = 0
	ENV              = 1
	ENTRY            = 2
	BLANK            = 3
	UNKNOWN          = 4
)

func getLineType(line string) (LineType, []string) {
	l := strings.TrimSpace(line)

	// Skip empty lines
	if len(l) == 0 {
		return BLANK, nil
	}

	// Skip comments
	if l[0] == '#' {
		return COMMENT, nil
	}

	// Parse environment variables
	matches := envRE.FindStringSubmatch(l)
	if matches != nil {
		return ENV, matches
	}

	matches = entryRE.FindStringSubmatch(l)
	if matches != nil {
		return ENTRY, matches
	}

	return UNKNOWN, nil
}

func extractEnv(matches []string) (k, v string) {
	return matches[1], matches[2]
}

func genBitSet(low, high uint) *BitSet {
	b := NewBitSet()
	for x := low; x <= high; x++ {
		b.Set(x, true)
	}

	return b
}

type scheduleOptions struct {
	dict   map[string]int
	allowL bool
	allowW bool
	allowH bool // allow hash
}

func extractSchedule(match string, low uint, high uint, opts scheduleOptions) (Schedule, error) {
	// Wildcard schedule
	if match == "*" {
		return Schedule{fields: *genBitSet(low, high)}, nil
	}

	pieces := strings.Split(match, ",")
	bitset := NewBitSet()

	// For each comma-separated value
	for _, piece := range pieces {
		// Check for a-b range values
		r := rangeRE.FindStringSubmatch(piece)
		if r != nil {
			start, _ := strconv.ParseUint(r[1], 10, 8)
			end, _ := strconv.ParseUint(r[2], 10, 8)

			if start > end || uint(start) < low || uint(end) > high {
				return Schedule{}, fmt.Errorf("bad range: %v-%v", start, end)
			}

			for i := start; i <= end; i++ {
				bitset.Set(uint(i), true)
			}
			continue
		}

		// Check for single value
		number, err := strconv.ParseUint(piece, 10, 8)
		if err != nil {
			return Schedule{}, fmt.Errorf("can't parse %v", match)
		}
		if uint(number) < low || uint(number) > high {
			return Schedule{}, fmt.Errorf("%v out of range (%v - %v)", number, low, high)
		}
		bitset.Set(uint(number), true)
	}

	return Schedule{fields: *bitset}, nil
}

func extractMinute(match string) (Schedule, error) {
	return extractSchedule(match, 0, 59, scheduleOptions{})
}

func extractHour(match string) (Schedule, error) {
	return extractSchedule(match, 0, 23, scheduleOptions{})
}

func extractDayOfMonth(match string) (Schedule, error) {
	return extractSchedule(match, 1, 31, scheduleOptions{allowW: true})
}

func extractDayOfWeek(match string) (Schedule, error) {
	return extractSchedule(match, 0, 6, scheduleOptions{
		dict: map[string]int{
			"sun": 0, "mon": 1, "tue": 2, "wed": 3,
			"thu": 4, "fri": 5, "sat": 6,
		},
		allowL: true,
		allowH: true,
	})
}

func extractMonth(match string) (Schedule, error) {
	return extractSchedule(match, 1, 12, scheduleOptions{
		dict: map[string]int{
			"jan": 1, "feb": 2, "mar": 3, "apr": 4, "may": 5,
			"jun": 6, "jul": 7, "aug": 8, "sep": 9, "oct": 10,
			"nov": 11, "dec": 12,
		}})
}

func extractEntry(matches []string) (Entry, []error) {
	ent := Entry{}
	errors := make([]error, 0)

	// parseMatch calls 'extractor' while checking for errors. Received
	// errors are appended to "errors". This is less tedious than checking
	// for error after each call.
	type extractorFunc func(string) (Schedule, error)
	parseMatch := func(extractor extractorFunc, match int) Schedule {
		value, err := extractor(matches[match])
		if err != nil {
			errors = append(errors, fmt.Errorf("field %v: %v", match, err))
			return Schedule{}
		}
		return value
	}

	ent.pieces = append(ent.pieces, []Schedule{
		parseMatch(extractMinute, 1),
		parseMatch(extractHour, 2),
		parseMatch(extractDayOfMonth, 3),
		parseMatch(extractMonth, 4),
		parseMatch(extractDayOfWeek, 5),
	}...)

	ent.cmd = matches[6]
	return ent, errors
}

func ParseConfig(data string) ([]Entry, Env, Stats, error) {
	if len(strings.TrimSpace(data)) == 0 {
		// Empty config data
		return []Entry{}, Env{}, Stats{}, nil
	}

	stats := Stats{}
	lines := eolRE.Split(data, -1)
	entries := make([]Entry, 0)
	env := make(Env)

	for i, line := range lines {
		stats.Lines++
		t, matches := getLineType(line)

		if t == UNKNOWN {
			return nil, nil, Stats{}, fmt.Errorf("unrecognized format in line %v: %v", i, line)
		}

		if t == COMMENT {
			stats.Comments++
			continue
		}

		if t == BLANK {
			stats.BlankLines++
			continue
		}

		if t == ENTRY {
			stats.Entries++
			ent, errs := extractEntry(matches)
			if len(errs) > 0 {
				err := fmt.Errorf("parse error(s) in line %v: ", i)
				for e := range errs {
					err = fmt.Errorf("%v %v", err, e)
				}
				return nil, nil, Stats{}, err
			}
			entries = append(entries, ent)
			continue
		}

		if t == ENV {
			stats.Envs++
			env[matches[1]] = matches[2]
			continue
		}
	}

	return entries, env, stats, nil
}
