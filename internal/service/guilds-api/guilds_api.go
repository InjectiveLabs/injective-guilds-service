package guildsapi

import (
	"context"

	svc "github.com/InjectiveLabs/injective-guilds-service/api/gen/guilds_service"
)

type GuildsAPI = svc.Service

type service struct {
	svc.Service
}

func NewService() (GuildsAPI, error) {
	return &service{}, nil
}

func (s *service) GetAllGuilds(context.Context) (res *svc.GetAllGuildsResult, err error) {
	return &svc.GetAllGuildsResult{}, nil
}

func (s *service) GetSingleGuild(context.Context, *svc.GetSingleGuildPayload) (res *svc.GetSingleGuildResult, err error) {
	return &svc.GetSingleGuildResult{}, nil
}

// Get members
func (s *service) GetGuildMembers(context.Context, *svc.GetGuildMembersPayload) (res *svc.GetGuildMembersResult, err error) {
	return &svc.GetGuildMembersResult{}, nil
}

// Get master address of given guild
func (s *service) GetGuildMasterAddress(context.Context, *svc.GetGuildMasterAddressPayload) (res *svc.GetGuildMasterAddressResult, err error) {
	return &svc.GetGuildMasterAddressResult{}, nil
}

// GetGuildDefaultMember implements GetGuildDefaultMember.
func (s *service) GetGuildDefaultMember(context.Context, *svc.GetGuildDefaultMemberPayload) (res *svc.GetGuildDefaultMemberResult, err error) {
	return &svc.GetGuildDefaultMemberResult{}, nil
}

// EnterGuild implements EnterGuild.
func (s *service) EnterGuild(context.Context, *svc.EnterGuildPayload) (res *svc.EnterGuildResult, err error) {
	return &svc.EnterGuildResult{}, nil
}

// LeaveGuild implements LeaveGuild.
func (s *service) LeaveGuild(context.Context, *svc.LeaveGuildPayload) (res *svc.LeaveGuildResult, err error) {
	return &svc.LeaveGuildResult{}, nil
}

// GetGuildMarkets implements GetGuildMarkets.
func (s *service) GetGuildMarkets(context.Context, *svc.GetGuildMarketsPayload) (res *svc.GetGuildMarketsResult, err error) {
	return &svc.GetGuildMarketsResult{}, nil
}

// GetAccountPortfolio implements GetAccountPortfolio.
func (s *service) GetAccountPortfolio(context.Context, *svc.GetAccountPortfolioPayload) (res *svc.GetAccountPortfolioResult, err error) {
	return &svc.GetAccountPortfolioResult{}, nil
}

// GetAccountPortfolios implements GetAccountPortfolios.
func (s *service) GetAccountPortfolios(context.Context, *svc.GetAccountPortfoliosPayload) (res *svc.GetAccountPortfoliosResult, err error) {
	return &svc.GetAccountPortfoliosResult{}, nil
}
