package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/HarukiIdo/go_sample/handler"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler.NewHandler(),
	}

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5+time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()
	fmt.Printf("start receiving at :8080")
	fmt.Fprintln(os.Stderr, server.ListenAndServe())
}
