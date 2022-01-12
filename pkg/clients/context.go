package clients

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	analyticsv1 "github.com/speechly/api/go/speechly/analytics/v1"
	configv1 "github.com/speechly/api/go/speechly/config/v1"
	salv1 "github.com/speechly/api/go/speechly/sal/v1"
	wluv1 "github.com/speechly/api/go/speechly/slu/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type contextKey int

const (
	keySpeechlyConfig contextKey = iota
	keyClientConnection
	keyFailFunc
)

var (
	ErrNoConfig        = errors.New("no Speechly project settings found")
	ErrNoContext       = errors.New("no project specified")
	ErrContextNotFound = errors.New("selected project not found in Speechly project settings file")

	ConnectionTimeout = 4 * time.Second
)

type connectionCache struct {
	sc   *SpeechlyContext
	conn *grpc.ClientConn
	ff   FailFunc
}

type FailFunc func(error)

func (cc *connectionCache) getConnection(ctx context.Context) *grpc.ClientConn {
	if cc.conn != nil {
		return cc.conn
	}
	if cc.sc.Host == "" {
		cc.ff(errors.New("no API host defined"))
		return nil
	}
	serverAddr := cc.sc.Host
	opts := []grpc.DialOption{grpc.WithBlock()}
	if strings.Contains(cc.sc.Host, "speechly.com") {
		// Always use TLS for Speechly hosts
		serverAddr = serverAddr + ":443"
		creds := credentials.NewTLS(&tls.Config{
			ServerName: cc.sc.Host,
		})
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	connCtx, cancel := context.WithTimeout(ctx, ConnectionTimeout)
	defer cancel()
	conn, err := grpc.DialContext(connCtx, serverAddr, opts...)
	if err != nil {
		cc.ff(fmt.Errorf("Connecting to host %s failed: %v", cc.sc.Host, err))
		return nil
	}
	cc.conn = conn
	return cc.conn
}

func NewContext(ff FailFunc) context.Context {
	ctx := context.Background()
	config, err := getSpeechlyConfig()
	if err != nil {
		config = &Config{}
	}
	ctx = context.WithValue(ctx, keySpeechlyConfig, config)
	ctx = context.WithValue(ctx, keyFailFunc, ff)

	sc := config.GetSpeechlyContext()
	if sc != nil {
		if ff == nil {
			ff = func(err error) {
				log.Fatalf("error: %v", err)
			}
		}
		ctx = context.WithValue(ctx, keyClientConnection, &connectionCache{sc: sc, ff: ff})
		md := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", sc.Apikey))
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	return ctx
}

func GetConfig(ctx context.Context) *Config {
	config, ok := ctx.Value(keySpeechlyConfig).(*Config)
	if !ok {
		return &Config{}
	}
	return config
}

func ConfigClient(ctx context.Context) (configv1.ConfigAPIClient, error) {
	cc, ok := ctx.Value(keyClientConnection).(*connectionCache)
	if !ok {
		return nil, errors.New("invalid context")
	}
	return configv1.NewConfigAPIClient(cc.getConnection(ctx)), nil
}

func AnalyticsClient(ctx context.Context) (analyticsv1.AnalyticsAPIClient, error) {
	cc, ok := ctx.Value(keyClientConnection).(*connectionCache)
	if !ok {
		return nil, errors.New("invalid context")
	}
	return analyticsv1.NewAnalyticsAPIClient(cc.getConnection(ctx)), nil
}

func CompileClient(ctx context.Context) (salv1.CompilerClient, error) {
	cc, ok := ctx.Value(keyClientConnection).(*connectionCache)
	if !ok {
		return nil, errors.New("invalid context")
	}
	return salv1.NewCompilerClient(cc.getConnection(ctx)), nil
}

func WLUClient(ctx context.Context) (wluv1.WLUClient, error) {
	cc, ok := ctx.Value(keyClientConnection).(*connectionCache)
	if !ok {
		return nil, errors.New("invalid context")
	}
	return wluv1.NewWLUClient(cc.getConnection(ctx)), nil
}
