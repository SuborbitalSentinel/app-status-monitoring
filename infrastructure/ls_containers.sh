#!/bin/sh
while true; do docker ps -a --format '{{.ID}} {{.Names}}' | nc -N -l -p 8000; done
