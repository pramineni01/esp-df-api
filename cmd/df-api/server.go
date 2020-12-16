package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jmoiron/sqlx"

	"bitbucket.org/antuitinc/esp-df-api/internal/dataloaders"
	"bitbucket.org/antuitinc/esp-df-api/internal/datamodels"
	graph "bitbucket.org/antuitinc/esp-df-api/internal/graph/df-api"
	"bitbucket.org/antuitinc/esp-df-api/internal/graph/df-api/generated"

	"bitbucket.org/antuitinc/esp-df-api/pkg/esputils"
	"github.com/go-redis/redis/v8"
)

const defaultPort = "8081"

// NewHandler returns a new graphql endpoint handler.
func NewHandler(repo datamodels.DBRepo, dl dataloaders.Retriever) http.Handler {
	return handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		DBRepo:      repo,
		Dataloaders: dl,
	},
	}))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	dbConfig := esputils.GetDBConfig()
	db, err := sqlx.Open("mysql", dbConfig.FormatDSN())

	if err != nil {
		log.Fatalln(err)
	}

	defer db.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "connection string",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	dbrepo := datamodels.NewRepo(db, rdb)
	dataloader := dataloaders.NewRetriever()
	//here we initialize the middleware
	dlMiddleware := dataloaders.Middleware(dbrepo)
	queryHandler := NewHandler(dbrepo, dataloader)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", dlMiddleware(queryHandler))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
