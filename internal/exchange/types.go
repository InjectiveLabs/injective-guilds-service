package exchange

import (
	"context"

	"github.com/shopspring/decimal"
	"google.golang.org/grpc"
)

// declare internal types, interface for this service
type Balance struct {
	Denom            string
	TotalBalance     decimal.Decimal
	AvailableBalance decimal.Decimal
}

// we only concern these fields in this service
type SpotOrder struct {
	MarketID string

	OrderHash    string
	FeeRecipient string
	OrderSide    string

	Price            decimal.Decimal
	UnfilledQuantity decimal.Decimal
}

type DerivativeOrder struct {
	MarketID     string
	OrderHash    string
	FeeRecipient string

	Margin decimal.Decimal
}

// to calculate unrealized pnl
type DerivativePosition struct {
	MarketID   string
	Direction  string
	Quantity   decimal.Decimal
	Margin     decimal.Decimal
	EntryPrice decimal.Decimal
	MarkPrice  decimal.Decimal
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

type BankAccountBalances struct {
	Balances []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	}
	Pagination struct {
		NextKey *string `json:"next_key"`
		Total   string  `json:"total"`
	}
}

type CoinPrice struct {
	ID                           string
	Symbol                       string
	Name                         string
	Image                        string
	CurrentPrice                 float64 `json:"current_price"`
	MarketCap                    float64 `json:"market_cap"`
	MarketCapRank                int     `json:"market_cap_rank"`
	TotalVolume                  float64 `json:"total_volume"`
	High24h                      float64 `json:"high_24h"`
	Low24h                       float64 `json:"low_24h"`
	PriceChange24h               float64 `json:"price_change_24h"`
	PriceChangePercentage24h     float64 `json:"price_change_percentage_24h"`
	MarketCapChange24h           float64 `json:"market_cap_change_24h"`
	MarketCapChangePercentage24h float64 `json:"market_cap_change_percentage_24h"`
	CirculatingSupply            float64 `json:"circulating_supply"`
	TotalSupply                  float64 `json:"total_supply"`
	MaxSupply                    float64 `json:"max_supply"`
	Ath                          float64
	AthChangePercentage          float64 `json:"ath_change_percentage"`
	AthDate                      string  `json:"ath_date"`
	Atl                          float64
	AtlChangePercentage          float64     `json:"atl_change_percentage"`
	AtlDate                      string      `json:"atl_date"`
	Roi                          interface{} `json:"roi"` // we don't use this
	LastUpdated                  string      `json:"last_updated"`
}

type CoinPriceResult struct {
	Data []*CoinPrice `json:"data"`
}

type DataProvider interface {
	GetSubaccountBalances(ctx context.Context, subaccount string) ([]*Balance, error)
	GetSpotOrders(ctx context.Context, marketIDs []string, subaccount string) ([]*SpotOrder, error)
	GetDerivativeOrders(ctx context.Context, marketIDs []string, subaccount string) ([]*DerivativeOrder, error)
	GetPositions(ctx context.Context, subaccount string) ([]*DerivativePosition, error)

	GetGrants(ctx context.Context, granter, grantee string) (*Grants, error)
	GetBankBalance(ctx context.Context, address string) (*BankAccountBalances, error)
	GetPriceUSD(ctx context.Context, coinIDs []string) ([]*CoinPrice, error)

	GetExchangeConn() *grpc.ClientConn
	// close provider
	Close() error
}
