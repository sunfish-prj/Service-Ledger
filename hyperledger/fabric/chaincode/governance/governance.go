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

// This chaincode allow to define a database to store
// contract proposal submitted by members and to vote for them

// "Proposal" the following structures need to support JSON building
type Proposal struct {
	// Requestor           *Requestor `json:"votes"`
	Requestor           string `json:"requestor"`
	ProposalID          string `json:"ID"`
	ProposalType        string `json:"type"`
	ProposalDescription string `json:"description"`
	ProposalQuorum      string `json:"quorum"`
	ProposalVoters      string `json:"votersNumber"`
	ProposalStatus      string `json:"status"`
}

// type Requestor struct {
// 	RequestorID          string `json:"ID"`
// 	RequestorDescription string `json:"description"`
// }

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("proposal chaincode Init")
	_, args := stub.GetFunctionAndParameters()

	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Expecting 0")
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("[E-VOTING CHAINCODE] e_voting chaincode Invoke")
	function, args := stub.GetFunctionAndParameters()

	// can be defined a new function here to support the vote
	if function == "submitProposal" {
		// put a new proposal in the bc
		return t.submitProposal(stub, args)
	} else if function == "getProposal" {
		// get the proposal stored with a given ID
		return t.getProposal(stub, args)
	} else if function == "vote" {
		// members can vote a proposal
		return t.vote(stub, args)
	} else if function == "countVote" {
		// start the votes counting to decide the proposal
		return t.countVote(stub, args)
	}
	return shim.Error("[E-VOTING CHAINCODE] Invalid invoke function name. Expecting \"submitProposal\" \"getProposal\" \"vote\" \"countVote\" ")
}

/* The function gives the possibility to store a proposal in the data structure
and are needed requestorID, propID, propInfo, propType, propQuorum, propVoters
to call it in a proper way*/
func (t *SimpleChaincode) submitProposal(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var requestorID, propID, propInfo, propType, propQuorum, propVoters string
	var propStatus string
	var err error

	// types of quorum that are accepted
	acceptedQuorum := map[string]bool{"unanimity": true, "majority": true, "oneThird": true}

	if len(args) != 6 {
		return shim.Error("[E-VOTING CHAINCODE][SubmitProposal] Incorrect number of arguments.\n	Expecting 6: requestorID, propID, propInfo, propType, propQuorum, propVoters")
	}

	// Setting chaincode args
	requestorID = args[0]
	propID = args[1]
	propInfo = args[2]
	propType = args[3]
	propQuorum = args[4]
	propVoters = args[5]
	propStatus = "pending"

	// check if the quorum is one of the supported
	if !acceptedQuorum[propQuorum] {
		return shim.Error("[E-VOTING CHAINCODE][SubmitProposal] Invalid quorum selected. Expecting \"unanimity\" \"majority\" \"oneThird\" ")
	}
	// create JSON structure of the proposal
	proposal := &Proposal{}
	proposal.Requestor = requestorID
	proposal.ProposalID = propID
	proposal.ProposalDescription = propInfo
	proposal.ProposalType = propType
	proposal.ProposalQuorum = propQuorum
	proposal.ProposalVoters = propVoters
	proposal.ProposalStatus = propStatus
	proposalJSON, _ := json.Marshal(proposal)

	// Check if a proposal with the same ID already exists TODO controllare il check
	propJSONasBytes, err := stub.GetState(propID)
	if propJSONasBytes != nil {
		return shim.Error("[E-VOTING CHAINCODE][SubmitProposal] A proposal ID" + propID + "already exists")
	}
	fmt.Printf("[E-VOTING CHAINCODE][SubmitProposal] Executing Put: key = %s, value = %s\n", propID, string(proposalJSON))

	// Write the state to the ledger
	err = stub.PutState(propID, proposalJSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("\033[1;36m[E-VOTING CHAINCODE][SubmitProposal] Proposal submission executed!\033[0m\n")
	return shim.Success(nil)
}

// ############################################################################################
// ############################################################################################

// Get a proposal description from a given proposal ID
func (t *SimpleChaincode) getProposal(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var propID string
	var err error

	if len(args) != 1 {
		return shim.Error("[E-VOTING CHAINCODE][GetProposal] Incorrect number of arguments. Expecting name of the person to query")
	}

	propID = args[0]

	fmt.Printf("[E-VOTING CHAINCODE][GetProposal] Executing Get: key = %s\n", propID)

	// Get the state from the ledger
	propJSONasBytes, err := stub.GetState(propID)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + propID + "\"}"
		return shim.Error(jsonResp)
	}

	// check if the get returns something
	if propJSONasBytes != nil {
		jsonResp := "{\"key\":\"" + propID + "\",\"value\":\"" + string(propJSONasBytes) + "\"}"
		fmt.Printf("\033[1;36m[E-VOTING CHAINCODE][GetProposal] Get proposal executed! The JSON of proposal %s is: %s\033[0m\n", propID, jsonResp)
		return shim.Success(propJSONasBytes)
	}
	fmt.Printf("\033[1;36m[E-VOTING CHAINCODE][GetProposal] Get proposal does not return anything!\033[0m\n")
	jsonResp := "\"Error\":\"Get proposal does not return anything!\""
	return shim.Error(jsonResp)

}

// ############################################################################################
// ############################################################################################

/* Write the vote of a member state to the ledger in a different way is needed.
There is as a 'pre-storing' phase to avoid probable concurrency issues
The votes are stored with CompositeKeys that are represented with the index
proposalID~voterID and they can be stored in the same data structure used for
storing Proposal. In that way it is possible to votes avoiding cuncurrency issues
that can be (each voter of given proposal has its slot to send the vote). */
func (t *SimpleChaincode) vote(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var voterID, propID, vote string
	var err error
	// map of valid votes
	acceptedVotValues := map[string]bool{"accept": true, "reject": true}

	if len(args) != 3 {
		jsonResp := "{\"Error\":\"Incorrect number of arguments. Expecting 3\"}"
		return shim.Error(jsonResp)
	}

	// Setting chaincode args
	propID = args[0]
	voterID = args[1]
	vote = args[2]

	// create JSON structure
	proposal := &Proposal{}

	// open the proposal and decompose the structure tu read and update its fields
	propJSONasBytes, err := stub.GetState(propID)
	// check if the get returns something
	if propJSONasBytes == nil || err != nil {
		jsonResp := "{\"Error\":\"Problem with the get of the proposal\"}"
		return shim.Error(jsonResp)
	}

	// to process and update data we need to unmarshal the JSON file
	err = json.Unmarshal(propJSONasBytes, proposal)
	if err != nil {
		jsonResp := "{\"Error\":\"Unmarshal process error\"}"
		return shim.Error(jsonResp)
	}
	fmt.Printf("\033[1;36m[E-VOTING CHAINCODE][Vote] JSON value of the proposal %s Unmarshal done.\033[0m\n", propID)

	if proposal.Requestor == voterID {
		jsonResp := "{\"Error\":\"The proposal can not be voted by its requestor}"
		return shim.Error(jsonResp)
	}

	//check vote value
	if acceptedVotValues[vote] {
		fmt.Println("[E-VOTING CHAINCODE][Vote] Value Proposed as vote is correct")
	} else {
		jsonResp := "\"Error\":\"Incorrect vote value. Expecting \"accept\", \"reject\"}"
		return shim.Error(jsonResp)
	}
	fmt.Printf("[E-VOTING CHAINCODE][Vote] Storing Vote of PROPID = %s, VOTERID = %s\n", propID, voterID)

	// create a composite key to store as a "cupled" key propID-voterId
	keyVote, err := stub.CreateCompositeKey("proposal~voter", []string{propID, voterID})
	if err != nil {
		jsonResp := "\033[1;31m{\"Error\":\"Failed to create the composite key\"}\033[0m\n"
		return shim.Error(jsonResp)
	}

	// Takes the iterator of all the compositekeys contains propID
	keyResultsIterator, err := stub.GetStateByPartialCompositeKey("proposal~voter", []string{propID})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer keyResultsIterator.Close()
	// This is needed to check if a member has already voted

	// Iterate through result set and for each prop found, to check the positive votes
	var i int
	for i = 0; keyResultsIterator.HasNext(); i++ {
		// Note that we don't get the value (2nd return variable), we'll just get the voterID from the composite key
		currentCompositeKey, err := keyResultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		// get the color and name from color~name composite key
		objectType, compositeKeyParts, err := stub.SplitCompositeKey(currentCompositeKey.Key)
		if err != nil {
			return shim.Error(err.Error())
		}
		// returnedProposal := compositeKeyParts[0]
		returnedVoter := compositeKeyParts[1]
		if returnedVeer == voterID {
			// If the member has been found ghet the value alreay voted and send an error
			fmt.Print("\033[1;31m[E-VOTING CHAINCODE][Vote] The index " + string(objectType) + " already contains the voter" + string(returnedVoter + "\033[0m\n")
			jsonResp := "{\"Error\":\"the member " + string(voterID) + " has already voted.\"}"
			return shim.Error(jsonResp)
		}
	}

	err = stub.PutState(keyVote, []byte(vote))
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("[E-VOTING CHAINCODE][Vote] Vote stored\n")
	return shim.Success(nil)
}

// #############################################################################################
// #############################################################################################

/* The countVote fuction takes as input the proposal ID to validate;
It interrogate the data structure to rertieve the information of
the proposal to validate and the composite keys that has as prefix
it (proposalID~voterID). The composite keys contains the votes to count
to validate the proposal as accepted or rejected. */
func (t *SimpleChaincode) countVote(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var propID string
	var acceptCounter, rejectCounter int
	var err error
	//var votes, keySet []string //this is the arey of all the key of the voter of a proposal
	// map of valid votes
	acceptedVotValues := map[string]bool{"accept": true, "reject": true}

	if len(args) != 1 {
		jsonResp := "{\"Error\":\"Incorrect number of arguments. Expecting 1: the propID to count\"}"
		return shim.Error(jsonResp)
	}

	// Setting chaincode args
	propID = args[0]
	fmt.Printf("[E-VOTING CHAINCODE][CountVote] Proposal to check: %s\n", propID)

	// create JSON structure
	proposal := &Proposal{}

	// open the proposal and decompose the structure tu read and update its fields
	propJSONasBytes, err := stub.GetState(propID)
	// check if the get returns something
	if propJSONasBytes == nil || err != nil {
		jsonResp := "{\"Error\":\"Problem with the get of the proposal\"}"
		return shim.Error(jsonResp)
	}

	// to process and update data we need to unmarshal the JSON file
	err = json.Unmarshal(propJSONasBytes, proposal)
	if err != nil {
		jsonResp := "{\"Error\":\"Unmarshal process error\"}"
		return shim.Error(jsonResp)
	}
	fmt.Printf("\033[1;36m[E-VOTING CHAINCODE][CountVote] JSON value of the proposal %s Unmarshal done.\033[0m\n", propID)

	if proposal.ProposalStatus != "pending" {
		jsonResp := "{\"Error\":\"Proposal already voted.\n The proposal " + propID + " has been validated as " + proposal.ProposalStatus + " \"}"
		return shim.Error(jsonResp)
	}

	//initialize counter for collect votes
	acceptCounter = 0
	rejectCounter = 0
	// takes the iterator of all the compositekeys contains propID
	keyResultsIterator, err := stub.GetStateByPartialCompositeKey("proposal~voter", []string{propID})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer keyResultsIterator.Close()

	// get the proposal id and voter id from proposal-voter composite key
	// #####################################################################################
	// Iterate through result set and for each prop found, to check the positive votes
	var i int
	for i = 0; keyResultsIterator.HasNext(); i++ {
		// Note that we don't get the value (2nd return variable), we'll just get the voterID from the composite key
		currentCompositeKey, err := keyResultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		// get the propsal and voter from proposal~voter composite key
		objectType, compositeKeyParts, err := stub.SplitCompositeKey(currentCompositeKey.Key)
		if err != nil {
			return shim.Error(err.Error())
		}
		returnedProposal := compositeKeyParts[0]
		returnedVoter := compositeKeyParts[1]
		fmt.Printf("[E-VOTING CHAINCODE][CountVote] found a vote from index:%s Proposal:%s Voter:%s\n", objectType, returnedProposal, returnedVoter)

		// get the value using the composite key
		value := currentCompositeKey.Value
		if err != nil {
			jsonResp := "{\"Error\":\"Failed to get state\"}"
			return shim.Error(jsonResp)
		}

		// in this section the counter are updated on the basis of the given keys' values
		if acceptedVotValues[string(value)] {
			if string(value) == "accept" {
				acceptCounter++
			} else if string(value) == "reject" {
				rejectCounter++
			}
		}
	}

	// having the counters we can update the status of the proposal basing on the quorum to reach
	// if the the quorum is reached, the proposal is accepted; rejected otherwise
	voterInt, err := strconv.Atoi(proposal.ProposalVoters)
	switch {
	case proposal.ProposalQuorum == "unanimity":
		if acceptCounter+rejectCounter >= voterInt {
			if acceptCounter >= voterInt {
				proposal.ProposalStatus = "accepted"
			} else {
				proposal.ProposalStatus = "rejected"
			}
		} else {
			fmt.Printf("\033[1;36m[E-VOTING CHAINCODE][CountVote] Some votes left!\033[0m\n")
			return shim.Success(nil)
		}
	case proposal.ProposalQuorum == "majority":
		if acceptCounter+rejectCounter >= voterInt {
			if acceptCounter > rejectCounter {
				proposal.ProposalStatus = "accepted"
			} else {
				proposal.ProposalStatus = "rejected"
			}
		} else {
			fmt.Printf("\033[1;36m[E-VOTING CHAINCODE][CountVote] Some votes left!\033[0m\n")
			return shim.Success(nil)
		}
	case proposal.ProposalQuorum == "oneThird":
		if acceptCounter+rejectCounter >= voterInt {
			if acceptCounter > ((acceptCounter + rejectCounter) / 3) {
				proposal.ProposalStatus = "accepted"
			} else {
				proposal.ProposalStatus = "rejected"
			}
		} else {
			fmt.Printf("\033[1;36m[E-VOTING CHAINCODE][CountVote] Some votes left!\033[0m\n")
			return shim.Success(nil)
		}
	}
	fmt.Printf("[E-VOTING CHAINCODE][CountVote] Counted \"accept\" votes: %d \n", acceptCounter)
	fmt.Printf("[E-VOTING CHAINCODE][CountVote] Counted \"reject\" votes: %d \n", rejectCounter)

	// after the update of the status we need to marshal the new json and save it
	proposalJSON, _ := json.Marshal(proposal)
	if err != nil {
		jsonResp := "{\"Error\": \"Something goes wrong during Marshal process\"}"
		return shim.Error(jsonResp)
	}

	// Write the state to the ledger
	err = stub.PutState(propID, proposalJSON)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("\033[1;36m[E-VOTING CHAINCODE][CountVote] Proposal validation executed!\033[0m\n")
	return shim.Success(nil)

	// #####################################################################################

}

// STRINGED VERSION: { "requestorInfo" : {"requestorID":"string", "requestorDescription":"string"},	"proposalType":"string", "proposalQuorum":"string", "proposalInfo":{"proposalDescription":"string","proposalID":"string"}}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("[E-VOTING CHAINCODE] Error starting Simple chaincode: %s", err)
	}
}
