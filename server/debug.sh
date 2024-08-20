#!/bin/sh 

docker compose --profile default down 
docker compose --profile db up --force-recreate