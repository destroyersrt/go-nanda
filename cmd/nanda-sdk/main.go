package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"go-nanda/internal/sdk"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// .env file is optional, so we don't fail if it doesn't exist
		log.Printf("No .env file found or error loading it: %v", err)
	}

	var rootCmd = &cobra.Command{
		Use:   "nanda-sdk",
		Short: "NANDA SDK for setting up Internet of Agents servers",
		Long: `A Go SDK for setting up Internet of Agents servers. 
This tool automates the process of configuring servers with DNS records, SSL certificates, and required software.`,
	}

	var anthropicKey string
	var domain string
	var smitheryKey string
	var agentID int
	var numAgents int
	var registryURL string
	var verbose bool

	// Set default values from environment variables (from .env file)
	anthropicKey = getEnvOrDefault("ANTHROPIC_API_KEY", "")
	domain = getEnvOrDefault("DOMAIN", "")
	smitheryKey = getEnvOrDefault("SMITHERY_API_KEY", "")
	agentID = getEnvIntOrDefault("AGENT_ID", 0)
	numAgents = getEnvIntOrDefault("NUM_AGENTS", 1)
	registryURL = getEnvOrDefault("REGISTRY_URL", "https://chat.nanda-registry.com:6900")
	verbose = getEnvBoolOrDefault("VERBOSE", false)

	rootCmd.Flags().StringVar(&anthropicKey, "anthropic-key", anthropicKey, "Anthropic API key for the agent (required)")
	rootCmd.Flags().StringVar(&domain, "domain", domain, "Complete domain name (e.g., myapp.example.com) (required)")
	rootCmd.Flags().StringVar(&smitheryKey, "smithery-key", smitheryKey, "Optional Smithery API key for the MCP connections")
	rootCmd.Flags().IntVar(&agentID, "agent-id", agentID, "Optional agent ID (if not provided, will generate one)")
	rootCmd.Flags().IntVar(&numAgents, "num-agents", numAgents, "Optional number of agents (if not provided, will default to one)")
	rootCmd.Flags().StringVar(&registryURL, "registry-url", registryURL, "URL of the NANDA registry")
	rootCmd.Flags().BoolVar(&verbose, "verbose", verbose, "Enable verbose output for Ansible playbook")

	// Only mark flags as required if they weren't provided via .env
	if anthropicKey == "" {
		rootCmd.MarkFlagRequired("anthropic-key")
	}
	if domain == "" {
		rootCmd.MarkFlagRequired("domain")
	}

	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		// Validate required parameters (whether from .env or flags)
		if anthropicKey == "" {
			return fmt.Errorf("anthropic API key is required (set ANTHROPIC_API_KEY in .env or use --anthropic-key)")
		}
		if domain == "" {
			return fmt.Errorf("domain is required (set DOMAIN in .env or use --domain)")
		}

		// Use default smithery key if not provided
		if smitheryKey == "" {
			smitheryKey = "b4e92d35-0034-43f0-beff-042466777ada"
		}

		// Create SDK instance
		nandaSDK := sdk.NewNandaSDK(domain, numAgents, registryURL, agentID)

		// Run setup
		if err := nandaSDK.Setup(anthropicKey, smitheryKey, verbose); err != nil {
			fmt.Printf("Setup failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Setup completed successfully")
		return nil
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Helper functions to get environment variables with defaults
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBoolOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
} 