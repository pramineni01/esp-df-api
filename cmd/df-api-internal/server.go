package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jmoiron/sqlx"

	"bitbucket.org/antuitinc/esp-df-api/internal/datamodels"
	graph "bitbucket.org/antuitinc/esp-df-api/internal/graph/df-api-internal"
	"bitbucket.org/antuitinc/esp-df-api/internal/graph/df-api-internal/generated"

	"bitbucket.org/antuitinc/esp-df-api/pkg/esputils"
)

const defaultPort = "8082"

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

	dbrepo := datamodels.NewRepo(db, nil)
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		DBRepo: dbrepo,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://0.0.0.0:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
