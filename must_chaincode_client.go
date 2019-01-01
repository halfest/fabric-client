package fabclient

// MustCreateChaincodeClient is the same as CreateChaincodeClient but panics in case of error
func MustCreateChaincodeClient(configPath string, ordererHost string, channelID string, chaincodeID string, name string, organization string) *ChaincodeClient {
	result, err := CreateChaincodeClient(configPath, ordererHost, channelID, chaincodeID, name, organization)
	if err != nil {
		panic(err)
	}
	return result
}

// MustInvoke is the same as Invoke but panics in case of error
func (c *ChaincodeClient) MustInvoke(functionName string, args [][]byte) []byte {
	result, err := c.Invoke(functionName, args)
	if err != nil {
		panic(err)
	}
	return result
}

// MustQuery is the same as Query but panics in case of error
func (c *ChaincodeClient) MustQuery(functionName string, args [][]byte) []byte {
	result, err := c.Query(functionName, args)
	if err != nil {
		panic(err)
	}
	return result
}
