package main

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

type EKSCluster struct {
	ClusterName string
	Profile     string
}

// Main menu function
func mainMenu(config Config) {
	// Main menu options
	options := []string{"AWS Login", "EKS Login", "Configuration", "Exit"}
	var selectedOption string

	prompt := &survey.Select{
		Message: "Main Menu:",
		Options: options,
	}

	// Run the main menu prompt
	err := survey.AskOne(prompt, &selectedOption)
	if err != nil {
		if err.Error() == "interrupt" {
			return // If interrupted, just return
		}
		fmt.Println("Error during prompt:", err)
		return
	}

	switch selectedOption {
	case "AWS Login":
		awsLoginMenu(config)
	case "EKS Login":
		eksLoginMenu(config)
	case "Configuration":
		configMenu(&config)
	case "Exit":
		fmt.Println("Exiting application...")
		os.Exit(0)
	}
}

// EKS login menu function
func eksLoginMenu(config Config) {
	if len(config.EKSClusters) == 0 {
		fmt.Println("No EKS clusters configured.")
		return
	}

	clusterNames := make([]string, len(config.EKSClusters))
	for i, cluster := range config.EKSClusters {
		clusterNames[i] = cluster.ClusterName
	}

	var selectedClusterName string
	clusterSelectPrompt := &survey.Select{
		Message: "Select the EKS cluster to login:",
		Options: clusterNames,
	}
	err := survey.AskOne(clusterSelectPrompt, &selectedClusterName)
	if err != nil {
		fmt.Println("Error during prompt:", err)
		return
	}

	var selectedCluster *EKSCluster
	for _, cluster := range config.EKSClusters {
		if cluster.ClusterName == selectedClusterName {
			selectedCluster = &cluster
			break
		}
	}

	if selectedCluster == nil {
		fmt.Println("Selected cluster not found.")
		return
	}

	// First, execute the AWS login command
	err = executeAWSLogin(config.LoginCommand, selectedCluster.Profile)
	if err != nil {
		fmt.Println("Error executing AWS login command:", err)
		return
	}

	// If AWS login is successful, execute the EKS login command
	err = executeEKSLogin(selectedCluster.ClusterName, selectedCluster.Profile)
	if err != nil {
		fmt.Println("Error executing EKS login command:", err)
	} else {
		fmt.Println("Successfully logged in with EKS cluster:", selectedCluster.ClusterName)
	}
}

// AWS login menu function
func awsLoginMenu(config Config) {
	profiles, err := getAWSProfiles()
	if err != nil {
		fmt.Println("Error reading AWS profiles:", err)
		return
	}

	if len(profiles) == 0 {
		fmt.Println("No profiles found in AWS config file.")
		return
	}

	profiles = append([]string{"All"}, profiles...)
	profiles = append(profiles, "Back")

	for {
		var selectedProfile string
		prompt := &survey.Select{
			Message: "AWS Login Menu:",
			Options: profiles,
		}

		err := survey.AskOne(prompt, &selectedProfile)
		if err != nil {
			if err.Error() == "interrupt" {
				fmt.Println("\nReturning to the main menu...")
				return
			}
			fmt.Println("Error during prompt:", err)
			return
		}

		if selectedProfile == "Back" {
			fmt.Println("Returning to the main menu...")
			break
		}

		if selectedProfile == "All" {
			for _, profile := range profiles[1 : len(profiles)-1] {
				fmt.Printf("Logging in with profile: %s\n", profile)
				err = executeAWSLogin(config.LoginCommand, profile)
				if err != nil {
					fmt.Printf("Error executing login command for profile %s: %s\n", profile, err)
				} else {
					fmt.Printf("Successfully logged in with profile: %s\n", profile)
				}
			}
			fmt.Println("All logins attempted. Returning to the main menu...")
			break
		}

		err = executeAWSLogin(config.LoginCommand, selectedProfile)
		if err != nil {
			fmt.Println("Error executing login command:", err)
		} else {
			fmt.Println("Successfully logged in with profile:", selectedProfile)
		}

		fmt.Println("Returning to the main menu...")
		break
	}
}

// Configuration menu for updating settings
func configMenu(config *Config) {
	for {
		options := []string{"Configure AWS Login Command", "Add EKS Cluster", "Edit EKS Cluster", "Exit"}
		prompt := &survey.Select{
			Message: "Configuration Menu:",
			Options: options,
		}

		var choice string
		err := survey.AskOne(prompt, &choice)
		if err != nil {
			fmt.Println("Error during prompt:", err)
			return
		}

		switch choice {
		case "Configure AWS Login Command":
			inputPrompt := &survey.Input{
				Message: "Enter the AWS login command to use:",
				Default: config.LoginCommand,
			}

			err := survey.AskOne(inputPrompt, &config.LoginCommand)
			if err != nil {
				fmt.Println("Error during prompt:", err)
				return
			}

			saveConfig(*config)
			fmt.Println("Configuration updated successfully!")

		case "Add EKS Cluster":
			var clusterName string
			clusterPrompt := &survey.Input{
				Message: "Enter the EKS cluster name:",
			}
			err := survey.AskOne(clusterPrompt, &clusterName)
			if err != nil {
				fmt.Println("Error during prompt:", err)
				return
			}

			profiles, err := getAWSProfiles()
			if err != nil {
				fmt.Println("Error reading AWS profiles:", err)
				return
			}

			var selectedProfile string
			profilePrompt := &survey.Select{
				Message: "Select the profile to use for this cluster:",
				Options: profiles,
			}
			err = survey.AskOne(profilePrompt, &selectedProfile)
			if err != nil {
				fmt.Println("Error during prompt:", err)
				return
			}

			newCluster := EKSCluster{
				ClusterName: clusterName,
				Profile:     selectedProfile,
			}
			config.EKSClusters = append(config.EKSClusters, newCluster)

			saveConfig(*config)
			fmt.Println("EKS Cluster configuration updated successfully!")

		case "Edit EKS Cluster":
			if len(config.EKSClusters) == 0 {
				fmt.Println("No EKS clusters to edit.")
				continue
			}

			clusterNames := make([]string, len(config.EKSClusters))
			for i, cluster := range config.EKSClusters {
				clusterNames[i] = cluster.ClusterName
			}

			var selectedClusterName string
			clusterSelectPrompt := &survey.Select{
				Message: "Select the EKS cluster to edit:",
				Options: clusterNames,
			}
			err := survey.AskOne(clusterSelectPrompt, &selectedClusterName)
			if err != nil {
				fmt.Println("Error during prompt:", err)
				return
			}

			var selectedCluster *EKSCluster
			for i, cluster := range config.EKSClusters {
				if cluster.ClusterName == selectedClusterName {
					selectedCluster = &config.EKSClusters[i]
					break
				}
			}

			if selectedCluster == nil {
				fmt.Println("Selected cluster not found.")
				continue
			}

			clusterNamePrompt := &survey.Input{
				Message: "Enter the new EKS cluster name:",
				Default: selectedCluster.ClusterName,
			}
			err = survey.AskOne(clusterNamePrompt, &selectedCluster.ClusterName)
			if err != nil {
				fmt.Println("Error during prompt:", err)
				return
			}

			profiles, err := getAWSProfiles()
			if err != nil {
				fmt.Println("Error reading AWS profiles:", err)
				return
			}

			profilePrompt := &survey.Select{
				Message: "Select the new profile to use for this cluster:",
				Options: profiles,
				Default: selectedCluster.Profile,
			}
			err = survey.AskOne(profilePrompt, &selectedCluster.Profile)
			if err != nil {
				fmt.Println("Error during prompt:", err)
				return
			}

			saveConfig(*config)
			fmt.Println("EKS Cluster configuration updated successfully!")

		case "Exit":
			return
		}
	}
}
