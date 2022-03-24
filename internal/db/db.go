package db

import (
	"context"

	"github.com/InjectiveLabs/injective-guilds-service/internal/db/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DBService interface {
	ListAllGuilds(ctx context.Context) ([]*model.Guild, error)
	GetSingleGuild(ctx context.Context, guildID string) (*model.Guild, error)
	ListGuildPortfolios(ctx context.Context, filter model.GuildPortfoliosFilter) ([]*model.GuildPortfolio, error)
	// TODO: *primitive.ObjectID -> string
	AddGuild(ctx context.Context, guild *model.Guild) (*primitive.ObjectID, error)
	AddGuildPortfolios(ctx context.Context, portfolios []*model.GuildPortfolio) error
	DeleteGuild(ctx context.Context, guildID string) error

	// members
	ListGuildMembers(ctx context.Context, filter model.MemberFilter) ([]*model.GuildMember, error)
	AddMember(ctx context.Context, guildID string, address model.Address, isDefaultMember bool) error
	RemoveMember(ctx context.Context, guildID string, address model.Address) error
	GetAccountPortfolio(ctx context.Context, address model.Address) (*model.AccountPortfolio, error)
	ListAccountPortfolios(ctx context.Context, address model.Address) ([]*model.AccountPortfolio, error)
	AddAccountPortfolios(ctx context.Context, portfolios []*model.AccountPortfolio) error

	// denom
	// TODO: Remove this ugly code
	ListDenomCoinID(ctx context.Context) ([]*model.DenomCoinID, error)
	Disconnect(ctx context.Context) error
}
