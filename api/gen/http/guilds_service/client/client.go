// Code generated by goa v3.6.2, DO NOT EDIT.
//
// GuildsService client HTTP transport
//
// Command:
// $ goa gen github.com/InjectiveLabs/injective-guilds-service/api/design -o ../

package client

import (
	"context"
	"net/http"

	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// Client lists the GuildsService service endpoint HTTP clients.
type Client struct {
	// GetAllGuilds Doer is the HTTP client used to make requests to the
	// GetAllGuilds endpoint.
	GetAllGuildsDoer goahttp.Doer

	// GetSingleGuild Doer is the HTTP client used to make requests to the
	// GetSingleGuild endpoint.
	GetSingleGuildDoer goahttp.Doer

	// GetGuildMembers Doer is the HTTP client used to make requests to the
	// GetGuildMembers endpoint.
	GetGuildMembersDoer goahttp.Doer

	// GetGuildMasterAddress Doer is the HTTP client used to make requests to the
	// GetGuildMasterAddress endpoint.
	GetGuildMasterAddressDoer goahttp.Doer

	// GetGuildDefaultMember Doer is the HTTP client used to make requests to the
	// GetGuildDefaultMember endpoint.
	GetGuildDefaultMemberDoer goahttp.Doer

	// EnterGuild Doer is the HTTP client used to make requests to the EnterGuild
	// endpoint.
	EnterGuildDoer goahttp.Doer

	// LeaveGuild Doer is the HTTP client used to make requests to the LeaveGuild
	// endpoint.
	LeaveGuildDoer goahttp.Doer

	// GetGuildMarkets Doer is the HTTP client used to make requests to the
	// GetGuildMarkets endpoint.
	GetGuildMarketsDoer goahttp.Doer

	// GetGuildPortfolios Doer is the HTTP client used to make requests to the
	// GetGuildPortfolios endpoint.
	GetGuildPortfoliosDoer goahttp.Doer

	// GetAccountInfo Doer is the HTTP client used to make requests to the
	// GetAccountInfo endpoint.
	GetAccountInfoDoer goahttp.Doer

	// GetAccountPortfolio Doer is the HTTP client used to make requests to the
	// GetAccountPortfolio endpoint.
	GetAccountPortfolioDoer goahttp.Doer

	// GetAccountPortfolios Doer is the HTTP client used to make requests to the
	// GetAccountPortfolios endpoint.
	GetAccountPortfoliosDoer goahttp.Doer

	// CORS Doer is the HTTP client used to make requests to the  endpoint.
	CORSDoer goahttp.Doer

	// RestoreResponseBody controls whether the response bodies are reset after
	// decoding so they can be read again.
	RestoreResponseBody bool

	scheme  string
	host    string
	encoder func(*http.Request) goahttp.Encoder
	decoder func(*http.Response) goahttp.Decoder
}

// NewClient instantiates HTTP clients for all the GuildsService service
// servers.
func NewClient(
	scheme string,
	host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restoreBody bool,
) *Client {
	return &Client{
		GetAllGuildsDoer:          doer,
		GetSingleGuildDoer:        doer,
		GetGuildMembersDoer:       doer,
		GetGuildMasterAddressDoer: doer,
		GetGuildDefaultMemberDoer: doer,
		EnterGuildDoer:            doer,
		LeaveGuildDoer:            doer,
		GetGuildMarketsDoer:       doer,
		GetGuildPortfoliosDoer:    doer,
		GetAccountInfoDoer:        doer,
		GetAccountPortfolioDoer:   doer,
		GetAccountPortfoliosDoer:  doer,
		CORSDoer:                  doer,
		RestoreResponseBody:       restoreBody,
		scheme:                    scheme,
		host:                      host,
		decoder:                   dec,
		encoder:                   enc,
	}
}

// GetAllGuilds returns an endpoint that makes HTTP requests to the
// GuildsService service GetAllGuilds server.
func (c *Client) GetAllGuilds() goa.Endpoint {
	var (
		decodeResponse = DecodeGetAllGuildsResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildGetAllGuildsRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.GetAllGuildsDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("GuildsService", "GetAllGuilds", err)
		}
		return decodeResponse(resp)
	}
}

// GetSingleGuild returns an endpoint that makes HTTP requests to the
// GuildsService service GetSingleGuild server.
func (c *Client) GetSingleGuild() goa.Endpoint {
	var (
		decodeResponse = DecodeGetSingleGuildResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildGetSingleGuildRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.GetSingleGuildDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("GuildsService", "GetSingleGuild", err)
		}
		return decodeResponse(resp)
	}
}

// GetGuildMembers returns an endpoint that makes HTTP requests to the
// GuildsService service GetGuildMembers server.
func (c *Client) GetGuildMembers() goa.Endpoint {
	var (
		decodeResponse = DecodeGetGuildMembersResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildGetGuildMembersRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.GetGuildMembersDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("GuildsService", "GetGuildMembers", err)
		}
		return decodeResponse(resp)
	}
}

// GetGuildMasterAddress returns an endpoint that makes HTTP requests to the
// GuildsService service GetGuildMasterAddress server.
func (c *Client) GetGuildMasterAddress() goa.Endpoint {
	var (
		decodeResponse = DecodeGetGuildMasterAddressResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildGetGuildMasterAddressRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.GetGuildMasterAddressDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("GuildsService", "GetGuildMasterAddress", err)
		}
		return decodeResponse(resp)
	}
}

// GetGuildDefaultMember returns an endpoint that makes HTTP requests to the
// GuildsService service GetGuildDefaultMember server.
func (c *Client) GetGuildDefaultMember() goa.Endpoint {
	var (
		decodeResponse = DecodeGetGuildDefaultMemberResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildGetGuildDefaultMemberRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.GetGuildDefaultMemberDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("GuildsService", "GetGuildDefaultMember", err)
		}
		return decodeResponse(resp)
	}
}

// EnterGuild returns an endpoint that makes HTTP requests to the GuildsService
// service EnterGuild server.
func (c *Client) EnterGuild() goa.Endpoint {
	var (
		encodeRequest  = EncodeEnterGuildRequest(c.encoder)
		decodeResponse = DecodeEnterGuildResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildEnterGuildRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.EnterGuildDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("GuildsService", "EnterGuild", err)
		}
		return decodeResponse(resp)
	}
}

// LeaveGuild returns an endpoint that makes HTTP requests to the GuildsService
// service LeaveGuild server.
func (c *Client) LeaveGuild() goa.Endpoint {
	var (
		encodeRequest  = EncodeLeaveGuildRequest(c.encoder)
		decodeResponse = DecodeLeaveGuildResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildLeaveGuildRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.LeaveGuildDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("GuildsService", "LeaveGuild", err)
		}
		return decodeResponse(resp)
	}
}

// GetGuildMarkets returns an endpoint that makes HTTP requests to the
// GuildsService service GetGuildMarkets server.
func (c *Client) GetGuildMarkets() goa.Endpoint {
	var (
		decodeResponse = DecodeGetGuildMarketsResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildGetGuildMarketsRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.GetGuildMarketsDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("GuildsService", "GetGuildMarkets", err)
		}
		return decodeResponse(resp)
	}
}

// GetGuildPortfolios returns an endpoint that makes HTTP requests to the
// GuildsService service GetGuildPortfolios server.
func (c *Client) GetGuildPortfolios() goa.Endpoint {
	var (
		encodeRequest  = EncodeGetGuildPortfoliosRequest(c.encoder)
		decodeResponse = DecodeGetGuildPortfoliosResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildGetGuildPortfoliosRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.GetGuildPortfoliosDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("GuildsService", "GetGuildPortfolios", err)
		}
		return decodeResponse(resp)
	}
}

// GetAccountInfo returns an endpoint that makes HTTP requests to the
// GuildsService service GetAccountInfo server.
func (c *Client) GetAccountInfo() goa.Endpoint {
	var (
		decodeResponse = DecodeGetAccountInfoResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildGetAccountInfoRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.GetAccountInfoDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("GuildsService", "GetAccountInfo", err)
		}
		return decodeResponse(resp)
	}
}

// GetAccountPortfolio returns an endpoint that makes HTTP requests to the
// GuildsService service GetAccountPortfolio server.
func (c *Client) GetAccountPortfolio() goa.Endpoint {
	var (
		decodeResponse = DecodeGetAccountPortfolioResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildGetAccountPortfolioRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.GetAccountPortfolioDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("GuildsService", "GetAccountPortfolio", err)
		}
		return decodeResponse(resp)
	}
}

// GetAccountPortfolios returns an endpoint that makes HTTP requests to the
// GuildsService service GetAccountPortfolios server.
func (c *Client) GetAccountPortfolios() goa.Endpoint {
	var (
		encodeRequest  = EncodeGetAccountPortfoliosRequest(c.encoder)
		decodeResponse = DecodeGetAccountPortfoliosResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildGetAccountPortfoliosRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.GetAccountPortfoliosDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("GuildsService", "GetAccountPortfolios", err)
		}
		return decodeResponse(resp)
	}
}
