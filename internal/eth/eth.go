package eth

import (
	"context"
	"eth_data/internal/rate_limiting"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/rpc"
)

const eth_blockNumberCredits = 1
const eth_getBlockByNumberTrueCredits = 80 * 2

type EthClient struct {
	rpcClient   *rpc.Client
	rateLimiter *rate_limiting.RateLimiter
}

func NewEthClient(secretUrl string) *EthClient {
	// Connect to the Ethereum node (use your local node or an Infura endpoint)
	rpcClient, err := rpc.Dial(secretUrl)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	client := &EthClient{
		rpcClient:   rpcClient,
		rateLimiter: rate_limiting.NewRateLimiter(),
	}

	return client
}

func (c *EthClient) LastestBlockNumber() (int64, error) {
	if c.rateLimiter.AllowedAfter(eth_blockNumberCredits) == 0 {
		var blockNumber string

		// Fetch the latest block number
		err := c.rpcClient.CallContext(context.Background(), &blockNumber, "eth_blockNumber")
		if err != nil {
			log.Fatalf("Failed to retrieve the latest block number: %v", err)
		}

		return ParseHex(blockNumber), nil
	} else {
		return -1, rpc.HTTPError{StatusCode: 429}
	}
}

func (c *EthClient) BlockByNumber(number string) (*EthBlock, error) {
	if c.rateLimiter.AllowedAfter(eth_getBlockByNumberTrueCredits) == 0 {
		var block map[string]interface{}

		// Fetch last finalized block
		err := c.rpcClient.CallContext(context.Background(), &block, "eth_getBlockByNumber", number, true)
		if err != nil {
			log.Fatalf("Failed to retrieve block: %v", err)
		}

		fmt.Println("block:", ParseHex(block["number"].(string)))
		return NewBlock(block), nil
	} else {
		return nil, rpc.HTTPError{StatusCode: 429}
	}
}

func (c *EthClient) Close() {
	c.rpcClient.Close()
}
