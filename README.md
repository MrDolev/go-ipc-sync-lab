# go-ipc-sync-lab

Practical Go implementations of IPC and synchronization primitives.

This project originated from project work for the **Fondamenti di Informatica** module (Bachelor's Degree), where I contributed to improving the theoretical descriptions of IPC mechanisms on the Italian Wikipedia page for [Comunicazione tra processi](https://it.wikipedia.org/wiki/Comunicazione_tra_processi).

> **Coming from another language?** See how Go primitives map to [Python, Java, and C++](#conceptual-map-for-polyglot-developers).

---

## Beginner On-Ramp: What is a Goroutine?

Before diving into synchronization, you need to know about the **Goroutine** — the building block of Go concurrency.

*   **Definition**: A goroutine is a function that runs concurrently with other functions.
*   **Analogy**: If a standard program is a single chef cooking one meal at a time, a program with goroutines is like having **multiple chefs** in the same kitchen.
*   **Syntax**: Just add the word `go` before any function call.

```go
func sayHello() {
    fmt.Println("Hello!")
}

func main() {
    go sayHello() // This runs in the background!
    fmt.Println("Main function continues...")
    
    // Important: If we don't wait, the program might finish 
    // before the background chef (goroutine) can say "Hello!"
    time.Sleep(time.Millisecond * 10) 
}
```

**Why do we need synchronization?**
In the example above, `time.Sleep` is a "hack." In real programs, we use the patterns below (Mutex, Semaphore, Channels) to ensure goroutines coordinate perfectly without guessing how much time they need. When multiple "chefs" (goroutines) share the same "ingredients" (data), these patterns are the **rules of the kitchen** that keep everything running safely.

---

## The Go Concurrency Toolkit

Before looking at the patterns, here are the four main tools we use in Go:

1.  **`sync.Mutex` (The Lock)**: Used to protect a single piece of data. Only one goroutine can hold the lock at a time.
2.  **`sync.WaitGroup` (The Counter)**: Used to wait for a collection of goroutines to finish. It's like a guest list; the program won't close until everyone has "checked out."
3.  **Unbuffered Channel (`chan T`)**: A pipe with **zero capacity**. It forces a "Direct Handover" where both the sender and receiver must be ready at the same instant.
4.  **Buffered Channel (`chan T, n`)**: A pipe with a **fixed capacity**. It acts like a queue, allowing the sender to keep working until the pipe is full.

---

## Patterns

| Pattern | Concept (Analogy) | Go Primitive |
|---|---|---|
| **Mutex** | **The Key**: Only one person can enter the room at a time. | `sync.Mutex` |
| **Semaphore** | **The Parking Lot**: Limited spaces; wait if the lot is full. | `chan struct{}` |
| **Prod-Cons** | **The Phone Call**: Direct, synchronous data handover. | `chan T` |

> **Pro-Tip**: Use **Channels** to distribute work or orchestrate data flow (Messaging). Use **Mutexes** to protect a single piece of shared data (Memory Sharing).

---

## Conceptual Map for Polyglot Developers

If you already know threads and locks from Python, Java, or C++, you don't need to learn Go from scratch — you just need a translation layer.

### 1. The Core Unit: Goroutines
A **goroutine** is a lightweight thread managed by the Go runtime.

| Feature | Thread (Python/Java/C++) | Goroutine (Go) |
|---|---|---|
| **Spawn Cost** | ~1–8 MB stack | ~2 KB stack |
| **Scheduling** | OS Kernel (Expensive) | Go Runtime (Cheap) |
| **Creation** | `new Thread(fn).start()` | `go function()` |

### 2. Language-Specific Mapping

<details>
<summary><b>Python Developers</b></summary>

| Python | Go equivalent | Key Difference |
|---|---|---|
| `threading.Lock()` | `sync.Mutex` | Use `defer mu.Unlock()` for the same safety as `with lock:`. |
| `threading.Semaphore(n)` | `make(chan struct{}, n)` | No special type; a buffered channel *is* the semaphore. |
| `queue.Queue()` | `chan T` | Built-in primitive; unbuffered = synchronous handover. |

</details>

<details>
<summary><b>Java Developers</b></summary>

| Java | Go equivalent | Key Difference |
|---|---|---|
| `synchronized(this)` | `mu.Lock()` / `mu.Unlock()` | Always explicit; no method-level locking sugar. |
| `new Semaphore(n)` | `make(chan struct{}, n)` | No `java.util.concurrent` import needed. |
| `BlockingQueue` | `chan T` | Channels are first-class citizens in the language. |

</details>

<details>
<summary><b>C++ Developers</b></summary>

| C++ | Go equivalent | Key Difference |
|---|---|---|
| `std::lock_guard` | `defer mu.Unlock()` | `defer` is Go's RAII (runs at function exit). |
| `std::condition_variable` | `chan struct{}` | No spurious wakeups; no predicate loops required. |
| `std::thread::join()` | `sync.WaitGroup` | Synchronization is handled via WaitGroups or channels. |

</details>

### 3. The Signal Value: `struct{}{}`
In the semaphore pattern, you'll see `make(chan struct{}, n)`.

`struct{}` is a **zero-byte type**. It carries no data and uses no memory. Think of it as a **Signal Light**: we don't care *what* color the light is, only that it is **ON** (a value is in the channel) or **OFF** (the channel is empty). We use it to notify the system that a slot is taken without wasting memory on unnecessary data.

---

## Visualizing Concurrency

### 1. Mutex (Mutual Exclusion)
> **The Concept**: Think of a single **Key** to a restricted room. To enter and perform work, you must hold the key. If someone else has it, you must wait in line.

*   **The Problem**: **Data Races**. When multiple threads (goroutines) try to write to the same memory location at once, the data becomes corrupted.
*   **The Solution**: The `sync.Mutex` ensures **Atomic Access**. It locks the data so only one worker can touch it at a time, making the operation "thread-safe."
*   **Avoiding Deadlocks**: Always use `defer mu.Unlock()` immediately after `mu.Lock()`. This ensures the key is returned even if the function crashes or returns early.
*   **Real-World Scenario**: Imagine 5,000 concurrent users trying to buy the last seat on a flight; a Mutex ensures only one user succeeds and the seat isn't sold twice.

![Mutex](./docs/img/mutex.png)

**Key Logic**:
*   `mu.Lock()`: Acquire the key (blocks if unavailable).
*   `mu.Unlock()`: Return the key (signals the next waiter).

**Example Code**:
```go
var (
    mu    sync.Mutex
    count int
)

// Safe update across all goroutines
mu.Lock()
defer mu.Unlock()
count++ 
```

*   **Goal**: Prevent data corruption in shared state.

### 2. Semaphore (Bounded Concurrency)
> **The Concept**: Think of a **Parking Lot** with 10 spaces. Cars can enter as long as there is a free spot. If the lot is full, new cars must wait until a space is vacated.

*   **The Solution**: A **Buffered Channel** acts as a counter. By setting a capacity (N), you strictly limit how many goroutines are allowed to run the "heavy" code at once.
*   **Avoiding Starvation**: Use `defer func() { <-sem }()` immediately after acquiring a slot. This guarantees that even if the worker panics, the slot is returned to the pool, preventing others from waiting forever.
*   **Real-World Scenario**: Imagine 10,000 workers needing to call an external API that only allows 50 concurrent requests; a Semaphore prevents your app from being blocked or banned.

![Semaphore](./docs/img/semaphore.png)

**Key Logic**:
*   **Acquire**: `sem <- struct{}{}` — Take a space in the lot.
*   **Release**: `<-sem` — Leave the space for the next car.

**Example Code**:
```go
// Limit to 3 concurrent workers
sem := make(chan struct{}, 3) // struct{}{} is a zero-byte token

for i := 0; i < 10; i++ {
    go func() {
        sem <- struct{}{}        // Acquire slot
        defer func() { <-sem }() // Always release slot
        doWork()
    }()
}
```

*   **Goal**: Limit resource usage and prevent system overload.

### 3. Producer-Consumer (Synchronous Handover)
> **The Concept**: Think of a **Phone Call**. A conversation can only happen if both the caller and the receiver are on the line at the same time. No data is stored "in the middle."

*   **The Problem**: **Execution Coordination**. You need to pass data from a "source" to a "sink" without creating a permanent storage buffer, ensuring both sides are perfectly synchronized.
*   **The Solution**: An **Unbuffered Channel** (Capacity 0). It forces a **Direct Handover** (where both sides must meet at the same time) for the data packet.
*   **Avoiding Deadlocks**: The Producer *must* `close(ch)` when finished. If the channel is never closed, the Consumer's `for msg := range ch` loop will wait forever (deadlock) after the last item is sent.
*   **Real-World Scenario**: Imagine a log-processing pipeline where one service generates logs and another pushes them to a database; a Channel ensures logs are passed safely without needing a complex database locking mechanism.

![Producer-Consumer](./docs/img/prod-cons.png)

**Key Logic**:
*   **Send**: `ch <- data` — Wait for a receiver to pick up.
*   **Receive**: `data := <-ch` — Wait for a sender to dial in.
*   **Close**: `close(ch)` — Hang up the line; signals the end of communication.

**Example Code**:
```go
ch := make(chan string) // Zero capacity (Direct Handover)

// Producer (Sender)
go func() {
    ch <- "Signal" 
    close(ch)
}()

// Consumer (Receiver)
for msg := range ch {
    fmt.Println("Received:", msg)
}
```

*   **Goal**: Decouple logic and synchronize flow without intermediate storage.

---

## Common Pitfalls: How Deadlocks Happen

To truly understand synchronization, you must understand how it fails. A **Deadlock** occurs when goroutines are stuck waiting for each other, and none can proceed.

| Pitfall | Consequence | How to Fix |
|---|---|---|
| **Forgetting to Unlock** | Other goroutines wait forever for the Mutex. | Use `defer mu.Unlock()` immediately after locking. |
| **Forgetting to Close** | The Consumer `range` loop waits forever for more data. | Always `close(ch)` in the Producer when done. |
| **Circular Waiting** | G1 waits for G2, and G2 waits for G1. | Always acquire multiple locks in the same consistent order (e.g., always Lock A then Lock B). |
| **Sending to Full Channel** | If no one is receiving, the Producer blocks forever. | Ensure a Consumer is always active or use a large enough buffer. |

---

## Prerequisites

To run this project, you need the following tools installed:

*   **Go 1.24+**: The programming language used for the implementation.
*   **Make**: A build automation tool used to run the tasks defined in the `Makefile`.
*   **Git**: For version control and repository management.

---

## Usage

The project uses a `Makefile` to provide a consistent interface for development.

### Setup & Verification
```bash
make check         # Verify your environment (Go, Make, tools)
make install-tools # Install required development tools (govulncheck)
make run           # Build and run the concurrency simulation (or: go run ./cmd/main.go)
```

### Development & Execution
```bash
make fmt    # Format all Go source files (or: go fmt ./...)
make build  # Compile the binary without running (or: go build -o bin/go-ipc ./cmd)
```

### Testing & Verification
```bash
make test   # Run all tests and race checks (or: go test -v ./...)
make race   # Specifically run with the Race Detector (or: go run -race ./cmd/main.go)
make vet    # Run 'go vet' (or: go vet ./...)
make vuln   # Check for vulnerabilities (or: govulncheck ./...)
```

### Maintenance
```bash
make clean  # Remove build artifacts and binary files
```

---

## References & Further Learning

### 🎓 Go Concurrency Tutorials
*   **[A Tour of Go: Goroutines](https://go.dev/tour/concurrency/1)** - The official interactive introduction to concurrent Go.
*   **[A Tour of Go: Channels](https://go.dev/tour/concurrency/2)** - Learn the mechanics of unbuffered and buffered channels.
*   **[The Nature of Channels in Go](https://www.ardanlabs.com/blog/2014/02/the-nature-of-channels-in-go.html)** - A deep dive into the philosophy and mechanics of Go channels.
*   **[Go by Example: Mutexes](https://gobyexample.com/mutexes)** - Hands-on examples of state protection.

### 🛡️ Safety & Best Practices
*   **[Detecting Data Races](https://yourbasic.org/golang/detect-data-races/)** - How to use the Go Race Detector effectively.
*   **[Common Concurrency Mistakes](https://go101.org/article/concurrent-common-mistakes.html)** - A deep dive into traps like leaks, panics, and deadlocks.
*   **[Context Package](https://pkg.go.dev/context)** - For production apps, learn how to handle timeouts and cancellations.

### 🛠️ Tooling & Project Structure
*   **[Golang Makefile Guide](https://earthly.dev/blog/golang-makefile/)** - Best practices for writing Makefiles for Go projects.
*   **[A Time-Saving Makefile for Go](https://www.alexedwards.net/blog/a-time-saving-makefile-for-your-go-projects)** - Practical tips for automating your workflow.

---
*This lab is a living document. See the `pkg/` directory for implementation details of each pattern.*
