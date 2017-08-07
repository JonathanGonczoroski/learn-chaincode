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
const CREATED = "created"
const CANCEL = "cancel"
const SUCESS = "success"

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Volume struct {
	TrackId									string `json: trackId`
	Owner									string `json: owner`
	Shipper									string `json: owner`
	Status									string `json: status`
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
	} /*else if function == "shipperToLogisticProvider" {
        return t.shipperToLogisticProvider(stub, args)
    } else if function == "LogisticProviderToCustomer" {
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
	
 	fmt.Println("[Volume]: " + v.TrackId)

	bytes, err = json.Marshal(v)

	err = stub.PutState(v.TrackId, bytes)

	if err != nil { return nil, errors.New("Unable to put the state") }

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

func (t *SimpleChaincode) get_volumes(stub shim.ChaincodeStubInterface) ([]byte, error) {
	//select range
	resultsIterator, err := stub.RangeQueryState("volume-0", "volume-9999999999")

	if err != nil {
		return nil, errors.New("[IP][Query] Unknown error")
	}

	// clientId := args[0]
	// logisticProviderId := args[1]
	// pendingOrder, err := strconv.ParseBool(args[2])
	// findAll, err := strconv.ParseBool(args[3])

	hasResult := false

	defer resultsIterator.Close()

	result := "{\"volumes\": ["
	
	for resultsIterator.HasNext() {
		queryKeyAsStr, queryValAsBytes, err := resultsIterator.Next()

		fmt.Println("[IP][Query] hack: " + queryKeyAsStr)

		if err != nil {
			return nil, errors.New("[IP][Query] Unknown error")
		}

		var volume Volume
		json.Unmarshal(queryValAsBytes, &volume)

		// clientIdOk := clientId == "-1" || volume.ClientId == clientId			 				
		// logisticProviderIdOk := logisticProviderId == "-1" || volume.LogisticProviderId == logisticProviderId
		
		// var findPendingOk bool

		// if pendingOrder {
		// 	findPendingOk = (volume.LogisticProviderFinalShippingCost == 0)
		// } else {
		// 	findPendingOk = (volume.LogisticProviderFinalShippingCost != 0)
		// }

		// fmt.Println("[IP][Query] ClientId: " + clientId + " | LogisticProviderId: " + logisticProviderId + " | PendingOrder: " + strconv.FormatBool(pendingOrder) + " | ClientIdOk: " + strconv.FormatBool(clientIdOk) + " | LogisticProviderIdOk: " + strconv.FormatBool(logisticProviderIdOk) + " | FindPendingOk: " + strconv.FormatBool(findPendingOk) + " | FindAll: " + strconv.FormatBool(findAll))

		result += string(queryValAsBytes) + ","
		hasResult = true
	}

	if hasResult {
		result = result[:len(result)-1] + "]}"
	} else {
		result = result + "]}"
	}

	return []byte(result), nil
}