package model

import (
	"encoding/binary"
	"time"

	cosmtypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

type Address struct {
	cosmtypes.AccAddress
}

func (a Address) IsZero() bool {
	return a.AccAddress.Empty()
}

func (a Address) MarshalJSON() ([]byte, error) {
	encoding := a.AccAddress.String()
	buf := make([]byte, 0, len(encoding)+2)
	buf = append(buf, '"')
	buf = append(buf, encoding...)
	buf = append(buf, '"')
	return buf, nil
}

func (a Address) MarshalBSONValue() (bsontype.Type, []byte, error) {
	buf := bsoncore.AppendString(nil, a.AccAddress.String())
	return bsontype.String, buf, nil
}

type Hash struct {
	common.Hash
}

func (h *Hash) To4Uint64() [4]uint64 {
	var res [4]uint64
	for i := 0; i < 4; i++ {
		res[i] = binary.LittleEndian.Uint64(h.Hash[8*i : 8*(i+1)])
	}
	return res
}

func (h *Hash) UnmarshalKey(key string) error {
	h.Hash = common.HexToHash(key)
	return nil
}

func (h Hash) MarshalKey() (key string, err error) {
	return h.Hex(), nil
}

var zeroHash = common.Hash{}

func (h Hash) IsZero() bool {
	return h.Hash == zeroHash
}

func (h Hash) MarshalJSON() ([]byte, error) {
	hex := h.Hash.Hex()
	buf := make([]byte, 0, len(hex)+2)
	buf = append(buf, '"')
	buf = append(buf, hex...)
	buf = append(buf, '"')
	return buf, nil
}

func (h *Hash) UnmarshalJSON(b []byte) error {
	hash := common.HexToHash(string(b)[1 : len(string(b))-1])
	*h = Hash{
		Hash: hash,
	}

	return nil
}

var _ bson.ValueMarshaler = Hash{}

func (h Hash) MarshalBSONValue() (bsontype.Type, []byte, error) {
	buf := bsoncore.AppendString(nil, h.Hash.Hex())
	return bsontype.String, buf, nil
}

func (h *Hash) UnmarshalBSONValue(t bsontype.Type, src []byte) error {
	if t != bsontype.String {
		return errors.Errorf("bsontype(%s) not allowed in BigNum.UnmarshalBSONValue", t.String())
	}

	v, _, ok := bsoncore.ReadString(src)
	if !ok {
		return errors.Errorf("bsoncore failed to read String")
	}

	*h = Hash{
		Hash: common.HexToHash(v),
	}
	return nil
}

type GuildMarket struct {
	MarketID    Hash `bson:"market_id" json:"market_id"`
	IsPerpetual bool `bson:"is_perpetual" json:"is_perpetual"`
}

type Guild struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"guild_id"`

	Name          string  `bson:"name" json:"name"`
	MasterAddress Address `bson:"master_address" json:"master_address"`

	SpotBaseRequirement        primitive.Decimal128 `bson:"spot_base_requirement" json:"spot_base_requirement"`
	SpotQuoteRequirement       primitive.Decimal128 `bson:"spot_quote_requirement" json:"spot_quote_requirement"`
	DerivativeQuoteRequirement primitive.Decimal128 `bson:"derivative_quote_requirement" json:"derivative_quote_requirement"`
	StakingRequirement         primitive.Decimal128 `bson:"staking_requirement" json:"staking_requirement"`

	Capacity int `bson:"capacity" json:"capacity"`
	// since number of markets is limited, we can embeded here:
	Markets []GuildMarket `bson:"markets" json:"markets"`
}

type AccountPortfolio struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	InjectiveAddress Address              `bson:"injective_address" json:"injective_address"`
	Denom            string               `bson:"denom" json:"denom"`
	TotalBalance     primitive.Decimal128 `bson:"total_balance" json:"total_balance"`
	AvailableBalance primitive.Decimal128 `bson:"available_balance" json:"available_balance"`
	UnrealizedPNL    primitive.Decimal128 `bson:"unrealized_pnl" json:"unrealized_pnl"`
	MarginHold       primitive.Decimal128 `bson:"margin_hold" json:"margin_hold"`

	// timestamp when this gets update
	UpdatedAt time.Time `bson:"time" json:"time"`
}

// Design for future, when 1 guild -> 10k+ guild member
type GuildMember struct {
	GuildID primitive.ObjectID `bson:"guild_id" json:"guild_id"`

	InjectiveAddress     Address `bson:"injective_address" json:"injective_address"`
	IsDefaultGuildMember bool    `bson:"is_default_guild_member" json:"is_default_guild_member"` // json might not need here
}
