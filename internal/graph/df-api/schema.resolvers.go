package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"log"

	"bitbucket.org/antuitinc/esp-df-api/internal/datamodels"
	generated1 "bitbucket.org/antuitinc/esp-df-api/internal/graph/df-api/generated"
	model1 "bitbucket.org/antuitinc/esp-df-api/internal/graph/df-api/model"
	"bitbucket.org/antuitinc/esp-df-api/pkg/esputils"
)

func (r *dFForecastResolver) Scenarios(ctx context.Context, obj *datamodels.DFForecast, scope []*model1.DFScenarioScopeEntryInput) ([]*datamodels.DFScenario, error) {
	return r.DBRepo.FindScenarios(ctx, obj.ID, scope)
}

func (r *dFForecastResolver) Scenario(ctx context.Context, obj *datamodels.DFForecast, id string) (*datamodels.DFScenario, error) {
	return r.DBRepo.FindScenario(ctx, obj.ID, id)
}

func (r *dFScenarioResolver) Comments(ctx context.Context, obj *datamodels.DFScenario) ([]*datamodels.DFScenarioComment, error) {
	return r.Dataloaders.Retrieve(ctx).CommentByCommentID.Load(obj.ID)
}

func (r *dFScenarioResolver) Tags(ctx context.Context, obj *datamodels.DFScenario) ([]*model1.Tag, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *dFScenarioResolver) User(ctx context.Context, obj *datamodels.DFScenario) (*model1.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *dFScenarioResolver) LatestRun(ctx context.Context, obj *datamodels.DFScenario) (*datamodels.DFScenarioRun, error) {
	return r.DBRepo.FetchLatestRun(ctx, obj.ID)
}

func (r *dFScenarioResolver) AllRuns(ctx context.Context, obj *datamodels.DFScenario) ([]*datamodels.DFScenarioRun, error) {
	return r.DBRepo.FetchAllRuns(ctx, obj.ID)
}

func (r *dFScenarioCommentResolver) User(ctx context.Context, obj *datamodels.DFScenarioComment) (*model1.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *dFScenarioRunResolver) User(ctx context.Context, obj *datamodels.DFScenarioRun) (*model1.User, error) {
	if obj.UserID != nil {
		return &model1.User{
			UserID: *obj.UserID,
		}, nil

	}
	return &model1.User{}, nil
}

func (r *mutationResolver) DfCreateScenario(ctx context.Context, forecastID string, scenarioName string, daBranchID string, tagIDs []int, scope []*model1.DFScenarioScopeEntryInput, comment *string) (*datamodels.DFScenarioRun, error) {
	userID := ctx.Value("userId").(string)

	// Create the Scenario
	scenario, err := r.DBRepo.CreateScenario(ctx, datamodels.CreateScenarioParams{
		ForecastID:   forecastID,
		ScenarioName: scenarioName,
		TagIDs:       tagIDs,
		Comment:      comment,
		DaBranchID:   daBranchID,
		UserID:       userID,
		Scope:        scope,
	})
	if err != nil {
		log.Printf("ERROR: schema.resolvers.go/DfCreateScenario: %s", err)
		return nil, err
	}
	return scenario, nil
}

func (r *mutationResolver) DfReRunScenario(ctx context.Context, scenarioID string) (*datamodels.DFScenarioRun, error) {
	userID := ctx.Value("userId").(string)
	return r.DBRepo.InsertScenarioRun(ctx, scenarioID, userID)
}

func (r *mutationResolver) DfPromoteScenario(ctx context.Context, scenarioID string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DfDeleteScenario(ctx context.Context, scenarioID string) (bool, error) {
	userID, err := esputils.GetUserIDFromContext(ctx)
	if err != nil {
		return false, err
	}

	err = r.DBRepo.DeleteScenarioById(ctx, datamodels.DeleteScenarioByIdParams{
		UserID:     userID,
		ScenarioID: scenarioID,
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DfTagScenario(ctx context.Context, scenarioID string, tagIds []int) (*datamodels.DFScenario, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DfAddComentToScenario(ctx context.Context, scenarioID string, comment string) (*datamodels.DFScenarioComment, error) {
	scenarioComment, err := r.DBRepo.CreateScenarioComment(ctx, datamodels.CreateScenarioCommentParams{
		ScenarioID: scenarioID,
		Comment:    comment,
	})
	if err != nil {
		log.Printf("ERROR: schema.resolvers.go/DfAddComentToScenario: %s", err)
		return nil, err
	}
	return scenarioComment, nil
}

func (r *mutationResolver) DfCreateDataFilter(ctx context.Context, name string, definition string) (*datamodels.DFDataFilter, error) {
	//Create the DataFilter
	userID := ctx.Value("userId").(string)
	dataFilter, err := r.DBRepo.CreateDataFilter(ctx, datamodels.CreateDataFilterParams{
		Name:       name,
		Definition: definition,
		UserID:     userID,
	})
	if err != nil {
		log.Printf("ERROR: schema.resolvers.go/DfCreateDataFilter: %s", err)
		return nil, err
	}
	return dataFilter, nil
}

func (r *mutationResolver) DfUpdateDataFilter(ctx context.Context, id string, name string, definition string) (*datamodels.DFDataFilter, error) {
	//Update the DataFilter
	userID := ctx.Value("userId").(string)
	dataFilter, err := r.DBRepo.UpdateDataFilter(ctx, datamodels.UpdateDataFilterParams{
		Id:         id,
		Name:       name,
		Definition: definition,
		UserID:     userID,
	})
	if err != nil {
		log.Printf("ERROR: schema.resolvers.go/DfUpdateDataFilter: %s", err)
		return nil, err
	}
	return dataFilter, nil
}

func (r *mutationResolver) DfDeleteDataFilter(ctx context.Context, id string) (bool, error) {
	//Delete the DataFilter
	userID := ctx.Value("userId").(string)
	err := r.DBRepo.DeleteDataFilter(ctx, datamodels.DeleteDataFilterParams{
		Id:     id,
		UserID: userID,
	})
	if err != nil {
		log.Printf("ERROR: schema.resolvers.go/DfDeleteDataFilter: %s", err)
		return false, err
	}
	return true, nil
}

func (r *queryResolver) DfForecasts(ctx context.Context) ([]*datamodels.DFForecast, error) {
	forecasts, err := r.DBRepo.FindForecasts(ctx)
	if err != nil {
		return nil, err
	}

	return forecasts, nil
}

func (r *queryResolver) DfForecast(ctx context.Context, id string) (*datamodels.DFForecast, error) {
	return r.DBRepo.FindForecastById(ctx, id)
}

func (r *queryResolver) DfForecastVersions(ctx context.Context, limit int) ([]*datamodels.DFForecastVersion, error) {
	forecast_versions, err := r.DBRepo.GetForecastVersions(ctx, datamodels.GetForecastVersionsParams{
		Limit: limit,
	})
	if err != nil {
		return nil, err
	}

	return forecast_versions, nil
}

func (r *queryResolver) DfForecastVersion(ctx context.Context, id string) (*datamodels.DFForecastVersion, error) {
	forecast_version, err := r.DBRepo.GetForecastVersion(ctx, id)
	if err != nil {
		return nil, err
	}

	return forecast_version, nil
}

func (r *queryResolver) DfDataFilters(ctx context.Context) ([]*datamodels.DFDataFilter, error) {
	//Get Data filters
	userID := ctx.Value("userId").(string)
	dataFilters, err := r.DBRepo.FindDataFilters(ctx, datamodels.FindDataFiltersParams{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	return dataFilters, nil
}

// DFForecast returns generated1.DFForecastResolver implementation.
func (r *Resolver) DFForecast() generated1.DFForecastResolver { return &dFForecastResolver{r} }

// DFScenario returns generated1.DFScenarioResolver implementation.
func (r *Resolver) DFScenario() generated1.DFScenarioResolver { return &dFScenarioResolver{r} }

// DFScenarioComment returns generated1.DFScenarioCommentResolver implementation.
func (r *Resolver) DFScenarioComment() generated1.DFScenarioCommentResolver {
	return &dFScenarioCommentResolver{r}
}

// DFScenarioRun returns generated1.DFScenarioRunResolver implementation.
func (r *Resolver) DFScenarioRun() generated1.DFScenarioRunResolver { return &dFScenarioRunResolver{r} }

// Mutation returns generated1.MutationResolver implementation.
func (r *Resolver) Mutation() generated1.MutationResolver { return &mutationResolver{r} }

// Query returns generated1.QueryResolver implementation.
func (r *Resolver) Query() generated1.QueryResolver { return &queryResolver{r} }

type dFForecastResolver struct{ *Resolver }
type dFScenarioResolver struct{ *Resolver }
type dFScenarioCommentResolver struct{ *Resolver }
type dFScenarioRunResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
