package gtfs

import (
	"math"
	"strings"
	"testing"
)

func TestParseFrequenciesHeadwayBounds(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    int32
		wantRow bool
	}{
		{name: "one second", input: "1", want: 1, wantRow: true},
		{name: "normal headway", input: "600", want: 600, wantRow: true},
		{name: "trimmed", input: " 600 ", want: 600, wantRow: true},
		{name: "maximum int32", input: "2147483647", want: math.MaxInt32, wantRow: true},
		{name: "above int32", input: "2147483648"},
		{name: "previously wrapped to positive", input: "4294967297"},
		{name: "maximum int64", input: "9223372036854775807"},
		{name: "above int64", input: "9223372036854775808"},
		{name: "minimum int32", input: "-2147483648"},
		{name: "below int32", input: "-2147483649"},
		{name: "negative", input: "-1"},
		{name: "zero", input: "0"},
		{name: "empty", input: ""},
		{name: "not a number", input: "invalid"},
		{name: "fraction", input: "1.5"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			feed := &Feed{Freqs: make(map[string][]FreqEntry)}
			input := "trip_id,start_time,end_time,headway_secs\ntrip,23:00:00,25:00:00," + tc.input + "\n"
			if err := feed.parseFrequencies(newCSVReader(strings.NewReader(input))); err != nil {
				t.Fatal(err)
			}
			entries := feed.Freqs["trip"]
			if !tc.wantRow {
				if len(entries) != 0 {
					t.Fatalf("invalid headway %q produced entries: %#v", tc.input, entries)
				}
				return
			}
			if len(entries) != 1 {
				t.Fatalf("expected one frequency entry, got %#v", entries)
			}
			if entries[0].HeadwaySec != tc.want || entries[0].StartTime != 23*3600 || entries[0].EndTime != 25*3600 {
				t.Fatalf("unexpected frequency entry: %#v", entries[0])
			}
		})
	}
}

func TestParseFrequenciesSkipsInvalidRowsAndKeepsFollowingRows(t *testing.T) {
	feed := &Feed{Freqs: make(map[string][]FreqEntry)}
	input := "trip_id,start_time,end_time,headway_secs\n" +
		"trip,08:00:00,09:00:00,4294967297\n" +
		"trip,08:00:00,09:00:00,invalid\n" +
		",08:00:00,09:00:00,600\n" +
		"trip,08:00:00,09:00:00,600\n"
	if err := feed.parseFrequencies(newCSVReader(strings.NewReader(input))); err != nil {
		t.Fatal(err)
	}
	entries := feed.Freqs["trip"]
	if len(entries) != 1 || entries[0].HeadwaySec != 600 {
		t.Fatalf("expected only the valid following row, got %#v", feed.Freqs)
	}
}

func TestExpandFrequenciesPreservesDeparturesAndStopOffsets(t *testing.T) {
	feed := &Feed{
		Freqs: map[string][]FreqEntry{"trip": {{StartTime: 23*3600 + 50*60, EndTime: 24*3600 + 30*60, HeadwaySec: 600}}},
		Trips: map[string]*RawTrip{"trip": {ServiceID: "service", RouteID: "route", ShapeID: "shape"}},
		StopTimes: map[string][]RawStopTime{"trip": {
			{StopID: "first", Sequence: 1, Arrival: 100, Dep: 100},
			{StopID: "second", Sequence: 2, Arrival: 160, Dep: 170},
		}},
	}
	trips := ExpandFrequencies(feed)
	if len(trips) != 4 {
		t.Fatalf("expected 4 departures before exclusive end time, got %d", len(trips))
	}
	for i, trip := range trips {
		departure := int32(23*3600 + 50*60 + i*600)
		if trip.ServiceID != "service" || trip.RouteID != "route" || trip.ShapeID != "shape" {
			t.Fatalf("trip metadata changed: %#v", trip)
		}
		if len(trip.StopTimes) != 2 || trip.StopTimes[0].Dep != departure || trip.StopTimes[1].Arrival != departure+60 || trip.StopTimes[1].Dep != departure+70 {
			t.Fatalf("departure/stop offsets changed: %#v", trip.StopTimes)
		}
	}
}

func TestExpandFrequenciesStopsBeforeInt32Overflow(t *testing.T) {
	cases := []struct {
		name  string
		entry FreqEntry
		want  []int32
	}{
		{name: "maximum headway", entry: FreqEntry{StartTime: 3600, EndTime: 7200, HeadwaySec: math.MaxInt32}, want: []int32{3600}},
		{name: "near maximum departure", entry: FreqEntry{StartTime: math.MaxInt32 - 3, EndTime: math.MaxInt32, HeadwaySec: 2}, want: []int32{math.MaxInt32 - 3, math.MaxInt32 - 1}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			feed := &Feed{
				Freqs:     map[string][]FreqEntry{"trip": {tc.entry}},
				Trips:     map[string]*RawTrip{"trip": {ServiceID: "service", RouteID: "route"}},
				StopTimes: map[string][]RawStopTime{"trip": {{StopID: "first", Sequence: 1}}},
			}
			trips := ExpandFrequencies(feed)
			if len(trips) != len(tc.want) {
				t.Fatalf("expected %d departures, got %d", len(tc.want), len(trips))
			}
			for i, want := range tc.want {
				if trips[i].StopTimes[0].Dep != want {
					t.Fatalf("departure %d = %d, want %d", i, trips[i].StopTimes[0].Dep, want)
				}
			}
		})
	}
}
