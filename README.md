# Go NANDA SDK

A Go SDK for setting up Internet of Agents servers. This tool automates the process of configuring servers with DNS records, SSL certificates, and required software.

## 🛠️ Setup Prerequisites

Before running the SDK, make sure you have the following:

### 1. ✅ AWS Account with a Running EC2 Linux Instance

Create an AWS account: https://aws.amazon.com
Launch an EC2 instance (any Linux distro, e.g., Amazon Linux, Ubuntu, Debian)
Allow the following ports in the security group:
22 (SSH), 80 (HTTP), 443 (HTTPS), 3000, 5001, 6000-6200, 8080, 6900
Save your .pem key file during instance creation — you'll need it to SSH.

### 2. 🌐 Domain or Subdomain with A Record

Register a domain (or use a subdomain) via Namecheap, GoDaddy, etc.
In your domain registrar's DNS settings, create an A Record pointing to your EC2 instance's public IPv4 address.
For root domains, use @ as the host.
For subdomains, use something like agent.yourdomain.com.

### 3. 🔑 Anthropic API Key

Sign up and request your API key from: https://www.anthropic.com

Once all the above is ready, proceed with installing and running the SDK below.

## Installation

### Prerequisites

SSH into the servers and install Go:

```bash
# For Ubuntu/Debian:
sudo apt update && sudo apt install -y golang-go

# For RHEL/CentOS/Fedora(Amazon Linux):
sudo dnf install -y golang
```

### Building the SDK

```bash
# Clone the repository
git clone https://github.com/destroyersrt/go-nanda.git
cd go-nanda

# Build the binary
go build -o go-nanda-sdk cmd/nanda-sdk/main.go

# Make it executable
chmod +x go-nanda-sdk
```

## Quick Setup Guide

### 1. Build the SDK
```bash
go build -o go-nanda-sdk cmd/nanda-sdk/main.go
```

### 2. Run the Setup
The setup requires two mandatory parameters:
- `--anthropic-key`: Your Anthropic API key
- `--domain`: Your complete domain name (e.g., myapp.example.com)

Optional parameters:
- `--smithery-key`: Your Smithery API key for connecting to MCP servers. A default key will be provided by application for connectivity
- `--agent-id`: A specific agent ID (if not provided, a random 6-digit number will be generated)
- `--num-agents`: Number of agents to set up (defaults to 1 if not specified)
- `--registry-url`: If the registry url needs to be changed. Default to https://chat.nanda-registry.com:6900
- `--verbose`: Enable verbose output for Ansible playbook

Example commands:
```bash
# Basic setup with random agent ID
./go-nanda-sdk --anthropic-key <your_anthropic_api_key> --domain <myapp.example.com> 

# Setup with specific agent ID
./go-nanda-sdk --anthropic-key <your_anthropic_api_key> --domain <myapp.example.com> --agent-id 123456

# Setup with multiple agents
./go-nanda-sdk --anthropic-key <your_anthropic_api_key> --domain <myapp.example.com> --num-agents 3

# Setup with your own smithery key
./go-nanda-sdk --anthropic-key <your_anthropic_api_key> --domain <myapp.example.com> --smithery-key <your_smithery_api_key>

# Setup with your own registry
./go-nanda-sdk --anthropic-key <your_anthropic_api_key> --domain <myapp.example.com> --registry-url <https://your-domain.com>

# Setup with verbose output
./go-nanda-sdk --anthropic-key <your_anthropic_api_key> --domain <myapp.example.com> --verbose
```

### 3. Verify Installation
After setup completes, verify your agent is running:

```bash
# Check service status
systemctl status internet_of_agents

# View logs
journalctl -u internet_of_agents -f
```

## Development

### Project Structure
```
go-nanda/
├── cmd/
│   └── nanda-sdk/
│       └── main.go          # Main entry point
├── internal/
│   └── sdk/
│       └── nanda_sdk.go     # Core SDK implementation
├── ansible/
│   ├── playbook.yml         # Ansible playbook
│   └── templates/           # Ansible templates
├── go.mod                   # Go module file
└── README.md               # This file
```

### Building for Different Platforms

```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o go-nanda-sdk-linux cmd/nanda-sdk/main.go

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o go-nanda-sdk-macos cmd/nanda-sdk/main.go

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o go-nanda-sdk.exe cmd/nanda-sdk/main.go
```

## Features

- **Public IP Detection**: Automatically detects the server's public IP address using multiple services
- **Ansible Integration**: Uses Ansible playbooks for server configuration
- **SSL Certificate Management**: Automatically configures SSL certificates via Let's Encrypt
- **Multi-Agent Support**: Configure multiple agents on a single server
- **Flexible Configuration**: Support for custom registry URLs and API keys
- **Verbose Logging**: Optional detailed output for debugging

## Dependencies

- Go 1.19 or later
- Ansible (must be installed on the target server)
- Internet connectivity for IP detection and certificate generation

## License

This project is licensed under the MIT License. 