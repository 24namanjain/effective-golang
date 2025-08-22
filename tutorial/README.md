# Go Tutorial

A comprehensive Go tutorial covering fundamentals, concurrency, testing, and best practices.

## 🚀 Quick Start

```bash
# Run all tutorials
go run .

# Run specific tutorial
go run hello.go datatypes.go
```

## 📁 Tutorial Files

### Core Concepts
- **`hello.go`** - Variables, assignment, multiple file execution
- **`datatypes.go`** - Data types, arrays, structs, maps
- **`errorHandling.go`** - Error handling patterns
- **`jsonInGo.go`** - JSON encoding/decoding

### Advanced Topics
- **`concurrencyAndGoroutines.go`** - Concurrency, goroutines, context, worker pools
- **`channels.go`** - Channel patterns, pipelines, fan-out/fan-in
- **`table_driven_tests.go`** - Table-driven testing patterns
- **`benchmarks.go`** - Performance benchmarking

## 🧪 Testing

### Run Tests
```bash
# Run all tests
go test .

# Run specific test
go test -run TestReverseString

# Run with verbose output
go test -v .
```

### Run Benchmarks
```bash
# Run all benchmarks
go test -bench=. .

# Run specific benchmark
go test -bench=BenchmarkStringConcatenation .

# Run with memory info
go test -bench=. . -benchmem
```

## 🔗 Concurrency Examples

### Context & Goroutines
- Context cancellation and deadlines
- Goroutine leak prevention
- Mutex and channel protection
- Worker pools

### Channel Patterns
- Basic channel communication
- Buffered vs unbuffered channels
- Select statements
- Rate limiting
- Pipeline patterns
- Fan-out/Fan-in

## 📊 Testing Patterns

### Table-Driven Tests
```go
tests := []struct {
    name     string
    input    string
    expected string
}{
    {"empty", "", ""},
    {"simple", "hello", "olleh"},
}
```

### Benchmarks
```go
func BenchmarkFunction(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // Code to benchmark
    }
}
```

## 🎯 Key Commands

```bash
# Build all files
go build .

# Run specific files
go run file1.go file2.go

# Test specific function
go test -run TestFunctionName

# Benchmark with memory
go test -bench=. -benchmem

# Format code
go fmt .

# Vet code
go vet .
```

## 📝 Best Practices

### Concurrency
- Always use `context.Context` for cancellation
- Use `defer` for cleanup
- Protect shared resources with mutex or channels
- Use worker pools for high concurrency

### Testing
- Use table-driven tests for multiple cases
- Include edge cases and error conditions
- Use descriptive test names
- Benchmark performance-critical code

### Channels
- Close channels when done sending
- Use buffered channels when appropriate
- Use `select` for multiple channels
- Range over channels until closed

## 🔧 Project Structure

```
tutorial/
├── README.md
├── go.mod
├── hello.go              # Variables, assignment
├── datatypes.go          # Data types, arrays, structs
├── errorHandling.go      # Error handling
├── jsonInGo.go           # JSON operations
├── concurrencyAndGoroutines.go  # Concurrency patterns
├── channels.go           # Channel patterns
├── table_driven_tests.go # Testing patterns
└── benchmarks.go         # Performance testing
```

## 🎓 Learning Path

1. **Start with `hello.go`** - Basic syntax and variables
2. **Explore `datatypes.go`** - Data structures
3. **Learn `errorHandling.go`** - Error patterns
4. **Practice `jsonInGo.go`** - Data serialization
5. **Master `concurrencyAndGoroutines.go`** - Concurrency
6. **Understand `channels.go`** - Channel patterns
7. **Test with `table_driven_tests.go`** - Testing
8. **Optimize with `benchmarks.go`** - Performance

## 🚨 Common Issues

- **Undefined function**: Run with `go run .` to include all files
- **Test not found**: Use `go test -run TestName`
- **Benchmark not found**: Use `go test -bench=BenchmarkName`
- **Import errors**: Check `go.mod` and run `go mod tidy`

## 📚 Resources

- [Go Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Testing](https://golang.org/pkg/testing/)
- [Go Concurrency](https://golang.org/doc/effective_go.html#concurrency)
