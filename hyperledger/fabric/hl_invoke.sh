#!/bin/bash
PEER=$1
CHANNEL=$2
CHAINCODE_NAME=$3
ARGS=$4
DOCKER_ID=$5

docker exec $DOCKER_ID /bin/bash -c ". /opt/gopath/src/github.com/hyperledger/fabric/peer/scripts/invoke.sh $PEER $CHANNEL $CHAINCODE_NAME $ARGS $DOCKER_ID"

exit 0
