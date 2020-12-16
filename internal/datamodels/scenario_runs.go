package datamodels

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"bitbucket.org/antuitinc/esp-df-api/pkg/esputils"
	"github.com/jmoiron/sqlx"
)

// Model
type DFScenarioRun struct {
	ID         string              `db:"scenario_run_id"`
	ScenarioID string              `db:"scenario_id"`
	UserID     *string             `db:"user_id"`
	Status     DFScenarioRunStatus `db:"scenario_run_status"`
	DaVersion  *string             `db:"da_version_id"`

	ElapsedTime       *float64     `db:"elapsed_time"`
	RunStartTimestamp sql.NullTime `db:"run_start_timestamp"`
	RunEndTimestamp   sql.NullTime `db:"run_end_timestamp"`
}

// Queries

type UpdateScenarioRunParams struct {
	ScenarioRunID string
	Status        DFScenarioRunStatus
	DaVersion     *string
}

func (q *Queries) UpdateScenarioRun(ctx context.Context, args UpdateScenarioRunParams) error {
	// Validate the scenarioRun
	scenarioRun, err := q.findScenarioRunById(ctx, args.ScenarioRunID)
	if err != nil {
		log.Printf("ERROR: datamodels/scenarios.go/UpdateScenarioRun: %s", err)
		return err
	}
	if scenarioRun == nil {
		log.Printf("ERROR: datamodels/scenarios.go/UpdateScenarioRun: ScenarioRunID=%s does not exist", args.ScenarioRunID)
		return esputils.ErrDBNoRows("UpdateScenarioRun")
	}

	switch args.Status {
	case DFScenarioRunStatusScheduled:
		err := fmt.Errorf("Set status to %s is not allowed", DFScenarioRunStatusScheduled)
		log.Printf("ERROR: datamodels/scenarios.go/UpdateScenarioRun: %s", err)
		return err

	case DFScenarioRunStatusInProgress:
		if scenarioRun.Status == DFScenarioRunStatusForecasted || scenarioRun.Status == DFScenarioRunStatusError {
			err := fmt.Errorf("Change status from %s to %s is not allowed", scenarioRun.Status, DFScenarioRunStatusInProgress)
			log.Printf("ERROR: datamodels/scenarios.go/UpdateScenarioRun: %s", err)
			return err
		}
		query := `UPDATE scenario_runs
					SET scenario_run_status = ?,
						run_start_timestamp = NOW()
					WHERE scenario_run_id = ?`

		_, err = q.db.ExecContext(ctx, query,
			args.Status,
			args.ScenarioRunID,
		)
		if err != nil {
			log.Printf("ERROR: datamodels/scenarios.go/UpdateScenarioRun: [ScenarioRunID=%s Status=%s] %s", args.ScenarioRunID, DFScenarioRunStatusInProgress, err)
			return err
		}

	case DFScenarioRunStatusForecasted:
		if scenarioRun.Status == DFScenarioRunStatusScheduled || scenarioRun.Status == DFScenarioRunStatusError {
			err := fmt.Errorf("Change status from %s to %s is not allowed", scenarioRun.Status, DFScenarioRunStatusForecasted)
			log.Printf("ERROR: datamodels/scenarios.go/UpdateScenarioRun: %s", err)
			return err
		}
		if args.DaVersion == nil {
			err := fmt.Errorf("daVersion param is required")
			log.Printf("ERROR: datamodels/scenarios.go/UpdateScenarioRun: %s", err)
			return err
		}

		query := `UPDATE scenario_runs
				  SET scenario_run_status = ?,
					  run_end_timestamp   = NOW(),
					  da_version_id       = ?
				  WHERE scenario_run_id   = ?`

		_, err = q.db.ExecContext(ctx, query,
			args.Status,
			*args.DaVersion,
			args.ScenarioRunID,
		)
		if err != nil {
			log.Printf("ERROR: datamodels/scenarios.go/UpdateScenarioRun: [ScenarioRunID=%s Status=%s] %s", args.ScenarioRunID, DFScenarioRunStatusForecasted, err)
			return err
		}

	case DFScenarioRunStatusError:
		if scenarioRun.Status == DFScenarioRunStatusScheduled || scenarioRun.Status == DFScenarioRunStatusForecasted {
			err := fmt.Errorf("Change status from %s to %s is not allowed", scenarioRun.Status, DFScenarioRunStatusError)
			log.Printf("ERROR: datamodels/scenarios.go/UpdateScenarioRun: %s", err)
			return err
		}

		query := `UPDATE scenario_runs
				  SET scenario_run_status = ?,
					  run_end_timestamp   = NOW()
				  WHERE scenario_run_id   = ?`

		_, err = q.db.ExecContext(ctx, query,
			DFScenarioRunStatusError,
			args.ScenarioRunID,
		)
		if err != nil {
			log.Printf("ERROR: datamodels/scenarios.go/UpdateScenarioRun: [ScenarioRunID=%s Status=%s] %s", args.ScenarioRunID, DFScenarioRunStatusError, err)
			return err
		}

	default:
		err := fmt.Errorf("Unknown status=%s", args.Status)
		log.Printf("ERROR: datamodels/scenarios.go/UpdateScenarioRun: %s", err)
		return err
	}

	return nil
}

const creatScenarioRun = `
INSERT INTO scenario_runs(scenario_id, user_id)
VALUES(?, ?)
RETURNING
scenario_run_id,
scenario_id,
scenario_run_status,
user_id,
da_version_id

`

func (q *Queries) InsertScenarioRun(ctx context.Context, scenarioID string, userID string) (*DFScenarioRun, error) {
	scenarioRun := DFScenarioRun{}
	err := q.db.QueryRowxContext(ctx, creatScenarioRun, scenarioID, userID).StructScan(&scenarioRun)
	if err != nil {

		log.Printf("ERROR: datamodels/scenarios.go/insertScenarioRun: %s", err)
		return &scenarioRun, err
	}

	return &scenarioRun, nil

}

const scenarioRunQuery = `
SELECT scenario_run_id,
       scenario_id,
       scenario_run_status,
       user_id,
       da_version_id
  FROM scenario_runs
 WHERE scenario_id = ?
 ORDER BY scenario_run_id DESC

`

func (q *Queries) FetchLatestRun(ctx context.Context, scenarioID string) (*DFScenarioRun, error) {
	scenarioRun := DFScenarioRun{}
	query := scenarioRunQuery + " LIMIT 1"

	if err := q.db.QueryRowxContext(ctx, query, scenarioID).StructScan(&scenarioRun); err != nil {
		log.Printf("ERROR: datamodels/scenarios.go/FetchLatestRun: %s", err)
		return nil, err
	}
	return &scenarioRun, nil
}

func (q *Queries) FetchAllRuns(ctx context.Context, scenarioID string) ([]*DFScenarioRun, error) {
	scenarioRuns := []*DFScenarioRun{}
	rows, err := q.db.QueryxContext(ctx, scenarioRunQuery, scenarioID)
	if err != nil {

		log.Printf("ERROR: datamodels/scenarios.go/FetchAllRuns: %s", err)
	}
	for rows.Next() {
		scenarioRun := DFScenarioRun{}
		if err := rows.StructScan(&scenarioRun); err != nil {

			log.Printf("ERROR: datamodels/scenarios.go/FetchAllRuns: %s", err)
			return nil, err
		}
		scenarioRuns = append(scenarioRuns, &scenarioRun)

	}
	return scenarioRuns, nil
}

// Private Methods ------------------------------------------------------------

func (q *Queries) findScenarioRunById(ctx context.Context, scenarioRunId string) (*DFScenarioRun, error) {
	row := q.db.QueryRowxContext(ctx, "SELECT * FROM scenario_runs	WHERE scenario_run_id = ?", scenarioRunId)
	return rowToScenarioRun(row)
}

func rowToScenarioRun(row *sqlx.Row) (*DFScenarioRun, error) {
	scenarioRun := DFScenarioRun{}
	err := row.StructScan(&scenarioRun)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &scenarioRun, nil
}
