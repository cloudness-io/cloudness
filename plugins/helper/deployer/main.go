package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log := NewLogger(os.Getenv("VERBOSE") == "true")

	// Load configuration
	cfg, err := LoadConfigFromEnv()
	if err != nil {
		log.Error("Configuration error: %v", err)
		os.Exit(1)
	}

	// Create deployer
	deployer, err := NewDeployer(cfg)
	if err != nil {
		log.Error("Failed to initialize deployer: %v", err)
		os.Exit(1)
	}

	// Setup context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Warn("Received shutdown signal")
		cancel()
	}()

	// Run deployment
	if err := deployer.Deploy(ctx); err != nil {
		log.Error("Deployment failed: %v", err)
		os.Exit(1)
	}
}
