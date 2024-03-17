package service

import (
	"fmt"

	"github.com/tak1827/light-nft-indexer/util"
)

func (r *ListHolderNftTokenRequest) ValidateBasic() error {
	if err := util.ValidateEthAddress(r.GetWalletAddress()); err != nil {
		return fmt.Errorf("wallet address(=%s) is invalid format: %w", r.GetWalletAddress(), err)
	}
	if err := util.ValidateEthAddress(r.GetContractAddress()); err != nil {
		return fmt.Errorf("contract address(=%s) is invalid format: %w", r.GetContractAddress(), err)
	}

	return nil
}

func (r *ListHolderAllNftTokenRequest) ValidateBasic() error {
	if err := util.ValidateEthAddress(r.GetWalletAddress()); err != nil {
		return fmt.Errorf("wallet address(=%s) is invalid format: %w", r.GetWalletAddress(), err)
	}

	return nil
}
