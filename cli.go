package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/joho/godotenv"
)

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	apiKey := os.Getenv("INFURA_API_KEY")
	if apiKey == "" {
		log.Fatal("Missing INFURA_API_KEY env var")
	}

	// Connect to the Ethereum node (use your local node or an Infura endpoint)
	client, err := rpc.Dial("https://mainnet.infura.io/v3/" + apiKey)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()

	// Fetch the latest block number
	var blockNumber string
	err = client.CallContext(context.Background(), &blockNumber, "eth_blockNumber")
	if err != nil {
		log.Fatalf("Failed to retrieve the latest block number: %v", err)
	}

	fmt.Println("Latest block number:", parseHex(blockNumber))

	// Fetch last finalized block
	var block map[string]interface{}
	err = client.CallContext(context.Background(), &block, "eth_getBlockByNumber", "finalized", true)
	if err != nil {
		log.Fatalf("Failed to retrieve block: %v", err)
	}

	var blockKeys []string

	for key := range block {
		blockKeys = append(blockKeys, key)
	}

	// Print block details
	fmt.Println("Keys:", blockKeys)
	fmt.Println("Timestamp:", dateTimeFormat(parseHex(block["timestamp"].(string))))

	transactions := block["transactions"].([]interface{})
	lastTransaction := transactions[len(transactions)-1].(map[string]interface{})

	fmt.Println("Last finalized block:", parseHex(lastTransaction["blockNumber"].(string)))
	fmt.Printf("Last Transaction: %+v\n", lastTransaction)
	fmt.Printf("TransactionRoot: %+v\n", block["transactionsRoot"])
}

func parseHex(hexStr string) int64 {
	intValue, err := strconv.ParseInt(hexStr[2:], 16, 64)
	if err != nil {
		// Handle the error
		fmt.Println("Error:", err)
		return -1
	}

	return intValue
}

func dateTimeFormat(timestamp int64) string {
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}
