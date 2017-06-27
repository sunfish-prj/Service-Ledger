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

// Access structre for storing access data
type Access struct {
	LoggerID  string
	TimeStamp string
	Token     string
	DataType  string
	Data      string
}

//AccessLogs combined access logs for a specific ID
type AccessLogs struct {
	ID   string
	Logs []Access
}

var count int

// Init the methoed called while initialising the code
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### access chaincode Init ###########")
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
	fmt.Println("########### Invoke ###########")
	function, args := stub.GetFunctionAndParameters()

	if function != "invoke" {
		return shim.Error("Unknown function call")
	}

	if args[0] == "store" {
		return t.store(stub, args)
	}

	if args[0] == "read" {
		return t.read(stub, args)
	}

	return shim.Error("Unknown action, check the first argument, must be one of 'saveLog', 'returnIndex', or 'retrieveLogs'")
}

func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error

	//fmt.Println("Index at read:", args[1])

	returnData, err := stub.GetState(args[1])

	if err != nil {
		return shim.Error("Failed to get state")
	}
	if returnData == nil {
		return shim.Error("Entity not found")
	}

	return shim.Success(returnData)
}

func (t *SimpleChaincode) store(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	loggerID := args[1]
	timeStamp := args[2]
	token := args[3]
	dataType := args[4]
	data := args[5]
	ID := args[6]

	var tempAccess Access

	var tempAccessLogs AccessLogs

	tempAccess.LoggerID = loggerID
	tempAccess.TimeStamp = timeStamp
	tempAccess.Token = token
	tempAccess.DataType = dataType
	tempAccess.Data = data

	savedData, err := stub.GetState(ID)

	if err != nil {
		return shim.Error("Failed to get state")
	}
	if savedData == nil {
		tempAccessLogs.ID = ID
		tempAccessLogs.Logs = append(tempAccessLogs.Logs, tempAccess)
	} else {
		var savedAccessLogs AccessLogs
		if err := json.Unmarshal(savedData, &savedAccessLogs); err != nil {
			jsonResp := "At store access: error in unmarsharring JSON!"
			fmt.Println("At store access: error in unmarsharring JSON:", err)
			return shim.Error(jsonResp)
		}
		tempAccessLogs.ID = savedAccessLogs.ID
		tempAccessLogs.Logs = append(savedAccessLogs.Logs, tempAccess)
	}

	raw, err := json.Marshal(tempAccessLogs)
	if err != nil {
		jsonResp := "At store access: error in marsharring JSON!"
		fmt.Println("At store access: error in marsharring JSON:", err)
		return shim.Error(jsonResp)
	}

	err = stub.PutState(ID, raw)

	if err != nil {
		jsonResp := "Error creating the access entry!"
		fmt.Println("Error creating the access entry:", err)
		return shim.Error(jsonResp)
	}
	err = stub.SetEvent("storeAccess", []byte(ID))

	fmt.Println("Access data has been successfully stored for id:", ID)

	count++

	arr := []byte("Access data has been successfully stored!!")

	return shim.Success(arr)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
