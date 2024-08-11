# grocerycompare

## Run 

```
docker compose up
```

## DB Migration: 

```
migrate -database "postgres://postgres:example@localhost:5432/grocit?sslmode=disable" -path migrations up
```