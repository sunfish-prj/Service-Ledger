# Distributed Infrastructure Underlying FaaS Cloud Federations

Node.js server acting as a facade for the distributed infrastructure underlying a FaaS Cloud Federation. Such infrastructure relies on blockchain system distributed among the federation peers and featuring smart contract. 

This respository contains the Node.js files implementing a REST *server acting as medium for the blockchain system*. Indeed, *different blockchain system implementations can be used*, the configuartion paramenters used by the APIs just need to be changed. 

The server has been currently configured for the blockchain system [Hyperledger Fabric](https://hyperledger-fabric.readthedocs.io/en/latest/) v1.0.0. Its installation and deployment instruction can be found in the guide. Wheras, the smart contracts to deployed, which actually implement the functionality of the infrastructure, are instead available in the repository. 

For testing purpose, part of the infrastructure functionality are also implemented via [MongoDB](https://www.mongodb.com/en). Clearly, it only implements functionality related to management of data, not computational one. 

Full documentation is reported in the official [SUNFISH Manual](http://sunfish-platform-docs.readthedocs.io/)

## Installation Guide



