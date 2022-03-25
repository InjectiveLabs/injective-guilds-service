package design

import (
	. "goa.design/goa/v3/dsl"
	_ "goa.design/plugins/v3/docs"
)

var Market = Type("Market", func() {
	Description("Market supported by guild")
	Field(1, "market_id", String)
	Field(2, "is_perpetual", Boolean)

	Required("market_id")
	Required("is_perpetual")
})

var Requirement = Type("Requirement", func() {
	Field(1, "denom", String)
	Field(2, "min_amount_usd", Float64)

	Required("denom")
	Required("min_amount_usd")
})

var Guild = Type("Guild", func() {
	Description("Guild info")

	Field(1, "id", String)
	Field(2, "name", String)
	Field(3, "description", String)

	Field(4, "master_address", String)
	Field(5, "requirements", ArrayOf(Requirement))

	Field(6, "staking_requirement", String)
	Field(7, "capacity", Int)
	Field(8, "member_count", Int)
	Field(9, "current_portfolio", SingleGuildPortfolio)

	Required("id")
	Required("name")
	Required("description")
	Required("master_address")
	Required("requirements")
	Required("staking_requirement")
	Required("capacity")
	Required("member_count")
})

var Balance = Type("Balance", func() {
	Field(1, "denom", String)
	Field(2, "total_balance", String)
	Field(3, "available_balance", String)
	Field(4, "unrealized_pnl", String)
	Field(5, "margin_hold", String)
	Field(6, "price_usd", Float64)

	Required("denom")
	Required("total_balance")
	Required("available_balance")
	Required("unrealized_pnl")
	Required("margin_hold")
	Required("price_usd")
})

var SingleAccountPortfolio = Type("SingleAccountPortfolio", func() {
	Description("Single account portfio snapshot")
	Field(1, "injective_address", String)
	Field(2, "balances", ArrayOf(Balance))
	Field(3, "updated_at", Int64)
	Required("injective_address")
	Required("balances")
	Required("updated_at")
})

var SingleGuildPortfolio = Type("SingleGuildPortfolio", func() {
	Description("Single guild portfolio snapshot")
	Field(0, "guild_id", String)
	Field(1, "balances", ArrayOf(Balance))
	Field(2, "updated_at", Int64)
	Required("balances")
	Required("updated_at")
})

var GuildMember = Type("GuildMember", func() {
	Description("Guild member metadata")
	Field(1, "injective_address", String)
	Field(2, "is_default_guild_member", Boolean)
	Field(2, "since", Int64)

	Required("injective_address")
	Required("is_default_guild_member")
	Required("since")
})
