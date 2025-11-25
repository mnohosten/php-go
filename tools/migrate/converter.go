// Package migrate provides tools for migrating PHP applications to PHP-Go
package migrate

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Config represents a PHP-Go configuration
type Config struct {
	PHP              PHPSettings          `yaml:"php"`
	Extensions       ExtensionSettings    `yaml:"extensions"`
	Parallelization  ParallelSettings     `yaml:"parallelization,omitempty"`
	Logging          LogSettings          `yaml:"logging,omitempty"`
	Metrics          MetricsSettings      `yaml:"metrics,omitempty"`
	Health           HealthSettings       `yaml:"health,omitempty"`
	GoIntegration    GoIntegrationSettings `yaml:"go_integration,omitempty"`
}

// PHPSettings contains PHP-compatible configuration
type PHPSettings struct {
	ErrorReporting      string      `yaml:"error_reporting"`
	DisplayErrors       bool        `yaml:"display_errors"`
	LogErrors           bool        `yaml:"log_errors"`
	ErrorLog            string      `yaml:"error_log,omitempty"`
	MemoryLimit         string      `yaml:"memory_limit"`
	MaxExecutionTime    int         `yaml:"max_execution_time"`
	PostMaxSize         string      `yaml:"post_max_size"`
	UploadMaxFilesize   string      `yaml:"upload_max_filesize"`
	Date                DateSettings `yaml:"date"`
}

// DateSettings contains date/time configuration
type DateSettings struct {
	Timezone string `yaml:"timezone"`
}

// ExtensionSettings contains extension configuration
type ExtensionSettings struct {
	Enabled []string `yaml:"enabled"`
	Paths   []string `yaml:"paths,omitempty"`
}

// ParallelSettings contains parallelization configuration
type ParallelSettings struct {
	Enabled         bool                   `yaml:"enabled"`
	MaxWorkers      int                    `yaml:"max_workers"`
	ArrayOperations ArrayParallelSettings  `yaml:"array_operations"`
}

// ArrayParallelSettings contains array parallelization settings
type ArrayParallelSettings struct {
	Enabled bool `yaml:"enabled"`
	MinSize int  `yaml:"min_size"`
}

// LogSettings contains logging configuration
type LogSettings struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	Output string `yaml:"output"`
}

// MetricsSettings contains metrics configuration
type MetricsSettings struct {
	Enabled bool   `yaml:"enabled"`
	Port    int    `yaml:"port"`
	Path    string `yaml:"path"`
}

// HealthSettings contains health check configuration
type HealthSettings struct {
	Enabled bool   `yaml:"enabled"`
	Port    int    `yaml:"port"`
	Path    string `yaml:"path"`
}

// GoIntegrationSettings contains Go integration configuration
type GoIntegrationSettings struct {
	Enabled         bool   `yaml:"enabled"`
	PluginDirectory string `yaml:"plugin_directory,omitempty"`
	ImportPath      string `yaml:"import_path,omitempty"`
}

// ConfigConverter converts php.ini files to PHP-Go YAML configuration
type ConfigConverter struct {
	// Input is the parsed INI data
	Input map[string]map[string]string

	// Extensions tracks all extension directives found
	Extensions []string

	// Warnings contains any warnings during conversion
	Warnings []string
}

// NewConfigConverter creates a new config converter
func NewConfigConverter() *ConfigConverter {
	return &ConfigConverter{
		Input:      make(map[string]map[string]string),
		Extensions: []string{},
		Warnings:   []string{},
	}
}

// ParseINI parses a php.ini file
func (c *ConfigConverter) ParseINI(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentSection := "PHP" // Default section
	c.Input[currentSection] = make(map[string]string)

	// Regex patterns
	sectionPattern := regexp.MustCompile(`^\[([^\]]+)\]`)
	settingPattern := regexp.MustCompile(`^([a-zA-Z0-9_.]+)\s*=\s*(.*)$`)
	commentPattern := regexp.MustCompile(`^\s*[;#]`)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || commentPattern.MatchString(line) {
			continue
		}

		// Check for section header
		if matches := sectionPattern.FindStringSubmatch(line); matches != nil {
			currentSection = matches[1]
			if c.Input[currentSection] == nil {
				c.Input[currentSection] = make(map[string]string)
			}
			continue
		}

		// Parse setting
		if matches := settingPattern.FindStringSubmatch(line); matches != nil {
			key := strings.TrimSpace(matches[1])
			value := strings.TrimSpace(matches[2])

			// Remove quotes if present
			if len(value) >= 2 && (value[0] == '"' && value[len(value)-1] == '"' ||
			                       value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}

			// Remove inline comments
			if idx := strings.Index(value, ";"); idx != -1 {
				value = strings.TrimSpace(value[:idx])
			}

			// Handle extension directives specially to collect all of them
			if key == "extension" && value != "" {
				// Remove .so or .dll extension if present
				extName := strings.TrimSuffix(value, ".so")
				extName = strings.TrimSuffix(extName, ".dll")
				extName = strings.TrimSpace(extName)
				if extName != "" {
					c.Extensions = append(c.Extensions, extName)
				}
			} else if key == "zend_extension" {
				c.Warnings = append(c.Warnings, fmt.Sprintf("zend_extension in section [%s] is not supported", currentSection))
			}

			c.Input[currentSection][key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	return nil
}

// Convert converts the parsed INI data to PHP-Go configuration
func (c *ConfigConverter) Convert() (*Config, error) {
	config := &Config{
		PHP: PHPSettings{
			ErrorReporting:    "E_ALL",
			DisplayErrors:     true,
			LogErrors:         true,
			MemoryLimit:       "128M",
			MaxExecutionTime:  30,
			PostMaxSize:       "8M",
			UploadMaxFilesize: "2M",
			Date: DateSettings{
				Timezone: "UTC",
			},
		},
		Extensions: ExtensionSettings{
			Enabled: []string{},
			Paths:   []string{},
		},
	}

	// Process PHP section
	if phpSection, ok := c.Input["PHP"]; ok {
		c.processPHPSettings(phpSection, &config.PHP)
	}

	// Process extensions
	c.processExtensions(&config.Extensions)

	// Add default PHP-Go specific settings with conservative defaults
	config.Parallelization = ParallelSettings{
		Enabled:    false, // Disabled by default for compatibility
		MaxWorkers: 4,
		ArrayOperations: ArrayParallelSettings{
			Enabled: false,
			MinSize: 1000,
		},
	}

	config.Logging = LogSettings{
		Level:  "info",
		Format: "text",
		Output: "",
	}

	config.Metrics = MetricsSettings{
		Enabled: false,
		Port:    9090,
		Path:    "/metrics",
	}

	config.Health = HealthSettings{
		Enabled: false,
		Port:    8080,
		Path:    "/health",
	}

	config.GoIntegration = GoIntegrationSettings{
		Enabled: false,
	}

	return config, nil
}

// processPHPSettings processes PHP section settings
func (c *ConfigConverter) processPHPSettings(settings map[string]string, php *PHPSettings) {
	for key, value := range settings {
		switch key {
		case "error_reporting":
			php.ErrorReporting = c.convertErrorReporting(value)
		case "display_errors":
			php.DisplayErrors = c.convertBool(value)
		case "log_errors":
			php.LogErrors = c.convertBool(value)
		case "error_log":
			php.ErrorLog = value
		case "memory_limit":
			php.MemoryLimit = c.normalizeMemorySize(value)
		case "max_execution_time":
			if intVal, err := strconv.Atoi(value); err == nil {
				php.MaxExecutionTime = intVal
			}
		case "post_max_size":
			php.PostMaxSize = c.normalizeMemorySize(value)
		case "upload_max_filesize":
			php.UploadMaxFilesize = c.normalizeMemorySize(value)
		case "date.timezone":
			php.Date.Timezone = value
		default:
			// Check if it's a date setting without section prefix
			if strings.HasPrefix(key, "date.") {
				subkey := strings.TrimPrefix(key, "date.")
				if subkey == "timezone" {
					php.Date.Timezone = value
				}
			}
		}
	}

	// Check date section separately
	if dateSection, ok := c.Input["Date"]; ok {
		if tz, ok := dateSection["date.timezone"]; ok {
			php.Date.Timezone = tz
		} else if tz, ok := dateSection["timezone"]; ok {
			php.Date.Timezone = tz
		}
	}
}

// processExtensions processes extension settings
func (c *ConfigConverter) processExtensions(ext *ExtensionSettings) {
	// Use the extensions collected during parsing
	// Remove duplicates using a map
	extensions := make(map[string]bool)
	for _, extName := range c.Extensions {
		extensions[extName] = true
	}

	// Convert map to list
	for name := range extensions {
		ext.Enabled = append(ext.Enabled, name)
	}
}

// convertErrorReporting converts error_reporting value to a readable format
func (c *ConfigConverter) convertErrorReporting(value string) string {
	// If it's already a constant name, keep it
	if strings.Contains(value, "E_") {
		return value
	}

	// If it's a number, try to convert to constant names
	if intVal, err := strconv.Atoi(value); err == nil {
		// Common error_reporting values
		switch intVal {
		case 32767:
			return "E_ALL"
		case 30719:
			return "E_ALL & ~E_DEPRECATED & ~E_STRICT"
		case 22527:
			return "E_ALL & ~E_NOTICE"
		case 0:
			return "0"
		default:
			return value
		}
	}

	return value
}

// convertBool converts string boolean values to bool
func (c *ConfigConverter) convertBool(value string) bool {
	value = strings.ToLower(value)
	return value == "on" || value == "1" || value == "true" || value == "yes"
}

// normalizeMemorySize normalizes memory size values
func (c *ConfigConverter) normalizeMemorySize(value string) string {
	value = strings.ToUpper(strings.TrimSpace(value))

	// Check if it already has a suffix
	if len(value) > 0 {
		lastChar := value[len(value)-1]
		if lastChar == 'K' || lastChar == 'M' || lastChar == 'G' {
			return value
		}
	}

	// If it's just a number, assume bytes
	if _, err := strconv.Atoi(value); err == nil {
		// Convert to appropriate unit
		if intVal, _ := strconv.ParseInt(value, 10, 64); intVal > 0 {
			if intVal >= 1024*1024*1024 {
				return fmt.Sprintf("%dG", intVal/(1024*1024*1024))
			} else if intVal >= 1024*1024 {
				return fmt.Sprintf("%dM", intVal/(1024*1024))
			} else if intVal >= 1024 {
				return fmt.Sprintf("%dK", intVal/1024)
			}
		}
	}

	return value
}

// ToYAML converts the configuration to YAML format
func (c *Config) ToYAML() string {
	var sb strings.Builder

	sb.WriteString("# PHP-Go Configuration\n")
	sb.WriteString("# Generated from php.ini\n\n")

	// PHP settings
	sb.WriteString("# PHP-compatible settings\n")
	sb.WriteString("php:\n")
	sb.WriteString(fmt.Sprintf("  error_reporting: %s\n", c.PHP.ErrorReporting))
	sb.WriteString(fmt.Sprintf("  display_errors: %t\n", c.PHP.DisplayErrors))
	sb.WriteString(fmt.Sprintf("  log_errors: %t\n", c.PHP.LogErrors))
	if c.PHP.ErrorLog != "" {
		sb.WriteString(fmt.Sprintf("  error_log: %s\n", c.PHP.ErrorLog))
	}
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("  memory_limit: %s\n", c.PHP.MemoryLimit))
	sb.WriteString(fmt.Sprintf("  max_execution_time: %d\n", c.PHP.MaxExecutionTime))
	sb.WriteString(fmt.Sprintf("  post_max_size: %s\n", c.PHP.PostMaxSize))
	sb.WriteString(fmt.Sprintf("  upload_max_filesize: %s\n", c.PHP.UploadMaxFilesize))
	sb.WriteString("\n")
	sb.WriteString("  date:\n")
	sb.WriteString(fmt.Sprintf("    timezone: %s\n", c.PHP.Date.Timezone))

	// Extensions
	sb.WriteString("\n# Extensions\n")
	sb.WriteString("extensions:\n")
	if len(c.Extensions.Enabled) > 0 {
		sb.WriteString("  enabled:\n")
		for _, ext := range c.Extensions.Enabled {
			sb.WriteString(fmt.Sprintf("    - %s\n", ext))
		}
	} else {
		sb.WriteString("  enabled: []\n")
	}

	if len(c.Extensions.Paths) > 0 {
		sb.WriteString("\n  paths:\n")
		for _, path := range c.Extensions.Paths {
			sb.WriteString(fmt.Sprintf("    - %s\n", path))
		}
	}

	// PHP-Go specific features
	sb.WriteString("\n# PHP-Go specific features\n")
	sb.WriteString("# Note: These settings can be enabled to take advantage of PHP-Go's\n")
	sb.WriteString("# additional features. They are disabled by default for compatibility.\n\n")

	sb.WriteString("parallelization:\n")
	sb.WriteString(fmt.Sprintf("  enabled: %t\n", c.Parallelization.Enabled))
	sb.WriteString(fmt.Sprintf("  max_workers: %d\n", c.Parallelization.MaxWorkers))
	sb.WriteString("  array_operations:\n")
	sb.WriteString(fmt.Sprintf("    enabled: %t\n", c.Parallelization.ArrayOperations.Enabled))
	sb.WriteString(fmt.Sprintf("    min_size: %d\n", c.Parallelization.ArrayOperations.MinSize))

	sb.WriteString("\nlogging:\n")
	sb.WriteString(fmt.Sprintf("  level: %s\n", c.Logging.Level))
	sb.WriteString(fmt.Sprintf("  format: %s  # text or json\n", c.Logging.Format))
	if c.Logging.Output != "" {
		sb.WriteString(fmt.Sprintf("  output: %s\n", c.Logging.Output))
	}

	sb.WriteString("\nmetrics:\n")
	sb.WriteString(fmt.Sprintf("  enabled: %t\n", c.Metrics.Enabled))
	sb.WriteString(fmt.Sprintf("  port: %d\n", c.Metrics.Port))
	sb.WriteString(fmt.Sprintf("  path: %s\n", c.Metrics.Path))

	sb.WriteString("\nhealth:\n")
	sb.WriteString(fmt.Sprintf("  enabled: %t\n", c.Health.Enabled))
	sb.WriteString(fmt.Sprintf("  port: %d\n", c.Health.Port))
	sb.WriteString(fmt.Sprintf("  path: %s\n", c.Health.Path))

	sb.WriteString("\ngo_integration:\n")
	sb.WriteString(fmt.Sprintf("  enabled: %t\n", c.GoIntegration.Enabled))
	if c.GoIntegration.PluginDirectory != "" {
		sb.WriteString(fmt.Sprintf("  plugin_directory: %s\n", c.GoIntegration.PluginDirectory))
	}
	if c.GoIntegration.ImportPath != "" {
		sb.WriteString(fmt.Sprintf("  import_path: %s\n", c.GoIntegration.ImportPath))
	}

	return sb.String()
}

// ConvertFile converts a php.ini file to PHP-Go YAML configuration
func ConvertFile(iniPath string) (*Config, []string, error) {
	converter := NewConfigConverter()

	if err := converter.ParseINI(iniPath); err != nil {
		return nil, nil, err
	}

	config, err := converter.Convert()
	if err != nil {
		return nil, nil, err
	}

	return config, converter.Warnings, nil
}
