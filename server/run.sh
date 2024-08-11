#!/bin/sh 

docker compose down && docker image rm grocit -f &&  docker compose up --force-recreate