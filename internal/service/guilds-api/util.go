package guildsapi

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	svc "github.com/InjectiveLabs/injective-guilds-service/api/gen/guilds_service"
	"github.com/InjectiveLabs/injective-guilds-service/internal/config"
	"github.com/InjectiveLabs/injective-guilds-service/internal/db/model"
	"github.com/InjectiveLabs/sdk-go/chain/crypto/ethsecp256k1"
	cosmtypes "github.com/cosmos/cosmos-sdk/types"
)

type MemberMessage struct {
	Action    string `json:"action"`
	ExpiredAt int64  `json:"expired_at"` // unix timestamp, second
}

func modelGuildToResponse(m *model.Guild, portfolio *model.GuildPortfolio, defaultMember *model.GuildMember) *svc.Guild {
	var (
		requirements    []*svc.Requirement
		balances        []*svc.Balance
		denomToUsdPrice = make(map[string]float64)
	)

	for _, b := range portfolio.Balances {
		balances = append(balances, &svc.Balance{
			Denom:            b.Denom,
			TotalBalance:     b.TotalBalance.String(),
			AvailableBalance: b.AvailableBalance.String(),
			UnrealizedPnl:    b.UnrealizedPNL.String(),
			MarginHold:       b.MarginHold.String(),
			PriceUsd:         b.PriceUSD,
		})

		denomToUsdPrice[b.Denom] = b.PriceUSD
	}

	for _, req := range m.Requirements {
		displayDecimal := config.DenomConfigs[req.Denom].DisplayDecimal
		priceUsd := denomToUsdPrice[req.Denom]
		roundedFloat := math.Ceil(req.MinAmountUSD*math.Pow10(displayDecimal)/priceUsd) / math.Pow10(displayDecimal)

		// IMPORTANT: We want to return price, so that FE and BE will be sync-ed for result
		requirements = append(requirements, &svc.Requirement{
			Denom:        req.Denom,
			MinAmountUsd: req.MinAmountUSD,
			MinAmount:    roundedFloat,
		})
	}

	var currentPortfolio *svc.SingleGuildPortfolio
	if len(balances) != 0 {
		currentPortfolio = &svc.SingleGuildPortfolio{
			Balances:  balances,
			UpdatedAt: portfolio.UpdatedAt.UnixMilli(),
		}
	}

	return &svc.Guild{
		ID:                   m.ID.Hex(),
		Name:                 m.Name,
		Description:          m.Description,
		MasterAddress:        m.MasterAddress.String(),
		Requirements:         requirements,
		Capacity:             m.Capacity,
		MemberCount:          m.MemberCount,
		CurrentPortfolio:     currentPortfolio,
		DefaultMemberAddress: defaultMember.InjectiveAddress.String(),
	}
}

// input are base64 strings
func verifySigAndExtractInfo(
	publicKeyBase64 string,
	signatureBase64 string,
	messageBase64 string,
) (cosmtypes.AccAddress, MemberMessage, error) {
	now := time.Now()

	signatureBytes, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		return nil, MemberMessage{}, fmt.Errorf("bad signature: %w", err)
	}

	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return nil, MemberMessage{}, fmt.Errorf("bad publicKey: %w", err)
	}

	messageBytes, err := base64.StdEncoding.DecodeString(messageBase64)
	if err != nil {
		return nil, MemberMessage{}, fmt.Errorf("bad message: %w", err)
	}

	var payload MemberMessage
	err = json.Unmarshal(messageBytes, &payload)
	if err != nil {
		return nil, MemberMessage{}, fmt.Errorf("cannot parse message: %w", err)
	}

	if payload.Action == "" || time.Unix(payload.ExpiredAt, 0).Before(now) {
		return nil, MemberMessage{}, errors.New("invalid action or expired timestamp")
	}

	// verify message, reference: https://pkg.go.dev/github.com/cosmos/cosmos-sdk/crypto/keyring
	pubKey := &ethsecp256k1.PubKey{
		Key: publicKeyBytes,
	}

	if !pubKey.VerifySignature(messageBytes, signatureBytes) {
		return nil, MemberMessage{}, fmt.Errorf("cannot verify message and signature")
	}

	// derive address
	return cosmtypes.AccAddress(pubKey.Address()), payload, nil
}
