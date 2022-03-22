package guildsapi

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	svc "github.com/InjectiveLabs/injective-guilds-service/api/gen/guilds_service"
	"github.com/InjectiveLabs/injective-guilds-service/internal/db"
	"github.com/InjectiveLabs/injective-guilds-service/internal/db/model"
	"github.com/InjectiveLabs/injective-guilds-service/internal/exchange"
	guildsprocess "github.com/InjectiveLabs/injective-guilds-service/internal/service/guilds-process"
	secp256k1 "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cosmtypes "github.com/cosmos/cosmos-sdk/types"
	log "github.com/xlab/suplog"
)

const expirationTimeLayout = "2006-01-02T15:04:05.000Z"

type GuildsAPI = svc.Service

type service struct {
	svc.Service
	exchangeProvider exchange.DataProvider
	dbSvc            db.DBService
	portfolioHelper  *guildsprocess.PortfolioHelper
	logger           log.Logger
	// TODO: Load as env var
	grants []string
}

func NewService(ctx context.Context, dbSvc db.DBService, exchangeProvider exchange.DataProvider) (GuildsAPI, error) {
	helper, err := guildsprocess.NewPortfolioHelper(ctx, dbSvc, exchangeProvider)
	if err != nil {
		return nil, err
	}

	return &service{
		dbSvc:            dbSvc,
		exchangeProvider: exchangeProvider,
		portfolioHelper:  helper,
		logger:           log.WithField("svc", "guilds_api"),
		grants: []string{
			// TODO: Double check with Peiyun for these message
			"/injective.exchange.v1beta1.MsgCreateSpotLimitOrder",
			"/injective.exchange.v1beta1.MsgCreateDerivativeLimitOrder",
			"/injective.exchange.v1beta1.MsgCancelDerivativeLimitOrder",
		},
	}, nil
}

// TODO: Refine error handling

func (s *service) GetAllGuilds(ctx context.Context) (res *svc.GetAllGuildsResult, err error) {
	guilds, err := s.dbSvc.ListAllGuilds(ctx)
	if err != nil {
		return nil, svc.MakeInternal(err)
	}

	var result []*svc.Guild
	for _, g := range guilds {
		result = append(result, modelGuildToResponse(g))
	}

	return &svc.GetAllGuildsResult{Guilds: result}, nil
}

func (s *service) GetSingleGuild(ctx context.Context, payload *svc.GetSingleGuildPayload) (res *svc.GetSingleGuildResult, err error) {
	guild, err := s.dbSvc.GetSingleGuild(ctx, payload.GuildID)
	if err != nil {
		return nil, svc.MakeInternal(err)
	}

	return &svc.GetSingleGuildResult{
		Guild: modelGuildToResponse(guild),
	}, nil
}

// Get members
func (s *service) GetGuildMembers(ctx context.Context, payload *svc.GetGuildMembersPayload) (res *svc.GetGuildMembersResult, err error) {
	isDefaultMember := false
	members, err := s.dbSvc.ListGuildMembers(
		ctx,
		model.MemberFilter{
			GuildID:         &payload.GuildID,
			IsDefaultMember: &isDefaultMember,
		},
	)
	if err != nil {
		return nil, svc.MakeInternal(err)
	}

	var result []*svc.GuildMember
	for _, m := range members {
		result = append(result, &svc.GuildMember{
			InjectiveAddress:     m.InjectiveAddress.String(),
			IsDefaultGuildMember: m.IsDefaultGuildMember,
		})
	}

	return &svc.GetGuildMembersResult{
		Members: result,
	}, nil
}

// Get master address of given guild
func (s *service) GetGuildMasterAddress(ctx context.Context, payload *svc.GetGuildMasterAddressPayload) (res *svc.GetGuildMasterAddressResult, err error) {
	guild, err := s.dbSvc.GetSingleGuild(ctx, payload.GuildID)
	if err != nil {
		return nil, svc.MakeInternal(err)
	}
	address := guild.MasterAddress.String()

	return &svc.GetGuildMasterAddressResult{
		MasterAddress: &address,
	}, nil
}

func (s *service) GetGuildDefaultMember(ctx context.Context, payload *svc.GetGuildDefaultMemberPayload) (res *svc.GetGuildDefaultMemberResult, err error) {
	isDefaultMember := true
	defaultMember, err := s.dbSvc.ListGuildMembers(
		ctx,
		model.MemberFilter{
			GuildID:         &payload.GuildID,
			IsDefaultMember: &isDefaultMember,
		},
	)
	if err != nil {
		return nil, svc.MakeInternal(err)
	}

	if len(defaultMember) == 0 {
		// 1 guild should has a default member, if not found then it should be internal error
		return nil, svc.MakeInternal(errors.New("default member not found"))
	}

	return &svc.GetGuildDefaultMemberResult{
		DefaultMember: &svc.GuildMember{
			InjectiveAddress:     defaultMember[0].InjectiveAddress.String(),
			IsDefaultGuildMember: defaultMember[0].IsDefaultGuildMember,
		},
	}, nil
}

func (s *service) isAddressQualified(ctx context.Context, guild *model.Guild, address string) (bool, error) {
	// TODO: Check balances
	// Currently, we can handle it on UI (discussed) => skip for now
	// check grants
	grants, err := s.exchangeProvider.GetGrants(ctx, address, guild.MasterAddress.String())
	if err != nil {
		return false, err
	}

	var msgToExpiration map[string]time.Time

	for _, g := range grants.Grants {
		t, err := time.Parse(expirationTimeLayout, g.Expiration)
		if err != nil {
			return false, fmt.Errorf("time parse err: %w", err)
		}

		msgToExpiration[g.Authorization.Msg] = t
	}

	// all expected grants must be provided
	now := time.Now()
	for _, expectedMsg := range s.grants {
		expiration, ok := msgToExpiration[expectedMsg]
		if !ok {
			return false, nil
		}

		if expiration.Before(now) {
			return false, nil
		}
	}

	return true, nil
}

// input are base64 strings
func (s *service) verifySigAndDeriveAddress(
	publicKey string,
	signature string,
	message string,
) (cosmtypes.AccAddress, error) {
	// parse
	// TODO: Add timestamp check
	messageBytes, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return nil, svc.MakeInvalidArg(errors.New("bad message"))
	}

	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return nil, svc.MakeInvalidArg(errors.New("bad signature"))
	}

	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, svc.MakeInvalidArg(errors.New("bad publicKey"))
	}

	// verify message, reference: https://pkg.go.dev/github.com/cosmos/cosmos-sdk/crypto/keyring
	pubKey := new(secp256k1.PubKey)
	if err := pubKey.Unmarshal(publicKeyBytes); err != nil {
		return nil, svc.MakeInvalidArg(errors.New("bad pubKey"))
	}

	if !pubKey.VerifySignature(messageBytes, signatureBytes) {
		return nil, svc.MakeInvalidArg(errors.New("cannot verify message and signature"))
	}

	// derive address
	return cosmtypes.AccAddress(pubKey.Address().Bytes()), nil
}

func (s *service) EnterGuild(ctx context.Context, payload *svc.EnterGuildPayload) (res *svc.EnterGuildResult, err error) {
	accAddress, err := s.verifySigAndDeriveAddress(payload.PublicKey, payload.Signature, payload.Message)
	if err != nil {
		return nil, err
	}

	guild, err := s.dbSvc.GetSingleGuild(ctx, payload.GuildID)
	if err != nil {
		return nil, svc.MakeInternal(fmt.Errorf("guild error: %w", err))
	}

	// check qualification
	qualified, err := s.isAddressQualified(ctx, guild, accAddress.String())
	if err != nil {
		return nil, svc.MakeInternal(fmt.Errorf("check qualification error: %w", err))
	}

	if !qualified {
		joinStatus := "address_not_qualified"
		return &svc.EnterGuildResult{
			JoinStatus: &joinStatus,
		}, nil
	}

	// get portfolio
	portfolio, err := s.portfolioHelper.CaptureSingleMemberPortfolio(
		ctx,
		guild,
		&model.GuildMember{
			GuildID:          guild.ID,
			InjectiveAddress: model.Address{AccAddress: accAddress},
		},
		true,
	)
	if err != nil {
		return nil, svc.MakeInternal(fmt.Errorf("capture portfolio error: %w", err))
	}

	// add to database
	err = s.dbSvc.AddMember(ctx, payload.GuildID, model.Address{AccAddress: accAddress})
	if err != nil {
		s.logger.WithError(err).Errorln("cannot add member")
		return nil, svc.MakeInternal(err)
	}

	// TODO: transaction
	err = s.dbSvc.AddAccountPortfolios(ctx, guild.ID.Hex(), []*model.AccountPortfolio{portfolio})
	if err != nil {
		// This account now joined guild, this error is not fatal, portfolio can be captured later
		s.logger.WithError(err).Warningln("cannot write account portfolio to db")
	}

	joinStatus := "success"
	return &svc.EnterGuildResult{
		JoinStatus: &joinStatus,
	}, nil
}

func (s *service) LeaveGuild(ctx context.Context, payload *svc.LeaveGuildPayload) (res *svc.LeaveGuildResult, err error) {
	accAddress, err := s.verifySigAndDeriveAddress(payload.PublicKey, payload.Signature, payload.Message)
	if err != nil {
		return nil, err
	}

	// remove member from database
	err = s.dbSvc.RemoveMember(ctx, payload.GuildID, model.Address{AccAddress: accAddress})
	if err != nil {
		return nil, svc.MakeInternal(err)
	}

	leaveStatus := "success"
	return &svc.LeaveGuildResult{
		LeaveStatus: &leaveStatus,
	}, nil
}

// GetGuildMarkets implements GetGuildMarkets.
func (s *service) GetGuildMarkets(ctx context.Context, payload *svc.GetGuildMarketsPayload) (res *svc.GetGuildMarketsResult, err error) {
	guild, err := s.dbSvc.GetSingleGuild(ctx, payload.GuildID)
	if err != nil {
		return nil, svc.MakeInternal(err)
	}

	var markets []*svc.Market
	for _, m := range guild.Markets {
		markets = append(markets, &svc.Market{
			MarketID:    m.MarketID.Hex(),
			IsPerpetual: m.IsPerpetual,
		})
	}

	return &svc.GetGuildMarketsResult{
		Markets: markets,
	}, nil
}

// GetAccountPortfolio implements GetAccountPortfolio.
func (s *service) GetAccountPortfolio(ctx context.Context, payload *svc.GetAccountPortfolioPayload) (res *svc.GetAccountPortfolioResult, err error) {
	address, err := cosmtypes.AccAddressFromBech32(payload.InjectiveAddress)
	if err != nil {
		return nil, svc.MakeInvalidArg(err)
	}

	portfolio, err := s.dbSvc.GetAccountPortfolio(ctx, payload.GuildID, model.Address{
		AccAddress: address,
	})
	if err != nil {
		return nil, svc.MakeInternal(err)
	}

	var (
		balances  []*svc.Balance
		updatedAt time.Time
	)

	for _, b := range portfolio.Balances {
		balances = append(balances, &svc.Balance{
			Denom:            b.Denom,
			TotalBalance:     b.TotalBalance.String(),
			AvailableBalance: b.AvailableBalance.String(),
			UnrealizedPnl:    b.UnrealizedPNL.String(),
			MarginHold:       b.MarginHold.String(),
		})
	}
	updatedAt = portfolio.UpdatedAt

	return &svc.GetAccountPortfolioResult{
		Data: &svc.SingleAccountPortfolio{
			InjectiveAddress: address.String(),
			Balances:         balances,
			UpdatedAt:        updatedAt.String(),
		},
	}, nil
}

func (s *service) GetAccountPortfolios(ctx context.Context, payload *svc.GetAccountPortfoliosPayload) (res *svc.GetAccountPortfoliosResult, err error) {
	address, err := cosmtypes.AccAddressFromBech32(payload.InjectiveAddress)
	if err != nil {
		return nil, svc.MakeInvalidArg(err)
	}

	portfolios, err := s.dbSvc.ListAccountPortfolios(ctx, payload.GuildID, model.Address{
		AccAddress: address,
	})
	if err != nil {
		return nil, svc.MakeInternal(err)
	}

	result := make([]*svc.SingleAccountPortfolio, 0)
	// expected result to be sort by timestamp
	for _, p := range portfolios {
		var balances []*svc.Balance
		for _, b := range p.Balances {
			balances = append(balances, &svc.Balance{
				Denom:            b.Denom,
				TotalBalance:     b.TotalBalance.String(),
				AvailableBalance: b.AvailableBalance.String(),
				UnrealizedPnl:    b.UnrealizedPNL.String(),
				MarginHold:       b.MarginHold.String(),
			})
		}

		result = append(result, &svc.SingleAccountPortfolio{
			InjectiveAddress: address.String(),
			Balances:         balances,
			UpdatedAt:        p.UpdatedAt.String(),
		})
	}

	return &svc.GetAccountPortfoliosResult{
		Portfolios: result,
	}, nil
}
