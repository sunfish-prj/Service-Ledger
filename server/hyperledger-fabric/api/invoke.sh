#!/bin/bash

PEER=$1
CHANNEL=$2
CHAINCODE_NAME=$3
ARGS=$4
DOCKER_ID=$5
#SLEEP=5 #second to wait the commit on the blockchain before printing the response

if [ -z "$1" ] ; then
	echo "No argument supplied. Setting Default to 0."
	PEER=0
fi

echo "Setting environment for peer $PEER..."

# setting orderer variables
export ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/orderer/localMspConfig/cacerts/ordererOrg0.pem

# setting peer variables
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peer/peer$PEER/localMspConfig
export CORE_PEER_ADDRESS=peer$PEER:7051

if [ $PEER -eq 0 -o $PEER -eq 1 ] ; then
	export CORE_PEER_LOCALMSPID="Org0MSP"
	export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peer/peer$PEER/localMspConfig/cacerts/peerOrg0.pem
else
	export CORE_PEER_LOCALMSPID="Org1MSP"
	export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peer/peer$PEER/localMspConfig/cacerts/peerOrg1.pem
fi

env | grep CORE

ARGS_SPLIT=$(echo \"$ARGS\" | sed 's/\,/","/g')

ARGSHL={\"Args\":[$ARGS_SPLIT]}

echo "[invoke] args splitted: $ARGS_SPLIT"
echo "[invoke] args for fabric: $ARGSHL"

exec /usr/local/bin/peer chaincode invoke -o orderer0:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL -n $CHAINCODE_NAME -c $ARGSHL >&scripts/result.log & echo "[invoke] executing invoke on the peer..."

exit 0


