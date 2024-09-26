package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2" // Import survey package
)

const configFileName = ".aws-tools"

type Config struct {
	LoginCommand string            `json:"login_command"`
	EKSClusters  []EKSCluster      `json:"eks_clusters"`
	OtherOptions map[string]string `json:"other_options"` // Add any additional options as needed
}

// Function to load or create the configuration in JSON format
func loadConfig() Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		os.Exit(1)
	}

	configPath := filepath.Join(homeDir, configFileName)

	// Check if the config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		var config Config
		fmt.Println("No configuration file found.")
		prompt := &survey.Input{
			Message: "Please enter the AWS login command to use:",
			Default: "aws-google-login -p PROFILE_NAME", // Provide a default command
		}
		err := survey.AskOne(prompt, &config.LoginCommand)
		if err != nil {
			fmt.Println("Error during prompt:", err)
			os.Exit(1)
		}

		// Initialize other options if needed
		config.OtherOptions = make(map[string]string)

		// Save the config to the JSON file
		saveConfig(config)

		fmt.Printf("Configuration file created at %s\n", configPath)
		return config
	}

	// Read the config from the JSON file
	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Println("Error reading configuration file:", err)
		os.Exit(1)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Error unmarshaling JSON config:", err)
		os.Exit(1)
	}

	return config
}

// Function to save the updated configuration to the JSON file
func saveConfig(config Config) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}

	configPath := filepath.Join(homeDir, configFileName)

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling config to JSON:", err)
		return
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		fmt.Println("Error writing configuration file:", err)
	}
}
