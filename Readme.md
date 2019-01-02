# fabric-client
It is wrapper on fabric-sdk-go with different API to interact with.

It consists of 3 clients:
1. Configuration client used for creating, joining channels, installing, instantiating chaincodes.
2. User client used for invoking and querying transactions.
3. Channel client. It is similar to user client but for one chaincode only.

## Get
```go
go get -u github.com/halfest/fabric-client
```

## Features
There is difference in packing chaincode during it's installation between fabric-sdk-go and fabric peer cli. Former does not pack dependencies alongside with chaincode, so function 
```go
// Generates a deployment payload for GOLANG as a series of src/$pkg entries in .tar.gz format
func (goPlatform *Platform) GetDeploymentPayload(path string) ([]byte, error) {
```
from fabric peer cli was used.

## Usage

### Init fabric client
```go
import (
	fabclient "github.com/halfest/fabric-client"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/logging"
)

logging.SetLevel("fabclient", logging.DEBUG)
fabricClient, err := fabclient.CreateFabricClient("config file for fabric-sdk-go", "orderer host")
// Must version is also available
```

### Configuration client

#### Create configuration client
```go
configurationClient, err := fabricClient.CreateConfigurationClient("userName", "orgTitle")
// or without reuse of fabric client
configurationClient, err := fabclient.CreateConfigurationClient("config file for fabric-sdk-go", "orderer host", "userName", "orgTitle")
// Must versions is also available
```

#### Create channel
```go
err = configurationClient.CreateChannel("channelID", "pathToChannelTx")
// Must version is also available
```

#### Join channel
```go
err = configurationClient.JoinChannel("channelID")
// Must version is also available
```

#### Install chaincode
```go
err = configurationClient.InstallChaincode("chaincodeID", "chaincodePath", "chaincodeVersion")
// Must version is also available
```

#### Instanciate chaincode
```go
err = configurationClient.InstanciateChaincode("channelID", "chaincodeID", "chaincodePath", "chaincodeVersion", [][]byte{[]byte("instantiate"), []byte("args")}, "chaincodePolicy")
// Must version is also available
```

### User client

#### Create user client
```go
userClient, err := fabricClient.CreateUserClient("userName", "orgTitle", "channelID")
// or without reuse of fabric client
userClient, err := fabclient.CreateUserClient("config file for fabric-sdk-go", "orderer host", "userName", "orgTitle", "channelID")
// Must versions is also available
```

#### Invoke transaction
```go
response, err = userClient.Invoke("chaincodeID", "chaincodeMethod", [][]byte{[]byte("method"), []byte("args")})
// Must version is also available
```

#### Query transaction (transaction won't be recorded to blockchain)
```go
response, err = userClient.Query("chaincodeID", "chaincodeMethod", [][]byte{[]byte("method"), []byte("args")})
// Must version is also available
```

#### Query Int. The same as Query but converts response to int
```go
response, err = userClient.QueryInt("chaincodeID", "chaincodeMethod", [][]byte{[]byte("method"), []byte("args")})
// Must version is also available
```

### Chaincode client

#### Create user client
```go
chaincodeClient, err := fabricClient.CreateChaincodeClient("channelID", "chaincodeID", "userName", "orgTitle")
// or without reuse of fabric client
chaincodeClient, err := fabclient.CreateChaincodeClient("config file for fabric-sdk-go", "orderer host", "channelID", "chaincodeID", "userName", "orgTitle")
// Must versions is also available
```

#### Invoke transaction
```go
response, err = chaincodeClient.Invoke("chaincodeMethod", [][]byte{[]byte("method"), []byte("args")})
// Must version is also available
```

#### Query transaction (transaction won't be recorded to blockchain)
```go
response, err = chaincodeClient.Query("chaincodeMethod", [][]byte{[]byte("method"), []byte("args")})
// Must version is also available
```

#### Query Int. The same as Query but converts response to int
```go
response, err = chaincodeClient.QueryInt("chaincodeMethod", [][]byte{[]byte("method"), []byte("args")})
// Must version is also available
```
