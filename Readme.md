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

## Usage

### Init fabric client
```go
fabricClient, err := fabclient.CreateFabricClient("config file for fabric-sdk-go", "orderer host")
// Must version is also available
```

### Configuration client

#### Create configuration client
```go
import (
	fabclient "github.com/halfest/fabric-client"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/logging"
)

logging.SetLevel("fabclient", logging.DEBUG)
configurationClient, err := fabricClient.CreateConfigurationClient("userName", "orgTitle")
// Must version is also available
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
// Must version is also available
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

### Chaincode client

#### Create user client
```go
chaincodeClient, err := fabricClient.CreateChaincodeClient("channelID", "chaincodeID", "userName", "orgTitle")
// Must version is also available
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
