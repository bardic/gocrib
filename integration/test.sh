#!/bin/sh

migrate -path ../queries/migrations -database postgres://postgres:example@localhost:5432/cribbage?sslmode=disable down
migrate -path ../queries/migrations -database postgres://postgres:example@localhost:5432/cribbage?sslmode=disable up

ijhttp http/match.http
