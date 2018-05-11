package main

import (
	"sort"
	"time"
)

type SubResult struct {
	Jurisdiction string
	District     string
	TotalVotes   int
}

type OptionResult struct {
	Text       string
	TotalVotes int
	SubResults []SubResult
	// Party should be here, but it's a primary, so we can omit for now
	order int
}

type Result struct {
	LastUpdate           time.Time
	Contest              string
	District             string
	Jurisdiction         string
	Party                string
	VoteFor              int
	PrimaryDescription   string
	SecondaryDescription string
	FullDescription      string

	Options []OptionResult
	om      map[OptionID]int
}

func MapContestResults(m *Metadata, rc *ResultsContainer) map[ContestID]*Result {
	contests := make(map[ContestID]*Result)
	for _, rawResult := range rc.Results {
		cid := m.OptionParents[rawResult.OptionID]
		contest := m.Contests[cid]
		dist, jur := contest.DistrictJurisdiction(m)
		result, ok := contests[cid]
		if !ok {
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
				om:                   make(map[OptionID]int),
			}
			contests[cid] = result
		}

		pos, ok := result.om[rawResult.OptionID]
		if !ok {
			opt := m.Options[rawResult.OptionID]
			result.Options = append(result.Options, OptionResult{
				Text:       opt.Text,
				SubResults: []SubResult{},
				order:      opt.Order,
			})
			pos = len(result.Options) - 1
			result.om[rawResult.OptionID] = pos
		}
		option := &result.Options[pos]
		if rawResult.DistrictID == contest.District &&
			rawResult.JurisdictionID == contest.JurisdictionID(m) {
			option.TotalVotes = rawResult.TotalVotes
		} else {
			subDist := rawResult.DistrictID.From(m).Name
			subJur := rawResult.JurisdictionID.From(m).Name
			option.SubResults = append(option.SubResults, SubResult{
				District:     subDist,
				Jurisdiction: subJur,
				TotalVotes:   rawResult.TotalVotes,
			})
		}
	}

	// sort options by BoE order
	for _, result := range contests {
		sort.Slice(result.Options, func(i, j int) bool {
			return result.Options[i].order < result.Options[j].order
		})
	}
	return contests
}
