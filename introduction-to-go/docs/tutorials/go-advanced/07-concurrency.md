# Concurrency Patterns in Go

This document explains the concurrency patterns and best practices demonstrated in this project.

## 1. Goroutines

Goroutines are lightweight threads managed by the Go runtime.

### Basic Goroutine Usage
```go
// Starting a goroutine
go func() {
    // This runs concurrently
    processData()
}()

// Goroutine with parameters
go processUser(userID)
```

### Example from Game Service
```go
// Start server in a goroutine
go func() {
    if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Printf("Server error: %v", err)
    }
}()
```

**Key Points:**
- Goroutines are cheap (2KB stack initially)
- Use `go` keyword to start a goroutine
- Main function doesn't wait for goroutines by default
- Use synchronization primitives to coordinate

## 2. Channels

Channels are Go's primary mechanism for communication between goroutines.

### Channel Types
```go
// Unbuffered channel (synchronous)
ch := make(chan int)

// Buffered channel (asynchronous)
ch := make(chan int, 100)

// Send-only channel
var sendCh chan<- int

// Receive-only channel
var recvCh <-chan int
```

### Example: Event Processing
```go
type GameService struct {
    eventQueue chan *GameEvent
    // ...
}

// Sending to channel
func (s *GameService) QueueEvent(event *GameEvent) error {
    select {
    case s.eventQueue <- event:
        return nil
    default:
        return ErrEventQueueFull
    }
}

// Receiving from channel
func (ep *EventProcessor) processEvents() {
    defer ep.wg.Done()
    
    for {
        select {
        case event := <-ep.queue:
            ep.wg.Add(1)
            go ep.processEvent(event)
        case <-ep.ctx.Done():
            return
        }
    }
}
```

**Key Points:**
- Unbuffered channels block until sender and receiver are ready
- Buffered channels block only when full
- Use `select` for non-blocking operations
- Close channels to signal completion

## 3. Select Statements

Select allows you to wait on multiple channel operations.

### Basic Select
```go
select {
case msg := <-ch1:
    fmt.Println("Received from ch1:", msg)
case ch2 <- value:
    fmt.Println("Sent to ch2")
case <-time.After(time.Second):
    fmt.Println("Timeout")
default:
    fmt.Println("No channels ready")
}
```

### Example: Context Cancellation
```go
func (ep *EventProcessor) processEvents() {
    defer ep.wg.Done()
    
    for {
        select {
        case event := <-ep.queue:
            ep.wg.Add(1)
            go ep.processEvent(event)
        case <-ep.ctx.Done():
            return
        }
    }
}
```

**Key Points:**
- Select blocks until one case is ready
- Use `default` for non-blocking behavior
- Multiple cases can be ready (random selection)
- Use `time.After()` for timeouts

## 4. Context Package

Context provides request-scoped values, cancellation, and deadlines.

### Context Creation
```go
// Background context
ctx := context.Background()

// With cancellation
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// With deadline
deadline := time.Now().Add(5 * time.Second)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()
```

### Example: Graceful Shutdown
```go
func (app *Application) Start() error {
    // Start server in goroutine
    go func() {
        if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Printf("Server error: %v", err)
        }
    }()
    
    // Wait for shutdown signal
    <-app.shutdownCh
    log.Println("Shutdown signal received")
    
    return app.Shutdown()
}

func (app *Application) Shutdown() error {
    // Cancel context
    app.cancel()
    
    // Create shutdown context with timeout
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Shutdown HTTP server
    if err := app.server.Shutdown(shutdownCtx); err != nil {
        log.Printf("Server shutdown error: %v", err)
    }
    
    return nil
}
```

**Key Points:**
- Always pass context as first parameter
- Check for cancellation in long-running operations
- Use context for timeouts and deadlines
- Don't store context in structs

## 5. Sync Package

The sync package provides synchronization primitives.

### Mutex (Mutual Exclusion)
```go
type InMemoryUserRepository struct {
    users map[string]*User
    mutex sync.RWMutex
}

func (r *InMemoryUserRepository) GetByID(ctx context.Context, id string) (*User, error) {
    r.mutex.RLock()  // Read lock
    defer r.mutex.RUnlock()
    
    user, exists := r.users[id]
    if !exists {
        return nil, models.ErrUserNotFound
    }
    return user, nil
}

func (r *InMemoryUserRepository) Create(ctx context.Context, user *User) error {
    r.mutex.Lock()   // Write lock
    defer r.mutex.Unlock()
    
    r.users[user.ID] = user
    return nil
}
```

### WaitGroup
```go
type EventProcessor struct {
    wg sync.WaitGroup
}

func (ep *EventProcessor) Start() {
    ep.wg.Add(1)
    go ep.processEvents()
}

func (ep *EventProcessor) Stop() {
    ep.cancel()
    ep.wg.Wait()  // Wait for all goroutines to finish
}

func (ep *EventProcessor) processEvents() {
    defer ep.wg.Done()
    
    for {
        select {
        case event := <-ep.queue:
            ep.wg.Add(1)
            go ep.processEvent(event)
        case <-ep.ctx.Done():
            return
        }
    }
}
```

**Key Points:**
- Use `sync.RWMutex` for read-write locks
- Use `sync.WaitGroup` to wait for goroutines
- Always call `defer` with mutex unlocks
- Use `Add()` before starting goroutines

## 6. Worker Pools

Worker pools manage a fixed number of goroutines to process work.

### Example: Event Processing Pool
```go
type EventProcessor struct {
    workers chan struct{}  // Semaphore for limiting workers
    queue   chan *GameEvent
    gameSvc *GameService
    ctx     context.Context
    cancel  context.CancelFunc
    wg      sync.WaitGroup
}

func (ep *EventProcessor) processEvent(event *GameEvent) {
    defer ep.wg.Done()
    
    // Acquire worker slot
    <-ep.workers
    defer func() { ep.workers <- struct{}{} }()
    
    ctx := context.Background()
    
    switch event.EventType {
    case "game_started":
        ep.handleGameStarted(ctx, event)
    case "score_updated":
        ep.handleScoreUpdated(ctx, event)
    case "game_ended":
        ep.handleGameEnded(ctx, event)
    default:
        fmt.Printf("Unknown event type: %s\n", event.EventType)
    }
}
```

**Key Points:**
- Use buffered channels as semaphores
- Limit the number of concurrent workers
- Always release worker slots
- Handle panics in worker goroutines

## 7. Channel Patterns

### Fan-Out Pattern
```go
// Multiple workers processing from same channel
func fanOut(input <-chan *GameEvent, numWorkers int) {
    for i := 0; i < numWorkers; i++ {
        go worker(input)
    }
}

func worker(events <-chan *GameEvent) {
    for event := range events {
        processEvent(event)
    }
}
```

### Fan-In Pattern
```go
// Multiple channels feeding into one
func fanIn(channels ...<-chan *GameEvent) <-chan *GameEvent {
    out := make(chan *GameEvent)
    var wg sync.WaitGroup
    
    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan *GameEvent) {
            defer wg.Done()
            for event := range c {
                out <- event
            }
        }(ch)
    }
    
    go func() {
        wg.Wait()
        close(out)
    }()
    
    return out
}
```

### Pipeline Pattern
```go
func pipeline(input <-chan *GameEvent) <-chan *ProcessedEvent {
    out := make(chan *ProcessedEvent)
    
    go func() {
        defer close(out)
        for event := range input {
            processed := processEvent(event)
            select {
            case out <- processed:
            case <-ctx.Done():
                return
            }
        }
    }()
    
    return out
}
```

## 8. Best Practices

### Goroutine Management
```go
// Always ensure goroutines exit
func (ep *EventProcessor) Stop() {
    ep.cancel()  // Signal all goroutines to stop
    ep.wg.Wait() // Wait for all to finish
}

// Use context for cancellation
func worker(ctx context.Context, jobs <-chan Job) {
    for {
        select {
        case job := <-jobs:
            processJob(job)
        case <-ctx.Done():
            return
        }
    }
}
```

### Channel Safety
```go
// Don't close channels from receivers
// Only close from senders

// Use select for non-blocking sends
select {
case ch <- value:
    // Sent successfully
default:
    // Channel is full, handle accordingly
}

// Check for closed channels
for {
    select {
    case value, ok := <-ch:
        if !ok {
            // Channel is closed
            return
        }
        process(value)
    }
}
```

### Memory Management
```go
// Use buffered channels to prevent blocking
eventQueue := make(chan *GameEvent, 100)

// Limit goroutine creation
maxWorkers := 10
workers := make(chan struct{}, maxWorkers)

// Use object pools for frequently allocated objects
var eventPool = sync.Pool{
    New: func() interface{} {
        return &GameEvent{}
    },
}
```

## 9. Common Pitfalls

### Goroutine Leaks
```go
// BAD: Goroutine that never exits
go func() {
    for {
        // This runs forever
        processData()
    }
}()

// GOOD: Goroutine with exit condition
go func() {
    defer wg.Done()
    for {
        select {
        case <-ctx.Done():
            return
        case data := <-dataCh:
            processData(data)
        }
    }
}()
```

### Race Conditions
```go
// BAD: Unsafe concurrent access
type Counter struct {
    count int
}

func (c *Counter) Increment() {
    c.count++ // Race condition!
}

// GOOD: Safe concurrent access
type Counter struct {
    count int
    mutex sync.Mutex
}

func (c *Counter) Increment() {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    c.count++
}
```

### Channel Deadlocks
```go
// BAD: Potential deadlock
ch := make(chan int)
ch <- 1  // Blocks forever if no receiver

// GOOD: Use select with default
select {
case ch <- 1:
    // Sent successfully
default:
    // Handle full channel
}
```

## 10. Testing Concurrent Code

### Testing Goroutines
```go
func TestWorkerPool(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    jobs := make(chan int, 10)
    results := make(chan int, 10)
    
    // Start workers
    var wg sync.WaitGroup
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go worker(ctx, jobs, results, &wg)
    }
    
    // Send jobs
    go func() {
        defer close(jobs)
        for i := 0; i < 5; i++ {
            jobs <- i
        }
    }()
    
    // Wait for completion
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // Collect results
    var results []int
    for result := range results {
        results = append(results, result)
    }
    
    if len(results) != 5 {
        t.Errorf("Expected 5 results, got %d", len(results))
    }
}
```

**Key Points:**
- Use timeouts in tests
- Test both success and failure scenarios
- Use race detector (`go test -race`)
- Test cancellation and cleanup
