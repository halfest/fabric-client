package fabclient

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
)

type LedgerClient struct {
	Client *ledger.Client
}

func CreateLedgerClient(configPath string, ordererHost string, channelID string, name string, organization string) (*LedgerClient, error) {
	fabricClient, err := CreateFabricClient(configPath, ordererHost)
	if err != nil {
		return nil, err
	}
	return fabricClient.CreateLedgerClient(channelID, name, organization)
}
