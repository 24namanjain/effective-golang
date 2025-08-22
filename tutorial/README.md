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
# Copy file to have _test.go suffix (required for benchmarks)
cp benchmarks.go benchmarks_test.go

# Run all benchmarks
go test -bench=. benchmarks_test.go

# Run specific benchmark
go test -bench=BenchmarkStringConcatenation benchmarks_test.go

# Run with memory info
go test -bench=. benchmarks_test.go -benchmem

# Clean up
rm benchmarks_test.go
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

**Example Results:**
```
BenchmarkStringConcatenation/Strings_Join-8    24150760    44.30 ns/op
BenchmarkStringConcatenation/Strings_Builder-8  14668633    80.19 ns/op
BenchmarkStringConcatenation/String_Plus-8     11230993   116.6 ns/op
```

## 🎯 Key Commands

```bash
# Build all files
go build .

# Run specific files
go run file1.go file2.go

# Run all tutorials
go run .

# Test specific function
go test -run TestFunctionName

# Run benchmarks (requires _test.go suffix)
cp benchmarks.go benchmarks_test.go
go test -bench=BenchmarkName benchmarks_test.go
rm benchmarks_test.go

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
- Use `_test.go` suffix for benchmark files

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

## 📊 Benchmark Examples

### String Concatenation Performance
- `strings.Join`: Fastest (44.30 ns/op)
- `strings.Builder`: 2x faster than `+` (80.19 ns/op)
- `+` concatenation: Slowest (116.6 ns/op)

### Memory Allocation
- **Reuse allocation**: 3x faster than new allocation
- **Memory reuse**: Eliminates allocations (0 B/op vs 8192 B/op)

### Concurrency vs Sequential
- **Sequential**: Often faster for simple operations
- **Concurrency overhead**: Can make simple tasks slower

## 🚨 Common Issues

- **Undefined function**: Run with `go run .` to include all files
- **Test not found**: Use `go test -run TestName`
- **Benchmark not found**: Copy file to `_test.go` suffix first
- **Import errors**: Check `go.mod` and run `go mod tidy`
- **"no test files" error**: Use `_test.go` suffix for benchmark files

## 📚 Resources

- [Go Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Testing](https://golang.org/pkg/testing/)
- [Go Concurrency](https://golang.org/doc/effective_go.html#concurrency)
