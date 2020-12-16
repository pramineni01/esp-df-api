package datamodels

import (
	"context"
	"database/sql"
	"log"
	"strconv"

	"bitbucket.org/antuitinc/esp-df-api/pkg/esputils"
	"github.com/jmoiron/sqlx"
)

// Model

type DFForecast struct {
	ID                             string  `db:"forecast_id"`
	DatasetID                      string  `db:"dataset_id"`
	Name                           string  `db:"forecast_name"`
	LatestVersionDimensionMemberID *string `db:"latest_version_dimension_member_id"`
}

func (DFForecast) IsEntity() {}

// Queries
func rowToForecast(row *sqlx.Row) (*DFForecast, error) {
	forecast := DFForecast{}
	err := row.StructScan(&forecast)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &forecast, nil
}

//Create Forecast
const createForecast = `
INSERT INTO forecasts (
	dataset_id,
	name,
	latest_version_dimension_member_id
) VALUES (
	?, ?
) RETURNING *
`

type CreateForecastParams struct {
	DatasetID                      string
	Name                           string
	LatestVersionDimensionMemberID sql.NullInt32
}

func (q *Queries) CreateForecast(ctx context.Context, args CreateForecastParams) (*DFForecast, error) {
	forecast := DFForecast{}
	row := q.db.QueryRowxContext(ctx, createForecast,
		args.DatasetID,
		args.Name,
		args.LatestVersionDimensionMemberID,
	)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	err = row.StructScan(&forecast)
	return &forecast, err
}

type UpdateForecastParams struct {
	ForecastID                     string
	LatestVersionDimensionMemberID string
}

func (q *Queries) UpdateForecast(ctx context.Context, args UpdateForecastParams) error {
	// Validate the forecast
	forecast, err := q.FindForecastById(ctx, args.ForecastID)
	if err != nil {
		log.Printf("ERROR: datamodels/forecasts.go/UpdateForecast: %s", err)
		return err
	}
	if forecast == nil {
		log.Printf("ERROR: datamodels/forecasts.go/UpdateForecast: ForecastID=%s does not exist", args.ForecastID)
		return esputils.ErrDBNoRows("Update Forecast")
	}

	query := `UPDATE forecasts
			  SET latest_version_dimension_member_id = ?
			  WHERE forecast_id = ?`

	_, err = q.db.ExecContext(ctx, query,
		args.LatestVersionDimensionMemberID,
		args.ForecastID,
	)
	if err != nil {
		log.Printf("ERROR: datamodels/forecasts.go/UpdateForecast: %s", err)
		return err
	}

	return nil
}

// Find Forecast By ID
const findForecastById = `
SELECT f.forecast_id,
       f.dataset_id,
       f.latest_version_dimension_member_id,
       ft.forecast_name
  FROM forecasts f
         INNER JOIN forecast_translations ft
             ON f.forecast_id = ft.forecast_id
 WHERE f.forecast_id = ? AND ft.locale_id = ?
`

type FindForecastByIdParams struct {
	ForecastID string
}

func (q *Queries) FindForecastById(ctx context.Context, forecastID string) (*DFForecast, error) {
	var localeId int
	var err error
	if ctx.Value("userLocale").(string) == "" {
		localeId = esputils.DEFAULT_LOCALE_ID
	} else {
		localeId, err = strconv.Atoi(ctx.Value("userLocale").(string))
		if err != nil {
			return nil, err
		}
	}

	row := q.db.QueryRowxContext(ctx, findForecastById,
		forecastID,
		localeId,
	)
	return rowToForecast(row)
}

// Find Forecasts
const findForecasts = `
SELECT f.*, ft.forecast_name FROM forecasts f
INNER JOIN forecast_translations ft ON f.forecast_id = ft.forecast_id
WHERE ft.locale_id = ?`

func (q *Queries) FindForecasts(ctx context.Context) ([]*DFForecast, error) {
	var localeId int
	var err error
	if ctx.Value("userLocale").(string) == "" {
		localeId = esputils.DEFAULT_LOCALE_ID
	} else {
		localeId, err = strconv.Atoi(ctx.Value("userLocale").(string))
		if err != nil {
			return nil, err
		}
	}

	var listForecast []*DFForecast
	rows, err := q.db.QueryxContext(ctx, findForecasts, localeId)
	if err != nil {
		return nil, err
	}
	log.Printf(findForecasts)
	for rows.Next() {
		forecast := DFForecast{}

		err := rows.StructScan(&forecast)
		if err != nil {
			return nil, err
		}
		listForecast = append(listForecast, &forecast)
	}
	return listForecast, nil
}

// Find Forecasts by IDs
const findForecastsByIds = `
SELECT * FROM forecasts	WHERE forecast_id IN (?)
`

type FindForecastsByIdsParams struct {
	ForecastIDs []string
}

func (q *Queries) FindForecastsByIds(ctx context.Context, args FindForecastsByIdsParams) ([]*DFForecast, error) {
	_query := findForecastsByIds
	_in_args := []interface{}{args.ForecastIDs}

	var forecasts []*DFForecast
	query, inargs, err := sqlx.In(_query, _in_args...)

	if err != nil {
		return []*DFForecast{}, err
	}
	query = q.db.Rebind(query)
	rows, err := q.db.QueryxContext(ctx, query, inargs...)

	log.Printf(findForecastsByIds)

	if err != nil {
		return []*DFForecast{}, err
	}

	for rows.Next() {
		forecast := DFForecast{}
		err := rows.StructScan(&forecast)
		if err != nil {
			return []*DFForecast{}, err
		}
		forecasts = append(forecasts, &forecast)
	}
	return forecasts, nil
}
