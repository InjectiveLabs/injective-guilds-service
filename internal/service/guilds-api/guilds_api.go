package guildsapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	svc "github.com/InjectiveLabs/injective-guilds-service/api/gen/guilds_service"
	"github.com/InjectiveLabs/injective-guilds-service/internal/config"
	"github.com/InjectiveLabs/injective-guilds-service/internal/db"
	"github.com/InjectiveLabs/injective-guilds-service/internal/db/model"
	"github.com/InjectiveLabs/injective-guilds-service/internal/exchange"
	guildsprocess "github.com/InjectiveLabs/injective-guilds-service/internal/service/guilds-process"
	cosmtypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/shopspring/decimal"
	log "github.com/xlab/suplog"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	expirationTimeLayout = "2006-01-02T15:04:05Z"
	ActionEnterGuild     = "enter-guild"
	ActionLeaveGuild     = "leave-guild"

	StatusQualified   = "qualified"
	StatusUnqualified = "unqualified"
)

type GuildsAPI = svc.Service

type qualificationResult struct {
	status string
	detail string
}
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
			"/injective.exchange.v1beta1.MsgCreateSpotLimitOrder",
			"/injective.exchange.v1beta1.MsgCreateSpotMarketOrder",
			"/injective.exchange.v1beta1.MsgCancelSpotOrder",
			"/injective.exchange.v1beta1.MsgBatchUpdateOrders",
			"/injective.exchange.v1beta1.MsgBatchCancelSpotOrders",
			"/injective.exchange.v1beta1.MsgDeposit",
			"/injective.exchange.v1beta1.MsgWithdraw",
			"/injective.exchange.v1beta1.MsgCreateDerivativeLimitOrder",
			"/injective.exchange.v1beta1.MsgCreateDerivativeMarketOrder",
			"/injective.exchange.v1beta1.MsgCancelDerivativeOrder",
			"/injective.exchange.v1beta1.MsgBatchUpdateOrders",
			"/injective.exchange.v1beta1.MsgBatchCancelDerivativeOrders",
			"/injective.exchange.v1beta1.MsgDeposit",
			"/injective.exchange.v1beta1.MsgWithdraw",
		},
	}, nil
}

// TODO: Refine error handling
func (s *service) GetAllGuilds(ctx context.Context) (res *svc.GetAllGuildsResult, err error) {
	guilds, err := s.dbSvc.ListAllGuilds(ctx)
	if err != nil {
		return nil, svc.MakeInternal(err)
	}

	var (
		result []*svc.Guild
		limit  = int64(1)
	)

	for _, g := range guilds {
		portfolios, err := s.dbSvc.ListGuildPortfolios(ctx, model.GuildPortfoliosFilter{
			GuildID: g.ID.Hex(),
			Limit:   &limit,
		})
		if err != nil {
			return nil, svc.MakeInternal(err)
		}

		var portfolio model.GuildPortfolio
		if len(portfolios) > 0 {
			portfolio = *portfolios[0]
		}

		guildID := g.ID.Hex()
		isDefaultMember := true
		defaultMember, err := s.dbSvc.ListGuildMembers(ctx, model.MemberFilter{
			GuildID:         &guildID,
			IsDefaultMember: &isDefaultMember,
		})

		if err != nil {
			return nil, svc.MakeInternal(err)
		}

		if len(defaultMember) == 0 {
			return nil, svc.MakeInternal(errors.New("guild has no default member"))
		}

		result = append(result, modelGuildToResponse(g, &portfolio, defaultMember[0]))
	}

	return &svc.GetAllGuildsResult{Guilds: result}, nil
}

func (s *service) GetSingleGuild(ctx context.Context, payload *svc.GetSingleGuildPayload) (res *svc.GetSingleGuildResult, err error) {
	guild, err := s.dbSvc.GetSingleGuild(ctx, payload.GuildID)
	if err != nil {
		return nil, svc.MakeInternal(err)
	}

	var (
		limit = int64(1)
		to    = time.Now()
	)

	portfolios, err := s.dbSvc.ListGuildPortfolios(ctx, model.GuildPortfoliosFilter{
		GuildID: guild.ID.Hex(),
		EndTime: &to,
		Limit:   &limit,
	})
	if err != nil {
		return nil, svc.MakeInternal(err)
	}

	var portfolio model.GuildPortfolio
	if len(portfolios) > 0 {
		portfolio = *portfolios[0]
	}

	guildID := guild.ID.Hex()
	isDefaultMember := true
	defaultMember, err := s.dbSvc.ListGuildMembers(ctx, model.MemberFilter{
		GuildID:         &guildID,
		IsDefaultMember: &isDefaultMember,
	})
	if err != nil {
		return nil, svc.MakeInternal(err)
	}

	if len(defaultMember) == 0 {
		return nil, svc.MakeInternal(errors.New("guild has no default member"))
	}

	return &svc.GetSingleGuildResult{
		Guild: modelGuildToResponse(guild, &portfolio, defaultMember[0]),
	}, nil
}

// Get members
func (s *service) GetGuildMembers(ctx context.Context, payload *svc.GetGuildMembersPayload) (res *svc.GetGuildMembersResult, err error) {
	members, err := s.dbSvc.ListGuildMembers(
		ctx,
		model.MemberFilter{
			GuildID: &payload.GuildID,
		},
	)
	if err != nil {
		return nil, svc.MakeInternal(err)
	}

	var result []*svc.GuildMember
	for _, m := range members {
		result = append(result, &svc.GuildMember{
			GuildID:              &payload.GuildID,
			InjectiveAddress:     m.InjectiveAddress.String(),
			IsDefaultGuildMember: m.IsDefaultGuildMember,
			Since:                m.Since.UnixMilli(),
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
		s.logger.WithError(err).Error("list guild members error")
		return nil, svc.MakeInternal(err)
	}

	if len(defaultMember) == 0 {
		s.logger.WithField("guildID", payload.GuildID).Error("default member not found")
		return nil, svc.MakeNotFound(errors.New("default member not found"))
	}

	return &svc.GetGuildDefaultMemberResult{
		DefaultMember: &svc.GuildMember{
			InjectiveAddress:     defaultMember[0].InjectiveAddress.String(),
			IsDefaultGuildMember: defaultMember[0].IsDefaultGuildMember,
		},
	}, nil
}

func (s *service) checkGrants(ctx context.Context, guild *model.Guild, address string) (*qualificationResult, error) {
	// Currently, we can handle it on UI (discussed) => skip for now
	// check grants
	grants, err := s.exchangeProvider.GetGrants(ctx, address, guild.MasterAddress.String())
	if err != nil {
		return nil, err
	}

	msgToExpiration := make(map[string]time.Time)

	for _, g := range grants.Grants {
		t, err := time.Parse(expirationTimeLayout, g.Expiration)
		if err != nil {
			return nil, fmt.Errorf("time parse err: %w", err)
		}

		msgToExpiration[g.Authorization.Msg] = t
	}

	// all expected grants must be provided
	now := time.Now()
	for _, expectedMsg := range s.grants {
		expiration, ok := msgToExpiration[expectedMsg]
		if !ok {
			return &qualificationResult{
				status: StatusUnqualified,
				detail: fmt.Sprintf("%s not granted", expectedMsg),
			}, nil
		}

		if expiration.Before(now) {
			return &qualificationResult{
				status: StatusUnqualified,
				detail: fmt.Sprintf("%s expired", expectedMsg),
			}, nil
		}
	}

	return &qualificationResult{
		status: StatusQualified,
	}, nil
}

func (s *service) getLatestGuildPortfolio(ctx context.Context, guildID string) (*model.GuildPortfolio, error) {
	limit := int64(1)
	portfolios, err := s.dbSvc.ListGuildPortfolios(ctx, model.GuildPortfoliosFilter{
		GuildID: guildID,
		Limit:   &limit,
	})
	if err != nil {
		return nil, err
	}

	if len(portfolios) == 0 {
		return nil, errors.New("portfolio not found")
	}

	return portfolios[0], nil
}

func (s *service) checkBalances(ctx context.Context, guild *model.Guild, snapshot *model.AccountPortfolio) (*qualificationResult, error) {
	if snapshot == nil {
		return nil, fmt.Errorf("no snapshot found to check")
	}

	// get latest portfolio to have price in usd, we will compare min price with this snapshot
	portfolio, err := s.getLatestGuildPortfolio(ctx, guild.ID.Hex())
	if err != nil {
		return nil, fmt.Errorf("get latest portfolio err: %w", err)
	}

	denomToUsdPrice := make(map[string]float64)
	for _, b := range portfolio.Balances {
		if _, isStableCoin := config.StableCoinDenoms[b.Denom]; isStableCoin {
			// price can be fluctuate, let's consider it 1$ for stable coins
			denomToUsdPrice[b.Denom] = 1
			continue
		}

		denomToUsdPrice[b.Denom] = b.PriceUSD
	}

	denomToDecimal := make(map[string]int)
	for _, market := range guild.Markets {
		if market.BaseTokenMeta != nil {
			denomToDecimal[market.BaseDenom] = market.BaseTokenMeta.Decimals
		}

		if market.QuoteTokenMeta != nil {
			denomToDecimal[market.QuoteDenom] = market.QuoteTokenMeta.Decimals
		}
	}

	denomToMinAmount := make(map[string]float64)
	for _, req := range guild.Requirements {
		denomToMinAmount[req.Denom] = req.MinAmountUSD
	}

	// we will use price usd from latest guild portfolio snapshot
	for _, b := range snapshot.Balances {
		availBalance, _ := decimal.NewFromString(b.AvailableBalance.String())
		dec, exist := denomToDecimal[b.Denom]
		if !exist {
			return nil, fmt.Errorf("failed check denom not belongs to market")
		}

		min, exist := denomToMinAmount[b.Denom]
		if !exist {
			return nil, fmt.Errorf("failed check denom not belongs to market")
		}

		priceUsd, exist := denomToUsdPrice[b.Denom]
		if !exist {
			return nil, fmt.Errorf("failed to check denom price in usd")
		}

		usdInDecimal := decimal.NewFromFloat(priceUsd)
		availBalanceFloat := availBalance.Shift(int32(dec)).Mul(usdInDecimal)
		if !availBalanceFloat.GreaterThanOrEqual(decimal.NewFromFloat(min)) {
			return &qualificationResult{
				status: StatusUnqualified,
				detail: fmt.Sprintf("%s has balance %s <= min %.2f", b.Denom, availBalanceFloat.String(), min),
			}, nil
		}
	}
	return &qualificationResult{
		status: StatusQualified,
	}, nil
}

func (s *service) checkAddressQualification(
	ctx context.Context,
	guild *model.Guild, portfolio *model.AccountPortfolio,
) (*qualificationResult, error) {
	// Total Balance == Available Balance
	balanceQualifyResult, err := s.checkBalances(ctx, guild, portfolio)
	if err != nil {
		return nil, err
	}

	if balanceQualifyResult.status != StatusQualified {
		return balanceQualifyResult, nil
	}

	return s.checkGrants(ctx, guild, portfolio.InjectiveAddress.String())
}

func (s *service) EnterGuild(ctx context.Context, payload *svc.EnterGuildPayload) (res *svc.EnterGuildResult, err error) {
	accAddress, msgInfo, err := verifySigAndExtractInfo(payload.PublicKey, payload.Signature, payload.Message)
	if err != nil {
		return nil, svc.MakeInvalidArg(err)
	}

	if msgInfo.Action != ActionEnterGuild {
		return nil, svc.MakeInvalidArg(fmt.Errorf("invalid action, should be %s", ActionEnterGuild))
	}

	guild, err := s.dbSvc.GetSingleGuild(ctx, payload.GuildID)
	if err != nil {
		return nil, svc.MakeInternal(fmt.Errorf("guild error: %w", err))
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

	// check qualification
	qualificationResult, err := s.checkAddressQualification(ctx, guild, portfolio)
	if err != nil {
		return nil, svc.MakeInternal(fmt.Errorf("check qualification error: %w", err))
	}

	if qualificationResult.status != StatusQualified {
		joinStatus := "not_qualified"
		return &svc.EnterGuildResult{
			JoinStatus: &joinStatus,
			Message:    &qualificationResult.detail,
		}, nil
	}

	// add to database
	err = s.dbSvc.AddMember(ctx, payload.GuildID, model.Address{AccAddress: accAddress}, false)
	if err != nil {
		s.logger.WithError(err).Errorln("cannot add member")
		return nil, svc.MakeInternal(err)
	}

	// TODO: transaction
	err = s.dbSvc.AddAccountPortfolios(ctx, []*model.AccountPortfolio{portfolio})
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
	accAddress, msgInfo, err := verifySigAndExtractInfo(payload.PublicKey, payload.Signature, payload.Message)
	if err != nil {
		return nil, svc.MakeInvalidArg(err)
	}

	if msgInfo.Action != ActionLeaveGuild {
		return nil, svc.MakeInvalidArg(fmt.Errorf("invalid action, should be %s", ActionLeaveGuild))
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

func (s *service) GetGuildPortfolios(
	ctx context.Context,
	payload *svc.GetGuildPortfoliosPayload,
) (res *svc.GetGuildPortfoliosResult, err error) {
	filter := model.GuildPortfoliosFilter{
		GuildID: payload.GuildID,
	}

	if payload.EndTime != nil {
		to := time.UnixMilli(*payload.EndTime)
		filter.EndTime = &to
	}

	var from time.Time
	if payload.StartTime != nil {
		from = time.UnixMilli(*payload.StartTime)
		filter.StartTime = &from
	}

	portfolios, err := s.dbSvc.ListGuildPortfolios(ctx, filter)
	if err != nil {
		return nil, svc.MakeInternal(fmt.Errorf("list guild portfolio err: %w", err))
	}

	result := make([]*svc.SingleGuildPortfolio, 0)
	// expected result to be sort by timestamp
	for _, p := range portfolios {
		if len(p.BankBalances) > 0 && p.BankBalances[0].Denom == config.DEMOM_INJ {
			p.Balances = addINJToBalances(p.Balances, p.BankBalances[0].Balance)
		}

		var balances []*svc.Balance
		for _, b := range p.Balances {
			balances = append(balances, &svc.Balance{
				Denom:            b.Denom,
				PriceUsd:         b.PriceUSD,
				TotalBalance:     b.TotalBalance.String(),
				AvailableBalance: b.AvailableBalance.String(),
				UnrealizedPnl:    b.UnrealizedPNL.String(),
				MarginHold:       b.MarginHold.String(),
			})
		}

		result = append(result, &svc.SingleGuildPortfolio{
			GuildID:   &payload.GuildID,
			Balances:  balances,
			UpdatedAt: p.UpdatedAt.UnixMilli(),
		})
	}

	return &svc.GetGuildPortfoliosResult{
		Portfolios: result,
	}, nil
}

func (s *service) GetAccountInfo(ctx context.Context, payload *svc.GetAccountInfoPayload) (res *svc.GetAccountInfoResult, err error) {
	address, err := cosmtypes.AccAddressFromBech32(payload.InjectiveAddress)
	if err != nil {
		return nil, svc.MakeInvalidArg(err)
	}

	members, err := s.dbSvc.ListGuildMembers(
		ctx,
		model.MemberFilter{
			InjectiveAddress: &model.Address{AccAddress: address},
		},
	)
	if err != nil {
		return nil, svc.MakeInternal(err)
	}

	if len(members) == 0 {
		return nil, svc.MakeNotFound(errors.New("member not found"))
	}
	guildID := members[0].GuildID.Hex()

	return &svc.GetAccountInfoResult{
		Data: &svc.GuildMember{
			GuildID:              &guildID,
			InjectiveAddress:     members[0].InjectiveAddress.String(),
			IsDefaultGuildMember: members[0].IsDefaultGuildMember,
			Since:                members[0].Since.UnixMilli(),
		},
	}, nil
}

func addINJToBalances(balance []*model.Balance, injAmount primitive.Decimal128) []*model.Balance {
	for _, b := range balance {
		if b.Denom == config.DEMOM_INJ {
			b.TotalBalance = sum(b.TotalBalance, injAmount)
			b.AvailableBalance = sum(b.AvailableBalance, injAmount)
			return balance
		}
	}

	// if not found then append inj denom
	balance = append(balance, &model.Balance{
		Denom:            config.DEMOM_INJ,
		TotalBalance:     injAmount,
		AvailableBalance: injAmount,
	})
	return balance
}

func (s *service) GetAccountPortfolio(ctx context.Context, payload *svc.GetAccountPortfolioPayload) (res *svc.GetAccountPortfolioResult, err error) {
	address, err := cosmtypes.AccAddressFromBech32(payload.InjectiveAddress)
	if err != nil {
		return nil, svc.MakeInvalidArg(err)
	}

	portfolio, err := s.dbSvc.GetAccountPortfolio(ctx, model.Address{
		AccAddress: address,
	})
	if err != nil {
		return nil, svc.MakeInternal(err)
	}

	var (
		balances  []*svc.Balance
		updatedAt time.Time
	)

	if len(portfolio.BankBalances) > 0 && portfolio.BankBalances[0].Denom == config.DEMOM_INJ {
		portfolio.Balances = addINJToBalances(portfolio.Balances, portfolio.BankBalances[0].Balance)
	}

	for _, b := range portfolio.Balances {
		aBalance := &svc.Balance{
			Denom:            b.Denom,
			PriceUsd:         b.PriceUSD,
			TotalBalance:     b.TotalBalance.String(),
			AvailableBalance: b.AvailableBalance.String(),
			UnrealizedPnl:    b.UnrealizedPNL.String(),
			MarginHold:       b.MarginHold.String(),
		}
		balances = append(balances, aBalance)
	}
	updatedAt = portfolio.UpdatedAt

	return &svc.GetAccountPortfolioResult{
		Data: &svc.SingleAccountPortfolio{
			InjectiveAddress: address.String(),
			Balances:         balances,
			UpdatedAt:        updatedAt.UnixMilli(),
		},
	}, nil
}

func (s *service) GetAccountPortfolios(ctx context.Context, payload *svc.GetAccountPortfoliosPayload) (res *svc.GetAccountPortfoliosResult, err error) {
	address, err := cosmtypes.AccAddressFromBech32(payload.InjectiveAddress)
	if err != nil {
		return nil, svc.MakeInvalidArg(err)
	}

	filter := model.AccountPortfoliosFilter{
		InjectiveAddress: model.Address{AccAddress: address},
	}

	if payload.EndTime != nil {
		endTime := time.UnixMilli(*payload.EndTime)
		filter.EndTime = &endTime
	}

	if payload.StartTime != nil {
		startTime := time.UnixMilli(*payload.StartTime)
		filter.StartTime = &startTime
	}

	portfolios, err := s.dbSvc.ListAccountPortfolios(ctx, filter)
	if err != nil {
		return nil, svc.MakeInternal(err)
	}

	result := make([]*svc.SingleAccountPortfolio, 0)
	// expected result to be sorted by timestamp
	for _, p := range portfolios {
		if len(p.BankBalances) > 0 && p.BankBalances[0].Denom == config.DEMOM_INJ {
			p.Balances = addINJToBalances(p.Balances, p.BankBalances[0].Balance)
		}

		var balances []*svc.Balance
		for _, b := range p.Balances {
			balances = append(balances, &svc.Balance{
				Denom:            b.Denom,
				PriceUsd:         b.PriceUSD,
				TotalBalance:     b.TotalBalance.String(),
				AvailableBalance: b.AvailableBalance.String(),
				UnrealizedPnl:    b.UnrealizedPNL.String(),
				MarginHold:       b.MarginHold.String(),
			})
		}

		result = append(result, &svc.SingleAccountPortfolio{
			InjectiveAddress: address.String(),
			Balances:         balances,
			UpdatedAt:        p.UpdatedAt.UnixMilli(),
		})
	}

	return &svc.GetAccountPortfoliosResult{
		Portfolios: result,
	}, nil
}
