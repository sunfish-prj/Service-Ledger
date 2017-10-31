## Anonymisation Service in Hyperledger 

A sample Node.js app for anonymisation service in SUNFISH 

### Prerequisites and setup:

* [Docker](https://www.docker.com/products/overview) - v1.12 or higher
* [Docker Compose](https://docs.docker.com/compose/overview/) - v1.8 or higher
* [Git client](https://git-scm.com/downloads) - needed for clone commands
* **Node.js** v6.9.0 - 6.10.0 ( __Node v7+ is not supported__ )
* [Download Platform-specific Docker images]
```
curl -sSL https://goo.gl/Q3YRTi | bash
```

Once you have completed the above setup, you will have provisioned a local network with the following docker container configuration:

* 2 CAs
* A SOLO orderer
* 4 peers (2 peers per Org)

#### Artifacts
* Crypto material has been generated using the **cryptogen** tool from Hyperledger Fabric and mounted to all peers, the orderering node and CA containers. More details regarding the cryptogen tool are available [here](http://hyperledger-fabric.readthedocs.io/en/latest/build_network.html#crypto-generator).
* An Orderer genesis block (genesis.block) and channel configuration transaction (mychannel.tx) has been pre generated using the **configtxgen** tool from Hyperledger Fabric and placed within the artifacts folder. More details regarding the configtxgen tool are available [here](http://hyperledger-fabric.readthedocs.io/en/latest/build_network.html#configuration-transaction-generator).
* the anonymisation chaincode is in ./artifacts/src/github.com/anonymisation_cc

## Running the sample program

##### Terminal Window 1

```
cd fabric-samples/anonymisation

./runApp.sh

```
* This lauches the required network on your local machine
* Installs the fabric-client and fabric-ca-client node modules
* And, starts the node app on PORT 4000

##### Terminal Window 2

* Execute the REST APIs from Registry API using
  the provided postman collection


ToDO: license issue
