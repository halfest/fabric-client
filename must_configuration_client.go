package fabclient

// MustCreateConfigurationClient is the same as CreateConfigurationClient but panics in case of error
func MustCreateConfigurationClient(configPath string, ordererHost string, name string, organization string) *ConfigurationClient {
	result, err := CreateConfigurationClient(configPath, ordererHost, name, organization)
	if err != nil {
		panic(err)
	}
	return result
}

// MustInstallChaincode is the same as InstallChaincode but panics in case of error
func (c *ConfigurationClient) MustInstallChaincode(chaincodeID string, chaincodePath string, version string) {
	err := c.InstallChaincode(chaincodeID, chaincodePath, version)
	if err != nil {
		panic(err)
	}
}

// MustInstanciateChaincode is the same as InstanciateChaincode but panics in case of error
func (c *ConfigurationClient) MustInstanciateChaincode(channelID string, chaincodeID string, chaincodePath string, version string, args [][]byte, policy string) {
	err := c.InstanciateChaincode(channelID, chaincodeID, chaincodePath, version, args, policy)
	if err != nil {
		panic(err)
	}
}

// MustCreateChannel is the same as CreateChannel but panics in case of error
func (c *ConfigurationClient) MustCreateChannel(channelID string, channelConfigPath string) {
	err := c.CreateChannel(channelID, channelConfigPath)
	if err != nil {
		panic(err)
	}
}

// MustJoinChannel is the same as JoinChannel but panics in case of error
func (c *ConfigurationClient) MustJoinChannel(channelID string) {
	err := c.JoinChannel(channelID)
	if err != nil {
		panic(err)
	}
}

// MustCreateAndJoinChannel is the same as CreateAndJoinChannel but panics in case of error
func (c *ConfigurationClient) MustCreateAndJoinChannel(channelID string, channelConfigPath string) {
	err := c.CreateAndJoinChannel(channelID, channelConfigPath)
	if err != nil {
		panic(err)
	}
}
