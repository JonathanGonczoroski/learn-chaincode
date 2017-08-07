package main

/* import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
) */

import (
	"errors"
	"fmt"
	//"strconv"
	//"strings"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
	//"regexp"
	//"time"
	"crypto/rand"
	//"io"
)

var logger = shim.NewLogger("CLDChaincode")

// Participant
const	SHIPPER      =  		"shipper"
const	LOGISTIC_PROVIDER   =  	"logistic_provider"
const	INSURENCE_COMPANY = 	"insurence_company"

// Status
const CREATED = "CRIADO"
const CANCEL = "cancel"
const SUCESS = "success"
const DISPATCHED = "dispatched"

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Volume struct {
	TrackId									string `json: trackId`
	Owner									string `json: owner`
	Shipper									string `json: owner`
	Status									string `json: status`
	LogisticProvider						string `json: logisticProvider`
}

func main() {
	fmt.Println("Start Contract")

	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args [] string) ([]byte, error) {
    fmt.Println("invoke is running " + function)

    // Handle different functions
    if function == "init" {
        return t.Init(stub, "init", args)
	} else if function == "CreateVolume" {
		return t.CreateVolume(stub, args[0])
	} else if function == "shipperToLogisticProvider" {
        return t.ShipperToLogisticProvider(stub, args)
    } /*else if function == "LogisticProviderToCustomer" {
		return t.LogisticProviderToCustomer(stub, args)
	} else if function == "LogisticProviderToLogisticProvider" {
		return t.LogisticProviderToLogisticProvider(stub, args)
	} else if function == "LogisticProviderToShipper" {
		return t.LogisticProviderToShipper(stub, args)
	} */

    fmt.Println("invoke did not find func: " + function)
	logger.Debug("invoke did not find func: ", function)

    return nil, errors.New("Received unknown function invocation")
}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Debug("query function: ", function)
	fmt.Println(args)

	if function == "get_volumes" {
		return t.get_volumes(stub)
	}

	return nil, errors.New("Received unknown function invocation " + function)
}

// Functions to Write
func (t *SimpleChaincode) CreateVolume(stub shim.ChaincodeStubInterface, shipper string) ([]byte, error) {
	var v Volume

	var err error
	var bytes []byte

	v.TrackId = GenerateRandomString(10)
	v.Owner = "SHIPPER"
	v.Shipper = shipper
	v.Status = CREATED
	
 	fmt.Println("[Volume]: " + v.TrackId)

	bytes, err = json.Marshal(v)

	err = stub.PutState(v.TrackId, bytes)

	if err != nil { return nil, errors.New("Unable to put the state") }

	return nil, nil
}

func (t *SimpleChaincode) ShipperToLogisticProvider(stub shim.ChaincodeStubInterface, args [] string) ([]byte, error) {
	fmt.Println("shipper to logistic provider running");

	var trackId = args[0];
	var logisticProvider = args[1]
	var v Volume

	fmt.Println("shipper to logistic provider: get state");

	bytes, err := stub.GetState(trackId)
	if err != nil { return nil, errors.New("could not possible do getState ") }

	err = json.Unmarshal(bytes, &v)
	if err != nil { return nil, errors.New("could not possible do unmarshal ") }

	fmt.Println("shipper to logistic provider: got state");

	v.Status = DISPATCHED
	v.LogisticProvider = logisticProvider

	bytes, err = json.Marshal(v)
	if err != nil { return nil, errors.New("could not possible do marshal ") }

	err = stub.PutState(trackId, bytes)

	if err != nil { return nil, errors.New("could not possible do putState ") }

	return nil, nil
}

func GenerateRandomString(n int) (string) {
    const letters = "0123456789"
    bytes := GenerateRandomBytes(n)
    
    for i, b := range bytes {
        bytes[i] = letters[b%byte(len(letters))]
    }
	
    return (string(bytes))
}

func GenerateRandomBytes(n int) ([]byte) {
    b := make([]byte, n)
    _, err := rand.Read(b)

    if err != nil {
        return nil
    }

    return b
}