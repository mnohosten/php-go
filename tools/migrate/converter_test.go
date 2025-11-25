package migrate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseINI(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected map[string]map[string]string
	}{
		{
			name: "basic settings",
			content: `[PHP]
memory_limit = 128M
max_execution_time = 30
display_errors = On`,
			expected: map[string]map[string]string{
				"PHP": {
					"memory_limit":       "128M",
					"max_execution_time": "30",
					"display_errors":     "On",
				},
			},
		},
		{
			name: "multiple sections",
			content: `[PHP]
memory_limit = 256M

[Date]
date.timezone = UTC`,
			expected: map[string]map[string]string{
				"PHP": {
					"memory_limit": "256M",
				},
				"Date": {
					"date.timezone": "UTC",
				},
			},
		},
		{
			name: "comments and empty lines",
			content: `; This is a comment
[PHP]

# Another comment
memory_limit = 512M
; Inline comment after setting
display_errors = Off`,
			expected: map[string]map[string]string{
				"PHP": {
					"memory_limit":   "512M",
					"display_errors": "Off",
				},
			},
		},
		{
			name: "quoted values",
			content: `[PHP]
error_log = "/var/log/php.log"
date.timezone = 'America/New_York'`,
			expected: map[string]map[string]string{
				"PHP": {
					"error_log":      "/var/log/php.log",
					"date.timezone":  "America/New_York",
				},
			},
		},
		{
			name: "extensions",
			content: `[PHP]
extension=json
extension=hash.so
extension = mysqli`,
			expected: map[string]map[string]string{
				"PHP": {
					"extension": "mysqli", // Last one wins in the simple map
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpFile := filepath.Join(t.TempDir(), "test.ini")
			if err := os.WriteFile(tmpFile, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}

			// Parse the file
			converter := NewConfigConverter()
			if err := converter.ParseINI(tmpFile); err != nil {
				t.Fatalf("ParseINI failed: %v", err)
			}

			// Verify sections exist
			for section := range tt.expected {
				if _, ok := converter.Input[section]; !ok {
					t.Errorf("Expected section [%s] not found", section)
				}
			}

			// Verify key values (note: extension= keys are special, only check non-extension keys)
			for section, settings := range tt.expected {
				for key, expectedValue := range settings {
					if key == "extension" {
						continue // Skip extension checks in this simple test
					}
					if actualValue, ok := converter.Input[section][key]; !ok {
						t.Errorf("Expected key %s in section [%s] not found", key, section)
					} else if actualValue != expectedValue {
						t.Errorf("Section [%s], key %s: expected %q, got %q", section, key, expectedValue, actualValue)
					}
				}
			}
		})
	}
}

func TestConvertBool(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"On", true},
		{"on", true},
		{"ON", true},
		{"Off", false},
		{"off", false},
		{"1", true},
		{"0", false},
		{"true", true},
		{"false", false},
		{"yes", true},
		{"no", false},
		{"", false},
		{"invalid", false},
	}

	converter := NewConfigConverter()
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := converter.convertBool(tt.input)
			if result != tt.expected {
				t.Errorf("convertBool(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestConvertErrorReporting(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"E_ALL", "E_ALL"},
		{"E_ALL & ~E_NOTICE", "E_ALL & ~E_NOTICE"},
		{"32767", "E_ALL"},
		{"30719", "E_ALL & ~E_DEPRECATED & ~E_STRICT"},
		{"22527", "E_ALL & ~E_NOTICE"},
		{"0", "0"},
		{"12345", "12345"},
	}

	converter := NewConfigConverter()
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := converter.convertErrorReporting(tt.input)
			if result != tt.expected {
				t.Errorf("convertErrorReporting(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNormalizeMemorySize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"128M", "128M"},
		{"256m", "256M"},
		{"1G", "1G"},
		{"1024K", "1024K"},
		{"134217728", "128M"},  // 128MB in bytes
		{"1073741824", "1G"},   // 1GB in bytes
		{"1024", "1K"},         // 1KB in bytes
		{"-1", "-1"},           // Unlimited
		{"0", "0"},
	}

	converter := NewConfigConverter()
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := converter.normalizeMemorySize(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeMemorySize(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestConvert(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		validate func(*testing.T, *Config)
	}{
		{
			name: "basic PHP settings",
			content: `[PHP]
memory_limit = 256M
max_execution_time = 60
display_errors = On
log_errors = On
error_log = /var/log/php.log
error_reporting = E_ALL
date.timezone = America/New_York`,
			validate: func(t *testing.T, cfg *Config) {
				if cfg.PHP.MemoryLimit != "256M" {
					t.Errorf("Expected memory_limit 256M, got %s", cfg.PHP.MemoryLimit)
				}
				if cfg.PHP.MaxExecutionTime != 60 {
					t.Errorf("Expected max_execution_time 60, got %d", cfg.PHP.MaxExecutionTime)
				}
				if !cfg.PHP.DisplayErrors {
					t.Error("Expected display_errors to be true")
				}
				if !cfg.PHP.LogErrors {
					t.Error("Expected log_errors to be true")
				}
				if cfg.PHP.ErrorLog != "/var/log/php.log" {
					t.Errorf("Expected error_log /var/log/php.log, got %s", cfg.PHP.ErrorLog)
				}
				if cfg.PHP.ErrorReporting != "E_ALL" {
					t.Errorf("Expected error_reporting E_ALL, got %s", cfg.PHP.ErrorReporting)
				}
				if cfg.PHP.Date.Timezone != "America/New_York" {
					t.Errorf("Expected timezone America/New_York, got %s", cfg.PHP.Date.Timezone)
				}
			},
		},
		{
			name: "extensions",
			content: `[PHP]
extension=json
extension=hash
extension=mysqli.so`,
			validate: func(t *testing.T, cfg *Config) {
				if len(cfg.Extensions.Enabled) == 0 {
					t.Error("Expected extensions to be enabled")
				}
				// Check that extensions are present (order may vary)
				hasJson := false
				hasHash := false
				hasMysqli := false
				for _, ext := range cfg.Extensions.Enabled {
					switch ext {
					case "json":
						hasJson = true
					case "hash":
						hasHash = true
					case "mysqli":
						hasMysqli = true
					}
				}
				if !hasJson {
					t.Error("Expected json extension")
				}
				if !hasHash {
					t.Error("Expected hash extension")
				}
				if !hasMysqli {
					t.Error("Expected mysqli extension")
				}
			},
		},
		{
			name: "error reporting as number",
			content: `[PHP]
error_reporting = 32767`,
			validate: func(t *testing.T, cfg *Config) {
				if cfg.PHP.ErrorReporting != "E_ALL" {
					t.Errorf("Expected error_reporting E_ALL, got %s", cfg.PHP.ErrorReporting)
				}
			},
		},
		{
			name: "memory size conversion",
			content: `[PHP]
memory_limit = 134217728`,
			validate: func(t *testing.T, cfg *Config) {
				if cfg.PHP.MemoryLimit != "128M" {
					t.Errorf("Expected memory_limit 128M, got %s", cfg.PHP.MemoryLimit)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpFile := filepath.Join(t.TempDir(), "test.ini")
			if err := os.WriteFile(tmpFile, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}

			// Convert the file
			config, _, err := ConvertFile(tmpFile)
			if err != nil {
				t.Fatalf("ConvertFile failed: %v", err)
			}

			// Validate the result
			tt.validate(t, config)
		})
	}
}

func TestToYAML(t *testing.T) {
	config := &Config{
		PHP: PHPSettings{
			ErrorReporting:    "E_ALL",
			DisplayErrors:     true,
			LogErrors:         true,
			ErrorLog:          "/var/log/php.log",
			MemoryLimit:       "256M",
			MaxExecutionTime:  30,
			PostMaxSize:       "8M",
			UploadMaxFilesize: "2M",
			Date: DateSettings{
				Timezone: "UTC",
			},
		},
		Extensions: ExtensionSettings{
			Enabled: []string{"json", "hash"},
			Paths:   []string{"/usr/local/lib/php-go/extensions"},
		},
		Parallelization: ParallelSettings{
			Enabled:    true,
			MaxWorkers: 8,
			ArrayOperations: ArrayParallelSettings{
				Enabled: true,
				MinSize: 1000,
			},
		},
		Logging: LogSettings{
			Level:  "info",
			Format: "json",
			Output: "/var/log/php-go.log",
		},
		Metrics: MetricsSettings{
			Enabled: true,
			Port:    9090,
			Path:    "/metrics",
		},
		Health: HealthSettings{
			Enabled: true,
			Port:    8080,
			Path:    "/health",
		},
		GoIntegration: GoIntegrationSettings{
			Enabled:         true,
			PluginDirectory: "/usr/local/lib/php-go/plugins",
			ImportPath:      "/usr/local/lib/php-go/imports",
		},
	}

	yaml := config.ToYAML()

	// Check that YAML contains expected values
	expectedStrings := []string{
		"php:",
		"error_reporting: E_ALL",
		"display_errors: true",
		"log_errors: true",
		"error_log: /var/log/php.log",
		"memory_limit: 256M",
		"max_execution_time: 30",
		"timezone: UTC",
		"extensions:",
		"- json",
		"- hash",
		"parallelization:",
		"enabled: true",
		"max_workers: 8",
		"logging:",
		"level: info",
		"format: json",
		"metrics:",
		"port: 9090",
		"path: /metrics",
		"health:",
		"port: 8080",
		"path: /health",
		"go_integration:",
		"plugin_directory: /usr/local/lib/php-go/plugins",
		"import_path: /usr/local/lib/php-go/imports",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(yaml, expected) {
			t.Errorf("Expected YAML to contain %q", expected)
		}
	}
}

func TestProcessExtensions(t *testing.T) {
	converter := NewConfigConverter()

	// Simulate extensions collected during parsing
	converter.Extensions = []string{"json", "hash", "mysqli"}

	ext := &ExtensionSettings{}
	converter.processExtensions(ext)

	if len(ext.Enabled) < 2 {
		t.Errorf("Expected at least 2 extensions, got %d", len(ext.Enabled))
	}

	// Check that extensions are present
	hasJson := false
	hasHash := false
	hasMysqli := false
	for _, e := range ext.Enabled {
		switch e {
		case "json":
			hasJson = true
		case "hash":
			hasHash = true
		case "mysqli":
			hasMysqli = true
		}

		// Check that .so suffix is removed (shouldn't happen with our test data)
		if strings.HasSuffix(e, ".so") {
			t.Errorf("Extension %q should not have .so suffix", e)
		}
		if strings.HasSuffix(e, ".dll") {
			t.Errorf("Extension %q should not have .dll suffix", e)
		}
	}

	if !hasJson || !hasHash || !hasMysqli {
		t.Errorf("Expected json, hash, and mysqli extensions. Got: %v", ext.Enabled)
	}
}

func TestConvertFile(t *testing.T) {
	// Create a comprehensive php.ini file
	content := `; PHP Configuration

[PHP]
; Error handling
error_reporting = E_ALL
display_errors = On
log_errors = On
error_log = /var/log/php.log

; Resource limits
memory_limit = 256M
max_execution_time = 60
post_max_size = 64M
upload_max_filesize = 32M

; Date/time
date.timezone = UTC

; Extensions
extension=json
extension=hash
extension=mysqli

[Date]
date.timezone = America/Los_Angeles  ; This should override

[Session]
session.save_handler = files
session.save_path = /tmp
`

	tmpFile := filepath.Join(t.TempDir(), "php.ini")
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	config, warnings, err := ConvertFile(tmpFile)
	if err != nil {
		t.Fatalf("ConvertFile failed: %v", err)
	}

	// Verify configuration
	if config.PHP.MemoryLimit != "256M" {
		t.Errorf("Expected memory_limit 256M, got %s", config.PHP.MemoryLimit)
	}

	if config.PHP.MaxExecutionTime != 60 {
		t.Errorf("Expected max_execution_time 60, got %d", config.PHP.MaxExecutionTime)
	}

	if !config.PHP.DisplayErrors {
		t.Error("Expected display_errors to be true")
	}

	if len(config.Extensions.Enabled) < 2 {
		t.Errorf("Expected at least 2 extensions, got %d", len(config.Extensions.Enabled))
	}

	// Check timezone (Date section should override PHP section)
	if config.PHP.Date.Timezone != "America/Los_Angeles" {
		t.Errorf("Expected timezone America/Los_Angeles, got %s", config.PHP.Date.Timezone)
	}

	// Warnings can be empty, but should be a slice
	if warnings == nil {
		t.Error("Expected warnings slice, got nil")
	}
}

func TestDefaultConfiguration(t *testing.T) {
	// Test with minimal php.ini
	content := `[PHP]
; Just one setting
display_errors = Off
`

	tmpFile := filepath.Join(t.TempDir(), "php.ini")
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	config, _, err := ConvertFile(tmpFile)
	if err != nil {
		t.Fatalf("ConvertFile failed: %v", err)
	}

	// Check that defaults are set
	if config.PHP.MemoryLimit == "" {
		t.Error("Expected default memory_limit to be set")
	}

	if config.PHP.MaxExecutionTime == 0 {
		t.Error("Expected default max_execution_time to be set")
	}

	if config.PHP.ErrorReporting == "" {
		t.Error("Expected default error_reporting to be set")
	}

	// PHP-Go specific settings should have defaults
	if config.Parallelization.MaxWorkers == 0 {
		t.Error("Expected default max_workers to be set")
	}

	if config.Metrics.Port == 0 {
		t.Error("Expected default metrics port to be set")
	}

	if config.Health.Port == 0 {
		t.Error("Expected default health port to be set")
	}
}

func TestZendExtensionWarning(t *testing.T) {
	content := `[PHP]
zend_extension=opcache
zend_extension=/path/to/xdebug.so
`

	tmpFile := filepath.Join(t.TempDir(), "php.ini")
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	_, warnings, err := ConvertFile(tmpFile)
	if err != nil {
		t.Fatalf("ConvertFile failed: %v", err)
	}

	// Should have warnings about zend_extension
	if len(warnings) == 0 {
		t.Error("Expected warnings about zend_extension")
	}

	hasZendWarning := false
	for _, warning := range warnings {
		if strings.Contains(warning, "zend_extension") {
			hasZendWarning = true
			break
		}
	}

	if !hasZendWarning {
		t.Error("Expected warning about zend_extension not supported")
	}
}
