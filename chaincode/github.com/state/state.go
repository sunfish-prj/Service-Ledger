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

// State structre for Federation data
type State struct {
	ServiceID string
	URL       string
	Name      string
}

// List structre for Federation data
type List struct {
	Name  string
	Value string
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

	if args[0] == "store" {
		return t.store(stub, args)
	}

	if args[0] == "update" {
		return t.update(stub, args)
	}

	if args[0] == "read" {
		return t.read(stub, args)
	}

	if args[0] == "delete" {
		return t.delete(stub, args)
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

func (t *SimpleChaincode) update(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	serviceID := args[1]
	list := args[2]

	fmt.Println("List:", list)

	var tempList []List

	if err := json.Unmarshal([]byte(list), &tempList); err != nil {
		jsonResp := "At state update: error in unmarsharring JSON!"
		fmt.Println("At state update: error in unmarsharring JSON:", err)
		return shim.Error(jsonResp)
	}

	var tempState State

	returnFed, err := stub.GetState(serviceID)

	if err := json.Unmarshal(returnFed, &tempState); err != nil {
		jsonResp := "At state update: error in unmarsharring JSON!"
		fmt.Println("At state update: error in unmarsharring JSON:", err)
		return shim.Error(jsonResp)
	}

	for index, value := range tempList {
		if index == 0 {
			tempState.Name = value.Value
		}
		if index == 1 {
			tempState.URL = value.Value
		}

	}

	raw, err := json.Marshal(tempState)
	if err != nil {
		jsonResp := "At state update: error in marsharring JSON!"
		fmt.Println("At state update: error in marsharring JSON:", err)
		return shim.Error(jsonResp)
	}

	err = stub.PutState(serviceID, raw)

	if err != nil {
		jsonResp := "Error creating the federation entry!"
		fmt.Println("Error creating the federation entry:", err)
		return shim.Error(jsonResp)
	}

	err = stub.SetEvent("updateState", []byte(args[1]))

	return shim.Success([]byte("The new member has been successfully added into the federation!"))
}

func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("ID at remove:", args[1])

	err := stub.DelState(args[1])

	if err != nil {
		return shim.Error("Failed to delete the state data")
	}

	err = stub.SetEvent("deleteState", []byte(args[1]))

	return shim.Success([]byte("The federation entry has been successfully removed!"))
}

func (t *SimpleChaincode) store(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	serviceID := args[1]
	URL := args[2]
	name := args[3]

	var tempState State

	tempState.ServiceID = serviceID
	tempState.URL = URL
	tempState.Name = name

	raw, err := json.Marshal(tempState)
	if err != nil {
		jsonResp := "At create: error in marsharring JSON!"
		fmt.Println("At create: error in marsharring JSON:", err)
		return shim.Error(jsonResp)
	}

	err = stub.PutState(serviceID, raw)

	if err != nil {
		jsonResp := "Error storing state entries!"
		fmt.Println("Error storing state entries:", err)
		return shim.Error(jsonResp)
	}

	err = stub.SetEvent("storeState", []byte(serviceID))

	fmt.Println("New state entries have been successfully stored!")

	arr := []byte("New state entries have been successfully stored!")
	return shim.Success(arr)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
