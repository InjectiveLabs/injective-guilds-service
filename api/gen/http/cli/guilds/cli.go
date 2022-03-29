// Code generated by goa v3.6.2, DO NOT EDIT.
//
// Guilds HTTP client CLI support package
//
// Command:
// $ goa gen github.com/InjectiveLabs/injective-guilds-service/api/design -o ../

package cli

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	guildsservicec "github.com/InjectiveLabs/injective-guilds-service/api/gen/http/guilds_service/client"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// UsageCommands returns the set of commands and sub-commands using the format
//
//    command (subcommand1|subcommand2|...)
//
func UsageCommands() string {
	return `guilds-service (get-all-guilds|get-single-guild|get-guild-members|get-guild-master-address|get-guild-default-member|enter-guild|leave-guild|get-guild-markets|get-guild-portfolios|get-account-info|get-account-portfolio|get-account-portfolios)
`
}

// UsageExamples produces an example of a valid invocation of the CLI tool.
func UsageExamples() string {
	return os.Args[0] + ` guilds-service get-all-guilds` + "\n" +
		""
}

// ParseEndpoint returns the endpoint and payload as specified on the command
// line.
func ParseEndpoint(
	scheme, host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restore bool,
) (goa.Endpoint, interface{}, error) {
	var (
		guildsServiceFlags = flag.NewFlagSet("guilds-service", flag.ContinueOnError)

		guildsServiceGetAllGuildsFlags = flag.NewFlagSet("get-all-guilds", flag.ExitOnError)

		guildsServiceGetSingleGuildFlags       = flag.NewFlagSet("get-single-guild", flag.ExitOnError)
		guildsServiceGetSingleGuildGuildIDFlag = guildsServiceGetSingleGuildFlags.String("guild-id", "REQUIRED", "")

		guildsServiceGetGuildMembersFlags       = flag.NewFlagSet("get-guild-members", flag.ExitOnError)
		guildsServiceGetGuildMembersGuildIDFlag = guildsServiceGetGuildMembersFlags.String("guild-id", "REQUIRED", "")

		guildsServiceGetGuildMasterAddressFlags       = flag.NewFlagSet("get-guild-master-address", flag.ExitOnError)
		guildsServiceGetGuildMasterAddressGuildIDFlag = guildsServiceGetGuildMasterAddressFlags.String("guild-id", "REQUIRED", "")

		guildsServiceGetGuildDefaultMemberFlags       = flag.NewFlagSet("get-guild-default-member", flag.ExitOnError)
		guildsServiceGetGuildDefaultMemberGuildIDFlag = guildsServiceGetGuildDefaultMemberFlags.String("guild-id", "REQUIRED", "")

		guildsServiceEnterGuildFlags       = flag.NewFlagSet("enter-guild", flag.ExitOnError)
		guildsServiceEnterGuildBodyFlag    = guildsServiceEnterGuildFlags.String("body", "REQUIRED", "")
		guildsServiceEnterGuildGuildIDFlag = guildsServiceEnterGuildFlags.String("guild-id", "REQUIRED", "")

		guildsServiceLeaveGuildFlags                = flag.NewFlagSet("leave-guild", flag.ExitOnError)
		guildsServiceLeaveGuildGuildIDFlag          = guildsServiceLeaveGuildFlags.String("guild-id", "REQUIRED", "")
		guildsServiceLeaveGuildInjectiveAddressFlag = guildsServiceLeaveGuildFlags.String("injective-address", "REQUIRED", "")

		guildsServiceGetGuildMarketsFlags       = flag.NewFlagSet("get-guild-markets", flag.ExitOnError)
		guildsServiceGetGuildMarketsGuildIDFlag = guildsServiceGetGuildMarketsFlags.String("guild-id", "REQUIRED", "")

		guildsServiceGetGuildPortfoliosFlags         = flag.NewFlagSet("get-guild-portfolios", flag.ExitOnError)
		guildsServiceGetGuildPortfoliosGuildIDFlag   = guildsServiceGetGuildPortfoliosFlags.String("guild-id", "REQUIRED", "")
		guildsServiceGetGuildPortfoliosStartTimeFlag = guildsServiceGetGuildPortfoliosFlags.String("start-time", "", "")
		guildsServiceGetGuildPortfoliosEndTimeFlag   = guildsServiceGetGuildPortfoliosFlags.String("end-time", "", "")

		guildsServiceGetAccountInfoFlags                = flag.NewFlagSet("get-account-info", flag.ExitOnError)
		guildsServiceGetAccountInfoInjectiveAddressFlag = guildsServiceGetAccountInfoFlags.String("injective-address", "REQUIRED", "")

		guildsServiceGetAccountPortfolioFlags                = flag.NewFlagSet("get-account-portfolio", flag.ExitOnError)
		guildsServiceGetAccountPortfolioInjectiveAddressFlag = guildsServiceGetAccountPortfolioFlags.String("injective-address", "REQUIRED", "")

		guildsServiceGetAccountPortfoliosFlags                = flag.NewFlagSet("get-account-portfolios", flag.ExitOnError)
		guildsServiceGetAccountPortfoliosInjectiveAddressFlag = guildsServiceGetAccountPortfoliosFlags.String("injective-address", "REQUIRED", "")
		guildsServiceGetAccountPortfoliosStartTimeFlag        = guildsServiceGetAccountPortfoliosFlags.String("start-time", "", "")
		guildsServiceGetAccountPortfoliosEndTimeFlag          = guildsServiceGetAccountPortfoliosFlags.String("end-time", "", "")
	)
	guildsServiceFlags.Usage = guildsServiceUsage
	guildsServiceGetAllGuildsFlags.Usage = guildsServiceGetAllGuildsUsage
	guildsServiceGetSingleGuildFlags.Usage = guildsServiceGetSingleGuildUsage
	guildsServiceGetGuildMembersFlags.Usage = guildsServiceGetGuildMembersUsage
	guildsServiceGetGuildMasterAddressFlags.Usage = guildsServiceGetGuildMasterAddressUsage
	guildsServiceGetGuildDefaultMemberFlags.Usage = guildsServiceGetGuildDefaultMemberUsage
	guildsServiceEnterGuildFlags.Usage = guildsServiceEnterGuildUsage
	guildsServiceLeaveGuildFlags.Usage = guildsServiceLeaveGuildUsage
	guildsServiceGetGuildMarketsFlags.Usage = guildsServiceGetGuildMarketsUsage
	guildsServiceGetGuildPortfoliosFlags.Usage = guildsServiceGetGuildPortfoliosUsage
	guildsServiceGetAccountInfoFlags.Usage = guildsServiceGetAccountInfoUsage
	guildsServiceGetAccountPortfolioFlags.Usage = guildsServiceGetAccountPortfolioUsage
	guildsServiceGetAccountPortfoliosFlags.Usage = guildsServiceGetAccountPortfoliosUsage

	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		return nil, nil, err
	}

	if flag.NArg() < 2 { // two non flag args are required: SERVICE and ENDPOINT (aka COMMAND)
		return nil, nil, fmt.Errorf("not enough arguments")
	}

	var (
		svcn string
		svcf *flag.FlagSet
	)
	{
		svcn = flag.Arg(0)
		switch svcn {
		case "guilds-service":
			svcf = guildsServiceFlags
		default:
			return nil, nil, fmt.Errorf("unknown service %q", svcn)
		}
	}
	if err := svcf.Parse(flag.Args()[1:]); err != nil {
		return nil, nil, err
	}

	var (
		epn string
		epf *flag.FlagSet
	)
	{
		epn = svcf.Arg(0)
		switch svcn {
		case "guilds-service":
			switch epn {
			case "get-all-guilds":
				epf = guildsServiceGetAllGuildsFlags

			case "get-single-guild":
				epf = guildsServiceGetSingleGuildFlags

			case "get-guild-members":
				epf = guildsServiceGetGuildMembersFlags

			case "get-guild-master-address":
				epf = guildsServiceGetGuildMasterAddressFlags

			case "get-guild-default-member":
				epf = guildsServiceGetGuildDefaultMemberFlags

			case "enter-guild":
				epf = guildsServiceEnterGuildFlags

			case "leave-guild":
				epf = guildsServiceLeaveGuildFlags

			case "get-guild-markets":
				epf = guildsServiceGetGuildMarketsFlags

			case "get-guild-portfolios":
				epf = guildsServiceGetGuildPortfoliosFlags

			case "get-account-info":
				epf = guildsServiceGetAccountInfoFlags

			case "get-account-portfolio":
				epf = guildsServiceGetAccountPortfolioFlags

			case "get-account-portfolios":
				epf = guildsServiceGetAccountPortfoliosFlags

			}

		}
	}
	if epf == nil {
		return nil, nil, fmt.Errorf("unknown %q endpoint %q", svcn, epn)
	}

	// Parse endpoint flags if any
	if svcf.NArg() > 1 {
		if err := epf.Parse(svcf.Args()[1:]); err != nil {
			return nil, nil, err
		}
	}

	var (
		data     interface{}
		endpoint goa.Endpoint
		err      error
	)
	{
		switch svcn {
		case "guilds-service":
			c := guildsservicec.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "get-all-guilds":
				endpoint = c.GetAllGuilds()
				data = nil
			case "get-single-guild":
				endpoint = c.GetSingleGuild()
				data, err = guildsservicec.BuildGetSingleGuildPayload(*guildsServiceGetSingleGuildGuildIDFlag)
			case "get-guild-members":
				endpoint = c.GetGuildMembers()
				data, err = guildsservicec.BuildGetGuildMembersPayload(*guildsServiceGetGuildMembersGuildIDFlag)
			case "get-guild-master-address":
				endpoint = c.GetGuildMasterAddress()
				data, err = guildsservicec.BuildGetGuildMasterAddressPayload(*guildsServiceGetGuildMasterAddressGuildIDFlag)
			case "get-guild-default-member":
				endpoint = c.GetGuildDefaultMember()
				data, err = guildsservicec.BuildGetGuildDefaultMemberPayload(*guildsServiceGetGuildDefaultMemberGuildIDFlag)
			case "enter-guild":
				endpoint = c.EnterGuild()
				data, err = guildsservicec.BuildEnterGuildPayload(*guildsServiceEnterGuildBodyFlag, *guildsServiceEnterGuildGuildIDFlag)
			case "leave-guild":
				endpoint = c.LeaveGuild()
				data, err = guildsservicec.BuildLeaveGuildPayload(*guildsServiceLeaveGuildGuildIDFlag, *guildsServiceLeaveGuildInjectiveAddressFlag)
			case "get-guild-markets":
				endpoint = c.GetGuildMarkets()
				data, err = guildsservicec.BuildGetGuildMarketsPayload(*guildsServiceGetGuildMarketsGuildIDFlag)
			case "get-guild-portfolios":
				endpoint = c.GetGuildPortfolios()
				data, err = guildsservicec.BuildGetGuildPortfoliosPayload(*guildsServiceGetGuildPortfoliosGuildIDFlag, *guildsServiceGetGuildPortfoliosStartTimeFlag, *guildsServiceGetGuildPortfoliosEndTimeFlag)
			case "get-account-info":
				endpoint = c.GetAccountInfo()
				data, err = guildsservicec.BuildGetAccountInfoPayload(*guildsServiceGetAccountInfoInjectiveAddressFlag)
			case "get-account-portfolio":
				endpoint = c.GetAccountPortfolio()
				data, err = guildsservicec.BuildGetAccountPortfolioPayload(*guildsServiceGetAccountPortfolioInjectiveAddressFlag)
			case "get-account-portfolios":
				endpoint = c.GetAccountPortfolios()
				data, err = guildsservicec.BuildGetAccountPortfoliosPayload(*guildsServiceGetAccountPortfoliosInjectiveAddressFlag, *guildsServiceGetAccountPortfoliosStartTimeFlag, *guildsServiceGetAccountPortfoliosEndTimeFlag)
			}
		}
	}
	if err != nil {
		return nil, nil, err
	}

	return endpoint, data, nil
}

// guilds-serviceUsage displays the usage of the guilds-service command and its
// subcommands.
func guildsServiceUsage() {
	fmt.Fprintf(os.Stderr, `Service supports trading guild queries
Usage:
    %[1]s [globalflags] guilds-service COMMAND [flags]

COMMAND:
    get-all-guilds: Get all guilds
    get-single-guild: Get a single guild base on ID
    get-guild-members: Get all members a given guild (include default member)
    get-guild-master-address: Get master address of given guild
    get-guild-default-member: Get default guild member
    enter-guild: Enter the guild
    leave-guild: Leave the guild, guildID
    get-guild-markets: Get the guild markets
    get-guild-portfolios: Get the guild markets
    get-account-info: Get current account member status
    get-account-portfolio: Get current account portfolio snapshot
    get-account-portfolios: Get current account portfolios snapshots all the time

Additional help:
    %[1]s guilds-service COMMAND --help
`, os.Args[0])
}
func guildsServiceGetAllGuildsUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] guilds-service get-all-guilds

Get all guilds

Example:
    %[1]s guilds-service get-all-guilds
`, os.Args[0])
}

func guildsServiceGetSingleGuildUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] guilds-service get-single-guild -guild-id STRING

Get a single guild base on ID
    -guild-id STRING: 

Example:
    %[1]s guilds-service get-single-guild --guild-id "Magni nihil voluptatem voluptatibus nostrum magnam."
`, os.Args[0])
}

func guildsServiceGetGuildMembersUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] guilds-service get-guild-members -guild-id STRING

Get all members a given guild (include default member)
    -guild-id STRING: 

Example:
    %[1]s guilds-service get-guild-members --guild-id "Alias voluptas soluta id quo dolorem aut."
`, os.Args[0])
}

func guildsServiceGetGuildMasterAddressUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] guilds-service get-guild-master-address -guild-id STRING

Get master address of given guild
    -guild-id STRING: 

Example:
    %[1]s guilds-service get-guild-master-address --guild-id "Adipisci reiciendis minima id."
`, os.Args[0])
}

func guildsServiceGetGuildDefaultMemberUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] guilds-service get-guild-default-member -guild-id STRING

Get default guild member
    -guild-id STRING: 

Example:
    %[1]s guilds-service get-guild-default-member --guild-id "Adipisci sed libero a nam consectetur."
`, os.Args[0])
}

func guildsServiceEnterGuildUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] guilds-service enter-guild -body JSON -guild-id STRING

Enter the guild
    -body JSON: 
    -guild-id STRING: 

Example:
    %[1]s guilds-service enter-guild --body '{
      "injective_address": "Voluptates est."
   }' --guild-id "Harum eum vel quia."
`, os.Args[0])
}

func guildsServiceLeaveGuildUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] guilds-service leave-guild -guild-id STRING -injective-address STRING

Leave the guild, guildID
    -guild-id STRING: 
    -injective-address STRING: 

Example:
    %[1]s guilds-service leave-guild --guild-id "Quia natus iste eaque." --injective-address "Vel provident odio deserunt quas."
`, os.Args[0])
}

func guildsServiceGetGuildMarketsUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] guilds-service get-guild-markets -guild-id STRING

Get the guild markets
    -guild-id STRING: 

Example:
    %[1]s guilds-service get-guild-markets --guild-id "Molestiae cum."
`, os.Args[0])
}

func guildsServiceGetGuildPortfoliosUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] guilds-service get-guild-portfolios -guild-id STRING -start-time INT64 -end-time INT64

Get the guild markets
    -guild-id STRING: 
    -start-time INT64: 
    -end-time INT64: 

Example:
    %[1]s guilds-service get-guild-portfolios --guild-id "Accusantium nobis quia." --start-time 8038232831277701001 --end-time 6787478751059554332
`, os.Args[0])
}

func guildsServiceGetAccountInfoUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] guilds-service get-account-info -injective-address STRING

Get current account member status
    -injective-address STRING: 

Example:
    %[1]s guilds-service get-account-info --injective-address "Et ipsum."
`, os.Args[0])
}

func guildsServiceGetAccountPortfolioUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] guilds-service get-account-portfolio -injective-address STRING

Get current account portfolio snapshot
    -injective-address STRING: 

Example:
    %[1]s guilds-service get-account-portfolio --injective-address "Voluptatum magnam dolorem nostrum velit non ipsa."
`, os.Args[0])
}

func guildsServiceGetAccountPortfoliosUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] guilds-service get-account-portfolios -injective-address STRING -start-time INT64 -end-time INT64

Get current account portfolios snapshots all the time
    -injective-address STRING: 
    -start-time INT64: 
    -end-time INT64: 

Example:
    %[1]s guilds-service get-account-portfolios --injective-address "Totam fugit possimus et nostrum voluptatem." --start-time 6149130476220960367 --end-time 8186644215274380022
`, os.Args[0])
}
