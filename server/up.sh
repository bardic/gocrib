#!/bin/sh

migrate -database "postgres://postgres:example@db:5432/grocit?sslmode=disable" -path migrations up
./grocit