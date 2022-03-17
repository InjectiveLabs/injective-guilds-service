package mongoimpl

import (
	"context"
	"time"

	"github.com/InjectiveLabs/injective-guilds-service/internal/db"
	"github.com/InjectiveLabs/injective-guilds-service/internal/db/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	connectionTimeout       = 30 * time.Second
	GuildCollectionName     = "guilds"
	MemberCollectionName    = "members"
	PortfolioCollectionName = "portfolios"
)

type MongoService struct {
	db.DBService

	client              *mongo.Client
	guildColletion      *mongo.Collection
	memberCollection    *mongo.Collection
	portfolioCollection *mongo.Collection
}

func NewService(ctx context.Context, connectionURL, databaseName string) (db.DBService, error) {
	ctx, cancel := context.WithTimeout(ctx, connectionTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionURL))
	if err != nil {
		return nil, err
	}

	return &MongoService{
		client:              client,
		guildColletion:      client.Database(databaseName).Collection(GuildCollectionName),
		memberCollection:    client.Database(databaseName).Collection(MemberCollectionName),
		portfolioCollection: client.Database(databaseName).Collection(PortfolioCollectionName),
	}, nil
}

func makeIndex(unique bool, keys interface{}) mongo.IndexModel {
	idx := mongo.IndexModel{
		Keys:    keys,
		Options: options.Index(),
	}
	if unique {
		idx.Options = idx.Options.SetUnique(true)
	}

	return idx
}

func (s *MongoService) EnsureIndex(ctx context.Context) error {
	// use CreateMany here for future custom
	// TODO: Index for faster query
	return nil
}

func (s *MongoService) ListAllGuilds() ([]*model.Guild, error) {
	return nil, nil
}

func (s *MongoService) GetSingleGuild(guildID string) (*model.Guild, error) {
	return nil, nil
}

// members
func (s *MongoService) GetGuildMembers(guildID string, isDefaultMember bool) ([]*model.GuildMember, error) {
	return nil, nil
}

func (s *MongoService) AddMember(guildID string, address model.Address) error {
	return nil
}

func (s *MongoService) RemoveMember(guildID string, address model.Address) error {
	return nil
}

// account portfolio
func (s *MongoService) GetAccountPortfolio(guildID string, address model.Address) (*model.AccountPortfolio, error) {
	return nil, nil
}

func (s *MongoService) ListAccountPortfolios(guildID string, address model.Address) ([]*model.AccountPortfolio, error) {
	return nil, nil
}

func (s *MongoService) Disconnect(ctx context.Context) error {
	return s.client.Disconnect(ctx)
}
