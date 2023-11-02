package eth

import (
	"math/big"

	"github.com/dwasse/w3/internal/module"
	"github.com/dwasse/w3/w3types"
)

// GasPrice requests the current gas price in wei.
func GasPrice() w3types.CallerFactory[big.Int] {
	return module.NewFactory(
		"eth_gasPrice",
		nil,
		module.WithRetWrapper(module.HexBigRetWrapper),
	)
}
