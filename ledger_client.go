package fabclient

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

type LedgerClient struct {
	client *ledger.Client
}

func CreateLedgerClient(configPath string, ordererHost string, channelID string, name string, organization string) (*LedgerClient, error) {
	fabricClient, err := CreateFabricClient(configPath, ordererHost)
	if err != nil {
		return nil, err
	}
	return fabricClient.CreateLedgerClient(channelID, name, organization)
}

// GetBlockchainInfoReponse request
func (l *LedgerClient) GetBlockchainInfoReponse() (*fab.BlockchainInfoResponse, error) {
	resp, err := l.client.QueryInfo(ledger.WithMaxTargets(3))
	if err != nil {
		return nil, fmt.Errorf("Failed to get GetBlockchainInfoReponse\n Error: %v", err)
	}
	if resp == nil {
		return nil, fmt.Errorf("GetBlockchainInfoReponse is nil")
	}
	return resp, nil
}

//GetBlochchainHeight get blockchain height
func (l *LedgerClient) GetBlochchainHeight() (int, error) {
	var chainHeight int
	resp, err := l.GetBlockchainInfoReponse()
	if err != nil {
		return chainHeight, err
	}
	chainHeight = int(resp.BCI.GetHeight())
	return chainHeight, nil
}
