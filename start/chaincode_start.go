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
const DISPATCHED = "DESPACHADO"

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Volume struct {
	TrackId									string `json: trackId`
	Owner									string `json: owner`
	Shipper									string `json: owner`
	Status									string `json: status`
	LogisticProvider						string `json: logisticProvider`
	PropertiesLogisticProvider				string `json: propertisLogisticProvider`
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
    } else if function == "LogisticProviderToLogisticProvider" {
		return t.LogisticProviderToLogisticProvider(stub, args)
	}

    fmt.Println("invoke did not find func: " + function)
	logger.Debug("invoke did not find func: ", function)

    return nil, errors.New("Received unknown function invocation")
}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Debug("query function: ", function)
	fmt.Println(args)

	if function == "get_volumes" {
		return t.get_volumes(stub, args)
	} else if function == "get_tracker" {
		return t.get_tracker(stub, args[0])
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

	fmt.Println("[debug] - trackId: " + trackId);
	var logisticProvider = args[1]
	fmt.Println("[debug] - logisticProvider: " + logisticProvider);
	var v Volume

	fmt.Println("shipper to logistic provider: get state");

	bytes, err := stub.GetState(trackId)
	fmt.Println("[debug] - bytes ");

	if err != nil { return nil, errors.New("could not possible do getState ") }

	err = json.Unmarshal(bytes, &v)
	fmt.Println("[debug] - v before atribute: " + v.TrackId);

	if err != nil { return nil, errors.New("could not possible do unmarshal ") }

	fmt.Println("shipper to logistic provider: got state");

	v.Status = DISPATCHED
	v.LogisticProvider = logisticProvider

	fmt.Println("[debug] - v after attribute: " + v.Status);

	bytes, err = json.Marshal(v)
	if err != nil { return nil, errors.New("could not possible do marshal ") }
	fmt.Println("[debug] - bytes");

	err = stub.PutState(trackId, bytes)

	if err != nil { return nil, errors.New("could not possible do putState ") }

	return nil, nil
}

func (t *SimpleChaincode) LogisticProviderToLogisticProvider(stub shim.ChaincodeStubInterface, args [] string) ([]byte, error) {
	fmt.Println("logistic provider to logistic provider running");

	var trackId = args[0]
	var status = args[1]
	var properties = args[2]

	var v Volume

	fmt.Println("logistic provider to logistic provider: get state");

	bytes, err := stub.GetState(trackId)
	if err != nil { return nil, errors.New("could not possible do getState ") }

	err = json.Unmarshal(bytes, &v)
	if err != nil { return nil, errors.New("could not possible do unmarshal ") }

	fmt.Println("shipper to logistic provider: got state");

	v.Status = status
	v.PropertiesLogisticProvider = properties

	bytes, err = json.Marshal(v)
	if err != nil { return nil, errors.New("could not possible do marshall ") }

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

func (t *SimpleChaincode) get_volumes(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//select range
	resultsIterator, err := stub.RangeQueryState("0", "9999999999")

	if err != nil {
		return nil, errors.New("[IP][Query] Unknown error")
	}	

	field := args[0]
	value := args[1]
	fmt.Println(field)
	fmt.Println(value)
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
		fmt.Println(string(queryValAsBytes));

		if field == "trackerId" {
			if volume.TrackId == value {
				result += string(queryValAsBytes) + ","
				hasResult = true
			}
		} else if field == "shipper" {
			if volume.Shipper == value && volume.Status == "CRIADO" {
				result += string(queryValAsBytes) + ","
				hasResult = true
			}
		} else if field == "logistic_provider" {
			if volume.LogisticProvider == value {
				result += string(queryValAsBytes) + ","
				hasResult = true
			}
		} else {
			result += string(queryValAsBytes) + ","
			hasResult = true
		}

		// clientIdOk := clientId == "-1" || volume.ClientId == clientId			 				
		// logisticProviderIdOk := logisticProviderId == "-1" || volume.LogisticProviderId == logisticProviderId
		
		// var findPendingOk bool

		// if pendingOrder {
		// 	findPendingOk = (volume.LogisticProviderFinalShippingCost == 0)
		// } else {
		// 	findPendingOk = (volume.LogisticProviderFinalShippingCost != 0)
		// }

		// fmt.Println("[IP][Query] ClientId: " + clientId + " | LogisticProviderId: " + logisticProviderId + " | PendingOrder: " + strconv.FormatBool(pendingOrder) + " | ClientIdOk: " + strconv.FormatBool(clientIdOk) + " | LogisticProviderIdOk: " + strconv.FormatBool(logisticProviderIdOk) + " | FindPendingOk: " + strconv.FormatBool(findPendingOk) + " | FindAll: " + strconv.FormatBool(findAll))
	}

	if hasResult {
		result = result[:len(result)-1] + "]}"
	} else {
		result = result + "]}"
	}

	return []byte(result), nil
}

func (t *SimpleChaincode) get_volume(stub shim.ChaincodeStubInterface, trackerId string) (Volume, error) {
	//select range
	resultsIterator, err := stub.RangeQueryState("0", "9999999999")
    var volume Volume
	if err != nil {
		return volume, errors.New("[IP][Query] Unknown error")
	}		

	defer resultsIterator.Close()
	
	for resultsIterator.HasNext() {
		queryKeyAsStr, queryValAsBytes, err := resultsIterator.Next()

		fmt.Println("[IP][Query] hack: " + queryKeyAsStr)

		if err != nil {
			return volume, errors.New("[IP][Query] Unknown error")
		}


		json.Unmarshal(queryValAsBytes, &volume)
		fmt.Println(string(queryValAsBytes));

		return volume, nil
	}

	return volume, nil
}

func (t *SimpleChaincode) get_tracker(stub shim.ChaincodeStubInterface, trackerId string) ([]byte, error) {
	tracker, err := stub.GetState(trackerId)

	if err != nil { return nil, errors.New("Couldn't retrieve tracker for ID " + trackerId) }

	return tracker, nil
}