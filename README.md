# gocrib

Async TUI multiplayer cribbage game written in Go

## Server

### Run DB

```
docker compose --profile db up -d
```

### DB Migration: 

```
migrate -database "postgres://postgres:example@localhost:5432/cribbage?sslmode=disable" -path migrations up
```

### Build

```
cd server
go build .
```

## Client

###

```
cd cli
go build .
```