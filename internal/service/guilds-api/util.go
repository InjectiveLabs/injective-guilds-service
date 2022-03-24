package guildsapi

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	svc "github.com/InjectiveLabs/injective-guilds-service/api/gen/guilds_service"
	"github.com/InjectiveLabs/injective-guilds-service/internal/db/model"
	"github.com/InjectiveLabs/sdk-go/chain/crypto/ethsecp256k1"
	cosmtypes "github.com/cosmos/cosmos-sdk/types"
)

type MemberMessage struct {
	Action    string `json:"action"`
	ExpiredAt int64  `json:"expired_at"` // unix timestamp, second
}

func modelGuildToResponse(m *model.Guild) *svc.Guild {
	return &svc.Guild{
		ID:            m.ID.Hex(),
		Name:          m.Name,
		Description:   m.Description,
		MasterAddress: m.MasterAddress.String(),
		Capacity:      m.Capacity,
		MemberCount:   m.MemberCount,
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
