#!/bin/sh
while true; do echo "HTTP/1.1 200 OK\n\n$(docker ps -a --format '{{.ID}} {{.Names}}')" | nc -N -l -p 8000; done
