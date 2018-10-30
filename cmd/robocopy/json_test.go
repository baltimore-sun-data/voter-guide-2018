package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func readSource(t *testing.T, name string) []byte {
	t.Helper()
	rc, err := readFrom(name)
	if err != nil {
		t.Fatalf("unexpected error opening source in testing: %v ", err)
	}
	defer rc.Close()

	b, err := ioutil.ReadAll(BOMReader(rc))
	if err != nil {
		t.Fatalf("unexpected error reading source in testing: %v ", err)
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
		{"old-results", "test/Results.js", &ResultsContainer{}},
		{"new-precinct-results", "test/GP18-PrecinctResults.js", &PrecinctResults{}},
		{"new-results", "test/GP18-Results.js", &ResultsContainer{}},
		{"live-metadata", metadata18url, &Metadata{}},
		{"live-precinct-results", precinctResults18url, &PrecinctResults{}},
		{"live-results", results18url, &ResultsContainer{}},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			b := readSource(t, tc.file)
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
		{"live-metadata", metadata18url},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			b := readSource(t, tc.file)

			var m Metadata
			if err := json.Unmarshal(b, &m); err != nil {
				t.Fatalf("could not unmarshall JSON: %v ", err)
			}

			for _, c := range m.Contests {
				if c.Type != Race {
					if c.PartyID != 0 && c.PartyID.From(&m).Code != "Non Partisan" {
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

func TestResultsMetadata(t *testing.T) {
	var tcs = []struct {
		name         string
		metadatafile string
		resultsfile  string
	}{
		{"old-files", "test/Metadata.js", "test/Results.js"},
		{"new-files", "test/GP18-Metadata.js", "test/GP18-Results.js"},
		{"live-files", metadata18url, results18url},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			b := readSource(t, tc.metadatafile)
			var m Metadata
			if err := json.Unmarshal(b, &m); err != nil {
				t.Fatalf("could not unmarshall JSON: %v ", err)
			}

			b = readSource(t, tc.resultsfile)
			var rc ResultsContainer
			if err := json.Unmarshal(b, &rc); err != nil {
				t.Fatalf("could not unmarshall JSON: %v ", err)
			}

			for _, r := range rc.Results {
				if _, ok := m.Options[r.OptionID]; !ok {
					t.Fatalf("results have missing option: %d", r.OptionID)
				}
				cid := m.OptionParents[r.OptionID]
				contest, ok := m.Contests[cid]
				if !ok {
					t.Fatalf("bad contest id %d", cid)

				}
				dist, ok := m.Districts[r.DistrictID]
				if !ok {
					t.Fatalf("bad district %d", r.DistrictID)
				}

				if r.DistrictID != contest.District {
					if dist.Parent != m.Districts[contest.District].Parent {
						t.Fatalf("option %d has unexpected district parents: %d != %d",
							r.OptionID, dist.Parent, m.Districts[contest.District].Parent)
					}
				}
				if dist.Parent != 0 && r.JurisdictionID != dist.Parent {
					t.Fatalf("option %d has unexpected jurisdiction: %d != %d",
						r.OptionID, r.JurisdictionID, dist.Parent)
				}
			}
		})
	}
}
