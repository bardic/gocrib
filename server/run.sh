#!/bin/sh 

docker compose down 
docker image rm cribbage-server -f
docker compose up --force-recreate