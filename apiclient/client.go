package apiclient

import (
	"errors"
)

var (
	ErrConnectRefused = errors.New("connection endppint refused")
)

type ChainClient interface {
	// Start()
	// Close()

	// Batch() Batch

	// SendETH(ctx context.Context, req cdata.SendETHRequest) (resp cdata.SendETHResponse, err error)
	// DeployFactory(ctx context.Context, req cdata.DeployFCRequest) (resp cdata.DeployFCResponse, err error)
	// CreateContracts(ctx context.Context, req cdata.CreateContractsRequest) (resp cdata.CreateContractsResponse, err error)
	// IssueSecurityToken(ctx context.Context, req cdata.IssueRequest) (resp cdata.IssueResponse, err error)
	// RegisterWalletComplianceService(ctx context.Context, req cdata.RegisterWalletRequest) (resp cdata.RegisterWalletResponse, err error)
	// TransferSecurityToken(ctx context.Context, req cdata.TransferRequest) (resp cdata.TransferResponse, err error)
	// BalanceOfSecurityToken(ctx context.Context, req cdata.BalanceOfRequest) (resp cdata.BalanceOfResponse, err error)
	// TotalSupplySecurityToken(ctx context.Context, req cdata.TotalSupplyRequest) (resp cdata.TotalSupplyResponse, err error)
	// GrantRole(ctx context.Context, req cdata.GrantRoleRequest) (resp cdata.GrantRoleResponse, err error)
	// HasRole(ctx context.Context, req cdata.HasRoleRequest) (resp cdata.HasRoleResponse, err error)
}
