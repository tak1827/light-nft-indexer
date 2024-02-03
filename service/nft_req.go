package service

import (
	"errors"
	"fmt"

	"github.com/tak1827/light-nft-indexer/util"
)

const (
	MAX_CONTENTS_SIZE = 64
)

var (
	ErrInvalidSequence = errors.New("err invalid sequence")
)

func (r *GetNftContractRequest) ValidateBasic() error {
	if err := util.ValidateEthAddress(r.GetContractAddress()); err != nil {
		return fmt.Errorf("contract address(=%s) is invalid format: %w", r.GetContractAddress(), err)
	}

	return nil
}

func (r *GetNftTokenRequest) ValidateBasic() error {
	if err := util.ValidateEthAddress(r.GetContractAddress()); err != nil {
		return fmt.Errorf("contract address(=%s) is invalid format: %w", r.GetContractAddress(), err)
	}

	if err := util.ValidateOnlyNumber(r.GetTokenId()); err != nil {
		return fmt.Errorf("token id(=%s) is invalid format: %w", r.GetTokenId(), err)
	}

	return nil
}

func (r *ListAllNftTokenRequest) ValidateBasic() error {
	if err := util.ValidateEthAddress(r.GetContractAddress()); err != nil {
		return fmt.Errorf("contract address(=%s) is invalid format: %w", r.GetContractAddress(), err)
	}

	return nil
}

// func (r *GenerateTransferSigRequest) ValidateBasic() error {
// 	if err := util.ValidateEthAddress(r.GetFrom()); err != nil {
// 		return errors.Wrapf(err, "from(=%s) is invalid format", r.GetFrom())
// 	}

// 	if err := util.ValidateEthAddress(r.GetTo()); err != nil {
// 		return errors.Wrapf(err, "to(=%s) is invalid format", r.GetTo())
// 	}

// 	if err := util.ValidateOnlyAlphaNumber(r.GetId()); err != nil {
// 		return errors.Wrapf(err, "id(=%s) should be number string", r.GetId())
// 	}

// 	if err := util.ValidateOnlyAlphaNumber(r.GetAmount()); err != nil {
// 		return errors.Wrapf(err, "amount(=%s) should be number string", r.GetAmount())
// 	}

// 	return nil
// }
