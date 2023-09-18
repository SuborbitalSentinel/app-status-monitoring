#!/bin/sh
# This is an example of what a healthcheck script from a service container could look like
# This could also probably be inlined in the Docker or the Compose file...be careful with variable expansion
curl -X POST  -d "service-id=$HOSTNAME" -d "relationship=$RELATIONSHIP" -d "parent-key=$PARENT_KEY" $SERVICE_MONITOR_URL
