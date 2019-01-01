package fabclient

// MustCreateUserClient is the same as CreateUserClient but panics in case of error
func MustCreateUserClient(configPath string, ordererHost string, channelID string, name string, organization string) *UserClient {
	result, err := CreateUserClient(configPath, ordererHost, channelID, name, organization)
	if err != nil {
		panic(err)
	}
	return result
}

// MustInvoke is the same as Invoke but panics in case of error
func (c *UserClient) MustInvoke(chaincodeID string, functionName string, args [][]byte) []byte {
	result, err := c.Invoke(chaincodeID, functionName, args)
	if err != nil {
		panic(err)
	}
	return result
}

// MustQuery is the same as Query but panics in case of error
func (c *UserClient) MustQuery(chaincodeID string, functionName string, args [][]byte) []byte {
	result, err := c.Query(chaincodeID, functionName, args)
	if err != nil {
		panic(err)
	}
	return result
}
