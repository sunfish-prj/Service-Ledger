## Anonymisation Service in Hyperledger 

A sample Node.js app for anonymisation service in SUNFISH 

### Prerequisites and setup:

* [Docker](https://www.docker.com/products/overview) - v1.12 or higher
* [Docker Compose](https://docs.docker.com/compose/overview/) - v1.8 or higher
* [Git client](https://git-scm.com/downloads) - needed for clone commands
* **Node.js** v6.9.0 - 6.10.0 ( __Node v7+ is not supported__ )
* [Download Docker images](http://hyperledger-fabric.readthedocs.io/en/latest/samples.html#binaries)

```
cd fabric-samples/anonymisation/
```

Once you have completed the above setup, you will have provisioned a local network with the following docker container configuration:

* 2 CAs
* A SOLO orderer
* 4 peers (2 peers per Org)

#### Artifacts
* Crypto material has been generated using the **cryptogen** tool from Hyperledger Fabric and mounted to all peers, the orderering node and CA containers. More details regarding the cryptogen tool are available [here](http://hyperledger-fabric.readthedocs.io/en/latest/build_network.html#crypto-generator).
* An Orderer genesis block (genesis.block) and channel configuration transaction (mychannel.tx) has been pre generated using the **configtxgen** tool from Hyperledger Fabric and placed within the artifacts folder. More details regarding the configtxgen tool are available [here](http://hyperledger-fabric.readthedocs.io/en/latest/build_network.html#configuration-transaction-generator).

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

* Execute the REST APIs from the section [Sample REST APIs Requests]


```

## Sample REST APIs Requests

### Login Request

* Register and enroll new users in Organization - **Org1**:

`curl -s -X POST http://localhost:4000/users -H "content-type: application/x-www-form-urlencoded" -d 'username=Jim&orgName=org1'`

**OUTPUT:**

```
{
  "success": true,
  "secret": "RaxhMgevgJcm",
  "message": "Jim enrolled Successfully",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTQ4NjU1OTEsInVzZXJuYW1lIjoiSmltIiwib3JnTmFtZSI6Im9yZzEiLCJpYXQiOjE0OTQ4NjE5OTF9.yWaJhFDuTvMQRaZIqg20Is5t-JJ_1BP58yrNLOKxtNI"
}
```

The response contains the success/failure status, an **enrollment Secret** and a **JSON Web Token (JWT)** that is a required string in the Request Headers for subsequent requests.

### Create Channel request

```
curl -s -X POST \
  http://localhost:4000/channels \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTQ4NjU1OTEsInVzZXJuYW1lIjoiSmltIiwib3JnTmFtZSI6Im9yZzEiLCJpYXQiOjE0OTQ4NjE5OTF9.yWaJhFDuTvMQRaZIqg20Is5t-JJ_1BP58yrNLOKxtNI" \
  -H "content-type: application/json" \
  -d '{
	"channelName":"mychannel",
	"channelConfigPath":"../artifacts/channel/mychannel.tx"
}'
```

Please note that the Header **authorization** must contain the JWT returned from the `POST /users` call

### Join Channel request

```
curl -s -X POST \
  http://localhost:4000/channels/mychannel/peers \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTQ4NjU1OTEsInVzZXJuYW1lIjoiSmltIiwib3JnTmFtZSI6Im9yZzEiLCJpYXQiOjE0OTQ4NjE5OTF9.yWaJhFDuTvMQRaZIqg20Is5t-JJ_1BP58yrNLOKxtNI" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer1","peer2"]
}'
```
### Install chaincode

```
curl -s -X POST \
  http://localhost:4000/chaincodes \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTQ4NjU1OTEsInVzZXJuYW1lIjoiSmltIiwib3JnTmFtZSI6Im9yZzEiLCJpYXQiOjE0OTQ4NjE5OTF9.yWaJhFDuTvMQRaZIqg20Is5t-JJ_1BP58yrNLOKxtNI" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer1","peer2"],
	"chaincodeName":"mycc",
	"chaincodePath":"github.com/example_cc",
	"chaincodeVersion":"v0"
}'
```

### Instantiate chaincode

```
curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTQ4NjU1OTEsInVzZXJuYW1lIjoiSmltIiwib3JnTmFtZSI6Im9yZzEiLCJpYXQiOjE0OTQ4NjE5OTF9.yWaJhFDuTvMQRaZIqg20Is5t-JJ_1BP58yrNLOKxtNI" \
  -H "content-type: application/json" \
  -d '{
	"chaincodeName":"mycc",
	"chaincodeVersion":"v0",
	"args":["a","100","b","200"]
}'
```

### Invoke request

```
curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/mycc \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTQ4NjU1OTEsInVzZXJuYW1lIjoiSmltIiwib3JnTmFtZSI6Im9yZzEiLCJpYXQiOjE0OTQ4NjE5OTF9.yWaJhFDuTvMQRaZIqg20Is5t-JJ_1BP58yrNLOKxtNI" \
  -H "content-type: application/json" \
  -d '{
	"fcn":"move",
	"args":["a","b","10"]
}'
```
**NOTE:** Ensure that you save the Transaction ID from the response in order to pass this string in the subsequent query transactions.

### Chaincode Query

```
curl -s -X GET \
  "http://localhost:4000/channels/mychannel/chaincodes/mycc?peer=peer1&fcn=query&args=%5B%22a%22%5D" \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTQ4NjU1OTEsInVzZXJuYW1lIjoiSmltIiwib3JnTmFtZSI6Im9yZzEiLCJpYXQiOjE0OTQ4NjE5OTF9.yWaJhFDuTvMQRaZIqg20Is5t-JJ_1BP58yrNLOKxtNI" \
  -H "content-type: application/json"
```

### Query Block by BlockNumber

```
curl -s -X GET \
  "http://localhost:4000/channels/mychannel/blocks/1?peer=peer1" \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTQ4NjU1OTEsInVzZXJuYW1lIjoiSmltIiwib3JnTmFtZSI6Im9yZzEiLCJpYXQiOjE0OTQ4NjE5OTF9.yWaJhFDuTvMQRaZIqg20Is5t-JJ_1BP58yrNLOKxtNI" \
  -H "content-type: application/json"
```

### Query Transaction by TransactionID

```
curl -s -X GET http://localhost:4000/channels/mychannel/transactions/TRX_ID?peer=peer1 \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTQ4NjU1OTEsInVzZXJuYW1lIjoiSmltIiwib3JnTmFtZSI6Im9yZzEiLCJpYXQiOjE0OTQ4NjE5OTF9.yWaJhFDuTvMQRaZIqg20Is5t-JJ_1BP58yrNLOKxtNI" \
  -H "content-type: application/json"
```
**NOTE**: Here the TRX_ID can be from any previous invoke transaction


### Query ChainInfo

```
curl -s -X GET \
  "http://localhost:4000/channels/mychannel?peer=peer1" \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTQ4NjU1OTEsInVzZXJuYW1lIjoiSmltIiwib3JnTmFtZSI6Im9yZzEiLCJpYXQiOjE0OTQ4NjE5OTF9.yWaJhFDuTvMQRaZIqg20Is5t-JJ_1BP58yrNLOKxtNI" \
  -H "content-type: application/json"
```

### Query Installed chaincodes

```
curl -s -X GET \
  "http://localhost:4000/chaincodes?peer=peer1&type=installed" \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTQ4NjU1OTEsInVzZXJuYW1lIjoiSmltIiwib3JnTmFtZSI6Im9yZzEiLCJpYXQiOjE0OTQ4NjE5OTF9.yWaJhFDuTvMQRaZIqg20Is5t-JJ_1BP58yrNLOKxtNI" \
  -H "content-type: application/json"
```

### Query Instantiated chaincodes

```
curl -s -X GET \
  "http://localhost:4000/chaincodes?peer=peer1&type=instantiated" \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTQ4NjU1OTEsInVzZXJuYW1lIjoiSmltIiwib3JnTmFtZSI6Im9yZzEiLCJpYXQiOjE0OTQ4NjE5OTF9.yWaJhFDuTvMQRaZIqg20Is5t-JJ_1BP58yrNLOKxtNI" \
  -H "content-type: application/json"
```

### Query Channels

```
curl -s -X GET \
  "http://localhost:4000/channels?peer=peer1" \
  -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTQ4NjU1OTEsInVzZXJuYW1lIjoiSmltIiwib3JnTmFtZSI6Im9yZzEiLCJpYXQiOjE0OTQ4NjE5OTF9.yWaJhFDuTvMQRaZIqg20Is5t-JJ_1BP58yrNLOKxtNI" \
  -H "content-type: application/json"
```

ToDO: license issue
