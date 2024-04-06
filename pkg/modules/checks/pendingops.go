package checks

import (
	"fmt"
	"math/big"

	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

var (
	minPriceBump                = int64(10)
	ErrReplacementOpUnderpriced = fmt.Errorf(
		"pending ops: replacement op must increase maxFeePerGas and MaxPriorityFeePerGas by >= %d%%",
		minPriceBump,
	)
)

// calcNewThresholds returns new threshold values where newFee = oldFee  * (100 + minPriceBump) / 100.
func calcNewThresholds(cap *big.Int, tip *big.Int) (newCap *big.Int, newTip *big.Int) {
	a := big.NewInt(100 + minPriceBump)
	aFeeCap := big.NewInt(0).Mul(a, cap)
	aTip := big.NewInt(0).Mul(a, tip)

	b := big.NewInt(100)
	newCap = aFeeCap.Div(aFeeCap, b)
	newTip = aTip.Div(aTip, b)

	return newCap, newTip
}

// ValidatePendingOps checks the pending UserOperations by the same sender and only passes if:
//
//  1. Sender doesn't have another UserOp already present in the pool.
//  2. It replaces an existing UserOp with same nonce and higher fee.
func ValidatePendingOps(
	op *userop.UserOp,
	penOps []*userop.UserOp,
) error {
	if len(penOps) > 0 {
		var oldOp *userop.UserOp
		for _, penOp := range penOps {
			if op.Nonce.Cmp(penOp.Nonce) == 0 {
				oldOp = penOp
			}
		}

		if oldOp != nil {
			newMf, newMpf := calcNewThresholds(oldOp.MaxFeePerGas, oldOp.MaxPriorityFeePerGas)

			if op.MaxFeePerGas.Cmp(newMf) < 0 || op.MaxPriorityFeePerGas.Cmp(newMpf) < 0 {
				return ErrReplacementOpUnderpriced
			}
		}
	}
	return nil
}
