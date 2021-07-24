package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/chain-bot/prices/app/configs"
	"github.com/chain-bot/prices/app/pkg/server/routes"
	"go.uber.org/fx"
)

func Run(
	lc fx.Lifecycle,
	routes *routes.Handler,
	secrets *configs.Secrets,
) {
	var httpSrv http.Server
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			httpSrv = *listenAndServe(routes, secrets)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			shutdown(&httpSrv)
			return nil
		},
	})
}

func listenAndServe(
	routes *routes.Handler,
	secrets *configs.Secrets,
) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/", logMiddleware(http.HandlerFunc(routes.Ping)))
	mux.Handle("/candles", logMiddleware(http.HandlerFunc(routes.GetCandles)))
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", secrets.ServerConfig.Port),
		Handler: mux,
	}
	log.WithField("port", secrets.ServerConfig.Port).Infof("starting server")
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithField("err", err.Error()).Fatalf("server listening")
		}
	}()
	return srv
}

func shutdown(
	httpSrv *http.Server,
) error {
	log.Infof("stopping server")
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	if err := httpSrv.Shutdown(ctxShutDown); err != nil {
		log.WithField("err", err.Error()).Fatalf("failed to shutdown server")
		return err
	}

	log.Infof("server exited")
	return nil
}

func logMiddleware(
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithField("params", r.URL.Query()).Infof("handling %s", r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
