// Developed by Bombs King, Taipei Codefest 2026

package raptor

import (
	"archive/zip"
	"bytes"
	"testing"

	"TaipeiCityDashboardBE/app/services/isochrone/gtfs"
)

func TestBuildDoesNotScheduleInvalidFrequencyTemplates(t *testing.T) {
	cases := []struct {
		name     string
		headways []string
		wantDeps []int32
	}{
		{name: "zero", headways: []string{"0"}, wantDeps: []int32{28800}},
		{name: "negative", headways: []string{"-1"}, wantDeps: []int32{28800}},
		{name: "malformed", headways: []string{"invalid"}, wantDeps: []int32{28800}},
		{name: "out of range", headways: []string{"4294967297"}, wantDeps: []int32{28800}},
		{name: "valid then invalid", headways: []string{"600", "invalid"}, wantDeps: []int32{28800, 32400, 33000}},
		{name: "invalid then valid", headways: []string{"invalid", "600"}, wantDeps: []int32{28800, 32400, 33000}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			frequencies := "trip_id,start_time,end_time,headway_secs\n"
			for _, headway := range tc.headways {
				frequencies += "frequency,09:00:00,09:20:00," + headway + "\n"
			}
			files := map[string]string{
				"stops.txt": "stop_id,stop_name,stop_lat,stop_lon\n" +
					"a,First,25.0478,121.5174\nb,Second,25.0520,121.5200\n",
				"routes.txt": "route_id,route_type\nr,3\n",
				"trips.txt":  "route_id,service_id,trip_id\nr,svc,frequency\nr,svc,scheduled\n",
				"stop_times.txt": "trip_id,arrival_time,departure_time,stop_id,stop_sequence\n" +
					"frequency,06:00:00,06:00:00,a,1\nfrequency,06:05:00,06:05:00,b,2\n" +
					"scheduled,08:00:00,08:00:00,a,1\nscheduled,08:05:00,08:05:00,b,2\n",
				"frequencies.txt": frequencies,
			}
			var archive bytes.Buffer
			writer := zip.NewWriter(&archive)
			for name, contents := range files {
				file, err := writer.Create(name)
				if err != nil {
					t.Fatal(err)
				}
				if _, err := file.Write([]byte(contents)); err != nil {
					t.Fatal(err)
				}
			}
			if err := writer.Close(); err != nil {
				t.Fatal(err)
			}
			feed, err := gtfs.LoadFeedFromZip(archive.Bytes(), "bus:")
			if err != nil {
				t.Fatal(err)
			}
			rd, err := Build([]*gtfs.Feed{feed})
			if err != nil {
				t.Fatal(err)
			}
			if len(rd.Routes) != 1 || len(rd.Routes[0].Trips) != len(tc.wantDeps) {
				t.Fatalf("invalid frequency template was scheduled or valid trips lost: %#v", rd.Routes)
			}
			for i, want := range tc.wantDeps {
				if got := rd.Routes[0].Trips[i].Times[0].Departure; got != want {
					t.Fatalf("departure %d = %d, want %d", i, got, want)
				}
			}
		})
	}
}

func TestBuildMergesNearbySameNameBusStops(t *testing.T) {
	feed := testFeed("bus:", map[string]*gtfs.RawStop{
		"a": {ID: "a", Name: "Taipei Main", Lat: 25.04780, Lon: 121.51740},
		"b": {ID: "b", Name: "Taipei Main", Lat: 25.04790, Lon: 121.51750},
		"c": {ID: "c", Name: "Zhongshan", Lat: 25.05200, Lon: 121.52000},
	}, []gtfs.RawStopTime{
		{StopID: "a", Sequence: 1, Arrival: 28800, Dep: 28800},
		{StopID: "c", Sequence: 2, Arrival: 29100, Dep: 29100},
	}, []gtfs.RawStopTime{
		{StopID: "b", Sequence: 1, Arrival: 28800, Dep: 28800},
		{StopID: "c", Sequence: 2, Arrival: 29100, Dep: 29100},
	})

	rd, err := Build([]*gtfs.Feed{feed})
	if err != nil {
		t.Fatal(err)
	}

	if rd.StopIndex["bus:a"] != rd.StopIndex["bus:b"] {
		t.Fatalf("expected same-name nearby bus stops to merge, got %d and %d", rd.StopIndex["bus:a"], rd.StopIndex["bus:b"])
	}
	if len(rd.Stops) != 2 {
		t.Fatalf("expected 2 merged stops, got %d", len(rd.Stops))
	}
}

func TestBuildKeepsDistantSameNameBusStopsSeparate(t *testing.T) {
	feed := testFeed("bus:", map[string]*gtfs.RawStop{
		"a": {ID: "a", Name: "Same Name", Lat: 25.04780, Lon: 121.51740},
		"b": {ID: "b", Name: "Same Name", Lat: 25.14780, Lon: 121.61740},
	}, []gtfs.RawStopTime{
		{StopID: "a", Sequence: 1, Arrival: 28800, Dep: 28800},
		{StopID: "b", Sequence: 2, Arrival: 29100, Dep: 29100},
	})

	rd, err := Build([]*gtfs.Feed{feed})
	if err != nil {
		t.Fatal(err)
	}

	if rd.StopIndex["bus:a"] == rd.StopIndex["bus:b"] {
		t.Fatalf("expected distant same-name bus stops to stay separate")
	}
	if len(rd.Stops) != 2 {
		t.Fatalf("expected 2 stops, got %d", len(rd.Stops))
	}
}

func testFeed(prefix string, stops map[string]*gtfs.RawStop, trips ...[]gtfs.RawStopTime) *gtfs.Feed {
	feed := &gtfs.Feed{
		Prefix:    prefix,
		Stops:     stops,
		Routes:    map[string]*gtfs.RawRoute{"r": {ID: "r"}},
		Trips:     make(map[string]*gtfs.RawTrip),
		StopTimes: make(map[string][]gtfs.RawStopTime),
		Calendar:  make(map[string]*gtfs.ServicePattern),
		CalDates:  make(map[string][]gtfs.CalDateException),
		Freqs:     make(map[string][]gtfs.FreqEntry),
	}
	for i, stopTimes := range trips {
		tripID := string(rune('a' + i))
		feed.Trips[tripID] = &gtfs.RawTrip{ID: tripID, RouteID: "r", ServiceID: "svc"}
		feed.StopTimes[tripID] = stopTimes
	}
	return feed
}
