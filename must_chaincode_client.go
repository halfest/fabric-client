package fabclient

func (c *ChaincodeClient) MustInvoke(functionName string, args [][]byte) []byte {
	result, err := c.Invoke(functionName, args)
	if err != nil {
		panic(err)
	}
	return result
}

func (c *ChaincodeClient) MustQuery(functionName string, args [][]byte) []byte {
	result, err := c.Query(functionName, args)
	if err != nil {
		panic(err)
	}
	return result
}
