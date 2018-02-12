# API

This section contains the api to access hyperledger fabric's chaincodes.
You need to put all those scripts into all Hyperledger machines. The recommended path is */opt/gopath/src/github.com/hyperledger/fabric/examples/sunfish/scripts*. If you need to change this path remember to change it in line 12 of *hl_invoke.sh* and to modify the *hl_default_script_path* variable in the *Service-Ledger/server/config/default.yaml* file.