#!/bin/sh

# This is just a stupid script to create a dead simple http server that responds with all docker container ID:Name pairs
while true; do echo "HTTP/1.1 200 OK\n\n$(docker ps -a --format '{{.ID}} {{.Names}}')" | nc -N -l -p 8000; done
