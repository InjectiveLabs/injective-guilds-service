package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GuildMarket struct {
	MarketID    Hash `bson:"market_id" json:"market_id"`
	IsPerpetual bool `bson:"is_perpetual" json:"is_perpetual"`

	// Support information
	BaseDenom     string     `bson:"base_denom" json:"base_denom"`           // leave empty if this is perp market
	BaseTokenMeta *TokenMeta `bson:"base_token_meta" json:"base_token_meta"` // leave null if this is perp market

	QuoteDenom     string     `bson:"quote_denom" json:"quote_denom"`
	QuoteTokenMeta *TokenMeta `bson:"quote_token_meta" json:"quote_token_meta"`

	// these rate may be changed through a proposal
	MakerFeeRate primitive.Decimal128 `bson:"maker_fee_rate" json:"maker_fee_rate"`
	TakerFeeRate primitive.Decimal128 `bson:"taker_fee_rate" json:"taker_fee_rate"`
}

type DenomCoinID struct {
	Denom  string `bson:"denom" json:"denom"`
	CoinID string `bson:"coin_id" json:"coin_id"`
}

type TokenMeta struct {
	Name      string `bson:"name" json:"name"`
	Address   string `bson:"address" json:"address"`
	Symbol    string `bson:"symbol" json:"symbol"`
	Decimals  int    `bson:"decimal" json:"decimal"`
	UpdatedAt int64  `bson:"updated_at" json:"updated_at"`
}

type DenomRequirement struct {
	Denom        string  `bson:"denom" json:"denom"`
	MinAmountUSD float64 `bson:"min_amount_usd" json:"min_amount_usd"`
}

type Guild struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"guild_id"`

	Name          string  `bson:"name" json:"name"`
	Description   string  `bson:"description" json:"description"`
	MasterAddress Address `bson:"master_address" json:"master_address"`

	// Requirements are in USD
	Requirements []*DenomRequirement `bson:"denom_requirements" json:"denom_requirements"`

	Capacity    int `bson:"capacity" json:"capacity"`
	MemberCount int `bson:"member_count" json:"member_count"`

	// since number of markets is limited, we can embeded here:
	Markets []*GuildMarket `bson:"markets" json:"markets"`
}

// AccountPortfolio snapshot
type Balance struct {
	Denom    string  `bson:"denom" json:"denom"`
	PriceUSD float64 `bson:"price_usd" json:"price_usd"`

	TotalBalance     primitive.Decimal128 `bson:"total_balance" json:"total_balance"`
	AvailableBalance primitive.Decimal128 `bson:"available_balance" json:"available_balance"`
	UnrealizedPNL    primitive.Decimal128 `bson:"unrealized_pnl" json:"unrealized_pnl"`
	MarginHold       primitive.Decimal128 `bson:"margin_hold" json:"margin_hold"`
}

type BankBalance struct {
	Denom    string  `bson:"denom" json:"denom"`
	PriceUSD float64 `bson:"price_usd" json:"price_usd"`

	Balance primitive.Decimal128 `bson:"balance" json:"balance"`
}

type AccountPortfolio struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	// GuildID stores guildID at the time this portfolio is captured
	// Which supports the case that an address can leave a guild and join another guild
	GuildID primitive.ObjectID `bson:"guild_id" json:"guild_id"`

	InjectiveAddress Address `bson:"injective_address" json:"injective_address"`
	// Store default trading account's balance
	Balances []*Balance `bson:"balances" json:"balances"`
	// Store account's inj amount atm
	BankBalances []*BankBalance `bson:"bank_balances" json:"bank_balances"`
	// timestamp when this gets update
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type GuildPortfolio struct {
	GuildID     primitive.ObjectID `bson:"guild_id" json:"guild_id"`
	MemberCount int                `bson:"member_count" json:"member_count"`
	Balances    []*Balance         `bson:"balances" json:"balances"`
	// Store account's inj amount atm
	BankBalances []*BankBalance `bson:"bank_balances" json:"bank_balances"`
	// timestamp when this gets update
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// This stores current guild of an address
type GuildMember struct {
	GuildID primitive.ObjectID `bson:"guild_id" json:"guild_id"`

	Params               string    `bson:"params,omitempty" json:"params"`
	InjectiveAddress     Address   `bson:"injective_address" json:"injective_address"`
	IsDefaultGuildMember bool      `bson:"is_default_guild_member" json:"is_default_guild_member"` // json might not need here
	Since                time.Time `bson:"since" json:"since"`
}

func GetGuildDenoms(guild *Guild) []string {
	// get denoms from markets
	denomMap := make(map[string]bool)
	for _, m := range guild.Markets {
		denomMap[m.BaseDenom] = true
		denomMap[m.QuoteDenom] = true
	}

	denoms := make([]string, 0)
	for k := range denomMap {
		if k != "" {
			denoms = append(denoms, k)
		}
	}
	return denoms
}

func (a *AccountPortfolio) Copy() *AccountPortfolio {
	newPortfolio := &AccountPortfolio{
		GuildID:          a.GuildID,
		InjectiveAddress: a.InjectiveAddress,
		UpdatedAt:        a.UpdatedAt,
		Balances:         make([]*Balance, len(a.Balances)),
		BankBalances:     make([]*BankBalance, len(a.BankBalances)),
	}

	for i, b := range a.Balances {
		newPortfolio.Balances[i] = &Balance{
			Denom:    b.Denom,
			PriceUSD: b.PriceUSD,

			TotalBalance:     b.TotalBalance,
			AvailableBalance: b.AvailableBalance,
			UnrealizedPNL:    b.UnrealizedPNL,
			MarginHold:       b.MarginHold,
		}
	}

	for i, b := range a.BankBalances {
		newPortfolio.BankBalances[i] = &BankBalance{
			Denom:    b.Denom,
			PriceUSD: b.PriceUSD,
			Balance:  b.Balance,
		}
	}
	return newPortfolio
}
