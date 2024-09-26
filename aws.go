package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

// Function to read profiles from the AWS config file
func getAWSProfiles() ([]string, error) {
	awsConfigPath := os.ExpandEnv("$HOME/.aws/config")

	// Open the file
	file, err := os.Open(awsConfigPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a regular expression to match profile names
	profileRegex := regexp.MustCompile(`\[profile (.*?)\]`)

	// Create a slice to hold the profiles
	var profiles []string

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if matches := profileRegex.FindStringSubmatch(line); matches != nil {
			profiles = append(profiles, matches[1])
		}
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return profiles, nil
}

// Function to execute the AWS login command
func executeAWSLogin(command string, profile string) error {
	// Replace PROFILE_NAME with the actual profile in the command
	cmdString := regexp.MustCompile(`PROFILE_NAME`).ReplaceAllString(command, profile)

	// Construct the command to execute
	cmd := exec.Command("sh", "-c", cmdString)

	// Set the command's output to the terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	return cmd.Run()
}

// Function to execute the EKS login command
func executeEKSLogin(clusterName string, profile string) error {
	// Construct the command string
	cmdString := fmt.Sprintf("aws eks update-kubeconfig --region eu-west-1 --name %s --profile %s", clusterName, profile)

	// Construct the command to execute
	cmd := exec.Command("sh", "-c", cmdString)

	// Set the command's output to the terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	return cmd.Run()
}
