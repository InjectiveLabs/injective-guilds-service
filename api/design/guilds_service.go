package design

import (
	. "goa.design/goa/v3/dsl"
	cors "goa.design/plugins/v3/cors/dsl"
	_ "goa.design/plugins/v3/docs"
)

// Repo: Injective guilds service
// This service will be named Guilds
// API params/result should be snake_case

var _ = Service("GuildsService", func() {
	Description("Service supports trading guild queries")

	cors.Origin("*", func() {
		cors.Methods("GET", "POST", "DELETE")
		cors.Headers("Content-Type")
	})

	Error("not_found", ErrorResult, "Not found")
	Error("invalid_arg", ErrorResult, "Invalid argument")
	Error("internal", ErrorResult, "Internal Server Error")

	Method("GetAllGuilds", func() {
		Description("Get all guilds")

		Result(func() {
			Field(1, "guilds", ArrayOf(Guild), func() {
				Description("Existing guilds")
			})
		})

		HTTP(func() {
			GET("/guilds")

			Response(CodeOK)
			Response("not_found", StatusNotFound)
			Response("internal", StatusInternalServerError)
		})
	})

	Method("GetSingleGuild", func() {
		// TODO: add example later
		Description("Get a single guild base on ID")

		Payload(func() {
			Field(1, "guildID", String)
			Required("guildID")
		})

		Result(func() {
			Field(1, "guild", Guild, func() {
				Description("Existing guilds")
			})
		})

		HTTP(func() {
			GET("/guilds/{guildID}")

			Response(CodeOK)
			Response("not_found", StatusNotFound)
			Response("internal", StatusInternalServerError)
		})
	})

	Method("GetGuildMembers", func() {
		Description("Get all members a given guild (include default member)")

		Payload(func() {
			// TODO: Basic validation
			Field(1, "guildID", String)
			Required("guildID")
		})

		Result(func() {
			Field(1, "members", ArrayOf(GuildMember), func() {
				Description("Member of given guild")
			})
		})

		HTTP(func() {
			GET("/guilds/{guildID}/members")

			Response(CodeOK)
			Response("not_found", StatusNotFound)
			Response("internal", StatusInternalServerError)
		})
	})

	Method("GetGuildMasterAddress", func() {
		Description("Get master address of given guild")

		Payload(func() {
			// TODO: Basic validation
			Field(1, "guildID", String)
			Required("guildID")
		})

		Result(func() {
			Field(1, "master_address", String)
		})

		HTTP(func() {
			GET("/guilds/{guildID}/guild-master")

			Response(CodeOK)
			Response("not_found", StatusNotFound)
			Response("internal", StatusInternalServerError)
		})
	})

	Method("GetGuildDefaultMember", func() {
		Description("Get default guild member")
		// TODO: Need address only ?
		Payload(func() {
			// TODO: Basic validation
			Field(1, "guildID", String)
			Required("guildID")
		})

		Result(func() {
			Field(1, "default_member", GuildMember)
		})

		HTTP(func() {
			GET("/guilds/{guildID}/default-member")

			Response(CodeOK)
			Response("not_found", 404)
			Response("internal", 500)
		})
	})

	Method("EnterGuild", func() {
		Description("Enter the guild")

		Payload(func() {
			Field(0, "guildID", String)
			Field(1, "injective_address", String)
			Field(2, "params", String)

			Required("guildID")
			Required("injective_address")
		})

		Result(func() {
			Field(1, "join_status", String)
			Field(2, "message", String)
		})

		HTTP(func() {
			POST("/guilds/{guildID}/member")

			Response(CodeOK)
			Response("not_found", StatusNotFound)
			Response("internal", StatusInternalServerError)
		})
	})

	Method("LeaveGuild", func() {
		Description("Leave the guild, guildID")

		Payload(func() {
			Field(0, "guildID", String)
			Field(1, "injective_address", String)

			Required("guildID")
			Required("injective_address")
		})

		Result(func() {
			Field(1, "leave_status", String)
		})
		HTTP(func() {
			DELETE("/guilds/{guildID}/member/{injective_address}")

			Response(CodeOK)
			Response("not_found", StatusNotFound)
			Response("internal", StatusInternalServerError)
		})
	})

	Method("GetGuildMarkets", func() {
		Description("Get the guild markets")

		Payload(func() {
			// TODO: Basic validation
			Field(1, "guildID", String)
			Required("guildID")
		})

		Result(func() {
			Field(1, "markets", ArrayOf(Market))
		})

		HTTP(func() {
			GET("/guilds/{guildID}/markets")

			Response(CodeOK)
			Response("not_found", StatusNotFound)
			Response("internal", StatusInternalServerError)
		})
	})

	Method("GetGuildPortfolios", func() {
		Description("Get the guild markets")

		Payload(func() {
			// TODO: Basic validation
			Field(1, "guildID", String)
			Field(2, "start_time", Int64)
			Field(3, "end_time", Int64)
			Required("guildID")
		})

		Result(func() {
			Field(1, "portfolios", ArrayOf(SingleGuildPortfolio))
		})

		HTTP(func() {
			GET("/guilds/{guildID}/portfolios")
			Param("start_time")
			Param("end_time")

			Response(CodeOK)
			Response("not_found", StatusNotFound)
			Response("internal", StatusInternalServerError)
		})
	})

	Method("GetAccountInfo", func() {
		Description("Get current account member status")

		Payload(func() {
			Field(1, "injective_address", String)
			Required("injective_address")
		})

		Result(func() {
			Field(1, "data", GuildMember)
		})

		HTTP(func() {
			GET("/members/{injective_address}")

			Response(CodeOK)
			Response("not_found", StatusNotFound)
			Response("internal", StatusInternalServerError)
		})
	})

	Method("GetAccountPortfolio", func() {
		Description("Get current account portfolio snapshot")

		Payload(func() {
			Field(1, "injective_address", String)
			Required("injective_address")
		})

		Result(func() {
			Field(1, "data", SingleAccountPortfolio)
		})

		HTTP(func() {
			GET("/members/{injective_address}/portfolio")

			Response(CodeOK)
			Response("not_found", StatusNotFound)
			Response("internal", StatusInternalServerError)
		})
	})

	Method("GetAccountPortfolios", func() {
		Description("Get current account portfolios snapshots all the time")

		Payload(func() {
			Field(1, "injective_address", String)
			Field(2, "start_time", Int64)
			Field(3, "end_time", Int64)
			Required("injective_address")
		})

		Result(func() {
			Field(1, "portfolios", ArrayOf(SingleAccountPortfolio))
		})

		HTTP(func() {
			GET("/members/{injective_address}/portfolios")
			Param("start_time")
			Param("end_time")

			Response(CodeOK)
			Response("not_found", StatusNotFound)
			Response("internal", StatusInternalServerError)
		})
	})

	Method("GetAccountMonthlyPortfolios", func() {
		Description("Get current account portfolios monthly snapshots, including start_time, end_time snapshots")

		Payload(func() {
			Field(1, "injective_address", String)
			Field(2, "start_time", Int64)
			Field(3, "end_time", Int64)
			Required("injective_address")
		})

		Result(func() {
			Field(1, "portfolios", ArrayOf(MonthlyAccountPortfolio))
		})

		HTTP(func() {
			GET("/members/{injective_address}/monthly-portfolios")
			Param("start_time")
			Param("end_time")

			Response(CodeOK)
			Response("not_found", StatusNotFound)
			Response("internal", StatusInternalServerError)
		})
	})
})
