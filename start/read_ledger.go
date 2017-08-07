package main

import (
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
)

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Debug("query function: ", function)
	fmt.Println(args)

	if function == "get_volume" {
		return t.get_volume(stub)
	}

	return nil, errors.New("Received unknown function invocation " + function)
}

func (t *SimpleChaincode) get_volume(stub shim.ChaincodeStubInterface) ([]byte, error) {
	bytes, err := stub.GetState()

	return bytes, err
}