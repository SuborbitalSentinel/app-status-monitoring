#!/bin/sh
curl -X POST  -d "service=$HOSTNAME" -d "relationship=$RELATIONSHIP" -d "parent-key=$PARENT_KEY" $SERVICE_MONITOR_URL
