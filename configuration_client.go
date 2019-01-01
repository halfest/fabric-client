package fabclient

import (
	"fmt"
	"os"

	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/resource"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	pb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
	platform "github.com/hyperledger/fabric/core/chaincode/platforms/golang"
)

// ConfigurationClient
type ConfigurationClient struct {
	resMgmtClient *resmgmt.Client
	name          string
	organization  string
	fabricClient  *FabricClient
}

// ChannelParameters contains data used to call functions that requires struct as argument
type ChannelParameters struct {
	ChannelID         string
	ChannelConfigPath string
}

// ChaincodeParameters is representation for parameters used to interact with chaincode
type ChaincodeParameters struct {
	ChaincodeID   string
	ChaincodePath string
	Version       string
	ArgsForInit   [][]byte
	Policy        string
}

// CreateChannelParameters constructs ChannelParameters
func CreateChannelParameters(channelID string, channelConfigPath string) *ChannelParameters {
	return &ChannelParameters{
		ChannelID:         channelID,
		ChannelConfigPath: os.Getenv("GOPATH") + channelConfigPath,
	}
}

// CreateChaincodeParameters used to construct ChaincodeParameters
func CreateChaincodeParameters(chaincodeID string, chaincodePath string, version string, argsForInit [][]byte, policy string) *ChaincodeParameters {
	return &ChaincodeParameters{
		ChaincodeID:   chaincodeID,
		ChaincodePath: chaincodePath,
		Version:       version,
		ArgsForInit:   argsForInit,
		Policy:        policy,
	}
}

// CreateConfigurationClient is the same as  (c *FabricClient) CreateConfigurationClient(channelID string, name string, organization string) but it does not reuse Fabric Client
func CreateConfigurationClient(configPath string, ordererHost string, name string, organization string) (*ConfigurationClient, error) {
	fabricClient, err := CreateFabricClient(configPath, ordererHost)
	if err != nil {
		return nil, err
	}
	return fabricClient.CreateConfigurationClient(name, organization)
}

// CreateChannelFromStructure the sames as CreateChannel but accepts ChannelParameters struct
func (c *ConfigurationClient) CreateChannelFromStructure(channelParameters *ChannelParameters) error {
	return c.CreateChannel(channelParameters.ChannelID, channelParameters.ChannelConfigPath)
}

// CreateChannel creates channel
func (c *ConfigurationClient) CreateChannel(channelID string, channelConfigPath string) error {
	// logger.Debugf("Creating channel %s", channelID)
	mspClient, err := mspclient.New(c.fabricClient.sdk.Context(), mspclient.WithOrg(c.organization))
	if err != nil {
		return fmt.Errorf("Failed to create msp client with organisation %s.\n Error: %s", c.organization, err)
	}
	userIdentity, err := mspClient.GetSigningIdentity(c.name)
	if err != nil {
		return fmt.Errorf("Failed to get signing identity %s while creating channel [%s].\n Error: %v", c.name, channelID, err)
	}
	req := resmgmt.SaveChannelRequest{ChannelID: channelID, ChannelConfigPath: channelConfigPath, SigningIdentities: []msp.SigningIdentity{userIdentity}}
	txID, err := c.resMgmtClient.SaveChannel(req, resmgmt.WithOrdererEndpoint(c.fabricClient.ordererHost))
	if err != nil || txID.TransactionID == "" {
		return fmt.Errorf("Failed to save channel %s.\n Error: %s", channelID, err)
	}
	logger.Debugf("Channel %s created", channelID)
	return nil
}

// InstallChaincodeFromStructure the sames as InstallChaincode but accepts ChaincodeParameters struct
func (c *ConfigurationClient) InstallChaincodeFromStructure(chaincodeParameters *ChaincodeParameters) error {
	return c.InstallChaincode(chaincodeParameters.ChaincodeID, chaincodeParameters.ChaincodePath, chaincodeParameters.Version)
}

// InstallChaincode installs chaincode
func (c *ConfigurationClient) InstallChaincode(chaincodeID string, chaincodePath string, version string) error {
	// logger.Debugf("Installing chaincode %s version %s", chaincodeID, version)
	goPlatform := platform.Platform{}
	payload, err := goPlatform.GetDeploymentPayload(chaincodePath)
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
	logger.Debugf("Chaincode %s version %s installed", chaincodeID, version)
	return nil
}

// InstanciateChaincodeFromStructure the sames as InstanciateChaincode but accepts ChaincodeParameters struct
func (c *ConfigurationClient) InstanciateChaincodeFromStructure(channelID string, chaincodeParameters *ChaincodeParameters) error {
	return c.InstanciateChaincode(channelID, chaincodeParameters.ChaincodeID, chaincodeParameters.ChaincodePath, chaincodeParameters.Version, chaincodeParameters.ArgsForInit, chaincodeParameters.Policy)
}

// InstanciateChaincode instantiates chaincode
func (c *ConfigurationClient) InstanciateChaincode(channelID string, chaincodeID string, chaincodePath string, version string, args [][]byte, policy string) error {
	// logger.Debugf("Instantiating chaincode %s version %s", chaincodeID, version)
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
	logger.Debugf("Chaincode %s version %s instantiated", chaincodeID, version)
	return nil
}

// JoinChannelFromStructure the sames as JoinChannel but accepts ChannelParameters struct
func (c *ConfigurationClient) JoinChannelFromStructure(channelParameters *ChannelParameters) error {
	return c.JoinChannel(channelParameters.ChannelID)
}

// JoinChannel joins channel
func (c *ConfigurationClient) JoinChannel(channelID string) error {
	// logger.Debugf("Joining channel %s", channelID)
	if err := c.resMgmtClient.JoinChannel(channelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(c.fabricClient.ordererHost)); err != nil {
		return fmt.Errorf("Failed to join channel %s.\n Error: %v", channelID, err)
	}
	logger.Debugf("Channel %s joined", channelID)
	return nil
}

// CreateAndJoinChannelFromStructure the sames as CreateAndJoinChannel but accepts ChannelParameters struct
func (c *ConfigurationClient) CreateAndJoinChannelFromStructure(channelParameters *ChannelParameters) error {
	var err error
	err = c.CreateChannelFromStructure(channelParameters)
	if err != nil {
		return fmt.Errorf("Failed to create channel with structure %+v.\n Error: %v", channelParameters, err)
	}
	err = c.JoinChannelFromStructure(channelParameters)
	if err != nil {
		return fmt.Errorf("Failed to join channel with structure %+v.\n Error: %v", channelParameters, err)
	}
	return nil
}

// CreateAndJoinChannel creates and joins channel
func (c *ConfigurationClient) CreateAndJoinChannel(channelID string, channelConfigPath string) error {
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

func (c *ConfigurationClient) initResourceMgmtClient(sdk *fabsdk.FabricSDK) error {
	// logger.Debug("Creating ressource management client")
	// The resource management client is responsible for managing channels (create/update channel)
	resourceManagerClientContext := c.fabricClient.sdk.Context(fabsdk.WithUser(c.name), fabsdk.WithOrg(c.organization))
	resMgmtClient, err := resmgmt.New(resourceManagerClientContext)
	if err != nil {
		return fmt.Errorf("Failed to create channel management client with user %s and organisation %s.\n Error: %v", c.name, c.organization, err)
	}
	c.resMgmtClient = resMgmtClient
	logger.Debug("Ressource management client created")
	return nil
}
