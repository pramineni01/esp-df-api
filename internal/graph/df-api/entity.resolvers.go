package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"bitbucket.org/antuitinc/esp-df-api/internal/datamodels"
	generated1 "bitbucket.org/antuitinc/esp-df-api/internal/graph/df-api/generated"
)

func (r *entityResolver) FindDFDataFilterByID(ctx context.Context, id string) (*datamodels.DFDataFilter, error) {
	return r.DBRepo.FindDataFilterById(ctx, datamodels.FindDataFilterByIdParams{
		DataFilterID: id,
	})
}

func (r *entityResolver) FindDFForecastByID(ctx context.Context, id string) (*datamodels.DFForecast, error) {
	return r.DBRepo.FindForecastById(ctx, id)
}

func (r *entityResolver) FindDFScenarioByID(ctx context.Context, id string) (*datamodels.DFScenario, error) {
	return r.DBRepo.FindScenarioById(ctx, datamodels.FindScenarioByIdParams{
		ScenarioID: id,
	})
}

// Entity returns generated1.EntityResolver implementation.
func (r *Resolver) Entity() generated1.EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }
