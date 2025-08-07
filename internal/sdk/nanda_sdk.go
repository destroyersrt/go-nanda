package sdk

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// NandaSDK represents the NANDA SDK instance
type NandaSDK struct {
	Domain       string
	AgentID      int
	NumAgents    int
	RegistryURL  string
}

// GroupVars represents the Ansible group variables
type GroupVars struct {
	AnthropicAPIKey string `yaml:"anthropic_api_key"`
	SmitheryAPIKey  string `yaml:"smithery_api_key"`
	DomainName      string `yaml:"domain_name"`
	AgentIDPrefix   int    `yaml:"agent_id_prefix"`
	GithubRepo      string `yaml:"github_repo"`
	NumAgents       int    `yaml:"num_agents"`
	RegistryURL     string `yaml:"registry_url"`
}

// NewNandaSDK creates a new NandaSDK instance
func NewNandaSDK(domain string, numAgents int, registryURL string, agentID int) *NandaSDK {
	if agentID == 0 {
		agentID = generateAgentID()
	}

	log.Printf("Using agent ID: %d", agentID)
	log.Printf("Using domain: %s", domain)
	log.Printf("Using num_agents: %d", numAgents)
	log.Printf("Using registry URL: %s", registryURL)

	return &NandaSDK{
		Domain:      domain,
		AgentID:     agentID,
		NumAgents:   numAgents,
		RegistryURL: registryURL,
	}
}

// generateAgentID generates a random 6-digit agent ID
func generateAgentID() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(900000) + 100000 // 100000 to 999999
}

// getPublicIP gets the server's public IP address
func (n *NandaSDK) getPublicIP() (string, error) {
	ipServices := []string{
		"https://api.ipify.org",
		"https://ifconfig.me/ip",
		"https://icanhazip.com",
	}

	for _, service := range ipServices {
		resp, err := http.Get(service)
		if err != nil {
			log.Printf("Failed to get IP from %s: %v", service, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Failed to read response from %s: %v", service, err)
				continue
			}

			ip := string(body)
			log.Printf("Successfully detected public IP: %s", ip)
			return ip, nil
		}
	}

	return "", fmt.Errorf("failed to detect public IP from any service")
}

// createAnsibleInventory creates the Ansible inventory file
func (n *NandaSDK) createAnsibleInventory() (string, error) {
	serverIP, err := n.getPublicIP()
	if err != nil {
		return "", fmt.Errorf("failed to get public IP: %v", err)
	}

	inventoryContent := fmt.Sprintf(`[servers]
server ansible_host=%s

[all:vars]
ansible_user=root
ansible_connection=local
domain_name=%s
agent_id_prefix=%d
github_repo=https://github.com/aidecentralized/nanda-agent.git
registry_url=%s
`, serverIP, n.Domain, n.AgentID, n.RegistryURL)

	inventoryPath := "./ioa_inventory.ini"
	err = os.WriteFile(inventoryPath, []byte(inventoryContent), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write inventory file: %v", err)
	}

	return inventoryPath, nil
}

// setupServer sets up the server using Ansible
func (n *NandaSDK) setupServer(anthropicAPIKey string, smitheryAPIKey string, verbose bool) error {
	inventoryPath := ""
	groupVarsDir := ""

	defer func() {
		// Clean up temporary files
		if inventoryPath != "" {
			os.Remove(inventoryPath)
		}
		if groupVarsDir != "" {
			os.Remove(filepath.Join(groupVarsDir, "all.yml"))
			os.RemoveAll(groupVarsDir)
		}
	}()

	// Create Ansible inventory
	var err error
	inventoryPath, err = n.createAnsibleInventory()
	if err != nil {
		return fmt.Errorf("failed to create inventory: %v", err)
	}
	log.Printf("Created inventory file at %s", inventoryPath)

	// Create group_vars directory and file
	groupVarsDir = "./group_vars"
	err = os.MkdirAll(groupVarsDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create group_vars directory: %v", err)
	}
	log.Printf("Created group_vars directory at %s", groupVarsDir)

	// Create group_vars/all.yml
	groupVars := GroupVars{
		AnthropicAPIKey: anthropicAPIKey,
		SmitheryAPIKey:  smitheryAPIKey,
		DomainName:      n.Domain,
		AgentIDPrefix:   n.AgentID,
		GithubRepo:      "https://github.com/aidecentralized/nanda-agent.git",
		NumAgents:       n.NumAgents,
		RegistryURL:     n.RegistryURL,
	}

	groupVarsData, err := yaml.Marshal(groupVars)
	if err != nil {
		return fmt.Errorf("failed to marshal group vars: %v", err)
	}

	groupVarsPath := filepath.Join(groupVarsDir, "all.yml")
	err = os.WriteFile(groupVarsPath, groupVarsData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write group vars file: %v", err)
	}

	// Get the path to the local Ansible playbook
	// Look for playbook relative to the SDK package or in common locations
	playbookPath := n.findPlaybookPath()
	if playbookPath == "" {
		return fmt.Errorf("ansible playbook not found")
	}
	log.Printf("Using playbook at: %s", playbookPath)

	// Run Ansible playbook with optional verbose output
	verboseFlag := ""
	if verbose {
		verboseFlag = "-vvv"
	}

	cmd := exec.Command("ansible-playbook", "-i", inventoryPath, playbookPath)
	if verboseFlag != "" {
		cmd.Args = append(cmd.Args, verboseFlag)
	}

	log.Printf("Running command: %v", cmd.Args)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Ansible playbook error: %s", string(output))
		return fmt.Errorf("ansible playbook failed: %v", err)
	}

	log.Printf("Ansible playbook output: %s", string(output))

	// Check if the playbook failed by looking for "failed=1" in the output
	if strings.Contains(string(output), "failed=1") {
		return fmt.Errorf("ansible playbook failed with errors")
	}

	log.Println("Server setup completed successfully")
	return nil
}

// findPlaybookPath locates the Ansible playbook file
func (n *NandaSDK) findPlaybookPath() string {
	// Try multiple locations for the Go SDK playbook
	possiblePaths := []string{
		"ansible/playbook.yml",                        // Current directory (development)
		"../ansible/playbook.yml",                    // Parent directory
		"/usr/local/share/go-nanda-sdk/ansible/playbook.yml", // System installation
		"/opt/go-nanda-sdk/ansible/playbook.yml",     // Alternative system location
		"./ansible/playbook.yml",                     // Explicit current directory
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			absPath, _ := filepath.Abs(path)
			return absPath
		}
	}

	return ""
}

// Setup performs the complete setup process
func (n *NandaSDK) Setup(anthropicAPIKey, smitheryAPIKey string, verbose bool) error {
	if err := n.setupServer(anthropicAPIKey, smitheryAPIKey, verbose); err != nil {
		return fmt.Errorf("setup server failed: %v", err)
	}

	log.Println("Setup completed successfully")
	return nil
} 