package main

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	phuslog "github.com/phuslu/log"
	ustrg "github.com/ulibaysya/casenya/internal/modules/user/userstorage"
	"github.com/valkey-io/valkey-go"
)

func main() {
	valkeyClient, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{"/run/valkey/valkey.sock"},
		DialCtxFn: func(ctx context.Context, s string, d *net.Dialer, c *tls.Config) (conn net.Conn, err error) {
			return d.DialContext(ctx, "unix", s)
		},
	})
	if err != nil {
		phuslog.Fatal().Err(err).Msg("")
	}

	userStorage := ustrg.New(valkeyClient, time.Now)
	_ = userStorage
}
