package main

import (
	"encoding/json"
	"sort"
	"time"
)

func (c Contest) Party(m *Metadata) string {
	return c.PartyID.From(m).Code
}

func (c Contest) DistrictJurisdiction(m *Metadata) (dist, jur string) {
	d := c.District.From(m)
	p := d.Parent.From(m)
	return d.Name, p.Name
}

func (c Contest) DistrictName(m *Metadata) string {
	d := c.District.From(m)
	return d.Name
}

func (c Contest) Jurisdiction(m *Metadata) string {
	d := c.District.From(m)
	p := d.Parent.From(m)
	if p.Name == d.Name {
		return ""
	}
	return p.Name
}

func (m *Metadata) MarshalJSON() (b []byte, err error) {
	// Make a list of all contests for use in template
	type raceReturnJSON struct {
		Name         string
		Jurisdiction string
		District     string
		Party        string
		ID           int
	}
	type metadataReturnJSON struct {
		ElectionDate *time.Time
		ElectionType *string
		IsPrimary    *bool
		AllContests  []raceReturnJSON
	}
	var r = metadataReturnJSON{
		ElectionDate: &m.ElectionDate,
		ElectionType: &m.ElectionType,
		IsPrimary:    &m.IsPrimary,
	}

	// sort by BoE order
	cids := make([]ContestID, 0, len(m.Contests))
	for cid := range m.Contests {
		cids = append(cids, cid)
	}
	sort.Slice(cids, func(i, j int) bool {
		return m.Contests[cids[i]].Order < m.Contests[cids[j]].Order
	})

	for _, cid := range cids {
		con := m.Contests[cid]
		dist := m.Districts[con.District]
		jur := m.Jurisdictions[dist.Parent]
		party := m.Parties[con.PartyID]
		r.AllContests = append(r.AllContests, raceReturnJSON{
			Name:         con.Name,
			Jurisdiction: jur.Name,
			District:     dist.Name,
			Party:        party.Description,
			ID:           int(cid),
		})
	}
	return json.Marshal(&r)
}
