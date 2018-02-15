package main

import (
	"bytes"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// This chaincode allow to define a key-value store database
// which map a value to a key

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("key_value_store chaincode Init")
	_, args := stub.GetFunctionAndParameters()

	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Expecting 0")
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("key_value_store chaincode Invoke")
	function, args := stub.GetFunctionAndParameters()

	if function == "put" {
		// put a key-value couple
		return t.put(stub, args)
	} else if function == "delete" {
		// deletes a key-value couple
		return t.delete(stub, args)
	} else if function == "get" {
		// get a value associated to a given key
		return t.get(stub, args)
	} else if function == "getAll" {
		// get all the keys
		return t.getAll(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"put\" \"delete\" \"get\"")
}

// Put a key-value couple
func (t *SimpleChaincode) put(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key, value string
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	// Setting chaincode args
	key = args[0]
	value = args[1]

	fmt.Printf("Executing Put: key = %s, value = %s\n", key, value)

	// Write the state to the ledger
	err = stub.PutState(key, []byte(value))
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("Put executed!\n")
	return shim.Success(nil)
}

// Deletes a key-value couple
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	key := args[0]

	fmt.Printf("Executing Delete: key = %s\n", key)

	// Delete the key from the state in ledger
	err := stub.DelState(key)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	fmt.Printf("Delete executed!\n")
	return shim.Success(nil)
}

// Get a value froma a key
func (t *SimpleChaincode) get(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	key = args[0]

	fmt.Printf("Executing Get: key = %s\n", key)

	// Get the state from the ledger
	value, err := stub.GetState(key)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + key + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"key\":\"" + key + "\",\"value\":\"" + string(value) + "\"}"
	fmt.Printf("Get executed! Response: %s\n", jsonResp)
	return shim.Success(value)
}

// Get all function to retrieve all the keys
func (t *SimpleChaincode) getAll(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var startKey string
	var endKey string
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2: startKey and endKey. Leave both null to get all the keys.")
	}
	startKey = args[0]
	endKey = args[1]

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")
	buffer.WriteString("{\"Keys\":")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("}")
	buffer.WriteString("]")

	fmt.Printf("get all the keys :\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
