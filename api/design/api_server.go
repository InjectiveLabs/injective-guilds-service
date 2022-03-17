//go:generate goa gen github.com/InjectiveLabs/injective-guilds-service/api/design -o ../
package design

import (
	. "goa.design/goa/v3/dsl"
	_ "goa.design/plugins/v3/docs"
)

// API describes the global properties of the API server.
var _ = API("Injective Guilds Service", func() {
	Title("Injective Guilds Service")
	Description("HTTP server for the Trading guilds query")
	Server("Guilds", func() {
		Host("0.0.0.0", func() { URI("http://0.0.0.0:9930") })
	})
})
