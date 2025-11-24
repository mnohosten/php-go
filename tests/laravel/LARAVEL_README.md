# Laravel Testing for PHP-Go

This directory contains a Laravel installation for testing PHP-Go's compatibility with the Laravel framework.

## Installation Details

**Installation Date**: November 24, 2025
**Laravel Version**: v12.10.1 (Laravel Framework v12.39.0)
**PHP Version Required**: ^8.2
**Total PHP Files**: 7,482
**Total Lines of PHP Code**: ~441,000

## Installed Packages

### Production Dependencies
- `laravel/framework`: ^12.0 (v12.39.0)
- `laravel/tinker`: ^2.10.1

### Development Dependencies
- `fakerphp/faker`: ^1.23
- `laravel/pail`: ^1.2.2
- `laravel/pint`: ^1.24
- `laravel/sail`: ^1.41
- `mockery/mockery`: ^1.6
- `nunomaduro/collision`: ^8.6
- `phpunit/phpunit`: ^11.5.3

Total packages installed: 111 (including all dependencies)

## Directory Structure

```
tests/laravel/
├── app/              # Application source code
├── bootstrap/        # Framework bootstrap files
├── config/           # Configuration files
├── database/         # Database migrations, factories, and seeders
│   └── database.sqlite  # SQLite database (created during setup)
├── public/           # Public web root
├── resources/        # Views, raw assets
├── routes/           # Route definitions
├── storage/          # Compiled templates, logs, cache
├── tests/            # Test files
├── vendor/           # Composer dependencies
├── .env              # Environment configuration
├── artisan           # CLI tool
└── composer.json     # Composer configuration
```

## Setup Status

The following setup steps have been completed automatically:
- ✅ Composer dependencies installed
- ✅ `.env` file created from `.env.example`
- ✅ Application key generated
- ✅ SQLite database created (`database/database.sqlite`)
- ✅ Database migrations run (users, cache, jobs tables created)

## Running with Standard PHP

To verify the Laravel installation works with standard PHP:

```bash
cd tests/laravel

# Start the development server
php artisan serve

# Run tests
php artisan test

# Access the application at http://localhost:8000
```

## Testing with PHP-Go

### Prerequisites

Before testing with PHP-Go, ensure the following are implemented:

**Critical Requirements**:
- [ ] Composer autoloader support (PSR-4)
- [ ] Namespace/use statements (✅ Already implemented)
- [ ] Class autoloading
- [ ] Database support (SQLite/MySQL/PostgreSQL)
- [ ] Session handling
- [ ] File system operations
- [ ] HTTP request/response handling

**Standard Library Functions** (estimate ~400+ needed):
- String functions (str_*, mb_*)
- Array functions (array_*)
- File I/O functions (file_*, fopen, fread, etc.)
- JSON functions (✅ Already implemented)
- Hash functions (✅ Partially implemented)
- Regular expressions (preg_*)
- Date/Time functions (date, strtotime, DateTime class)
- URL functions (parse_url, http_build_query, etc.)
- Encryption functions (openssl_*)

**Extensions**:
- PDO/SQLite (CRITICAL)
- OpenSSL
- Ctype (✅ Implemented)
- Mbstring
- Session
- Reflection (✅ Partially implemented)

### Parse Testing

To test if PHP-Go can parse Laravel files:

```bash
# Build PHP-Go
go build -o php-go ./cmd/php-go

# Create a parse test script
cat > test-parse.php << 'EOF'
<?php
$dir = 'tests/laravel';
$files = new RecursiveIteratorIterator(
    new RecursiveDirectoryIterator($dir)
);

$total = 0;
$success = 0;
$errors = [];

foreach ($files as $file) {
    if ($file->isFile() && pathinfo($file, PATHINFO_EXTENSION) === 'php') {
        $total++;
        $path = $file->getPathname();

        // Try to parse with PHP-Go
        exec("./php-go parse " . escapeshellarg($path) . " 2>&1", $output, $code);

        if ($code === 0) {
            $success++;
        } else {
            $errors[$path] = implode("\n", $output);
        }
    }
}

echo "Parse Results:\n";
echo "Total files: $total\n";
echo "Successful: $success (" . round($success/$total*100, 2) . "%)\n";
echo "Failed: " . ($total - $success) . "\n";

// Save detailed errors
file_put_contents('parse-errors.log', json_encode($errors, JSON_PRETTY_PRINT));
EOF

php test-parse.php
```

### Execution Testing

Once parser and runtime support are complete:

```bash
# Run Laravel with PHP-Go
./php-go tests/laravel/artisan serve

# Run tests with PHP-Go
./php-go tests/laravel/artisan test
```

## Expected Challenges

Based on the WordPress testing experience, expect the following issues:

### Parser Issues (Priority: CRITICAL)
1. Advanced PHP features (match expressions, arrow functions, etc.)
2. Alternative syntax (if/endif, foreach/endforeach) - ✅ Already fixed
3. Inline HTML/PHP template mixing - ✅ Already fixed
4. Namespace and use statements - ✅ Already implemented

### Runtime Issues (Priority: HIGH)
1. **Autoloading**: Laravel uses PSR-4 autoloading extensively
2. **Service Container**: Dependency injection and service resolution
3. **Routing**: URL routing and middleware
4. **Eloquent ORM**: Database access and query building
5. **Blade Templates**: Template compilation and rendering
6. **Facades**: Static proxy pattern for services

### Missing Features (Priority: MEDIUM-HIGH)
1. **Reflection**: Heavily used for dependency injection
2. **Attributes**: PHP 8.0+ attributes used throughout
3. **Generators**: Used for lazy loading and iterators
4. **Closures**: Extensive use in routing, queries, middleware
5. **Anonymous classes**: Used in various places

### Extension Requirements
- PDO/Database drivers (CRITICAL - Laravel cannot run without it)
- OpenSSL (for encryption, hashing)
- Session (for web application state)
- Mbstring (for multi-byte string handling)
- Ctype (for character type checking)
- JSON (✅ Already implemented)
- Hash (✅ Partially implemented)

## Testing Strategy

### Phase 1: Parse Testing (Current)
- Test parsing of all 7,482 PHP files
- Identify and fix parser bugs
- Target: 95%+ parse success rate

### Phase 2: Basic Execution
- Implement autoloading support
- Run simple artisan commands
- Test basic route definitions

### Phase 3: Framework Features
- Database connection and queries
- Middleware execution
- Controller dispatch
- View rendering

### Phase 4: Full Test Suite
- Run Laravel's PHPUnit tests
- Fix discovered runtime issues
- Performance optimization

## Performance Targets

Based on PHP 8.4 benchmarks:
- **Cold start**: < 100ms
- **Request handling**: < 10ms average
- **Artisan commands**: Comparable to PHP 8.4
- **Test suite**: Complete in < 5 minutes

## Resources

- Laravel Documentation: https://laravel.com/docs
- Laravel Source: https://github.com/laravel/framework
- PHP Compatibility: https://laravel.com/docs/12.x/releases#support-policy

## Progress Tracking

- [x] Install Laravel (2h)
- [ ] Run with PHP-Go (3h)
- [ ] Run test suite (4h)
- [ ] Fix issues (5h)
- [ ] Performance testing (2h)

Total estimated effort: 16 hours
