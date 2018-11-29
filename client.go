package utari

import (
	aero "github.com/aerospike/aerospike-client-go"
)

const (
	AEROSPIKE_HOST        = "127.0.0.1"
	AEROSPIKE_PORT        = 3000
	AEROSPIKE_NAMESPACE   = "test"
	AEROSPIKE_TX_TABLE    = "TxTable"
	AEROSPIKE_BLOCL_TABLE = "BlockTable"
)

type aeroSpikeClient struct {
	client *aero.Client
}

type IAeroSpikeClinet interface {
	PutBlock(Block) error
	PutTransaction(Transaction) error
	GetBlock(string) (Block, error)
	GetTransaction(string) (Transaction, error)
	DeleteBlock(string) error
	DeleteTransaction(string) error
}

func NewAeroSpikeClient(host string, port int) (IAeroSpikeClinet, error) {
	cli, err := aero.NewClient(host, port)
	if err != nil {
		return nil, err
	}
	return aeroSpikeClient{
		client: cli,
	}, nil
}

// PutBlock は 引数のblockをaerospikeに格納するメソッドです
func (a aeroSpikeClient) PutBlock(block Block) error {
	// hash値の取得
	hash := getHash(block)

	// aerospike用のkey構造体を取得
	key, err := getBlockKey(hash)
	if err != nil {
		return err
	}

	// dataをbinmap(aerospikeに挿入可能な形)へ変換
	data := blockToBinMap(block)

	// データの格納
	err = a.client.Put(nil, key, data)
	if err != nil {
		return err
	}
	return nil
}

// PutTransaction は 引数Transactionをaerospikeに格納するメソッドです
func (a aeroSpikeClient) PutTransaction(tx Transaction) error {
	// hash値の取得
	hash := getHash(tx)

	// aerospike用のkey構造体を取得
	key, err := getTransactionKey(hash)
	if err != nil {
		return err
	}

	// dataをbinmap(aerospikeに挿入可能な形)へ変換
	data := transactionToBinMap(tx)

	// データの格納
	err = a.client.Put(nil, key, data)
	if err != nil {
		return err
	}
	return nil
}

// GetBlockは 引数であるhashをキーに Blockを取得するメソッドです
func (a aeroSpikeClient) GetBlock(hash string) (Block, error) {
	key, err := getBlockKey(hash)
	if err != nil {
		return Block{}, err
	}
	// レコードの取得
	record, err := a.client.Get(nil, key)
	if err != nil {
		return Block{}, err
	}

	// binmap to block
	block, err := binMapToBlock(record)
	if err != nil {
		return Block{}, err
	}

	return block, nil
}

// GetTransactionは 引数であるhashをキーに Transactionを取得するメソッドです
func (a aeroSpikeClient) GetTransaction(hash string) (Transaction, error) {
	key, err := getTransactionKey(hash)
	if err != nil {
		return Transaction{}, err
	}

	// レコードの取得
	record, err := a.client.Get(nil, key)
	if err != nil {
		return Transaction{}, err
	}

	// binmap to tx
	tx, err := binMapToTransaction(record)
	if err != nil {
		return Transaction{}, err
	}

	return tx, nil
}

// DeleteBlock は hashをキーに blockを削除するメソッドです
func (a aeroSpikeClient) DeleteBlock(hash string) error {
	key, err := getBlockKey(hash)
	_, err = a.client.Delete(nil, key)
	if err != nil {
		return err
	}
	return nil
}

// DeleteTransaction は hashをキーに transactionを削除するメソッドです
func (a aeroSpikeClient) DeleteTransaction(hash string) error {
	key, err := getTransactionKey(hash)
	_, err = a.client.Delete(nil, key)
	if err != nil {
		return err
	}
	return nil
}
