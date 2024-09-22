package tasks

import (
	"eth_data/internal/eth"
	"sync"
)

type TaskKind int

const (
	BlockByNumber = iota
)

type Task struct {
	ID   int64
	Kind int
	Args []string
}

type Result struct {
	TaskId int64
	Err    error
	Block  *eth.EthBlock
}

func ExecuteMany(ec *eth.EthClient, tasks []*Task, resultsCh chan<- Result) {
	var wg sync.WaitGroup
	wg.Add(len(tasks))

	for _, task := range tasks {
		go Execute(ec, task, resultsCh, &wg)
	}
	wg.Wait()
	close(resultsCh)
}

func Execute(ec *eth.EthClient, task *Task, resultsCh chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	switch task.Kind {
	case BlockByNumber:
		block, err := ec.BlockByNumber(task.Args[0])

		if err == nil {
			resultsCh <- Result{TaskId: task.ID, Block: block}
		} else {
			resultsCh <- Result{TaskId: task.ID, Err: err}
		}
	}
}
