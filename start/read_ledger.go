package main

import (
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    logger.Debug("query function: ", function)

	if function == "get_volume" {
		return t.get_volume(stub, args[0])
	}

	return nil, errors.New("Received unknown function invocation " + function)
}

func (t *SimpleChaincode) get_volume(stub shim.ChaincodeStubInterface, tracker_id string) ([]byte, error) {
	bytes, err := stub.GetState(tracker_id)

	return bytes, err
}