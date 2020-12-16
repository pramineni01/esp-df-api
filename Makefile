build:
	cd cmd/df-api && go build
	cd cmd/df-api-internal && go build
	cd cmd/df-api-internal-client && go build
fmt:
	cd cmd/df-api && go fmt
	cd cmd/df-api-internal && go fmt
generate:
	cd cmd/df-api && go run github.com/99designs/gqlgen generate
	cd cmd/df-api-internal && go run github.com/99designs/gqlgen generate
run:
	cd cmd/df-api && go run server.go
run-internal:
	cd cmd/df-api-internal && go run server.go
run-internal-client:
	cd cmd/df-api-internal-client && go run df-api-internal-client.go $(ARGS)
# Invoke: make run-internal-client ARGS='updateForecast -forecastId=1 -latestVersionDimensionMemberId=1'
