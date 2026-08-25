package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"minebelt/internal/console"
	"minebelt/internal/store"
)

func main() {
	cfg := LoadConfig()
	st, err := store.NewStore(cfg.DataDir)
	if err != nil {
		log.Fatalf("init store: %v", err)
	}
	log.Printf("data directory: %s", st.Dir())
	app := console.Build(st)
	addr := ":" + cfg.Port
	srv := &http.Server{Addr: addr, Handler: app.Router()}
	go func() {
		log.Printf("minebelt console listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("serve: %v", err)
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
