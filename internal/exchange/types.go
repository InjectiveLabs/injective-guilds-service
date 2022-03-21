package exchange

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// declare internal types, interface for this service
type Balance struct {
	Denom            string
	TotalBalance     primitive.Decimal128
	AvailableBalance primitive.Decimal128
}

// we only concern these fields in this service
type SpotOrder struct {
	OrderHash    string
	FeeRecipient string
}

type DerivativeOrder struct {
	OrderHash    string
	FeeRecipient string
}

// to calculate unrealized pnl
type DerivativePosition struct {
	MarketID   string
	Direction  string
	Quantity   primitive.Decimal128
	Margin     primitive.Decimal128
	EntryPrice primitive.Decimal128
}

type Grants struct {
	Grants []struct {
		Authorization struct {
			Type string `json:"@type"` // e.g "@type": "/cosmos.authz.v1beta1.GenericAuthorization",
			Msg  string `json:"msg"`   // e.g "msg": "/injective.exchange.v1beta1.MsgCreateSpotLimitOrder"
		} `json:"authorization"`
		Expiration string `json:"expiration"`
	} `json:"grants"`

	Pagination struct {
		NextKey string `json:"next_key"`
		Total   string `json:"total"`
	}
}

type DataProvider interface {
	GetSubaccountBalances(ctx context.Context, subaccount string) ([]*Balance, error)

	GetSpotOrders(ctx context.Context, marketIDs []string, subaccount string) ([]*SpotOrder, error)
	GetDerivativeOrders(ctx context.Context, marketIDs []string, subaccount string) ([]*DerivativeOrder, error)

	GetPositions(ctx context.Context, marketID string, subaccount string) ([]*DerivativePosition, error)
	GetGrants(ctx context.Context, granter, grantee string) (*Grants, error)
	// close provider
	Close() error
}
