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
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// PolList The array of policies..
type PolList struct {
	Id     string
	Policy string
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

	if args[0] == "save" {
		return t.savePolicy(stub, args)
	}

	if args[0] == "update" {
		return t.update(stub, args)
	}

	if args[0] == "service" {
		return t.service(stub, args)
	}

	if args[0] == "delete" {
		return t.delete(stub, args)
	}

	if args[0] == "read" {
		return t.retrievePolicy(stub, args)
	}
	return shim.Error("Unknown action, check the first argument, must be one of 'saveLog', 'returnIndex', or 'retrieveLogs'")
}

func (t *SimpleChaincode) retrievePolicy(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error

	// if len(args) != 2 {
	// 	return shim.Error("Incorrect number of arguments. Expecting 2, function followed by an index")
	// }

	fmt.Println("Index:", args[1])

	returnLog, err := stub.GetState(args[1])

	fmt.Println("Return log at retrievePolicy:", returnLog)

	if err != nil {
		return shim.Error("Failed to get state")
	}
	if returnLog == nil {
		return shim.Error("Entity not found")
	}

	return shim.Success(returnLog)
}

func (t *SimpleChaincode) update(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	// if len(args) != 3 {
	// 	return shim.Error("Incorrect number of arguments. Expecting 2, function followed by an index")
	// }

	fmt.Println("Index at update:", args[1])

	returnLog, err := stub.GetState(args[1])

	//err = stub.SetEvent("evtsender", []byte(tosend))

	polArray := strings.Split(string(returnLog), ",")
	fmt.Println("polArray at update:", polArray)
	polString := polArray[0] + "," + polArray[1] + "," + polArray[2] + "," + polArray[3] + "," + args[2]

	fmt.Println("polString at update:", polString)

	err = stub.PutState(args[1], []byte(polString))

	if err != nil {
		return shim.Error("Failed to get state")
	}
	if returnLog == nil {
		return shim.Error("Entity not found")
	}

	err = stub.SetEvent("updatePolicy", []byte(args[1]))

	return shim.Success(returnLog)
}

func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// if len(args) != 2 {
	// 	return shim.Error("Incorrect number of arguments. Expecting 2, function followed by an index")
	// }

	fmt.Println("Index at delete:", args[1])

	delPolData, errPol := stub.GetState(args[1])

	if errPol != nil {
		return shim.Error("Failed to get the service state!")
	}

	fmt.Println("At delete: delPolData:", delPolData)

	polDataArray := strings.Split(string(delPolData), ",")

	fmt.Println("At delete: polDataArray:", polDataArray)

	serviceID := polDataArray[1]
	policyType := polDataArray[3]

	serviceData, errServ := stub.GetState("serv-" + serviceID + "-" + policyType)

	fmt.Println("At delete: serviceData:", serviceData)

	if errServ != nil {
		return shim.Error("Failed to get the service state!")
	}

	serviceDataArray := strings.Split(string(serviceData), ",")

	fmt.Println("At delete: serviceDataArray:", serviceDataArray)

	newServString := ""

	for _, polID := range serviceDataArray {
		if polID != args[1] {
			newServString += polID + ","
		}
	}

	fmt.Println("First service string:", newServString)

	if len(newServString) > 0 {
		newServString = newServString[0 : len(newServString)-1]
	}

	fmt.Println("Updated service string:", newServString)

	errStoreServ := stub.PutState("serv-"+serviceID+"-"+policyType, []byte(newServString))

	if errStoreServ != nil {
		fmt.Println("Error saving the updated service array at delete!")
	}

	err := stub.DelState(args[1])

	if err != nil {
		return shim.Error("Failed to delete a state")
	}

	err = stub.SetEvent("deletePolicy", []byte(args[1]))

	return shim.Success([]byte("The policy has been successfully deleted!"))
}

func (t *SimpleChaincode) service(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error

	// if len(args) != 2 {
	// 	return shim.Error("Incorrect number of arguments. Expecting 2, function followed by an index")
	// }

	serviceID := args[1]
	policyType := args[2]

	fmt.Println("Index at service:", args[1])

	returnLog, err := stub.GetState("serv-" + serviceID + "-" + policyType)

	if err != nil {
		return shim.Error("Failed to get state for service")
	}
	if returnLog == nil {
		return shim.Error("Service id not found!")
	}

	serviceString := string(returnLog)

	fmt.Println("Service string:", serviceString)

	serviceStringArray := strings.Split(serviceString, ",")

	fmt.Println("Service string array:", serviceStringArray)

	// the following code did not work, convert the data using string with this pattern: polID, pol + polID, pol...
	var policies []PolList

	for _, polID := range serviceStringArray {

		returnPol, errPol := stub.GetState(polID)

		fmt.Println("Pol ID:", polID)

		fmt.Println("Return Pol:", returnPol)

		if errPol != nil {
			return shim.Error("Policy not found!")
		}

		if returnPol == nil {
			continue
		}

		// if errPol == nil {
		// 	return shim.Error("Policy not found!")
		// }
		polArray := strings.Split(string(returnPol), ",")

		fmt.Println("Pol Array:", polArray)

		pol := polArray[4]

		var tempPol PolList

		tempPol.Id = polID
		tempPol.Policy = pol

		policies = append(policies, tempPol)
	}

	fmt.Println("Polices:", policies)

	returnPol, errJSON := json.Marshal(policies)

	fmt.Println("Return pol:", returnPol)

	returnPolString := string(returnPol)

	fmt.Println("Return pol String:", returnPolString)

	if errJSON != nil {
		fmt.Println("Marshalling error:", errJSON)
		return shim.Error("Marshalling error!")
	}

	return shim.Success([]byte(returnPolString))
}

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) savePolicy(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	id := args[1]
	serviceID := args[2]
	expirationTime := args[3]
	policyType := args[4]
	encPolicy := args[5]

	polString := id + "," + serviceID + "," + expirationTime + "," + policyType + "," + encPolicy

	fmt.Println("Being saved again:", (args[1] + "," + args[2] + "," + args[3] + "," + args[4] + "," + args[5]))

	polStringArray := []byte(polString)

	err := stub.PutState(id, polStringArray)

	if err != nil {
		jsonResp := "Error saving the log!"
		fmt.Println("Error storing logs at index:", count)
		fmt.Println("The error is:", err)
		return shim.Error(jsonResp)
	}

	returnService, errServ := stub.GetState("serv-" + serviceID + "-" + policyType)

	if errServ != nil {
		return shim.Error("Failed to get the service state!")
	}

	if returnService == nil {
		errStoreServ := stub.PutState("serv-"+serviceID+"-"+policyType, []byte(id))
		if errStoreServ != nil {
			jsonResp := "Error saving the service array!"
			fmt.Println("Error saving the service array!")
			fmt.Println("The error is:", err)
			return shim.Error(jsonResp)
		}
	} else {
		newArray := string(returnService) + "," + id

		errStoreServ := stub.PutState("serv-"+serviceID+"-"+policyType, []byte(newArray))
		if errStoreServ != nil {
			jsonResp := "Error saving the service array!"
			fmt.Println("Error saving the service array!")
			fmt.Println("The error is:", err)
			return shim.Error(jsonResp)
		}
	}

	err = stub.SetEvent("storePolicy", []byte(id))

	jsonResp := "Log successfully stored!"

	fmt.Println("Log successfully stored!")

	arr := []byte(jsonResp)
	return shim.Success(arr)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
