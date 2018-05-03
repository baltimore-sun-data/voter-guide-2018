package main

import (
	"encoding/json"
	"fmt"
	"time"
)

/*
   {
         ElectionType:<”Gubernatorial”,”Presidential”, “Special Legislative”, etc>,
         ElectionDate:<Dateof Election>,
         IsPrimary:<if this is a primary>,
         Contests: [
         {
                     ID: <unique contest id>,//a unique id for this contest, across all ballot styles and jurisdictions
                     Order: <Universal Ballot Order>,//Universal Ballot order across all ballot styles and jurisdictions, i.e. web order
                     Name:<Contest Name>,//Contest Name
                     Primary Description:<Contest Name>,//Candidate Name for type 1, and question type description for type 2
                     Secondary Description:<Contest Name>,//”for continuance in office” for type 1, question short description (if any) for type 2
                     Full Description:<Contest Name>,//Full question text for type 2 (no formatting except for for line breaks)
                     PartyID:<Contest Party if Applicable>,//Usually Ballot Party, Except for when NON, which is actually on every party’s ballot
                     ContestJurisdiction:<Jurisdiction ID>,//Jursidiction this Race is in (State, County, or Baltimore City)
                     ContestDistrict:<Contest District ID>,//District this Contest appliesto/is for
                     BallotDistrict:<Ballot District ID>,//District in which this contest appears on the ballot ifdifferent from ContestDistrict
                     Type:<Contest type, 0=Race, 1=Continuance, 2=Question>,//if this contest is a race,forcontinuance in office,or is question/referendum etc
                     VoteFor:<Vote For No More Than…>,//Max votes allowed for this contest
                     Options:[
                           {
                                 // an option represents a single option on a ballot for this contest, i.e. a yes/for, or a no/against for a question/continuance, or a candidate name for a race
                                 // with the exception of write-in options, which represent a certified write-in, a sum of non-certified write-ins, or a sum of all write-ins
                                 ID: <Unique option id>,//a unique id for this option, across all contests
                                 Text:<Option Text>,//Candidate name or question/continuance answer
                                 PartyID:<Party Affiliation if any>,//Only applicable to candidates, and only in general elections, this is the party of this particular candidate
                                 IsWriteIn:<Y/N/O/A>//Only for a candidate in a general –
                                                         //a registered write-in candidate (Y),
                                                         //a regular candidate (N),
                                                         //a sum of non-certified write-ins(O),
                                                         //a sum of all write-ins(A)
                           },
                           …
                     ]
               },
               …
         ],
         DistrictTypes:[
                 {
                           //Types of Districts
                           ID: <unique district type id>,
                           Name:<District Type Name>,
                 },
                 …
         ],
         Parties:[
                 {
                           //Political/Ballot Parties
                           ID: <unique party type id>,
                           Code:<NON, REP, DEM, etc>,
                           Description:<Non-Partisan, Republican, Democrat, etc>,
                 },
                 …
         ],
         Districts:[
                 {
                           //Race & Result Districts (Congressional District 1, Allegany County, Cumberland, State Of Maryland, etc)
                           ID: <unique district id>,
                           Name:<District Name>,
                           JurisdictionID:<Parent Jurisdiction ID>,
                           TypeID:<District Type ID>
                 },
                 …
         ],
         Precincts:[
                 {
                           //Precincts
                           ID: <unique precinct id>,
                           Name:<Precinct Name i.e “001-001”>,
                           ActiveQualified:<Total Number of Active Qualified Voters (aka with a ballot) as of Close of Registration>,
                           JurisdictionID:<Precinct JurisdictionID>
                           VoteCenters:[<Precinct Vote Center ID>, <Early Voting Site 1 ID>,…],
                           RaceDistricts:[< district id (i.e. cong, leg, county, municipal, statewide >, …]
                 },
                 …
         ],
         VoteCenter:[
                 {
                           //Polling Places and Early Vote Centers
                           ID: <unique polling place id>,
                           Name:<Polling Place Name>,
                           Address:<Polling Place Address>,
                           IsEarly:<If this an early vote center>,
                           JurisdictionID:<Vote Center JurisdictionID>
                 },
                 …
         ],
         Jurisdictions:[
               {
                           // State, Counties, and Baltimore city
                           ID: <unique jurisdiction id>,
                           Name:<Jurisdiction Name(string)>,
                           AllDistricts: [< district id (i.e. cong, leg, county, municipal, statewide >, …]//All districts used in this election for any/all precincts in this jurisdiction
               },
               …
         ]
   }
*/
type metadataJSON struct {
	Contests []struct {
		BallotDistrict int    `json:"BallotDistrict"`
		District       int    `json:"District"`
		ID             int    `json:"ID"`
		Jurisdiction   int    `json:"Jurisdiction"`
		Name           string `json:"Name"`
		Options        []struct {
			ID      int    `json:"ID"`
			Order   int    `json:"Order"`
			PartyID int    `json:"PartyID,string"`
			Text    string `json:"Text"`
			WriteIn string `json:"WriteIn"`
		} `json:"Options"`
		Order                int    `json:"Order"`
		PartyID              int    `json:"PartyID,string"`
		Type                 int    `json:"Type"`
		VoteFor              int    `json:"VoteFor"`
		PrimaryDescription   string `json:"PrimaryDescription,omitempty"`
		SecondaryDescription string `json:"SecondaryDescription,omitempty"`
		FullDescription      string `json:"FullDescription,omitempty"`
	} `json:"Contests"`
	DistrictTypes []struct {
		Description string `json:"Description"`
		ID          int    `json:"ID"`
	} `json:"DistrictTypes"`
	Districts []struct {
		ID             int    `json:"ID"`
		JurisdictionID int    `json:"JurisdictionID"`
		Name           string `json:"Name"`
		TypeID         int    `json:"TypeID"`
	} `json:"Districts"`
	ElectionDate  string `json:"ElectionDate"`
	ElectionType  string `json:"ElectionType"`
	IsPrimary     string `json:"IsPrimary"`
	Jurisdictions []struct {
		AllDistricts []int  `json:"AllDistricts"`
		ID           int    `json:"ID"`
		Name         string `json:"Name"`
	} `json:"Jurisdictions"`
	Parties []struct {
		Code        string `json:"Code"`
		Description string `json:"Description"`
		ID          int    `json:"ID"`
	} `json:"Parties"`
	Precincts []struct {
		ActiveQualified map[string]int `json:"ActiveQualified"`
		ID              int            `json:"ID"`
		JurisdictionID  int            `json:"JurisdictionID"`
		Name            string         `json:"Name"`
		RaceDistricts   []int          `json:"RaceDistricts"`
		VoteCenters     []int          `json:"VoteCenters"`
	} `json:"Precincts"`
	VoteCenters []struct {
		Address        string `json:"Address"`
		ID             int    `json:"ID"`
		IsEarly        int    `json:"IsEarly"`
		JurisdictionID int    `json:"JurisdictionID"`
		Name           string `json:"Name"`
	} `json:"VoteCenters"`
}

type (
	OptionID int
	Option   struct {
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

const (
	Race        ContestType = 0
	Continuance ContestType = 1
	Question    ContestType = 2

	RegularOption                 = 'N'
	RegisteredWriteInCandidate    = 'Y'
	UnregisteredWriteInCandidates = 'O'
	AllWriteInCandidates          = 'A'
)

func (id ContestID) From(m Metadata) Contest {
	return m.Contests[id]
}

func (id DistrictTypeID) From(m Metadata) DistrictType {
	return m.DistrictTypes[id]
}

func (id DistrictID) From(m Metadata) District {
	return m.Districts[id]
}

func (id JurisdictionID) From(m Metadata) Jurisdiction {
	return m.Jurisdictions[id]
}

func (id PartyID) From(m Metadata) Party {
	return m.Parties[id]
}

type Metadata struct {
	ElectionDate  time.Time
	ElectionType  string
	IsPrimary     bool
	Contests      map[ContestID]Contest
	DistrictTypes map[DistrictTypeID]DistrictType
	Districts     map[DistrictID]District
	Jurisdictions map[JurisdictionID]Jurisdiction
	Parties       map[PartyID]Party
}

func MetadataFrom(name string) (m *Metadata, err error) {
	rc, err := readFrom(name)
	if err != nil {
		return nil, fmt.Errorf("could not open metadata: %v", err)
	}
	defer deferClose(&err, rc.Close)

	m = &Metadata{}
	dec := json.NewDecoder(BOMReader(rc))
	if err = dec.Decode(&m); err != nil {
		return nil, fmt.Errorf("could not decode metadata: %v", err)
	}
	return m, err
}

func (m *Metadata) UnmarshalJSON(b []byte) error {
	var raw metadataJSON
	err := json.Unmarshal(b, &raw)
	if err != nil {
		return err
	}

	m.ElectionDate, err = time.Parse("2006-01-02", raw.ElectionDate)
	if err != nil {
		return err
	}
	m.ElectionType = raw.ElectionType
	m.IsPrimary = raw.IsPrimary == "Yes"

	m.Contests = make(map[ContestID]Contest, len(raw.Contests))
	for _, c := range raw.Contests {
		nc := Contest{
			BallotDistrict:       DistrictID(c.BallotDistrict),
			District:             DistrictID(c.District),
			ID:                   ContestID(c.ID),
			Type:                 ContestType(c.Type),
			PartyID:              PartyID(c.PartyID),
			Name:                 c.Name,
			Order:                c.Order,
			VoteFor:              c.VoteFor,
			PrimaryDescription:   c.PrimaryDescription,
			SecondaryDescription: c.SecondaryDescription,
			FullDescription:      c.FullDescription,
		}
		for _, o := range c.Options {
			if len(o.WriteIn) != 1 {
				return fmt.Errorf("Unexpected WriteIn value: %s", o.WriteIn)
			}
			no := Option{
				ID:      OptionID(o.ID),
				PartyID: PartyID(o.PartyID),
				Order:   o.Order,
				Text:    o.Text,
				WriteIn: o.WriteIn[0],
			}
			nc.Options = append(nc.Options, no)
		}
		m.Contests[nc.ID] = nc
	}
	if len(raw.Contests) != len(m.Contests) {
		return fmt.Errorf("unexpected ContestID repeat: %d != %d",
			len(raw.Contests), len(m.Contests))
	}

	m.DistrictTypes = make(map[DistrictTypeID]DistrictType, len(raw.DistrictTypes))
	for _, dt := range raw.DistrictTypes {
		ndt := DistrictType{
			ID:          DistrictTypeID(dt.ID),
			Description: dt.Description,
		}
		m.DistrictTypes[ndt.ID] = ndt
	}
	if len(raw.DistrictTypes) != len(m.DistrictTypes) {
		return fmt.Errorf("unexpected DistrictTypeID repeat: %d != %d",
			len(raw.DistrictTypes), len(m.DistrictTypes))
	}

	m.Districts = make(map[DistrictID]District, len(raw.Districts))
	for _, d := range raw.Districts {
		nd := District{
			ID:     DistrictID(d.ID),
			Parent: JurisdictionID(d.JurisdictionID),
			Type:   DistrictTypeID(d.TypeID),
			Name:   d.Name,
		}
		m.Districts[nd.ID] = nd
	}
	if len(raw.Districts) != len(m.Districts) {
		return fmt.Errorf("unexpected DistrictID repeat: %d != %d",
			len(raw.Districts), len(m.Districts))
	}

	m.Jurisdictions = make(map[JurisdictionID]Jurisdiction, len(raw.Jurisdictions))
	for _, j := range raw.Jurisdictions {
		nj := Jurisdiction{
			ID:           JurisdictionID(j.ID),
			Name:         j.Name,
			AllDistricts: make([]DistrictID, len(j.AllDistricts)),
		}
		for i := range j.AllDistricts {
			nj.AllDistricts[i] = DistrictID(j.AllDistricts[i])
		}
		m.Jurisdictions[nj.ID] = nj
	}
	if len(raw.Jurisdictions) != len(m.Jurisdictions) {
		return fmt.Errorf("unexpected JurisdictionID repeat: %d != %d",
			len(raw.Jurisdictions), len(m.Jurisdictions))
	}

	m.Parties = make(map[PartyID]Party, len(raw.Parties))
	for _, p := range raw.Parties {
		np := Party{
			ID:          PartyID(p.ID),
			Code:        p.Code,
			Description: p.Description,
		}
		m.Parties[np.ID] = np
	}
	if len(raw.Parties) != len(m.Parties) {
		return fmt.Errorf("unexpected PartyID repeat: %d != %d",
			len(raw.Parties), len(m.Parties))
	}

	return err
}

type Results struct {
	Test       bool
	LastUpdate time.Time
	Results    []struct {
		OptionID, JurisdictionID, DistrictID, TotalVotes int
	}
	DistrictReporting []struct {
		JurisdictionID, DistrictID, PrecinctsReporting, TotalPrecincts int
		PercentCounted                                                 float64
	}
	Reporting []struct {
		PrecinctID       int
		PercentReporting float64
		BallotsCounted   int
	}
	EarlyVotingReporting []struct {
		EarlyVotingCenterID int
		PercentReporting    float64
		BallotsCounted      int
	}
	CanvasCounts []struct {
		JurisdictionID int
		BallotsCounted int
	}
}

func (r *Results) UnmarshalJSON(b []byte) error {
	/*
	   Test:<If this is a test feed>,
	   LastUpdate:<Dateand Time of Last Updatein the format “YYYY-MM-DD HH:mm:SS”>,
	   Results: [[<Option ID>, <Jurisdiction ID>, <District ID>, <Total Votes (within specified district within specified jurisdiction)>], …],
	   DReporting: [[<JurisdictionID>, <DistrictID>, <Total Precincts Currently Reporting>, <Total Precincts (in this district in this jurisdiction)>, <Percent Counted>], …],
	   Reporting: [[<Precinct ID>, <Percent Reporting>, <Ballots Counted>], …],
	   EVReporting: [[<EV Vote Center ID>, <Percent Reporting>, <Ballots Counted>], …],
	   CanvasCounts: [[<Jurisdiction ID>, <Ballots Counted>], …],
	*/
	var raw struct {
		Test         int          `json:"Test"`
		LastUpdate   string       `json:"LastUpdate"`
		Results      [][4]int     `json:"Results"`
		DReporting   [][5]float64 `json:"DReporting"`
		Reporting    [][3]float64 `json:"Reporting"`
		EVReporting  [][3]float64 `json:"EVReporting"`
		CanvasCounts [][2]int     `json:"CanvasCounts"`
	}
	err := json.Unmarshal(b, &raw)
	if err != nil {
		return err
	}
	r.Test = raw.Test == 0
	r.LastUpdate, err = time.Parse("2006/01/02 15:04:05", raw.LastUpdate)
	if err != nil {
		return err
	}
	r.Results = make([]struct {
		OptionID, JurisdictionID, DistrictID, TotalVotes int
	}, len(raw.Results))
	for i := range raw.Results {
		r.Results[i].OptionID = raw.Results[i][0]
		r.Results[i].JurisdictionID = raw.Results[i][1]
		r.Results[i].DistrictID = raw.Results[i][2]
		r.Results[i].TotalVotes = raw.Results[i][3]
	}
	r.DistrictReporting = make([]struct {
		JurisdictionID, DistrictID, PrecinctsReporting, TotalPrecincts int
		PercentCounted                                                 float64
	}, len(raw.DReporting))
	for i := range raw.DReporting {
		r.DistrictReporting[i].JurisdictionID = int(raw.DReporting[i][0])
		r.DistrictReporting[i].DistrictID = int(raw.DReporting[i][1])
		r.DistrictReporting[i].PrecinctsReporting = int(raw.DReporting[i][2])
		r.DistrictReporting[i].TotalPrecincts = int(raw.DReporting[i][3])
		r.DistrictReporting[i].PercentCounted = raw.DReporting[i][4]
	}
	r.Reporting = make([]struct {
		PrecinctID       int
		PercentReporting float64
		BallotsCounted   int
	}, len(raw.Reporting))
	for i := range raw.Reporting {
		r.Reporting[i].PrecinctID = int(raw.Reporting[i][0])
		r.Reporting[i].PercentReporting = raw.Reporting[i][1]
		r.Reporting[i].BallotsCounted = int(raw.Reporting[i][2])
	}
	r.EarlyVotingReporting = make([]struct {
		EarlyVotingCenterID int
		PercentReporting    float64
		BallotsCounted      int
	}, len(raw.EVReporting))
	for i := range raw.EVReporting {
		r.EarlyVotingReporting[i].EarlyVotingCenterID = int(raw.EVReporting[i][0])
		r.EarlyVotingReporting[i].PercentReporting = raw.EVReporting[i][1]
		r.EarlyVotingReporting[i].BallotsCounted = int(raw.EVReporting[i][2])
	}
	r.CanvasCounts = make([]struct {
		JurisdictionID int
		BallotsCounted int
	}, len(raw.CanvasCounts))
	for i := range raw.CanvasCounts {
		r.CanvasCounts[i].JurisdictionID = raw.CanvasCounts[i][0]
		r.CanvasCounts[i].BallotsCounted = raw.CanvasCounts[i][1]
	}

	return err
}

type PrecinctResults struct {
	Test       bool
	LastUpdate time.Time
	Results    []struct {
		OptionID, PrecinctID, TotalVotes int
	}
}

func (p *PrecinctResults) UnmarshalJSON(b []byte) error {
	/*
	   {
	    Test:<If this is a test feed>,
	    LastUpdate:<Dateand Time of Last Updatein the format “YYYY-MM-DD HH:mm:SS”>,
	    Results: [[<Option ID>, <Precinct ID>, <Total Votes>], …]
	   }
	*/
	var raw struct {
		Test            int      `json:"Test"`
		LastUpdate      string   `json:"LastUpdate"`
		PrecinctResults [][3]int `json:"PrecinctResults"`
	}
	err := json.Unmarshal(b, &raw)
	if err != nil {
		return err
	}
	p.Test = raw.Test == 0
	p.LastUpdate, err = time.Parse("2006/01/02 15:04:05", raw.LastUpdate)
	if err != nil {
		return err
	}
	p.Results = make([]struct {
		OptionID, PrecinctID, TotalVotes int
	}, len(raw.PrecinctResults))
	for i := range raw.PrecinctResults {
		p.Results[i].OptionID = raw.PrecinctResults[i][0]
		p.Results[i].PrecinctID = raw.PrecinctResults[i][1]
		p.Results[i].TotalVotes = raw.PrecinctResults[i][2]
	}
	return err
}
