package fabclient

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/logging"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

var logger = logging.NewLogger("fabclient")

// FabricClient contains FabricSDK and used to interact with fabric system
type FabricClient struct {
	sdk         *fabsdk.FabricSDK
	ordererHost string
}

// CreateFabricClient creates new Fabric Client
func CreateFabricClient(configPath string, ordererHost string) (*FabricClient, error) {
	var err error
	cp := config.FromFile(configPath)
	sdk, err := fabsdk.New(cp)
	if err != nil {
		return nil, fmt.Errorf("Failed to read fabric SDK config file: %s", err)
	}
	fabricClient := CreateFabricClientFromSDK(sdk, ordererHost)
	return fabricClient, nil
}

// CreateFabricClientFromSDK creates new Fabric Client based on passed sdk
func CreateFabricClientFromSDK(sdk *fabsdk.FabricSDK, ordererHost string) *FabricClient {
	FabricClient := FabricClient{
		ordererHost: ordererHost,
		sdk:         sdk,
	}
	logger.Debug("fabric-client created")
	return &FabricClient
}

// CreateConfigurationClient creates new Configuration Client
func (c *FabricClient) CreateConfigurationClient(name string, organization string) (*ConfigurationClient, error) {
	var err error
	configurationClient := &ConfigurationClient{
		name:         name,
		organization: organization,
		fabricClient: c,
	}
	err = configurationClient.initResourceMgmtClient(c.sdk)
	if err != nil {
		return nil, err
	}
	logger.Debugf("Configuration client for user: %s and organization: %s created", name, organization)
	return configurationClient, nil
}

// CreateUserClient creates new User Client
func (c *FabricClient) CreateUserClient(channelID string, name string, organization string) (*UserClient, error) {
	var err error
	userClient := &UserClient{
		name:         name,
		organization: organization,
		channelID:    channelID,
	}

	channelProvider := c.sdk.ChannelContext(userClient.channelID, fabsdk.WithUser(userClient.name), fabsdk.WithOrg(userClient.organization))
	clientInstance, err := channel.New(channelProvider)
	if err != nil {
		return nil, fmt.Errorf("Failed to create user client with channel id %s, user name %s and organization %s.\n Error: %v", userClient.channelID, userClient.name, userClient.organization, err)
	}
	userClient.channelClient = clientInstance

	userClient.signingIdentity, err = c.getUserIdentity(userClient.name, userClient.organization)
	if err != nil {
		return nil, err
	}
	logger.Debugf("User client for channelID: %s, user: %s and organization: %screated", channelID, name, organization)
	return userClient, nil
}

// CreateLedgerClient creates new User Ledger Client
func (c *FabricClient) CreateLedgerClient(channelID string, name string, organization string) (*LedgerClient, error) {
	var err error
	ledgerClient := &LedgerClient{}
	channelProvider := c.sdk.ChannelContext(channelID, fabsdk.WithUser(name), fabsdk.WithOrg(organization))
	l, err := ledger.New(channelProvider)
	if err != nil {
		return nil, fmt.Errorf("Failed to create ledger client with channel id %s, user name %s and organization %s.\n Error: %v", channelID, name, organization, err)
	}
	ledgerClient.client = l
	logger.Debugf("Ledger client for channelID: %s, user: %s and organization: %screated", channelID, name, organization)
	return ledgerClient, nil
}

// CreateChaincodeClient creates new Chaincode Client
func (c *FabricClient) CreateChaincodeClient(channelID string, chaincodeID string, name string, organization string) (*ChaincodeClient, error) {
	var err error
	chaincodeClient := &ChaincodeClient{
		chaincodeID: chaincodeID,
	}
	chaincodeClient.userClient, err = c.CreateUserClient(channelID, name, organization)
	if err != nil {
		return nil, fmt.Errorf("Failed to create user client with channel id %s, user name %s and organization %s.\n Error: %v", channelID, name, organization, err)
	}
	logger.Debugf("Chaincode client for channelID: %s, chaincodeID: %s, user: %s and organization: %screated", channelID, chaincodeID, name, organization)
	return chaincodeClient, nil
}

func (c *FabricClient) getUserIdentity(name string, organization string) (msp.SigningIdentity, error) {
	mspClient, err := mspclient.New(c.sdk.Context(), mspclient.WithOrg(organization))
	if err != nil {
		return nil, fmt.Errorf("Failed to create msp client with organisation %s.\n Error: %v", name, err)
	}
	userIdentity, err := mspClient.GetSigningIdentity(name)
	if err != nil {
		return nil, fmt.Errorf("Failed to get user signing identity with name: %s.\n Error: %v", name, err)
	}
	return userIdentity, nil
}
