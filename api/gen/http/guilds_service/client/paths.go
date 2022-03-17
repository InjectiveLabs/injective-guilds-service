// Code generated by goa v3.6.2, DO NOT EDIT.
//
// HTTP request path constructors for the GuildsService service.
//
// Command:
// $ goa gen github.com/InjectiveLabs/injective-guilds-service/api/design -o ../

package client

import (
	"fmt"
)

// GetAllGuildsGuildsServicePath returns the URL path to the GuildsService service GetAllGuilds HTTP endpoint.
func GetAllGuildsGuildsServicePath() string {
	return "/guilds"
}

// GetSingleGuildGuildsServicePath returns the URL path to the GuildsService service GetSingleGuild HTTP endpoint.
func GetSingleGuildGuildsServicePath(guildID string) string {
	return fmt.Sprintf("/guilds/%v", guildID)
}

// GetGuildMembersGuildsServicePath returns the URL path to the GuildsService service GetGuildMembers HTTP endpoint.
func GetGuildMembersGuildsServicePath(guildID string) string {
	return fmt.Sprintf("/guilds/%v/members", guildID)
}

// GetGuildMasterAddressGuildsServicePath returns the URL path to the GuildsService service GetGuildMasterAddress HTTP endpoint.
func GetGuildMasterAddressGuildsServicePath(guildID string) string {
	return fmt.Sprintf("/guilds/%v/guild-master", guildID)
}

// GetGuildDefaultMemberGuildsServicePath returns the URL path to the GuildsService service GetGuildDefaultMember HTTP endpoint.
func GetGuildDefaultMemberGuildsServicePath(guildID string) string {
	return fmt.Sprintf("/guilds/%v/default-guild-member", guildID)
}

// EnterGuildGuildsServicePath returns the URL path to the GuildsService service EnterGuild HTTP endpoint.
func EnterGuildGuildsServicePath(guildID string) string {
	return fmt.Sprintf("/guilds/%v/member", guildID)
}

// LeaveGuildGuildsServicePath returns the URL path to the GuildsService service LeaveGuild HTTP endpoint.
func LeaveGuildGuildsServicePath(guildID string) string {
	return fmt.Sprintf("/guilds/%v/member", guildID)
}

// GetGuildMarketsGuildsServicePath returns the URL path to the GuildsService service GetGuildMarkets HTTP endpoint.
func GetGuildMarketsGuildsServicePath(guildID string) string {
	return fmt.Sprintf("/guilds/%v/markets", guildID)
}

// GetAccountPortfolioGuildsServicePath returns the URL path to the GuildsService service GetAccountPortfolio HTTP endpoint.
func GetAccountPortfolioGuildsServicePath(guildID string) string {
	return fmt.Sprintf("/guilds/%v/portfolio", guildID)
}

// GetAccountPortfoliosGuildsServicePath returns the URL path to the GuildsService service GetAccountPortfolios HTTP endpoint.
func GetAccountPortfoliosGuildsServicePath(guildID string) string {
	return fmt.Sprintf("/guilds/%v/portfolios", guildID)
}
