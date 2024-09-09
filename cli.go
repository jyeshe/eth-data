package main

import (
	"bufio"
	"fmt"
	"internal/eth"
	"log"
	"math"
	"os"
	"os/exec"
	"strings"
	"time"

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

	ethClient := eth.NewEthClient(apiKey)
	defer ethClient.Close()

	promptLoop(ethClient)
}

func promptLoop(ethClient *eth.EthClient) {
	reader := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for reader.Scan() {
		text := strings.TrimSpace(reader.Text())
		if text == "exit" {
			// Close the program
			return
		} else {
			handleCommand(text, ethClient)
		}
		fmt.Print("> ")
	}
}

// Executes the given command
func handleCommand(text string, ethClient *eth.EthClient) {
	switch text {
	case "clear":
		clearScreen()
	case "help":
		displayHelp()
	case "blockNumber":
		blockNumber := ethClient.LastestBlockNumber()
		fmt.Println("Latest block number:", blockNumber)

	case "lastFinalizedBlock":
		block := ethClient.BlockByNumber("finalized")
		fmt.Println("Keys:", block.BlockKeys())
		fmt.Println("Timestamp:", dateTimeFormat(block.Timestamp()))

	case "lastFinalizedTx":
		block := ethClient.BlockByNumber("finalized")
		lastTransaction := block.LastTransaction()
		fmt.Println("Last finalized block number:", lastTransaction.BlockNumber())
		fmt.Printf("Dump: %+v\n", lastTransaction)
		fmt.Printf("Keys: %+v\n", lastTransaction.TxKeys())
		fmt.Println("Hash:", lastTransaction.GetString("hash"))
		fmt.Println("Gas:", lastTransaction.Gas())
		fmt.Println("GasPrice:", lastTransaction.GasPrice())
		fmt.Println("BaseFeePerGas:", block.BaseFeePerGas())
		fmt.Println("MaxTipPerGas:", eth.ParseHex(lastTransaction.GetString("maxPriorityFeePerGas")))
		fmt.Println("MaxFeePerGas:", eth.ParseHex(lastTransaction.GetString("maxFeePerGas")))

		gasPrice1 := block.BaseFeePerGas() + eth.ParseHex(lastTransaction.GetString("maxPriorityFeePerGas"))
		gasPrice2 := eth.ParseHex(lastTransaction.GetString("maxFeePerGas"))

		fmt.Println("EffectiveGasPrice:", math.Min(float64(gasPrice1), float64(gasPrice2)))
		fmt.Printf("Input: %+v\n", lastTransaction.GetString("input"))
	default:
		fmt.Println(text, ": command not found")
	}
}

// Shows the available commands
func displayHelp() {
	fmt.Printf("This is a CLI to retrieve Ethereum data.\n\n")
	fmt.Println("Available commands:")

	fmt.Println("blockNumber        - Gets the latest block number")
	fmt.Println("lastFinalizedBlock - Gets the last finalized block")
	fmt.Println("lastFinalizedTx    - Gets the last tx from finalized number")
	fmt.Println("clear - Clear the terminal screen")
	fmt.Println("help  - Show available commands")
	fmt.Println("exit  - Closes this program")
}

// Clears the terminal screen
func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func dateTimeFormat(timestamp int64) string {
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}
