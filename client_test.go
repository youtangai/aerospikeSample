package utari

import (
	"testing"
)

func TestPutBalance(t *testing.T) {
	cli, err := NewAeroSpikeClient("127.0.0.1", 3000)
	if err != nil {
		t.Fatal(err)
	}

	err = cli.PutBalance("yota", 1234.56)
}

func TestGetGetBalanceByAddress(t *testing.T) {
	cli, err := NewAeroSpikeClient("127.0.0.1", 3000)
	if err != nil {
		t.Fatal(err)
	}

	balance, err := cli.GetBalanceByAddress("yota")
	if err != nil {
		t.Fatal(err)
	}
	expected := float64(1234.56)
	if balance != expected {
		t.Fatal(balance)
	}
}
