package eth_test

import (
	"testing"

	"github.com/dwasse/w3/module/eth"
	"github.com/dwasse/w3/rpctest"
)

func TestChainID(t *testing.T) {
	tests := []rpctest.TestCase[uint64]{
		{
			Golden:  "chain_id",
			Call:    eth.ChainID(),
			WantRet: 1,
		},
	}

	rpctest.RunTestCases(t, tests)
}
