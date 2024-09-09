package eth

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/rpc"
)

type EthClient struct {
	rpcClient *rpc.Client
}

func NewEthClient(apiKey string) *EthClient {
	// Connect to the Ethereum node (use your local node or an Infura endpoint)
	rpcClient, err := rpc.Dial("https://mainnet.infura.io/v3/" + apiKey)

	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	client := &EthClient{
		rpcClient: rpcClient,
	}

	return client
}

func (c *EthClient) LastestBlockNumber() int64 {
	// Fetch the latest block number
	var blockNumber string
	err := c.rpcClient.CallContext(context.Background(), &blockNumber, "eth_blockNumber")
	if err != nil {
		log.Fatalf("Failed to retrieve the latest block number: %v", err)
	}

	return ParseHex(blockNumber)
}

func (c *EthClient) BlockByNumber(number string) *EthBlock {
	// Fetch last finalized block
	var block map[string]interface{}
	err := c.rpcClient.CallContext(context.Background(), &block, "eth_getBlockByNumber", number, true)
	if err != nil {
		log.Fatalf("Failed to retrieve block: %v", err)
	}

	return NewBlock(block)
}

func (c *EthClient) Close() {
	c.rpcClient.Close()
}
