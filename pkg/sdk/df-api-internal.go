package sdk

import (
	"context"
	"fmt"
	"log"

	"github.com/jaylane/graphql"
)

// External structures --------------------------

// Internal structures --------------------------

type Client struct {
	endpoint  string
	ctx       context.Context
	clientGQL *graphql.Client
}

// Logic ----------------------------------------

func NewClient(ctx context.Context, endpoint string) *Client {
	// Create GraphQL client
	clientGQL := graphql.NewClient(endpoint)
	client := &Client{
		endpoint:  endpoint,
		ctx:       ctx,
		clientGQL: clientGQL,
	}

	return client
}

func (c *Client) UpdateForecast(forecastId string, latestVersionDimensionMemberId string) (*bool, error) {

	// Fill parameters of request string
	requestString := fmt.Sprintf(`
		mutation UpdateForecast {
			dfUpdateForecast (forecastId: %s, latestVersionDimensionMemberId: %s)
		}
	`, forecastId, latestVersionDimensionMemberId)

	// Define expected response format
	type updateForecastResponse struct {
		UpdateForecast bool `json:"dfUpdateForecast"`
	}

	var responseData updateForecastResponse

	// Make the request
	req := graphql.NewRequest(requestString)
	if err := c.clientGQL.Run(c.ctx, req, &responseData); err != nil {
		log.Printf("ERROR sdk/df-api-internal.go/UpdateForecast: %s [c.clientGQL.Run]", err)
		return nil, err
	}

	// Handle the response
	log.Printf("Response (UpdateForecast): %t", responseData.UpdateForecast)

	return &responseData.UpdateForecast, nil
}

func (c *Client) UpdateScenarioRun(scenarioRunID string, status string, daVersion *string) (*bool, error) {

	// Fill parameters of request string
	daVersionRequest := ""
	if daVersion != nil && len(*daVersion) > 0 {
		daVersionRequest = fmt.Sprintf(", daVersion: %s", *daVersion)
	}
	requestString := fmt.Sprintf(`
		mutation UpdateScenarioRun {
			dfUpdateScenarioRun (scenarioRunID: %s, status: %s%s)
		}
	`, scenarioRunID, status, daVersionRequest)

	// Define expected response format
	type updateScenarioRunResponse struct {
		UpdateScenarioRun bool `json:"dfUpdateScenarioRun"`
	}

	var responseData updateScenarioRunResponse

	// Make the request
	req := graphql.NewRequest(requestString)
	if err := c.clientGQL.Run(c.ctx, req, &responseData); err != nil {
		log.Printf("ERROR sdk/df-api-internal.go/UpdateScenarioRun: %s [c.clientGQL.Run]", err)
		return nil, err
	}

	// Handle the response
	log.Printf("Response (UpdateScenarioRun): %t", responseData.UpdateScenarioRun)

	return &responseData.UpdateScenarioRun, nil
}
