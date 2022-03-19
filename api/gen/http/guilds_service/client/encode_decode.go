// Code generated by goa v3.6.2, DO NOT EDIT.
//
// GuildsService HTTP client encoders and decoders
//
// Command:
// $ goa gen github.com/InjectiveLabs/injective-guilds-service/api/design -o ../

package client

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"

	guildsservice "github.com/InjectiveLabs/injective-guilds-service/api/gen/guilds_service"
	goahttp "goa.design/goa/v3/http"
)

// BuildGetAllGuildsRequest instantiates a HTTP request object with method and
// path set to call the "GuildsService" service "GetAllGuilds" endpoint
func (c *Client) BuildGetAllGuildsRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: GetAllGuildsGuildsServicePath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("GuildsService", "GetAllGuilds", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeGetAllGuildsResponse returns a decoder for responses returned by the
// GuildsService GetAllGuilds endpoint. restoreBody controls whether the
// response body should be restored after having been read.
// DecodeGetAllGuildsResponse may return the following errors:
//	- "not_found" (type *goa.ServiceError): 5
//	- "internal" (type *goa.ServiceError): 13
//	- error: internal error
func DecodeGetAllGuildsResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body GetAllGuildsResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("GuildsService", "GetAllGuilds", err)
			}
			err = ValidateGetAllGuildsResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("GuildsService", "GetAllGuilds", err)
			}
			res := NewGetAllGuildsResultOK(&body)
			return res, nil
		case 5:
			var (
				body GetAllGuildsNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("GuildsService", "GetAllGuilds", err)
			}
			err = ValidateGetAllGuildsNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("GuildsService", "GetAllGuilds", err)
			}
			return nil, NewGetAllGuildsNotFound(&body)
		case 13:
			var (
				body GetAllGuildsInternalResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("GuildsService", "GetAllGuilds", err)
			}
			err = ValidateGetAllGuildsInternalResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("GuildsService", "GetAllGuilds", err)
			}
			return nil, NewGetAllGuildsInternal(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("GuildsService", "GetAllGuilds", resp.StatusCode, string(body))
		}
	}
}

// BuildGetSingleGuildRequest instantiates a HTTP request object with method
// and path set to call the "GuildsService" service "GetSingleGuild" endpoint
func (c *Client) BuildGetSingleGuildRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		guildID string
	)
	{
		p, ok := v.(*guildsservice.GetSingleGuildPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("GuildsService", "GetSingleGuild", "*guildsservice.GetSingleGuildPayload", v)
		}
		guildID = p.GuildID
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: GetSingleGuildGuildsServicePath(guildID)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("GuildsService", "GetSingleGuild", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeGetSingleGuildResponse returns a decoder for responses returned by the
// GuildsService GetSingleGuild endpoint. restoreBody controls whether the
// response body should be restored after having been read.
// DecodeGetSingleGuildResponse may return the following errors:
//	- "not_found" (type *goa.ServiceError): 5
//	- "internal" (type *goa.ServiceError): 13
//	- error: internal error
func DecodeGetSingleGuildResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body GetSingleGuildResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("GuildsService", "GetSingleGuild", err)
			}
			err = ValidateGetSingleGuildResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("GuildsService", "GetSingleGuild", err)
			}
			res := NewGetSingleGuildResultOK(&body)
			return res, nil
		case 5:
			var (
				body GetSingleGuildNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("GuildsService", "GetSingleGuild", err)
			}
			err = ValidateGetSingleGuildNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("GuildsService", "GetSingleGuild", err)
			}
			return nil, NewGetSingleGuildNotFound(&body)
		case 13:
			var (
				body GetSingleGuildInternalResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("GuildsService", "GetSingleGuild", err)
			}
			err = ValidateGetSingleGuildInternalResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("GuildsService", "GetSingleGuild", err)
			}
			return nil, NewGetSingleGuildInternal(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("GuildsService", "GetSingleGuild", resp.StatusCode, string(body))
		}
	}
}

// BuildGetGuildMembersRequest instantiates a HTTP request object with method
// and path set to call the "GuildsService" service "GetGuildMembers" endpoint
func (c *Client) BuildGetGuildMembersRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		guildID string
	)
	{
		p, ok := v.(*guildsservice.GetGuildMembersPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("GuildsService", "GetGuildMembers", "*guildsservice.GetGuildMembersPayload", v)
		}
		guildID = p.GuildID
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: GetGuildMembersGuildsServicePath(guildID)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("GuildsService", "GetGuildMembers", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeGetGuildMembersResponse returns a decoder for responses returned by
// the GuildsService GetGuildMembers endpoint. restoreBody controls whether the
// response body should be restored after having been read.
// DecodeGetGuildMembersResponse may return the following errors:
//	- "not_found" (type *goa.ServiceError): 5
//	- "internal" (type *goa.ServiceError): 13
//	- error: internal error
func DecodeGetGuildMembersResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body GetGuildMembersResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("GuildsService", "GetGuildMembers", err)
			}
			err = ValidateGetGuildMembersResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("GuildsService", "GetGuildMembers", err)
			}
			res := NewGetGuildMembersResultOK(&body)
			return res, nil
		case 5:
			var (
				body GetGuildMembersNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("GuildsService", "GetGuildMembers", err)
			}
			err = ValidateGetGuildMembersNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("GuildsService", "GetGuildMembers", err)
			}
			return nil, NewGetGuildMembersNotFound(&body)
		case 13:
			var (
				body GetGuildMembersInternalResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("GuildsService", "GetGuildMembers", err)
			}
			err = ValidateGetGuildMembersInternalResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("GuildsService", "GetGuildMembers", err)
			}
			return nil, NewGetGuildMembersInternal(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("GuildsService", "GetGuildMembers", resp.StatusCode, string(body))
		}
	}
}

// BuildGetGuildMasterAddressRequest instantiates a HTTP request object with
// method and path set to call the "GuildsService" service
// "GetGuildMasterAddress" endpoint
func (c *Client) BuildGetGuildMasterAddressRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		guildID string
	)
	{
		p, ok := v.(*guildsservice.GetGuildMasterAddressPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("GuildsService", "GetGuildMasterAddress", "*guildsservice.GetGuildMasterAddressPayload", v)
		}
		guildID = p.GuildID
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: GetGuildMasterAddressGuildsServicePath(guildID)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("GuildsService", "GetGuildMasterAddress", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeGetGuildMasterAddressResponse returns a decoder for responses returned
// by the GuildsService GetGuildMasterAddress endpoint. restoreBody controls
// whether the response body should be restored after having been read.
// DecodeGetGuildMasterAddressResponse may return the following errors:
//	- "not_found" (type *goa.ServiceError): 5
//	- "internal" (type *goa.ServiceError): 13
//	- error: internal error
func DecodeGetGuildMasterAddressResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body GetGuildMasterAddressResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("GuildsService", "GetGuildMasterAddress", err)
			}
			res := NewGetGuildMasterAddressResultOK(&body)
			return res, nil
		case 5:
			var (
				body GetGuildMasterAddressNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("GuildsService", "GetGuildMasterAddress", err)
			}
			err = ValidateGetGuildMasterAddressNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("GuildsService", "GetGuildMasterAddress", err)
			}
			return nil, NewGetGuildMasterAddressNotFound(&body)
		case 13:
			var (
				body GetGuildMasterAddressInternalResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("GuildsService", "GetGuildMasterAddress", err)
			}
			err = ValidateGetGuildMasterAddressInternalResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("GuildsService", "GetGuildMasterAddress", err)
			}
			return nil, NewGetGuildMasterAddressInternal(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("GuildsService", "GetGuildMasterAddress", resp.StatusCode, string(body))
		}
	}
}

// BuildGetGuildDefaultMemberRequest instantiates a HTTP request object with
// method and path set to call the "GuildsService" service
// "GetGuildDefaultMember" endpoint
func (c *Client) BuildGetGuildDefaultMemberRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		guildID string
	)
	{
		p, ok := v.(*guildsservice.GetGuildDefaultMemberPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("GuildsService", "GetGuildDefaultMember", "*guildsservice.GetGuildDefaultMemberPayload", v)
		}
		guildID = p.GuildID
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: GetGuildDefaultMemberGuildsServicePath(guildID)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("GuildsService", "GetGuildDefaultMember", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeGetGuildDefaultMemberResponse returns a decoder for responses returned
// by the GuildsService GetGuildDefaultMember endpoint. restoreBody controls
// whether the response body should be restored after having been read.
// DecodeGetGuildDefaultMemberResponse may return the following errors:
//	- "not_found" (type *goa.ServiceError): 5
//	- "internal" (type *goa.ServiceError): 13
//	- error: internal error
func DecodeGetGuildDefaultMemberResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body GetGuildDefaultMemberResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("GuildsService", "GetGuildDefaultMember", err)
			}
			err = ValidateGetGuildDefaultMemberResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("GuildsService", "GetGuildDefaultMember", err)
			}
			res := NewGetGuildDefaultMemberResultOK(&body)
			return res, nil
		case 5:
			var (
				body GetGuildDefaultMemberNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("GuildsService", "GetGuildDefaultMember", err)
			}
			err = ValidateGetGuildDefaultMemberNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("GuildsService", "GetGuildDefaultMember", err)
			}
			return nil, NewGetGuildDefaultMemberNotFound(&body)
		case 13:
			var (
				body GetGuildDefaultMemberInternalResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("GuildsService", "GetGuildDefaultMember", err)
			}
			err = ValidateGetGuildDefaultMemberInternalResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("GuildsService", "GetGuildDefaultMember", err)
			}
			return nil, NewGetGuildDefaultMemberInternal(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("GuildsService", "GetGuildDefaultMember", resp.StatusCode, string(body))
		}
	}
}

// BuildEnterGuildRequest instantiates a HTTP request object with method and
// path set to call the "GuildsService" service "EnterGuild" endpoint
func (c *Client) BuildEnterGuildRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		guildID string
	)
	{
		p, ok := v.(*guildsservice.EnterGuildPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("GuildsService", "EnterGuild", "*guildsservice.EnterGuildPayload", v)
		}
		if p.GuildID != nil {
			guildID = *p.GuildID
		}
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: EnterGuildGuildsServicePath(guildID)}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("GuildsService", "EnterGuild", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeEnterGuildRequest returns an encoder for requests sent to the
// GuildsService EnterGuild server.
func EncodeEnterGuildRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*guildsservice.EnterGuildPayload)
		if !ok {
			return goahttp.ErrInvalidType("GuildsService", "EnterGuild", "*guildsservice.EnterGuildPayload", v)
		}
		body := NewEnterGuildRequestBody(p)
		if err := encoder(req).Encode(&body); err != nil {
			return goahttp.ErrEncodingError("GuildsService", "EnterGuild", err)
		}
		return nil
	}
}

// DecodeEnterGuildResponse returns a decoder for responses returned by the
// GuildsService EnterGuild endpoint. restoreBody controls whether the response
// body should be restored after having been read.
// DecodeEnterGuildResponse may return the following errors:
//	- "not_found" (type *goa.ServiceError): 5
//	- "internal" (type *goa.ServiceError): 13
//	- error: internal error
func DecodeEnterGuildResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body EnterGuildResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("GuildsService", "EnterGuild", err)
			}
			res := NewEnterGuildResultOK(&body)
			return res, nil
		case 5:
			var (
				body EnterGuildNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("GuildsService", "EnterGuild", err)
			}
			err = ValidateEnterGuildNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("GuildsService", "EnterGuild", err)
			}
			return nil, NewEnterGuildNotFound(&body)
		case 13:
			var (
				body EnterGuildInternalResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("GuildsService", "EnterGuild", err)
			}
			err = ValidateEnterGuildInternalResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("GuildsService", "EnterGuild", err)
			}
			return nil, NewEnterGuildInternal(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("GuildsService", "EnterGuild", resp.StatusCode, string(body))
		}
	}
}

// BuildLeaveGuildRequest instantiates a HTTP request object with method and
// path set to call the "GuildsService" service "LeaveGuild" endpoint
func (c *Client) BuildLeaveGuildRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		guildID string
	)
	{
		p, ok := v.(*guildsservice.LeaveGuildPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("GuildsService", "LeaveGuild", "*guildsservice.LeaveGuildPayload", v)
		}
		if p.GuildID != nil {
			guildID = *p.GuildID
		}
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: LeaveGuildGuildsServicePath(guildID)}
	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("GuildsService", "LeaveGuild", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeLeaveGuildRequest returns an encoder for requests sent to the
// GuildsService LeaveGuild server.
func EncodeLeaveGuildRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*guildsservice.LeaveGuildPayload)
		if !ok {
			return goahttp.ErrInvalidType("GuildsService", "LeaveGuild", "*guildsservice.LeaveGuildPayload", v)
		}
		body := NewLeaveGuildRequestBody(p)
		if err := encoder(req).Encode(&body); err != nil {
			return goahttp.ErrEncodingError("GuildsService", "LeaveGuild", err)
		}
		return nil
	}
}

// DecodeLeaveGuildResponse returns a decoder for responses returned by the
// GuildsService LeaveGuild endpoint. restoreBody controls whether the response
// body should be restored after having been read.
// DecodeLeaveGuildResponse may return the following errors:
//	- "not_found" (type *goa.ServiceError): 5
//	- "internal" (type *goa.ServiceError): 13
//	- error: internal error
func DecodeLeaveGuildResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body LeaveGuildResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("GuildsService", "LeaveGuild", err)
			}
			res := NewLeaveGuildResultOK(&body)
			return res, nil
		case 5:
			var (
				body LeaveGuildNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("GuildsService", "LeaveGuild", err)
			}
			err = ValidateLeaveGuildNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("GuildsService", "LeaveGuild", err)
			}
			return nil, NewLeaveGuildNotFound(&body)
		case 13:
			var (
				body LeaveGuildInternalResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("GuildsService", "LeaveGuild", err)
			}
			err = ValidateLeaveGuildInternalResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("GuildsService", "LeaveGuild", err)
			}
			return nil, NewLeaveGuildInternal(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("GuildsService", "LeaveGuild", resp.StatusCode, string(body))
		}
	}
}

// BuildGetGuildMarketsRequest instantiates a HTTP request object with method
// and path set to call the "GuildsService" service "GetGuildMarkets" endpoint
func (c *Client) BuildGetGuildMarketsRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		guildID string
	)
	{
		p, ok := v.(*guildsservice.GetGuildMarketsPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("GuildsService", "GetGuildMarkets", "*guildsservice.GetGuildMarketsPayload", v)
		}
		guildID = p.GuildID
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: GetGuildMarketsGuildsServicePath(guildID)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("GuildsService", "GetGuildMarkets", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeGetGuildMarketsResponse returns a decoder for responses returned by
// the GuildsService GetGuildMarkets endpoint. restoreBody controls whether the
// response body should be restored after having been read.
// DecodeGetGuildMarketsResponse may return the following errors:
//	- "not_found" (type *goa.ServiceError): 5
//	- "internal" (type *goa.ServiceError): 13
//	- error: internal error
func DecodeGetGuildMarketsResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body GetGuildMarketsResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("GuildsService", "GetGuildMarkets", err)
			}
			err = ValidateGetGuildMarketsResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("GuildsService", "GetGuildMarkets", err)
			}
			res := NewGetGuildMarketsResultOK(&body)
			return res, nil
		case 5:
			var (
				body GetGuildMarketsNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("GuildsService", "GetGuildMarkets", err)
			}
			err = ValidateGetGuildMarketsNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("GuildsService", "GetGuildMarkets", err)
			}
			return nil, NewGetGuildMarketsNotFound(&body)
		case 13:
			var (
				body GetGuildMarketsInternalResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("GuildsService", "GetGuildMarkets", err)
			}
			err = ValidateGetGuildMarketsInternalResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("GuildsService", "GetGuildMarkets", err)
			}
			return nil, NewGetGuildMarketsInternal(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("GuildsService", "GetGuildMarkets", resp.StatusCode, string(body))
		}
	}
}

// BuildGetAccountPortfolioRequest instantiates a HTTP request object with
// method and path set to call the "GuildsService" service
// "GetAccountPortfolio" endpoint
func (c *Client) BuildGetAccountPortfolioRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		guildID string
	)
	{
		p, ok := v.(*guildsservice.GetAccountPortfolioPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("GuildsService", "GetAccountPortfolio", "*guildsservice.GetAccountPortfolioPayload", v)
		}
		guildID = p.GuildID
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: GetAccountPortfolioGuildsServicePath(guildID)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("GuildsService", "GetAccountPortfolio", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeGetAccountPortfolioRequest returns an encoder for requests sent to the
// GuildsService GetAccountPortfolio server.
func EncodeGetAccountPortfolioRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*guildsservice.GetAccountPortfolioPayload)
		if !ok {
			return goahttp.ErrInvalidType("GuildsService", "GetAccountPortfolio", "*guildsservice.GetAccountPortfolioPayload", v)
		}
		values := req.URL.Query()
		values.Add("injective_address", p.InjectiveAddress)
		req.URL.RawQuery = values.Encode()
		return nil
	}
}

// DecodeGetAccountPortfolioResponse returns a decoder for responses returned
// by the GuildsService GetAccountPortfolio endpoint. restoreBody controls
// whether the response body should be restored after having been read.
// DecodeGetAccountPortfolioResponse may return the following errors:
//	- "not_found" (type *goa.ServiceError): 5
//	- "internal" (type *goa.ServiceError): 13
//	- error: internal error
func DecodeGetAccountPortfolioResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body GetAccountPortfolioResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("GuildsService", "GetAccountPortfolio", err)
			}
			err = ValidateGetAccountPortfolioResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("GuildsService", "GetAccountPortfolio", err)
			}
			res := NewGetAccountPortfolioResultOK(&body)
			return res, nil
		case 5:
			var (
				body GetAccountPortfolioNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("GuildsService", "GetAccountPortfolio", err)
			}
			err = ValidateGetAccountPortfolioNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("GuildsService", "GetAccountPortfolio", err)
			}
			return nil, NewGetAccountPortfolioNotFound(&body)
		case 13:
			var (
				body GetAccountPortfolioInternalResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("GuildsService", "GetAccountPortfolio", err)
			}
			err = ValidateGetAccountPortfolioInternalResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("GuildsService", "GetAccountPortfolio", err)
			}
			return nil, NewGetAccountPortfolioInternal(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("GuildsService", "GetAccountPortfolio", resp.StatusCode, string(body))
		}
	}
}

// BuildGetAccountPortfoliosRequest instantiates a HTTP request object with
// method and path set to call the "GuildsService" service
// "GetAccountPortfolios" endpoint
func (c *Client) BuildGetAccountPortfoliosRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		guildID string
	)
	{
		p, ok := v.(*guildsservice.GetAccountPortfoliosPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("GuildsService", "GetAccountPortfolios", "*guildsservice.GetAccountPortfoliosPayload", v)
		}
		guildID = p.GuildID
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: GetAccountPortfoliosGuildsServicePath(guildID)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("GuildsService", "GetAccountPortfolios", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeGetAccountPortfoliosRequest returns an encoder for requests sent to
// the GuildsService GetAccountPortfolios server.
func EncodeGetAccountPortfoliosRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*guildsservice.GetAccountPortfoliosPayload)
		if !ok {
			return goahttp.ErrInvalidType("GuildsService", "GetAccountPortfolios", "*guildsservice.GetAccountPortfoliosPayload", v)
		}
		values := req.URL.Query()
		values.Add("injective_address", p.InjectiveAddress)
		req.URL.RawQuery = values.Encode()
		return nil
	}
}

// DecodeGetAccountPortfoliosResponse returns a decoder for responses returned
// by the GuildsService GetAccountPortfolios endpoint. restoreBody controls
// whether the response body should be restored after having been read.
// DecodeGetAccountPortfoliosResponse may return the following errors:
//	- "not_found" (type *goa.ServiceError): 5
//	- "internal" (type *goa.ServiceError): 13
//	- error: internal error
func DecodeGetAccountPortfoliosResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body GetAccountPortfoliosResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("GuildsService", "GetAccountPortfolios", err)
			}
			err = ValidateGetAccountPortfoliosResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("GuildsService", "GetAccountPortfolios", err)
			}
			res := NewGetAccountPortfoliosResultOK(&body)
			return res, nil
		case 5:
			var (
				body GetAccountPortfoliosNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("GuildsService", "GetAccountPortfolios", err)
			}
			err = ValidateGetAccountPortfoliosNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("GuildsService", "GetAccountPortfolios", err)
			}
			return nil, NewGetAccountPortfoliosNotFound(&body)
		case 13:
			var (
				body GetAccountPortfoliosInternalResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("GuildsService", "GetAccountPortfolios", err)
			}
			err = ValidateGetAccountPortfoliosInternalResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("GuildsService", "GetAccountPortfolios", err)
			}
			return nil, NewGetAccountPortfoliosInternal(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("GuildsService", "GetAccountPortfolios", resp.StatusCode, string(body))
		}
	}
}

// unmarshalGuildResponseBodyToGuildsserviceGuild builds a value of type
// *guildsservice.Guild from a value of type *GuildResponseBody.
func unmarshalGuildResponseBodyToGuildsserviceGuild(v *GuildResponseBody) *guildsservice.Guild {
	if v == nil {
		return nil
	}
	res := &guildsservice.Guild{
		ID:                         *v.ID,
		Name:                       *v.Name,
		Description:                *v.Description,
		MasterAddress:              *v.MasterAddress,
		SpotBaseRequirement:        *v.SpotBaseRequirement,
		SpotQuoteRequirement:       *v.SpotQuoteRequirement,
		DerivativeQuoteRequirement: *v.DerivativeQuoteRequirement,
		StakingRequirement:         *v.StakingRequirement,
		Capacity:                   *v.Capacity,
		MemberCount:                *v.MemberCount,
	}

	return res
}

// unmarshalGuildMemberResponseBodyToGuildsserviceGuildMember builds a value of
// type *guildsservice.GuildMember from a value of type
// *GuildMemberResponseBody.
func unmarshalGuildMemberResponseBodyToGuildsserviceGuildMember(v *GuildMemberResponseBody) *guildsservice.GuildMember {
	if v == nil {
		return nil
	}
	res := &guildsservice.GuildMember{
		InjectiveAddress:     *v.InjectiveAddress,
		IsDefaultGuildMember: *v.IsDefaultGuildMember,
	}

	return res
}

// unmarshalMarketResponseBodyToGuildsserviceMarket builds a value of type
// *guildsservice.Market from a value of type *MarketResponseBody.
func unmarshalMarketResponseBodyToGuildsserviceMarket(v *MarketResponseBody) *guildsservice.Market {
	if v == nil {
		return nil
	}
	res := &guildsservice.Market{
		MarketID:    *v.MarketID,
		IsPerpetual: *v.IsPerpetual,
	}

	return res
}

// unmarshalSingleAccountPortfolioResponseBodyToGuildsserviceSingleAccountPortfolio
// builds a value of type *guildsservice.SingleAccountPortfolio from a value of
// type *SingleAccountPortfolioResponseBody.
func unmarshalSingleAccountPortfolioResponseBodyToGuildsserviceSingleAccountPortfolio(v *SingleAccountPortfolioResponseBody) *guildsservice.SingleAccountPortfolio {
	if v == nil {
		return nil
	}
	res := &guildsservice.SingleAccountPortfolio{
		InjectiveAddress: *v.InjectiveAddress,
		Denom:            *v.Denom,
		TotalBalance:     *v.TotalBalance,
		AvailableBalance: *v.AvailableBalance,
		UnrealizedPnl:    *v.UnrealizedPnl,
		MarginHold:       *v.MarginHold,
		UpdatedAt:        *v.UpdatedAt,
	}

	return res
}

// unmarshalAccountPorfoliosResponseBodyToGuildsserviceAccountPorfolios builds
// a value of type *guildsservice.AccountPorfolios from a value of type
// *AccountPorfoliosResponseBody.
func unmarshalAccountPorfoliosResponseBodyToGuildsserviceAccountPorfolios(v *AccountPorfoliosResponseBody) *guildsservice.AccountPorfolios {
	if v == nil {
		return nil
	}
	res := &guildsservice.AccountPorfolios{
		InjectiveAddress: *v.InjectiveAddress,
	}
	res.Portfolios = make([]*guildsservice.EmbededAccountPortfolio, len(v.Portfolios))
	for i, val := range v.Portfolios {
		res.Portfolios[i] = unmarshalEmbededAccountPortfolioResponseBodyToGuildsserviceEmbededAccountPortfolio(val)
	}

	return res
}

// unmarshalEmbededAccountPortfolioResponseBodyToGuildsserviceEmbededAccountPortfolio
// builds a value of type *guildsservice.EmbededAccountPortfolio from a value
// of type *EmbededAccountPortfolioResponseBody.
func unmarshalEmbededAccountPortfolioResponseBodyToGuildsserviceEmbededAccountPortfolio(v *EmbededAccountPortfolioResponseBody) *guildsservice.EmbededAccountPortfolio {
	res := &guildsservice.EmbededAccountPortfolio{
		Denom:            *v.Denom,
		TotalBalance:     *v.TotalBalance,
		AvailableBalance: *v.AvailableBalance,
		UnrealizedPnl:    *v.UnrealizedPnl,
		MarginHold:       *v.MarginHold,
		UpdatedAt:        *v.UpdatedAt,
	}

	return res
}
