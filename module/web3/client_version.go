package web3

import (
	"github.com/dwasse/w3/internal/module"
	"github.com/dwasse/w3/w3types"
)

// ClientVersion requests the endpoints client version.
func ClientVersion() w3types.CallerFactory[string] {
	return module.NewFactory[string](
		"web3_clientVersion",
		nil,
	)
}
