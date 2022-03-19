// Code generated by goa v3.6.2, DO NOT EDIT.
//
// GuildsService HTTP server encoders and decoders
//
// Command:
// $ goa gen github.com/InjectiveLabs/injective-guilds-service/api/design -o ../

package server

import (
	"context"
	"errors"
	"io"
	"net/http"

	guildsservice "github.com/InjectiveLabs/injective-guilds-service/api/gen/guilds_service"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// EncodeGetAllGuildsResponse returns an encoder for responses returned by the
// GuildsService GetAllGuilds endpoint.
func EncodeGetAllGuildsResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.(*guildsservice.GetAllGuildsResult)
		enc := encoder(ctx, w)
		body := NewGetAllGuildsResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// EncodeGetAllGuildsError returns an encoder for errors returned by the
// GetAllGuilds GuildsService endpoint.
func EncodeGetAllGuildsError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en ErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "not_found":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewGetAllGuildsNotFoundResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(5)
			return enc.Encode(body)
		case "internal":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewGetAllGuildsInternalResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(13)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeGetSingleGuildResponse returns an encoder for responses returned by
// the GuildsService GetSingleGuild endpoint.
func EncodeGetSingleGuildResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.(*guildsservice.GetSingleGuildResult)
		enc := encoder(ctx, w)
		body := NewGetSingleGuildResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeGetSingleGuildRequest returns a decoder for requests sent to the
// GuildsService GetSingleGuild endpoint.
func DecodeGetSingleGuildRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			guildID string

			params = mux.Vars(r)
		)
		guildID = params["guildID"]
		payload := NewGetSingleGuildPayload(guildID)

		return payload, nil
	}
}

// EncodeGetSingleGuildError returns an encoder for errors returned by the
// GetSingleGuild GuildsService endpoint.
func EncodeGetSingleGuildError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en ErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "not_found":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewGetSingleGuildNotFoundResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(5)
			return enc.Encode(body)
		case "internal":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewGetSingleGuildInternalResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(13)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeGetGuildMembersResponse returns an encoder for responses returned by
// the GuildsService GetGuildMembers endpoint.
func EncodeGetGuildMembersResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.(*guildsservice.GetGuildMembersResult)
		enc := encoder(ctx, w)
		body := NewGetGuildMembersResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeGetGuildMembersRequest returns a decoder for requests sent to the
// GuildsService GetGuildMembers endpoint.
func DecodeGetGuildMembersRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			guildID string

			params = mux.Vars(r)
		)
		guildID = params["guildID"]
		payload := NewGetGuildMembersPayload(guildID)

		return payload, nil
	}
}

// EncodeGetGuildMembersError returns an encoder for errors returned by the
// GetGuildMembers GuildsService endpoint.
func EncodeGetGuildMembersError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en ErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "not_found":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewGetGuildMembersNotFoundResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(5)
			return enc.Encode(body)
		case "internal":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewGetGuildMembersInternalResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(13)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeGetGuildMasterAddressResponse returns an encoder for responses
// returned by the GuildsService GetGuildMasterAddress endpoint.
func EncodeGetGuildMasterAddressResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.(*guildsservice.GetGuildMasterAddressResult)
		enc := encoder(ctx, w)
		body := NewGetGuildMasterAddressResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeGetGuildMasterAddressRequest returns a decoder for requests sent to
// the GuildsService GetGuildMasterAddress endpoint.
func DecodeGetGuildMasterAddressRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			guildID string

			params = mux.Vars(r)
		)
		guildID = params["guildID"]
		payload := NewGetGuildMasterAddressPayload(guildID)

		return payload, nil
	}
}

// EncodeGetGuildMasterAddressError returns an encoder for errors returned by
// the GetGuildMasterAddress GuildsService endpoint.
func EncodeGetGuildMasterAddressError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en ErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "not_found":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewGetGuildMasterAddressNotFoundResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(5)
			return enc.Encode(body)
		case "internal":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewGetGuildMasterAddressInternalResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(13)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeGetGuildDefaultMemberResponse returns an encoder for responses
// returned by the GuildsService GetGuildDefaultMember endpoint.
func EncodeGetGuildDefaultMemberResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.(*guildsservice.GetGuildDefaultMemberResult)
		enc := encoder(ctx, w)
		body := NewGetGuildDefaultMemberResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeGetGuildDefaultMemberRequest returns a decoder for requests sent to
// the GuildsService GetGuildDefaultMember endpoint.
func DecodeGetGuildDefaultMemberRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			guildID string

			params = mux.Vars(r)
		)
		guildID = params["guildID"]
		payload := NewGetGuildDefaultMemberPayload(guildID)

		return payload, nil
	}
}

// EncodeGetGuildDefaultMemberError returns an encoder for errors returned by
// the GetGuildDefaultMember GuildsService endpoint.
func EncodeGetGuildDefaultMemberError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en ErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "not_found":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewGetGuildDefaultMemberNotFoundResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(5)
			return enc.Encode(body)
		case "internal":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewGetGuildDefaultMemberInternalResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(13)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeEnterGuildResponse returns an encoder for responses returned by the
// GuildsService EnterGuild endpoint.
func EncodeEnterGuildResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.(*guildsservice.EnterGuildResult)
		enc := encoder(ctx, w)
		body := NewEnterGuildResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeEnterGuildRequest returns a decoder for requests sent to the
// GuildsService EnterGuild endpoint.
func DecodeEnterGuildRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			body EnterGuildRequestBody
			err  error
		)
		err = decoder(r).Decode(&body)
		if err != nil {
			if err == io.EOF {
				return nil, goa.MissingPayloadError()
			}
			return nil, goa.DecodePayloadError(err.Error())
		}
		err = ValidateEnterGuildRequestBody(&body)
		if err != nil {
			return nil, err
		}

		var (
			guildID string

			params = mux.Vars(r)
		)
		guildID = params["guildID"]
		payload := NewEnterGuildPayload(&body, guildID)

		return payload, nil
	}
}

// EncodeEnterGuildError returns an encoder for errors returned by the
// EnterGuild GuildsService endpoint.
func EncodeEnterGuildError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en ErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "not_found":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewEnterGuildNotFoundResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(5)
			return enc.Encode(body)
		case "internal":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewEnterGuildInternalResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(13)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeLeaveGuildResponse returns an encoder for responses returned by the
// GuildsService LeaveGuild endpoint.
func EncodeLeaveGuildResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.(*guildsservice.LeaveGuildResult)
		enc := encoder(ctx, w)
		body := NewLeaveGuildResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeLeaveGuildRequest returns a decoder for requests sent to the
// GuildsService LeaveGuild endpoint.
func DecodeLeaveGuildRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			body LeaveGuildRequestBody
			err  error
		)
		err = decoder(r).Decode(&body)
		if err != nil {
			if err == io.EOF {
				return nil, goa.MissingPayloadError()
			}
			return nil, goa.DecodePayloadError(err.Error())
		}
		err = ValidateLeaveGuildRequestBody(&body)
		if err != nil {
			return nil, err
		}

		var (
			guildID string

			params = mux.Vars(r)
		)
		guildID = params["guildID"]
		payload := NewLeaveGuildPayload(&body, guildID)

		return payload, nil
	}
}

// EncodeLeaveGuildError returns an encoder for errors returned by the
// LeaveGuild GuildsService endpoint.
func EncodeLeaveGuildError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en ErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "not_found":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewLeaveGuildNotFoundResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(5)
			return enc.Encode(body)
		case "internal":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewLeaveGuildInternalResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(13)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeGetGuildMarketsResponse returns an encoder for responses returned by
// the GuildsService GetGuildMarkets endpoint.
func EncodeGetGuildMarketsResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.(*guildsservice.GetGuildMarketsResult)
		enc := encoder(ctx, w)
		body := NewGetGuildMarketsResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeGetGuildMarketsRequest returns a decoder for requests sent to the
// GuildsService GetGuildMarkets endpoint.
func DecodeGetGuildMarketsRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			guildID string

			params = mux.Vars(r)
		)
		guildID = params["guildID"]
		payload := NewGetGuildMarketsPayload(guildID)

		return payload, nil
	}
}

// EncodeGetGuildMarketsError returns an encoder for errors returned by the
// GetGuildMarkets GuildsService endpoint.
func EncodeGetGuildMarketsError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en ErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "not_found":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewGetGuildMarketsNotFoundResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(5)
			return enc.Encode(body)
		case "internal":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewGetGuildMarketsInternalResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(13)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeGetAccountPortfolioResponse returns an encoder for responses returned
// by the GuildsService GetAccountPortfolio endpoint.
func EncodeGetAccountPortfolioResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.(*guildsservice.GetAccountPortfolioResult)
		enc := encoder(ctx, w)
		body := NewGetAccountPortfolioResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeGetAccountPortfolioRequest returns a decoder for requests sent to the
// GuildsService GetAccountPortfolio endpoint.
func DecodeGetAccountPortfolioRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			guildID          string
			injectiveAddress string
			err              error

			params = mux.Vars(r)
		)
		guildID = params["guildID"]
		injectiveAddress = r.URL.Query().Get("injective_address")
		if injectiveAddress == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("injective_address", "query string"))
		}
		if err != nil {
			return nil, err
		}
		payload := NewGetAccountPortfolioPayload(guildID, injectiveAddress)

		return payload, nil
	}
}

// EncodeGetAccountPortfolioError returns an encoder for errors returned by the
// GetAccountPortfolio GuildsService endpoint.
func EncodeGetAccountPortfolioError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en ErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "not_found":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewGetAccountPortfolioNotFoundResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(5)
			return enc.Encode(body)
		case "internal":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewGetAccountPortfolioInternalResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(13)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeGetAccountPortfoliosResponse returns an encoder for responses returned
// by the GuildsService GetAccountPortfolios endpoint.
func EncodeGetAccountPortfoliosResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.(*guildsservice.GetAccountPortfoliosResult)
		enc := encoder(ctx, w)
		body := NewGetAccountPortfoliosResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeGetAccountPortfoliosRequest returns a decoder for requests sent to the
// GuildsService GetAccountPortfolios endpoint.
func DecodeGetAccountPortfoliosRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			guildID          string
			injectiveAddress string
			err              error

			params = mux.Vars(r)
		)
		guildID = params["guildID"]
		injectiveAddress = r.URL.Query().Get("injective_address")
		if injectiveAddress == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("injective_address", "query string"))
		}
		if err != nil {
			return nil, err
		}
		payload := NewGetAccountPortfoliosPayload(guildID, injectiveAddress)

		return payload, nil
	}
}

// EncodeGetAccountPortfoliosError returns an encoder for errors returned by
// the GetAccountPortfolios GuildsService endpoint.
func EncodeGetAccountPortfoliosError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en ErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "not_found":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewGetAccountPortfoliosNotFoundResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(5)
			return enc.Encode(body)
		case "internal":
			var res *goa.ServiceError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			var body interface{}
			if formatter != nil {
				body = formatter(res)
			} else {
				body = NewGetAccountPortfoliosInternalResponseBody(res)
			}
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(13)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// marshalGuildsserviceGuildToGuildResponseBody builds a value of type
// *GuildResponseBody from a value of type *guildsservice.Guild.
func marshalGuildsserviceGuildToGuildResponseBody(v *guildsservice.Guild) *GuildResponseBody {
	if v == nil {
		return nil
	}
	res := &GuildResponseBody{
		ID:                         v.ID,
		Name:                       v.Name,
		Description:                v.Description,
		MasterAddress:              v.MasterAddress,
		SpotBaseRequirement:        v.SpotBaseRequirement,
		SpotQuoteRequirement:       v.SpotQuoteRequirement,
		DerivativeQuoteRequirement: v.DerivativeQuoteRequirement,
		StakingRequirement:         v.StakingRequirement,
		Capacity:                   v.Capacity,
		MemberCount:                v.MemberCount,
	}

	return res
}

// marshalGuildsserviceGuildMemberToGuildMemberResponseBody builds a value of
// type *GuildMemberResponseBody from a value of type
// *guildsservice.GuildMember.
func marshalGuildsserviceGuildMemberToGuildMemberResponseBody(v *guildsservice.GuildMember) *GuildMemberResponseBody {
	if v == nil {
		return nil
	}
	res := &GuildMemberResponseBody{
		InjectiveAddress:     v.InjectiveAddress,
		IsDefaultGuildMember: v.IsDefaultGuildMember,
	}

	return res
}

// marshalGuildsserviceMarketToMarketResponseBody builds a value of type
// *MarketResponseBody from a value of type *guildsservice.Market.
func marshalGuildsserviceMarketToMarketResponseBody(v *guildsservice.Market) *MarketResponseBody {
	if v == nil {
		return nil
	}
	res := &MarketResponseBody{
		MarketID:    v.MarketID,
		IsPerpetual: v.IsPerpetual,
	}

	return res
}

// marshalGuildsserviceSingleAccountPortfolioToSingleAccountPortfolioResponseBody
// builds a value of type *SingleAccountPortfolioResponseBody from a value of
// type *guildsservice.SingleAccountPortfolio.
func marshalGuildsserviceSingleAccountPortfolioToSingleAccountPortfolioResponseBody(v *guildsservice.SingleAccountPortfolio) *SingleAccountPortfolioResponseBody {
	if v == nil {
		return nil
	}
	res := &SingleAccountPortfolioResponseBody{
		InjectiveAddress: v.InjectiveAddress,
		UpdatedAt:        v.UpdatedAt,
	}
	if v.Balances != nil {
		res.Balances = make([]*BalanceResponseBody, len(v.Balances))
		for i, val := range v.Balances {
			res.Balances[i] = marshalGuildsserviceBalanceToBalanceResponseBody(val)
		}
	}

	return res
}

// marshalGuildsserviceBalanceToBalanceResponseBody builds a value of type
// *BalanceResponseBody from a value of type *guildsservice.Balance.
func marshalGuildsserviceBalanceToBalanceResponseBody(v *guildsservice.Balance) *BalanceResponseBody {
	res := &BalanceResponseBody{
		Denom:            v.Denom,
		TotalBalance:     v.TotalBalance,
		AvailableBalance: v.AvailableBalance,
		UnrealizedPnl:    v.UnrealizedPnl,
		MarginHold:       v.MarginHold,
	}

	return res
}
