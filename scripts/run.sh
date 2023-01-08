#!/bin/bash

NAME="url-collector"
DOCKERFILE="./deployment/url-collector.Dockerfile"

LINES=`docker ps -all --filter "name=${NAME}"|wc -l`
if [ ${LINES} -gt 1 ]; then
    docker rm -f ${NAME}
fi

docker image build -t ${NAME} -f ${DOCKERFILE} .

docker run -d -p 8080:8080 --name $NAME $NAME
