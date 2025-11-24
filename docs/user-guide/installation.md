# PHP-Go Installation Guide

This guide covers all the methods for installing PHP-Go on various platforms and environments.

## Table of Contents

1. [System Requirements](#system-requirements)
2. [Installation Methods](#installation-methods)
3. [Platform-Specific Instructions](#platform-specific-instructions)
4. [Building from Source](#building-from-source)
5. [Docker Installation](#docker-installation)
6. [Verification](#verification)
7. [Updating PHP-Go](#updating-php-go)
8. [Uninstallation](#uninstallation)
9. [Troubleshooting](#troubleshooting)

## System Requirements

### Minimum Requirements

- **CPU**: x86_64 (AMD64) or ARM64 architecture
- **Memory**: 256 MB RAM (1 GB+ recommended for production)
- **Disk Space**: 50 MB for the binary, additional space for your PHP applications
- **Operating System**:
  - Linux (kernel 3.2+)
  - macOS 10.13+ (High Sierra or later)
  - Windows 10+ or Windows Server 2016+
  - FreeBSD 11+

### For Building from Source

- **Go**: Version 1.21 or later (Go 1.25+ recommended)
- **Git**: For cloning the repository
- **Make**: Optional, but recommended for simplified builds
- **GCC/Clang**: Required for CGo-enabled builds (optional feature)

### Recommended System Configuration

For optimal performance:
- **CPU**: Multi-core processor (2+ cores to benefit from parallelization)
- **Memory**: 2 GB+ RAM for production workloads
- **Disk**: SSD storage for better I/O performance
- **Network**: 100 Mbps+ for web applications

## Installation Methods

PHP-Go can be installed using several methods:

1. **Pre-built Binaries** (Recommended for users) - Coming soon
2. **Building from Source** (Recommended for developers)
3. **Docker Container** (Recommended for deployment)
4. **Package Managers** - Coming soon
   - Homebrew (macOS/Linux)
   - apt/yum/dnf (Linux)
   - Chocolatey/Scoop (Windows)

## Platform-Specific Instructions

### Linux

#### Ubuntu/Debian

**Method 1: Pre-built Binary** (Coming soon)

```bash
# Download the latest release
wget https://github.com/krizos/php-go/releases/latest/download/php-go-linux-amd64.tar.gz

# Extract the archive
tar -xzf php-go-linux-amd64.tar.gz

# Move to system path
sudo mv php-go /usr/local/bin/

# Make executable
sudo chmod +x /usr/local/bin/php-go

# Verify installation
php-go --version
```

**Method 2: Build from Source**

```bash
# Install prerequisites
sudo apt update
sudo apt install -y git golang-go

# Verify Go version (should be 1.21+)
go version

# Clone repository
git clone https://github.com/krizos/php-go.git
cd php-go

# Build
go build -o php-go ./cmd/php-go

# Install globally (optional)
sudo cp php-go /usr/local/bin/

# Verify
php-go --version
```

#### RHEL/CentOS/Fedora

```bash
# Install prerequisites
sudo dnf install -y git golang

# Or for older systems (RHEL/CentOS 7)
# sudo yum install -y git golang

# Clone and build
git clone https://github.com/krizos/php-go.git
cd php-go
go build -o php-go ./cmd/php-go

# Install globally
sudo cp php-go /usr/local/bin/

# Verify
php-go --version
```

#### Arch Linux

```bash
# Install prerequisites
sudo pacman -S git go

# Clone and build
git clone https://github.com/krizos/php-go.git
cd php-go
go build -o php-go ./cmd/php-go

# Install globally
sudo cp php-go /usr/local/bin/

# Verify
php-go --version
```

### macOS

#### Using Homebrew (Coming soon)

```bash
# Install Homebrew if not already installed
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Install PHP-Go
brew tap krizos/php-go
brew install php-go

# Verify
php-go --version
```

#### Build from Source

```bash
# Install Xcode Command Line Tools (if not already installed)
xcode-select --install

# Install Go using Homebrew
brew install go

# Or download from https://go.dev/dl/

# Clone repository
git clone https://github.com/krizos/php-go.git
cd php-go

# Build
go build -o php-go ./cmd/php-go

# Install globally (optional)
sudo cp php-go /usr/local/bin/

# Verify
php-go --version
```

### Windows

#### Using Pre-built Binary (Coming soon)

1. Download the latest Windows release from [GitHub Releases](https://github.com/krizos/php-go/releases)
2. Extract the ZIP file
3. Add the extracted directory to your PATH:
   - Right-click "This PC" → Properties
   - Click "Advanced system settings"
   - Click "Environment Variables"
   - Under "System Variables", find "Path" and click "Edit"
   - Click "New" and add the directory containing php-go.exe
4. Open a new Command Prompt or PowerShell window
5. Verify: `php-go --version`

#### Build from Source (PowerShell)

```powershell
# Install Go from https://go.dev/dl/
# Install Git from https://git-scm.com/download/win

# Clone repository
git clone https://github.com/krizos/php-go.git
cd php-go

# Build
go build -o php-go.exe ./cmd/php-go

# Add to PATH or copy to a directory in PATH
# Example: copy php-go.exe to C:\Windows\System32 (requires admin)
# Or add current directory to PATH

# Verify
php-go --version
```

#### Using WSL (Windows Subsystem for Linux)

If you're using WSL, follow the Linux instructions above for your specific distribution.

### FreeBSD

```bash
# Install prerequisites
sudo pkg install git go

# Clone and build
git clone https://github.com/krizos/php-go.git
cd php-go
go build -o php-go ./cmd/php-go

# Install globally
sudo cp php-go /usr/local/bin/

# Verify
php-go --version
```

## Building from Source

### Quick Build

The simplest way to build PHP-Go:

```bash
# Clone the repository
git clone https://github.com/krizos/php-go.git
cd php-go

# Build
go build -o php-go ./cmd/php-go

# Run
./php-go --version
```

### Development Build

For development with debugging symbols:

```bash
# Clone the repository
git clone https://github.com/krizos/php-go.git
cd php-go

# Build with race detector and debugging info
go build -race -gcflags="all=-N -l" -o php-go ./cmd/php-go

# Run tests
go test ./...

# Run with coverage
go test -cover ./...
```

### Optimized Production Build

For production deployments with maximum performance:

```bash
# Clone the repository
git clone https://github.com/krizos/php-go.git
cd php-go

# Build with optimizations and without debug symbols
go build -ldflags="-s -w" -trimpath -o php-go ./cmd/php-go

# Optional: Further compress with UPX (if installed)
upx --best --lzma php-go
```

Build flags explained:
- `-ldflags="-s -w"`: Strip debug information and symbol table
- `-trimpath`: Remove file system paths from binary
- `-race`: Enable data race detector (development only, adds overhead)
- `-gcflags="all=-N -l"`: Disable optimizations for debugging

### Cross-Compilation

Build for different platforms:

```bash
# Build for Linux AMD64
GOOS=linux GOARCH=amd64 go build -o php-go-linux-amd64 ./cmd/php-go

# Build for Linux ARM64
GOOS=linux GOARCH=arm64 go build -o php-go-linux-arm64 ./cmd/php-go

# Build for macOS AMD64 (Intel)
GOOS=darwin GOARCH=amd64 go build -o php-go-darwin-amd64 ./cmd/php-go

# Build for macOS ARM64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o php-go-darwin-arm64 ./cmd/php-go

# Build for Windows AMD64
GOOS=windows GOARCH=amd64 go build -o php-go-windows-amd64.exe ./cmd/php-go

# Build for FreeBSD AMD64
GOOS=freebsd GOARCH=amd64 go build -o php-go-freebsd-amd64 ./cmd/php-go
```

### Build with Version Information

Embed version information into the binary:

```bash
# Set version variables
VERSION=$(git describe --tags --always --dirty)
COMMIT=$(git rev-parse --short HEAD)
BUILD_DATE=$(date -u '+%Y-%m-%d_%H:%M:%S')

# Build with version info
go build -ldflags="-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.buildDate=${BUILD_DATE}" \
  -o php-go ./cmd/php-go

# Verify version info
./php-go --version
```

### Custom Build Configurations

#### Minimal Build (Smaller Binary)

```bash
# Build without certain features to reduce size
go build -tags=no_cgo,no_extensions -ldflags="-s -w" -o php-go ./cmd/php-go
```

#### Full-Featured Build

```bash
# Build with all optional features
go build -tags=full,cgo,extensions,profiling -o php-go ./cmd/php-go
```

## Docker Installation

### Using Official Docker Image (Coming soon)

```bash
# Pull the latest image
docker pull krizos/php-go:latest

# Run a PHP script
docker run --rm -v $(pwd):/app krizos/php-go:latest run /app/script.php

# Interactive shell
docker run --rm -it krizos/php-go:latest sh

# Built-in web server
docker run --rm -p 8000:8000 -v $(pwd):/app krizos/php-go:latest serve -S 0.0.0.0:8000 -t /app
```

### Building Your Own Docker Image

Create a `Dockerfile`:

```dockerfile
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /build

# Copy source
COPY . .

# Build
RUN go build -ldflags="-s -w" -trimpath -o php-go ./cmd/php-go

# Create minimal runtime image
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /build/php-go /usr/local/bin/php-go

# Set working directory
WORKDIR /app

# Default command
ENTRYPOINT ["php-go"]
CMD ["--help"]
```

Build and run:

```bash
# Build image
docker build -t php-go:local .

# Run
docker run --rm php-go:local --version

# Run script
docker run --rm -v $(pwd):/app php-go:local run /app/script.php
```

### Docker Compose

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  php-go:
    image: krizos/php-go:latest
    volumes:
      - ./:/app
    working_dir: /app
    command: serve -S 0.0.0.0:8000 -t /app
    ports:
      - "8000:8000"
    environment:
      - PHP_GO_ENV=production
    restart: unless-stopped
```

Run:

```bash
docker-compose up -d
```

## Verification

After installation, verify that PHP-Go is working correctly:

### Basic Verification

```bash
# Check version
php-go --version

# Should output something like:
# PHP-Go v0.0.1-dev (built with Go 1.25.4)

# Check help
php-go --help

# Should display usage information
```

### Test Execution

Create a test script `test.php`:

```php
<?php
echo "PHP-Go is working!\n";
echo "PHP Version: " . phpversion() . "\n";
echo "System: " . php_uname() . "\n";
```

Run it:

```bash
php-go run test.php

# Expected output:
# PHP-Go is working!
# PHP Version: 8.4.0
# System: [your system info]
```

### Advanced Tests

```bash
# Test lexer
php-go lex test.php

# Test parser
php-go parse test.php

# Run demo
php-go demo

# Run benchmarks (if available)
cd php-go
go test -bench=. ./pkg/...
```

## Updating PHP-Go

### From Pre-built Binaries

```bash
# Download latest release
wget https://github.com/krizos/php-go/releases/latest/download/php-go-linux-amd64.tar.gz

# Replace existing binary
tar -xzf php-go-linux-amd64.tar.gz
sudo mv php-go /usr/local/bin/

# Verify new version
php-go --version
```

### From Source

```bash
# Navigate to repository
cd php-go

# Pull latest changes
git pull origin main

# Rebuild
go build -o php-go ./cmd/php-go

# Reinstall
sudo cp php-go /usr/local/bin/

# Verify
php-go --version
```

### Using Package Managers (Coming soon)

```bash
# Homebrew (macOS/Linux)
brew upgrade php-go

# apt (Debian/Ubuntu)
sudo apt update && sudo apt upgrade php-go

# yum/dnf (RHEL/CentOS/Fedora)
sudo dnf update php-go
```

## Uninstallation

### Manual Installation

```bash
# Remove binary
sudo rm /usr/local/bin/php-go

# Remove repository (if cloned)
rm -rf ~/php-go

# Remove configuration (if any)
rm -rf ~/.php-go
```

### Package Manager Installation (Coming soon)

```bash
# Homebrew
brew uninstall php-go

# apt
sudo apt remove php-go

# yum/dnf
sudo dnf remove php-go
```

### Docker

```bash
# Remove Docker image
docker rmi krizos/php-go:latest

# Remove all PHP-Go images
docker rmi $(docker images krizos/php-go -q)
```

## Troubleshooting

### Common Issues

#### Issue: "command not found: php-go"

**Solution**: The binary is not in your PATH.

```bash
# Check if binary exists
ls -l /usr/local/bin/php-go

# If it exists, add to PATH
export PATH=$PATH:/usr/local/bin

# Make permanent (add to ~/.bashrc or ~/.zshrc)
echo 'export PATH=$PATH:/usr/local/bin' >> ~/.bashrc
source ~/.bashrc
```

#### Issue: "permission denied" when running php-go

**Solution**: Make the binary executable.

```bash
sudo chmod +x /usr/local/bin/php-go
```

#### Issue: Build fails with "go: cannot find main module"

**Solution**: Ensure you're in the correct directory with go.mod.

```bash
cd php-go
ls go.mod  # Should exist
go build ./cmd/php-go
```

#### Issue: Build fails with Go version error

**Solution**: Update Go to version 1.21 or later.

```bash
# Check current version
go version

# Update Go (varies by platform)
# Linux: Download from https://go.dev/dl/
# macOS: brew upgrade go
# Windows: Download installer from https://go.dev/dl/
```

#### Issue: "too many open files" error

**Solution**: Increase file descriptor limit.

```bash
# Temporary fix (current session)
ulimit -n 4096

# Permanent fix (Linux)
echo "* soft nofile 4096" | sudo tee -a /etc/security/limits.conf
echo "* hard nofile 65536" | sudo tee -a /etc/security/limits.conf
```

#### Issue: Slow performance

**Solution**: Ensure you're using an optimized build.

```bash
# Rebuild with optimizations
go build -ldflags="-s -w" -o php-go ./cmd/php-go

# Check if running in development mode
php-go --version  # Should not show debug info in production
```

#### Issue: Docker container exits immediately

**Solution**: Provide a command or script to run.

```bash
# Bad
docker run krizos/php-go:latest

# Good
docker run krizos/php-go:latest run script.php
# Or
docker run -it krizos/php-go:latest sh
```

### Platform-Specific Issues

#### macOS: "Developer Tools not found"

```bash
# Install Xcode Command Line Tools
xcode-select --install
```

#### macOS: "cannot be opened because the developer cannot be verified"

```bash
# Allow the binary
xattr -d com.apple.quarantine /usr/local/bin/php-go
```

#### Windows: "The term 'php-go' is not recognized"

```powershell
# Add directory to PATH
$env:Path += ";C:\path\to\php-go"

# Make permanent
[System.Environment]::SetEnvironmentVariable("Path", $env:Path, [System.EnvironmentVariableTarget]::User)
```

#### Linux: SELinux prevents execution

```bash
# Check SELinux status
getenforce

# Allow execution (temporary)
sudo setenforce 0

# Or set correct context (permanent)
sudo chcon -t bin_t /usr/local/bin/php-go
```

### Getting Help

If you encounter issues not covered here:

1. **Check documentation**: Read the [User Guide](README.md)
2. **Search issues**: Check [GitHub Issues](https://github.com/krizos/php-go/issues)
3. **Ask questions**: Open a new issue with:
   - Your operating system and version
   - Go version (`go version`)
   - PHP-Go version (`php-go --version`)
   - Complete error message
   - Steps to reproduce
4. **Debug mode**: Run with verbose output:
   ```bash
   php-go --debug run script.php
   ```

### Diagnostic Information

Collect diagnostic information for bug reports:

```bash
# System information
uname -a

# Go version
go version

# PHP-Go version
php-go --version

# Environment
env | grep -i php

# Check for conflicts
which php
which php-go
```

## Next Steps

After successful installation:

1. Read the [User Guide](README.md)
2. Try the [Getting Started](getting-started.md) tutorial
3. Explore [Configuration Options](configuration.md)
4. Review [Performance Tuning](performance.md)
5. Check out [Example Applications](../examples/)

## Additional Resources

- [PHP-Go Repository](https://github.com/krizos/php-go)
- [Release Notes](https://github.com/krizos/php-go/releases)
- [Contributing Guide](../../CONTRIBUTING.md)
- [Architecture Documentation](../02-go-architecture.md)
- [Development Guide](../internals/development.md)
