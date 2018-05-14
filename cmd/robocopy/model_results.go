package main

import (
	"fmt"
	"time"
)

type ResultsContainer struct {
	Test       bool
	LastUpdate time.Time
	Results    []struct {
		OptionID
		JurisdictionID
		DistrictID
		TotalVotes int
	}
	DistrictReporting []struct {
		JurisdictionID
		DistrictID
		PrecinctsReporting, TotalPrecincts int
		PercentCounted                     float64
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

func ResultsContainerFrom(name string) (*ResultsContainer, error) {
	var r ResultsContainer
	if err := unmarshalFrom(name, &r); err != nil {
		return nil, fmt.Errorf("could not read results data: %v", err)
	}
	return &r, nil
}
