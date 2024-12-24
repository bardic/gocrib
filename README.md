# gocrib

Async TUI multiplayer cribbage game written in Go

## Generate / Update queries 

```
dagger call sqlc --src=. export --path=sql/queries
```

## Server

```
dagger call game-server --src=. up
```

## Integration Test

```
dagger call http --src=.
```

## Client

```
cd cli
go run .
```