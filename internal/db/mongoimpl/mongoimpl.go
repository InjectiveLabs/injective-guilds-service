package mongoimpl

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/InjectiveLabs/injective-guilds-service/internal/db"
	"github.com/InjectiveLabs/injective-guilds-service/internal/db/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	connectionTimeout       = 30 * time.Second
	GuildCollectionName     = "guilds"
	MemberCollectionName    = "members"
	PortfolioCollectionName = "portfolios"
)

var (
	ErrNotFound        = errors.New("dberr: not found")
	ErrMemberExceedCap = errors.New("member exceeds cap")
	ErrAlreadyMember   = errors.New("already member")
)

type MongoImpl struct {
	db.DBService

	client  *mongo.Client
	session mongo.Session

	guildCollection     *mongo.Collection
	memberCollection    *mongo.Collection
	portfolioCollection *mongo.Collection
}

func NewService(ctx context.Context, connectionURL, databaseName string) (db.DBService, error) {
	ctx, cancel := context.WithTimeout(ctx, connectionTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionURL))
	if err != nil {
		return nil, fmt.Errorf("connect mongo err: %w", err)
	}

	session, err := client.StartSession()
	if err != nil {
		return nil, fmt.Errorf("new session err: %w", err)
	}

	return &MongoImpl{
		client:              client,
		session:             session,
		guildCollection:     client.Database(databaseName).Collection(GuildCollectionName),
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

func (s *MongoImpl) EnsureIndex(ctx context.Context) error {
	// use CreateMany here for future custom
	// TODO: Index for faster query
	return nil
}

func (s *MongoImpl) ListAllGuilds(ctx context.Context) (result []*model.Guild, err error) {
	filter := bson.M{}
	cur, err := s.guildCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var guild model.Guild
		err := cur.Decode(&guild)
		if err != nil {
			return nil, err
		}

		result = append(result, &guild)
	}

	return result, nil
}

func (s *MongoImpl) GetSingleGuild(ctx context.Context, guildID string) (*model.Guild, error) {
	guildObjectID, err := primitive.ObjectIDFromHex(guildID)
	if err != nil {
		return nil, fmt.Errorf("cannot parse guildID: %w", err)
	}

	filter := bson.M{
		"_id": guildObjectID,
	}

	res := s.guildCollection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		return nil, err
	}

	var guild model.Guild
	if err := res.Decode(&guild); err != nil {
		return nil, err
	}

	return &guild, nil
}

func (s *MongoImpl) GetGuildMembers(
	ctx context.Context,
	guildID string,
	isDefaultMember bool,
) (result []*model.GuildMember, err error) {
	guildObjectID, err := primitive.ObjectIDFromHex(guildID)
	if err != nil {
		return nil, fmt.Errorf("cannot parse guildID: %w", err)
	}

	filter := bson.M{
		"guild_id":                guildObjectID,
		"is_default_guild_member": isDefaultMember,
	}

	cur, err := s.memberCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var member model.GuildMember
		err := cur.Decode(&member)
		if err != nil {
			return nil, err
		}

		result = append(result, &member)
	}

	return result, nil
}

func (s *MongoImpl) upsertMember(
	ctx context.Context,
	guildID primitive.ObjectID,
	address model.Address,
) (*mongo.UpdateResult, error) {
	filter := bson.M{
		"injective_address": address.String(),
	}
	upd := bson.M{
		"$set": bson.M{
			"guild_id": guildID,
		},
	}
	updOpt := &options.UpdateOptions{}
	updOpt.SetUpsert(true)

	return s.memberCollection.UpdateOne(ctx, filter, upd, updOpt)
}

func (s *MongoImpl) deleteMember(
	ctx context.Context,
	guildID primitive.ObjectID,
	address model.Address,
) (*mongo.DeleteResult, error) {
	filter := bson.M{
		"guild_id":          guildID,
		"injective_address": address.String(),
	}

	return s.memberCollection.DeleteOne(ctx, filter)
}

func (s *MongoImpl) adjustMemberCount(
	ctx context.Context,
	guildID primitive.ObjectID,
	increment int,
) (*mongo.UpdateResult, error) {
	filter := bson.M{
		"guild_id": guildID,
	}
	upd := bson.M{
		"$inc": bson.M{
			"member_count": increment,
		},
	}
	return s.guildCollection.UpdateOne(ctx, filter, upd)
}

func (s *MongoImpl) AddMember(ctx context.Context, guildID string, address model.Address) error {
	guildObjectID, err := primitive.ObjectIDFromHex(guildID)
	if err != nil {
		return fmt.Errorf("cannot parse guildID: %w", err)
	}

	_, err = s.session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		guild, err := s.GetSingleGuild(sessCtx, guildID)
		if err != nil {
			return nil, err
		}

		if guild.MemberCount >= guild.Capacity {
			return nil, ErrMemberExceedCap
		}

		_, err = s.adjustMemberCount(sessCtx, guildObjectID, 1)
		if err != nil {
			return nil, err
		}

		upsertRes, err := s.upsertMember(sessCtx, guildObjectID, address)
		if err != nil {
			return nil, err
		}

		// duplicate member
		if upsertRes.UpsertedCount < 1 {
			return nil, ErrAlreadyMember
		}

		return nil, nil
	})

	return err
}

func (s *MongoImpl) RemoveMember(ctx context.Context, guildID string, address model.Address) error {
	guildObjectID, err := primitive.ObjectIDFromHex(guildID)
	if err != nil {
		return fmt.Errorf("cannot parse guildID: %w", err)
	}

	_, err = s.session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		deleteRes, err := s.deleteMember(ctx, guildObjectID, address)
		if err != nil {
			return nil, err
		}

		// expected to have 1 account deleted
		if deleteRes.DeletedCount != 1 {
			return nil, errors.New("cannot delete")
		}

		_, err = s.adjustMemberCount(sessCtx, guildObjectID, -1)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})

	return err
}

// account portfolio gets latest account portfolio
func (s *MongoImpl) GetAccountPortfolio(ctx context.Context, guildID string, address model.Address) ([]*model.AccountPortfolio, error) {
	// fetch guild supported markets
	guild, err := s.GetSingleGuild(ctx, guildID)
	if err != nil {
		return nil, fmt.Errorf("cannot find guild: %w", err)
	}

	denoms := make([]string, 0)
	for _, m := range guild.Markets {
		denoms = append(denoms, m.Denom)
	}

	filter := bson.M{
		"injective_address": address.String(),
		"denom":             bson.M{"$in": denoms},
	}

	opts := &options.FindOptions{}
	opts.SetSort(bson.M{"updated_at": -1})

	cur, err := s.portfolioCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	var (
		recentTime time.Time
		result     []*model.AccountPortfolio
	)

	for cur.Next(ctx) {
		var portfolio model.AccountPortfolio
		err := cur.Decode(&portfolio)
		if err != nil {
			return nil, err
		}

		if recentTime.IsZero() {
			recentTime = portfolio.UpdatedAt
		} else if recentTime != portfolio.UpdatedAt {
			break
		}

		result = append(result, &portfolio)
	}

	return result, nil
}

func (s *MongoImpl) ListAccountPortfolios(
	ctx context.Context,
	guildID string,
	address model.Address,
) (result []*model.AccountPortfolio, err error) {
	// fetch guild supported markets
	guild, err := s.GetSingleGuild(ctx, guildID)
	if err != nil {
		return nil, fmt.Errorf("cannot find guild: %w", err)
	}

	denoms := make([]string, 0)
	for _, m := range guild.Markets {
		denoms = append(denoms, m.Denom)
	}

	filter := bson.M{
		"injective_address": address.String(),
		"denom":             bson.M{"$in": denoms},
	}
	opts := &options.FindOptions{}
	opts.SetSort(bson.M{"updated_at": -1})

	cur, err := s.portfolioCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var portfolio model.AccountPortfolio
		err := cur.Decode(&portfolio)
		if err != nil {
			return nil, err
		}

		result = append(result, &portfolio)
	}
	return result, nil
}

func (s *MongoImpl) Disconnect(ctx context.Context) error {
	return s.client.Disconnect(ctx)
}