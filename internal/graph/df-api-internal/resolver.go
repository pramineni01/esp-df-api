package graph

import (
	"bitbucket.org/antuitinc/esp-df-api/internal/dataloaders"
	"bitbucket.org/antuitinc/esp-df-api/internal/datamodels"
)

// Resolver allows for dependency injection (into generated resolver methods)
type Resolver struct {
	DBRepo      datamodels.DBRepo
	Dataloaders dataloaders.Retriever
}
