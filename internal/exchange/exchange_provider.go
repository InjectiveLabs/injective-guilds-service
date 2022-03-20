package exchange

import (
	"context"
	"fmt"
	"net"
	"strings"

	accountPB "github.com/InjectiveLabs/sdk-go/exchange/accounts_rpc/pb"
	derivativeExchangePB "github.com/InjectiveLabs/sdk-go/exchange/derivative_exchange_rpc/pb"
	spotExchangePB "github.com/InjectiveLabs/sdk-go/exchange/spot_exchange_rpc/pb"
	"github.com/pkg/errors"
	log "github.com/xlab/suplog"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

// TODO: improve sdk-go with new functions
// NOW:
// - Decided to use a stable release now, v1.31.0
// - Atm, we can clone functions into this repo first
type exchangeProvider struct {
	DataProvider

	conn                     *grpc.ClientConn
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

// derives from current `master` of sdk-go
// vvvvv

func NewExchangeProvider(protoAddr string, options ...ClientOption) (DataProvider, error) {
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
		conn, err = grpc.Dial(protoAddr, grpc.WithTransportCredentials(opts.TLSCert), grpc.WithContextDialer(DialerFunc))
	} else {
		conn, err = grpc.Dial(protoAddr, grpc.WithInsecure(), grpc.WithContextDialer(DialerFunc))
	}

	if err != nil {
		err := errors.Wrapf(err, "failed to connect to the gRPC: %s", protoAddr)
		return nil, err
	}

	cc := &exchangeProvider{
		conn:                     conn,
		accountClient:            accountPB.NewInjectiveAccountsRPCClient(conn),
		spotExchangeClient:       spotExchangePB.NewInjectiveSpotExchangeRPCClient(conn),
		derivativeExchangeClient: derivativeExchangePB.NewInjectiveDerivativeExchangeRPCClient(conn),
	}

	return cc, nil
}

func (p *exchangeProvider) GetDefaultSubaccountBalances(ctx context.Context, subaccount string) (result []*Balance, err error) {
	// get all denoms
	req := &accountPB.SubaccountBalancesListRequest{
		SubaccountId: subaccount,
	}

	var header metadata.MD
	res, err := p.accountClient.SubaccountBalancesList(ctx, req, grpc.Header(&header))
	if err != nil {
		return nil, err
	}

	for _, b := range res.Balances {
		totalBalance, err := primitive.ParseDecimal128(b.Deposit.TotalBalance)
		if err != nil {
			return nil, fmt.Errorf("parse total balance err: %w", err)
		}

		availBalance, err := primitive.ParseDecimal128(b.Deposit.AvailableBalance)
		if err != nil {
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
