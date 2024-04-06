package useropV07

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
	"math/big"
)

var (
	address, _ = abi.NewType("address", "", nil)
	uint256, _ = abi.NewType("uint256", "", nil)
	bytes32, _ = abi.NewType("bytes32", "", nil)

	// PackedUserOpPrimitives is the primitive ABI types for each UserOp field.
	PackedUserOpPrimitives = []abi.ArgumentMarshaling{
		{Name: "sender", InternalType: "Sender", Type: "address"},
		{Name: "nonce", InternalType: "Nonce", Type: "uint256"},
		{Name: "initCode", InternalType: "InitCode", Type: "bytes"},
		{Name: "callData", InternalType: "CallData", Type: "bytes"},
		{Name: "accountGasLimits", InternalType: "AccountGasLimits", Type: "bytes32"},
		{Name: "preVerificationGas", InternalType: "PreVerificationGas", Type: "uint256"},
		{Name: "gasFees", InternalType: "GasFees", Type: "bytes32"},
		{Name: "paymasterAndData", InternalType: "PaymasterAndData", Type: "bytes"},
		{Name: "signature", InternalType: "Signature", Type: "bytes"},
	}

	// UserOpType is the ABI type of a UserOp.
	UserOpType, _ = abi.NewType("tuple", "op", PackedUserOpPrimitives)

	// UserOpArr is the ABI type for an array of UserOperations.
	UserOpArr, _ = abi.NewType("tuple[]", "ops", PackedUserOpPrimitives)
)

// PackedUserOperation represents an EIP-4337 style transaction for a smart contract account. (V0.7)
type PackedUserOperation struct {
	userop.BasisUserOperation
	AccountGasLimits   [32]byte `json:"accountGasLimits"   mapstructure:"accountGasLimits"   validate:"required"`
	PreVerificationGas *big.Int `json:"preVerificationGas" mapstructure:"preVerificationGas" validate:"required"`
	GasFees            [32]byte `json:"gasFees"            mapstructure:"gasFees"            validate:"required"`
}

// GetMaxGasAvailable implement UserOp interface
func (op *PackedUserOperation) GetMaxGasAvailable() *big.Int {
	panic("implement me")
}

// GetMaxPrefund implement UserOp interface
func (op *PackedUserOperation) GetMaxPrefund() *big.Int {
	panic("implement me")
}

// GetDynamicGasPrice implement UserOp interface
func (op *PackedUserOperation) GetDynamicGasPrice(basefee *big.Int) *big.Int {
	panic("implement me")
}

// Pack implement UserOp interface
func (op *PackedUserOperation) Pack() []byte {
	panic("implement me")
}

// PackForSignature implement UserOp interface
func (op *PackedUserOperation) PackForSignature() []byte {
	panic("implement me")
}

// GetUserOpHash implement UserOp interface
func (op *PackedUserOperation) GetUserOpHash(entryPoint common.Address, chainID *big.Int) common.Hash {
	panic("implement me")
}

// MarshalJSON implement UserOp interface
func (op *PackedUserOperation) MarshalJSON() ([]byte, error) {
	panic("implement me")
}

// ToMap implement UserOp interface
func (op *PackedUserOperation) ToMap() (map[string]interface{}, error) {
	panic("implement me")
}
