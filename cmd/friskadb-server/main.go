package main

import (
	"flag"
	"friskadb/pkg/server"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func main() {
	// Parse flags
	host := flag.String("host", "0.0.0.0", "Server host")
	port := flag.Int("port", 7171, "Server port")
	dbName := flag.String("db", "mydb", "Database name")
	dataDir := flag.String("dir", getDefaultDataDir(), "Data directory")
	adminUser := flag.String("user", "admin", "Admin username")
	adminPass := flag.String("pass", "friska123", "Admin password")
	flag.Parse()

	// Create server config
	cfg := &server.Config{
		Host:          *host,
		Port:          *port,
		DatabaseName:  *dbName,
		DataDir:       *dataDir,
		AdminUser:     *adminUser,
		AdminPassword: *adminPass,
	}

	// Create server
	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to create server: %v", err)
	}

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("❌ Server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("\n💾 Shutting down server...")

	// Stop server gracefully
	if err := srv.Stop(); err != nil {
		log.Printf("❌ Error during shutdown: %v", err)
	}
}

func getDefaultDataDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./friskadb_data"
	}
	return filepath.Join(homeDir, ".friskadb")
}
