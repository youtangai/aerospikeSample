package utari

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

const (
	prefix = "utari"
	tag    = "utari-aerospike-client"
)

// Config は このパッケージのコンフィグを管理する構造体です
type Config struct {
	AerospikeNamespace  string `default:"test"`
	AerospikeBlockTable string `default:"BlockTable"`
	AerospikeTxTable    string `default:"TxTable"`
}

var (
	c Config
)

func init() {
	envconfig.MustProcess(prefix, &c)
	fmt.Printf("%s: namespace: %s, blocktable: %s, txtable: %s\n", tag, c.AerospikeNamespace, c.AerospikeBlockTable, c.AerospikeTxTable)
}

func getAerospikeNamespace() string {
	return c.AerospikeNamespace
}

func getAerospikeBlockTable() string {
	return c.AerospikeBlockTable
}

func getAerospikeTxTable() string {
	return c.AerospikeTxTable
}
