package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Version information (set by build flags)
var (
	Version   = "dev"
	BuildTime = "unknown"
)

func main() {
	var (
		configPath  = flag.String("config", "config.yaml", "Path to configuration file")
		dryRun      = flag.Bool("dry-run", false, "Perform a dry run without making changes")
		showVersion = flag.Bool("version", false, "Show version information")
		showHelp    = flag.Bool("help", false, "Show help information")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("gh-deployer version %s\nBuilt: %s\n", Version, BuildTime)
		os.Exit(0)
	}

	if *showHelp {
		fmt.Println("GitHub Release Deployer - Blue/Green deployment tool")
		fmt.Println("")
		fmt.Println("Usage:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	// Load configuration
	config, err := LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Setup logging
	logger := setupLogging(config)

	// Create deployer
	deployer, err := NewDeployer(config, logger, *dryRun)
	if err != nil {
		logger.Fatalf("Failed to create deployer: %v", err)
	}

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChan
		logger.Println("Received shutdown signal")
		cancel()
	}()

	// Start the deployer
	logger.Println("Starting GitHub Release Deployer")
	if err := deployer.Run(ctx); err != nil {
		logger.Fatalf("Deployer failed: %v", err)
	}

	logger.Println("Deployer stopped gracefully")
}

func setupLogging(config *Config) *log.Logger {
	// Set log output
	if config.Logging.File != "" {
		file, err := os.OpenFile(config.Logging.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
		if err == nil {
			return log.New(file, "", log.LstdFlags|log.Lshortfile)
		} else {
			log.Printf("Failed to open log file %s: %v", config.Logging.File, err)
		}
	}

	return log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
}
