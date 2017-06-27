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
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// Anon structre for Federation data
type Anon struct {
	DataProvider string
	DataConsumer string
	TimeStamp    string
	DataID       string
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

// Invoke wrapper method
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### example_cc Invoke ###########")
	function, args := stub.GetFunctionAndParameters()

	if function != "invoke" {
		return shim.Error("Unknown function call")
	}

	if args[0] == "register" {
		return t.register(stub, args)
	}

	if args[0] == "read" {
		return t.read(stub, args)
	}

	return shim.Error("Unknown action, check the first argument, must be one of 'saveLog', 'returnIndex', or 'retrieveLogs'")
}

func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error

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

func (t *SimpleChaincode) register(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	timeStamp := args[1]
	dataProvider := args[2]
	dataConsumer := args[3]
	dataID := args[4]

	var tempAnon Anon

	tempAnon.DataProvider = dataProvider
	tempAnon.TimeStamp = timeStamp
	tempAnon.DataConsumer = dataConsumer
	tempAnon.DataID = dataID

	raw, err := json.Marshal(tempAnon)
	if err != nil {
		jsonResp := "At register anon: error in marsharring JSON!"
		fmt.Println("At register anon: error in marsharring JSON:", err)
		return shim.Error(jsonResp)
	}

	err = stub.PutState(dataConsumer, raw)

	if err != nil {
		jsonResp := "Error creating the anon entry!"
		fmt.Println("Error creating the anon entry:", err)
		return shim.Error(jsonResp)
	}

	err = stub.SetEvent("registerAnon", []byte(dataConsumer))

	fmt.Println("Anonymisation metrics have been successfully registered!")

	count++

	arr := []byte("Anonymisation metrics have been successfully registered!")

	return shim.Success(arr)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
