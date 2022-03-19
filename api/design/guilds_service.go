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
			Response("not_found", CodeNotFound)
			Response("internal", CodeInternal)
		})
	})

	Method("GetSingleGuild", func() {
		// TODO: add example later
		Description("Get a single guild")

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
			Response("not_found", CodeNotFound)
			Response("internal", CodeInternal)
		})
	})

	Method("GetGuildMembers", func() {
		Description("Get members")

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
			Response("not_found", CodeNotFound)
			Response("internal", CodeInternal)
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
			Response("not_found", CodeNotFound)
			Response("internal", CodeInternal)
		})
	})

	Method("GetGuildDefaultMember", func() {
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
			GET("/guilds/{guildID}/default-guild-member")

			Response(CodeOK)
			Response("not_found", CodeNotFound)
			Response("internal", CodeInternal)
		})
	})

	Method("EnterGuild", func() {
		Payload(func() {
			Field(0, "guildID", String)
			Field(1, "public_key", String)
			Field(2, "message", String)
			Field(3, "signature", String)

			Required("guildID")
			Required("public_key")
			Required("message")
			Required("signature")
		})

		Result(func() {
			Field(1, "join_status", String)
		})

		HTTP(func() {
			POST("/guilds/{guildID}/member")

			Response(CodeOK)
			Response("not_found", CodeNotFound)
			Response("internal", CodeInternal)
		})
	})

	Method("LeaveGuild", func() {
		Payload(func() {
			Field(0, "guildID", String)
			Field(1, "public_key", String)
			Field(2, "message", String)
			Field(3, "signature", String)

			Required("guildID")
			Required("public_key")
			Required("message")
			Required("signature")
		})

		Result(func() {
			Field(1, "leave_status", String)
		})
		HTTP(func() {
			DELETE("/guilds/{guildID}/member")

			Response(CodeOK)
			Response("not_found", CodeNotFound)
			Response("internal", CodeInternal)
		})
	})

	// Markets
	Method("GetGuildMarkets", func() {
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
			Response("not_found", CodeNotFound)
			Response("internal", CodeInternal)
		})
	})

	// Account's Porfolio(s)
	Method("GetAccountPortfolio", func() {
		Payload(func() {
			Field(1, "guildID", String)
			Field(2, "injective_address", String)

			Required("guildID")
			Required("injective_address")
		})

		Result(func() {
			Field(1, "data", SingleAccountPortfolio)
		})

		HTTP(func() {
			GET("/guilds/{guildID}/portfolio")
			Param("injective_address")

			Response(CodeOK)
			Response("not_found", CodeNotFound)
			Response("internal", CodeInternal)
		})
	})

	Method("GetAccountPortfolios", func() {
		Payload(func() {
			Field(1, "guildID", String)
			Field(2, "injective_address", String)

			Required("guildID")
			Required("injective_address")
		})

		Result(func() {
			Field(1, "portfolios", ArrayOf(SingleAccountPortfolio))
		})

		HTTP(func() {
			GET("/guilds/{guildID}/portfolios")
			Param("injective_address")

			Response(CodeOK)
			Response("not_found", CodeNotFound)
			Response("internal", CodeInternal)
		})
	})
})
