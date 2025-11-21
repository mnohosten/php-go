# Development Guide

Quick reference for developing PHP-Go.

## Getting Started

### Prerequisites
- Go 1.21 or later
- Git
- 16GB+ RAM recommended
- Multi-core CPU for parallel testing

### Setup
```bash
# Clone repository
cd /path/to/php-go

# Build
go build -o php-go ./cmd/php-go

# Run
./php-go --version
```

## Master Task List

**Primary Reference**: `TODO.md`

This is your single source of truth for all development tasks. It contains:
- All tasks from all 10 phases
- Effort estimates
- References to detailed docs
- Progress tracking with checkboxes

### Using TODO.md

**Daily Workflow**:
1. Open `TODO.md`
2. Find next unchecked `[ ]` task
3. Read referenced phase doc for details
4. Implement the task
5. Write tests
6. Mark task complete `[x]`
7. Commit

**Example**:
```bash
# You see in TODO.md:
# Phase 1: Foundation
# ### 1.1 Token System (6h)
# - [ ] Define Token struct with type, literal, position (2h)

# 1. Read details
cat docs/phases/01-foundation/README.md

# 2. Implement
vim pkg/lexer/token.go

# 3. Test
go test ./pkg/lexer/...

# 4. Mark complete
# Change [ ] to [x] in TODO.md

# 5. Commit
git add .
git commit -m "Implement Token struct (Phase 1, Task 1.1)"
```

## Project Structure

```
php-go/
├── TODO.md              ← MASTER TASK LIST (start here!)
├── ROADMAP.md           ← Timeline and milestones
├── README.md            ← Project overview
├── cmd/php-go/          ← CLI executable
├── pkg/                 ← Core packages
│   ├── lexer/           ← Phase 1
│   ├── parser/          ← Phase 1
│   ├── ast/             ← Phase 1
│   ├── compiler/        ← Phase 2
│   ├── vm/              ← Phase 3
│   ├── types/           ← Phase 3-4
│   ├── runtime/         ← Phase 3
│   ├── stdlib/          ← Phase 6
│   ├── parallel/        ← Phase 7
│   └── goext/           ← Phase 8
├── docs/                ← Detailed documentation
│   ├── 00-project-overview.md
│   ├── 01-php-analysis.md
│   ├── 02-go-architecture.md
│   └── phases/          ← Phase-specific details
│       ├── 01-foundation/
│       ├── 02-compiler/
│       ├── 03-runtime-vm/
│       ├── 04-data-structures/
│       ├── 05-objects/
│       ├── 06-stdlib/
│       ├── 07-parallelization/
│       ├── 08-go-integration/
│       ├── 09-advanced/
│       └── 10-testing/
├── tests/               ← Test suite
└── benchmarks/          ← Performance tests
```

## Development Flow

### Phase-Based Development

**Current Phase**: Phase 1 (Foundation)

1. **Read Phase Doc**: `docs/phases/01-foundation/README.md`
2. **Check TODO.md**: Find Phase 1 tasks
3. **Implement**: Work through tasks sequentially
4. **Test**: Maintain 85%+ coverage
5. **Complete**: Mark tasks done in TODO.md

### Task Details

Each task in TODO.md has:
- **Title**: What to build
- **Effort**: Hours estimate
- **Files**: Where to write code
- **Reference**: Link to detailed docs

**Example**:
```markdown
### 1.1 Token System (6h)
- [ ] Define Token struct with type, literal, position (2h)

**Files**: `pkg/lexer/token.go`
```

Then check `docs/phases/01-foundation/README.md` for:
- Detailed requirements
- Code examples
- Testing requirements
- Success criteria

## Testing

### Unit Tests
```bash
# Test single package
go test ./pkg/lexer/

# Test with coverage
go test -cover ./pkg/lexer/

# Test all
go test ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Benchmarks
```bash
# Run benchmarks
go test -bench=. ./pkg/lexer/

# With memory stats
go test -bench=. -benchmem ./pkg/lexer/
```

### Target Coverage
- Overall: 85%+
- Critical components: 90%+
- New code: Always include tests

## Code Style

### Go Conventions
- Follow official Go style
- Use `gofmt` / `goimports`
- Clear naming
- Document exported items

### Example
```go
// Token represents a lexical token in PHP source code.
type Token struct {
    Type    TokenType  // Token type (keyword, operator, etc.)
    Literal string     // Actual token text
    Pos     Position   // Source position
}

// String returns a human-readable representation of the token.
func (t Token) String() string {
    return fmt.Sprintf("%s(%q) at %s", t.Type, t.Literal, t.Pos)
}
```

## Commit Messages

### Format
```
<type>(<phase>): <subject>

<body>

Refs: #<issue> (if applicable)
```

### Examples
```
feat(phase1): Add Token struct and type definitions

Implemented Token struct with type, literal, and position fields.
Added all PHP 8.4 token type constants.

Refs: TODO.md Phase 1, Task 1.1
```

```
test(phase1): Add lexer unit tests

Added comprehensive tests for Token and Position structs.
Coverage: 92%

Refs: TODO.md Phase 1, Task 1.12
```

### Types
- `feat`: New feature
- `fix`: Bug fix
- `test`: Add tests
- `docs`: Documentation
- `refactor`: Code refactoring
- `perf`: Performance improvement
- `chore`: Maintenance

## Documentation

### When to Document

**Always**:
- Exported functions/types
- Complex algorithms
- Design decisions
- Non-obvious code

**In docs/**:
- Architecture changes
- New features
- Implementation notes

### Example
```go
// parseExpression parses a PHP expression with operator precedence.
//
// This uses Pratt parsing (precedence climbing) to handle PHP's
// complex operator precedence rules. See docs/internals/parser.md
// for algorithm details.
//
// Returns the parsed expression AST node or an error if the
// expression is invalid.
func (p *Parser) parseExpression(precedence int) (ast.Expr, error) {
    // ...
}
```

## Progress Tracking

### Check Progress
```bash
# Completed tasks
grep -c "\[x\]" TODO.md

# Total tasks
grep -c "\[ \]" TODO.md
grep -c "\[x\]" TODO.md

# Current phase
grep "Phase [0-9].*⬜" TODO.md | head -1
```

### Update Status
1. Mark tasks complete: `[ ]` → `[x]`
2. Update phase status in TODO.md
3. Update ROADMAP.md milestones
4. Commit changes

### Milestone Checklist

**Phase 1 Complete**:
- [ ] All Phase 1 tasks checked in TODO.md
- [ ] Tests passing with 85%+ coverage
- [ ] Can parse real PHP files
- [ ] CLI `lex` and `parse` commands work
- [ ] Documentation updated
- [ ] Update TODO.md phase status: ⬜ → ✅

## Common Tasks

### Add New File
```bash
# Create file
vim pkg/lexer/token.go

# Add to git
git add pkg/lexer/token.go

# Test
go test ./pkg/lexer/

# Update TODO.md
# Mark task as [x] complete
```

### Run Examples
```bash
# Build
go build -o php-go ./cmd/php-go

# Run (Phase 1+)
./php-go lex test.php
./php-go parse test.php

# Run (Phase 3+)
./php-go test.php
```

### Profile Performance
```bash
# CPU profile
go test -cpuprofile=cpu.prof -bench=. ./pkg/lexer/
go tool pprof cpu.prof

# Memory profile
go test -memprofile=mem.prof -bench=. ./pkg/lexer/
go tool pprof mem.prof
```

### Check Dependencies
```bash
# List dependencies
go list -m all

# Update dependencies
go get -u ./...
go mod tidy
```

## Getting Help

### Documentation
1. **TODO.md** - What to build next
2. **docs/phases/** - How to build it
3. **docs/01-php-analysis.md** - PHP reference
4. **docs/02-go-architecture.md** - Go design
5. **php-src/** - PHP source code (reference)

### When Stuck
1. Read relevant phase doc
2. Check PHP source code
3. Look at similar implementations
4. Ask for help (GitHub issues)
5. Document the solution

## Phase 1 Quick Start

**Right Now** (Phase 1 starting):

1. **Read Phase 1 Doc**:
   ```bash
   cat docs/phases/01-foundation/README.md
   ```

2. **Start Task 1.1**:
   ```bash
   vim pkg/lexer/token.go
   ```

3. **Implement Token struct** (see phase doc for details)

4. **Write tests**:
   ```bash
   vim pkg/lexer/token_test.go
   go test ./pkg/lexer/
   ```

5. **Mark complete**:
   - Edit TODO.md
   - Change `[ ]` to `[x]` for Task 1.1

6. **Commit**:
   ```bash
   git add .
   git commit -m "feat(phase1): Add Token struct and constants"
   ```

7. **Next task**: Move to Task 1.2 in TODO.md

## Tips

### Stay Organized
- ✅ Always check TODO.md first
- ✅ Read phase docs for context
- ✅ Test as you go
- ✅ Commit often
- ✅ Update progress

### Avoid Common Mistakes
- ❌ Don't skip tests
- ❌ Don't work on multiple phases at once
- ❌ Don't forget to update TODO.md
- ❌ Don't commit broken code
- ❌ Don't optimize prematurely

### Performance
- Profile before optimizing
- Focus on hot paths
- Measure improvements
- Document tradeoffs

### Code Quality
- Write clear code
- Add comments for complex logic
- Keep functions small
- Use descriptive names
- Follow Go idioms

## Resources

### PHP Reference
- `php-src/` - PHP source code
- [PHP Language Spec](https://github.com/php/php-langspec)
- [PHP Internals Book](http://www.phpinternalsbook.com/)

### Go Resources
- [Go Documentation](https://go.dev/doc/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Tools
- `gofmt` - Format code
- `goimports` - Organize imports
- `go vet` - Check for issues
- `golangci-lint` - Comprehensive linter

## Next Steps

**Right now**:
1. ✅ You've read this guide
2. → Open `TODO.md`
3. → Read `docs/phases/01-foundation/README.md`
4. → Start Task 1.1: Define Token struct
5. → Begin coding!

**Good luck building PHP-Go!** 🚀

---

**Last Updated**: 2025-11-21
**Current Phase**: Phase 1 Starting
**Next Task**: TODO.md - Phase 1, Task 1.1
