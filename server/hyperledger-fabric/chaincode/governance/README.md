# Governance Chain-code

This section contains "governance" chaincodes for the SUNFISH federation:
The chain-code support four different function in order to achieve a voting session. The function inputs are as follows: 

- requestorID : it represents the ID of the member who submits a proposal
- proposalID : it is the ID of the proposal
- proposalDescription : it is the description of the proposal
- proposalType : it represents the proposal typology (one among *join, leave, update*)
- proposalQuorum : the type of quorum needed to validate the proposal (i.e. ...)
- proposalStatus : it is the current status of the proposal (one among *pending, accepted, rejected*)
- votersNumber : it is the number of voters needed to make the proposal validable


## submitProposal

The function is used for submitting a new proposal to vote. It takes in input "requestorID, proposalID, proposalDescription, proposalType, proposalQuorum, proposalStatus, votersNumber"

By reluing on the Service Ledger *invoke* API, the correponding invocation is as follows:

    {
    "channel": "sunfish-channel",
    "peer": "peer01",
    "chaincodeName": "governance.go",
    "fcn": "submitProposal",
    "args": [
        "member01",
        "proposal01",
        "this is a test proposal description",
        "join",
        "unanimity",
        "pending",
        "10"
        ]
    } 

## getProposal

The function returns detail about the corresponding proposal taken as input "proposalID" (which must corresponds to an existing proposal)

The corresponding *invoke* call is as follows
    {
    "channel": "sunfish-channel",
    "peer": "peer01",
    "chaincodeName": "governance.go",
    "fcn": "getProposal",
    "args": [
        "proposal01"
        ]
    } 

## vote

The function to vote---"accept" or "reject"---a proposal. It takes as input "memberID, proposalID, vote"

The corresponding *invoke* call is as follows
    {
    "channel": "sunfish-channel",
    "peer": "peer01",
    "chaincodeName": "governance.go",
    "fcn": "vote",
    "args": [
        "member01",
        "proposal01",
        "accept"
        ]
    } 

## countVote

The function triggers the validation process for a proposal. It takes as input "proposalID".

The corresponding *invoke* call is as follows
    {
    "channel": "sunfish-channel",
    "peer": "peer01",
    "chaincodeName": "governance.go",
    "fcn": "countVote",
    "args": [
        "proposal01"
        ]
    } 


