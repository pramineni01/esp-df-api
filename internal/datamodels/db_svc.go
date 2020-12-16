package datamodels

import (
	"context"
	"database/sql"
	"time"

	model1 "bitbucket.org/antuitinc/esp-df-api/internal/graph/df-api/model"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PreparexContext(context.Context, string) (*sqlx.Stmt, error)
	QueryxContext(context.Context, string, ...interface{}) (*sqlx.Rows, error)
	QueryRowxContext(context.Context, string, ...interface{}) *sqlx.Row

	Rebind(string) string
}

type RDBTX interface {
	Set(context.Context, string, interface{}, time.Duration) *redis.StatusCmd
	Get(context.Context, string) *redis.StringCmd
	Incr(context.Context, string) *redis.IntCmd
	ZRevRange(context.Context, string, int64, int64) *redis.StringSliceCmd
	ZAdd(context.Context, string, ...*redis.Z) *redis.IntCmd
	ZRemRangeByRank(context.Context, string, int64, int64) *redis.IntCmd
}

type DBRepo interface {
	CreateForecast(context.Context, CreateForecastParams) (*DFForecast, error)
	UpdateForecast(context.Context, UpdateForecastParams) error

	CreateScenario(context.Context, CreateScenarioParams) (*DFScenarioRun, error)
	CreateScenarioComment(context.Context, CreateScenarioCommentParams) (*DFScenarioComment, error)
	UpdateScenarioRun(context.Context, UpdateScenarioRunParams) error

	DeleteScenarioById(context.Context, DeleteScenarioByIdParams) error

	ListScenariosByIds(context.Context, ListScenariosByIdsParams) ([]*DFScenario, error)

	InsertScenarioRun(context.Context, string, string) (*DFScenarioRun, error)
	FindScenario(context.Context, string, string) (*DFScenario, error)
	FindScenarios(context.Context, string, []*model1.DFScenarioScopeEntryInput) ([]*DFScenario, error)
	FindScenarioById(context.Context, FindScenarioByIdParams) (*DFScenario, error)
	FindScenariosByIds(context.Context, FindScenariosByIdsParams) ([]*DFScenario, error)
	FindScenarioCommentById(context.Context, FindScenarioCommentByIdParams) (*DFScenarioComment, error)
	FindScenarioCommentsByScenarioIds(context.Context, FindScenarioCommentsByScenarioIdsParams) ([]*DFScenarioComment, error)
	FetchAllRuns(context.Context, string) ([]*DFScenarioRun, error)
	FetchLatestRun(context.Context, string) (*DFScenarioRun, error)
	FindForecastById(context.Context, string) (*DFForecast, error)
	FindForecasts(context.Context) ([]*DFForecast, error)
	FindForecastsByIds(context.Context, FindForecastsByIdsParams) ([]*DFForecast, error)
	FindScenarioComments(context.Context) ([]*DFScenarioComment, error)
	ScheduleForecastRun(context.Context, ScheduleForecastRunParams) (*DFScenario, error)
	FindDataFilterById(context.Context, FindDataFilterByIdParams) (*DFDataFilter, error)
	CreateDataFilter(context.Context, CreateDataFilterParams) (*DFDataFilter, error)
	UpdateDataFilter(context.Context, UpdateDataFilterParams) (*DFDataFilter, error)
	DeleteDataFilter(context.Context, DeleteDataFilterParams) error
	FindDataFilters(context.Context, FindDataFiltersParams) ([]*DFDataFilter, error)

	GetForecastVersions(context.Context, GetForecastVersionsParams) ([]*DFForecastVersion, error)
	GetForecastVersion(context.Context, string) (*DFForecastVersion, error)
}

type repoSvc struct {
	*Queries
	db  *sqlx.DB
	rdb *redis.Client
}

func NewRepo(db *sqlx.DB, rdb *redis.Client) DBRepo {
	return &repoSvc{
		Queries: New(db),
		db:      db,
		rdb:     rdb,
	}
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db  DBTX
	rdb RDBTX
}
