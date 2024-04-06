// Package userop provides the base transaction object used throughout the stackup-bundler.
package userop

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type UserOp interface {
	GetPaymaster() common.Address
	GetFactory() common.Address
	GetFactoryData() []byte
	GetMaxGasAvailable() *big.Int
	GetMaxPrefund() *big.Int
	GetDynamicGasPrice(basefee *big.Int) *big.Int
	Pack() []byte
	PackForSignature() []byte
	GetUserOpHash(entryPoint common.Address, chainID *big.Int) common.Hash
	MarshalJSON() ([]byte, error)
	ToMap() (map[string]any, error)
}

type BasisUserOperation struct {
	Sender           common.Address `json:"sender"               mapstructure:"sender"               validate:"required"`
	Nonce            *big.Int       `json:"nonce"                mapstructure:"nonce"                validate:"required"`
	InitCode         []byte         `json:"initCode"             mapstructure:"initCode"             validate:"required"`
	CallData         []byte         `json:"callData"             mapstructure:"callData"             validate:"required"`
	PaymasterAndData []byte         `json:"paymasterAndData"     mapstructure:"paymasterAndData"     validate:"required"`
	Signature        []byte         `json:"signature"            mapstructure:"signature"            validate:"required"`
}

// GetPaymaster returns the address portion of PaymasterAndData if applicable. Otherwise, it returns the zero
// address.
func (op *BasisUserOperation) GetPaymaster() common.Address {
	if len(op.PaymasterAndData) < common.AddressLength {
		return common.HexToAddress("0x")
	}

	return common.BytesToAddress(op.PaymasterAndData[:common.AddressLength])
}

// GetFactory returns the address portion of InitCode if applicable. Otherwise, it returns the zero address.
func (op *BasisUserOperation) GetFactory() common.Address {
	if len(op.InitCode) < common.AddressLength {
		return common.HexToAddress("0x")
	}

	return common.BytesToAddress(op.InitCode[:common.AddressLength])
}

// GetFactoryData returns the data portion of InitCode if applicable. Otherwise, it returns an empty byte
// array.
func (op *BasisUserOperation) GetFactoryData() []byte {
	if len(op.InitCode) < common.AddressLength {
		return []byte{}
	}

	return op.InitCode[common.AddressLength:]
}
