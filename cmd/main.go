package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	
	"github.com/usa4ev/gotracer/internal/mongo/connector"
	"github.com/usa4ev/gotracer/internal/mongo/mongohook"
	"github.com/usa4ev/gotracer/internal/mongo/provider"
	"github.com/usa4ev/gotracer/internal/resources"
	"github.com/usa4ev/gotracer/internal/router"
	"github.com/usa4ev/gotracer/internal/tracer"
)

var (
	// command-line options:
	httpServerEndpoint = flag.String("http-server-endpoint", ":8080", "HTTP server endpoint")
	mongouri = flag.String("mongobd_uri", "mongodb://localhost:27017", "Mongo db URI")
	resourcePath = flag.String("resource_path", "./config/resources", "Path to url list")
)

func main(){

	// Read url list from file 
	resources,err := resources.ReadResources(*resourcePath)
	if err != nil{
		log.Fatal(err)
	}

	tracker := tracer.New(resources)
	
	// Create logger
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{
		DisableTimestamp: true,
	}

	// Connect to mongo
	conn,err := connector.New(*mongouri)
	if err != nil{
		log.Fatal(err)
	}

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
			log.Printf("HTTP server Shutdown: %v", err)
		}

		log.Printf("graceful shutdown, got call: %v\n", call.String())
	}()

	if err := srv.ListenAndServe(); err != nil{
		log.Fatal(err)
	}
}