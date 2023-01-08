#!/bin/bash

NAME="url-collector"

set -e

LINES=`docker ps -all --filter "name=${NAME}"|wc -l`
if [ ${LINES} -gt 1 ]; then
    docker rm -f ${NAME}
fi
