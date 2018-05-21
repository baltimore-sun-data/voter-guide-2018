package main

import (
	"fmt"
	"time"
)

const (
	Race        ContestType = 0
	Continuance ContestType = 1
	Question    ContestType = 2

	RegularOption                 = 'N'
	RegisteredWriteInCandidate    = 'Y'
	UnregisteredWriteInCandidates = 'O'
	AllWriteInCandidates          = 'A'
)

type Metadata struct {
	ElectionDate  time.Time
	ElectionType  string
	IsPrimary     bool
	Contests      map[ContestID]Contest
	Options       map[OptionID]*Option
	OptionParents map[OptionID]ContestID
	DistrictTypes map[DistrictTypeID]DistrictType
	Districts     map[DistrictID]District
	Jurisdictions map[JurisdictionID]Jurisdiction
	Parties       map[PartyID]Party
}

func MetadataFrom(name string) (*Metadata, error) {
	var m Metadata
	if err := unmarshalFrom(name, &m); err != nil {
		return nil, fmt.Errorf("could not read metadata: %v", err)
	}
	return &m, nil
}

type (
	OptionID int
	Option   struct {
		ContestID
		ID      OptionID
		PartyID PartyID
		Order   int
		Text    string
		WriteIn byte
	}
	ContestID   int
	ContestType int
	Contest     struct {
		BallotDistrict       DistrictID
		District             DistrictID
		ID                   ContestID
		Type                 ContestType
		PartyID              PartyID
		Name                 string
		Order                int
		VoteFor              int
		PrimaryDescription   string
		SecondaryDescription string
		FullDescription      string
		Options              []Option
	}
	DistrictTypeID int
	DistrictType   struct {
		ID          DistrictTypeID
		Description string
	}
	DistrictID int
	District   struct {
		ID     DistrictID
		Parent JurisdictionID
		Type   DistrictTypeID
		Name   string
	}
	JurisdictionID int
	Jurisdiction   struct {
		AllDistricts []DistrictID
		ID           JurisdictionID
		Name         string
	}
	PartyID int
	Party   struct {
		ID                PartyID
		Code, Description string
	}
)

func (id ContestID) From(m *Metadata) Contest {
	return m.Contests[id]
}

func (id DistrictTypeID) From(m *Metadata) DistrictType {
	return m.DistrictTypes[id]
}

func (id DistrictID) From(m *Metadata) District {
	return m.Districts[id]
}

func (id JurisdictionID) From(m *Metadata) Jurisdiction {
	return m.Jurisdictions[id]
}

func (id PartyID) From(m *Metadata) Party {
	return m.Parties[id]
}

func (c Contest) Party(m *Metadata) string {
	return c.PartyID.From(m).Code
}

func (c Contest) JurisdictionID(m *Metadata) JurisdictionID {
	d := c.District.From(m)
	return d.Parent
}

func (c Contest) DistrictJurisdiction(m *Metadata) (dist, jur string) {
	d := c.District.From(m)
	p := d.Parent.From(m)
	return d.Name, p.Name
}

func (o Option) Party(m *Metadata) string {
	return o.PartyID.From(m).Description
}
