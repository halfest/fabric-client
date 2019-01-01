package fabclient

import (
	"fmt"
)

type ChaincodeClient struct {
	ChaincodeID string
	UserClient  *UserClient
}

func (c *ChaincodeClient) Invoke(functionName string, args [][]byte) ([]byte, error) {
	resp, err := c.UserClient.Invoke(c.ChaincodeID, functionName, args)
	if err != nil {
		return nil, fmt.Errorf("Failed to invoke chaincode %s with funactions %s and arguments %v.\n Error: %v", c.ChaincodeID, functionName, args, err)
	}
	logger.Debugf("Response on invoke chaincode: %s\n", resp)
	return resp, nil
}

func (c *ChaincodeClient) Query(functionName string, args [][]byte) ([]byte, error) {
	resp, err := c.UserClient.Query(c.ChaincodeID, functionName, args)
	if err != nil {
		return nil, fmt.Errorf("Failed to query chaincode %s with funactions %s and arguments %v.\n Error: %v", c.ChaincodeID, functionName, args, err)
	}
	logger.Debugf("Response on query chaincode: %s\n", resp)
	return resp, nil
}
