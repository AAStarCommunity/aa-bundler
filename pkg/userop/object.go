// Package userop provides the base transaction object used throughout the stackup-bundler.
package userop

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	address, _ = abi.NewType("address", "", nil)
	uint256, _ = abi.NewType("uint256", "", nil)
	bytes32, _ = abi.NewType("bytes32", "", nil)

	// UserOpPrimitives is the primitive ABI types for each UserOperation field.
	UserOpPrimitives = []abi.ArgumentMarshaling{
		{Name: "sender", InternalType: "Sender", Type: "address"},
		{Name: "nonce", InternalType: "Nonce", Type: "uint256"},
		{Name: "initCode", InternalType: "InitCode", Type: "bytes"},
		{Name: "callData", InternalType: "CallData", Type: "bytes"},
		{Name: "callGasLimit", InternalType: "CallGasLimit", Type: "uint256"},
		{Name: "verificationGasLimit", InternalType: "VerificationGasLimit", Type: "uint256"},
		{Name: "preVerificationGas", InternalType: "PreVerificationGas", Type: "uint256"},
		{Name: "maxFeePerGas", InternalType: "MaxFeePerGas", Type: "uint256"},
		{Name: "maxPriorityFeePerGas", InternalType: "MaxPriorityFeePerGas", Type: "uint256"},
		{Name: "paymasterAndData", InternalType: "PaymasterAndData", Type: "bytes"},
		{Name: "signature", InternalType: "Signature", Type: "bytes"},
	}

	// UserOpType is the ABI type of a UserOperation.
	UserOpType, _ = abi.NewType("tuple", "op", UserOpPrimitives)

	// UserOpArr is the ABI type for an array of UserOperations.
	UserOpArr, _ = abi.NewType("tuple[]", "ops", UserOpPrimitives)
)

// UserOperation represents an EIP-4337 style transaction for a smart contract account.
type UserOperation struct {
	Sender             common.Address `json:"sender"               mapstructure:"sender"               validate:"required"`
	Nonce              *big.Int       `json:"nonce"                mapstructure:"nonce"                validate:"required"`
	InitCode           []byte         `json:"initCode"             mapstructure:"initCode"             validate:"required"`
	CallData           []byte         `json:"callData"             mapstructure:"callData"             validate:"required"`
	PaymasterAndData   []byte         `json:"paymasterAndData"     mapstructure:"paymasterAndData"     validate:"required"`
	Signature          []byte         `json:"signature"            mapstructure:"signature"            validate:"required"`
	AccountGasLimits   [32]byte       `json:"accountGasLimits"     mapstructure:"accountGasLimits"     validate:"required"`
	PreVerificationGas *big.Int       `json:"preVerificationGas"   mapstructure:"preVerificationGas"   validate:"required"`
	GasFees            [32]byte       `json:"gasFees"              mapstructure:"gasFees"              validate:"required"`
}

// GetPaymaster returns the address portion of PaymasterAndData if applicable. Otherwise, it returns the zero
// address.
func (op *UserOperation) GetPaymaster() common.Address {
	if len(op.PaymasterAndData) < common.AddressLength {
		return common.HexToAddress("0x")
	}

	return common.BytesToAddress(op.PaymasterAndData[:common.AddressLength])
}

// GetFactory returns the address portion of InitCode if applicable. Otherwise, it returns the zero address.
func (op *UserOperation) GetFactory() common.Address {
	if len(op.InitCode) < common.AddressLength {
		return common.HexToAddress("0x")
	}

	return common.BytesToAddress(op.InitCode[:common.AddressLength])
}

// GetFactoryData returns the data portion of InitCode if applicable. Otherwise, it returns an empty byte
// array.
func (op *UserOperation) GetFactoryData() []byte {
	if len(op.InitCode) < common.AddressLength {
		return []byte{}
	}

	return op.InitCode[common.AddressLength:]
}

// GetMaxGasAvailable returns the max amount of gas that can be consumed by this UserOperation.
func (op *UserOperation) GetMaxGasAvailable() *big.Int {
	// TODO: Multiplier logic might change in v0.7
	//mul := big.NewInt(1)
	//paymaster := op.GetPaymaster()
	//if paymaster != common.HexToAddress("0x") {
	//	mul = big.NewInt(3)
	//}

	//return big.NewInt(0).Add(
	//	big.NewInt(0).Mul(op.VerificationGasLimit, mul),
	//	big.NewInt(0).Add(op.PreVerificationGas, op.CallGasLimit),
	//)

	panic("not implemented")
}

// GetMaxPrefund returns the max amount of wei required to pay for gas fees by either the sender or
// paymaster.
func (op *UserOperation) GetMaxPrefund() *big.Int {
	//return big.NewInt(0).Mul(op.GetMaxGasAvailable(), op.MaxFeePerGas)
	panic("not implemented")
}

// GetDynamicGasPrice returns the effective gas price paid by the UserOperation given a basefee. If basefee is
// nil, it will assume a value of 0.
func (op *UserOperation) GetDynamicGasPrice(basefee *big.Int) *big.Int {
	//bf := basefee
	//if bf == nil {
	//	bf = big.NewInt(0)
	//}
	//
	//gp := big.NewInt(0).Add(bf, op.MaxPriorityFeePerGas)
	//if gp.Cmp(op.MaxFeePerGas) == 1 {
	//	return op.MaxFeePerGas
	//}
	//return gp
	panic("not implemented")
}

// Pack returns a standard message of the userOp. This cannot be used to generate a userOpHash.
func (op *UserOperation) Pack() []byte {
	args := abi.Arguments{
		{Name: "UserOp", Type: UserOpType},
	}
	packed, _ := args.Pack(&struct {
		Sender             common.Address
		Nonce              *big.Int
		InitCode           []byte
		CallData           []byte
		PaymasterAndData   []byte
		Signature          []byte
		AccountGasLimits   [32]byte
		PreVerificationGas *big.Int
		GasFees            [32]byte
	}{
		op.Sender,
		op.Nonce,
		op.InitCode,
		op.CallData,
		op.PaymasterAndData,
		op.Signature,
		op.AccountGasLimits,
		op.PreVerificationGas,
		op.GasFees,
	})

	enc := hexutil.Encode(packed)
	enc = "0x" + enc[66:]
	return hexutil.MustDecode(enc)
}

// PackForSignature returns a minimal message of the userOp. This can be used to generate a userOpHash.
func (op *UserOperation) PackForSignature() []byte {
	args := abi.Arguments{
		{Name: "sender", Type: address},
		{Name: "nonce", Type: uint256},
		{Name: "hashInitCode", Type: bytes32},
		{Name: "hashCallData", Type: bytes32},
		{Name: "hashPaymasterAndData", Type: bytes32},
		{Name: "accountGasLimits", Type: uint256},
		{Name: "preVerificationGas", Type: uint256},
		{Name: "gasFees", Type: uint256},
	}
	packed, _ := args.Pack(
		op.Sender,
		op.Nonce,
		crypto.Keccak256Hash(op.InitCode),
		crypto.Keccak256Hash(op.CallData),
		crypto.Keccak256Hash(op.PaymasterAndData),
		op.AccountGasLimits,
		op.PreVerificationGas,
		op.GasFees,
	)

	return packed
}

// GetUserOpHash returns the hash of the userOp + entryPoint address + chainID.
func (op *UserOperation) GetUserOpHash(entryPoint common.Address, chainID *big.Int) common.Hash {
	return crypto.Keccak256Hash(
		crypto.Keccak256(op.PackForSignature()),
		common.LeftPadBytes(entryPoint.Bytes(), 32),
		common.LeftPadBytes(chainID.Bytes(), 32),
	)
}

// MarshalJSON returns a JSON encoding of the UserOperation.
func (op *UserOperation) MarshalJSON() ([]byte, error) {
	// Note: The bundler spec test requires the address portion of the initCode to include the checksum.
	ic := "0x"
	if fa := op.GetFactory(); fa != common.HexToAddress("0x") {
		ic = fmt.Sprintf("%s%s", fa, common.Bytes2Hex(op.GetFactoryData()))
	}

	return json.Marshal(&struct {
		Sender             string `json:"sender"`
		Nonce              string `json:"nonce"`
		InitCode           string `json:"initCode"`
		CallData           string `json:"callData"`
		PaymasterAndData   string `json:"paymasterAndData"`
		Signature          string `json:"signature"`
		AccountGasLimits   string `json:"accountGasLimits"`
		PreVerificationGas string `json:"preVerificationGas"`
		GasFees            string `json:"gasFees"`
	}{
		Sender:             op.Sender.String(),
		Nonce:              hexutil.EncodeBig(op.Nonce),
		InitCode:           ic,
		CallData:           hexutil.Encode(op.CallData),
		PaymasterAndData:   hexutil.Encode(op.PaymasterAndData),
		Signature:          hexutil.Encode(op.Signature),
		AccountGasLimits:   hexutil.Encode(op.AccountGasLimits[:]),
		PreVerificationGas: hexutil.EncodeBig(op.PreVerificationGas),
		GasFees:            hexutil.Encode(op.GasFees[:]),
	})
}

// ToMap returns the current UserOp struct as a map type.
func (op *UserOperation) ToMap() (map[string]any, error) {
	data, err := op.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var opData map[string]any
	if err := json.Unmarshal(data, &opData); err != nil {
		return nil, err
	}
	return opData, nil
}
