# Symfony Testing for PHP-Go

This directory contains a Symfony 7.3 installation for testing PHP-Go compatibility with modern Symfony applications.

## Installation Summary

- **Symfony Version**: 7.3.7
- **PHP Required**: >=8.2
- **Total PHP Files**: 1,541
- **Total Lines of Code**: ~207,286
- **Total Packages**: 96 (31 direct dependencies)
- **Installation Date**: November 24, 2025

## Directory Structure

```
symfony-app/
├── bin/            - Console commands (bin/console)
├── config/         - Configuration files (YAML/PHP)
│   ├── bundles.php
│   ├── routes/
│   ├── packages/
│   └── services.yaml
├── public/         - Web root
│   └── index.php   - Front controller
├── src/            - Application code
│   └── Kernel.php  - Application kernel
├── var/            - Cache and logs (180 PHP files, auto-generated)
└── vendor/         - Dependencies (1,357 PHP files)
```

## File Statistics

| Directory | PHP Files | Lines of Code |
|-----------|-----------|---------------|
| bin       | 0         | 0             |
| config    | 2         | 10            |
| public    | 1         | 9             |
| src       | 1         | 11            |
| var       | 180       | 17,242        |
| vendor    | 1,357     | 190,014       |
| **Total** | **1,541** | **207,286**   |

## Main Symfony Components

The installation includes 27 core Symfony packages:

- **symfony/cache** - PSR-6/PSR-16 caching
- **symfony/console** - CLI applications
- **symfony/dependency-injection** - Dependency injection container
- **symfony/event-dispatcher** - Event dispatcher
- **symfony/framework-bundle** - Core framework integration
- **symfony/http-foundation** - HTTP abstraction
- **symfony/http-kernel** - HTTP kernel
- **symfony/routing** - URL routing
- **symfony/config** - Configuration loading
- **symfony/finder** - File/directory finder
- **symfony/filesystem** - File system utilities
- **symfony/var-dumper** - Variable dumping
- **symfony/yaml** - YAML parser
- Plus 14 more Symfony components

## Testing Scripts

### 1. `inventory.sh`
Generates detailed statistics about the Symfony installation.

```bash
./inventory.sh
```

Output includes:
- Version information
- File and directory counts
- Symfony package list
- Configuration file counts
- Key file verification

### 2. `test-parse.php`
Quick test of critical Symfony files using standard PHP.

```bash
php test-parse.php
```

Tests:
- `bin/console` - CLI entry point
- `public/index.php` - Web entry point
- `src/Kernel.php` - Application kernel
- `config/bundles.php` - Bundle configuration
- `config/services.yaml` - Service configuration

### 3. `run-tests.sh`
Tests critical Symfony files with PHP-Go parser.

```bash
./run-tests.sh
```

Validates that PHP-Go can successfully parse key framework files.

### 4. `identify-issues.sh`
Comprehensive analysis of all Symfony PHP files.

```bash
./identify-issues.sh
```

Generates:
- Parse success rate across all 1,541 files
- Detailed error logs
- Error pattern analysis
- Compatibility report in Markdown format

Output files created in `reports/` directory:
- `issues_<timestamp>.md` - Main compatibility report
- `parse_errors_<timestamp>.log` - Detailed parse errors
- `failed_files_<timestamp>.txt` - List of files that failed
- `error_summary_<timestamp>.txt` - Error frequency analysis

## Critical Files for Testing

These are the essential files that must parse successfully for basic Symfony functionality:

1. **public/index.php** - Front controller (web entry point)
2. **bin/console** - CLI entry point
3. **src/Kernel.php** - Application kernel
4. **config/bundles.php** - Bundle registration

## Symfony Features to Test

### PHP 8.2+ Features Used
- Constructor property promotion
- Union types
- Readonly properties
- Attributes (extensively used)
- Named arguments
- Match expressions
- Nullsafe operator (`?->`)
- First-class callables

### Symfony-Specific Patterns
- Dependency injection
- Service configuration (YAML)
- Event dispatching
- Routing
- Bundle system
- Console commands
- HTTP kernel

## Expected Challenges

Based on Symfony's architecture and PHP 8.2+ requirement:

1. **Attributes** - Symfony heavily uses PHP 8.0+ attributes for configuration
2. **Advanced Type Declarations** - Union types, intersection types, DNF types
3. **Readonly Properties** - PHP 8.1+ feature used extensively
4. **Enums** - PHP 8.1+ enums in routing and configuration
5. **First-Class Callables** - `$fn(...)` syntax (PHP 8.1+)
6. **Constructor Property Promotion** - Modern PHP syntax
7. **Named Arguments** - Used in service configuration

## Comparison with Laravel

| Metric | Symfony 7.3 | Laravel 12 |
|--------|-------------|------------|
| PHP Files | 1,541 | 7,482 |
| Lines of Code | ~207k | ~930k |
| PHP Requirement | >=8.2 | >=8.2 |
| Modern Features | Very Heavy | Heavy |
| Framework Focus | Components | Full-Stack |

Symfony is more compact but uses more modern PHP features intensively.

## Testing Workflow

1. **Build PHP-Go**:
   ```bash
   cd ../..
   go build -o php-go ./cmd/php-go
   ```

2. **Run inventory**:
   ```bash
   ./inventory.sh
   ```

3. **Quick parse test**:
   ```bash
   ./run-tests.sh
   ```

4. **Full analysis**:
   ```bash
   ./identify-issues.sh
   ```

5. **Review results**:
   ```bash
   cat reports/issues_*.md
   ```

## Next Steps

1. Run parse tests to establish baseline compatibility
2. Identify top parsing errors from error analysis
3. Implement missing language features in PHP-Go
4. Re-test and measure improvement
5. Focus on critical framework files first
6. Expand to full vendor directory testing

## Notes

- Symfony skeleton is minimal - contains only core components
- Most code is in `vendor/` (190k lines across 1,357 files)
- Application code in `src/` is minimal (11 lines) - standard for skeleton
- Generated cache files in `var/` (180 files, 17k lines) are auto-created
- All configuration can be in YAML, PHP, or XML - this installation uses YAML primarily

## Resources

- [Symfony Documentation](https://symfony.com/doc)
- [Symfony GitHub](https://github.com/symfony/symfony)
- [PHP 8.2 Features](https://www.php.net/releases/8.2/en.php)
- [PHP-Go Project](../../README.md)
