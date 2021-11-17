package server

import (
	"avito-intership/balance"
	balanceHttp "avito-intership/balance/delivery/http"
	"avito-intership/balance/repository/postgres"
	"avito-intership/balance/usecase"
	"avito-intership/db"
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type App struct {
	httpServer *http.Server

	balance balance.UseCase
}

func NewApp() *App {
	balanceRepo := postgres.NewBalanceRepository(db.GetDB())

	return &App{
		balance: usecase.NewBalanceUseCase(balanceRepo),
	}
}

func (a *App) Run(port string) error {
	router := mux.NewRouter()

	balanceHttp.RegisterEndpoints(router, a.balance)

	router.Use(mux.CORSMethodMiddleware(router))
	a.httpServer = &http.Server{
		Addr:           port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}