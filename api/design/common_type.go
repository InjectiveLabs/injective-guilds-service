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

var Guild = Type("Guild", func() {
	Description("Guild info")

	Field(1, "id", String)
	Field(2, "name", String)
	Field(3, "description", String)

	Field(4, "master_address", String)
	Field(5, "spot_base_requirement", String)
	Field(6, "spot_quote_requirement", String)
	Field(7, "derivative_quote_requirement", String)

	Field(8, "staking_requirement", String)
	Field(9, "capacity", Int)
	Field(10, "member_count", Int)

	Required("id")
	Required("name")
	Required("description")
	Required("master_address")
	Required("spot_base_requirement")
	Required("spot_quote_requirement")
	Required("derivative_quote_requirement")

	Required("staking_requirement")
	Required("capacity")
	Required("member_count")
})

var EmbededAccountPortfolio = Type("EmbededAccountPortfolio", func() {
	Field(1, "denom", String)
	Field(2, "total_balance", String)
	Field(3, "available_balance", String)
	Field(4, "unrealized_pnl", String)
	Field(5, "margin_hold", String)

	Field(6, "updated_at", String, func() {
		Format(FormatDateTime)
	})

	Required("denom")
	Required("total_balance")
	Required("available_balance")
	Required("unrealized_pnl")
	Required("margin_hold")
	Required("updated_at")
})

var SingleAccountPortfolio = Type("SingleAccountPortfolio", func() {
	Description("Account portfio snapshot")
	Field(1, "injective_address", String)
	Field(1, "denom", String)
	Field(2, "total_balance", String)
	Field(3, "available_balance", String)
	Field(4, "unrealized_pnl", String)
	Field(5, "margin_hold", String)

	Field(6, "updated_at", String, func() {
		Format(FormatDateTime)
	})

	Required("injective_address")
	Required("denom")
	Required("total_balance")
	Required("available_balance")
	Required("unrealized_pnl")
	Required("margin_hold")
	Required("updated_at")
})

var AccountPortfolios = Type("AccountPorfolios", func() {
	Description("Account portfio snapshots")
	Field(1, "injective_address", String)
	Field(2, "portfolios", ArrayOf(EmbededAccountPortfolio))

	Required("injective_address")
	Required("portfolios")
})

var GuildMember = Type("GuildMember", func() {
	Description("Guild member metadata")
	Field(1, "injective_address", String)
	Field(2, "is_default_guild_member", Boolean)
})
