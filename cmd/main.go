package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/usa4ev/gotracer/internal/mongo/connector"
	"github.com/usa4ev/gotracer/internal/mongo/mongohook"
	"github.com/usa4ev/gotracer/internal/mongo/provider"
	"github.com/usa4ev/gotracer/internal/resources"
	"github.com/usa4ev/gotracer/internal/router"
	"github.com/usa4ev/gotracer/internal/tracer"
)

func main(){
	// command-line options:
	httpServerEndpoint := flag.String("http-server-endpoint", ":8080", "HTTP server endpoint")
	mongouri := flag.String("mongobd-uri", "mongodb://mongo:27017", "Mongo db URI")
	resourcePath := flag.String("resource-path", "./config/resources", "Path to url list")

	flag.Parse()

	// Create logger
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{
		DisableTimestamp: true,
	}

	// Read url list from file 
	resources,err := resources.ReadResources(*resourcePath)
	if err != nil{
		logger.Fatal(err)
	}

	tracker := tracer.New(resources)

	// Connect to mongo
	conn,err := connector.New(*mongouri)
	if err != nil{
		logger.Fatal(err)
	}

	defer func(){
		ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
		err := conn.Disconnect(ctx)
		if err != nil{
			logger.Errorf("failed to close mongo connection: %v", ctx)
		}
	}()

	// Create hook that writes logs to mongo
	hook := mongohook.New(conn)
	logger.Hooks.Add(hook)

	// Provider queries mongo
	prov := provider.New(conn)	

	// Create server
	r := router.New(tracker, logger, prov)
	srv := &http.Server{Addr: *httpServerEndpoint, Handler: r}

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		call := <-sig

		if err := srv.Shutdown(context.Background()); err != nil {
			// failed to close listener
			logger.Infof("HTTP server Shutdown: %v", err)
		}

		logger.Errorf("graceful shutdown, got call: %v\n", call.String())
	}()

	if err := srv.ListenAndServe(); err != nil{
		logger.Fatal(err)
	}
}