package main

import (
	"fmt"
	"time"
)

// =============================================================================
// CHANNELS TUTORIAL - Understanding Go Channels
// =============================================================================

// Channels are Go's way of communicating between goroutines (concurrent functions)
// Think of them as pipes that allow data to flow from one goroutine to another

// =============================================================================
// 1. BASIC CHANNEL OPERATIONS
// =============================================================================

// Example 1: Simple channel communication
func basicChannelDemo() {
	fmt.Println("\n=== 1. Basic Channel Demo ===")

	// Create an unbuffered channel (capacity = 0)
	// Unbuffered channels block until sender and receiver are ready
	messageChan := make(chan string)

	// Start a goroutine to send a message
	go func() {
		fmt.Println("Sender: Preparing to send message...")
		time.Sleep(1 * time.Second)            // Simulate some work
		messageChan <- "Hello from goroutine!" // Send data to channel
		fmt.Println("Sender: Message sent!")
	}()

	// Main goroutine receives the message
	fmt.Println("Receiver: Waiting for message...")
	message := <-messageChan // Receive data from channel
	fmt.Printf("Receiver: Got message: %s\n", message)
}

// =============================================================================
// 2. BUFFERED CHANNELS
// =============================================================================

// Example 2: Buffered channels (with capacity)
func bufferedChannelDemo() {
	fmt.Println("\n=== 2. Buffered Channel Demo ===")

	// Create a buffered channel with capacity 4
	// Buffered channels can hold multiple values before blocking
	bufferedChan := make(chan int, 4)

	fmt.Println("Sending 3 values to buffered channel...")

	// Send multiple values without blocking
	bufferedChan <- 1
	fmt.Println("Sent: 1")
	bufferedChan <- 2
	fmt.Println("Sent: 2")
	bufferedChan <- 3
	fmt.Println("Sent: 3")

	// Try to send a 4th value - this would block if channel was full
	// But since we have capacity 3, it works fine
	bufferedChan <- 4
	fmt.Println("Sent: 4")

	fmt.Println("\nReceiving values from buffered channel...")

	// Receive all values
	fmt.Printf("Received: %d\n", <-bufferedChan)
	fmt.Printf("Received: %d\n", <-bufferedChan)
	fmt.Printf("Received: %d\n", <-bufferedChan)
	fmt.Printf("Received: %d\n", <-bufferedChan)

	// Close the channel to indicate no more values will be sent
	close(bufferedChan)
}

// =============================================================================
// 3. CHANNEL DIRECTION (SEND-ONLY, RECEIVE-ONLY)
// =============================================================================

// Example 3: Channel direction and function parameters
func channelDirectionDemo() {
	fmt.Println("\n=== 3. Channel Direction Demo ===")

	// Create a channel
	dataChan := make(chan int, 5)

	// Start producer (send-only channel)
	go producer(dataChan)

	// Start consumer (receive-only channel)
	go consumer(dataChan)

	// Wait for both to complete
	time.Sleep(2 * time.Second)
}

// producer sends data to a send-only channel (chan<- int)
func producer(sendChan chan<- int) {
	fmt.Println("Producer: Starting to produce data...")

	for i := 1; i <= 5; i++ {
		fmt.Printf("Producer: Sending %d\n", i)
		sendChan <- i // Can only send to this channel
		time.Sleep(200 * time.Millisecond)
	}

	close(sendChan) // Close channel when done
	fmt.Println("Producer: Finished producing data")
}

// consumer receives data from a receive-only channel (<-chan int)
func consumer(receiveChan <-chan int) {
	fmt.Println("Consumer: Starting to consume data...")

	for value := range receiveChan { // Range over channel until closed
		fmt.Printf("Consumer: Received %d\n", value)
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("Consumer: Finished consuming data")
}

// =============================================================================
// 4. SELECT STATEMENT - HANDLING MULTIPLE CHANNELS
// =============================================================================

// Example 4: Using select to handle multiple channels
func selectDemo() {
	fmt.Println("\n=== 4. Select Statement Demo ===")

	// Create multiple channels
	channel1 := make(chan string)
	channel2 := make(chan string)

	// Start goroutines that send to different channels
	go func() {
		time.Sleep(1 * time.Second)
		channel1 <- "Message from channel 1"
	}()

	go func() {
		time.Sleep(500 * time.Millisecond)
		channel2 <- "Message from channel 2"
	}()

	// Use select to handle whichever channel has data first
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-channel1:
			fmt.Printf("Received from channel1: %s\n", msg1)
		case msg2 := <-channel2:
			fmt.Printf("Received from channel2: %s\n", msg2)
		case <-time.After(2 * time.Second):
			fmt.Println("Timeout: No message received within 2 seconds")
			return
		}
	}
}

// =============================================================================
// 5. RATE LIMITING WITH CHANNELS
// =============================================================================

// Example 5: Rate limiting using time.Tick
func rateLimitingDemo() {
	fmt.Println("\n=== 5. Rate Limiting Demo ===")

	// Create a rate limiter that ticks every 300ms
	rateLimiter := time.Tick(300 * time.Millisecond)

	// Create channels for jobs and results
	jobs := make(chan int, 10)
	results := make(chan int, 10)

	// Start worker goroutines
	for workerID := 1; workerID <= 2; workerID++ {
		go rateLimitedWorker(workerID, jobs, results, rateLimiter)
	}

	// Send jobs
	go func() {
		for jobID := 1; jobID <= 6; jobID++ {
			fmt.Printf("Job Sender: Sending job %d\n", jobID)
			jobs <- jobID
		}
		close(jobs) // Close jobs channel when done
	}()

	// Collect results
	resultCount := 0
	for result := range results {
		fmt.Printf("Result Collector: Got result %d\n", result)
		resultCount++
		if resultCount >= 6 {
			break
		}
	}

	fmt.Printf("Rate Limiting Demo: Processed %d jobs\n", resultCount)
}

// rateLimitedWorker processes jobs at a controlled rate
func rateLimitedWorker(id int, jobs <-chan int, results chan<- int, rateLimiter <-chan time.Time) {
	for job := range jobs {
		// Wait for rate limiter tick before processing
		<-rateLimiter

		fmt.Printf("Worker %d: Processing job %d\n", id, job)
		time.Sleep(100 * time.Millisecond) // Simulate work

		result := job * 2
		results <- result
		fmt.Printf("Worker %d: Completed job %d, result: %d\n", id, job, result)
	}
}

// =============================================================================
// 6. PIPELINE PATTERN
// =============================================================================

// Example 6: Pipeline pattern with multiple stages
func pipelineDemo() {
	fmt.Println("\n=== 6. Pipeline Demo ===")

	// Stage 1: Generate numbers
	numbers := generateNumbers(5)

	// Stage 2: Square the numbers
	squared := squareNumbers(numbers)

	// Stage 3: Add 10 to the squared numbers
	final := addTen(squared)

	// Collect results
	fmt.Println("Pipeline results:")
	for result := range final {
		fmt.Printf("Final result: %d\n", result)
	}
}

// Stage 1: Generate numbers
func generateNumbers(count int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 1; i <= count; i++ {
			fmt.Printf("Generator: Sending %d\n", i)
			out <- i
			time.Sleep(200 * time.Millisecond)
		}
	}()
	return out
}

// Stage 2: Square numbers
func squareNumbers(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for num := range in {
			squared := num * num
			fmt.Printf("Squarer: %d -> %d\n", num, squared)
			out <- squared
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return out
}

// Stage 3: Add 10
func addTen(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for num := range in {
			result := num + 10
			fmt.Printf("Adder: %d -> %d\n", num, result)
			out <- result
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return out
}

// =============================================================================
// 7. FAN-OUT, FAN-IN PATTERN
// =============================================================================

// Example 7: Fan-out (multiple workers) and Fan-in (collect results)
func fanOutFanInDemo() {
	fmt.Println("\n=== 7. Fan-Out, Fan-In Demo ===")

	// Generate work
	work := generateWork(10)

	// Fan-out: Distribute work to multiple workers
	worker1 := worker("Worker-1", work)
	worker2 := worker("Worker-2", work)
	worker3 := worker("Worker-3", work)

	// Fan-in: Collect results from all workers
	results := fanIn(worker1, worker2, worker3)

	// Process results
	fmt.Println("Fan-Out, Fan-In results:")
	for i := 0; i < 10; i++ {
		result := <-results
		fmt.Printf("Result: %s\n", result)
	}
}

// Generate work items
func generateWork(count int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 1; i <= count; i++ {
			out <- i
		}
	}()
	return out
}

// Worker processes work items
func worker(name string, work <-chan int) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for w := range work {
			time.Sleep(100 * time.Millisecond) // Simulate work
			result := fmt.Sprintf("%s processed item %d", name, w)
			out <- result
		}
	}()
	return out
}

// Fan-in collects results from multiple channels
func fanIn(inputs ...<-chan string) <-chan string {
	out := make(chan string)

	// Start a goroutine for each input channel
	for _, input := range inputs {
		go func(in <-chan string) {
			for result := range in {
				out <- result
			}
		}(input)
	}

	return out
}

// =============================================================================
// MAIN CHANNELS DEMO FUNCTION
// =============================================================================

// Run all channel examples
func channelsDemo() {
	fmt.Println("🔗 CHANNELS TUTORIAL")
	fmt.Println("====================")

	// 1. Basic channel operations
	basicChannelDemo()

	// 2. Buffered channels
	bufferedChannelDemo()

	// 3. Channel direction
	channelDirectionDemo()

	// 4. Select statement
	selectDemo()

	// 5. Rate limiting
	rateLimitingDemo()

	// 6. Pipeline pattern
	pipelineDemo()

	// 7. Fan-out, Fan-in pattern
	fanOutFanInDemo()

	fmt.Println("\n✅ Channels tutorial completed!")
	fmt.Println("Key takeaways:")
	fmt.Println("- Channels are Go's way of communicating between goroutines")
	fmt.Println("- Unbuffered channels block until sender and receiver are ready")
	fmt.Println("- Buffered channels can hold multiple values")
	fmt.Println("- Use select to handle multiple channels")
	fmt.Println("- Pipelines and Fan-out/Fan-in are powerful patterns")
	fmt.Println("- Always close channels when done sending")
}
