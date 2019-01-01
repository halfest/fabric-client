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

type ConfigurationClient struct {
	resMgmtClient *resmgmt.Client
	Name          string
	Organization  string
	FabricClient  *FabricClient
}

type ChannelParameters struct {
	ChannelID         string
	ChannelConfigPath string
}

func CreateChannelParameters(channelID string, channelConfigPath string) *ChannelParameters {
	return &ChannelParameters{
		ChannelID:         channelID,
		ChannelConfigPath: os.Getenv("GOPATH") + channelConfigPath,
	}
}

func (c *ConfigurationClient) InstallChaincodeFromStructure(chaincodeParameters *ChaincodeParameters) error {
	return c.InstallChaincode(chaincodeParameters.ChaincodeID, chaincodeParameters.ChaincodePath, chaincodeParameters.Version)
}

func (c *ConfigurationClient) MustInstallChaincode(chaincodeID string, chaincodePath string, version string) {
	err := c.InstallChaincode(chaincodeID, chaincodePath, version)
	if err != nil {
		panic(err)
	}
}

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

func (c *ConfigurationClient) InstanciateChaincodeFromStructure(channelID string, chaincodeParameters *ChaincodeParameters) error {
	return c.InstanciateChaincode(channelID, chaincodeParameters.ChaincodeID, chaincodeParameters.ChaincodePath, chaincodeParameters.Version, chaincodeParameters.ArgsForInit, chaincodeParameters.Policy)
}

func (c *ConfigurationClient) MustInstanciateChaincode(channelID string, chaincodeID string, chaincodePath string, version string, args [][]byte, policy string) {
	err := c.InstanciateChaincode(channelID, chaincodeID, chaincodePath, version, args, policy)
	if err != nil {
		panic(err)
	}
}

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

func (c *ConfigurationClient) CreateChannelFromStructure(channelParameters *ChannelParameters) error {
	return c.CreateChannel(channelParameters.ChannelID, channelParameters.ChannelConfigPath)
}

func (c *ConfigurationClient) MustCreateChannel(channelID string, channelConfigPath string) {
	err := c.CreateChannel(channelID, channelConfigPath)
	if err != nil {
		panic(err)
	}
}

func (c *ConfigurationClient) CreateChannel(channelID string, channelConfigPath string) error {
	// logger.Debugf("Creating channel %s", channelID)
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
	logger.Debugf("Channel %s created", channelID)
	return nil
}

func (c *ConfigurationClient) JoinChannelFromStructure(channelParameters *ChannelParameters) error {
	return c.JoinChannel(channelParameters.ChannelID)
}

func (c *ConfigurationClient) MustJoinChannel(channelID string) {
	err := c.JoinChannel(channelID)
	if err != nil {
		panic(err)
	}
}

func (c *ConfigurationClient) JoinChannel(channelID string) error {
	// logger.Debugf("Joining channel %s", channelID)
	if err := c.resMgmtClient.JoinChannel(channelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(c.FabricClient.OrdererHost)); err != nil {
		return fmt.Errorf("Failed to join channel %s.\n Error: %v", channelID, err)
	}
	logger.Debugf("Channel %s joined", channelID)
	return nil
}

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

func (c *ConfigurationClient) MustCreateAndJoinChannel(channelID string, channelConfigPath string) {
	err := c.CreateAndJoinChannel(channelID, channelConfigPath)
	if err != nil {
		panic(err)
	}
}

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
	resourceManagerClientContext := c.FabricClient.SDK.Context(fabsdk.WithUser(c.Name), fabsdk.WithOrg(c.Organization))
	resMgmtClient, err := resmgmt.New(resourceManagerClientContext)
	if err != nil {
		return fmt.Errorf("Failed to create channel management client with user %s and organisation %s.\n Error: %v", c.Name, c.Organization, err)
	}
	c.resMgmtClient = resMgmtClient
	logger.Debug("Ressource management client created")
	return nil
}
