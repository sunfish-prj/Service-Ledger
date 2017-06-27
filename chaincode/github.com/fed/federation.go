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

// Member structre for a Federation member
type Member struct {
	Name      string
	TimeStamp string
}

// Fed structre for Federation data
type Fed struct {
	Name      string
	TimeStamp string
	Members   []Member
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

	if args[0] == "create" {
		return t.create(stub, args)
	}

	if args[0] == "add" {
		return t.add(stub, args)
	}

	if args[0] == "read" {
		return t.read(stub, args)
	}

	if args[0] == "remove" {
		return t.remove(stub, args)
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

func (t *SimpleChaincode) add(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	id := args[1]
	newMember := args[2]
	timeStamp := args[3]

	var tempFed Fed
	var tempMember Member

	tempMember.Name = newMember
	tempMember.TimeStamp = timeStamp

	fmt.Println("Federation Index at add:", id)

	returnFed, err := stub.GetState(args[1])

	if err := json.Unmarshal(returnFed, &tempFed); err != nil {
		jsonResp := "At add: error in unmarsharring JSON!"
		fmt.Println("At add: error in unmarsharring JSON:", err)
		return shim.Error(jsonResp)
	}

	tempFed.Members = append(tempFed.Members, tempMember)

	raw, err := json.Marshal(tempFed)
	if err != nil {
		jsonResp := "At add: error in marsharring JSON!"
		fmt.Println("At add: error in marsharring JSON:", err)
		return shim.Error(jsonResp)
	}

	err = stub.PutState(id, raw)

	if err != nil {
		jsonResp := "Error creating the federation entry!"
		fmt.Println("Error creating the federation entry:", err)
		return shim.Error(jsonResp)
	}

	err = stub.SetEvent("addFederation", []byte(args[1]))

	return shim.Success([]byte("The new member has been successfully added into the federation!"))
}

func (t *SimpleChaincode) remove(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("ID at remove:", args[1])

	err := stub.DelState(args[1])

	if err != nil {
		return shim.Error("Failed to delete the federation state")
	}

	err = stub.SetEvent("removeFederation", []byte(args[1]))

	return shim.Success([]byte("The federation entry has been successfully removed!"))
}

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) create(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	id := args[1]
	name := args[2]
	timeStamp := args[3]

	var tempFed Fed

	tempFed.Name = name
	tempFed.TimeStamp = timeStamp

	raw, err := json.Marshal(tempFed)
	if err != nil {
		jsonResp := "At create: error in marsharring JSON!"
		fmt.Println("At create: error in marsharring JSON:", err)
		return shim.Error(jsonResp)
	}

	err = stub.PutState(id, raw)

	if err != nil {
		jsonResp := "Error creating the federation entry!"
		fmt.Println("Error creating the federation entry:", err)
		return shim.Error(jsonResp)
	}

	err = stub.SetEvent("createFederation", []byte(args[1]))

	fmt.Println("New federation entry has been successfully created")

	arr := []byte("New federation entry has been successfully created")
	return shim.Success(arr)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
