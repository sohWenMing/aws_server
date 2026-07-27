package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/sohWenMing/aws_server/internal/awsserver"
	dbpostgres "github.com/sohWenMing/aws_server/internal/db_postgres"
	"github.com/sohWenMing/aws_server/internal/envloader"
)

const addr = ":8080"

func main() {
	project_root, err := get_project_root()
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
	conn, err := pool.Acquire(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	connErr := conn.Ping(context.TODO())
	if connErr != nil {
		log.Fatal(err)
	}
	if err != nil {
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

func get_project_root() (project_root string, err error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		_, err = os.Stat(fmt.Sprintf("%s/go.mod", cwd))
		// if we can find the go.mod file, this is the project root, exit
		if err == nil {
			return cwd, nil
		}
		if cwd == "/" {
			log.Fatal("project root could not be found")
		}
		if errors.Is(err, os.ErrNotExist) {
			cwd = filepath.Dir(cwd)
			continue
		} else {
			log.Fatal(fmt.Sprintf("project root could not be found. expected error: %v", err))
		}
	}
}
