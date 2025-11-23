# Parallelization Examples

⚠️ **Note**: These examples demonstrate **future functionality** that is currently being developed.

## Status

- **Phase 7 (Parallelization)**: Complete - Core parallelization infrastructure implemented
- **Phase 6 (Standard Library)**: In Progress - Many stdlib functions used in these examples are not yet implemented
- **Integration**: Pending - Full integration between parallelization and stdlib is ongoing

## Examples in This Directory

### parallel_examples.php
Comprehensive demonstration of parallelization features including:
- Automatic parallelization with `#[Parallel]` attribute
- Parallel array operations (`pmap`, `pfilter`, `preduce`)
- Concurrent HTTP requests
- Parallel database queries
- Worker pools
- Parallel file processing
- Map-reduce operations
- Graph processing
- Image processing pipeline
- Monte Carlo simulations

### ecommerce_parallel.php
Real-world e-commerce order processing example showing:
- Performance gains from parallelization
- Practical use cases (inventory, pricing, shipping, recommendations)
- Comparison between sequential and parallel execution

## When These Will Work

These examples will be fully functional once:
1. Standard library functions are implemented (Phase 6)
   - `array_map`, `array_filter`, `array_reduce`
   - `file_get_contents`, `file_put_contents`
   - HTTP and database functions
2. Integration testing is complete (Phase 10)

## Technical Documentation

For technical details on how the parallelization is implemented in Go, see:
- `docs/internals/parallelization.md` - Go implementation details
- `docs/phases/07-parallelization/` - Phase 7 documentation

## Running These Examples

Once the features are complete, you'll be able to run:

```bash
./php-go examples/parallel/parallel_examples.php
./php-go examples/parallel/ecommerce_parallel.php
```

For now, these serve as:
- Design specifications
- Integration targets
- Performance benchmarking reference
