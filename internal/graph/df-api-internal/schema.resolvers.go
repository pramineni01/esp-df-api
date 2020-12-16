package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"bitbucket.org/antuitinc/esp-df-api/internal/datamodels"
	"bitbucket.org/antuitinc/esp-df-api/internal/graph/df-api-internal/generated"
)

func (r *mutationResolver) DfUpdateForecast(ctx context.Context, forecastID string, latestVersionDimensionMemberID string) (*bool, error) {
	response := true
	err := r.DBRepo.UpdateForecast(ctx, datamodels.UpdateForecastParams{
		ForecastID:                     forecastID,
		LatestVersionDimensionMemberID: latestVersionDimensionMemberID,
	})
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (r *mutationResolver) DfUpdateScenarioRun(ctx context.Context, scenarioRunID string, status datamodels.DFScenarioRunStatus, daVersion *string) (*bool, error) {
	response := true
	err := r.DBRepo.UpdateScenarioRun(ctx, datamodels.UpdateScenarioRunParams{
		ScenarioRunID: scenarioRunID,
		Status:        status,
		DaVersion:     daVersion,
	})
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
