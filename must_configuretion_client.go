package fabclient

func (c *ConfigurationClient) MustInstallChaincode(chaincodeID string, chaincodePath string, version string) {
	err := c.InstallChaincode(chaincodeID, chaincodePath, version)
	if err != nil {
		panic(err)
	}
}

func (c *ConfigurationClient) MustInstanciateChaincode(channelID string, chaincodeID string, chaincodePath string, version string, args [][]byte, policy string) {
	err := c.InstanciateChaincode(channelID, chaincodeID, chaincodePath, version, args, policy)
	if err != nil {
		panic(err)
	}
}

func (c *ConfigurationClient) MustCreateChannel(channelID string, channelConfigPath string) {
	err := c.CreateChannel(channelID, channelConfigPath)
	if err != nil {
		panic(err)
	}
}

func (c *ConfigurationClient) MustJoinChannel(channelID string) {
	err := c.JoinChannel(channelID)
	if err != nil {
		panic(err)
	}
}

func (c *ConfigurationClient) MustCreateAndJoinChannel(channelID string, channelConfigPath string) {
	err := c.CreateAndJoinChannel(channelID, channelConfigPath)
	if err != nil {
		panic(err)
	}
}
