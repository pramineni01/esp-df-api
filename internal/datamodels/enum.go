package datamodels

type DFScenarioStatus string

const (
	ScenarioStatusCURRENT     = "CURRENT"
	ScenarioStatusDELETED     = "DELETED"
	ScenarioStatusPROMOTED    = "PROMOTED"
	ScenarioStatusSUPERSCEDED = "SUPERSCEDED"
)

type DFScenarioRunStatus string

const (
	DFScenarioRunStatusScheduled  = "SCHEDULED"
	DFScenarioRunStatusInProgress = "IN_PROGRESS"
	DFScenarioRunStatusForecasted = "FORECASTED"
	DFScenarioRunStatusError      = "ERROR"
)
