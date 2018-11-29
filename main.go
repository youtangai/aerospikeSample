package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	aero "github.com/aerospike/aerospike-client-go"
)

const (
	AEROSPIKE_HOST        = "127.0.0.1"
	AEROSPIKE_PORT        = 3000
	AEROSPIKE_NAMESPACE   = "test"
	AEROSPIKE_TX_TABLE    = "TxTable"
	AEROSPIKE_BLOCL_TABLE = "BlockTable"
)

type Block struct {
	Id         string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Version    int32    `protobuf:"varint,2,opt,name=version,proto3" json:"version,omitempty"`
	Prehash    string   `protobuf:"bytes,3,opt,name=prehash,proto3" json:"prehash,omitempty"`
	Merkleroot string   `protobuf:"bytes,4,opt,name=merkleroot,proto3" json:"merkleroot,omitempty"`
	Timestamp  string   `protobuf:"bytes,5,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Level      string   `protobuf:"bytes,6,opt,name=level,proto3" json:"level,omitempty"`
	Nonce      uint32   `protobuf:"varint,7,opt,name=nonce,proto3" json:"nonce,omitempty"`
	Size       int64    `protobuf:"varint,8,opt,name=size,proto3" json:"size,omitempty"`
	Txcount    int64    `protobuf:"varint,9,opt,name=txcount,proto3" json:"txcount,omitempty"`
	TxidList   []string `protobuf:"bytes,10,rep,name=txid_list,json=txidList,proto3" json:"txid_list,omitempty"`
}

type Transaction struct {
	Txid      string  `protobuf:"bytes,1,opt,name=txid,proto3" json:"txid,omitempty"`
	Output    string  `protobuf:"bytes,2,opt,name=output,proto3" json:"output,omitempty"`
	Input     string  `protobuf:"bytes,3,opt,name=input,proto3" json:"input,omitempty"`
	Amount    float64 `protobuf:"fixed64,4,opt,name=amount,proto3" json:"amount,omitempty"`
	Timestamp string  `protobuf:"bytes,5,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Sign      string  `protobuf:"bytes,6,opt,name=sign,proto3" json:"sign,omitempty"`
	Pubkey    string  `protobuf:"bytes,7,opt,name=pubkey,proto3" json:"pubkey,omitempty"`
}

func main() {
	cli, err := aero.NewClient(AEROSPIKE_HOST, AEROSPIKE_PORT)
	if err != nil {
		panic(err)
	}

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

	// hash値の取得
	hashKeyBlock := getHashForKey(block)
	hashKeyTx := getHashForKey(tx)

	// aerospike用のkey構造体を取得
	keyBlock, err := aero.NewKey(AEROSPIKE_NAMESPACE, AEROSPIKE_BLOCL_TABLE, hashKeyBlock)
	if err != nil {
		panic(err)
	}
	keyTx, err := aero.NewKey(AEROSPIKE_NAMESPACE, AEROSPIKE_TX_TABLE, hashKeyTx)
	if err != nil {
		panic(err)
	}

	// dataをbinmapへ変換
	dataBlock := blockToBinMap(block)
	dataTx := transactionToBinMap(tx)

	// データの挿入
	err = cli.Put(nil, keyBlock, dataBlock)
	if err != nil {
		panic(err)
	}
	err = cli.Put(nil, keyTx, dataTx)
	if err != nil {
		panic(err)
	}

	// レコードの取得
	recordBlock, err := cli.Get(nil, keyBlock)
	if err != nil {
		panic(err)
	}
	recordTx, err := cli.Get(nil, keyTx)
	if err != nil {
		panic(err)
	}

	// レコードをBlock型とTx型に変換
	blockRecv, err := binMapToBlock(recordBlock)
	if err != nil {
		panic(err)
	}
	txRecv, err := binMapToTransaction(recordTx)
	if err != nil {
		panic(err)
	}

	// データの確認
	fmt.Printf("block:%v\n", blockRecv)
	fmt.Printf("transaction:%v\n", txRecv)
}

// keyとして利用するhash値を取得する関数
func getHashForKey(v interface{}) string {
	// 構造体を[]byteに変換
	byteData, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	// ハッシュ関数にかける
	checksum := sha256.Sum256(byteData)
	// 文字列として取得する
	hash := fmt.Sprintf("%x", checksum)
	return hash
}

// Block構造体 to binmap
func blockToBinMap(b Block) aero.BinMap {
	return aero.BinMap{
		"Id":         b.Id,
		"Version":    b.Version,
		"Prehash":    b.Prehash,
		"Merkleroot": b.Merkleroot,
		"Timestamp":  b.Timestamp,
		"Level":      b.Level,
		"Nonce":      b.Nonce,
		"Size":       b.Size,
		"Txcount":    b.Txcount,
		"TxidList":   b.TxidList,
	}
}

// Transaction構造体 to binmap
func transactionToBinMap(t Transaction) aero.BinMap {
	return aero.BinMap{
		"Txid":      t.Txid,
		"Output":    t.Output,
		"Input":     t.Input,
		"Amount":    t.Amount,
		"Timestamp": t.Timestamp,
		"Sign":      t.Sign,
		"Pubkey":    t.Pubkey,
	}
}

// binmap to Block構造体
func binMapToBlock(record *aero.Record) (Block, error) {
	var block Block
	binMap := record.Bins

	// Idの型アサーション
	id, ok := binMap["Id"].(string)
	if !ok {
		return Block{}, fmt.Errorf("failed Id assertion")
	}
	block.Id = id

	// Versionの型アサーション
	version, ok := binMap["Version"].(int)
	if !ok {
		return Block{}, fmt.Errorf("failed Version assertion")
	}
	block.Version = int32(version)

	// Prehashの型アサーション
	prehash, ok := binMap["Prehash"].(string)
	if !ok {
		return Block{}, fmt.Errorf("failed Prehash assertion")
	}
	block.Prehash = prehash

	// Merklerootの型アサーション
	merkleroot, ok := binMap["Merkleroot"].(string)
	if !ok {
		return Block{}, fmt.Errorf("failed Merkleroot assertion")
	}
	block.Merkleroot = merkleroot

	// Timestampの型アサーション
	timestamp, ok := binMap["Timestamp"].(string)
	if !ok {
		return Block{}, fmt.Errorf("failed Timestamp assertion")
	}
	block.Timestamp = timestamp

	// Levelの型アサーション
	level, ok := binMap["Level"].(string)
	if !ok {
		return Block{}, fmt.Errorf("failed Level assertion")
	}
	block.Level = level

	// Nonceの型アサーション
	nonce, ok := binMap["Nonce"].(int)
	if !ok {
		return Block{}, fmt.Errorf("failed Nonce assertion")
	}
	block.Nonce = uint32(nonce)

	// Sizeの型アサーション
	size, ok := binMap["Size"].(int)
	if !ok {
		return Block{}, fmt.Errorf("failed Size assertion")
	}
	block.Size = int64(size)

	// Txcountの型アサーション
	txcount, ok := binMap["Txcount"].(int)
	if !ok {
		return Block{}, fmt.Errorf("failed Txcount assertion")
	}
	block.Txcount = int64(txcount)

	// TxidListの型アサーション
	var txidList []string
	// まずはスライスの型アサーション
	interfaceSlice, ok := binMap["TxidList"].([]interface{})
	if !ok {
		return Block{}, fmt.Errorf("failed TxidList assertion")
	}

	// スライスの中身を型アサーション
	for _, value := range interfaceSlice {
		txid, ok := value.(string)
		if !ok {
			return Block{}, fmt.Errorf("failed TxidList assertion")
		}
		txidList = append(txidList, txid)
	}
	block.TxidList = txidList

	return block, nil
}

// binmap to Transaction構造体
func binMapToTransaction(record *aero.Record) (Transaction, error) {
	var tx Transaction
	binMap := record.Bins

	// txidの型アサーション
	txid, ok := binMap["Txid"].(string)
	if !ok {
		return Transaction{}, fmt.Errorf("failed Txid assertion")
	}
	tx.Txid = txid

	// outputの型アサーション
	output, ok := binMap["Output"].(string)
	if !ok {
		return Transaction{}, fmt.Errorf("failed output assertion")
	}
	tx.Output = output

	// Inputの型アサーション
	input, ok := binMap["Input"].(string)
	if !ok {
		return Transaction{}, fmt.Errorf("failed input assertion")
	}
	tx.Input = input

	// Amountの型アサーション
	amount, ok := binMap["Amount"].(float64)
	if !ok {
		return Transaction{}, fmt.Errorf("failed Amount assertion")
	}
	tx.Amount = amount

	// Timestampの型アサーション
	timestamp, ok := binMap["Timestamp"].(string)
	if !ok {
		return Transaction{}, fmt.Errorf("failed Timestamp assertion")
	}
	tx.Timestamp = timestamp

	// Signの型アサーション
	sign, ok := binMap["Sign"].(string)
	if !ok {
		return Transaction{}, fmt.Errorf("failed sign assertion")
	}
	tx.Sign = sign

	// Pubkeyの型アサーション
	pubkey, ok := binMap["Pubkey"].(string)
	if !ok {
		return Transaction{}, fmt.Errorf("failed pubkey assertion")
	}
	tx.Pubkey = pubkey

	return tx, nil
}
