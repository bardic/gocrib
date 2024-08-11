FROM golang:1.22.0 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go get -u github.com/swaggo/echo-swagger

COPY model ./model
COPY route ./route
COPY *.go ./

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN swag init --parseDependency --parseInternal
RUN CGO_ENABLED=0 GOOS=linux go build -o ./grocit
COPY up.sh ./
COPY migrations/ ./migrations

ENTRYPOINT ["/bin/sh","./up.sh"]