package fabclient

// MustCreateFabricClient is the same as CreateFabricClient but panics in case of error
func MustCreateFabricClient(configPath string, ordererHost string) *FabricClient {
	result, err := CreateFabricClient(configPath, ordererHost)
	if err != nil {
		panic(err)
	}
	return result
}

// MustCreateConfigurationClient is the same as CreateConfigurationClient but panics in case of error
func (c *FabricClient) MustCreateConfigurationClient(name string, organization string) *ConfigurationClient {
	result, err := c.CreateConfigurationClient(name, organization)
	if err != nil {
		panic(err)
	}
	return result
}

// MustCreateUserClient is the same as CreateUserClient but panics in case of error
func (c *FabricClient) MustCreateUserClient(channelID string, name string, organization string) *UserClient {
	result, err := c.CreateUserClient(channelID, name, organization)
	if err != nil {
		panic(err)
	}
	return result
}

// MustCreateChaincodeClient is the same as CreateChaincodeClient but panics in case of error
func (c *FabricClient) MustCreateChaincodeClient(channelID string, chaincodeID string, name string, organization string) *ChaincodeClient {
	result, err := c.CreateChaincodeClient(channelID, chaincodeID, name, organization)
	if err != nil {
		panic(err)
	}
	return result
}
