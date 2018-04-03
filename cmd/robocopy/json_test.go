package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func readFile(t *testing.T, fn string) []byte {
	t.Helper()
	f, err := os.Open(fn)
	if err != nil {
		t.Fatalf("unexpected error opening file in testing: %v ", err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(BOMReader(f))
	if err != nil {
		t.Fatalf("unexpected read error in testing: %v ", err)
	}
	return b
}

func TestUnmarshallJSON(t *testing.T) {
	var tcs = []struct {
		name  string
		file  string
		value interface{}
	}{
		{"old-metadata", "test/Metadata.js", &Metadata{}},
		{"new-metadata", "test/GP18-Metadata.js", &Metadata{}},
		{"old-precinct-results", "test/PrecinctResults.js", &PrecinctResults{}},
		{"old-results", "test/Results.js", &Results{}},
		{"new-precinct-results", "test/GP18-PrecinctResults.js", &PrecinctResults{}},
		{"new-results", "test/GP18-Results.js", &Results{}},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			b := readFile(t, tc.file)
			if err := json.Unmarshal(b, &tc.value); err != nil {
				t.Fatalf("could not unmarshall JSON: %v ", err)
			}
		})
	}
}

func TestMetadataObj(t *testing.T) {
	var tcs = []struct {
		name string
		file string
	}{
		{"new-metadata", "test/GP18-Metadata.js"},
		{"old-metadata", "test/Metadata.js"},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			b := readFile(t, tc.file)

			var m Metadata
			if err := json.Unmarshal(b, &m); err != nil {
				t.Fatalf("could not unmarshall JSON: %v ", err)
			}

			for _, c := range m.Contests {
				if c.Type != Race {
					if c.PartyID != 0 {
						t.Errorf("PartyID present for non-race contest: %#v", c)
					}
				} else if _, ok := m.Parties[c.PartyID]; !ok {
					t.Errorf("contest has invalid party ID: %#v", c)
				}
				for _, o := range c.Options {
					if c.Type != Race && o.WriteIn != 'N' {
						t.Errorf("option with unexpected WriteIn type: %#v", o)
					}
				}
			}

			for _, d := range m.Districts {
				if _, ok := m.Jurisdictions[d.Parent]; !ok {
					t.Errorf("district has invalid parent: %#v", d)
				}
				if _, ok := m.DistrictTypes[d.Type]; !ok {
					t.Errorf("district has invalid type: %#v", d)
				}
			}

			for _, j := range m.Jurisdictions {
				var bad []DistrictID
				for _, d := range j.AllDistricts {
					if _, ok := m.Districts[d]; !ok {
						bad = append(bad, d)
					}
				}
				if len(bad) > 0 {
					t.Errorf("jurisdiction %v has invalid children: %v", j.Name, bad)
				}
			}
		})
	}
}
