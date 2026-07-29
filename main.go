package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sohWenMing/aws_server/internal/awsserver"
	dbpostgres "github.com/sohWenMing/aws_server/internal/db_postgres"
	"github.com/sohWenMing/aws_server/internal/envloader"
	"github.com/sohWenMing/aws_server/internal/utils"
)

const addr = ":8080"

func main() {
	project_root, err := utils.GetProjectRoot()
	if err != nil {
		log.Fatal(err)
	}
	envSettings, err := envloader.LoadEnvironment(fmt.Sprintf("%s/.env", project_root))
	if err != nil {
		log.Fatal(err)
	}
	pool, err := dbpostgres.ConnectToDB(envSettings, context.TODO())
	if err != nil {
		if envSettings.GetEnv() == "DEV" {
			log.Fatal("db connection failed in dev mode, is docker compose launched? Err: %v\n", err)
		} else {
			log.Fatal(err)
		}
	}
	if err := checkDBConnection(pool); err != nil {
		log.Fatal(err)
	}
	// fmt.Println("pool: ", pool)
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

// sets a 10 second timeout for the whole operation to be able to finish, where we can get the correct ping
func checkDBConnection(pool *pgxpool.Pool) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	doneChan := ctx.Done()
	defer cancelFunc()
	conn, err := pool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return err
	}

	connErr := conn.Ping(ctx)
	if connErr != nil {
		return err
	}
	for {
		select {
		case <-doneChan:
			return ctx.Err()
		default:
			if err != nil {
				fmt.Println("ping test to db connection not ok ... retrying ...")
				continue
			} else {
				fmt.Println("ping test to db connection is OK")
				return nil
			}
		}
	}
}
