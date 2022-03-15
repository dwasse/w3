package eth

import (
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/lmittmann/w3/core"
)

// Call requests the output data of the given message.
func Call(msg ethereum.CallMsg) *CallFactory {
	return &CallFactory{msg: msg}
}

type CallFactory struct {
	// args
	msg       ethereum.CallMsg
	atBlock   *big.Int
	overrides AccountOverrides

	// returns
	result  hexutil.Bytes
	returns *[]byte
}

func (f *CallFactory) AtBlock(blockNumber *big.Int) *CallFactory {
	f.atBlock = blockNumber
	return f
}

func (f *CallFactory) Overrides(overrides AccountOverrides) *CallFactory {
	f.overrides = overrides
	return f
}

func (f *CallFactory) Returns(output *[]byte) core.Caller {
	f.returns = output
	return f
}

// CreateRequest implements the core.RequestCreator interface.
func (f *CallFactory) CreateRequest() (rpc.BatchElem, error) {
	if f.overrides != nil {
		return rpc.BatchElem{
			Method: "eth_call",
			Args:   []any{toCallArg(f.msg), toBlockNumberArg(f.atBlock), f.overrides},
			Result: &f.result,
		}, nil
	}
	return rpc.BatchElem{
		Method: "eth_call",
		Args:   []any{toCallArg(f.msg), toBlockNumberArg(f.atBlock)},
		Result: &f.result,
	}, nil
}

// HandleResponse implements the core.ResponseHandler interface.
func (f *CallFactory) HandleResponse(elem rpc.BatchElem) error {
	if err := elem.Error; err != nil {
		return err
	}
	*f.returns = []byte(f.result)
	return nil
}

// CallFunc requests the returns of Func fn at common.Address contract with the
// given args.
func CallFunc(fn core.Func, contract common.Address, args ...any) *CallFuncFactory {
	return &CallFuncFactory{fn: fn, contract: contract, args: args}
}

type CallFuncFactory struct {
	// args
	fn        core.Func
	contract  common.Address
	args      []any
	from      *common.Address
	atBlock   *big.Int
	overrides AccountOverrides

	// returns
	result  hexutil.Bytes
	returns []any
}

func (f *CallFuncFactory) AtBlock(blockNumber *big.Int) *CallFuncFactory {
	f.atBlock = blockNumber
	return f
}

func (f *CallFuncFactory) Returns(returns ...any) core.Caller {
	f.returns = returns
	return f
}

func (f *CallFuncFactory) From(from common.Address) *CallFuncFactory {
	f.from = &from
	return f
}

func (f *CallFuncFactory) Overrides(overrides AccountOverrides) *CallFuncFactory {
	f.overrides = overrides
	return f
}

// CreateRequest implements the core.RequestCreator interface.
func (f *CallFuncFactory) CreateRequest() (rpc.BatchElem, error) {
	input, err := f.fn.EncodeArgs(f.args...)
	if err != nil {
		return rpc.BatchElem{}, err
	}

	msg := ethereum.CallMsg{
		To:   &f.contract,
		Data: input,
	}
	if f.from != nil {
		msg.From = *f.from
	}
	if f.overrides != nil {
		return rpc.BatchElem{
			Method: "eth_call",
			Args:   []any{toCallArg(msg), toBlockNumberArg(f.atBlock), f.overrides},
			Result: &f.result,
		}, nil
	}
	return rpc.BatchElem{
		Method: "eth_call",
		Args:   []any{toCallArg(msg), toBlockNumberArg(f.atBlock)},
		Result: &f.result,
	}, nil
}

// HandleResponse implements the core.ResponseHandler interface.
func (f *CallFuncFactory) HandleResponse(elem rpc.BatchElem) error {
	if err := elem.Error; err != nil {
		return err
	}
	output := []byte(f.result)
	if err := f.fn.DecodeReturns(output, f.returns...); err != nil {
		return err
	}
	return nil
}

func toCallArg(msg ethereum.CallMsg) any {
	arg := map[string]any{
		"to": msg.To,
	}
	if msg.From.Hash().Big().Sign() > 0 {
		arg["from"] = msg.From
	}
	if len(msg.Data) > 0 {
		arg["data"] = hexutil.Bytes(msg.Data)
	}
	if msg.Value != nil {
		arg["value"] = (*hexutil.Big)(msg.Value)
	}
	if msg.Gas != 0 {
		arg["gas"] = hexutil.Uint64(msg.Gas)
	}
	if msg.GasPrice != nil {
		arg["gasPrice"] = (*hexutil.Big)(msg.GasPrice)
	}
	return arg
}
