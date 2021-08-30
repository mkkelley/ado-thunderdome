package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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

const thunderdomeApiKeyHeaderName = "X-API-Key"

func getTimeSuffix() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func adoItemToPlan(item adoItem) Plan {
	return Plan{
		PlanName:           item.title,
		Type:               "Story",
		ReferenceID:        "",
		Link:               item.link,
		Description:        item.description,
		AcceptanceCriteria: item.acceptanceCriteria,
	}
}

func generateBattle(config *AppConfig, apiKey string, queryId string, battlePrefix string) (planResponse, error) {
	adoItems, err := getReadyForGroomingItems(config, queryId)
	if err != nil {
		return planResponse{}, err
	}
	plans := make([]*Plan, len(adoItems))
	for i, item := range adoItems {
		plan := adoItemToPlan(item)
		plans[i] = &plan
	}

	request := createBattleRequest{
		BattleName: fmt.Sprintf("%s %s", battlePrefix, getTimeSuffix()),
		PointValuesAllowed: []string{
			".5",
			"1",
			"2",
			"3",
			"5",
			"8",
			"13",
			"20",
			"?",
		},
		AutoFinishVoting:     true,
		Plans:                plans,
		PointAverageRounding: "ceil",
	}
	requestJson, err := json.Marshal(request)
	if err != nil {
		return planResponse{}, err
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/battle", config.Thunderdome.BaseUrl), bytes.NewReader(requestJson))
	if err != nil {
		return planResponse{}, err
	}
	req.Header.Set(thunderdomeApiKeyHeaderName, apiKey)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return planResponse{}, err
	}
	resp := planResponse{}
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		return planResponse{}, err
	}
	return resp, nil
}
