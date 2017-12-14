In this section there are stored the chaincodes supported by SUNFISH.
In the following there are shown the main folders and the instructions to set up a network, to install and instantiate the chaincodes. 

*Main directories* 

hyperledger directory:
/opt/gopath/src/github.com/hyperledger/fabric/

hyperledger working directory:
/opt/gopath/src/github.com/hyperledger/fabric/examples/e2e_cli

hyperledger script directory:
/opt/gopath/src/github.com/hyperledger/fabric/examples/e2e_cli/scripts

hyperledger chaincode directory:
/opt/gopath/src/github.com/hyperledger/fabric/examples/chaincode/go

serviceledger directory:
/opt/gopath/src/github.com/Service-Ledger/server




*Hyperledger Scripts and Commands*
Before each step there are what the command does ad there is also defined in which folder it must be executed.
    
From the "hyperledger working directory":

    1) create the network
    CHANNEL_NAME=mychannel docker-compose -f docker-compose-no-cli.yaml up -d 2>&1

    2) check deployed containers (outsidedocker container)
    docker ps

On a 2nd terminal inside the "hyperledger script directory"
    3) start a 
          ./ update-logs.sh
       or can be used:
          COMPOSE_HTTP_TIMEOUT=<set a desired timeout> docker-compose logs -f

On the 1st terminal, create the channel and join all peers to it (inside docker container)
    4) docker exec -it cli bash
       . scripts/script-solo-channel.sh

Install the desired chaincode "<chaincode name>" 
    5) . scripts/set-peer.sh 0 
       peer chaincode install -n <chaincode name> -v 1.0 -p github.com/hyperledger/fabric/examples/chaincode/go/<path of the chaincode> >&log.txt
       (e.g. per governance: peer chaincode install -n governance -v 1.0 -p github.com/hyperledger/fabric/examples/chaincode/governance >&log.txt)
    (note: the <chaincode name> used it is only a identifier that hyperledger uses for the chaincode contained in the selected folder )

Instantiate the chaincode
    6)  peer chaincode instantiate -o orderer0:7050 --tls true --cafile $ORDERER_CA -C mychannel -n <chaincode name> -v 1.0 -c '{"Args":["init"]}' -P "OR   ('Org0MSP.member','Org1MSP.member')" >&log.txt
    
    (e.g per governance: peer chaincode instantiate -o orderer0:7050 --tls true --cafile $ORDERER_CA -C mychannel -n governance -v 1.0 -c '{"Args":["init"]}' -P "OR ('Org0MSP.member','Org1MSP.member')")

On a 3rd terminal, attach to the log of the chaincode
    7)  docker ps 
        (to get the ID of the chaincode)

        docker attach <ID>
        (to attach to the desired code and view logs messages)

On the 1st terminal, (inside the docker cli) invoke the chaincode
    8)  peer chaincode invoke -o orderer0:7050  --tls true --cafile $ORDERER_CA -C mychannel -n <chaincode name> -c '{"Args":["function to call and args"]}' >&log.txt
        (per governance function "submitProposal":  peer chaincode invoke -o orderer0:7050  --tls true --cafile $ORDERER_CA -C mychannel -n governance -c '{"Args":["submitProposal","member01","proposal01", "testing proposal description", "JOIN", "majority", "5"]}' )


To shutdown all the network
    9)  ./network-setup.sh down
