package eth

type EthBlock struct {
	block map[string]interface{}
}

func NewBlock(block map[string]interface{}) *EthBlock {
	return &EthBlock{
		block: block,
	}
}

func (b *EthBlock) BlockKeys() []string {
	var blockKeys []string

	for key := range b.block {
		blockKeys = append(blockKeys, key)
	}

	return blockKeys
}

func (b *EthBlock) Timestamp() int64 {
	return parseHex(b.block["timestamp"].(string))
}

func (b *EthBlock) LastTransaction() *EthTx {
	transactions := b.block["transactions"].([]interface{})
	lastTransaction := transactions[len(transactions)-1].(map[string]interface{})

	return NewTx(lastTransaction)
}
