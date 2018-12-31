# fabric-client
It is wrapper on fabric-sdk-go with different API to interact with.

## Get
```
go get -u github.com/halfest/fabric-client
```

## Usage

### Init fabric client
```
fabricClient, err := fabclient.CreateFabricClient("config file for fabric-sdk-go", "orderer host")
// fabclient.MustCreateFabricClient also available
```

### Configuration client

#### Create configuration client
```
configurationClient, err := fabricClient.CreateConfigurationClient("userName", "orgTitle")
// fabricClient.MustCreateConfigurationClient also available
```

#### Create channel
```
err = configurationClient.CreateChannel("channelID", "pathToChannelTx")
// configurationClient.MustCreateChannel also available
```

#### Join channel
```
err = configurationClient.JoinChannel("channelID")
// configurationClient.MustJoinChannel also available
```

#### Install chaincode
```
err = configurationClient.InstallChaincode("chaincodeID", "chaincodePath", "chaincodeVersion")
// configurationClient.MustInstallChaincode also available
```

#### Instanciate chaincode
```
err = configurationClient.InstanciateChaincode(channelID, chaincodeID, chaincodePath, "chaincodeVersion", [][]byte{[]byte("instantiate"), []byte("args")}, "chaincodePolicy")
// configurationClient.MustInstanciateChaincode also available
```

### User client

#### Create user client
```
userClient, err := fabricClient.CreateUserClient("userName", "orgTitle", "channelID")
// fabricClient.MustCreateUserClient also available
```

#### Invoke transaction
```
response, err = userClient.Invoke("chaincodeID", "chaincodeMethod", [][]byte{[]byte("method"), []byte("args")})
// userClient.MustInvoke also available
```

#### Query transaction (transaction won't be recorded to blockchain)
```
response, err = userClient.Invoke("chaincodeID", "chaincodeMethod", [][]byte{[]byte("method"), []byte("args")})
// userClient.MustInvoke also available
```