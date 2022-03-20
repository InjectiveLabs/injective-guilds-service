package exchange

import (
	"context"
	"time"

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

type DerivativePosition struct {
}

type Grant struct {
	// "@type": "/cosmos.authz.v1beta1.GenericAuthorization",
	// "msg": "/injective.exchange.v1beta1.MsgCreateSpotLimitOrder"
	Authorization []struct {
		Type string `json:"@type"`
		Msg  string `json:"msg"`
	} `json:"authorization"`
	Expiration time.Time `json:"expiration"`
}

type DataProvider interface {
	GetDefaultSubaccountBalances(ctx context.Context, subaccount string) ([]*Balance, error)
	GetSpotOrders(ctx context.Context, subaccount string) ([]*SpotOrder, error)
	GetDerivativeOrders(ctx context.Context, subaccount string) ([]*DerivativeOrder, error)
	GetPositions(ctx context.Context, subaccount string) ([]*DerivativePosition, error)

	GetGrants(ctx context.Context, granter, grantee string) ([]*Grant, error)
	// close provider
	Close() error
}
