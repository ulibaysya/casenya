package valkey

import (
	"context"
	"crypto/tls"
	"net"
	"os"
	"testing"
	"time"

	"github.com/valkey-io/valkey-go"
)

func newValkey(t *testing.T) valkey.Client {
	ctx, cancel := context.WithTimeout(t.Context(), time.Second*5)
	defer cancel()

	connstr, ok := os.LookupEnv("TEST_VALKEY_CONNSTR")
	if !ok {
		t.Fatalf("TEST_VALKEY_CONNSTR is not set for integration test")
	}

	clientOption := valkey.ClientOption{
		InitAddress: []string{connstr},
	}

	_, unix := os.LookupEnv("TEST_VALKEY_UNIX")
	if unix {
		clientOption.DialCtxFn = func(ctx context.Context, s string, d *net.Dialer, c *tls.Config) (conn net.Conn, err error) {
			return d.DialContext(ctx, "unix", s)
		}
	}

	client, err := valkey.NewClient(clientOption)
	if err != nil {
		t.Fatalf("error while initializing connection with valkey: %v", err)
	}

	if err = client.Do(ctx, client.B().Ping().Build()).Error(); err != nil {
		t.Fatalf("error while pinging valkey: %v", err)
	}

	return client
}

func resetValkey(t *testing.T, client valkey.Client) {
	ctx, cancel := context.WithTimeout(t.Context(), time.Second*5)
	defer cancel()

	if err := client.Do(ctx, client.B().Flushall().Build()).Error(); err != nil {
		t.Fatalf("error while flushing all: %v", err)
	}
}
