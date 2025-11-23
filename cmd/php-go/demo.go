package main

import (
	"fmt"
	"time"

	"github.com/krizos/php-go/pkg/parallel"
)

func runParallelDemo() {
	fmt.Println("Demonstrating PHP-Go's parallelization features")

	// Demo 1: Worker Pool
	demo1WorkerPool()

	// Demo 2: COW Optimization
	demo2COWOptimization()

	// Demo 3: Parallel Array Map
	demo3ParallelArrayMap()

	// Demo 4: Request Isolation
	demo4RequestIsolation()

	fmt.Println("\n=== Demo Complete ===")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Run tests: go test ./pkg/parallel -v")
	fmt.Println("  2. Run benchmarks: go test ./pkg/parallel -bench=. -benchmem")
	fmt.Println("  3. Run with race detector: go test ./pkg/parallel -race")
	fmt.Println("  4. See examples in examples/parallel_examples.php")
}

func demo1WorkerPool() {
	fmt.Println("--- Demo 1: Worker Pool ---")
	fmt.Println("Creating worker pool with 4 workers...")

	pool := parallel.NewWorkerPool(4)
	pool.Start() // Must start the pool before submitting tasks
	defer pool.Shutdown()

	// Submit tasks
	numTasks := 10
	tasks := make([]*parallel.Task, numTasks)

	fmt.Printf("Submitting %d tasks...\n", numTasks)
	start := time.Now()

	for i := 0; i < numTasks; i++ {
		taskNum := i
		task := parallel.NewTask(fmt.Sprintf("task-%d", i), func() (interface{}, error) {
			// Simulate work
			time.Sleep(10 * time.Millisecond)
			return fmt.Sprintf("Result from task %d", taskNum), nil
		})
		pool.Submit(task)
		tasks[i] = task
	}

	// Wait for all tasks
	for i, task := range tasks {
		result, err := task.Wait()
		if err != nil {
			fmt.Printf("  Task %d failed: %v\n", i, err)
		} else {
			fmt.Printf("  Task %d: %v\n", i, result)
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("Completed %d tasks in %v (with 4 workers)\n", numTasks, elapsed)
	fmt.Printf("Sequential would take: ~%v\n", time.Duration(numTasks)*10*time.Millisecond)
	fmt.Printf("Speedup: ~%.1fx\n\n", float64(numTasks*10)/float64(elapsed.Milliseconds()))
}

func demo2COWOptimization() {
	fmt.Println("--- Demo 2: Copy-on-Write Optimization ---")
	fmt.Println("Creating large array with 10,000 elements...")

	// Create large array
	data := make([]interface{}, 10000)
	for i := 0; i < 10000; i++ {
		data[i] = i
	}

	cowArray := parallel.NewCOWArray(data)
	fmt.Printf("Created COW array, ref count: %d\n", cowArray.RefCount())

	// Clone it 4 times (for 4 workers)
	clones := make([]*parallel.COWArray, 4)
	for i := 0; i < 4; i++ {
		clones[i] = cowArray.Clone()
	}

	fmt.Printf("After 4 clones, original ref count: %d\n", cowArray.RefCount())
	fmt.Printf("Memory: 1 shared copy instead of 5 copies!\n")

	// Read from all clones (no copy triggered)
	fmt.Println("Reading from all clones (no copy)...")
	for i, clone := range clones {
		val := clone.Get(0)
		fmt.Printf("  Clone %d reads: %v\n", i, val)
	}

	// Demonstrate memory savings
	fmt.Println("\nMemory Savings Calculation:")
	elementSize := 8 // bytes per interface{}
	arraySize := 10000 * elementSize
	fmt.Printf("  Array size: %d bytes (%.1f KB)\n", arraySize, float64(arraySize)/1024)
	fmt.Printf("  Without COW: 5 copies × %d bytes = %d bytes (%.1f KB)\n",
		arraySize, arraySize*5, float64(arraySize*5)/1024)
	fmt.Printf("  With COW: 1 copy × %d bytes = %d bytes (%.1f KB)\n",
		arraySize, arraySize, float64(arraySize)/1024)
	fmt.Printf("  Memory saved: %d bytes (%.1f KB) = %.0f%%!\n\n",
		arraySize*4, float64(arraySize*4)/1024, 80.0)
}

func demo3ParallelArrayMap() {
	fmt.Println("--- Demo 3: Parallel Array Map ---")
	fmt.Println("Processing 1,000 elements with parallel map...")

	// Create test data
	arr := make([]interface{}, 1000)
	for i := 0; i < 1000; i++ {
		arr[i] = i
	}

	// Configure parallel map
	config := parallel.DefaultArrayMapConfig()
	config.MinSize = 100  // Parallelize arrays with 100+ elements
	config.NumWorkers = 4

	// Sequential version
	fmt.Println("Running sequential map...")
	seqStart := time.Now()
	seqResult, _ := parallel.ParallelArrayMap(arr, func(v interface{}) (interface{}, error) {
		return v.(int) * 2, nil
	}, parallel.ArrayMapConfig{MinSize: 10000}) // High threshold = sequential
	seqTime := time.Since(seqStart)

	// Parallel version
	fmt.Println("Running parallel map with 4 workers...")
	parStart := time.Now()
	parResult, _ := parallel.ParallelArrayMap(arr, func(v interface{}) (interface{}, error) {
		return v.(int) * 2, nil
	}, config)
	parTime := time.Since(parStart)

	fmt.Printf("Sequential: %d results in %v\n", len(seqResult), seqTime)
	fmt.Printf("Parallel:   %d results in %v\n", len(parResult), parTime)

	if parTime < seqTime {
		fmt.Printf("Speedup: %.2fx faster\n\n", float64(seqTime)/float64(parTime))
	} else {
		fmt.Printf("(For very fast operations, parallel overhead may not be worth it)\n\n")
	}
}

func demo4RequestIsolation() {
	fmt.Println("--- Demo 4: Request Isolation ---")
	fmt.Println("Simulating 3 concurrent HTTP requests...")

	manager := parallel.NewRequestManager(10, 30*time.Second)

	// Launch 3 concurrent "requests"
	for i := 0; i < 3; i++ {
		requestID := fmt.Sprintf("req-%d", i)

		go func(id string, num int) {
			ctx, err := manager.CreateRequest(id)
			if err != nil {
				fmt.Printf("Request %s failed to start: %v\n", id, err)
				return
			}
			defer manager.CompleteRequest(id)

			// Each request has isolated globals
			ctx.SetGlobal("REQUEST_ID", id)
			ctx.SetGlobal("REQUEST_NUM", num)

			// Simulate request processing
			time.Sleep(50 * time.Millisecond)

			// Write to isolated output buffer
			ctx.WriteString(fmt.Sprintf("Response from %s\n", id))

			reqID, _ := ctx.GetGlobal("REQUEST_ID")
			fmt.Printf("  %s completed (globals isolated)\n", reqID)
		}(requestID, i)
	}

	// Wait for requests to complete
	time.Sleep(200 * time.Millisecond)

	stats := manager.Stats()
	fmt.Printf("\nRequest Manager Stats:\n")
	fmt.Printf("  Total:     %d\n", stats.TotalRequests)
	fmt.Printf("  Active:    %d\n", stats.ActiveRequests)
	fmt.Printf("  Completed: %d\n", stats.CompletedRequests)
	fmt.Printf("\nEach request has isolated:\n")
	fmt.Printf("  - Globals ($_SERVER, $GLOBALS, etc.)\n")
	fmt.Printf("  - Output buffer\n")
	fmt.Printf("  - Error tracking\n")
	fmt.Printf("  - Like PHP-FPM but with goroutines!\n\n")
}
