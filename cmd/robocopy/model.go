package main

import (
	"encoding/json"
	"sort"
	"strings"
	"time"
)

func (c Contest) Party(m Metadata) string {
	return c.PartyID.From(m).Code
}

func (c Contest) DistrictName(m Metadata) string {
	d := c.District.From(m)
	return d.Name
}

func (c Contest) Jurisdiction(m Metadata) string {
	d := c.District.From(m)
	p := d.Parent.From(m)
	if p.Name == d.Name {
		return ""
	}
	return p.Name
}

func (m *Metadata) MarshalJSON() (b []byte, err error) {
	// This JSON is a turkducken of crud trying to invert the contests into something useful
	type contestReturnJSON struct {
		Name  string
		Party string
		ID    int
	}
	type districtReturnJSON struct {
		Name     string
		Contests []contestReturnJSON
	}
	type jurisdictionReturnJSON struct {
		Name      string
		Districts []districtReturnJSON
		dmap      map[DistrictID]*districtReturnJSON
	}
	type raceReturnJSON struct {
		Name          string
		Jurisdictions []jurisdictionReturnJSON
		jmap          map[JurisdictionID]*jurisdictionReturnJSON
	}
	type metadataReturnJSON struct {
		ElectionDate *time.Time
		ElectionType *string
		IsPrimary    *bool
		Contests     []raceReturnJSON
	}
	var r = metadataReturnJSON{
		ElectionDate: &m.ElectionDate,
		ElectionType: &m.ElectionType,
		IsPrimary:    &m.IsPrimary,
	}
	// Do a lot of work to group things together logically
	cids := make([]ContestID, 0, len(m.Contests))
	for cid := range m.Contests {
		cids = append(cids, cid)
	}
	sort.Slice(cids, func(i, j int) bool {
		return m.Contests[cids[i]].Order < m.Contests[cids[j]].Order
	})
	races := map[string]*raceReturnJSON{}
	for _, cid := range cids {
		raceName := m.Contests[cid].Name
		raceName = strings.TrimSuffix(raceName, " At Large")
		raceName = strings.TrimSuffix(raceName, " Male")
		raceName = strings.TrimSuffix(raceName, " Female")
		race, ok := races[raceName]
		if !ok {
			r.Contests = append(r.Contests, raceReturnJSON{
				Name: raceName,
				jmap: make(map[JurisdictionID]*jurisdictionReturnJSON),
			})
			race = &r.Contests[len(r.Contests)-1]
			races[raceName] = race
		}

		did := m.Contests[cid].District
		jid := m.Districts[did].Parent
		jur, ok := race.jmap[jid]
		if !ok {
			race.Jurisdictions = append(race.Jurisdictions, jurisdictionReturnJSON{
				Name: m.Jurisdictions[jid].Name,
				dmap: make(map[DistrictID]*districtReturnJSON),
			})
			jur = &race.Jurisdictions[len(race.Jurisdictions)-1]
			race.jmap[jid] = jur
		}

		dist, ok := jur.dmap[did]
		if !ok {
			jur.Districts = append(jur.Districts, districtReturnJSON{
				Name: m.Jurisdictions[jid].Name,
			})
			dist = &jur.Districts[len(jur.Districts)-1]
		}

		dist.Contests = append(dist.Contests, contestReturnJSON{
			Name:  m.Contests[cid].Name,
			Party: m.Parties[m.Contests[cid].PartyID].Code,
			ID:    int(m.Contests[cid].ID),
		})

	}
	return json.Marshal(&r)
}
