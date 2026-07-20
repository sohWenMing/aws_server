package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sohWenMing/aws_server/internal/awsserver"
)

const addr = ":8080"

func main() {
	sigChan := make(chan (os.Signal), 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	srv := awsserver.InitServer(addr)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	go awsserver.ShutdownOnSignal(srv, sigChan, ctx)
	fmt.Println("starting server and listening on addr", addr)
	if err := srv.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Println("server shutdown gracefully")
			os.Exit(0)
			return
		} else {
			fmt.Println("error when shutting down: ", err)
			os.Exit(1)
			return
		}
	}
}
