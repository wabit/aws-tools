package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Setup signal handling for graceful exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	// Load or create configuration
	config := loadConfig()

	// Run the main application loop
	go func() {
		<-c        // Wait for signal
		os.Exit(0) // Exit when an interrupt signal is received
	}()

	// Start the main menu
	for {
		mainMenu(config)
	}
}
