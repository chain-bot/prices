package main

import (
	"context"
	"os"

	"github.com/chain-bot/prices/app/configs"
	"github.com/chain-bot/prices/app/internal/data/influxdb"
	"github.com/chain-bot/prices/app/internal/data/psql"
	"github.com/chain-bot/prices/app/internal/repository"
	"github.com/chain-bot/prices/app/pkg/server"
	"github.com/chain-bot/prices/app/pkg/server/routes"
	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	// TODO(Zahin): Enable if deploying to allow integration with splunk/logstash
	// log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func main() {
	fxApp := fx.New(
		fx.Provide(configs.GetSecrets),
		fx.Provide(psql.NewDatabase),
		fx.Provide(influxdb.NewInfluxDBClient),
		fx.Provide(repository.NewRepository),
		fx.Provide(routes.NewHandler),
		fx.Invoke(
			psql.RunMigrations,
			server.Run,
		),
		fx.NopLogger,
	)
	if err := fxApp.Start(context.Background()); err != nil {
		log.WithField("err", err.Error()).Fatalf("starting fx app for api server")
	}
	<-fxApp.Done()
}
