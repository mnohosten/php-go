# WordPress Testing for PHP-Go

This directory contains WordPress 6.8.3 for testing PHP-Go's compatibility with a real-world PHP application.

## Installation

WordPress 6.8.3 has been installed in the `wordpress/` subdirectory.

### Structure

```
tests/wordpress/
├── README.md              # This file
├── wordpress.tar.gz       # Downloaded WordPress archive
├── wordpress/             # Extracted WordPress installation
│   ├── wp-admin/          # WordPress admin interface
│   ├── wp-content/        # Themes, plugins, uploads
│   ├── wp-includes/       # WordPress core files
│   ├── index.php          # Main entry point
│   ├── wp-config-sample.php
│   └── ...
└── test-config.php        # Test configuration for PHP-Go (to be created)
```

## Testing Goals

The WordPress test suite aims to:

1. **Compatibility Testing**: Verify that PHP-Go can execute WordPress core files
2. **Feature Coverage**: Test real-world usage of PHP features WordPress depends on
3. **Performance Baseline**: Establish performance metrics vs standard PHP
4. **Issue Identification**: Find and document incompatibilities

## WordPress Requirements

WordPress 6.8.3 requires:
- PHP 7.4+ (recommended PHP 8.0+)
- MySQL 5.7+ or MariaDB 10.4+
- Apache/Nginx with mod_rewrite

### PHP Extensions Used
WordPress uses these PHP extensions:
- **Core**: json, mysqli, curl, openssl, zip, gd/imagick
- **Optional**: bcmath, filter, hash, mbstring, xml, xmlrpc

### Key PHP Features Used
- Object-oriented programming (classes, inheritance, interfaces)
- Namespaces (in some modern plugins)
- Anonymous functions and closures
- Array manipulation (array_map, array_filter, etc.)
- String processing (str_replace, preg_*, etc.)
- File I/O operations
- Database abstraction (wpdb class)
- Hooks system (actions and filters)

## Test Phases

### Phase 1: Basic Execution (Current)
- Load wp-config.php
- Execute wp-settings.php
- Verify core files parse without errors

### Phase 2: Installation Process
- Run WordPress installation
- Create database tables
- Configure initial settings

### Phase 3: Core Functionality
- Load admin dashboard
- Create/edit posts
- Upload media
- Install themes/plugins

### Phase 4: Performance Testing
- Measure page load times
- Compare with standard PHP 8.4
- Profile hotspots

## Test Scripts

### Basic Load Test
```bash
# Test if WordPress core files can be loaded
./php-go tests/wordpress/wordpress/index.php
```

### Installation Test
```bash
# Run WordPress installation (requires database)
./php-go tests/wordpress/wordpress/wp-admin/install.php
```

## Expected Challenges

Based on current PHP-Go implementation status:

1. **Database Support**: WordPress requires mysqli or PDO_MySQL
   - Status: Not yet implemented in PHP-Go
   - Impact: Critical - WordPress cannot run without database

2. **File I/O**: WordPress uses extensive file operations
   - Status: Basic file functions implemented
   - Impact: Medium - may need more file functions

3. **HTTP Functions**: WordPress uses curl and HTTP streams
   - Status: Not yet implemented
   - Impact: High - needed for plugin/theme installation

4. **Image Processing**: GD or Imagick for media handling
   - Status: Not implemented
   - Impact: Medium - affects media uploads

5. **Session Handling**: Used by some plugins
   - Status: Not implemented
   - Impact: Low - core doesn't require it

6. **Output Buffering**: Used extensively in WordPress
   - Status: Implemented in PHP-Go
   - Impact: Should work

## Test Results

Results will be documented as testing progresses.

### Baseline Metrics (to be established)
- [ ] Parse success rate
- [ ] Execution success rate
- [ ] Performance vs PHP 8.4
- [ ] Memory usage
- [ ] Identified incompatibilities

## Database Setup

For full WordPress testing, you'll need a MySQL/MariaDB database:

```sql
CREATE DATABASE wordpress_phpgo_test;
CREATE USER 'wptest'@'localhost' IDENTIFIED BY 'wptest123';
GRANT ALL PRIVILEGES ON wordpress_phpgo_test.* TO 'wptest'@'localhost';
FLUSH PRIVILEGES;
```

## Configuration

A sample wp-config.php for testing:

```php
<?php
define( 'DB_NAME', 'wordpress_phpgo_test' );
define( 'DB_USER', 'wptest' );
define( 'DB_PASSWORD', 'wptest123' );
define( 'DB_HOST', 'localhost' );
define( 'DB_CHARSET', 'utf8mb4' );
define( 'DB_COLLATE', '' );

// Security keys (use WordPress.org secret key generator)
define( 'AUTH_KEY',         'put your unique phrase here' );
define( 'SECURE_AUTH_KEY',  'put your unique phrase here' );
define( 'LOGGED_IN_KEY',    'put your unique phrase here' );
define( 'NONCE_KEY',        'put your unique phrase here' );
define( 'AUTH_SALT',        'put your unique phrase here' );
define( 'SECURE_AUTH_SALT', 'put your unique phrase here' );
define( 'LOGGED_IN_SALT',   'put your unique phrase here' );
define( 'NONCE_SALT',       'put your unique phrase here' );

$table_prefix = 'wp_';

define( 'WP_DEBUG', true );
define( 'WP_DEBUG_LOG', true );
define( 'WP_DEBUG_DISPLAY', false );

if ( ! defined( 'ABSPATH' ) ) {
    define( 'ABSPATH', __DIR__ . '/' );
}

require_once ABSPATH . 'wp-settings.php';
```

## Next Steps

1. Create simple test scripts to verify basic WordPress file loading
2. Document which PHP features WordPress depends on most heavily
3. Prioritize missing PHP-Go features based on WordPress usage
4. Create incremental test cases (load config → load core → load admin)
5. Establish baseline metrics once database support is added

## References

- [WordPress Developer Resources](https://developer.wordpress.org/)
- [WordPress Core Handbook](https://make.wordpress.org/core/handbook/)
- [WordPress Coding Standards](https://developer.wordpress.org/coding-standards/)
