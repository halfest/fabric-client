package fabclient

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
)

type UserClient struct {
	name            string
	organization    string
	channelClient   *channel.Client
	channelID       string
	signingIdentity msp.SigningIdentity
}

// CreateUserClient is the same as  (c *FabricClient) CreateUserClient(channelID string, name string, organization string) but it does not reuse Fabric Client
func CreateUserClient(configPath string, ordererHost string, channelID string, name string, organization string) (*UserClient, error) {
	fabricClient, err := CreateFabricClient(configPath, ordererHost)
	if err != nil {
		return nil, err
	}
	return fabricClient.CreateUserClient(channelID, name, organization)
}

// Invoke triggers invokation of transaction
func (c *UserClient) Invoke(chaincodeID string, functionName string, args [][]byte) ([]byte, error) {
	resp, err := c.channelClient.Execute(channel.Request{ChaincodeID: chaincodeID, Fcn: functionName, Args: args},
		channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		return nil, fmt.Errorf("Failed to invoke chaincode %s with funactions %s and arguments %v.\n Error: %v", chaincodeID, functionName, args, err)
	}
	logger.Debugf("Response on invoke chaincode: %s\n", resp.Payload)
	return resp.Payload, nil
}

// Query is the same as Invoke but without sending transaction to orderer so tx does not added to blockchain history. It is used for querying data
func (c *UserClient) Query(chaincodeID string, functionName string, args [][]byte) ([]byte, error) {
	resp, err := c.channelClient.Query(channel.Request{ChaincodeID: chaincodeID, Fcn: functionName, Args: args},
		channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		return nil, fmt.Errorf("Failed to query chaincode %s with funactions %s and arguments %v.\n Error: %v", chaincodeID, functionName, args, err)
	}
	logger.Debugf("Response on query chaincode: %s\n", resp.Payload)
	return resp.Payload, nil
}

// GetSigningIdentity return SigningIdentity of user
func (c *UserClient) GetSigningIdentity() msp.SigningIdentity {
	return c.signingIdentity
}
