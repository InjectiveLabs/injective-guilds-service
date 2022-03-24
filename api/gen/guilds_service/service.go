// Code generated by goa v3.6.2, DO NOT EDIT.
//
// GuildsService service
//
// Command:
// $ goa gen github.com/InjectiveLabs/injective-guilds-service/api/design -o ../

package guildsservice

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Service supports trading guild queries
type Service interface {
	// Get all guilds
	GetAllGuilds(context.Context) (res *GetAllGuildsResult, err error)
	// Get a single guild base on ID
	GetSingleGuild(context.Context, *GetSingleGuildPayload) (res *GetSingleGuildResult, err error)
	// Get all members a given guild (include default member)
	GetGuildMembers(context.Context, *GetGuildMembersPayload) (res *GetGuildMembersResult, err error)
	// Get master address of given guild
	GetGuildMasterAddress(context.Context, *GetGuildMasterAddressPayload) (res *GetGuildMasterAddressResult, err error)
	// Get default guild member
	GetGuildDefaultMember(context.Context, *GetGuildDefaultMemberPayload) (res *GetGuildDefaultMemberResult, err error)
	// Enter the guild: Should supply public_key, message, signature in base64
	EnterGuild(context.Context, *EnterGuildPayload) (res *EnterGuildResult, err error)
	// Enter the guild: Should supply public_key, message, signature in base64
	LeaveGuild(context.Context, *LeaveGuildPayload) (res *LeaveGuildResult, err error)
	// Get the guild markets
	GetGuildMarkets(context.Context, *GetGuildMarketsPayload) (res *GetGuildMarketsResult, err error)
	// Get current account portfolio
	GetAccountPortfolio(context.Context, *GetAccountPortfolioPayload) (res *GetAccountPortfolioResult, err error)
	// Get current account portfolios snapshots all the time
	GetAccountPortfolios(context.Context, *GetAccountPortfoliosPayload) (res *GetAccountPortfoliosResult, err error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "GuildsService"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [10]string{"GetAllGuilds", "GetSingleGuild", "GetGuildMembers", "GetGuildMasterAddress", "GetGuildDefaultMember", "EnterGuild", "LeaveGuild", "GetGuildMarkets", "GetAccountPortfolio", "GetAccountPortfolios"}

type Balance struct {
	Denom            string
	TotalBalance     string
	AvailableBalance string
	UnrealizedPnl    string
	MarginHold       string
	PriceUsd         float64
}

// EnterGuildPayload is the payload type of the GuildsService service
// EnterGuild method.
type EnterGuildPayload struct {
	GuildID   string
	PublicKey string
	// Supply base64 json encoded string cointaining {"action": "enter-guild",
	// "expired_at": unixTimestamp }
	Message   string
	Signature string
}

// EnterGuildResult is the result type of the GuildsService service EnterGuild
// method.
type EnterGuildResult struct {
	JoinStatus *string
}

// GetAccountPortfolioPayload is the payload type of the GuildsService service
// GetAccountPortfolio method.
type GetAccountPortfolioPayload struct {
	GuildID          string
	InjectiveAddress string
}

// GetAccountPortfolioResult is the result type of the GuildsService service
// GetAccountPortfolio method.
type GetAccountPortfolioResult struct {
	Data *SingleAccountPortfolio
}

// GetAccountPortfoliosPayload is the payload type of the GuildsService service
// GetAccountPortfolios method.
type GetAccountPortfoliosPayload struct {
	GuildID          string
	InjectiveAddress string
}

// GetAccountPortfoliosResult is the result type of the GuildsService service
// GetAccountPortfolios method.
type GetAccountPortfoliosResult struct {
	Portfolios []*SingleAccountPortfolio
}

// GetAllGuildsResult is the result type of the GuildsService service
// GetAllGuilds method.
type GetAllGuildsResult struct {
	// Existing guilds
	Guilds []*Guild
}

// GetGuildDefaultMemberPayload is the payload type of the GuildsService
// service GetGuildDefaultMember method.
type GetGuildDefaultMemberPayload struct {
	GuildID string
}

// GetGuildDefaultMemberResult is the result type of the GuildsService service
// GetGuildDefaultMember method.
type GetGuildDefaultMemberResult struct {
	DefaultMember *GuildMember
}

// GetGuildMarketsPayload is the payload type of the GuildsService service
// GetGuildMarkets method.
type GetGuildMarketsPayload struct {
	GuildID string
}

// GetGuildMarketsResult is the result type of the GuildsService service
// GetGuildMarkets method.
type GetGuildMarketsResult struct {
	Markets []*Market
}

// GetGuildMasterAddressPayload is the payload type of the GuildsService
// service GetGuildMasterAddress method.
type GetGuildMasterAddressPayload struct {
	GuildID string
}

// GetGuildMasterAddressResult is the result type of the GuildsService service
// GetGuildMasterAddress method.
type GetGuildMasterAddressResult struct {
	MasterAddress *string
}

// GetGuildMembersPayload is the payload type of the GuildsService service
// GetGuildMembers method.
type GetGuildMembersPayload struct {
	GuildID string
}

// GetGuildMembersResult is the result type of the GuildsService service
// GetGuildMembers method.
type GetGuildMembersResult struct {
	// Member of given guild
	Members []*GuildMember
}

// GetSingleGuildPayload is the payload type of the GuildsService service
// GetSingleGuild method.
type GetSingleGuildPayload struct {
	GuildID string
}

// GetSingleGuildResult is the result type of the GuildsService service
// GetSingleGuild method.
type GetSingleGuildResult struct {
	// Existing guilds
	Guild *Guild
}

// Guild info
type Guild struct {
	ID                         string
	Name                       string
	Description                string
	MasterAddress              string
	SpotBaseRequirement        string
	SpotQuoteRequirement       string
	DerivativeQuoteRequirement string
	StakingRequirement         string
	Capacity                   int
	MemberCount                int
}

// Guild member metadata
type GuildMember struct {
	InjectiveAddress     string
	IsDefaultGuildMember bool
	Since                int64
}

// LeaveGuildPayload is the payload type of the GuildsService service
// LeaveGuild method.
type LeaveGuildPayload struct {
	GuildID   string
	PublicKey string
	// Supply base64 json encoded string cointaining {"action": "leave-guild",
	// "expired_at": unixTimestamp}
	Message   string
	Signature string
}

// LeaveGuildResult is the result type of the GuildsService service LeaveGuild
// method.
type LeaveGuildResult struct {
	LeaveStatus *string
}

// Market supported by guild
type Market struct {
	MarketID    string
	IsPerpetual bool
}

// Single account portfio snapshot
type SingleAccountPortfolio struct {
	InjectiveAddress string
	Balances         []*Balance
	UpdatedAt        int64
}

// MakeNotFound builds a goa.ServiceError from an error.
func MakeNotFound(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "not_found",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}

// MakeInvalidArg builds a goa.ServiceError from an error.
func MakeInvalidArg(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "invalid_arg",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}

// MakeInternal builds a goa.ServiceError from an error.
func MakeInternal(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "internal",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}
