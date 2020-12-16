package datamodels

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
)

// Model
type DFScenarioComment struct {
	ID         string `db:"scenario_comment_id"`
	ScenarioID string `db:"scenario_id"`
	Comment    string `db:"comment"`
}

func (DFScenarioComment) IsEntity() {}

// Queries
func rowToScenarioComment(row *sqlx.Row) (*DFScenarioComment, error) {
	scenarioComment := DFScenarioComment{}
	err := row.Err()
	if err != nil {
		return nil, err
	}
	row.StructScan(&scenarioComment)
	return &scenarioComment, nil
}

// Create Scenario Comment
const createScenarioComment = `
INSERT INTO scenario_comments (
	scenario_id,
	comment
) VALUES (
	?, ?
) RETURNING *
`

type CreateScenarioCommentParams struct {
	ScenarioID string
	Comment    string
}

func (q *Queries) CreateScenarioComment(ctx context.Context, args CreateScenarioCommentParams) (*DFScenarioComment, error) {
	scenarioComment := DFScenarioComment{}
	row := q.db.QueryRowxContext(ctx, createScenarioComment,
		args.ScenarioID,
		args.Comment,
	)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	err = row.StructScan(&scenarioComment)
	return &scenarioComment, err
}

// Find ScenarioComment By ID
const findScenarioCommentById = `
SELECT * FROM scenario_comments
	WHERE scenario_comment_id = ?
`

type FindScenarioCommentByIdParams struct {
	ScenarioCommentID string
}

func (q *Queries) FindScenarioCommentById(ctx context.Context, args FindScenarioCommentByIdParams) (*DFScenarioComment, error) {
	row := q.db.QueryRowxContext(ctx, findScenarioCommentById,
		args.ScenarioCommentID,
	)
	return rowToScenarioComment(row)
}

// Find ScenarioComments
const findScenarioComments = `
SELECT * FROM scenario_comments `

func (q *Queries) FindScenarioComments(ctx context.Context) ([]*DFScenarioComment, error) {
	var sc []*DFScenarioComment
	rows, err := q.db.QueryxContext(ctx, findScenarioComments)
	if err != nil {
		return nil, err
	}
	log.Printf(findScenarioComments)
	for rows.Next() {
		scenariocomment := DFScenarioComment{}

		err := rows.StructScan(&scenariocomment)
		if err != nil {
			return nil, err
		}
		sc = append(sc, &scenariocomment)
	}
	return sc, nil
}

// Find ScenarioComments by ScenarioIDs
const findScenarioCommentsByScenarioIds = `
SELECT * FROM scenario_comments	WHERE scenario_id IN (?)
`

type FindScenarioCommentsByScenarioIdsParams struct {
	ScenarioIDs []string
}

func (q *Queries) FindScenarioCommentsByScenarioIds(ctx context.Context, args FindScenarioCommentsByScenarioIdsParams) ([]*DFScenarioComment, error) {
	_query := findScenarioCommentsByScenarioIds
	_in_args := []interface{}{args.ScenarioIDs}

	var scenariocomments []*DFScenarioComment
	query, inargs, err := sqlx.In(_query, _in_args...)

	if err != nil {
		return []*DFScenarioComment{}, err
	}
	query = q.db.Rebind(query)
	rows, err := q.db.QueryxContext(ctx, query, inargs...)

	log.Printf(findScenarioCommentsByScenarioIds)

	if err != nil {
		return []*DFScenarioComment{}, err
	}

	for rows.Next() {
		scenariocomment := DFScenarioComment{}
		err := rows.StructScan(&scenariocomment)
		if err != nil {
			return []*DFScenarioComment{}, err
		}
		scenariocomments = append(scenariocomments, &scenariocomment)
	}
	return scenariocomments, nil
}
