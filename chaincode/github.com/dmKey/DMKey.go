/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// DMKey structre for Key data
type DMKey struct {
	Timestamp string
	DataID    string
	Key       string
}

var count int

// Init the methoed called while initialising the code
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### example_cc Init ###########")
	//_, args := stub.GetFunctionAndParameters()
	count = 0
	return shim.Success(nil)
}

// Query the method is invoked by the Query call, unused in this contract..
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Error("Unknown supported call")
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### example_cc Invoke ###########")
	function, args := stub.GetFunctionAndParameters()

	if function != "invoke" {
		return shim.Error("Unknown function call")
	}

	if args[0] == "register" {
		return t.register(stub, args)
	}

	if args[0] == "update" {
		return t.update(stub, args)
	}

	if args[0] == "delete" {
		return t.delete(stub, args)
	}

	if args[0] == "read" {
		return t.read(stub, args)
	}
	if args[0] == "returnIndex" {
		return t.returnIndex(stub, args)
	}
	return shim.Error("Unknown action, check the first argument, must be one of 'saveLog', 'returnIndex', or 'retrieveLogs'")
}

func (t *SimpleChaincode) returnIndex(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	index := "" + string(count)
	return shim.Success([]byte(index))
}

func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2, function followed by an index")
	}

	fmt.Println("Index at read:", args[1])

	returnData, err := stub.GetState(args[1])

	if err != nil {
		return shim.Error("Failed to get state")
	}
	if returnData == nil {
		return shim.Error("Entity not found")
	}

	return shim.Success(returnData)
}

func (t *SimpleChaincode) update(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	index := args[1]
	key := args[2]

	fmt.Println("Index at update:", index)

	returnLog, err := stub.GetState(index)

	newDMKey := DMKey{}
	json.Unmarshal(returnLog, &newDMKey)

	fmt.Println("Unmarshalled DMKey:", newDMKey)

	newDMKey.Key = key

	raw, err := json.Marshal(newDMKey)
	if err != nil {
		jsonResp := "At update: error in marsharring JSON!"
		fmt.Println("At update: error in marsharring JSON!", err)
		return shim.Error(jsonResp)
	}

	err = stub.PutState(args[1], []byte(raw))

	if err != nil {
		return shim.Error("Failed to get state")
	}
	if returnLog == nil {
		return shim.Error("Entity not found")
	}

	err = stub.SetEvent("updateKey", []byte(args[1]))

	return shim.Success([]byte("Key successfully updated!"))
}

func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2, function followed by an index")
	}

	fmt.Println("Index:", args[1])

	err := stub.DelState(args[1])

	if err != nil {
		return shim.Error("Failed to delete a state")
	}

	err = stub.SetEvent("deleteKey", []byte(args[1]))

	return shim.Success([]byte("The policy has been successfully deleted!"))
}

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) register(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	timestamp := args[1]
	dataID := args[2]
	key := args[3]

	var tempDMKey DMKey

	tempDMKey.Timestamp = timestamp
	tempDMKey.DataID = dataID
	tempDMKey.Key = key

	index := dataID + "-" + string(count)

	hasher := sha1.New()
	hasher.Write([]byte(index))
	shaIndex := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	count++

	raw, err := json.Marshal(tempDMKey)
	if err != nil {
		jsonResp := "At register: error in marsharring JSON!"
		fmt.Println("At register: error in marsharring JSON:", err)
		return shim.Error(jsonResp)
	}

	err = stub.PutState(shaIndex, raw)

	if err != nil {
		jsonResp := "Error saving the log!"
		fmt.Println("Error storing logs at index:", count)
		fmt.Println("The error is:", err)
		return shim.Error(jsonResp)
	}

	err = stub.SetEvent("registerKey", []byte(shaIndex))

	fmt.Println("Log successfully stored! Returning index:", shaIndex)

	arr := []byte(shaIndex)
	return shim.Success(arr)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
