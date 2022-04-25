package exchange

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/InjectiveLabs/metrics"
	accountPB "github.com/InjectiveLabs/sdk-go/exchange/accounts_rpc/pb"
	derivativeExchangePB "github.com/InjectiveLabs/sdk-go/exchange/derivative_exchange_rpc/pb"
	spotExchangePB "github.com/InjectiveLabs/sdk-go/exchange/spot_exchange_rpc/pb"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	log "github.com/xlab/suplog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

// TODO: improve sdk-go with new functions
// NOW:
// - Decided to use a stable release now, v1.31.0
// - Atm, we can clone functions into this repo first
// - We use response from our internal services
type exchangeProvider struct {
	DataProvider

	conn           *grpc.ClientConn
	lcdAddr        string
	assetPriceAddr string
	svcTags        metrics.Tags

	httpClient               *http.Client
	accountClient            accountPB.InjectiveAccountsRPCClient
	spotExchangeClient       spotExchangePB.InjectiveSpotExchangeRPCClient
	derivativeExchangeClient derivativeExchangePB.InjectiveDerivativeExchangeRPCClient
}

type ClientOptions struct {
	GasPrices string
	TLSCert   credentials.TransportCredentials
}

type ClientOption func(opts *ClientOptions) error

func OptionTLSCert(tlsCert credentials.TransportCredentials) ClientOption {
	return func(opts *ClientOptions) error {
		if tlsCert == nil {
			log.Infoln("Client does not use grpc secure transport")
		} else {
			log.Infoln("Succesfully load server TLS cert")
		}
		opts.TLSCert = tlsCert
		return nil
	}
}

// NewExchangeProvider returns DataProvider which fetches from many internal services
// (derives from current `master` of sdk-go)
func NewExchangeProvider(
	exchangeAddr string,
	lcdAddr string,
	assetPriceAddr string,
	options ...ClientOption,
) (DataProvider, error) {
	// process options
	opts := &ClientOptions{}
	for _, opt := range options {
		if err := opt(opts); err != nil {
			err = errors.Wrap(err, "error in client option")
			return nil, err
		}
	}

	var conn *grpc.ClientConn
	var err error
	if opts.TLSCert != nil {
		conn, err = grpc.Dial(exchangeAddr, grpc.WithTransportCredentials(opts.TLSCert), grpc.WithContextDialer(DialerFunc))
	} else {
		conn, err = grpc.Dial(exchangeAddr, grpc.WithInsecure(), grpc.WithContextDialer(DialerFunc))
	}

	if err != nil {
		err := errors.Wrapf(err, "failed to connect to the gRPC: %s", exchangeAddr)
		return nil, err
	}

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	cc := &exchangeProvider{
		conn:                     conn,
		lcdAddr:                  lcdAddr,
		httpClient:               httpClient,
		assetPriceAddr:           assetPriceAddr,
		accountClient:            accountPB.NewInjectiveAccountsRPCClient(conn),
		spotExchangeClient:       spotExchangePB.NewInjectiveSpotExchangeRPCClient(conn),
		derivativeExchangeClient: derivativeExchangePB.NewInjectiveDerivativeExchangeRPCClient(conn),
		svcTags: metrics.Tags{
			"svc": "data_providers",
		},
	}

	return cc, nil
}

func (p *exchangeProvider) GetSubaccountBalances(ctx context.Context, subaccount string) (result []*Balance, err error) {
	doneFn := metrics.ReportFuncTiming(p.svcTags)
	defer doneFn()
	metrics.ReportFuncCall(p.svcTags)

	// get all denoms
	req := &accountPB.SubaccountBalancesListRequest{
		SubaccountId: subaccount,
	}

	var header metadata.MD
	res, err := p.accountClient.SubaccountBalancesList(ctx, req, grpc.Header(&header))
	if err != nil {
		metrics.ReportFuncError(p.svcTags)
		return nil, fmt.Errorf("subaccount balance list err: %w", err)
	}

	for _, b := range res.GetBalances() {
		totalBalance, err := decimal.NewFromString(b.GetDeposit().GetTotalBalance())
		if err != nil {
			metrics.ReportFuncError(p.svcTags)
			return nil, fmt.Errorf("parse total balance err: %w", err)
		}

		availBalance, err := decimal.NewFromString(b.GetDeposit().GetAvailableBalance())
		if err != nil {
			metrics.ReportFuncError(p.svcTags)
			return nil, fmt.Errorf("parse avail balance err: %w", err)
		}

		result = append(result, &Balance{
			Denom:            b.Denom,
			TotalBalance:     totalBalance,
			AvailableBalance: availBalance,
		})
	}

	return result, nil
}

func (p *exchangeProvider) GetSpotOrders(ctx context.Context, marketIDs []string, subaccount string) (result []*SpotOrder, err error) {
	doneFn := metrics.ReportFuncTiming(p.svcTags)
	defer doneFn()
	metrics.ReportFuncCall(p.svcTags)

	req := &spotExchangePB.OrdersRequest{
		SubaccountId: subaccount,
	}

	var header metadata.MD
	res, err := p.spotExchangeClient.Orders(ctx, req, grpc.Header(&header))
	if err != nil {
		metrics.ReportFuncError(p.svcTags)
		return nil, fmt.Errorf("get spot orders err: %w", err)
	}

	// we will filter on client side, support on exchange api later
	m := make(map[string]bool)
	for _, id := range marketIDs {
		m[id] = true
	}

	isAllowAll := false
	if len(m) == 0 {
		isAllowAll = true
	}

	for _, o := range res.GetOrders() {
		_, exist := m[o.MarketId]
		if !isAllowAll && !exist {
			continue
		}

		price, err := decimal.NewFromString(o.GetPrice())
		if err != nil {
			metrics.ReportFuncError(p.svcTags)
			return nil, fmt.Errorf("parse price err: %w", err)
		}

		unfilledQuantity, err := decimal.NewFromString(o.GetUnfilledQuantity())
		if err != nil {
			metrics.ReportFuncError(p.svcTags)
			return nil, fmt.Errorf("parse unfilled quantity err: %w", err)
		}

		result = append(result, &SpotOrder{
			MarketID:     o.GetMarketId(),
			OrderHash:    o.GetOrderHash(),
			FeeRecipient: o.GetFeeRecipient(),
			OrderSide:    o.GetOrderSide(),

			Price:            price,
			UnfilledQuantity: unfilledQuantity,
		})
	}

	return result, nil
}

func (p *exchangeProvider) GetDerivativeOrders(ctx context.Context, marketIDs []string, subaccount string) (result []*DerivativeOrder, err error) {
	doneFn := metrics.ReportFuncTiming(p.svcTags)
	defer doneFn()
	metrics.ReportFuncCall(p.svcTags)

	req := &derivativeExchangePB.OrdersRequest{
		SubaccountId: subaccount,
	}

	var header metadata.MD
	res, err := p.derivativeExchangeClient.Orders(ctx, req, grpc.Header(&header))
	if err != nil {
		metrics.ReportFuncError(p.svcTags)
		return nil, fmt.Errorf("get derivative orders err: %w", err)
	}

	// we will filter on client side, support on exchange api later
	m := make(map[string]bool)
	for _, id := range marketIDs {
		m[id] = true
	}

	isAllowAll := false
	if len(m) == 0 {
		isAllowAll = true
	}

	for _, o := range res.GetOrders() {
		_, exist := m[o.MarketId]
		if !isAllowAll && !exist {
			continue
		}

		margin, err := decimal.NewFromString(o.GetMargin())
		if err != nil {
			metrics.ReportFuncError(p.svcTags)
			return nil, fmt.Errorf("parse margin err: %w", err)
		}

		result = append(result, &DerivativeOrder{
			MarketID:     o.GetMarketId(),
			OrderHash:    o.GetOrderHash(),
			FeeRecipient: o.GetFeeRecipient(),
			Margin:       margin,
		})
	}

	return result, nil
}

func (p *exchangeProvider) GetPositions(ctx context.Context, subaccount string) (result []*DerivativePosition, err error) {
	doneFn := metrics.ReportFuncTiming(p.svcTags)
	defer doneFn()
	metrics.ReportFuncCall(p.svcTags)

	req := &derivativeExchangePB.PositionsRequest{
		SubaccountId: subaccount,
	}

	var header metadata.MD
	res, err := p.derivativeExchangeClient.Positions(ctx, req, grpc.Header(&header))
	if err != nil {
		metrics.ReportFuncError(p.svcTags)
		return nil, fmt.Errorf("get position err: %w", err)
	}

	for _, pos := range res.GetPositions() {
		quantity, err := decimal.NewFromString(pos.GetQuantity())
		if err != nil {
			metrics.ReportFuncError(p.svcTags)
			return nil, fmt.Errorf("parse quantity err: %w", err)
		}

		margin, err := decimal.NewFromString(pos.GetMargin())
		if err != nil {
			metrics.ReportFuncError(p.svcTags)
			return nil, fmt.Errorf("parse quantity err: %w", err)
		}

		entryPrice, err := decimal.NewFromString(pos.GetEntryPrice())
		if err != nil {
			metrics.ReportFuncError(p.svcTags)
			return nil, fmt.Errorf("parse quantity err: %w", err)
		}

		markPrice, err := decimal.NewFromString(pos.GetMarkPrice())
		if err != nil {
			metrics.ReportFuncError(p.svcTags)
			return nil, fmt.Errorf("parse mark price err: %w", err)
		}

		result = append(result, &DerivativePosition{
			MarketID:   pos.GetMarketId(),
			Direction:  pos.GetDirection(),
			Quantity:   quantity,
			Margin:     margin,
			EntryPrice: entryPrice,
			MarkPrice:  markPrice,
		})
	}

	return result, nil
}

// GetGrants fetch first 100 grants atm, it should be engouh to check grants that bot needs
func (p *exchangeProvider) GetGrants(ctx context.Context, granter, grantee string) (*Grants, error) {
	doneFn := metrics.ReportFuncTiming(p.svcTags)
	defer doneFn()
	metrics.ReportFuncCall(p.svcTags)

	url := fmt.Sprintf(
		"%s/cosmos/authz/v1beta1/grants?granter=%s&grantee=%s&pagination.limit=100",
		p.lcdAddr, granter, grantee,
	)
	resp, err := p.httpClient.Get(url)
	if err != nil {
		metrics.ReportFuncError(p.svcTags)
		return nil, fmt.Errorf("request err: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		metrics.ReportFuncError(p.svcTags)
		return nil, fmt.Errorf("response bad status: %d", resp.StatusCode)
	}

	var res Grants
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		metrics.ReportFuncError(p.svcTags)
		return nil, fmt.Errorf("read request body err: %w", err)
	}

	if err := json.Unmarshal(bytes, &res); err != nil {
		metrics.ReportFuncError(p.svcTags)
		return nil, fmt.Errorf("marshal request body err: %w", err)
	}

	return &res, nil
}

func (p *exchangeProvider) GetPriceUSD(ctx context.Context, coinIDs []string) ([]*CoinPrice, error) {
	doneFn := metrics.ReportFuncTiming(p.svcTags)
	defer doneFn()
	metrics.ReportFuncCall(p.svcTags)

	coinList := strings.Join(coinIDs, ",")

	url := fmt.Sprintf(
		"%s/asset-price/v1/coin/price?coinIds=%s&currency=usd",
		p.assetPriceAddr, coinList,
	)
	resp, err := p.httpClient.Get(url)
	if err != nil {
		metrics.ReportFuncError(p.svcTags)
		return nil, fmt.Errorf("request err: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		metrics.ReportFuncError(p.svcTags)
		return nil, fmt.Errorf("request err: bad status: %d", resp.StatusCode)
	}

	var res CoinPriceResult
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		metrics.ReportFuncError(p.svcTags)
		return nil, fmt.Errorf("read request body err: %w", err)
	}

	if err := json.Unmarshal(bytes, &res); err != nil {
		metrics.ReportFuncError(p.svcTags)
		return nil, fmt.Errorf("marshal request body err: %w", err)
	}
	return res.Data, nil
}

func (p *exchangeProvider) GetBankBalance(ctx context.Context, address string) (*BankAccountBalances, error) {
	doneFn := metrics.ReportFuncTiming(p.svcTags)
	defer doneFn()
	metrics.ReportFuncCall(p.svcTags)

	url := fmt.Sprintf(
		"%s/cosmos/bank/v1beta1/balances/%s?pagination.limit=1000",
		p.lcdAddr, address,
	)
	resp, err := p.httpClient.Get(url)
	if err != nil {
		metrics.ReportFuncError(p.svcTags)
		return nil, fmt.Errorf("request err: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		metrics.ReportFuncError(p.svcTags)
		return nil, fmt.Errorf("request err: bad status: %d", resp.StatusCode)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		metrics.ReportFuncError(p.svcTags)
		return nil, fmt.Errorf("read request body err: %w", err)
	}

	var res BankAccountBalances
	if err := json.Unmarshal(bytes, &res); err != nil {
		metrics.ReportFuncError(p.svcTags)
		return nil, fmt.Errorf("failed marshal bankaccount resp: %w", err)
	}

	return &res, nil
}

func (p *exchangeProvider) GetExchangeConn() *grpc.ClientConn {
	return p.conn
}

func (p *exchangeProvider) Close() error {
	return p.conn.Close()
}

func DialerFunc(ctx context.Context, addr string) (net.Conn, error) {
	return Connect(addr)
}

// Connect dials the given address and returns a net.Conn. The protoAddr argument should be prefixed with the protocol,
// eg. "tcp://127.0.0.1:8080" or "unix:///tmp/test.sock"
func Connect(protoAddr string) (net.Conn, error) {
	proto, address := ProtocolAndAddress(protoAddr)
	conn, err := net.Dial(proto, address)
	return conn, err
}

// ProtocolAndAddress splits an address into the protocol and address components.
// For instance, "tcp://127.0.0.1:8080" will be split into "tcp" and "127.0.0.1:8080".
// If the address has no protocol prefix, the default is "tcp".
func ProtocolAndAddress(listenAddr string) (string, string) {
	protocol, address := "tcp", listenAddr
	parts := strings.SplitN(address, "://", 2)
	if len(parts) == 2 {
		protocol, address = parts[0], parts[1]
	}
	return protocol, address
}
