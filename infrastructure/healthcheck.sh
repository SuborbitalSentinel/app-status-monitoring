#!/bin/sh

# relationship: Choices ["", "parent", "child"]
# if parent, key will be registered so that children can group to it
# if child, key will group with parent
# if empty, service is standalone; key is ignored

# parent-key

curl -X POST -d "relationship=parent" -d "service=$HOSTNAME" -d "parent-key=domain" $SERVICE_MONITOR_URL
