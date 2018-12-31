package fabclient

import (
	"fmt"
	"os"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/logging"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

var logger = logging.NewLogger("fabclient")

type FabricClient struct {
	SDK         *fabsdk.FabricSDK
	OrdererHost string
	GoPath      string
}

type ChaincodeParameters struct {
	ChaincodeID   string
	ChaincodePath string
	Version       string
	ArgsForInit   [][]byte
	Policy        string
}

func CreateChaincodeParameters(chaincodeID string, chaincodePath string, version string, argsForInit [][]byte, policy string) *ChaincodeParameters {
	return &ChaincodeParameters{
		ChaincodeID:   chaincodeID,
		ChaincodePath: chaincodePath,
		Version:       version,
		ArgsForInit:   argsForInit,
		Policy:        policy,
	}
}

func MustCreateFabricClient(configPath string, ordererHost string) *FabricClient {
	result, err := CreateFabricClient(configPath, ordererHost)
	if err != nil {
		panic(err)
	}
	return result
}

func CreateFabricClient(configPath string, ordererHost string) (*FabricClient, error) {
	var err error
	FabricClient := FabricClient{
		OrdererHost: ordererHost,
		GoPath:      os.Getenv("GOPATH"),
	}
	cp := config.FromFile(configPath)
	FabricClient.SDK, err = fabsdk.New(cp)
	if err != nil {
		return nil, fmt.Errorf("Failed to read fabric SDK config file: %s", err)
	}
	logger.Debug("fabric-client created")
	return &FabricClient, nil
}

func (c *FabricClient) MustCreateConfigurationClient(name string, organization string) *ConfigurationClient {
	result, err := c.CreateConfigurationClient(name, organization)
	if err != nil {
		panic(err)
	}
	return result
}

func (c *FabricClient) CreateConfigurationClient(name string, organization string) (*ConfigurationClient, error) {
	var err error
	configurationClient := &ConfigurationClient{
		Name:         name,
		Organization: organization,
		FabricClient: c,
	}
	err = configurationClient.initResourceMgmtClient(c.SDK)
	if err != nil {
		return nil, err
	}
	logger.Debugf("Configuration client for user: %s and organization: %s created", name, organization)
	return configurationClient, nil
}

func (c *FabricClient) MustCreateUserClient(channelID string, name string, organization string) *UserClient {
	result, err := c.CreateUserClient(channelID, name, organization)
	if err != nil {
		panic(err)
	}
	return result
}

func (c *FabricClient) CreateUserClient(channelID string, name string, organization string) (*UserClient, error) {
	var err error
	userClient := &UserClient{
		Name:         name,
		Organization: organization,
		ChannelID:    channelID,
	}

	channelProvider := c.SDK.ChannelContext(userClient.ChannelID, fabsdk.WithUser(userClient.Name), fabsdk.WithOrg(userClient.Organization))
	clientInstance, err := channel.New(channelProvider)
	if err != nil {
		return nil, fmt.Errorf("Failed to create user client with channel id %s, user name %s and organization %s.\n Error: %v", userClient.ChannelID, userClient.Name, userClient.Organization, err)
	}
	userClient.ChannelClient = clientInstance

	userClient.SigningIdentity, err = c.getUserIdentity(userClient.Name, userClient.Organization)
	if err != nil {
		return nil, err
	}
	logger.Debugf("User client for channelID: %s, user: %s and organization: %screated", channelID, name, organization)
	return userClient, nil
}

func (c *FabricClient) MustCreateChaincodeClient(channelID string, chaincodeID string, name string, organization string) *ChaincodeClient {
	result, err := c.CreateChaincodeClient(channelID, chaincodeID, name, organization)
	if err != nil {
		panic(err)
	}
	return result
}

func (c *FabricClient) CreateChaincodeClient(channelID string, chaincodeID string, name string, organization string) (*ChaincodeClient, error) {
	var err error
	chaincodeClient := &ChaincodeClient{
		ChaincodeID: chaincodeID,
	}
	chaincodeClient.UserClient, err = c.CreateUserClient(channelID, name, organization)
	if err != nil {
		return nil, fmt.Errorf("Failed to create user client with channel id %s, user name %s and organization %s.\n Error: %v", channelID, name, organization, err)
	}
	logger.Debugf("Chaincode client for channelID: %s, chaincodeID: %s, user: %s and organization: %screated", channelID, chaincodeID, name, organization)
	return chaincodeClient, nil
}

func (c *FabricClient) getUserIdentity(name string, organization string) (msp.SigningIdentity, error) {
	mspClient, err := mspclient.New(c.SDK.Context(), mspclient.WithOrg(organization))
	if err != nil {
		return nil, fmt.Errorf("Failed to create msp client with organisation %s.\n Error: %v", name, err)
	}
	userIdentity, err := mspClient.GetSigningIdentity(name)
	if err != nil {
		return nil, fmt.Errorf("Failed to get user signing identity with name: %s.\n Error: %v", name, err)
	}
	return userIdentity, nil
}
