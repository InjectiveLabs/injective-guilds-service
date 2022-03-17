// Code generated by goa v3.6.2, DO NOT EDIT.
//
// HTTP request path constructors for the GuildService service.
//
// Command:
// $ goa gen github.com/InjectiveLabs/injective-guilds-service/api/design -o ../

package server

import (
	"fmt"
)

// GetAllGuildsGuildServicePath returns the URL path to the GuildService service GetAllGuilds HTTP endpoint.
func GetAllGuildsGuildServicePath() string {
	return "/guilds"
}

// GetSingleGuildGuildServicePath returns the URL path to the GuildService service GetSingleGuild HTTP endpoint.
func GetSingleGuildGuildServicePath(guildID string) string {
	return fmt.Sprintf("/guilds/%v", guildID)
}

// GetGuildMembersGuildServicePath returns the URL path to the GuildService service GetGuildMembers HTTP endpoint.
func GetGuildMembersGuildServicePath(guildID string) string {
	return fmt.Sprintf("/guilds/%v/members", guildID)
}

// GetGuildMasterAddressGuildServicePath returns the URL path to the GuildService service GetGuildMasterAddress HTTP endpoint.
func GetGuildMasterAddressGuildServicePath(guildID string) string {
	return fmt.Sprintf("/guilds/%v/guild-master", guildID)
}

// GetGuildDefaultMemberGuildServicePath returns the URL path to the GuildService service GetGuildDefaultMember HTTP endpoint.
func GetGuildDefaultMemberGuildServicePath(guildID string) string {
	return fmt.Sprintf("/guilds/%v/default-guild-member", guildID)
}

// EnterGuildGuildServicePath returns the URL path to the GuildService service EnterGuild HTTP endpoint.
func EnterGuildGuildServicePath(guildID string) string {
	return fmt.Sprintf("/guilds/%v/member", guildID)
}

// LeaveGuildGuildServicePath returns the URL path to the GuildService service LeaveGuild HTTP endpoint.
func LeaveGuildGuildServicePath(guildID string) string {
	return fmt.Sprintf("/guilds/%v/member", guildID)
}

// GetGuildMarketsGuildServicePath returns the URL path to the GuildService service GetGuildMarkets HTTP endpoint.
func GetGuildMarketsGuildServicePath(guildID string) string {
	return fmt.Sprintf("/guilds/%v/markets", guildID)
}

// GetAccountPortfolioGuildServicePath returns the URL path to the GuildService service GetAccountPortfolio HTTP endpoint.
func GetAccountPortfolioGuildServicePath(guildID string) string {
	return fmt.Sprintf("/guilds/%v/portfolio", guildID)
}

// GetAccountPortfoliosGuildServicePath returns the URL path to the GuildService service GetAccountPortfolios HTTP endpoint.
func GetAccountPortfoliosGuildServicePath(guildID string) string {
	return fmt.Sprintf("/guilds/%v/portfolio", guildID)
}
