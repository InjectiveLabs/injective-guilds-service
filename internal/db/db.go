package db

import (
	"context"
	"errors"

	"github.com/InjectiveLabs/injective-guilds-service/internal/db/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrMemberExceedCap = errors.New("max guild capacity has been reached")
	ErrAlreadyMember   = errors.New("already member")
)

type DBService interface {
	ListAllGuilds(ctx context.Context) ([]*model.Guild, error)
	GetSingleGuild(ctx context.Context, guildID string) (*model.Guild, error)
	ListGuildPortfolios(ctx context.Context, filter model.GuildPortfoliosFilter) ([]*model.GuildPortfolio, error)
	// TODO: *primitive.ObjectID -> string
	AddGuild(ctx context.Context, guild *model.Guild) (*primitive.ObjectID, error)
	SetGuildCap(ctx context.Context, guildID string, cap int) error
	AddGuildPortfolios(ctx context.Context, portfolios []*model.GuildPortfolio) error
	DeleteGuild(ctx context.Context, guildID string) error

	// members
	ListGuildMembers(ctx context.Context, filter model.MemberFilter) ([]*model.GuildMember, error)
	AddMember(ctx context.Context, guildID string, address model.Address, initialPortfolio *model.AccountPortfolio, isDefaultMember bool, params string) error
	RemoveMember(ctx context.Context, guildID string, address model.Address) error
	GetAccountPortfolio(ctx context.Context, address model.Address) (*model.AccountPortfolio, error)
	ListAccountPortfolios(ctx context.Context, filter model.AccountPortfoliosFilter) ([]*model.AccountPortfolio, error)
	AddAccountPortfolios(ctx context.Context, portfolios []*model.AccountPortfolio) error

	Disconnect(ctx context.Context) error
}
