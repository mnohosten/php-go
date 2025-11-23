package parallel

import (
	"testing"
)

// ============================================================================
// Performance Benchmarks - Run with: go test -bench=. -benchmem
// ============================================================================

// ============================================================================
// Worker Pool Benchmarks
// ============================================================================

func BenchmarkWorkerPoolSubmit(b *testing.B) {
	pool := NewWorkerPool(4)
	defer pool.Shutdown()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		task := NewTask("bench", func() (interface{}, error) {
			return nil, nil
		})
		pool.Submit(task)
	}
}

func BenchmarkWorkerPoolSubmitHeavy(b *testing.B) {
	pool := NewWorkerPool(4)
	defer pool.Shutdown()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		task := NewTask("bench", func() (interface{}, error) {
			sum := 0
			for j := 0; j < 100; j++ {
				sum += j
			}
			return sum, nil
		})
		pool.Submit(task)
	}
}

func BenchmarkWorkerPoolSizes(b *testing.B) {
	sizes := []int{1, 2, 4, 8, 16}

	for _, size := range sizes {
		b.Run(string(rune('0'+size)), func(b *testing.B) {
			pool := NewWorkerPool(size)
			defer pool.Shutdown()

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				task := NewTask("bench", func() (interface{}, error) {
					return nil, nil
				})
				pool.Submit(task)
			}
		})
	}
}

// ============================================================================
// COW Array Benchmarks
// ============================================================================

func BenchmarkCOWArrayClone(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}

	for _, size := range sizes {
		data := make([]interface{}, size)
		for i := 0; i < size; i++ {
			data[i] = i
		}
		cowArr := NewCOWArray(data)

		b.Run(string(rune('0'+size/10)), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = cowArr.Clone()
			}
		})
	}
}

func BenchmarkCOWArrayGet(b *testing.B) {
	data := make([]interface{}, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}
	cowArr := NewCOWArray(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cowArr.Get(i % 1000)
	}
}

func BenchmarkCOWArraySet(b *testing.B) {
	data := make([]interface{}, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cowArr := NewCOWArray(data)
		cowArr.Set(0, i)
	}
}

func BenchmarkCOWArraySetShared(b *testing.B) {
	data := make([]interface{}, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}
	cowArr := NewCOWArray(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		clone := cowArr.Clone()
		clone.Set(0, i) // Triggers copy
	}
}

func BenchmarkCOWArrayAppend(b *testing.B) {
	data := make([]interface{}, 100)
	for i := 0; i < 100; i++ {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cowArr := NewCOWArray(data)
		cowArr.Append(i)
	}
}

func BenchmarkCOWArrayToSlice(b *testing.B) {
	data := make([]interface{}, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}
	cowArr := NewCOWArray(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cowArr.ToSlice()
	}
}

// ============================================================================
// COW Map Benchmarks
// ============================================================================

func BenchmarkCOWMapClone(b *testing.B) {
	cowMap := NewCOWMap()
	for i := 0; i < 100; i++ {
		cowMap.Set(string(rune('a'+i%26)), i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cowMap.Clone()
	}
}

func BenchmarkCOWMapGet(b *testing.B) {
	cowMap := NewCOWMap()
	cowMap.Set("key1", "value1")
	cowMap.Set("key2", "value2")
	cowMap.Set("key3", "value3")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cowMap.Get("key1")
	}
}

func BenchmarkCOWMapSet(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cowMap := NewCOWMap()
		cowMap.Set("key", i)
	}
}

func BenchmarkCOWMapSetShared(b *testing.B) {
	cowMap := NewCOWMap()
	cowMap.Set("key1", 1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		clone := cowMap.Clone()
		clone.Set("key2", i) // Triggers copy
	}
}

func BenchmarkCOWMapKeys(b *testing.B) {
	cowMap := NewCOWMap()
	for i := 0; i < 100; i++ {
		cowMap.Set(string(rune('a'+i%26)), i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cowMap.Keys()
	}
}

func BenchmarkCOWMapToMap(b *testing.B) {
	cowMap := NewCOWMap()
	for i := 0; i < 100; i++ {
		cowMap.Set(string(rune('a'+i%26)), i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cowMap.ToMap()
	}
}

// ============================================================================
// COW String Benchmarks
// ============================================================================

func BenchmarkCOWStringClone(b *testing.B) {
	cowStr := NewCOWString("Hello, World! This is a test string for benchmarking.")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cowStr.Clone()
	}
}

func BenchmarkCOWStringString(b *testing.B) {
	cowStr := NewCOWString("Hello, World!")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cowStr.String()
	}
}

func BenchmarkCOWStringAppend(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cowStr := NewCOWString("Hello")
		cowStr.Append(" World")
	}
}

func BenchmarkCOWStringAppendShared(b *testing.B) {
	cowStr := NewCOWString("Hello")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		clone := cowStr.Clone()
		clone.Append(" World") // Triggers copy
	}
}

// ============================================================================
// Parallel Array Benchmarks
// ============================================================================

func BenchmarkParallelArrayMapSizes(b *testing.B) {
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		data := make([]interface{}, size)
		for i := 0; i < size; i++ {
			data[i] = i
		}

		b.Run(string(rune('0'+size/100)), func(b *testing.B) {
			config := DefaultArrayMapConfig()
			config.MinSize = 50
			config.NumWorkers = 4

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ParallelArrayMap(data, func(v interface{}) (interface{}, error) {
					return v.(int) * 2, nil
				}, config)
			}
		})
	}
}

func BenchmarkParallelArrayMapWorkers(b *testing.B) {
	data := make([]interface{}, 10000)
	for i := 0; i < 10000; i++ {
		data[i] = i
	}

	workers := []int{1, 2, 4, 8, 16}

	for _, numWorkers := range workers {
		b.Run(string(rune('0'+numWorkers)), func(b *testing.B) {
			config := DefaultArrayMapConfig()
			config.MinSize = 50
			config.NumWorkers = numWorkers

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ParallelArrayMap(data, func(v interface{}) (interface{}, error) {
					return v.(int) * 2, nil
				}, config)
			}
		})
	}
}

func BenchmarkParallelArrayMapVsSequential(b *testing.B) {
	data := make([]interface{}, 10000)
	for i := 0; i < 10000; i++ {
		data[i] = i
	}

	mapFn := func(v interface{}) (interface{}, error) {
		// Simulate some work
		val := v.(int)
		result := val * 2
		for j := 0; j < 10; j++ {
			result += j
		}
		return result, nil
	}

	b.Run("Parallel", func(b *testing.B) {
		config := DefaultArrayMapConfig()
		config.MinSize = 1000
		config.NumWorkers = 4

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ParallelArrayMap(data, mapFn, config)
		}
	})

	b.Run("Sequential", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sequentialArrayMap(data, mapFn)
		}
	})
}

func BenchmarkParallelArrayFilter(b *testing.B) {
	data := make([]interface{}, 10000)
	for i := 0; i < 10000; i++ {
		data[i] = i
	}

	config := DefaultArrayFilterConfig()
	config.MinSize = 1000
	config.NumWorkers = 4

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParallelArrayFilter(data, func(v interface{}) (bool, error) {
			return v.(int)%2 == 0, nil
		}, config)
	}
}

func BenchmarkParallelArrayFilterVsSequential(b *testing.B) {
	data := make([]interface{}, 10000)
	for i := 0; i < 10000; i++ {
		data[i] = i
	}

	filterFn := func(v interface{}) (bool, error) {
		return v.(int)%2 == 0, nil
	}

	b.Run("Parallel", func(b *testing.B) {
		config := DefaultArrayFilterConfig()
		config.MinSize = 1000
		config.NumWorkers = 4

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ParallelArrayFilter(data, filterFn, config)
		}
	})

	b.Run("Sequential", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sequentialArrayFilter(data, filterFn)
		}
	})
}

func BenchmarkParallelArrayReduce(b *testing.B) {
	data := make([]interface{}, 10000)
	for i := 0; i < 10000; i++ {
		data[i] = 1
	}

	config := DefaultArrayReduceConfig()
	config.MinSize = 2000
	config.NumWorkers = 4

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParallelArrayReduce(data, func(acc, val interface{}) (interface{}, error) {
			return acc.(int) + val.(int), nil
		}, 0, config)
	}
}

func BenchmarkParallelArrayReduceVsSequential(b *testing.B) {
	data := make([]interface{}, 10000)
	for i := 0; i < 10000; i++ {
		data[i] = 1
	}

	reduceFn := func(acc, val interface{}) (interface{}, error) {
		return acc.(int) + val.(int), nil
	}

	b.Run("Parallel", func(b *testing.B) {
		config := DefaultArrayReduceConfig()
		config.MinSize = 2000
		config.NumWorkers = 4

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ParallelArrayReduce(data, reduceFn, 0, config)
		}
	})

	b.Run("Sequential", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sequentialArrayReduce(data, reduceFn, 0)
		}
	})
}

func BenchmarkParallelArrayWalk(b *testing.B) {
	data := make([]interface{}, 10000)
	for i := 0; i < 10000; i++ {
		data[i] = i
	}

	config := DefaultArrayWalkConfig()
	config.MinSize = 1000
	config.NumWorkers = 4

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParallelArrayWalk(data, func(v interface{}, idx int) error {
			_ = v.(int) * 2
			return nil
		}, config)
	}
}

// ============================================================================
// Context and Metrics Benchmarks
// ============================================================================

func BenchmarkRequestContextCreate(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewRequestContext("bench")
	}
}

func BenchmarkRequestContextGetGlobal(b *testing.B) {
	ctx := NewRequestContext("bench")
	ctx.SetGlobal("key", "value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ctx.GetGlobal("key")
	}
}

func BenchmarkRequestContextSetGlobal(b *testing.B) {
	ctx := NewRequestContext("bench")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx.SetGlobal("key", i)
	}
}

func BenchmarkCOWManagerRecordShare(b *testing.B) {
	mgr := NewCOWManager()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mgr.RecordShare(1024)
	}
}

func BenchmarkCOWManagerRecordCopy(b *testing.B) {
	mgr := NewCOWManager()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mgr.RecordCopy(512)
	}
}

func BenchmarkCOWManagerGetStats(b *testing.B) {
	mgr := NewCOWManager()
	mgr.RecordShare(1024)
	mgr.RecordCopy(512)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mgr.GetStats()
	}
}

// ============================================================================
// Sync Primitives Benchmarks
// ============================================================================

func BenchmarkBarrierWait(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			barrier := NewBarrier(4)
			done := make(chan bool)

			for i := 0; i < 4; i++ {
				go func() {
					barrier.Wait()
					done <- true
				}()
			}

			for i := 0; i < 4; i++ {
				<-done
			}
		}
	})
}

func BenchmarkSemaphoreAcquireRelease(b *testing.B) {
	sem := NewSemaphore(10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sem.Acquire()
		sem.Release()
	}
}

func BenchmarkBarrierCreate(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBarrier(10)
	}
}

// ============================================================================
// Integration Benchmarks
// ============================================================================

func BenchmarkIntegrationPoolWithCOW(b *testing.B) {
	pool := NewWorkerPool(4)
	defer pool.Shutdown()

	data := make([]interface{}, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cowArr := NewCOWArray(data)

		for j := 0; j < 4; j++ {
			task := NewTask("bench", func() (interface{}, error) {
				clone := cowArr.Clone()
				return clone.Get(0), nil
			})
			pool.Submit(task)
		}
	}
}

func BenchmarkIntegrationPipelineWithCOW(b *testing.B) {
	pool := NewWorkerPool(4)
	defer pool.Shutdown()

	data := make([]interface{}, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cowArr := NewCOWArray(data)

		task := NewTask("bench", func() (interface{}, error) {
			clone := cowArr.Clone()

			config := DefaultArrayMapConfig()
			config.MinSize = 500
			config.NumWorkers = 2

			return ParallelArrayMap(clone.ToSlice(), func(v interface{}) (interface{}, error) {
				return v.(int) * 2, nil
			}, config)
		})
		pool.Submit(task)
	}
}

func BenchmarkIntegrationMapFilterReduce(b *testing.B) {
	data := make([]interface{}, 10000)
	for i := 0; i < 10000; i++ {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Map
		mapConfig := DefaultArrayMapConfig()
		mapConfig.MinSize = 1000
		mapped, _ := ParallelArrayMap(data, func(v interface{}) (interface{}, error) {
			return v.(int) * 2, nil
		}, mapConfig)

		// Filter
		filterConfig := DefaultArrayFilterConfig()
		filterConfig.MinSize = 1000
		filtered, _ := ParallelArrayFilter(mapped, func(v interface{}) (bool, error) {
			return v.(int) > 5000, nil
		}, filterConfig)

		// Reduce
		reduceConfig := DefaultArrayReduceConfig()
		reduceConfig.MinSize = 1000
		ParallelArrayReduce(filtered, func(acc, val interface{}) (interface{}, error) {
			return acc.(int) + val.(int), nil
		}, 0, reduceConfig)
	}
}

// ============================================================================
// Memory Benchmarks
// ============================================================================

func BenchmarkMemoryCOWArrayClone(b *testing.B) {
	data := make([]interface{}, 10000)
	for i := 0; i < 10000; i++ {
		data[i] = i
	}
	cowArr := NewCOWArray(data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cowArr.Clone()
	}
}

func BenchmarkMemoryCOWArrayCopy(b *testing.B) {
	data := make([]interface{}, 10000)
	for i := 0; i < 10000; i++ {
		data[i] = i
	}
	cowArr := NewCOWArray(data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		clone := cowArr.Clone()
		clone.Set(0, i) // Triggers copy
	}
}

func BenchmarkMemorySliceCopy(b *testing.B) {
	data := make([]interface{}, 10000)
	for i := 0; i < 10000; i++ {
		data[i] = i
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newData := make([]interface{}, len(data))
		copy(newData, data)
	}
}

func BenchmarkMemoryParallelArrayMap(b *testing.B) {
	data := make([]interface{}, 10000)
	for i := 0; i < 10000; i++ {
		data[i] = i
	}

	config := DefaultArrayMapConfig()
	config.MinSize = 1000
	config.NumWorkers = 4

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParallelArrayMap(data, func(v interface{}) (interface{}, error) {
			return v.(int) * 2, nil
		}, config)
	}
}
