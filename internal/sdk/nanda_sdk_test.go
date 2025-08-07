package sdk

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateAgentID(t *testing.T) {
	agentID := generateAgentID()
	
	// Check that agent ID is within expected range (6 digits)
	if agentID < 100000 || agentID > 999999 {
		t.Errorf("Agent ID %d is not a 6-digit number", agentID)
	}
}

func TestNewNandaSDK(t *testing.T) {
	domain := "test.example.com"
	numAgents := 2
	registryURL := "https://test-registry.com:6900"
	agentID := 123456

	sdk := NewNandaSDK(domain, numAgents, registryURL, agentID)

	if sdk.Domain != domain {
		t.Errorf("Expected domain %s, got %s", domain, sdk.Domain)
	}

	if sdk.AgentID != agentID {
		t.Errorf("Expected agent ID %d, got %d", agentID, sdk.AgentID)
	}

	if sdk.NumAgents != numAgents {
		t.Errorf("Expected num agents %d, got %d", numAgents, sdk.NumAgents)
	}

	if sdk.RegistryURL != registryURL {
		t.Errorf("Expected registry URL %s, got %s", registryURL, sdk.RegistryURL)
	}
}

func TestNewNandaSDKWithGeneratedAgentID(t *testing.T) {
	domain := "test.example.com"
	numAgents := 1
	registryURL := "https://test-registry.com:6900"

	sdk := NewNandaSDK(domain, numAgents, registryURL, 0) // 0 means generate

	if sdk.AgentID == 0 {
		t.Error("Expected generated agent ID, got 0")
	}

	if sdk.AgentID < 100000 || sdk.AgentID > 999999 {
		t.Errorf("Generated agent ID %d is not a 6-digit number", sdk.AgentID)
	}
}

func TestFindPlaybookPath(t *testing.T) {
	sdk := &NandaSDK{}

	// Create a temporary test playbook file
	tempDir := t.TempDir()
	playbookPath := filepath.Join(tempDir, "ansible", "playbook.yml")
	err := os.MkdirAll(filepath.Dir(playbookPath), 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	err = os.WriteFile(playbookPath, []byte("test playbook content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test playbook: %v", err)
	}

	// Change to temp directory to test relative path resolution
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tempDir)

	foundPath := sdk.findPlaybookPath()
	if foundPath == "" {
		t.Error("Expected to find playbook, got empty path")
	}

	// Verify the found path exists
	if _, err := os.Stat(foundPath); os.IsNotExist(err) {
		t.Errorf("Found path does not exist: %s", foundPath)
	}
}

func TestFindPlaybookPathNotFound(t *testing.T) {
	sdk := &NandaSDK{}

	// Change to a directory without playbook
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tempDir)

	foundPath := sdk.findPlaybookPath()
	if foundPath != "" {
		t.Errorf("Expected empty path when playbook not found, got: %s", foundPath)
	}
} 