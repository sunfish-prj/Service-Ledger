#!/bin/bash
PEER=$1
CHANNEL=$2
CHAINCODE_NAME=$3
KEY=$4
VALUE=$5
DOCKER_ID=$6

docker exec $DOCKER_ID /bin/bash -c ". /opt/gopath/src/github.com/hyperledger/fabric/peer/scripts/put.sh $PEER $CHANNEL $CHAINCODE_NAME $KEY $VALUE $DOCKER_ID"

exit 0
