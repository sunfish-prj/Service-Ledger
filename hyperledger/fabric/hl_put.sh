#!/bin/bash
CHANNEL=$1
CHAINCODE_NAME=$2
KEY=$3
VALUE=$4
DOCKER_ID=$5

docker exec $DOCKER_ID /bin/sh -c "peer chaincode invoke -o orderer0:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL -n $CHAINCODE_NAME -c '{"Args":["put","$KEY","$VALUE"]}'"

exit 0
