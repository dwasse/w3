package eth_test

import (
	"testing"

	"github.com/dwasse/w3"
	"github.com/dwasse/w3/module/eth"
	"github.com/dwasse/w3/rpctest"
	"github.com/ethereum/go-ethereum/common"
)

func TestStorageAt(t *testing.T) {
	tests := []rpctest.TestCase[common.Hash]{
		{
			Golden:  "get_storage_at",
			Call:    eth.StorageAt(w3.A("0x000000000000000000000000000000000000c0DE"), w3.H("0x0000000000000000000000000000000000000000000000000000000000000001"), nil),
			WantRet: w3.H("0x0000000000000000000000000000000000000000000000000000000000000042"),
		},
	}

	rpctest.RunTestCases(t, tests)
}
