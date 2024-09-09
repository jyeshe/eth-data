package eth

type EthTx struct {
	tx map[string]interface{}
}

func NewTx(tx map[string]interface{}) *EthTx {
	return &EthTx{
		tx: tx,
	}
}

func (t *EthTx) TxKeys() []string {
	var txKeys []string

	for key := range t.tx {
		txKeys = append(txKeys, key)
	}

	return txKeys
}

func (t *EthTx) BlockNumber() int64 {
	return ParseHex(t.tx["blockNumber"].(string))
}

func (t *EthTx) Gas() int64 {
	return ParseHex(t.tx["gas"].(string))
}

func (t *EthTx) GasPrice() int64 {
	return ParseHex(t.tx["gasPrice"].(string))
}

func (t *EthTx) GetString(field string) string {
	return t.tx[field].(string)
}
