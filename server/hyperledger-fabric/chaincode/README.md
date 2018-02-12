# Chain-code

This section contains chaincodes to be deployed on Hyperledger Fabric v1.0 for the SUNFISH federation:

- anonymisation
- governance
- keyvaluestore
- monitoring

To test a chaincode from ServiceLedger, go to the invoke API http://localhost:8090/docs/#!/invokeChaincode/invokePOST

and execute:

{
  "channel": "sunfish-channel",
  "peer": "3",
  "chaincodeName": "chaincode_sunfish_mef_min_01",
  "fcn": "getIrpef",
  "args": [
    "BRZPLA23M90E101D201710"
  ]
}