package main

import "time"

type SubResult struct {
	Jurisdiction string
	District     string
	TotalVotes   int
}

type OptionResult struct {
	Text         string
	TotalVotes   int
	Jurisdiction string
	District     string
	SubResults   []SubResult
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
	om      map[OptionID]*OptionResult
}

func MapContestResults(m *Metadata, rc *ResultsContainer) map[ContestID]*Result {
	contests := make(map[ContestID]*Result)
	for _, rawResult := range rc.Results {
		cid := m.OptionParents[rawResult.OptionID]
		contest := m.Contests[cid]

		result, ok := contests[cid]
		if !ok {
			result = &Result{}
			result.LastUpdate = rc.LastUpdate
			result.Contest = contest.Name
			dist := m.Districts[contest.District]
			result.District = dist.Name
			result.Jurisdiction = m.Jurisdictions[dist.Parent].Name
			result.Party = m.Parties[contest.PartyID].Description
			result.VoteFor = contest.VoteFor
			result.PrimaryDescription = contest.PrimaryDescription
			result.SecondaryDescription = contest.SecondaryDescription
			result.FullDescription = contest.FullDescription
			result.om = make(map[OptionID]*OptionResult)
			contests[cid] = result
		}

		option, ok := result.om[rawResult.OptionID]
		if !ok {
			did := contest.District
			jid := m.Districts[did].Parent
			result.Options = append(result.Options, OptionResult{
				Text:         m.Options[rawResult.OptionID].Text,
				District:     m.Districts[did].Name,
				Jurisdiction: m.Jurisdictions[jid].Name,
				SubResults:   []SubResult{},
			})

			option = &result.Options[len(result.Options)-1]
			result.om[rawResult.OptionID] = option
		}
		if contest.District == rawResult.DistrictID {
			option.TotalVotes = rawResult.TotalVotes
		} else {
			dist := rawResult.DistrictID.From(m)
			jur := rawResult.JurisdictionID.From(m)
			option.SubResults = append(option.SubResults, SubResult{
				District:     dist.Name,
				Jurisdiction: jur.Name,
				TotalVotes:   rawResult.TotalVotes,
			})
		}
	}
	return contests
}
