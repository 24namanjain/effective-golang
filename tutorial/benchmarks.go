package main

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// =============================================================================
// GO BENCHMARKING TUTORIAL
// =============================================================================

// Benchmarks in Go measure the performance of your code by running it multiple times
// and measuring how long it takes. They help you identify performance bottlenecks
// and compare different implementations.

// =============================================================================
// 1. BASIC BENCHMARK
// =============================================================================

// Function to benchmark: Simple string concatenation
func concatenateStrings(strs []string) string {
	result := ""
	for _, s := range strs {
		result += s
	}
	return result
}

// Basic benchmark function
func BenchmarkConcatenateStrings(b *testing.B) {
	// Prepare test data
	testData := []string{"hello", "world", "golang", "benchmarking"}

	// Reset the timer to exclude setup time
	b.ResetTimer()

	// Run the benchmark b.N times
	for i := 0; i < b.N; i++ {
		concatenateStrings(testData)
	}
}

// =============================================================================
// 2. COMPARING DIFFERENT IMPLEMENTATIONS
// =============================================================================

// Alternative implementation using strings.Builder
func concatenateWithBuilderBench(strs []string) string {
	var builder strings.Builder
	for _, s := range strs {
		builder.WriteString(s)
	}
	return builder.String()
}

// Alternative implementation using strings.Join
func concatenateWithJoin(strs []string) string {
	return strings.Join(strs, "")
}

// Benchmark comparing different concatenation methods
func BenchmarkStringConcatenation(b *testing.B) {
	testData := []string{"hello", "world", "golang", "benchmarking", "performance", "testing"}

	// Benchmark string concatenation with +
	b.Run("String_Plus", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			concatenateStrings(testData)
		}
	})

	// Benchmark using strings.Builder
	b.Run("Strings_Builder", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			concatenateWithBuilderBench(testData)
		}
	})

	// Benchmark using strings.Join
	b.Run("Strings_Join", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			concatenateWithJoin(testData)
		}
	})
}

// =============================================================================
// 3. BENCHMARKING WITH DIFFERENT INPUT SIZES
// =============================================================================

// Function to benchmark: Sum of integers
func sumIntegers(nums []int) int {
	sum := 0
	for _, num := range nums {
		sum += num
	}
	return sum
}

// Benchmark with different slice sizes
func BenchmarkSumIntegers(b *testing.B) {
	benchmarks := []struct {
		name string
		size int
	}{
		{"Small_10", 10},
		{"Medium_100", 100},
		{"Large_1000", 1000},
		{"Very_Large_10000", 10000},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			// Generate test data
			nums := make([]int, bm.size)
			for i := range nums {
				nums[i] = i
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				sumIntegers(nums)
			}
		})
	}
}

// =============================================================================
// 4. BENCHMARKING MEMORY ALLOCATIONS
// =============================================================================

// Function that allocates memory
func createSliceWithAllocation(size int) []int {
	slice := make([]int, size)
	for i := range slice {
		slice[i] = i
	}
	return slice
}

// Function that reuses memory
func createSliceWithReuse(size int, slice []int) []int {
	if cap(slice) < size {
		slice = make([]int, size)
	} else {
		slice = slice[:size]
	}
	for i := range slice {
		slice[i] = i
	}
	return slice
}

// Benchmark comparing memory allocation strategies
func BenchmarkMemoryAllocation(b *testing.B) {
	b.Run("New_Allocation", func(b *testing.B) {
		b.ReportAllocs() // Report memory allocations
		for i := 0; i < b.N; i++ {
			createSliceWithAllocation(1000)
		}
	})

	b.Run("Reuse_Allocation", func(b *testing.B) {
		b.ReportAllocs() // Report memory allocations
		var slice []int
		for i := 0; i < b.N; i++ {
			slice = createSliceWithReuse(1000, slice)
		}
	})
}

// =============================================================================
// 5. BENCHMARKING WITH SETUP AND TEARDOWN
// =============================================================================

// Function to benchmark: String search
func findString(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

// Benchmark with setup and teardown
func BenchmarkStringSearch(b *testing.B) {
	// Setup: Create test data
	haystack := make([]string, 1000)
	for i := range haystack {
		haystack[i] = fmt.Sprintf("string_%d", i)
	}
	needle := "string_999" // Last element

	b.ResetTimer()

	// Run benchmark
	for i := 0; i < b.N; i++ {
		findString(haystack, needle)
	}
}

// =============================================================================
// 6. BENCHMARKING CONCURRENT OPERATIONS
// =============================================================================

// Function to benchmark: Concurrent sum calculation
func sumConcurrent(nums []int, numGoroutines int) int {
	chunkSize := len(nums) / numGoroutines
	if chunkSize == 0 {
		chunkSize = 1
	}

	results := make(chan int, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(nums) {
			end = len(nums)
		}

		go func(nums []int) {
			sum := 0
			for _, num := range nums {
				sum += num
			}
			results <- sum
		}(nums[start:end])
	}

	total := 0
	for i := 0; i < numGoroutines; i++ {
		total += <-results
	}

	return total
}

// Benchmark concurrent vs sequential operations
func BenchmarkConcurrentSum(b *testing.B) {
	nums := make([]int, 10000)
	for i := range nums {
		nums[i] = i
	}

	b.Run("Sequential", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sumIntegers(nums)
		}
	})

	b.Run("Concurrent_2", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sumConcurrent(nums, 2)
		}
	})

	b.Run("Concurrent_4", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sumConcurrent(nums, 4)
		}
	})

	b.Run("Concurrent_8", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sumConcurrent(nums, 8)
		}
	})
}

// =============================================================================
// 7. BENCHMARKING WITH CUSTOM METRICS
// =============================================================================

// Function to benchmark: Fibonacci calculation
func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

// Benchmark with custom metrics
func BenchmarkFibonacci(b *testing.B) {
	benchmarks := []struct {
		name string
		n    int
	}{
		{"Fib_10", 10},
		{"Fib_20", 20},
		{"Fib_30", 30},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				fibonacci(bm.n)
			}
		})
	}
}

// =============================================================================
// 8. BENCHMARKING WITH TABLE-DRIVEN TESTS
// =============================================================================

// Function to benchmark: String reversal
func reverseStringBench(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Table-driven benchmark
func BenchmarkReverseString(b *testing.B) {
	benchmarks := []struct {
		name  string
		input string
	}{
		{"Short", "hello"},
		{"Medium", "hello world golang"},
		{"Long", "this is a very long string that we want to reverse for benchmarking purposes"},
		{"Unicode", "café résumé naïve"},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				reverseStringBench(bm.input)
			}
		})
	}
}

// =============================================================================
// 9. BENCHMARKING WITH PARALLEL EXECUTION
// =============================================================================

// Function to benchmark: CPU-intensive operation
func cpuIntensiveOperation(n int) int {
	result := 0
	for i := 0; i < n; i++ {
		result += i * i
	}
	return result
}

// Benchmark with parallel execution
func BenchmarkCPUIntensive(b *testing.B) {
	b.Run("Sequential", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cpuIntensiveOperation(1000)
		}
	})

	b.Run("Parallel", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				cpuIntensiveOperation(1000)
			}
		})
	})
}

// =============================================================================
// 10. BENCHMARKING WITH CUSTOM TIMERS
// =============================================================================

// Function to benchmark: Database-like operation simulation
func simulateDatabaseOperation() {
	// Simulate some work
	time.Sleep(1 * time.Millisecond)
}

// Benchmark with custom timing
func BenchmarkDatabaseOperation(b *testing.B) {
	b.ResetTimer()

	// Start custom timer
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		simulateDatabaseOperation()
	}

	// Stop custom timer
	b.StopTimer()
}

// =============================================================================
// DEMO FUNCTION TO SHOW BENCHMARKING CONCEPTS
// =============================================================================

// Demo function to explain benchmarking concepts
func benchmarkingDemo() {
	fmt.Println("⚡ GO BENCHMARKING TUTORIAL")
	fmt.Println("============================")

	fmt.Println("\n1. What are Benchmarks?")
	fmt.Println("   - Benchmarks measure code performance")
	fmt.Println("   - They run your code multiple times (b.N)")
	fmt.Println("   - They report time per operation")
	fmt.Println("   - They help identify performance bottlenecks")

	fmt.Println("\n2. Basic Benchmark Structure:")
	fmt.Println("   func BenchmarkFunctionName(b *testing.B) {")
	fmt.Println("       for i := 0; i < b.N; i++ {")
	fmt.Println("           // Code to benchmark")
	fmt.Println("       }")
	fmt.Println("   }")

	fmt.Println("\n3. Running Benchmarks:")
	fmt.Println("   go test -bench=.")
	fmt.Println("   go test -bench=BenchmarkFunctionName")
	fmt.Println("   go test -bench=BenchmarkFunctionName -benchmem")
	fmt.Println("   go test -bench=BenchmarkFunctionName -benchtime=5s")

	fmt.Println("\n4. Benchmark Output:")
	fmt.Println("   BenchmarkFunctionName-8   1000000   1234 ns/op")
	fmt.Println("   - Function name")
	fmt.Println("   - Number of CPU cores")
	fmt.Println("   - Number of iterations")
	fmt.Println("   - Time per operation")

	fmt.Println("\n5. Best Practices:")
	fmt.Println("   ✅ Use b.ResetTimer() to exclude setup time")
	fmt.Println("   ✅ Use b.ReportAllocs() to track memory usage")
	fmt.Println("   ✅ Use b.Run() for sub-benchmarks")
	fmt.Println("   ✅ Use b.RunParallel() for concurrent benchmarks")
	fmt.Println("   ✅ Compare different implementations")
	fmt.Println("   ✅ Test with different input sizes")

	fmt.Println("\n6. Common Benchmark Patterns:")
	fmt.Println("   - Comparing implementations")
	fmt.Println("   - Testing with different input sizes")
	fmt.Println("   - Measuring memory allocations")
	fmt.Println("   - Benchmarking concurrent operations")
	fmt.Println("   - Table-driven benchmarks")

	fmt.Println("\n✅ Benchmarks help you write faster, more efficient code!")
}
