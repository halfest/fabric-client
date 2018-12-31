# It is wrapper on fabric-sdk-go with different API to interact with.

## Get
```
go get -u github.com/halfest/fabric-client
```

## Usage

### Init fabric client (fabclient.MustCreateFabricClient also available)
```
fabricClient, err := fabclient.CreateFabricClient("config file for fabric-sdk-go", "orderer host")
```

### Configuration client

#### Create configuration client (fabricClient.MustCreateConfigurationClient also available)
```
configurationClient, err := fabricClient.CreateConfigurationClient("userName", "orgTitle")
```

#### Create channel (configurationClient.MustCreateChannel also available)
```
err = configurationClient.CreateChannel("channelID", "pathToChannelTx")
```

#### Join channel (configurationClient.MustJoinChannel also available)
```
err = configurationClient.JoinChannel("channelID")
```

#### Install chaincode (configurationClient.MustInstallChaincode also available)
```
err = configurationClient.InstallChaincode("chaincodeID", "chaincodePath", "chaincodeVersion")
```

#### Instanciate chaincode (configurationClient.MustInstanciateChaincode also available)
```
err = configurationClient.InstanciateChaincode(channelID, chaincodeID, chaincodePath, "chaincodeVersion", [][]byte{[]byte("instantiate"), []byte("args")}, "chaincodePolicy")
```

### User client

#### Create user client (fabricClient.MustCreateUserClient also available)
```
userClient, err := fabricClient.CreateUserClient("userName", "orgTitle", "channelID")
```

#### Invoke transaction (userClient.MustInvoke also available)
```
response, err = userClient.Invoke("chaincodeID", "chaincodeMethod", [][]byte{[]byte("method"), []byte("args")})
```

#### Query transaction (transaction won't be recorded to blockchain) (userClient.MustInvoke also available)
```
response, err = userClient.Invoke("chaincodeID", "chaincodeMethod", [][]byte{[]byte("method"), []byte("args")})
```