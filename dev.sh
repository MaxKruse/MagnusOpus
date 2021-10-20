#!/bin/bash

# Reset Everything
docker-compose down
docker system prune -f

rm -rf backend/storage/*.osu

docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build 
