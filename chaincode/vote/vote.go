package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {
}

type Candidate struct {
	Name string `json:"name"`
	Party string `json:"party"`
	Votes int `json:"votes"`
}

type Voter struct {
	Voted bool `json:"voted"`
}

type Election struct {
	Ended bool `json:"ended"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	function, args := APIstub.GetFunctionAndParameters()
	if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "queryCandidate" {
		return s.queryCandidate(APIstub, args)
	} else if function == "queryAllCandidates" {
		return s.queryAllCandidates(APIstub)
	} else if function == "queryVoter" {
		return s.queryVoter(APIstub,args)
	} else if function == "createCandidate" {
		return s.createCandidate(APIstub, args)
	} else if function == "startElection" {
		return s.startElection(APIstub)
	} else if function == "castVote" {
		return s.castVote(APIstub, args)
	} else if function == "endElection" {
		return s.endElection(APIstub)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryCandidate(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	fmt.Printf("\nFunction: queryCandidate\n")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	fmt.Printf("Args: %s\n", args[0])

	candidateAsBytes, err := APIstub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	candidate := Candidate{}
	json.Unmarshal(candidateAsBytes, &candidate)

	fmt.Printf("{ Name: %s Party: %s }\n",candidate.Name,candidate.Party)

	return shim.Success(candidateAsBytes)
}

func (s *SmartContract) queryVoter(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	fmt.Printf("\nFunction: queryVoter\n")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	fmt.Printf("Args: %s\n", args[0])

	voterAsBytes, err := APIstub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	voter := Voter{}
	json.Unmarshal(voterAsBytes, &voter)

	fmt.Printf("{ Id : %s Voted: %t }\n",args[0],voter.Voted)

	return shim.Success(voterAsBytes)
}

func (s *SmartContract) startElection(APIstub shim.ChaincodeStubInterface) sc.Response {

	fmt.Printf("\nFunction: startElection\n")

	electionAsBytes, _ := json.Marshal(Election{Ended:false})
	err := APIstub.PutState("ELECTION", electionAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	var buffer bytes.Buffer
	buffer.WriteString("Election started")
	fmt.Printf("%s\n", buffer.String())
	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) endElection(APIstub shim.ChaincodeStubInterface) sc.Response {

	fmt.Printf("\nFunction: startElection\n")

	electionAsBytes, err := APIstub.GetState("ELECTION")
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(electionAsBytes) == 0 {
		var buffer bytes.Buffer
		buffer.WriteString("Election has not started yet")
		fmt.Printf("%s\n", buffer.String())
		return shim.Success(buffer.Bytes())
	}
	election := Election{}

	json.Unmarshal(electionAsBytes, &election)
	election.Ended = true
	electionAsBytes, _ = json.Marshal(election)

	err = APIstub.PutState("ELECTION", electionAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	var buffer bytes.Buffer
	buffer.WriteString("Election ended")
	fmt.Printf("%s\n", buffer.String())
	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {

	fmt.Printf("\nFunction: initLedger\n")

	candidates := []Candidate{
		Candidate{Name: "Mark Zuckerberg", Party: "Facebook", Votes: 0},
		Candidate{Name: "Sundar Pitchai", Party: "Google", Votes: 0},
		Candidate{Name: "Satya Nadella", Party: "Microsoft", Votes: 0},
		Candidate{Name: "Tim Cook", Party: "Apple", Votes: 0},
		Candidate{Name: "Jeff Bezos", Party: "Amazon", Votes: 0},
	}

	i := 0
	for i < len(candidates) {
		candidateAsBytes, _ := json.Marshal(candidates[i])
		APIstub.PutState("CANDIDATE"+strconv.Itoa(i), candidateAsBytes)
		fmt.Printf("Added { Name: %s Party: %s }\n",candidates[i].Name,candidates[i].Party)
		i = i + 1
	}

	i = 0
	for i < 10 {
		voter := Voter{Voted:false}
		voterAsBytes, _ := json.Marshal(voter)
		APIstub.PutState("VOTER"+strconv.Itoa(i), voterAsBytes)
		fmt.Printf("Added { Id: VOTER%d }\n",i)
		i = i + 1
	}

	fmt.Println("Ledger initated with values")

	return shim.Success(nil)
}

func (s *SmartContract) createCandidate(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	fmt.Printf("\nFunction: createCandidate\n")

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	fmt.Printf("Args: %s %s %s\n", args[0], args[1], args[2])

	electionAsBytes, err := APIstub.GetState("ELECTION")
	if err != nil {
		return shim.Error(err.Error())
	}
	election := Election{}
	json.Unmarshal(electionAsBytes, &election)
	if election.Ended {
		var buffer bytes.Buffer
		buffer.WriteString("Election ended")
		fmt.Printf("%s\n", buffer.String())
		return shim.Success(buffer.Bytes())
	}

	var candidate = Candidate{Name: args[1], Party: args[2], Votes: 0}

	candidateAsBytes, _ := json.Marshal(candidate)
	APIstub.PutState(args[0], candidateAsBytes)

	fmt.Printf("{ Name: %s Party: %s }\n",candidate.Name,candidate.Party)

	var buffer bytes.Buffer
	buffer.WriteString("Successfully created candidate")
	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) queryAllCandidates(APIstub shim.ChaincodeStubInterface) sc.Response {

	fmt.Printf("\nFunction: queryAllCandidates\n")

	startKey := "CANDIDATE0"
	endKey := "CANDIDATE999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) castVote(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	fmt.Printf("\nFunction: castVote\n")
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	fmt.Printf("Args: %s %s\n", args[0], args[1])

	electionAsBytes, err := APIstub.GetState("ELECTION")
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(electionAsBytes) == 0 {
		var buffer bytes.Buffer
		buffer.WriteString("Election has not started yet")
		fmt.Printf("%s\n", buffer.String())
		return shim.Success(buffer.Bytes())
	}
	election := Election{}
	json.Unmarshal(electionAsBytes, &election)
	if election.Ended {
		var buffer bytes.Buffer
		buffer.WriteString("Election ended")
		fmt.Printf("%s\n", buffer.String())
		return shim.Success(buffer.Bytes())
	}

	voterAsBytes, err := APIstub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	if len(voterAsBytes) == 0 {
		var buffer bytes.Buffer
		buffer.WriteString("Invalid Voter ID")
		fmt.Printf("%s\n", buffer.String())
		return shim.Success(buffer.Bytes())
	}

	voter := Voter{}
	json.Unmarshal(voterAsBytes, &voter)
	if voter.Voted {
		var buffer bytes.Buffer
		buffer.WriteString("You have already voted")
		fmt.Printf("%s\n", buffer.String())
		return shim.Success(buffer.Bytes())
	}

	candidateAsBytes, err := APIstub.GetState(args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(candidateAsBytes) == 0 {
		var buffer bytes.Buffer
		buffer.WriteString("Invalid Candidate ID")
		fmt.Printf("%s\n", buffer.String())
		return shim.Success(buffer.Bytes())
	}

	voter.Voted=true
	newVoterAsBytes, _ := json.Marshal(voter)
	err = APIstub.PutState(args[0], newVoterAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	candidate := Candidate{}

	json.Unmarshal(candidateAsBytes, &candidate)
	candidate.Votes =  candidate.Votes + 1

	newCandidateAsBytes, _ := json.Marshal(candidate)
	err = APIstub.PutState(args[1], newCandidateAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	var buffer bytes.Buffer
	buffer.WriteString("Successfully voted for "+args[1])
	return shim.Success(buffer.Bytes())
}

func main() {

	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s\n", err)
	}
}
