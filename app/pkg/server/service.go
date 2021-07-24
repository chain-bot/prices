package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/chain-bot/prices/app/configs"
	"go.uber.org/fx"
)

func Run(
	lc fx.Lifecycle,
	routes *Routes,
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
	routes *Routes,
	secrets *configs.Secrets,
) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(routes.ping))
	mux.Handle("/candles", http.HandlerFunc(routes.getCandles))
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", secrets.ServerConfig.Port),
		Handler: mux,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen, err=%+s\n", err)
		}
	}()
	log.Printf(fmt.Sprintf("server started on port=%d", secrets.ServerConfig.Port))
	return srv
}

func shutdown(
	httpSrv *http.Server,
) error {
	log.Printf("server stopped")
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	if err := httpSrv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server Shutdown Failed, err=%+s", err)
		return err
	}

	log.Printf("server exited properly")
	return nil
}
