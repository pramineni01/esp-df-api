package datamodels

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

// Model
type DFDataFilter struct {
	ID         string `db:"data_filter_id"`
	Name       string `db:"data_filter_name"`
	Definition string `db:"data_filter_definition"`
	UserID     string `db:"user_id"`
}

func (DFDataFilter) IsEntity() {}

// Queries
func rowToDataFilter(row *sqlx.Row) (*DFDataFilter, error) {
	dataFilter := DFDataFilter{}
	err := row.Err()
	if err != nil {
		return nil, err
	}
	row.StructScan(&dataFilter)
	return &dataFilter, nil
}

//Create DataFilter
const createDataFilter = `
INSERT INTO data_filters (
	data_filter_name,
	data_filter_definition,
	user_id
) VALUES (
	?, ?, ?
) RETURNING *
`

type CreateDataFilterParams struct {
	Name       string
	Definition string
	UserID     string
}

func (q *Queries) CreateDataFilter(ctx context.Context, args CreateDataFilterParams) (*DFDataFilter, error) {

	dataFilter := DFDataFilter{}
	row := q.db.QueryRowxContext(ctx, createDataFilter,
		args.Name,
		args.Definition,
		args.UserID,
	)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	log.Printf(createDataFilter)
	err = row.StructScan(&dataFilter)
	return &dataFilter, err
}

//Update DataFilter
const updateDataFilter = `
UPDATE data_filters
SET data_filter_name = ?,
	data_filter_definition = ?
WHERE user_id = ? AND data_filter_id = ?
`

type UpdateDataFilterParams struct {
	Id         string
	Name       string
	Definition string
	UserID     string
}

func (q *Queries) UpdateDataFilter(ctx context.Context, args UpdateDataFilterParams) (*DFDataFilter, error) {

	dataFilter := DFDataFilter{ID: args.Id, Name: args.Name, Definition: args.Definition, UserID: args.UserID}
	result, err := q.db.ExecContext(ctx, updateDataFilter,
		args.Name,
		args.Definition,
		args.UserID,
		args.Id,
	)
	log.Printf(updateDataFilter)
	if err != nil {
		return nil, err
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affectedRows == 0 {
		return nil, errors.New(fmt.Sprintf("The data filter id %s is not found or does not belong to the user.", args.Id))
	}
	return &dataFilter, nil
}

//Delete DataFilter
const deleteDataFilter = `
DELETE FROM data_filters
WHERE user_id = ? AND data_filter_id = ?
`

type DeleteDataFilterParams struct {
	Id     string
	UserID string
}

func (q *Queries) DeleteDataFilter(ctx context.Context, args DeleteDataFilterParams) error {

	result, err := q.db.ExecContext(ctx, deleteDataFilter,
		args.UserID,
		args.Id,
	)
	log.Printf(deleteDataFilter)
	if err != nil {
		return err
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows == 0 {
		return errors.New(fmt.Sprintf("The data filter id %s is not found or does not belong to the user.", args.Id))
	}
	return nil
}

// Find DataFilter By ID
const findDataFilterById = `
SELECT * FROM data_filters
	WHERE data_filter_id = ?
`

type FindDataFilterByIdParams struct {
	DataFilterID string
}

func (q *Queries) FindDataFilterById(ctx context.Context, args FindDataFilterByIdParams) (*DFDataFilter, error) {
	row := q.db.QueryRowxContext(ctx, findDataFilterById,
		args.DataFilterID,
	)
	return rowToDataFilter(row)
}

// Find DataFilters
const findDataFilters = `
SELECT * FROM data_filters WHERE user_id = ? ORDER BY data_filter_name ASC `

type FindDataFiltersParams struct {
	UserID string
}

func (q *Queries) FindDataFilters(ctx context.Context, args FindDataFiltersParams) ([]*DFDataFilter, error) {

	var dataFilters []*DFDataFilter
	rows, err := q.db.QueryxContext(ctx, findDataFilters, args.UserID)
	if err != nil {
		return nil, err
	}
	log.Printf(findDataFilters)
	for rows.Next() {
		dataFilter := DFDataFilter{}

		err := rows.StructScan(&dataFilter)
		if err != nil {
			return nil, err
		}
		dataFilters = append(dataFilters, &dataFilter)
	}
	return dataFilters, nil
}
