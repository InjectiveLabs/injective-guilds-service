package db

import (
	"context"

	"github.com/InjectiveLabs/injective-guilds-service/internal/db/model"
)

type DBService interface {
	ListAllGuilds(ctx context.Context) ([]*model.Guild, error)
	GetSingleGuild(ctx context.Context, guildID string) (*model.Guild, error)

	// members
	ListGuildMembers(ctx context.Context, filter model.MemberFilter) ([]*model.GuildMember, error)
	AddMember(ctx context.Context, guildID string, address model.Address) error
	RemoveMember(ctx context.Context, guildID string, address model.Address) error

	// account portfolio
	GetAccountPortfolio(ctx context.Context, guildID string, address model.Address) (*model.AccountPortfolio, error)
	ListAccountPortfolios(ctx context.Context, guildID string, address model.Address) ([]*model.AccountPortfolio, error)
	AddAccountPortfolios(ctx context.Context, guildID string, portfolios []*model.AccountPortfolio) error

	Disconnect(ctx context.Context) error
}
