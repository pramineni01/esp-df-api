package datamodels

import (
	model1 "bitbucket.org/antuitinc/esp-df-api/internal/graph/df-api/model"
	"bitbucket.org/antuitinc/esp-df-api/pkg/esputils"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"strconv"
	"strings"
	"time"
)

// Model
type DFScenario struct {
	ID         string           `db:"scenario_id"`
	ForecastID string           `db:"forecast_id"`
	UserID     string           `db:"user_id"`
	Name       string           `db:"scenario_name"`
	ScopeID    string           `db:"scope_id"`
	DaBranchID string           `db:"da_branch_id"`
	Status     DFScenarioStatus `db:"scenario_status"`
	IsBlocked  bool             `db:"is_blocked"`
}

func (s *DFScenario) Scope() []*model1.DFScenarioScopeEntry {
	scopes := strings.Split(s.ScopeID, "&")
	scopeModelList := []*model1.DFScenarioScopeEntry{}
	for _, scope := range scopes {
		scopeVal := strings.Split(scope, "=")
		scopeModelList = append(scopeModelList, &model1.DFScenarioScopeEntry{
			DaDimLevelColumnName: scopeVal[0],
			DaDimMemberID:        scopeVal[1],
		})
	}
	return scopeModelList
}

func (DFScenario) IsEntity() {}

// Queries
func rowToScenario(row *sqlx.Row) (*DFScenario, error) {
	scenario := DFScenario{}
	err := row.Err()
	if err != nil {
		return nil, err
	}
	row.StructScan(&scenario)
	return &scenario, nil
}

//Create Scenario
const createScenario = `
INSERT INTO scenarios (
	forecast_id,
	scenario_name,
	da_branch_id,
	user_id,
        scope_id
) VALUES (
	?, ?, ?, ?, ?
)
`

type CreateScenarioParams struct {
	ForecastID   string
	ScenarioName string
	TagIDs       []int
	Comment      *string
	DaBranchID   string
	UserID       string
	Scope        []*model1.DFScenarioScopeEntryInput
}

func getScopeID(scope []*model1.DFScenarioScopeEntryInput) string {
	scopes_arr := []string{}
	for _, entry := range scope {
		scopes_arr = append(scopes_arr, fmt.Sprintf("%s=%s", entry.DaDimLevelColumnName, entry.DaDimMemberID))

	}
	return strings.Join(scopes_arr, "&")

}

func (q *Queries) CreateScenario(ctx context.Context, args CreateScenarioParams) (*DFScenarioRun, error) {
	// TODO: put this in a transaction
	result, err := q.db.ExecContext(ctx, createScenario,
		args.ForecastID,
		args.ScenarioName,
		args.DaBranchID,
		args.UserID,
		getScopeID(args.Scope),
	)
	if err != nil {
		log.Printf("ERROR: datamodels/scenarios.go/CreateScenario: %s", err)
		return nil, err
	}

	scenarioID, err := result.LastInsertId()
	if err != nil {
		log.Printf("ERROR: datamodels/branches.go/CreateScenario: %s", err)
		return nil, err
	}

	scenarioRun, err := q.InsertScenarioRun(ctx, strconv.FormatInt(scenarioID, 10), args.UserID)
	if err != nil {
		return nil, err
	}
	if len(args.TagIDs) > 0 {
		q.insertScenarioTags(ctx, scenarioID, &args.TagIDs)
	}

	if args.Comment != nil && len(*args.Comment) > 0 {
		q.insertScenarioComment(ctx, scenarioID, *args.Comment)
	}

	return scenarioRun, nil
}

func (q *Queries) insertScenarioTags(ctx context.Context, scenarioID int64, tagIDs *[]int) error {
	tags_query := `INSERT INTO scenario_tags VALUES `
	vals := []interface{}{}

	for _, tag_id := range *tagIDs {
		tags_query += "(?, ?),"
		vals = append(vals, scenarioID, tag_id)
	}

	tags_query = tags_query[0 : len(tags_query)-1]

	_, err := q.db.ExecContext(ctx, tags_query, vals...)
	if err != nil {
		log.Printf("ERROR: datamodels/scenarios.go/insertScenarioTags: %s", err)
		return err
	}

	return nil
}

func (q *Queries) insertScenarioComment(ctx context.Context, scenarioID int64, comment string) error {
	tags_query := `INSERT INTO scenario_comments (scenario_id, comment) 
		               VALUES (?, ?)`

	_, err := q.db.ExecContext(ctx, tags_query, scenarioID, comment)
	if err != nil {
		log.Printf("ERROR: datamodels/scenarios.go/insertScenarioComment: %s", err)
		return err
	}

	return nil
}

// Delete Scenario By ID
const deleteScenarioById = `
DELETE FROM scenarios
	WHERE scenario_id = ? AND user_id = ?
`

type DeleteScenarioByIdParams struct {
	ScenarioID string
	UserID     string
}

func (q *Queries) DeleteScenarioById(ctx context.Context, args DeleteScenarioByIdParams) error {
	res, err := q.db.ExecContext(ctx, deleteScenarioById,
		args.ScenarioID,
		args.UserID,
	)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows < 1 {
		return esputils.ErrDBNoRows("delete scenario")
	}
	return nil
}

// Find Scenario By ID
const findScenario = `
SELECT scenario_id,
       forecast_id,
       user_id,
       scenario_name,
       scope_id,
       da_branch_id,
       scenario_status,
       (SELECT CASE scenario_run_status
               WHEN "FORECASTED"
                 THEN FALSE
               WHEN "ERROR"
                 THEN FALSE
               ELSE TRUE
               END AS is_blocked
          FROM scenario_runs
         WHERE scenario_id = scenario_id
         ORDER BY scenario_run_id DESC LIMIT 1
       ) AS is_blocked

  FROM scenarios


`

type FindScenarioByIdParams struct {
	ScenarioID string
}

func (q *Queries) FindScenario(ctx context.Context, forecastID string, scenarioID string) (*DFScenario, error) {
	scenario := DFScenario{}
	query := findScenario + " WHERE scenario_id = ? AND forecast_id = ?"
	err := q.db.QueryRowxContext(ctx, query, scenarioID, forecastID).StructScan(&scenario)
	if err != nil {
		log.Printf("ERROR: datamodels/scenarios.go/FindScenario: %s", err)
		return nil, err

	}
	return &scenario, nil

}

func (q *Queries) FindScenarios(ctx context.Context, forecastID string, scope []*model1.DFScenarioScopeEntryInput) ([]*DFScenario, error) {
	clauses := []string{"forecast_id = ?"}
	values := []interface{}{forecastID}
	for _, scope := range scope {
		clauses = append(clauses, `scope_id LIKE ?`)
		values = append(values, fmt.Sprintf(`%%%s=%s%%`, scope.DaDimLevelColumnName, scope.DaDimMemberID))

	}
	query := findScenario + "WHERE " + strings.Join(clauses, " AND ")
	log.Printf(query)
	fmt.Println(values)
	rows, err := q.db.QueryxContext(ctx, query, values...)
	if err != nil {
		log.Printf("ERROR: datamodels/scenarios.go/FindScenarios: %s", err)
		return nil, err
	}
	scenarios := []*DFScenario{}
	for rows.Next() {
		scenario := DFScenario{}
		if err := rows.StructScan(&scenario); err != nil {
			log.Printf("ERROR: datamodels/scenarios.go/FindScenarios: %s", err)
			return nil, err

		}

		scenarios = append(scenarios, &scenario)

	}
	return scenarios, nil
}

func (q *Queries) FindScenarioById(ctx context.Context, arg FindScenarioByIdParams) (*DFScenario, error) {
	query := findScenario + "WHERE scenario_id = ?"
	scenario := DFScenario{}
	err := q.db.QueryRowxContext(ctx, query, arg.ScenarioID).StructScan(&scenario)
	if err != nil {
		log.Printf("ERROR: datamodels/scenarios.go/FindScenarioById: %s", err)
		return nil, err

	}

	return &scenario, err
}

// Find Scenarios by IDs
const findScenariosByIds = `
SELECT * FROM scenarios	WHERE scenario_id IN (?)
`

type FindScenariosByIdsParams struct {
	ScenarioIDs []string
}

func (q *Queries) FindScenariosByIds(ctx context.Context, args FindScenariosByIdsParams) ([]*DFScenario, error) {
	_query := findScenariosByIds
	_in_args := []interface{}{args.ScenarioIDs}

	var scenarios []*DFScenario
	query, inargs, err := sqlx.In(_query, _in_args...)

	if err != nil {
		return []*DFScenario{}, err
	}
	query = q.db.Rebind(query)
	rows, err := q.db.QueryxContext(ctx, query, inargs...)

	log.Printf(findScenariosByIds)

	if err != nil {
		return []*DFScenario{}, err
	}

	for rows.Next() {
		scenario := DFScenario{}
		err := rows.StructScan(&scenario)
		if err != nil {
			return []*DFScenario{}, err
		}
		scenarios = append(scenarios, &scenario)
	}
	return scenarios, nil
}

// List Scenarios by IDs and Status
const listScenariosByIds = `
SELECT * FROM scenarios
	WHERE scenario_status IN (?)
`

type ListScenariosByIdsParams struct {
	ScenariosStatus []DFScenarioStatus
	ScenariosIDs    []string
}

func (q *Queries) ListScenariosByIds(ctx context.Context, args ListScenariosByIdsParams) ([]*DFScenario, error) {
	_query := listScenariosByIds
	_in_args := []interface{}{args.ScenariosStatus}
	if len(args.ScenariosIDs) != 0 {
		_query = listScenariosByIds + " AND scenario_id IN (?) "
		_in_args = append(_in_args, args.ScenariosIDs)
	}
	var scenarios []*DFScenario
	query, inargs, err := sqlx.In(_query, _in_args...)
	if err != nil {
		return []*DFScenario{}, err
	}
	query = q.db.Rebind(query)
	rows, err := q.db.QueryxContext(ctx, query, inargs...)
	log.Printf(_query)

	if err != nil {
		return []*DFScenario{}, err
	}
	for rows.Next() {
		scenario := DFScenario{}
		err := rows.StructScan(&scenario)
		if err != nil {
			return []*DFScenario{}, err
		}
		scenarios = append(scenarios, &scenario)
	}
	return scenarios, nil
}

// Schedule Forecast Run
const scheduleForecastRun = `
UPDATE scenarios
SET run_scheduled_timestamp = ?
WHERE scenario_id = ?
`

type ScheduleForecastRunParams struct {
	ScenarioID            string
	RunScheduledTimestamp *time.Time
}

func (q *Queries) ScheduleForecastRun(ctx context.Context, args ScheduleForecastRunParams) (*DFScenario, error) {
	res, err := q.db.ExecContext(ctx, scheduleForecastRun,
		args.RunScheduledTimestamp,
		args.ScenarioID,
	)
	if err != nil {
		return nil, err
	}

	row, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if row < 1 {
		return nil, esputils.ErrDBNoRows("update scenario")
	}

	scenario, err := q.FindScenarioById(ctx, FindScenarioByIdParams{ScenarioID: args.ScenarioID})
	if err != nil {
		return nil, err
	}

	return scenario, nil
}
