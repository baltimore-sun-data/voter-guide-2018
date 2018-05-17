package main

import (
	"log"
	"sort"
	"time"
)

type JurisdictionDistrictID struct {
	JurisdictionID
	DistrictID
}

type Reporting struct {
	PrecinctsReporting, TotalPrecincts int
	PercentCounted                     float64
}

type SubResult struct {
	Jurisdiction string
	District     string
	Reporting
	Options []*OptionResult
}

type OptionResult struct {
	Text            string
	Party           string
	Order           int
	TotalVotes      int
	PercentageVotes float64
	FrontRunner     bool
}

type Result struct {
	LastUpdate           time.Time
	Contest              string
	Party                string
	TotalVotes           int
	VoteFor              int
	PrimaryDescription   string
	SecondaryDescription string
	FullDescription      string

	District     string
	Jurisdiction string
	Reporting

	Options    []*OptionResult
	SubResults []*SubResult
	srm        map[JurisdictionDistrictID]*SubResult
}

func MapContestResults(m *Metadata, rc *ResultsContainer) map[ContestID]*Result {
	// Get reporting info for later use
	reporting := map[JurisdictionDistrictID]Reporting{}
	for _, report := range rc.DistrictReporting {
		jdid := JurisdictionDistrictID{
			JurisdictionID: report.JurisdictionID,
			DistrictID:     report.DistrictID,
		}
		reporting[jdid] = Reporting{
			PrecinctsReporting: report.PrecinctsReporting,
			TotalPrecincts:     report.TotalPrecincts,
			PercentCounted:     report.PercentCounted,
		}
	}

	// Create entries for all known contests
	contests := make(map[ContestID]*Result, len(m.Contests))
	for cid, contest := range m.Contests {
		dist, jur := contest.DistrictJurisdiction(m)
		jdid := JurisdictionDistrictID{
			JurisdictionID: contest.JurisdictionID(m),
			DistrictID:     contest.District,
		}
		result := &Result{
			LastUpdate:           rc.LastUpdate,
			Contest:              contest.Name,
			District:             dist,
			Jurisdiction:         jur,
			Party:                contest.Party(m),
			VoteFor:              contest.VoteFor,
			PrimaryDescription:   contest.PrimaryDescription,
			SecondaryDescription: contest.SecondaryDescription,
			FullDescription:      contest.FullDescription,
			Reporting:            reporting[jdid],
			srm:                  make(map[JurisdictionDistrictID]*SubResult),
		}

		contests[cid] = result
	}

	// Keep track of repeated contests
	type rawIDs struct {
		OptionID
		JurisdictionID
		DistrictID
	}
	seen := map[rawIDs]*OptionResult{}

	for _, rawResult := range rc.Results {
		optM, ok := m.Options[rawResult.OptionID]
		if !ok {
			log.Printf("warning: unknown option: %d", rawResult.OptionID)
			continue
		}
		cid := m.OptionParents[rawResult.OptionID]
		result, ok := contests[cid]
		if !ok {
			log.Printf("warning: unknown contest: %d", cid)
			continue
		}

		sid := rawIDs{rawResult.OptionID, rawResult.JurisdictionID, rawResult.DistrictID}
		option := seen[sid]
		if option != nil {
			// We've seen this before!?
			// If the total votes don't match, I guess prefer the bigger one?
			if option.TotalVotes < rawResult.TotalVotes {
				option.TotalVotes = rawResult.TotalVotes
			}
			// Just keep going
			continue
		}

		option = &OptionResult{
			Text:       optM.Text,
			Party:      optM.Party(m),
			Order:      optM.Order,
			TotalVotes: rawResult.TotalVotes,
		}
		seen[sid] = option

		contestM := m.Contests[cid]
		if rawResult.DistrictID == contestM.District &&
			rawResult.JurisdictionID == contestM.JurisdictionID(m) {
			result.Options = append(result.Options, option)
		} else {
			jdid := JurisdictionDistrictID{
				JurisdictionID: rawResult.JurisdictionID,
				DistrictID:     rawResult.DistrictID,
			}
			subres, ok := result.srm[jdid]
			if !ok {
				subDist := rawResult.DistrictID.From(m).Name
				subJur := rawResult.JurisdictionID.From(m).Name
				subres = &SubResult{
					Jurisdiction: subJur,
					District:     subDist,
					Reporting:    reporting[jdid],
				}
				result.SubResults = append(result.SubResults, subres)
				result.srm[jdid] = subres
			}
			subres.Options = append(subres.Options, option)
		}
	}

	// set the total votes / percentage / front-runner
	// sort options by BoE order
	for _, result := range contests {
		// For some local races, the jurisdiction info seems wrong
		if len(result.Options) == 0 && len(result.SubResults) == 1 {
			result.Options = result.SubResults[0].Options
		}

		total := 0
		for _, o := range result.Options {
			total += o.TotalVotes
		}
		if total > 0 {
			result.TotalVotes = total
			tf := float64(total)
			for _, o := range result.Options {
				o.PercentageVotes = float64(o.TotalVotes) / tf * 100
			}
			// Sort to mark the top N front runners
			sort.Slice(result.Options, func(i, j int) bool {
				return result.Options[i].TotalVotes > result.Options[j].TotalVotes
			})
			// Bug: How to deal with ties?
			for i := 0; i < result.VoteFor && i < len(result.Options); i++ {
				result.Options[i].FrontRunner = true
			}
		}
		sort.Slice(result.Options, func(i, j int) bool {
			return result.Options[i].Order < result.Options[j].Order
		})
		sort.Slice(result.SubResults, func(i, j int) bool {
			return result.SubResults[i].District < result.SubResults[j].District
		})

		for _, subr := range result.SubResults {
			sort.Slice(subr.Options, func(i, j int) bool {
				return subr.Options[i].Order < subr.Options[j].Order
			})
		}
	}
	return contests
}
