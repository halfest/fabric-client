package fabclient

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
)

type UserClient struct {
	Name            string
	Organization    string
	ChannelClient   *channel.Client
	ChannelID       string
	SigningIdentity msp.SigningIdentity
}

func (c *UserClient) Invoke(chaincodeID string, functionName string, args [][]byte) ([]byte, error) {
	resp, err := c.ChannelClient.Execute(channel.Request{ChaincodeID: chaincodeID, Fcn: functionName, Args: args},
		channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		return nil, fmt.Errorf("Failed to invoke chaincode %s with funactions %s and arguments %v.\n Error: %v", chaincodeID, functionName, args, err)
	}
	logger.Debugf("Response on invoke chaincode: %s\n", resp.Payload)
	return resp.Payload, nil
}

func (c *UserClient) Query(chaincodeID string, functionName string, args [][]byte) ([]byte, error) {
	resp, err := c.ChannelClient.Query(channel.Request{ChaincodeID: chaincodeID, Fcn: functionName, Args: args},
		channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		return nil, fmt.Errorf("Failed to query chaincode %s with funactions %s and arguments %v.\n Error: %v", chaincodeID, functionName, args, err)
	}
	logger.Debugf("Response on query chaincode: %s\n", resp.Payload)
	return resp.Payload, nil
}
