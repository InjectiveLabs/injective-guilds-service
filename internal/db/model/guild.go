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

type TokenMeta struct {
	// Token full name
	Name string `bson:"name" json:"name"`
	// Token Ethereum contract address
	Address string `bson:"address" json:"address"`
	// Token symbol short name
	Symbol string `bson:"symbol" json:"symbol"`
	// URL to the logo image
	Logo *string `bson:"logo" json:"logo"`
	// Token decimals
	Decimals int `bson:"decimal" json:"decimal"`
	// Token metadata fetched timestamp in UNIX millis.
	UpdatedAt int64 `bson:"updated_at" json:"updated_at"`
}

type Guild struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"guild_id"`

	Name          string  `bson:"name" json:"name"`
	Description   string  `bson:"description" json:"description"`
	MasterAddress Address `bson:"master_address" json:"master_address"`

	SpotBaseRequirement        primitive.Decimal128 `bson:"spot_base_requirement" json:"spot_base_requirement"`
	SpotQuoteRequirement       primitive.Decimal128 `bson:"spot_quote_requirement" json:"spot_quote_requirement"`
	DerivativeQuoteRequirement primitive.Decimal128 `bson:"derivative_quote_requirement" json:"derivative_quote_requirement"`
	StakingRequirement         primitive.Decimal128 `bson:"staking_requirement" json:"staking_requirement"`

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

type AccountPortfolio struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	// GuildID stores guildID at the time this portfolio is captured
	// Which supports the case that an address can leave a guild and join another guild
	GuildID primitive.ObjectID `bson:"guild_id" json:"guild_id"`

	InjectiveAddress Address    `bson:"injective_address" json:"injective_address"`
	Balances         []*Balance `bson:"balances" json:"balances"`
	// timestamp when this gets update
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// This stores current guild of an address
type GuildMember struct {
	GuildID primitive.ObjectID `bson:"guild_id" json:"guild_id"`

	InjectiveAddress     Address `bson:"injective_address" json:"injective_address"`
	IsDefaultGuildMember bool    `bson:"is_default_guild_member" json:"is_default_guild_member"` // json might not need here
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
