package dataloaders

// go:generate go run github.com/vektah/dataloaden DatasetLoader int *bitbucket.org/antuitinc/esp-df-api/datamodels.Dataset

import (
	"context"
	"time"

	"bitbucket.org/antuitinc/esp-df-api/internal/datamodels"
)

type contextKey string

const key = contextKey("dataloaders")

// Loaders holds references to the individual dataloaders.
type Loaders struct {
	ScenarioByScenarioID *ScenarioLoader
	ForecastByForecastID *ForecastLoader
	CommentByCommentID   *CommentLoader
}

func newLoaders(ctx context.Context, dbrepo datamodels.DBRepo) *Loaders {
	return &Loaders{
		// individual loaders will be initialized here
		ScenarioByScenarioID: newScenarioById(ctx, dbrepo),
		ForecastByForecastID: newForecastByID(ctx, dbrepo),
		CommentByCommentID:   NewCommentByID(ctx, dbrepo),
	}
}

// Retriever retrieves dataloaders from the request context.
type Retriever interface {
	Retrieve(context.Context) *Loaders
}

type retriever struct {
	key contextKey
}

func (r *retriever) Retrieve(ctx context.Context) *Loaders {
	return ctx.Value(r.key).(*Loaders)
}

// NewRetriever instantiates a new implementation of Retriever.
func NewRetriever() Retriever {
	return &retriever{key: key}
}

func newScenarioById(ctx context.Context, dbrepo datamodels.DBRepo) *ScenarioLoader {
	return NewScenarioLoader(ScenarioLoaderConfig{
		MaxBatch: 100,
		Wait:     5 * time.Millisecond,
		Fetch: func(scenarioIDs []string) ([]*datamodels.DFScenario, []error) {
			// db query
			res, err := dbrepo.FindScenariosByIds(ctx, datamodels.FindScenariosByIdsParams{
				ScenarioIDs: scenarioIDs,
			})
			if err != nil {
				return nil, []error{err}
			}
			return res, nil
		},
	})
}

func newForecastByID(ctx context.Context, dbrepo datamodels.DBRepo) *ForecastLoader {
	return NewForecastLoader(ForecastLoaderConfig{
		MaxBatch: 100,
		Wait:     5 * time.Millisecond,
		Fetch: func(forecastIDs []string) ([]*datamodels.DFForecast, []error) {
			// db query
			res, err := dbrepo.FindForecastsByIds(ctx, datamodels.FindForecastsByIdsParams{
				ForecastIDs: forecastIDs,
			})
			if err != nil {
				return nil, []error{err}
			}
			return res, nil
		},
	})
}

func NewCommentByID(ctx context.Context, dbrepo datamodels.DBRepo) *CommentLoader {
	return NewCommentLoader(CommentLoaderConfig{
		MaxBatch: 100,
		Wait:     5 * time.Millisecond,
		Fetch: func(scenarioIDs []string) ([][]*datamodels.DFScenarioComment, []error) {

			// db query
			res, err := dbrepo.FindScenarioCommentsByScenarioIds(ctx, datamodels.FindScenarioCommentsByScenarioIdsParams{
				ScenarioIDs: scenarioIDs,
			})
			if err != nil {
				return nil, []error{err}
			}

			// group
			groupByScenarioID := make(map[string][]*datamodels.DFScenarioComment, len(scenarioIDs))
			for _, r := range res {
				groupByScenarioID[r.ScenarioID] = append(groupByScenarioID[r.ScenarioID], r)
			}

			// order
			result := make([][]*datamodels.DFScenarioComment, len(scenarioIDs))
			for i, scenarioID := range scenarioIDs {
				result[i] = groupByScenarioID[scenarioID]
			}

			return result, nil
		},
	})
}
