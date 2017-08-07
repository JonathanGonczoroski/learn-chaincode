package main

import (
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	"encoding/json"
)

func (t *SimpleChaincode) get_volumes(stub shim.ChaincodeStubInterface) ([]byte, error) {
	//select range
	resultsIterator, err := stub.RangeQueryState("order-0", "order-9999999999")

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