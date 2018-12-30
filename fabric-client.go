package fabclient

import (
	"fmt"
	"os"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/logging"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/resource"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	pb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
	platform "github.com/hyperledger/fabric/core/chaincode/platforms/golang"
)

var logger = logging.NewLogger("fabclient")

type FabricClient struct {
	SDK         *fabsdk.FabricSDK
	OrdererHost string
	GoPath      string
}

type ConfiguratorClient struct {
	resMgmtClient *resmgmt.Client
	Name          string
	Organization  string
	FabricClient  *FabricClient
}

type UserClient struct {
	Name            string
	Organization    string
	ChannelClient   *channel.Client
	ChannelID       string
	SigningIdentity msp.SigningIdentity
}

type ChannelParameters struct {
	ChannelID         string
	ChannelConfigPath string
}

type ChaincodeParameters struct {
	ChaincodeID   string
	ChaincodePath string
	Version       string
	ArgsForInit   [][]byte
	Policy        string
}

func CreateChannelParameters(channelID string, channelConfigPath string) *ChannelParameters {
	return &ChannelParameters{
		ChannelID:         channelID,
		ChannelConfigPath: os.Getenv("GOPATH") + channelConfigPath,
	}
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

func NewFabricClient(configPath string, ordererHost string) (*FabricClient, error) {
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
	logger.Info("SDK created")
	return &FabricClient, nil
}

func (f *FabricClient) CreateConfigurationClient(name string, organization string) (*ConfiguratorClient, error) {
	var err error
	newAdmin := &ConfiguratorClient{
		Name:         name,
		Organization: organization,
		FabricClient: f,
	}
	err = newAdmin.initResourceMgmtClient(f.SDK)
	if err != nil {
		return nil, err
	}
	return newAdmin, nil
}

func (c *ConfiguratorClient) InstallChaincodeFromStructure(chaincodeParameters *ChaincodeParameters) error {
	return c.InstallChaincode(chaincodeParameters.ChaincodeID, chaincodeParameters.ChaincodePath, chaincodeParameters.Version)
}

func (c *ConfiguratorClient) InstallChaincode(chaincodeID string, chaincodePath string, version string) error {
	platform := platform.Platform{}
	payload, err := platform.GetDeploymentPayload(chaincodePath)
	if err != nil {
		return fmt.Errorf("Failed to create chaincode package with chaincode path %s.\n Error: %v", chaincodePath, err)
	}
	ccPkg := &resource.CCPackage{Type: pb.ChaincodeSpec_GOLANG, Code: payload}
	// Install example cc to org peers
	installCCReq := resmgmt.InstallCCRequest{Name: chaincodeID, Path: chaincodePath, Version: version, Package: ccPkg}
	_, err = c.resMgmtClient.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return fmt.Errorf("Failed to install chaincode with chaincode id %s, chaincode path and version %s.\n Error: %v", chaincodeID, chaincodePath, version, err)
	}
	logger.Infof("Chaincode %s version %s installed", chaincodeID, version)
	return nil
}

func (c *ConfiguratorClient) InstanciateChaincodeFromStructure(channelID string, chaincodeParameters *ChaincodeParameters) error {
	return c.InstanciateChaincode(channelID, chaincodeParameters.ChaincodeID, chaincodeParameters.ChaincodePath, chaincodeParameters.Version, chaincodeParameters.ArgsForInit, chaincodeParameters.Policy)
}

func (c *ConfiguratorClient) InstanciateChaincode(channelID string, chaincodeID string, chaincodePath string, version string, args [][]byte, policy string) error {
	ccPolicy, err := cauthdsl.FromString(policy)
	if err != nil {
		return fmt.Errorf("Failed to construct signature policy from string %s.\n Error: %v", policy, err)
	}
	resp, err := c.resMgmtClient.InstantiateCC(channelID,
		resmgmt.InstantiateCCRequest{Name: chaincodeID, Path: chaincodePath, Version: version, Args: args, Policy: ccPolicy},
	)
	if err != nil || resp.TransactionID == "" {
		return fmt.Errorf("Failed to instantiate the chaincode with channelID: %s, chaincodeID: %s, chaincodePath: %s, version: %s, args: %v and signature policy: %s.\n Error: %v", channelID, chaincodeID, chaincodePath, version, args, policy, err)
	}
	logger.Infof("Chaincode %s version %s instantiated", chaincodeID, version)
	return nil
}

//getResourceManClient
//returns a resource management client instance.
func (c *ConfiguratorClient) initResourceMgmtClient(sdk *fabsdk.FabricSDK) error {
	// var err error
	// The resource management client is responsible for managing channels (create/update channel)
	resourceManagerClientContext := c.FabricClient.SDK.Context(fabsdk.WithUser(c.Name), fabsdk.WithOrg(c.Organization))
	resMgmtClient, err := resmgmt.New(resourceManagerClientContext)
	if err != nil {
		return fmt.Errorf("Failed to create channel management client with user %s and organisation %s.\n Error: %v", c.Name, c.Organization, err)
	}
	c.resMgmtClient = resMgmtClient
	logger.Info("Ressource management client created")
	return nil
}

func (c *ConfiguratorClient) CreateChannelWithStructure(channelParameters *ChannelParameters) error {
	return c.CreateChannel(channelParameters.ChannelID, channelParameters.ChannelConfigPath)
}

func (c *ConfiguratorClient) CreateChannel(channelID string, channelConfigPath string) error {
	mspClient, err := mspclient.New(c.FabricClient.SDK.Context(), mspclient.WithOrg(c.Organization))
	if err != nil {
		return fmt.Errorf("Failed to create msp client with organisation %s.\n Error: %s", c.Organization, err)
	}
	userIdentity, err := mspClient.GetSigningIdentity(c.Name)
	if err != nil {
		return fmt.Errorf("Failed to get signing identity %s while creating channel [%s].\n Error: %v", c.Name, channelID, err)
	}
	req := resmgmt.SaveChannelRequest{ChannelID: channelID, ChannelConfigPath: channelConfigPath, SigningIdentities: []msp.SigningIdentity{userIdentity}}
	txID, err := c.resMgmtClient.SaveChannel(req, resmgmt.WithOrdererEndpoint(c.FabricClient.OrdererHost))
	if err != nil || txID.TransactionID == "" {
		return fmt.Errorf("Failed to save channel %s.\n Error: %s", channelID, err)
	}
	logger.Infof("Channel %s created", channelID)
	return nil
}

func (c *ConfiguratorClient) JoinChannelWithStructure(channelParameters *ChannelParameters) error {
	return c.JoinChannel(channelParameters.ChannelID)
}

func (c *ConfiguratorClient) JoinChannel(channelID string) error {
	if err := c.resMgmtClient.JoinChannel(channelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(c.FabricClient.OrdererHost)); err != nil {
		return fmt.Errorf("Failed to join channel %s.\n Error: %v", channelID, err)
	}
	logger.Infof("Channel %s joined", channelID)
	return nil
}

func (c *ConfiguratorClient) CreateAndJoinChannelWithStructure(channelParameters *ChannelParameters) error {
	var err error
	err = c.CreateChannelWithStructure(channelParameters)
	if err != nil {
		return fmt.Errorf("Failed to create channel with structure %+v.\n Error: %v", channelParameters, err)
	}
	err = c.JoinChannelWithStructure(channelParameters)
	if err != nil {
		return fmt.Errorf("Failed to join channel with structure %+v.\n Error: %v", channelParameters, err)
	}
	return nil
}

func (c *ConfiguratorClient) CreateAndJoinChannel(channelID string, channelConfigPath string) error {
	var err error
	err = c.CreateChannel(channelID, channelConfigPath)
	if err != nil {
		return fmt.Errorf("Failed to create channel with channelID %s and channelConfigPath %s.\n Error: %v", channelID, channelConfigPath, err)
	}
	err = c.JoinChannel(channelID)
	if err != nil {
		return fmt.Errorf("Failed to join channel with channelID %s.\n Error: %v", channelID, err)
	}
	return nil
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

func (c *FabricClient) CreateUserClient(name string, organization string, channelID string) (*UserClient, error) {
	var err error
	userClient := &UserClient{
		Name:         name,
		Organization: organization,
		ChannelID:    channelID,
	}

	clientContext := c.SDK.ChannelContext(userClient.ChannelID, fabsdk.WithUser(userClient.Name), fabsdk.WithOrg(userClient.Organization))
	clientInstance, err := channel.New(clientContext)
	if err != nil {
		return nil, fmt.Errorf("Failed to create channel client with channel id %s, user name %s and organization %s.\n Error: %v", userClient.ChannelID, userClient.Name, userClient.Organization, err)
	}
	userClient.ChannelClient = clientInstance

	userClient.SigningIdentity, err = c.getUserIdentity(userClient.Name, userClient.Organization)
	if err != nil {
		return nil, err
	}
	logger.Infof("Created client for user %s", userClient.Name)
	return userClient, nil
}

func (c *UserClient) initClientInstance(sdk *fabsdk.FabricSDK) error {
	// Channel client is used to query and execute transactions
	channelProvider := sdk.ChannelContext(c.ChannelID, fabsdk.WithUser(c.Name), fabsdk.WithOrg(c.Organization))
	channelClient, err := channel.New(channelProvider)
	if err != nil {
		return fmt.Errorf("Failed to create new channel client.\n Error: %v", err)
	}
	c.ChannelClient = channelClient
	logger.Infof("In channel %s client %s created", c.ChannelID, c.Name)
	return nil
}

func (c *UserClient) Invoke(chaincodeID string, functionName string, args [][]byte) ([]byte, error) {
	resp, err := c.ChannelClient.Execute(channel.Request{ChaincodeID: chaincodeID, Fcn: functionName, Args: args},
		channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		return nil, fmt.Errorf("Failed to invoke chaincode %s with funactions %s and arguments %v.\n Error: %v", chaincodeID, functionName, args, err)
	}
	logger.Infof("Response on invoke chaincode: %s\n", resp.Payload)
	return resp.Payload, nil
}

func (c *UserClient) Query(chaincodeID string, functionName string, args [][]byte) ([]byte, error) {
	resp, err := c.ChannelClient.Query(channel.Request{ChaincodeID: chaincodeID, Fcn: functionName, Args: args},
		channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		return nil, fmt.Errorf("Failed to query chaincode %s with funactions %s and arguments %v.\n Error: %v", chaincodeID, functionName, args, err)
	}
	logger.Infof("Response on query chaincode: %s\n", resp.Payload)
	return resp.Payload, nil
}
