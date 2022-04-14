package model

import (
	"bytes"
	"encoding/binary"
	"time"

	cosmtypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
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

var emptyCosmosAddr = cosmtypes.AccAddress{}

func (a Address) IsEmpty() bool {
	return bytes.Equal(a.AccAddress, emptyCosmosAddr)
}

func (a *Address) UnmarshalBSONValue(t bsontype.Type, src []byte) error {
	if t != bsontype.String {
		return errors.Errorf("bsontype(%s) not allowed in Address.UnmarshalBSONValue", t.String())
	}

	v, _, ok := bsoncore.ReadString(src)
	if !ok {
		return errors.Errorf("bsoncore failed to read String")
	}

	accAddress, err := cosmtypes.AccAddressFromBech32(v)
	if err != nil {
		err = errors.Wrapf(err, "failed to unmarshal cosmos address from bech32: %s", v)
		return err
	}

	*a = Address{
		AccAddress: accAddress,
	}
	return nil
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

type MemberFilter struct {
	GuildID          *string
	IsDefaultMember  *bool
	InjectiveAddress *Address
}

type AccountPortfoliosFilter struct {
	InjectiveAddress Address
	StartTime        *time.Time
	EndTime          *time.Time
	Limit            *int64
}

type GuildPortfoliosFilter struct {
	GuildID   string
	StartTime *time.Time
	EndTime   *time.Time
	Limit     *int64
}
