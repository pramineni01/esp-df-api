package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"bitbucket.org/antuitinc/esp-df-api/internal/config"
	"bitbucket.org/antuitinc/esp-df-api/internal/datamodels"
	"bitbucket.org/antuitinc/esp-df-api/pkg/sdk"
)

// ./df-api-internal-client updateForecast -forecastId XXX -latestVersionDimensionMemberId YYY
// ./df-api-internal-client updateScenarioRun -scenarioRunID XXX -status YYY [-daVersion ZZZ]

func main() {

	ctx := context.Background()

	updateForecastCmd := flag.NewFlagSet("updateForecast", flag.ExitOnError)
	updateForecastID := updateForecastCmd.String("forecastId", "", "forecastId")
	updateForecastDimMemID := updateForecastCmd.String("latestVersionDimensionMemberId", "", "latestVersionDimensionMemberId")

	updateScenarioRunCmd := flag.NewFlagSet("updateScenarioRun", flag.ExitOnError)
	updateScenarioRunID := updateScenarioRunCmd.String("scenarioRunID", "", "scenarioRunID")
	updateScenarioRunStatus := updateScenarioRunCmd.String("status", "", "status")
	updateScenarioRunDaVersion := updateScenarioRunCmd.String("daVersion", "", "daVersion") // optional

	if len(os.Args) < 2 {
		fmt.Println("Expected 'updateForecast' or 'updateScenarioRun' subcommands")
		os.Exit(1)
	}

	client := sdk.NewClient(ctx, config.ENDPOINT_DF_API_INTERNAL())

	switch os.Args[1] {

	case "updateForecast":
		updateForecastCmd.Parse(os.Args[2:])
		fmt.Println("subcommand 'updateForecast'")
		if updateForecastID == nil || len(*updateForecastID) == 0 {
			fmt.Println("Expected 'forecastId' in 'updateForecast' subcommands")
			os.Exit(1)
		}
		if updateForecastDimMemID == nil || len(*updateForecastDimMemID) == 0 {
			fmt.Println("Expected 'latestVersionDimensionMemberId' in 'updateForecast' subcommands")
			os.Exit(1)
		}
		fmt.Println("  forecastId:", *updateForecastID)
		fmt.Println("  latestVersionDimensionMemberId:", *updateForecastDimMemID)

		result, err := client.UpdateForecast(*updateForecastID, *updateForecastDimMemID)
		if err != nil {
			log.Printf("ERROR: %s", err)
			os.Exit(1)
		}
		log.Printf("Result: %t", *result)
	case "updateScenarioRun":
		updateScenarioRunCmd.Parse(os.Args[2:])
		fmt.Println("subcommand 'updateScenarioRun'")
		if updateScenarioRunID == nil || len(*updateScenarioRunID) == 0 {
			fmt.Println("Expected 'scenarioRunID' in 'updateScenarioRun' subcommands")
			os.Exit(1)
		}
		if updateScenarioRunStatus == nil || len(*updateScenarioRunStatus) == 0 {
			fmt.Println("Expected 'status' in 'updateScenarioRun' subcommands")
			os.Exit(1)
		}

		switch *updateScenarioRunStatus {
		case datamodels.DFScenarioRunStatusScheduled:
		case datamodels.DFScenarioRunStatusInProgress:
		case datamodels.DFScenarioRunStatusForecasted:
		case datamodels.DFScenarioRunStatusError:
		default:
			fmt.Println("Invalid 'status' in 'updateScenarioRun': SCHEDULED, IN_PROGRESS, FORECASTED or ERROR")
			os.Exit(1)
		}

		fmt.Println("  scenarioRunID:", *updateScenarioRunID)
		fmt.Println("  status:", *updateScenarioRunStatus)
		fmt.Println("  daVersion (optional):", *updateScenarioRunDaVersion)
		result, err := client.UpdateScenarioRun(*updateScenarioRunID, *updateScenarioRunStatus, updateScenarioRunDaVersion)
		if err != nil {
			log.Printf("ERROR: %s", err)
			os.Exit(1)
		}
		log.Printf("Result: %t", *result)
	default:
		fmt.Println("Expected 'updateForecast' or 'updateScenarioRun' subcommands")
		os.Exit(1)
	}
}
