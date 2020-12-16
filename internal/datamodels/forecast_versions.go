package datamodels

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Model

type DFForecastVersion struct {
	ID           string `db:"forecast_version_id"`
	ForecastName string `db:"forecast_name"`
	DatasetID    string `db:"dataset_id"`
}

func (DFForecastVersion) IsEntity() {}

// Queries
func rowToForecastDataset(row *sqlx.Row) (*DFForecastVersion, error) {
	forecastV := DFForecastVersion{}
	err := row.StructScan(&forecastV)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &forecastV, nil
}

type GetForecastVersionsParams struct {
	Limit int
}

func (q *Queries) GetForecastVersions(ctx context.Context, args GetForecastVersionsParams) ([]*DFForecastVersion, error) {
	query := fmt.Sprintf(`
		SELECT
			forecast_version_id,
			forecast_name,
			dataset_id
		FROM forecast_versions
		ORDER BY forecast_name DESC
		LIMIT %d;`,
		args.Limit,
	)

	rows, err := q.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var listForecastV []*DFForecastVersion
	for rows.Next() {
		forecastV := DFForecastVersion{}
		err := rows.StructScan(&forecastV)
		if err != nil {
			return nil, err
		}
		listForecastV = append(listForecastV, &forecastV)
	}

	return listForecastV, nil
}

func (q *Queries) GetForecastVersion(ctx context.Context, id string) (*DFForecastVersion, error) {
	query := `
	SELECT
		forecast_version_id,
		forecast_name,
		dataset_id
	FROM forecast_versions
	WHERE forecast_version_id = ?
	LIMIT 1`

	row := q.db.QueryRowxContext(ctx, query, id)
	return rowToForecastDataset(row)
}
