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
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// Alert structre for alert data
type Alert struct {
	RequestorID string
	TimeStamp   string
	Token       string
	AlertID     string
	AlertType   string
	AlertSource string
	Alert       string
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

func (t *SimpleChaincode) store(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	requestorID := args[1]
	timeStamp := args[2]
	token := args[3]
	alertID := args[4]
	alertType := args[5]
	alertSource := args[6]
	alert := args[7]

	var tempAlert Alert

	tempAlert.RequestorID = requestorID
	tempAlert.TimeStamp = timeStamp
	tempAlert.Token = token
	tempAlert.AlertID = alertID
	tempAlert.AlertType = alertType
	tempAlert.AlertSource = alertSource
	tempAlert.Alert = alert

	raw, err := json.Marshal(tempAlert)
	if err != nil {
		jsonResp := "At register anon: error in marsharring JSON!"
		fmt.Println("At register anon: error in marsharring JSON:", err)
		return shim.Error(jsonResp)
	}

	err = stub.PutState(strconv.Itoa(count), raw)

	if err != nil {
		jsonResp := "Error creating the federation entry!"
		fmt.Println("Error creating the federation entry:", err)
		return shim.Error(jsonResp)
	}

	err = stub.SetEvent("storeAlert", []byte(strconv.Itoa(count)))

	fmt.Println("Alert has been successfully stored!")

	count++

	arr := []byte("Alert has been successfully stored!!")

	return shim.Success(arr)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

