package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	amzn "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kalpaj/verve/internal/config"
	"github.com/kalpaj/verve/internal/http"
	"github.com/kalpaj/verve/internal/http/router"
	"github.com/kalpaj/verve/pkg/aws"
	"github.com/kalpaj/verve/pkg/db/redis"
	"github.com/kalpaj/verve/pkg/job"
)

func main() {
	log.SetFlags(log.LstdFlags | log.LUTC | log.Llongfile)

	// Init application configuration
	cfg := config.New()

	// Init redis
	store := must(func() (*redis.Redis, error) {
		return redis.New(cfg.RedisEndpoint, cfg.RedisPassword, cfg.RedisMaxRetries, cfg.RedisMinIdleConnections)
	}, "redis")

	// Init server's routes
	router := must(func() (*router.Router, error) {
		return router.New(cfg, store)
	}, "router")

	// Init server's aws services
	awsConfig := must(func() (*amzn.Config, error) {
		return aws.ConfigWithSecretKey(cfg.AwsAccessKey, cfg.AwsSecretKey, cfg.AwsRegion)
	}, "aws configuration")

	// Init kinesis
	kinesis := aws.NewKinesis(awsConfig)

	// Unique cout publihser job
	ucJob := job.New(store, kinesis, cfg.AwsKinesisUniqueCountStream)

	// Init server
	server := must(func() (*http.Server, error) {
		return http.NewServer(router.GetHandler(), cfg.HTTPHost, cfg.HTTPPort, cfg.HTTPTimeout)
	}, "server")

	// Run server
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		ucJob.Start()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := server.Start(); err != nil {
			log.Fatalf("ApplicationStartupError: Server failed to start - %v\n", err)
		}
	}()
	log.Printf("Server successfully started on port %v\n", cfg.HTTPPort)

	// Wait stop signal to shut down server gracefully
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	<-sigchan
	log.Printf("Shutting down the server\n")

	server.Stop()
	wg.Wait()
}

func must[T any](i func() (T, error), name string) T {
	mod, err := i()

	if err != nil {
		log.Fatalf("ApplicationError: Failed to init %s: %s", name, err.Error())
	}

	return mod
}
