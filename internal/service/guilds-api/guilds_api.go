package guildsapi

import (
	"context"
	"errors"

	svc "github.com/InjectiveLabs/injective-guilds-service/api/gen/guilds_service"
	"github.com/InjectiveLabs/injective-guilds-service/internal/db"
)

type GuildsAPI = svc.Service

type service struct {
	svc.Service
	dbSvc db.DBService
}

func NewService(dbSvc db.DBService) (GuildsAPI, error) {
	return &service{
		dbSvc: dbSvc,
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
	members, err := s.dbSvc.GetGuildMembers(ctx, payload.GuildID, false)
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
	defaultMember, err := s.dbSvc.GetGuildMembers(ctx, payload.GuildID, true)
	if err != nil {
		return nil, svc.MakeInternal(err)
	}

	if len(defaultMember) == 0 {
		return nil, svc.MakeInternal(errors.New("default member notfound"))
	}

	return &svc.GetGuildDefaultMemberResult{
		DefaultMember: &svc.GuildMember{
			InjectiveAddress:     defaultMember[0].InjectiveAddress.String(),
			IsDefaultGuildMember: defaultMember[0].IsDefaultGuildMember,
		},
	}, nil
}

func (s *service) EnterGuild(context.Context, *svc.EnterGuildPayload) (res *svc.EnterGuildResult, err error) {
	return &svc.EnterGuildResult{}, nil
}

func (s *service) LeaveGuild(context.Context, *svc.LeaveGuildPayload) (res *svc.LeaveGuildResult, err error) {
	return &svc.LeaveGuildResult{}, nil
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
func (s *service) GetAccountPortfolio(context.Context, *svc.GetAccountPortfolioPayload) (res *svc.GetAccountPortfolioResult, err error) {
	return &svc.GetAccountPortfolioResult{}, nil
}

// GetAccountPortfolios implements GetAccountPortfolios.
func (s *service) GetAccountPortfolios(context.Context, *svc.GetAccountPortfoliosPayload) (res *svc.GetAccountPortfoliosResult, err error) {
	return &svc.GetAccountPortfoliosResult{}, nil
}
