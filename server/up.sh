#!/bin/sh

migrate -database "postgres://postgres:example@localhost:5432/cribbage?sslmode=disable" -path migrations up
./cribbage-server