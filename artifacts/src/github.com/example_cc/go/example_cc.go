package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type PayloadStruct struct {
	RecordId string `json:"recordId"`
	Payload  string `json:"payload"`
}

type PayloadPrivateStruct struct {
	RecordId  string `json:"recordId"`
	OwnerName string `json:"ownerName"`
}

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}

}

// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Init done ")
	return shim.Success(nil)
}

// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)
	action := args[0]
	fmt.Println("invoke action " + action)
	fmt.Println(args)
	if action == "init" { //initialize the chaincode state, used as reset
		return t.Init(stub)
	} else if action == "commit" {
		return t.EventHandler(stub, args)
	} else if action == "commitPrivate" {
		return t.EventHandlerPrivate(stub, args)
	} else if action == "query" {
		return t.Query(stub, args)
	} else if action == "queryPrivate" {
		return t.QueryPrivate(stub, args)
	}

	fmt.Println("invoke did not find func: " + action) //error

	return shim.Error("Received unknown function")
}

// ===== Example: Ad hoc rich query ========================================================
// Only available on state databases that support rich query (e.g. CouchDB)
// =========================================================================================
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	queryString := args[1]

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

// ===== Example: Ad hoc rich query ========================================================
// Only available on state databases that support rich query (e.g. CouchDB)
// =========================================================================================
func (t *SimpleChaincode) QueryPrivate(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	queryString := args[1]

	queryResults, err := getQueryPrivateResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

// =========================================================================================
// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func getQueryPrivateResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	//fmt.Println("GetQueryResultForQueryString() : getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetPrivateDataQueryResult("collectionUser", queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	fmt.Println(resultsIterator)
	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		queryResponseStr := string(queryResponse.Value)
		fmt.Println(queryResponseStr)
		buffer.WriteString(queryResponseStr)
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	//fmt.Println("GetQueryResultForQueryString(): getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

// =========================================================================================
// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	//fmt.Println("GetQueryResultForQueryString() : getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	fmt.Println(resultsIterator)
	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		queryResponseStr := string(queryResponse.Value)
		fmt.Println(queryResponseStr)
		buffer.WriteString(queryResponseStr)
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	//fmt.Println("GetQueryResultForQueryString(): getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

func (t *SimpleChaincode) EventHandler(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	uniqueid := args[1]

	fmt.Printf("Entering Invoke......\n")

	payloadValue := args[2]

	payloadDataEvent := &PayloadStruct{
		uniqueid,
		payloadValue}

	payloadDataEventJSONasBytes, err := json.Marshal(payloadDataEvent)

	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(stub.GetTxID(), payloadDataEventJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) EventHandlerPrivate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	recordId := args[1]
	ownerName := args[2]

	payloadPrivateStruct := &PayloadPrivateStruct{
		recordId,
		ownerName}
	payloadPrivateDataEventJSONasBytes, err := json.Marshal(payloadPrivateStruct)

	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutPrivateData("collectionUser", recordId, payloadPrivateDataEventJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
