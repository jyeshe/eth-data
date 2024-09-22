package cmd

import (
	"encoding/json"
	"eth_data/internal/db"
	"eth_data/internal/eth"
	"eth_data/internal/tasks"
	"fmt"

	"github.com/dgraph-io/badger"
)

func FetchWalletTxsCmd(ec *eth.EthClient, walletAddress string, firstBlock int64, lastBlock int64) {
	var tasksArray []*tasks.Task

	for blockNumber := firstBlock; blockNumber <= lastBlock; blockNumber++ {
		blockNumberHex := fmt.Sprintf("0x%x", blockNumber)
		task := &tasks.Task{
			ID:   blockNumber,
			Kind: tasks.BlockByNumber,
			Args: []string{blockNumberHex},
		}
		tasksArray = append(tasksArray, task)
	}

	executeAndSave(ec, walletAddress, tasksArray)
}

func executeAndSave(ec *eth.EthClient, walletAddress string, tasksArray []*tasks.Task) {
	resultsCh := make(chan tasks.Result, len(tasksArray))

	tasks.ExecuteMany(ec, tasksArray, resultsCh)

	var retryTasks []*tasks.Task
	tasksMap := make(map[int64]*tasks.Task)

	for result := range resultsCh {
		if result.Err != nil {
			if len(tasksMap) == 0 {
				for _, task := range tasksArray {
					tasksMap[task.ID] = task
				}
			}
			task := tasksMap[result.TaskId]
			retryTasks = append(retryTasks, task)
		} else {
			for _, tx := range result.Block.Transactions() {
				if tx.GetString("to") == walletAddress || tx.GetString("from") == walletAddress {
					saveWalletTx(walletAddress, tx)
				}
			}
		}
	}

	if len(retryTasks) > 0 {
		executeAndSave(ec, walletAddress, retryTasks)
	}
}

func saveWalletTx(walletAddress string, tx *eth.EthTx) {
	txHash := tx.GetString("hash")
	blockNumber := tx.GetString("blockNumber")
	key := fmt.Sprintf("wallet:%s:%s:%s", walletAddress, blockNumber, txHash)
	value, err := json.Marshal(tx.GetAll())

	fmt.Printf("Found: %s\n", txHash)

	if err != nil {
		fmt.Printf("Error: %+v\n", err)
	} else {
		err = db.GetKvDb().Update(func(txn *badger.Txn) error {
			return txn.Set([]byte(key), []byte(value))
		})

		if err != nil {
			fmt.Printf("Err: %+v\n", err)
		}
	}
}
