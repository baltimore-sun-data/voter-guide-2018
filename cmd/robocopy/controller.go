package main

import (
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
	TotalVotes   int
	Reporting
}

type OptionResult struct {
	Text            string
	TotalVotes      int
	PercentageVotes float64
	FrontRunner     bool
	SubResults      []SubResult
	// Party should be here, but it's a primary, so we can omit for now
	order int
}

type Result struct {
	LastUpdate           time.Time
	Contest              string
	District             string
	Jurisdiction         string
	Party                string
	TotalVotes           int
	VoteFor              int
	PrimaryDescription   string
	SecondaryDescription string
	FullDescription      string
	Reporting

	Options []*OptionResult
	om      map[OptionID]int
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

	contests := make(map[ContestID]*Result)
	for _, rawResult := range rc.Results {
		cid := m.OptionParents[rawResult.OptionID]
		contest := m.Contests[cid]
		dist, jur := contest.DistrictJurisdiction(m)
		result, ok := contests[cid]
		if !ok {
			jdid := JurisdictionDistrictID{
				JurisdictionID: contest.JurisdictionID(m),
				DistrictID:     contest.District,
			}
			result = &Result{
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
				om:                   make(map[OptionID]int),
			}
			contests[cid] = result
		}

		pos, ok := result.om[rawResult.OptionID]
		if !ok {
			opt := m.Options[rawResult.OptionID]
			result.Options = append(result.Options, &OptionResult{
				Text:       opt.Text,
				SubResults: []SubResult{},
				order:      opt.Order,
			})
			pos = len(result.Options) - 1
			result.om[rawResult.OptionID] = pos
		}
		option := result.Options[pos]
		if rawResult.DistrictID == contest.District &&
			rawResult.JurisdictionID == contest.JurisdictionID(m) {
			option.TotalVotes = rawResult.TotalVotes
		} else {
			jdid := JurisdictionDistrictID{
				JurisdictionID: rawResult.JurisdictionID,
				DistrictID:     rawResult.DistrictID,
			}

			subDist := rawResult.DistrictID.From(m).Name
			subJur := rawResult.JurisdictionID.From(m).Name
			option.SubResults = append(option.SubResults, SubResult{
				District:     subDist,
				Jurisdiction: subJur,
				TotalVotes:   rawResult.TotalVotes,
				Reporting:    reporting[jdid],
			})
		}
	}

	// set the total votes / percentage / front-runner
	// sort options by BoE order
	for _, result := range contests {
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
			return result.Options[i].order < result.Options[j].order
		})
	}
	return contests
}
