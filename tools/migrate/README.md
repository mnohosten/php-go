# PHP-Go Migration Tools

Tools to help migrate PHP applications to PHP-Go.

## Compatibility Analyzer

The compatibility analyzer scans PHP applications and generates a compatibility report identifying:

- Parse errors and syntax issues
- Unsupported PHP extensions
- Unsupported functions
- Language features that may not be supported
- Categorized issues with severity levels and suggestions

### Installation

Build the tool:

```bash
go build -o php-go-analyzer ./cmd/php-go-analyzer/
```

Or install it globally:

```bash
go install ./cmd/php-go-analyzer/
```

### Usage

```bash
# Analyze current directory
php-go-analyzer

# Analyze specific directory
php-go-analyzer -path /path/to/project

# Save report to file
php-go-analyzer -path /path/to/project -output report.txt

# Exclude patterns (default: vendor,node_modules,.git)
php-go-analyzer -exclude vendor,tests,cache

# Show help
php-go-analyzer -help
```

### Example Output

```
===========================================
PHP-Go Compatibility Analysis Report
===========================================

SUMMARY
-------
Total Files Analyzed:     1234
Successfully Parsed:      1150 (93.20%)
Parse Failures:           84 (6.80%)
Total Issues Found:       245
Unique Extensions Used:   12
Unsupported Extensions:   5
Unsupported Functions:    32

ISSUE SEVERITY BREAKDOWN
------------------------
Critical:  0
Error:     45
Warning:   150
Info:      50

UNSUPPORTED EXTENSIONS
----------------------
  mysqli               120 usage(s)
  curl                 45 usage(s)
  gd                   23 usage(s)
  mbstring             15 usage(s)
  xml                  12 usage(s)

TOP 10 UNSUPPORTED FUNCTIONS
----------------------------
  mysqli_connect                 45 usage(s)
  mysqli_query                   38 usage(s)
  curl_init                      25 usage(s)
  imagecreatetruecolor           15 usage(s)
  mb_strlen                      12 usage(s)

LANGUAGE FEATURES USED
----------------------
  ✓ enums                        15 usage(s)
  ✓ readonly                     8 usage(s)
  ✗ generators                   5 usage(s) [UNSUPPORTED]
  ✗ attributes                   3 usage(s) [UNSUPPORTED]

RECOMMENDATIONS
---------------
✓ GOOD: Parse success rate is acceptable (93.20%).
⚠ WARNING: 5 unsupported extensions detected.
  Action: Review extension usage and plan for alternatives or custom implementations.
⚠ WARNING: 245 compatibility issues found.
  Action: Review and prioritize issues before migration.

For detailed migration guidance, see: docs/user-guide/migration.md
===========================================
```

### Detected Extensions

The analyzer can detect usage of the following extensions:

**Database**:
- mysqli
- pdo
- pgsql
- sqlite3

**XML**:
- simplexml
- xml
- dom

**Images**:
- gd
- imagick

**Networking**:
- curl
- sockets
- ftp

**Compression**:
- zlib
- bz2
- zip

**Caching**:
- apcu
- opcache
- memcache

**Encryption**:
- openssl
- mcrypt

**String Handling**:
- mbstring
- iconv

**Other**:
- ctype
- filter
- intl

### Language Features

The analyzer detects the following PHP language features:

**Supported** (✓):
- Enums (PHP 8.1+)
- Readonly properties (PHP 8.1+)
- Named arguments (PHP 8.0+)

**Not Yet Supported** (✗):
- Generators/yield (planned for Phase 9)
- Arrow functions (planned for Phase 9)
- Attributes (planned for Phase 9)
- Match expressions (planned for Phase 9)
- Throw expressions (planned for Phase 9)

### Issue Severity Levels

- **CRITICAL**: Blocks migration completely
- **ERROR**: Must be fixed before migration
- **WARNING**: Should be reviewed, may need attention
- **INFO**: Informational, no action required

### Integration with Migration Workflow

1. **Pre-Migration Assessment**:
   ```bash
   php-go-analyzer -path /path/to/project -output assessment.txt
   ```

2. **Review Report**: Examine the compatibility report and identify blockers

3. **Plan Migration**:
   - Prioritize fixing CRITICAL and ERROR issues
   - Plan alternatives for unsupported extensions
   - Consider refactoring unsupported language features

4. **Iterative Testing**:
   - Fix issues
   - Re-run analyzer
   - Track progress

5. **Final Validation**: Run analyzer one more time before production migration

### Testing

Run the test suite:

```bash
go test ./tools/migrate/ -v
```

Run benchmarks:

```bash
go test ./tools/migrate/ -bench=. -benchmem
```

### Development

The analyzer consists of:

- **analyzer.go**: Core analyzer implementation
  - File scanning and parsing
  - Extension and function detection
  - Language feature detection
  - Report generation

- **analyzer_test.go**: Comprehensive test suite
  - 17 test cases
  - 2 benchmarks
  - Tests for all major functionality

- **main.go** (in cmd/php-go-analyzer/): CLI tool
  - Command-line argument parsing
  - User-friendly output
  - File and stdout output options

### Future Enhancements

Planned improvements:

- JSON output format for CI/CD integration
- HTML report generation
- Issue filtering by severity
- Custom rule definitions
- Integration with popular PHP frameworks
- Automatic fix suggestions
- Progressive migration tracking

### Related Documentation

- [Migration Guide](../../docs/user-guide/migration.md)
- [Configuration Guide](../../docs/user-guide/configuration.md)
- [Extension Development Guide](../../docs/user-guide/extension-development.md)

## Configuration Converter

The configuration converter converts php.ini files to PHP-Go YAML configuration format, preserving all PHP settings while adding PHP-Go specific features.

### Installation

Build the tool:

```bash
go build -o php-go-config ./cmd/php-go-config/
```

Or install it globally:

```bash
go install ./cmd/php-go-config/
```

### Usage

```bash
# Convert php.ini to stdout
php-go-config -input /etc/php.ini

# Convert and save to file
php-go-config -input /etc/php.ini -output php-go.yaml

# Convert from current directory
php-go-config

# Disable warnings
php-go-config -input php.ini -warnings=false

# Show help
php-go-config -help
```

### Features

**PHP Settings Conversion:**
- Error reporting (converts numeric values to constants)
- Display/log errors
- Memory limits (normalizes sizes: bytes → K/M/G)
- Execution time limits
- Upload/POST size limits
- Timezone configuration

**Extension Management:**
- Automatically collects all `extension=` directives
- Removes .so/.dll suffixes
- Detects and warns about unsupported `zend_extension` directives
- Deduplicates extensions

**PHP-Go Specific Settings:**
- Adds parallelization configuration (disabled by default)
- Adds logging configuration
- Adds metrics endpoint configuration
- Adds health check configuration
- Adds Go integration settings

### Example

Input `php.ini`:
```ini
[PHP]
error_reporting = 32767
display_errors = On
memory_limit = 256M
max_execution_time = 60
date.timezone = UTC

extension=json
extension=hash
extension=mysqli
```

Output `php-go.yaml`:
```yaml
# PHP-Go Configuration
# Generated from php.ini

php:
  error_reporting: E_ALL
  display_errors: true
  log_errors: true
  memory_limit: 256M
  max_execution_time: 60
  post_max_size: 8M
  upload_max_filesize: 2M

  date:
    timezone: UTC

extensions:
  enabled:
    - json
    - hash
    - mysqli

# PHP-Go specific features (disabled by default for compatibility)
parallelization:
  enabled: false
  max_workers: 4
  array_operations:
    enabled: false
    min_size: 1000

logging:
  level: info
  format: text

metrics:
  enabled: false
  port: 9090
  path: /metrics

health:
  enabled: false
  port: 8080
  path: /health

go_integration:
  enabled: false
```

### Testing

Run the test suite:

```bash
go test ./tools/migrate/ -v -run TestConvert
```

Check coverage:

```bash
go test -cover ./tools/migrate/
```

### Integration with Migration Workflow

1. **Pre-Migration Configuration Assessment**:
   ```bash
   php-go-config -input /etc/php/php.ini -output php-go-initial.yaml
   ```

2. **Review Generated Config**: Examine the output and adjust PHP-Go specific settings

3. **Enable PHP-Go Features**: Update the YAML to enable parallelization, metrics, etc.

4. **Test Configuration**:
   ```bash
   php-go -c php-go-initial.yaml your-script.php
   ```

5. **Iterate**: Adjust settings based on performance and compatibility testing

### Development

The converter consists of:

- **converter.go**: Core converter implementation
  - INI parsing with section support
  - Setting conversion and normalization
  - YAML generation
  - Extension tracking

- **converter_test.go**: Comprehensive test suite
  - 13 test functions
  - Tests for INI parsing, conversion logic, and YAML output
  - 83%+ code coverage

- **main.go** (in cmd/php-go-config/): CLI tool
  - Command-line argument parsing
  - File I/O handling
  - User-friendly output

### Future Enhancements

Planned improvements:

- INI file detection (auto-find php.ini)
- Batch conversion of multiple files
- Configuration validation
- Diff mode (compare php.ini with php-go.yaml)
- Interactive mode with prompts
- Configuration templates for common setups
- Environment variable substitution

## Automated Migration Tool

The automated migration tool (`php-go-migrate-auto`) orchestrates the entire migration process by running all the tools above in sequence and generating a comprehensive migration report.

### Installation

Build the tool:

```bash
go build -o php-go-migrate-auto ./cmd/php-go-migrate-auto/
```

Or install it globally:

```bash
go install ./cmd/php-go-migrate-auto/
```

### Usage

```bash
# Analyze current directory
php-go-migrate-auto

# Analyze specific project with custom name
php-go-migrate-auto -path /path/to/project -name MyProject

# Enable verbose output
php-go-migrate-auto -path /path/to/project -verbose

# Show help
php-go-migrate-auto -help
```

### What It Does

The automated migration tool performs the following steps:

1. **Creates Output Directory**: Creates `.php-go-migration/` in the project for all artifacts
2. **Runs Compatibility Analysis**: Scans all PHP files and generates compatibility report
3. **Converts Configuration**: Converts `php.ini` to PHP-Go YAML format (if found)
4. **Creates Migration Checklist**: Initializes a 93-task migration tracking checklist
5. **Tests Parse Success**: Tests PHP file parsing with `php-go` binary (if available)
6. **Generates Reports**: Creates comprehensive migration report in Markdown and JSON

### Output Files

All migration artifacts are saved to `<project>/.php-go-migration/`:

- **compatibility-report.txt**: Detailed compatibility analysis with issues and recommendations
- **migration-checklist.json**: Task tracking checklist for the migration process
- **php-go.yaml**: Converted PHP-Go configuration (if `php.ini` found)
- **parse-failures.txt**: List of files that failed to parse (if parse test ran)
- **migration-report.md**: Comprehensive migration report with summary and next steps
- **migration-report.json**: Machine-readable report for CI/CD integration

### Example Output

```
[INFO] Starting automated migration for project: MyProject
[INFO] Project path: /path/to/project
[INFO] Output directory: /path/to/project/.php-go-migration

[INFO] Running step: Create Output Directory
  Description: Create directory for migration artifacts
  Step completed in 1.2ms

[INFO] Running step: Run Compatibility Analyzer
  Description: Analyze codebase for compatibility issues
  Analyzing project for compatibility issues...
  Analysis complete:
    Total files: 1234
    Parse success rate: 93.20%
    Issues found: 245
    Report saved to: /path/to/project/.php-go-migration/compatibility-report.txt
  Step completed in 2.3s

[INFO] Running step: Convert Configuration
  Description: Convert php.ini to PHP-Go YAML format
  Converting configuration from: /etc/php.ini
  Configuration converted and saved to: /path/to/project/.php-go-migration/php-go.yaml
  Step completed in 45ms

[INFO] Running step: Create Migration Checklist
  Description: Initialize migration tracking checklist
  Creating migration checklist...
  Checklist created with 93 items
  Saved to: /path/to/project/.php-go-migration/migration-checklist.json
  Step completed in 12ms

[INFO] Running step: Test Parse Success Rate
  Description: Parse all PHP files and report success rate
  Testing parse success rate...
  Found 1234 PHP files to test
  Parse test complete:
    Success: 1150/1234 (93.20%)
    Failures: 84
  Parse failures saved to: /path/to/project/.php-go-migration/parse-failures.txt
  Step completed in 4.5s

[INFO] Running step: Generate Migration Report
  Description: Generate comprehensive migration report
  Generating migration report...
  JSON report saved to: /path/to/project/.php-go-migration/migration-report.json
  Migration report saved to: /path/to/project/.php-go-migration/migration-report.md
  Step completed in 23ms

[INFO] Migration process completed successfully!

═══════════════════════════════════════════════════
           MIGRATION PROCESS COMPLETE
═══════════════════════════════════════════════════

Steps Completed: 6/6

Parse Success Rate: 93.20%
Issues Found: 245

All migration artifacts saved to: /path/to/project/.php-go-migration

Next steps:
  1. Review the migration report: /path/to/project/.php-go-migration/migration-report.md
  2. Check migration progress: php-go-migrate status -file /path/to/project/.php-go-migration/migration-checklist.json
  3. Consult the migration guide: docs/user-guide/migration.md
```

### Generated Migration Report

The tool generates a comprehensive Markdown report (`migration-report.md`):

```markdown
# PHP-Go Migration Report

**Project**: MyProject
**Date**: 2025-11-24 10:30:00

## Summary

### Compatibility Analysis

- **Total Files**: 1234
- **Parse Success Rate**: 93.20%
- **Issues Found**: 245
- **Unsupported Extensions**: 5
- **Unsupported Functions**: 32
- **Report**: /path/to/project/.php-go-migration/compatibility-report.txt

### Parse Test Results

- **Total Files Tested**: 1234
- **Success Rate**: 93.20%
- **Successful**: 1150
- **Failed**: 84

### Configuration

- **Source**: /etc/php.ini
- **Destination**: /path/to/project/.php-go-migration/php-go.yaml

### Migration Checklist

- **Total Tasks**: 93
- **Checklist File**: /path/to/project/.php-go-migration/migration-checklist.json

## Next Steps

1. Review the compatibility report to understand migration challenges
2. Check the parse test results to identify files that need fixes
3. Review the generated configuration file and adjust as needed
4. Use the migration checklist to track progress:
   ```
   php-go-migrate status -file /path/to/project/.php-go-migration/migration-checklist.json
   ```
5. Address critical issues identified in the compatibility report
6. Run functional tests to verify application behavior
7. Consult the migration guide: `docs/user-guide/migration.md`
```

### Integration with Migration Workflow

The automated migration tool is designed to be the **first step** in your migration journey:

```bash
# Step 1: Run automated analysis
php-go-migrate-auto -path /path/to/project

# Step 2: Review the migration report
cat .php-go-migration/migration-report.md

# Step 3: Track progress with the checklist
php-go-migrate status -file .php-go-migration/migration-checklist.json

# Step 4: Mark tasks as complete
php-go-migrate complete pre-migration-1 -file .php-go-migration/migration-checklist.json

# Step 5: Generate progress report
php-go-migrate report -format markdown -file .php-go-migration/migration-checklist.json
```

### Testing

Run the test suite for the orchestrator:

```bash
go test ./tools/migrate/ -v -run Orchestrator
```

Run all migration tool tests:

```bash
go test ./tools/migrate/ -v
```

### Development

The orchestrator consists of:

- **orchestrator.go**: Core orchestration logic
  - Step management
  - Tool integration
  - Report generation
  - Error handling

- **orchestrator_test.go**: Comprehensive test suite
  - 18+ test functions
  - Tests for all major functionality
  - Integration tests

- **main.go** (in cmd/php-go-migrate-auto/): CLI tool
  - Command-line argument parsing
  - User-friendly output
  - Error handling

### Future Enhancements

Planned improvements:

- Interactive mode with prompts for user decisions
- CI/CD integration hooks
- Custom step plugins
- Parallel execution of independent steps
- Incremental migration support (resume from last step)
- Migration dry-run mode
- Automatic fix suggestions for common issues
- Integration with version control systems

## Migration Checklist Tool

The migration checklist tool (`php-go-migrate`) helps you track progress through the 93-task migration process.

See the full documentation in the main README section above for details on:
- Creating checklists (`init` command)
- Tracking progress (`status` command)
- Managing tasks (`complete`, `incomplete`, `note` commands)
- Generating reports (`report` command)

## License

Part of the PHP-Go project.
