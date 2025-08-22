package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Concurrency and Goroutines Tutorial
// This tutorial covers:
// 1. Using context.Context for cancellation and deadlines
// 2. Avoiding goroutine leaks
// 3. Protecting shared resources with sync.Mutex and channels
// 4. Using worker pools for high concurrency workloads

// =============================================================================
// 1. CONTEXT CANCELLATION AND DEADLINES
// =============================================================================

// Example 1: Basic context cancellation
func contextCancellationDemo() {
	fmt.Println("\n=== 1. Context Cancellation Demo ===")

	// Create a context with cancellation capability
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Always defer cancel to prevent goroutine leaks

	// Start a goroutine that listens for cancellation
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Goroutine: Received cancellation signal, exiting gracefully")
				return // Exit the goroutine
			default:
				fmt.Println("Goroutine: Working...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	// Let the goroutine run for a bit
	time.Sleep(2 * time.Second)

	// Cancel the context
	fmt.Println("Main: Sending cancellation signal")
	cancel()

	// Give the goroutine time to exit
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Main: Context cancelled successfully")
}

// Example 2: Context with deadline
func contextDeadlineDemo() {
	fmt.Println("\n=== 2. Context Deadline Demo ===")

	// Create a context with a 2-second deadline
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Start a goroutine that respects the deadline
	go func() {
		for {
			select {
			case <-ctx.Done():
				if ctx.Err() == context.DeadlineExceeded {
					fmt.Println("Goroutine: Deadline exceeded, exiting")
				} else {
					fmt.Println("Goroutine: Context cancelled, exiting")
				}
				return
			default:
				fmt.Println("Goroutine: Processing...")
				time.Sleep(300 * time.Millisecond)
			}
		}
	}()

	// Wait for the deadline to be reached
	time.Sleep(3 * time.Second)
	fmt.Println("Main: Deadline demo completed")
}

// =============================================================================
// 2. AVOIDING GOROUTINE LEAKS
// =============================================================================

// Example 3: Proper goroutine cleanup with channels
func goroutineLeakPreventionDemo() {
	fmt.Println("\n=== 3. Goroutine Leak Prevention Demo ===")

	// Use a done channel to signal goroutine termination
	done := make(chan struct{})

	// Start a goroutine with proper cleanup
	go func() {
		defer close(done) // Ensure channel is closed when goroutine exits

		for i := 0; i < 5; i++ {
			fmt.Printf("Goroutine: Processing item %d\n", i)
			time.Sleep(200 * time.Millisecond)
		}
		fmt.Println("Goroutine: Completed all work")
	}()

	// Wait for the goroutine to finish
	<-done
	fmt.Println("Main: Goroutine completed successfully")
}

// Example 4: Multiple goroutines with WaitGroup
func waitGroupDemo() {
	fmt.Println("\n=== 4. WaitGroup Demo ===")

	var wg sync.WaitGroup

	// Launch multiple goroutines
	for i := 1; i <= 3; i++ {
		wg.Add(1) // Increment counter before starting goroutine

		go func(id int) {
			defer wg.Done() // Decrement counter when goroutine finishes

			fmt.Printf("Worker %d: Starting work\n", id)
			time.Sleep(time.Duration(id) * 300 * time.Millisecond)
			fmt.Printf("Worker %d: Completed work\n", id)
		}(i)
	}

	// Wait for all goroutines to complete
	fmt.Println("Main: Waiting for all workers to complete...")
	wg.Wait()
	fmt.Println("Main: All workers completed")
}

// =============================================================================
// 3. PROTECTING SHARED RESOURCES
// =============================================================================

// Example 5: Using sync.Mutex to protect shared data
func mutexProtectionDemo() {
	fmt.Println("\n=== 5. Mutex Protection Demo ===")

	type SafeCounter struct {
		mu      sync.Mutex
		counter int
	}

	counter := &SafeCounter{}
	var wg sync.WaitGroup

	// Function to increment counter safely
	increment := func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			counter.mu.Lock() // Lock before accessing shared resource
			counter.counter++
			counter.mu.Unlock() // Unlock after accessing shared resource
		}
	}

	// Start multiple goroutines
	wg.Add(3)
	go increment()
	go increment()
	go increment()

	wg.Wait()
	fmt.Printf("Final counter value (should be 3000): %d\n", counter.counter)
}

// Example 6: Using channels for safe communication
func channelProtectionDemo() {
	fmt.Println("\n=== 6. Channel Protection Demo ===")

	// Create a channel for safe communication
	dataChan := make(chan int, 100)
	var wg sync.WaitGroup

	// Producer goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(dataChan) // Close channel when done producing

		for i := 1; i <= 10; i++ {
			fmt.Printf("Producer: Sending %d\n", i)
			dataChan <- i // Send data through channel
			time.Sleep(100 * time.Millisecond)
		}
		fmt.Println("Producer: Finished producing")
	}()

	// Consumer goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()

		for value := range dataChan { // Range over channel until closed
			fmt.Printf("Consumer: Received %d\n", value)
			time.Sleep(50 * time.Millisecond)
		}
		fmt.Println("Consumer: Finished consuming")
	}()

	wg.Wait()
	fmt.Println("Main: Channel communication completed")
}

// =============================================================================
// 4. WORKER POOLS FOR HIGH CONCURRENCY
// =============================================================================

// Example 7: Simple worker pool
func workerPoolDemo() {
	fmt.Println("\n=== 7. Worker Pool Demo ===")

	const (
		numWorkers = 3
		numJobs    = 10
	)

	// Create channels for jobs and results
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	// Worker function
	worker := func(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
		defer wg.Done()

		for job := range jobs {
			fmt.Printf("Worker %d: Processing job %d\n", id, job)
			time.Sleep(200 * time.Millisecond) // Simulate work
			result := job * job
			results <- result
			fmt.Printf("Worker %d: Completed job %d, result: %d\n", id, job, result)
		}
	}

	// Start workers
	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// Send jobs
	go func() {
		for j := 1; j <= numJobs; j++ {
			jobs <- j
		}
		close(jobs) // Close jobs channel when all jobs are sent
	}()

	// Collect results
	go func() {
		wg.Wait()
		close(results) // Close results channel when all workers are done
	}()

	// Print results
	fmt.Println("Results:")
	for result := range results {
		fmt.Printf("Result: %d\n", result)
	}
}

// Example 8: Advanced worker pool with context cancellation
func advancedWorkerPoolDemo() {
	fmt.Println("\n=== 8. Advanced Worker Pool with Context ===")

	const (
		numWorkers = 4
		numJobs    = 15
	)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	// Advanced worker with context awareness
	worker := func(id int, jobs <-chan int, results chan<- int, ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				fmt.Printf("Worker %d: Context cancelled, shutting down\n", id)
				return
			case job, ok := <-jobs:
				if !ok {
					fmt.Printf("Worker %d: No more jobs, shutting down\n", id)
					return
				}

				fmt.Printf("Worker %d: Processing job %d\n", id, job)

				// Simulate work with context awareness
				select {
				case <-ctx.Done():
					fmt.Printf("Worker %d: Cancelled while processing job %d\n", id, job)
					return
				case <-time.After(300 * time.Millisecond):
					result := job * job
					results <- result
					fmt.Printf("Worker %d: Completed job %d, result: %d\n", id, job, result)
				}
			}
		}
	}

	// Start workers
	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, ctx, &wg)
	}

	// Send jobs
	go func() {
		for j := 1; j <= numJobs; j++ {
			select {
			case <-ctx.Done():
				fmt.Println("Job sender: Context cancelled, stopping job distribution")
				return
			case jobs <- j:
				fmt.Printf("Job sender: Sent job %d\n", j)
			}
		}
		close(jobs)
	}()

	// Collect results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Print results
	fmt.Println("Results:")
	resultCount := 0
	for result := range results {
		fmt.Printf("Result: %d\n", result)
		resultCount++
	}
	fmt.Printf("Total results processed: %d\n", resultCount)
}

// =============================================================================
// MAIN TUTORIAL FUNCTION
// =============================================================================

// Run all concurrency tutorials
func concurrencyAndGoroutinesTutorial() {
	fmt.Println("🚀 CONCURRENCY & GOROUTINES TUTORIAL")
	fmt.Println("=====================================")

	// 1. Context cancellation and deadlines
	contextCancellationDemo()
	contextDeadlineDemo()

	// 2. Avoiding goroutine leaks
	goroutineLeakPreventionDemo()
	waitGroupDemo()

	// 3. Protecting shared resources
	mutexProtectionDemo()
	channelProtectionDemo()

	// 4. Worker pools for high concurrency
	workerPoolDemo()
	advancedWorkerPoolDemo()

	fmt.Println("\n✅ Concurrency tutorial completed!")
	fmt.Println("Key takeaways:")
	fmt.Println("- Always use context.Context for cancellation and deadlines")
	fmt.Println("- Use defer statements to ensure proper cleanup")
	fmt.Println("- Protect shared resources with sync.Mutex or channels")
	fmt.Println("- Use worker pools for high concurrency workloads")
	fmt.Println("- Always ensure goroutines exit properly to avoid leaks")
}
