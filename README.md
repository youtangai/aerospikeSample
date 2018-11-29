# example

```
func main() {
	// クライアントを取得する
	client, err := NewAeroSpikeClient(AEROSPIKE_HOST, AEROSPIKE_PORT)
	if err != nil {
		panic(err)
	}

	// ダミーデータ作成
	block := Block{
		Id:         "testid",
		Version:    12,
		Prehash:    "testprehash",
		Merkleroot: "testmerkleroot",
		Timestamp:  "test_timestamp",
		Level:      "test_level",
		Nonce:      123,
		Size:       1234,
		Txcount:    12345,
		TxidList:   []string{"testid1", "testid2"},
	}
	tx := Transaction{
		Txid:      "testtxid",
		Output:    "testoutput",
		Input:     "testinput",
		Amount:    12.34,
		Timestamp: "test_timestamp",
		Sign:      "test_sign",
		Pubkey:    "test_pubkey",
	}

	// データの格納
	err = client.PutBlock(block)
	if err != nil {
		panic(err)
	}
	err = client.PutTransaction(tx)
	if err != nil {
		panic(err)
	}

	// keyとして必要なハッシュ値を取得
	blockHash := getHash(block)
	txHash := getHash(tx)

	// レコードの取得
	blockRecv, err := client.GetBlock(blockHash)
	if err != nil {
		panic(err)
	}
	txRecv, err := client.GetTransaction(txHash)
	if err != nil {
		panic(err)
	}

	// データの確認
	fmt.Printf("block:%v\n", blockRecv)
	fmt.Printf("transaction:%v\n", txRecv)

	// データの削除
	err = client.DeleteBlock(blockHash)
	err = client.DeleteTransaction(txHash)
}

```