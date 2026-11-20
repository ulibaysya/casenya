package main

import (
	"context"
	cryptorand "crypto/rand"
	"crypto/tls"
	"math/rand/v2"
	"net"
	"sync"

	phuslog "github.com/phuslu/log"
	"github.com/valkey-io/valkey-go"
)

func main() {
	c, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{`/run/valkey/valkey.sock`},
		DialCtxFn: func(ctx context.Context, s string, d *net.Dialer, c *tls.Config) (conn net.Conn, err error) {
			return d.DialUnix(ctx, "unix", nil, &net.UnixAddr{
				Name: s,
			})
		},
	})
	if err != nil {
		phuslog.Fatal().Err(err).Msg("")
	}
	defer c.Close()

	if err = c.Do(context.Background(), c.B().Ping().Build()).Error(); err != nil {
		phuslog.Fatal().Err(err).Msg("")
	}

	for range 500 {
		if err = c.Do(context.Background(), c.B().Set().Key(cryptorand.Text()).Value("anya").Build()).Error(); err != nil {
			phuslog.Fatal().Err(err).Msg("")
		}
	}

	if err = c.Do(context.Background(), c.B().Set().Key("anya").Value("i love you").Build()).Error(); err != nil {
		phuslog.Fatal().Err(err).Msg("")
	}

	value, err := c.Do(context.Background(), c.B().Get().Key("anya").Build()).ToString()
	if err != nil {
		phuslog.Fatal().Err(err).Msg("")
	}

	phuslog.Info().Str("idi nahuy", value).Msg("")

	wg := sync.WaitGroup{}

	for i := range 1000 {
		wg.Go(func() { scan(i, c) })
	}
	wg.Wait()

	if err = c.Do(context.Background(), c.B().Flushall().Async().Build()).Error(); err != nil {
		phuslog.Fatal().Err(err).Msg("")
	}

}

func scan(scanIndex int, c valkey.Client) {
	unique := map[string]struct{}{}
	scanCursor := uint64(0)
	for i := 1; ; i += 1 {
		count := rand.Int64N(150) + 1
		scan, err := c.Do(context.TODO(), c.B().Scan().Cursor(scanCursor).Count(count).Build()).AsScanEntry()
		if err != nil {
			phuslog.Fatal().Err(err).Int("i", scanIndex).Uint64("cursor", scanCursor).Msg("")
		}
		phuslog.Info().Strs("entries", scan.Elements).Int("i", scanIndex).Int("entriesLen", len(scan.Elements)).Int64("count", count).Uint64("cursor", scan.Cursor).Msg("")
		for _, elem := range scan.Elements {
			if _, found := unique[elem]; found {
				phuslog.Error().Str("elem", elem).Int("i", scanIndex).Msg("found already returned element")
			}
			unique[elem] = struct{}{}
		}
		if scan.Cursor == 0 {
			phuslog.Info().Int("iter", i).Int("i", scanIndex).Int("uniqueLen", len(unique)).Msg("scan end")
			break
		}
		scanCursor = scan.Cursor
	}
}
