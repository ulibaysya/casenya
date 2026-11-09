package main

import (
	"context"
	"log/slog"
	"math/rand/v2"
	"net"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgx/v5/pgxpool"
	phuslog "github.com/phuslu/log"
	itemblobstrg "github.com/ulibaysya/casenya/internal/modules/item/blobstrg"
	itemh "github.com/ulibaysya/casenya/internal/modules/item/handler"
	"github.com/ulibaysya/casenya/internal/modules/item/service"
	itemstrg "github.com/ulibaysya/casenya/internal/modules/item/storage"
	itemtagh "github.com/ulibaysya/casenya/internal/modules/itemtag/handler"
	itemtagsrvc "github.com/ulibaysya/casenya/internal/modules/itemtag/service"
	itemtagstrg "github.com/ulibaysya/casenya/internal/modules/itemtag/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	logger := slog.New(phuslog.SlogNewJSONHandler(os.Stderr, &slog.HandlerOptions{
		// Level: slog.LevelDebug,
		Level: slog.LevelInfo,
	}))

	db, err := pgxpool.New(context.Background(), "dbname=casenya")
	if err != nil {
		phuslog.Fatal().Err(err).Msg("")
	}

	s3Cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		phuslog.Fatal().Err(err).Msg("")
	}
	presignClient := s3.NewPresignClient(
		s3.NewFromConfig(s3Cfg, func(o *s3.Options) {
			o.UsePathStyle = true
		}),
		func(po *s3.PresignOptions) {
			po.Expires = time.Second * 60
		},
	)

	rng := rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))

	serv := grpc.NewServer()

	itemtagStorage := itemtagstrg.New(db, logger)
	itemtagService := itemtagsrvc.New(itemtagStorage)
	itemtagh.Register(serv, itemtagService)

	itemBlobStorage := itemblobstrg.New("127.0.0.1:9000", "items", presignClient, nil)
	itemStorage := itemstrg.New(db, logger)
	itemService := service.New(itemStorage, itemBlobStorage, rng)
	itemh.Register(serv, itemService)

	reflection.Register(serv)

	l, err := net.Listen("tcp", ":50051")
	if err != nil {
		phuslog.Fatal().Err(err).Msg("")
	}
	defer l.Close()
	phuslog.Info().Msg("listening")
	phuslog.Fatal().Err(serv.Serve(l)).Msg("")
}
