package eth

import (
	"encoding/json"
	"math/big"

	"github.com/dwasse/w3/internal/module"
	"github.com/dwasse/w3/w3types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

// Call requests the output data of the given message at the given blockNumber.
// If blockNumber is nil, the output of the message at the latest known block is
// requested.
func Call(msg *w3types.Message, blockNumber *big.Int, overrides w3types.State) w3types.CallerFactory[[]byte] {
	args := []any{msg, module.BlockNumberArg(blockNumber)}
	if overrides != nil {
		args = append(args, overrides)
	}

	return module.NewFactory(
		"eth_call",
		args,
		module.WithArgsWrapper[[]byte](msgArgsWrapper),
		module.WithRetWrapper(module.HexBytesRetWrapper),
	)
}

// EstimateGas requests the estimated gas cost of the given message at the given
// blockNumber. If blockNumber is nil, the estimated gas cost of the message at
// the latest block is requested.
func EstimateGas(msg *w3types.Message, blockNumber *big.Int) w3types.CallerFactory[uint64] {
	return module.NewFactory(
		"eth_estimateGas",
		[]any{msg, module.BlockNumberArg(blockNumber)},
		module.WithArgsWrapper[uint64](msgArgsWrapper),
		module.WithRetWrapper(module.HexUint64RetWrapper),
	)
}

// AccessList requests the access list of the given message at the given
// blockNumber. If blockNumber is nil, the access list of the message at the
// latest block is requested.
func AccessList(msg *w3types.Message, blockNumber *big.Int) w3types.CallerFactory[AccessListResponse] {
	return module.NewFactory(
		"eth_createAccessList",
		[]any{msg, module.BlockNumberArg(blockNumber)},
		module.WithArgsWrapper[AccessListResponse](msgArgsWrapper),
	)
}

type AccessListResponse struct {
	AccessList types.AccessList
	GasUsed    uint64
}

// UnmarshalJSON implements the [json.Unmarshaler].
func (resp *AccessListResponse) UnmarshalJSON(data []byte) error {
	type accessListResponse struct {
		AccessList types.AccessList `json:"accessList"`
		GasUsed    hexutil.Uint64   `json:"gasUsed"`
	}

	var alResp accessListResponse
	if err := json.Unmarshal(data, &alResp); err != nil {
		return err
	}

	resp.AccessList = alResp.AccessList
	resp.GasUsed = uint64(alResp.GasUsed)
	return nil
}

func msgArgsWrapper(slice []any) ([]any, error) {
	msg := slice[0].(*w3types.Message)
	if msg.Input != nil || msg.Func == nil {
		return slice, nil
	}

	input, err := msg.Func.EncodeArgs(msg.Args...)
	if err != nil {
		return nil, err
	}
	msg.Input = input
	slice[0] = msg
	return slice, nil
}

// CallFunc requests the returns of Func fn at common.Address contract with the
// given args.
func CallFunc(fn w3types.Func, contract common.Address, args ...any) *CallFuncFactory {
	return &CallFuncFactory{fn: fn, contract: contract, args: args}
}

type CallFuncFactory struct {
	// args
	fn        w3types.Func
	contract  common.Address
	args      []any
	atBlock   *big.Int
	overrides w3types.State

	// returns
	result  []byte
	returns []any
}

func (f *CallFuncFactory) Returns(returns ...any) w3types.Caller {
	f.returns = returns
	return f
}

func (f *CallFuncFactory) AtBlock(blockNumber *big.Int) *CallFuncFactory {
	f.atBlock = blockNumber
	return f
}

func (f *CallFuncFactory) Overrides(overrides w3types.State) *CallFuncFactory {
	f.overrides = overrides
	return f
}

func (f *CallFuncFactory) CreateRequest() (rpc.BatchElem, error) {
	input, err := f.fn.EncodeArgs(f.args...)
	if err != nil {
		return rpc.BatchElem{}, err
	}

	args := []any{
		&w3types.Message{
			To:    &f.contract,
			Input: input,
		},
		module.BlockNumberArg(f.atBlock),
	}
	if len(f.overrides) > 0 {
		args = append(args, f.overrides)
	}

	return rpc.BatchElem{
		Method: "eth_call",
		Args:   args,
		Result: (*hexutil.Bytes)(&f.result),
	}, nil
}

func (f *CallFuncFactory) HandleResponse(elem rpc.BatchElem) error {
	if err := elem.Error; err != nil {
		return err
	}

	if err := f.fn.DecodeReturns(f.result, f.returns...); err != nil {
		return err
	}
	return nil
}
