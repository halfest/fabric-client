package fabclient

func (c *UserClient) MustInvoke(chaincodeID string, functionName string, args [][]byte) []byte {
	result, err := c.Invoke(chaincodeID, functionName, args)
	if err != nil {
		panic(err)
	}
	return result
}

func (c *UserClient) MustQuery(chaincodeID string, functionName string, args [][]byte) []byte {
	result, err := c.Query(chaincodeID, functionName, args)
	if err != nil {
		panic(err)
	}
	return result
}
