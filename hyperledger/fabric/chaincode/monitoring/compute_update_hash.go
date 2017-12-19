package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("proposal chaincode Init")
	_, args := stub.GetFunctionAndParameters()

	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Expecting 0")
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("[COMPUTE_UPDATE_HASH CHAINCODE] update_hash chaincode Invoke")
	function, args := stub.GetFunctionAndParameters()

	// can be defined a new function here to support the vote
	if function == "computeEquals" {
		return t.computeEquals(stub, args)
	} else if function == "getValue" {
		return t.getValue(stub, args)
	}

	return shim.Error("[COMPUTE_UPDATE_HASH CHAINCODE] Invalid invoke function name. Expecting \"computeEquals\"")
}

// the function compute two given string to check if they are equals.
//
func (t *SimpleChaincode) computeEquals(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key, string1, string2 string

	if len(args) != 3 {
		return shim.Error("[COMPUTE_UPDATE_HASH CHAINCODE] Incorrect number of arguments. Expecting 3")
	}

	key = args[0]
	string1 = args[1]
	string2 = args[2]

	// chek if the strings are equals
	if string1 == string2 {
		fmt.Print("[COMPUTE_UPDATE_HASH CHAINCODE] The givens strings are equals. Proceding to the update operation...")
		// retrieve the current value in the blockchain(if exist) and update the stored string
		currentValueAsByte, err := stub.GetState(key)
		if currentValueAsByte != nil {
			newValue := string(currentValueAsByte) + string1
			// there are no difference appending string 1 ore 2 since they are the same
			err = stub.PutState(key, []byte(newValue))
			if err != nil {
				return shim.Error(err.Error())
			}
		} else {
			// if there is no value with that key, simply put the new one
			err = stub.PutState(key, []byte(string1))
			if err != nil {
				return shim.Error(err.Error())
			}
		}
	} else {
		fmt.Print("\033[1;31m[COMPUTE_UPDATE_HASH CHAINCODE] The given strings are different\033[0m\n")
		return shim.Error("[COMPUTE_UPDATE_HASH CHAINCODE] The given strings are different")
	}
	fmt.Print("\033[1;36m[COMPUTE_UPDATE_HASH CHAINCODE] The key " + string(key) + " has been updated\033[0m\n")
	return shim.Success(nil)
}

func (t *SimpleChaincode) getValue(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key string
	var err error

	if len(args) != 1 {
		return shim.Error("[COMPUTE_UPDATE_HASH CHAINCODE] Incorrect number of arguments. Expecting name of the person to query")
	}

	key = args[0]

	fmt.Printf("COMPUTE_UPDATE_HASH CHAINCODE] Executing Get: key = %s\n", key)

	// Get the state from the ledger
	value, err := stub.GetState(key)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + key + "\"}"
		return shim.Error(jsonResp)
	}
	fmt.Printf("COMPUTE_UPDATE_HASH CHAINCODE] Executing Get: key = %s\n", string(value))
	return shim.Success(value)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("[COMPUTE_UPDATE_HASH CHAINCODE] Error starting Simple chaincode: %s", err)
	}
}
