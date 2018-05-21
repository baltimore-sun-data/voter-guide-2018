package main

import (
	"encoding/json"
	"testing"
)

func TestMapContestResults(t *testing.T) {
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

			results := MapContestResults(&m, &rc)
			for cid, result := range results {
				if result.Contest != m.Contests[cid].Name {
					t.Fatalf("did not set the name on %d", cid)
				}
				if len(result.Options) != len(m.Contests[cid].Options) {
					t.Errorf("contest %d has the wrong number of options %d != %d",
						cid, len(result.Options), len(m.Contests[cid].Options))
				}
				for _, subr := range result.SubResults {
					if len(subr.Options) != len(result.Options) {
						t.Errorf("contest %d has a mismatching number of subresult options %d != %d",
							cid, len(subr.Options), len(result.Options))
					}
				}
			}
		})
	}
}
