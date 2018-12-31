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
// Must version is also available
```

### Configuration client

#### Create configuration client
```
configurationClient, err := fabricClient.CreateConfigurationClient("userName", "orgTitle")
// Must version is also available
```

#### Create channel
```
err = configurationClient.CreateChannel("channelID", "pathToChannelTx")
// Must version is also available
```

#### Join channel
```
err = configurationClient.JoinChannel("channelID")
// Must version is also available
```

#### Install chaincode
```
err = configurationClient.InstallChaincode("chaincodeID", "chaincodePath", "chaincodeVersion")
// Must version is also available
```

#### Instanciate chaincode
```
err = configurationClient.InstanciateChaincode(channelID, chaincodeID, chaincodePath, "chaincodeVersion", [][]byte{[]byte("instantiate"), []byte("args")}, "chaincodePolicy")
// Must version is also available
```

### User client

#### Create user client
```
userClient, err := fabricClient.CreateUserClient("userName", "orgTitle", "channelID")
// Must version is also available
```

#### Invoke transaction
```
response, err = userClient.Invoke("chaincodeID", "chaincodeMethod", [][]byte{[]byte("method"), []byte("args")})
// Must version is also available
```

#### Query transaction (transaction won't be recorded to blockchain)
```
response, err = userClient.Invoke("chaincodeID", "chaincodeMethod", [][]byte{[]byte("method"), []byte("args")})
// Must version is also available
```