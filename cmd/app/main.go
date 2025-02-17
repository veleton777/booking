package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/veleton777/booking_api/internal/config"
	"github.com/veleton777/booking_api/internal/server"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		log.Fatalln("load config", err)
	}

	time.Local, err = time.LoadLocation("UTC")
	if err != nil {
		log.Fatalln("load location", err)
	}

	l := zerolog.New(os.Stdout).
		Level(conf.LogLevel()).
		With().Timestamp().Stack().Caller().
		Logger()

	ctx := l.WithContext(context.Background())
	exitCode := 0

	s, err := server.New(ctx, &conf, &l)
	if err != nil {
		l.Fatal().Err(err).Msg("create server")

		exitCode = 1
		os.Exit(exitCode)
	}

	if err = s.Run(ctx); err != nil {
		err = errors.Wrap(err, "run server")
		l.Err(err).Send()

		exitCode = 1
	}

	os.Exit(exitCode)
}
