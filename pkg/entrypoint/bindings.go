package entrypoint

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	entrypointV06 "github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/bindings/v06"
)
import entrypointV07 "github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/bindings/v07"

type Unknown interface {
	GetVersion() string
}

type BundlerClient struct {
	ver string
}

func (e *BundlerClient) GetVersion() string {
	return e.ver
}

// NewBundlerClient represent special version of entrypoint
func NewBundlerClient(ver string) *BundlerClient {
	return &BundlerClient{ver: ver}
}

func (e *BundlerClient) NewEntrypoint(address common.Address, backend bind.ContractBackend) (Entrypoint, error) {
	if f, ok := NewEntrypointFactory[e.ver]; ok {
		return f(address, backend)
	} else {
		panic(fmt.Sprintf("not supported version of '%s'", e.ver))
	}
}

type IStakeManagerDepositInfo interface {
	Unknown
}

type StakeManagerDepositInfoFactoryFunc func() IStakeManagerDepositInfo

var StakeManagerDepositInfoFactories = make(map[string]StakeManagerDepositInfoFactoryFunc)

func NewStakeManagerByVersion(version string) Entrypoint {
	if factory, ok := StakeManagerDepositInfoFactories[version]; ok {
		return factory()
	}
	return nil
}

type Entrypoint interface {
	Unknown
}

var NewEntrypointFactory = make(map[string]func(address common.Address, backend bind.ContractBackend) (Entrypoint, error))

// NewEntrypoint creates a new instance of Entrypoint, bound to a specific deployed contract.
func NewEntrypoint(address common.Address, backend bind.ContractBackend) (Entrypoint, error) {
	contract, err := bindEntrypoint(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Entrypoint{EntrypointCaller: EntrypointCaller{contract: contract}, EntrypointTransactor: EntrypointTransactor{contract: contract}, EntrypointFilterer: EntrypointFilterer{contract: contract}}, nil
}

func init() {
	StakeManagerDepositInfoFactories[entrypointV06.VERSION] = func() IStakeManagerDepositInfo {
		return &entrypointV06.IStakeManagerDepositInfo{}
	}
	StakeManagerDepositInfoFactories[entrypointV07.VERSION] = func() IStakeManagerDepositInfo {
		return &entrypointV07.IStakeManagerDepositInfo{}
	}

	NewEntrypointFactory[entrypointV06.VERSION] = func(address common.Address, backend bind.ContractBackend) (Entrypoint, error) {
		return entrypointV06.NewEntrypoint(address, backend)
	}

	NewEntrypointFactory[entrypointV06.VERSION] = func(address common.Address, backend bind.ContractBackend) (Entrypoint, error) {
		return entrypointV07.NewEntrypoint(address, backend)
	}
}
