FROM golang:latest as builder

WORKDIR /esp-df-api
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o esp-df-api cmd/df-api/server.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o esp-df-api-internal cmd/df-api-internal/server.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o esp-df-api-internal-client cmd/df-api-internal-client/df-api-internal-client.go


FROM alpine:latest

WORKDIR /opt/esp-df-api/
COPY --from=builder /esp-df-api/esp-df-api external-api
COPY --from=builder /esp-df-api/esp-df-api-internal internal-api
COPY --from=builder /esp-df-api/esp-df-api-internal-client internal-api-client

ENTRYPOINT ["./external-api"]