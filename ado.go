package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type adoItem struct {
	title              string
	link               string
	description        string
	acceptanceCriteria string
}

type adoQueryResponse struct {
	WorkItems []struct {
		Id int64 `json:"id"`
	} `json:"workItems"`
}

type adoWorkItemsResponse struct {
	Count int64 `json:"count"`
	Value []struct {
		Id     int64 `json:"id"`
		Fields struct {
			Description        string `json:"System.Description"`
			AcceptanceCriteria string `json:"Microsoft.VSTS.Common.AcceptanceCriteria"`
			Title              string `json:"System.Title"`
		}
	} `json:"value"`
}

func getUserLinkForWorkItem(id int64, config *AppConfig) string {
	return fmt.Sprintf("%s/%s/%s/_workItems/edit/%d", config.Ado.BaseUrl, config.Ado.Organization, config.Ado.Project, id)
}

func getAdoHttpRequest(method string, url string, config *AppConfig) (*http.Request, error) {
	// Username is not relevant in the basic auth, the user is inferred from the PAT
	basicAuthString := fmt.Sprintf("username:%s", config.Ado.PersonalAccessToken)
	encodedAuthString := []byte(base64.StdEncoding.EncodeToString([]byte(basicAuthString)))
	authHeader := fmt.Sprintf("Basic %s", encodedAuthString)
	req, err := http.NewRequest(method, url, bytes.NewReader(encodedAuthString))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", authHeader)

	return req, nil
}

func getWorkItemsById(ids []int64, config *AppConfig) ([]adoItem, error) {
	idStrings := make([]string, len(ids))
	for i, id := range ids {
		idStrings[i] = strconv.FormatInt(id, 10)
	}

	idQueryString := strings.Join(idStrings, ",")
	url := fmt.Sprintf("%s/%s/_apis/wit/workitems?ids=%s&apiiversion=6.0", config.Ado.BaseUrl, config.Ado.Organization, idQueryString)
	request, err := getAdoHttpRequest(http.MethodGet, url, config)
	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	workItemsResponse := adoWorkItemsResponse{}
	err = json.NewDecoder(response.Body).Decode(&workItemsResponse)
	if err != nil {
		return nil, err
	}

	adoItems := make([]adoItem, workItemsResponse.Count)
	for i, workItem := range workItemsResponse.Value {
		adoItems[i] = adoItem{
			title:              workItem.Fields.Title,
			link:               getUserLinkForWorkItem(workItem.Id, config),
			description:        workItem.Fields.Description,
			acceptanceCriteria: workItem.Fields.AcceptanceCriteria,
		}
	}
	return adoItems, nil
}

func getReadyForGroomingItems(config *AppConfig, queryId string) ([]adoItem, error) {
	url := fmt.Sprintf("%s/%s/%s/_apis/wit/wiql/%s?api-version=6.0", config.Ado.BaseUrl, config.Ado.Organization,
		config.Ado.Project, queryId)
	req, err := getAdoHttpRequest(http.MethodGet, url, config)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	queryResponse := adoQueryResponse{}
	err = json.NewDecoder(resp.Body).Decode(&queryResponse)
	if err != nil {
		return nil, err
	}

	ids := make([]int64, len(queryResponse.WorkItems))
	for i, workItem := range queryResponse.WorkItems {
		ids[i] = workItem.Id
	}

	return getWorkItemsById(ids, config)
}
