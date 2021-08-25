package main

import (
	"fmt"
	"time"
)

type planResponse struct {
	Id string `json:"id"`
}

func getUrlForPlan(response planResponse, config *AppConfig) string {
	return fmt.Sprintf("%s/game/%s", config.Thunderdome.BaseUrl, response.Id)

}

type createBattleRequest struct {
	BattleName           string   `json:"battleName"`
	PointValuesAllowed   []string `json:"pointValuesAllowed"`
	AutoFinishVoting     bool     `json:"autoFinishVoting"`
	Plans                []*Plan  `json:"plans"`
	PointAverageRounding string   `json:"pointAverageRounding"`
}

type Vote struct {
	UserID    string `json:"warriorId"`
	VoteValue string `json:"vote"`
}

type Plan struct {
	PlanID             string    `json:"id"`
	PlanName           string    `json:"name"`
	Type               string    `json:"type"`
	ReferenceID        string    `json:"referenceId"`
	Link               string    `json:"link"`
	Description        string    `json:"description"`
	AcceptanceCriteria string    `json:"acceptanceCriteria"`
	Votes              []*Vote   `json:"votes"`
	Points             string    `json:"points"`
	PlanActive         bool      `json:"active"`
	PlanSkipped        bool      `json:"skipped"`
	VoteStartTime      time.Time `json:"voteStartTime"`
	VoteEndTime        time.Time `json:"voteEndTime"`
}
