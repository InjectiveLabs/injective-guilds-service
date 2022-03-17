package db

import (
	"context"

	"github.com/InjectiveLabs/injective-guilds-service/internal/db/model"
)

type DBService interface {
	ListAllGuilds() ([]*model.Guild, error)
	GetSingleGuild(guildID string) (*model.Guild, error)

	// members
	GetGuildMembers(guildID string, isDefaultMember bool) ([]*model.GuildMember, error)
	AddMember(guildID string, address model.Address) error
	RemoveMember(guildID string, address model.Address) error

	// account portfolio
	GetAccountPortfolio(guildID string, address model.Address) (*model.AccountPortfolio, error)
	ListAccountPortfolios(guildID string, address model.Address) ([]*model.AccountPortfolio, error)
	Disconnect(ctx context.Context) error
}
