#!/bin/sh 

docker compose down 
docker image rm server-cribbage-server -f
docker compose --profile default up --force-recreate