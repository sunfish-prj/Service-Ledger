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

// Fed structre for Federation data
type Metrics struct {
	RequestorID  string
	TimeStamp    string
	Token        string
	Availability string
	ResponseTime string
	CPUPower     string
	DiskSpace    string
	MemorySize   string
	Bandwidth    string
	Throughput   string
	Connections  string
	Elasticity   string
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
	availability := args[4]
	responseTime := args[5]
	CPUPower := args[6]
	diskSpace := args[7]
	memorySize := args[8]
	bandwidth := args[9]
	throughput := args[10]
	connections := args[11]
	elasticity := args[12]

	var tempMetrics Metrics

	tempMetrics.RequestorID = requestorID
	tempMetrics.TimeStamp = timeStamp
	tempMetrics.Token = token
	tempMetrics.Availability = availability
	tempMetrics.ResponseTime = responseTime
	tempMetrics.CPUPower = CPUPower
	tempMetrics.DiskSpace = diskSpace
	tempMetrics.MemorySize = memorySize
	tempMetrics.Bandwidth = bandwidth
	tempMetrics.Throughput = throughput
	tempMetrics.Connections = connections
	tempMetrics.Elasticity = elasticity

	raw, err := json.Marshal(tempMetrics)
	if err != nil {
		jsonResp := "At store: error in marsharring JSON!"
		fmt.Println("At store: error in marsharring JSON:", err)
		return shim.Error(jsonResp)
	}

	err = stub.PutState(strconv.Itoa(count), raw)

	if err != nil {
		jsonResp := "Error creating the federation entry!"
		fmt.Println("Error creating the federation entry:", err)
		return shim.Error(jsonResp)
	}

	err = stub.SetEvent("storeMetrics", []byte(strconv.Itoa(count)))

	fmt.Println("SLA Metrics have been successfully stored at index:", count)

	count++

	arr := []byte("SLA Metrics have been successfully stored!")

	return shim.Success(arr)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

