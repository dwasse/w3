package eth_test

import (
	"math/big"
	"testing"

	"github.com/dwasse/w3"
	"github.com/dwasse/w3/module/eth"
	"github.com/dwasse/w3/rpctest"
)

func TestGasPrice(t *testing.T) {
	tests := []rpctest.TestCase[big.Int]{
		{
			Golden:  "gas_price",
			Call:    eth.GasPrice(),
			WantRet: *w3.I("0xc0fe"),
		},
	}

	rpctest.RunTestCases(t, tests)
}
