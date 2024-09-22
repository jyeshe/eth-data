package main

import (
	"bufio"
	"eth_data/internal/cmd"
	"eth_data/internal/db"
	"eth_data/internal/eth"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	db, err := db.Open("data")
	if err != nil {
		log.Fatalf("Failed to open the database: %v", err)
	}
	defer db.Close()

	// Load the .env file
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	secretUrl := os.Getenv("SECRET_URL")
	if secretUrl == "" {
		log.Fatal("Missing SECRET_URL env var")
	}

	ethClient := eth.NewEthClient(secretUrl)
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
	// text := "fetchWalletTxs 0x0839D7318603205fEF3340f365b087A5F92c6Df8 20804500 20804599"
	command, args := parseCommand(text)

	switch command {
	case "clear":
		clearScreen()
	case "help":
		displayHelp()
	case "blockNumber":
		blockNumber, err := ethClient.LastestBlockNumber()
		if err == nil {
			fmt.Println("Latest block number:", blockNumber)
		} else {
			fmt.Printf("Error: %+v\n", err)
		}

	case "blockByNumber":
		block, err := ethClient.BlockByNumber(args[0])
		if err == nil {
			fmt.Printf("Dump: %+v\n", block)
		} else {
			fmt.Printf("Error: %+v\n", err)
		}

	case "lastFinalizedBlock":
		block, err := ethClient.BlockByNumber("finalized")
		if err == nil {
			fmt.Println("Keys:", block.BlockKeys())
			fmt.Println("Timestamp:", dateTimeFormat(block.Timestamp()))
		} else {
			fmt.Printf("Error: %+v\n", err)
		}

	case "lastFinalizedTx":
		block, err := ethClient.BlockByNumber("finalized")
		if err == nil {
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
		} else {
			fmt.Printf("Error: %+v\n", err)
		}
	case "fetchWalletTxs":
		if len(args) != 3 {
			fmt.Println("Usage: fetchWalletTxs <walletAddress> <firstBlockInt> <lastBlockInt>")
		} else {
			walletAddress := strings.ToLower(args[0])
			firstBlock, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				fmt.Println("Usage: fetchWalletTxs <walletAddress> <firstBlockInt> <lastBlockInt>")
			} else {
				lastBlock, err := strconv.ParseInt(args[2], 10, 64)
				if err != nil {
					fmt.Println("Usage: fetchWalletTxs <walletAddress> <firstBlockInt> <lastBlockInt>")
				} else {
					cmd.FetchWalletTxsCmd(ethClient, walletAddress, firstBlock, lastBlock)
				}
			}
		}

	default:
		fmt.Println(text, ": command not found")
	}
}

func parseCommand(text string) (string, []string) {
	parts := strings.Split(text, " ")
	command := parts[0]
	args := parts[1:]
	return command, args
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
