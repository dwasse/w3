package eth

import (
	"github.com/dwasse/w3/internal/module"
	"github.com/dwasse/w3/w3types"
)

// ChainID requests the chains ID.
func ChainID() w3types.CallerFactory[uint64] {
	return module.NewFactory(
		"eth_chainId",
		nil,
		module.WithRetWrapper(module.HexUint64RetWrapper),
	)
}
