package awsserver

import (
	"context"
	"net/http"
	"os"

	"github.com/sohWenMing/aws_server/internal/handlers"
)

func InitServer(addr string) *http.Server {
	server := &http.Server{}
	server.Addr = addr
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handlers.HealthHandler)
	server.Handler = mux
	return server
}

func ShutdownOnSignal(server *http.Server, sigChan <-chan os.Signal, ctx context.Context) {
	<-sigChan
	server.Shutdown(ctx)
	return
}
