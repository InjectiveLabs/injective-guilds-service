package mongoimpl

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/InjectiveLabs/injective-guilds-service/internal/db"
	"github.com/InjectiveLabs/injective-guilds-service/internal/db/model"
	"github.com/InjectiveLabs/metrics"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	connectionTimeout              = 20 * time.Second
	GuildCollectionName            = "guilds"
	MemberCollectionName           = "members"
	AccountPortfolioCollectionName = "account_portfolios"
	GuildPortfolioCollectionName   = "guild_portfolios"
	DenomCollectionName            = "denoms"
)

type MongoImpl struct {
	db.DBService

	client  *mongo.Client
	session mongo.Session

	guildCollection            *mongo.Collection
	memberCollection           *mongo.Collection
	accountPortfolioCollection *mongo.Collection
	guildPortfolioCollection   *mongo.Collection
	denomCollection            *mongo.Collection
	svcTags                    metrics.Tags
}

func NewService(ctx context.Context, connectionURL, databaseName string) (db.DBService, error) {
	ctx, cancel := context.WithTimeout(ctx, connectionTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionURL))
	if err != nil {
		return nil, fmt.Errorf("connect mongo err: %w", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	session, err := client.StartSession()
	if err != nil {
		return nil, fmt.Errorf("new session err: %w", err)
	}

	return &MongoImpl{
		client:                     client,
		session:                    session,
		guildCollection:            client.Database(databaseName).Collection(GuildCollectionName),
		memberCollection:           client.Database(databaseName).Collection(MemberCollectionName),
		accountPortfolioCollection: client.Database(databaseName).Collection(AccountPortfolioCollectionName),
		guildPortfolioCollection:   client.Database(databaseName).Collection(GuildPortfolioCollectionName),
		denomCollection:            client.Database(databaseName).Collection(DenomCollectionName),
		svcTags: metrics.Tags{
			"svc": "db_svc",
		},
	}, nil
}

func makeIndex(unique bool, keys interface{}) mongo.IndexModel {
	idx := mongo.IndexModel{
		Keys:    keys,
		Options: options.Index().SetUnique(unique),
	}
	return idx
}

func (s *MongoImpl) EnsureIndex(ctx context.Context) error {
	// use CreateMany here for future custom
	_, err := s.guildCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		makeIndex(true, bson.D{{Key: "name", Value: 1}}),
	})
	if err != nil {
		return err
	}

	_, err = s.memberCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		makeIndex(true, bson.D{{Key: "injective_address", Value: 1}}),
		makeIndex(false, bson.D{{Key: "is_default_guild_member", Value: 1}}),
		makeIndex(false, bson.D{{Key: "guild_id", Value: 1}}),
	})
	if err != nil {
		return err
	}

	_, err = s.accountPortfolioCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		makeIndex(false, bson.D{{Key: "injective_address", Value: 1}, {Key: "updated_at", Value: -1}}),
		makeIndex(false, bson.D{{Key: "guild_id", Value: 1}, {Key: "updated_at", Value: -1}}),
		makeIndex(false, bson.D{{Key: "updated_at", Value: -1}}),
	})
	if err != nil {
		return err
	}

	_, err = s.guildPortfolioCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		makeIndex(false, bson.D{{Key: "guild_id", Value: 1}, {Key: "updated_at", Value: -1}}),
		makeIndex(false, bson.D{{Key: "updated_at", Value: -1}}),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *MongoImpl) ListGuildPortfolios(
	ctx context.Context,
	filter model.GuildPortfoliosFilter,
) (result []*model.GuildPortfolio, err error) {
	doneFn := metrics.ReportFuncTiming(s.svcTags)
	defer doneFn()
	metrics.ReportFuncCall(s.svcTags)

	guildObjectID, err := primitive.ObjectIDFromHex(filter.GuildID)
	if err != nil {
		metrics.ReportFuncError(s.svcTags)
		return nil, fmt.Errorf("cannot parse guildID: %w", err)
	}

	portfolioFilter := bson.M{
		"guild_id": guildObjectID,
	}

	var updatedAtFilter = make(bson.M)
	if filter.StartTime != nil {
		updatedAtFilter["$gte"] = *filter.StartTime
	}

	if filter.EndTime != nil {
		updatedAtFilter["$lt"] = *filter.EndTime
	}

	if len(updatedAtFilter) > 0 {
		portfolioFilter["updated_at"] = updatedAtFilter
	}

	opt := &options.FindOptions{}
	opt.SetSort(bson.M{"updated_at": -1})
	if filter.Limit != nil {
		opt.SetLimit(*filter.Limit)
	}

	cur, err := s.guildPortfolioCollection.Find(ctx, portfolioFilter, opt)
	if err != nil {
		metrics.ReportFuncError(s.svcTags)
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var guildPortfolio model.GuildPortfolio
		err := cur.Decode(&guildPortfolio)
		if err != nil {
			metrics.ReportFuncError(s.svcTags)
			return nil, err
		}

		result = append(result, &guildPortfolio)
	}

	return result, nil
}

func (s *MongoImpl) AddGuild(ctx context.Context, guild *model.Guild) (*primitive.ObjectID, error) {
	doneFn := metrics.ReportFuncTiming(s.svcTags)
	defer doneFn()
	metrics.ReportFuncCall(s.svcTags)

	insertOneRes, err := s.guildCollection.InsertOne(ctx, guild)
	if err != nil {
		metrics.ReportFuncError(s.svcTags)
		return nil, err
	}

	objID := insertOneRes.InsertedID.(primitive.ObjectID)
	return &objID, nil
}

func (s *MongoImpl) DeleteGuild(ctx context.Context, guildID string) error {
	doneFn := metrics.ReportFuncTiming(s.svcTags)
	defer doneFn()
	metrics.ReportFuncCall(s.svcTags)

	guildObjectID, err := primitive.ObjectIDFromHex(guildID)
	if err != nil {
		return fmt.Errorf("cannot parse guildID: %w", err)
	}

	_, err = s.session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		filter := bson.M{
			"_id": guildObjectID,
		}
		_, err := s.guildCollection.DeleteOne(sessCtx, filter)
		if err != nil {
			return nil, err
		}

		filter = bson.M{
			"guild_id": guildObjectID,
		}

		_, err = s.memberCollection.DeleteMany(sessCtx, filter)
		if err != nil {
			return nil, err
		}

		_, err = s.accountPortfolioCollection.DeleteMany(sessCtx, filter)
		if err != nil {
			return nil, err
		}

		_, err = s.guildPortfolioCollection.DeleteMany(sessCtx, filter)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		metrics.ReportFuncError(s.svcTags)
		return err
	}

	return nil
}

func (s *MongoImpl) ListAllGuilds(ctx context.Context) (result []*model.Guild, err error) {
	doneFn := metrics.ReportFuncTiming(s.svcTags)
	defer doneFn()
	metrics.ReportFuncCall(s.svcTags)

	filter := bson.M{}
	cur, err := s.guildCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var guild model.Guild
		err := cur.Decode(&guild)
		if err != nil {
			metrics.ReportFuncError(s.svcTags)
			return nil, err
		}

		result = append(result, &guild)
	}

	return result, nil
}

func (s *MongoImpl) GetSingleGuild(ctx context.Context, guildID string) (*model.Guild, error) {
	doneFn := metrics.ReportFuncTiming(s.svcTags)
	defer doneFn()
	metrics.ReportFuncCall(s.svcTags)

	guildObjectID, err := primitive.ObjectIDFromHex(guildID)
	if err != nil {
		return nil, fmt.Errorf("cannot parse guildID: %w", err)
	}

	filter := bson.M{
		"_id": guildObjectID,
	}

	res := s.guildCollection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		metrics.ReportFuncError(s.svcTags)

		if err == mongo.ErrNoDocuments {
			return nil, db.ErrNotFound
		}
		return nil, err
	}

	var guild model.Guild
	if err := res.Decode(&guild); err != nil {
		metrics.ReportFuncError(s.svcTags)
		return nil, err
	}

	return &guild, nil
}

func (s *MongoImpl) SetGuildCap(ctx context.Context, guildID string, cap int) error {
	doneFn := metrics.ReportFuncTiming(s.svcTags)
	defer doneFn()
	metrics.ReportFuncCall(s.svcTags)

	guildObjectID, err := primitive.ObjectIDFromHex(guildID)
	if err != nil {
		return fmt.Errorf("cannot parse guildID: %w", err)
	}

	filter := bson.M{
		"_id": guildObjectID,
	}

	upd := bson.M{
		"$set": bson.M{
			"capacity": cap,
		},
	}

	updateRes, err := s.guildCollection.UpdateOne(ctx, filter, upd)
	if err != nil {
		return err
	}

	if updateRes.ModifiedCount == 0 {
		return fmt.Errorf("not found guild to set cap")
	}

	return nil
}

func (s *MongoImpl) AddGuildPortfolios(ctx context.Context, portfolios []*model.GuildPortfolio) error {
	doneFn := metrics.ReportFuncTiming(s.svcTags)
	defer doneFn()
	metrics.ReportFuncCall(s.svcTags)

	docs := make([]interface{}, len(portfolios))
	for i, p := range portfolios {
		docs[i] = p
	}

	_, err := s.guildPortfolioCollection.InsertMany(ctx, docs)
	return err
}

func (s *MongoImpl) ListGuildMembers(
	ctx context.Context,
	memberFilter model.MemberFilter,
) (result []*model.GuildMember, err error) {
	doneFn := metrics.ReportFuncTiming(s.svcTags)
	defer doneFn()
	metrics.ReportFuncCall(s.svcTags)

	filter := bson.M{}
	if memberFilter.GuildID != nil {
		guildObjectID, err := primitive.ObjectIDFromHex(*memberFilter.GuildID)
		if err != nil {
			metrics.ReportFuncError(s.svcTags)
			return nil, fmt.Errorf("cannot parse guildID: %w", err)
		}
		filter["guild_id"] = guildObjectID
	}

	if memberFilter.IsDefaultMember != nil {
		filter["is_default_guild_member"] = *memberFilter.IsDefaultMember
	}

	if memberFilter.InjectiveAddress != nil {
		filter["injective_address"] = *memberFilter.InjectiveAddress
	}

	cur, err := s.memberCollection.Find(ctx, filter)
	if err != nil {
		metrics.ReportFuncError(s.svcTags)
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var member model.GuildMember
		err := cur.Decode(&member)
		if err != nil {
			metrics.ReportFuncError(s.svcTags)
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
	isDefaultMember bool,
	params string,
) (*mongo.UpdateResult, error) {
	doneFn := metrics.ReportFuncTiming(s.svcTags)
	defer doneFn()
	metrics.ReportFuncCall(s.svcTags)

	filter := bson.M{
		"injective_address": address.String(),
	}
	upd := bson.M{
		"$set": bson.M{
			"guild_id":                guildID,
			"is_default_guild_member": isDefaultMember,
			"since":                   time.Now(),
			"params":                  params,
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

// do we want to keep data for future analyze?
func (s *MongoImpl) deletePortfolios(
	ctx context.Context,
	guildID primitive.ObjectID,
	address model.Address,
) (*mongo.DeleteResult, error) {
	filter := bson.M{
		"guild_id":          guildID,
		"injective_address": address.String(),
	}

	return s.accountPortfolioCollection.DeleteMany(ctx, filter)
}

func (s *MongoImpl) adjustMemberCount(
	ctx context.Context,
	guildID primitive.ObjectID,
	increment int,
) (*mongo.UpdateResult, error) {
	filter := bson.M{
		"_id": guildID,
	}
	upd := bson.M{
		"$inc": bson.M{
			"member_count": increment,
		},
	}
	return s.guildCollection.UpdateOne(ctx, filter, upd)
}

func (s *MongoImpl) AddMember(
	ctx context.Context,
	guildID string,
	address model.Address,
	initialPortfolio *model.AccountPortfolio, isDefaultMember bool,
	params string,
) error {
	doneFn := metrics.ReportFuncTiming(s.svcTags)
	defer doneFn()
	metrics.ReportFuncCall(s.svcTags)

	guildObjectID, err := primitive.ObjectIDFromHex(guildID)
	if err != nil {
		metrics.ReportFuncError(s.svcTags)
		return fmt.Errorf("cannot parse guildID: %w", err)
	}

	_, err = s.session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		guild, err := s.GetSingleGuild(sessCtx, guildID)
		if err != nil {
			return nil, err
		}

		if guild.MemberCount >= guild.Capacity {
			return nil, db.ErrMemberExceedCap
		}

		_, err = s.adjustMemberCount(sessCtx, guildObjectID, 1)
		if err != nil {
			return nil, err
		}

		upsertRes, err := s.upsertMember(sessCtx, guildObjectID, address, isDefaultMember, params)
		if err != nil {
			return nil, err
		}

		// duplicate member, revert transaction
		if upsertRes.UpsertedCount < 1 {
			return nil, db.ErrAlreadyMember
		}

		limit := int64(1)
		latestGuildPortfolios, err := s.ListGuildPortfolios(sessCtx, model.GuildPortfoliosFilter{
			GuildID: guildID,
			Limit:   &limit,
		})
		if err != nil {
			return nil, err
		}

		if len(latestGuildPortfolios) > 0 {
			err = s.updateGuildPortfolio(sessCtx, latestGuildPortfolios[0], initialPortfolio, true)
		} else {
			guildPortfolio := &model.GuildPortfolio{
				GuildID:      guildObjectID,
				Balances:     initialPortfolio.Balances,
				BankBalances: initialPortfolio.BankBalances,
				UpdatedAt:    initialPortfolio.UpdatedAt,
			}
			err = s.AddGuildPortfolios(sessCtx, []*model.GuildPortfolio{guildPortfolio})
		}
		if err != nil {
			return nil, err
		}

		initialPortfolio.GuildID = guildObjectID
		err = s.AddAccountPortfolios(sessCtx, []*model.AccountPortfolio{initialPortfolio})
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		metrics.ReportFuncError(s.svcTags)
		return err
	}

	return nil
}

// return a + coef * b
func sumMul(a, b primitive.Decimal128, coef decimal.Decimal) primitive.Decimal128 {
	dA, _ := decimal.NewFromString(a.String())
	dB, _ := decimal.NewFromString(b.String())
	result, _ := primitive.ParseDecimal128(dA.Add(dB.Mul(coef)).String())
	return result
}

func (s *MongoImpl) updateGuildPortfolio(
	ctx context.Context,
	guildPorfolio *model.GuildPortfolio,
	accountPorfolio *model.AccountPortfolio, isAddition bool,
) error {
	coef := decimal.NewFromInt(1)
	if !isAddition {
		coef = decimal.NewFromInt(-1)
	}

	// at most 6 denoms -> 36 loop times
	for _, guildBalance := range guildPorfolio.Balances {
		for _, accBalance := range accountPorfolio.Balances {
			if guildBalance.Denom == accBalance.Denom {
				guildBalance.AvailableBalance = sumMul(guildBalance.AvailableBalance, accBalance.AvailableBalance, coef)
				guildBalance.TotalBalance = sumMul(guildBalance.TotalBalance, accBalance.TotalBalance, coef)
				guildBalance.MarginHold = sumMul(guildBalance.MarginHold, accBalance.MarginHold, coef)
				guildBalance.UnrealizedPNL = sumMul(guildBalance.UnrealizedPNL, accBalance.UnrealizedPNL, coef)
			}
		}
	}

	for _, guildBalance := range guildPorfolio.BankBalances {
		for _, accBalance := range accountPorfolio.BankBalances {
			if guildBalance.Denom == accBalance.Denom {
				guildBalance.Balance = sumMul(guildBalance.Balance, accBalance.Balance, coef)
			}
		}
	}

	filter := bson.M{
		"updated_at": guildPorfolio.UpdatedAt,
	}

	upd := bson.M{
		"$set": guildPorfolio,
	}

	_, err := s.guildPortfolioCollection.UpdateOne(ctx, filter, upd)
	if err != nil {
		return err
	}
	// don't handle case number of updated document == 0 since there might be update with +0 for all balances

	return nil
}

func (s *MongoImpl) RemoveMember(ctx context.Context, guildID string, address model.Address) error {
	doneFn := metrics.ReportFuncTiming(s.svcTags)
	defer doneFn()
	metrics.ReportFuncCall(s.svcTags)

	guildObjectID, err := primitive.ObjectIDFromHex(guildID)
	if err != nil {
		metrics.ReportFuncError(s.svcTags)
		return fmt.Errorf("cannot parse guildID: %w", err)
	}

	_, err = s.session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		deleteRes, err := s.deleteMember(sessCtx, guildObjectID, address)
		if err != nil {
			return nil, err
		}

		// expected to have 1 account deleted
		if deleteRes.DeletedCount != 1 {
			return nil, errors.New("cannot delete: no such member")
		}

		_, err = s.adjustMemberCount(sessCtx, guildObjectID, -1)
		if err != nil {
			return nil, err
		}

		limit := int64(1)
		latestAccountPortfolios, err := s.ListAccountPortfolios(sessCtx, model.AccountPortfoliosFilter{
			InjectiveAddress: address,
			Limit:            &limit,
		})
		if err != nil {
			return nil, err
		}

		if len(latestAccountPortfolios) > 0 {
			latestGuildPortfolios, err := s.ListGuildPortfolios(sessCtx, model.GuildPortfoliosFilter{
				GuildID: guildID,
				Limit:   &limit,
			})
			if err != nil {
				return nil, err
			}

			// TODO: Update snapshot timestamp upon join guild
			if len(latestGuildPortfolios) > 0 {
				err = s.updateGuildPortfolio(sessCtx, latestGuildPortfolios[0], latestAccountPortfolios[0], false)
				if err != nil {
					return nil, err
				}
			}
		}

		_, err = s.deletePortfolios(sessCtx, guildObjectID, address)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		metrics.ReportFuncError(s.svcTags)
		return err
	}

	return nil
}

// account portfolio gets latest account portfolio
// TODO: Unify getAccountPortfolio to 1 function
func (s *MongoImpl) GetAccountPortfolio(ctx context.Context, address model.Address) (*model.AccountPortfolio, error) {
	doneFn := metrics.ReportFuncTiming(s.svcTags)
	defer doneFn()
	metrics.ReportFuncCall(s.svcTags)

	filter := bson.M{
		"injective_address": address.String(),
	}

	opts := &options.FindOneOptions{}
	opts.SetSort(bson.M{"updated_at": -1})

	singleRow := s.accountPortfolioCollection.FindOne(ctx, filter, opts)
	if err := singleRow.Err(); err != nil {
		metrics.ReportFuncError(s.svcTags)
		return nil, err
	}

	var portfolio model.AccountPortfolio
	if err := singleRow.Decode(&portfolio); err != nil {
		metrics.ReportFuncError(s.svcTags)
		return nil, err
	}

	return &portfolio, nil
}

func (s *MongoImpl) ListAccountPortfolios(
	ctx context.Context,
	filter model.AccountPortfoliosFilter,
) (result []*model.AccountPortfolio, err error) {
	doneFn := metrics.ReportFuncTiming(s.svcTags)
	defer doneFn()
	metrics.ReportFuncCall(s.svcTags)

	portfolioFilter := bson.M{
		"injective_address": filter.InjectiveAddress.String(),
	}

	var updatedAtFilter = make(bson.M)
	if filter.StartTime != nil {
		updatedAtFilter["$gte"] = *filter.StartTime
	}

	if filter.EndTime != nil {
		updatedAtFilter["$lt"] = *filter.EndTime
	}

	if len(updatedAtFilter) > 0 {
		portfolioFilter["updated_at"] = updatedAtFilter
	}

	opts := &options.FindOptions{}
	opts.SetSort(bson.M{"updated_at": -1})
	if filter.Limit != nil {
		opts.SetLimit(*filter.Limit)
	}

	cur, err := s.accountPortfolioCollection.Find(ctx, portfolioFilter, opts)
	if err != nil {
		metrics.ReportFuncError(s.svcTags)
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var portfolio model.AccountPortfolio
		err := cur.Decode(&portfolio)
		if err != nil {
			metrics.ReportFuncError(s.svcTags)
			return nil, err
		}

		result = append(result, &portfolio)
	}
	return result, nil
}

// AddAccountPortfolios add portfolio snapshots in single write call
func (s *MongoImpl) AddAccountPortfolios(
	ctx context.Context,
	portfolios []*model.AccountPortfolio,
) error {
	doneFn := metrics.ReportFuncTiming(s.svcTags)
	defer doneFn()
	metrics.ReportFuncCall(s.svcTags)

	docs := make([]interface{}, len(portfolios))
	for i, p := range portfolios {
		docs[i] = p
	}

	if _, err := s.accountPortfolioCollection.InsertMany(ctx, docs); err != nil {
		metrics.ReportFuncError(s.svcTags)
		return err
	}

	return nil
}

func (s *MongoImpl) Disconnect(ctx context.Context) error {
	return s.client.Disconnect(ctx)
}

func (s *MongoImpl) GetClient() *mongo.Client {
	return s.client
}
