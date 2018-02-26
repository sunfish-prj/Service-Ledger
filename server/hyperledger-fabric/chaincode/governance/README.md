# Governance Chain-code

This section contains "governance" chaincodes for the SUNFISH federation:
The chain-code support four different function in order to achieve a voting session. The inputs that can thai takes are: 

- requestorID : represent the ID of the member that submit a proposal
- proposalID : it is the ID of the proposal
- proposalDescription : it is the description of the proposal
- proposalType : this field represent the proposal typology (e.g. join, leave, update)
- proposalQuorum : the typology of the quorum is needed to validate the proposal
- proposalStatus : this is the current status of the proposal (e.g. pending, accepted, rejected)
- votersNumber : this is the number of voters needed to make the proposal validable


## submitProposal

The function is used for submitting of new proposal to vote. It takes as input "requestorID, proposalID, proposalDescription, proposalType, proposalQuorum, proposalStatus, votersNumber"

An example of the invoke for it is :

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

The function returns detail about the corresponding proposal taken as input "proposalID".

An example instance of the invoke for it is :

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

The function to vote a proposal. It takes as input "memberID, proposalID, vote"

An example instance of the invoke for it is :

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

The function starts the validation process of a proposal. It takes as input "proposalID".

An example of the invoke for it is :

    {
    "channel": "sunfish-channel",
    "peer": "peer01",
    "chaincodeName": "governance.go",
    "fcn": "countVote",
    "args": [
        "proposal01"
        ]
    } 


