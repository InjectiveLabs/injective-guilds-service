// Code generated by goa v3.6.2, DO NOT EDIT.
//
// GuildsService HTTP client CLI support package
//
// Command:
// $ goa gen github.com/InjectiveLabs/injective-guilds-service/api/design -o ../

package client

import (
	"encoding/json"
	"fmt"

	guildsservice "github.com/InjectiveLabs/injective-guilds-service/api/gen/guilds_service"
)

// BuildGetSingleGuildPayload builds the payload for the GuildsService
// GetSingleGuild endpoint from CLI flags.
func BuildGetSingleGuildPayload(guildsServiceGetSingleGuildGuildID string) (*guildsservice.GetSingleGuildPayload, error) {
	var guildID string
	{
		guildID = guildsServiceGetSingleGuildGuildID
	}
	v := &guildsservice.GetSingleGuildPayload{}
	v.GuildID = guildID

	return v, nil
}

// BuildGetGuildMembersPayload builds the payload for the GuildsService
// GetGuildMembers endpoint from CLI flags.
func BuildGetGuildMembersPayload(guildsServiceGetGuildMembersGuildID string) (*guildsservice.GetGuildMembersPayload, error) {
	var guildID string
	{
		guildID = guildsServiceGetGuildMembersGuildID
	}
	v := &guildsservice.GetGuildMembersPayload{}
	v.GuildID = guildID

	return v, nil
}

// BuildGetGuildMasterAddressPayload builds the payload for the GuildsService
// GetGuildMasterAddress endpoint from CLI flags.
func BuildGetGuildMasterAddressPayload(guildsServiceGetGuildMasterAddressGuildID string) (*guildsservice.GetGuildMasterAddressPayload, error) {
	var guildID string
	{
		guildID = guildsServiceGetGuildMasterAddressGuildID
	}
	v := &guildsservice.GetGuildMasterAddressPayload{}
	v.GuildID = guildID

	return v, nil
}

// BuildGetGuildDefaultMemberPayload builds the payload for the GuildsService
// GetGuildDefaultMember endpoint from CLI flags.
func BuildGetGuildDefaultMemberPayload(guildsServiceGetGuildDefaultMemberGuildID string) (*guildsservice.GetGuildDefaultMemberPayload, error) {
	var guildID string
	{
		guildID = guildsServiceGetGuildDefaultMemberGuildID
	}
	v := &guildsservice.GetGuildDefaultMemberPayload{}
	v.GuildID = guildID

	return v, nil
}

// BuildEnterGuildPayload builds the payload for the GuildsService EnterGuild
// endpoint from CLI flags.
func BuildEnterGuildPayload(guildsServiceEnterGuildBody string, guildsServiceEnterGuildGuildID string) (*guildsservice.EnterGuildPayload, error) {
	var err error
	var body EnterGuildRequestBody
	{
		err = json.Unmarshal([]byte(guildsServiceEnterGuildBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"message\": \"Nihil quia quo eligendi omnis.\",\n      \"public_key\": \"Dolore atque aperiam officiis necessitatibus dicta est.\",\n      \"signature\": \"Adipisci sed libero a nam consectetur.\"\n   }'")
		}
	}
	var guildID string
	{
		guildID = guildsServiceEnterGuildGuildID
	}
	v := &guildsservice.EnterGuildPayload{
		PublicKey: body.PublicKey,
		Message:   body.Message,
		Signature: body.Signature,
	}
	v.GuildID = guildID

	return v, nil
}

// BuildLeaveGuildPayload builds the payload for the GuildsService LeaveGuild
// endpoint from CLI flags.
func BuildLeaveGuildPayload(guildsServiceLeaveGuildBody string, guildsServiceLeaveGuildGuildID string) (*guildsservice.LeaveGuildPayload, error) {
	var err error
	var body LeaveGuildRequestBody
	{
		err = json.Unmarshal([]byte(guildsServiceLeaveGuildBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"message\": \"At quasi necessitatibus maxime enim.\",\n      \"public_key\": \"Et placeat id.\",\n      \"signature\": \"Id alias at magnam ut.\"\n   }'")
		}
	}
	var guildID string
	{
		guildID = guildsServiceLeaveGuildGuildID
	}
	v := &guildsservice.LeaveGuildPayload{
		PublicKey: body.PublicKey,
		Message:   body.Message,
		Signature: body.Signature,
	}
	v.GuildID = guildID

	return v, nil
}

// BuildGetGuildMarketsPayload builds the payload for the GuildsService
// GetGuildMarkets endpoint from CLI flags.
func BuildGetGuildMarketsPayload(guildsServiceGetGuildMarketsGuildID string) (*guildsservice.GetGuildMarketsPayload, error) {
	var guildID string
	{
		guildID = guildsServiceGetGuildMarketsGuildID
	}
	v := &guildsservice.GetGuildMarketsPayload{}
	v.GuildID = guildID

	return v, nil
}

// BuildGetAccountPortfolioPayload builds the payload for the GuildsService
// GetAccountPortfolio endpoint from CLI flags.
func BuildGetAccountPortfolioPayload(guildsServiceGetAccountPortfolioGuildID string, guildsServiceGetAccountPortfolioInjectiveAddress string) (*guildsservice.GetAccountPortfolioPayload, error) {
	var guildID string
	{
		guildID = guildsServiceGetAccountPortfolioGuildID
	}
	var injectiveAddress string
	{
		injectiveAddress = guildsServiceGetAccountPortfolioInjectiveAddress
	}
	v := &guildsservice.GetAccountPortfolioPayload{}
	v.GuildID = guildID
	v.InjectiveAddress = injectiveAddress

	return v, nil
}

// BuildGetAccountPortfoliosPayload builds the payload for the GuildsService
// GetAccountPortfolios endpoint from CLI flags.
func BuildGetAccountPortfoliosPayload(guildsServiceGetAccountPortfoliosGuildID string, guildsServiceGetAccountPortfoliosInjectiveAddress string) (*guildsservice.GetAccountPortfoliosPayload, error) {
	var guildID string
	{
		guildID = guildsServiceGetAccountPortfoliosGuildID
	}
	var injectiveAddress string
	{
		injectiveAddress = guildsServiceGetAccountPortfoliosInjectiveAddress
	}
	v := &guildsservice.GetAccountPortfoliosPayload{}
	v.GuildID = guildID
	v.InjectiveAddress = injectiveAddress

	return v, nil
}
