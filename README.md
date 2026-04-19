# go-ipc-sync-lab

Practical Go implementations of IPC and synchronization primitives.

This project originated from project work for the **Fondamenti di Informatica** module (Bachelor's Degree), where I contributed to improving the theoretical descriptions of IPC mechanisms on the Italian Wikipedia page for [Comunicazione tra processi](https://it.wikipedia.org/wiki/Comunicazione_tra_processi).

---

## Patterns

| Pattern | Concept (Analogy) | Go Primitive |
|---|---|---|
| **Mutex** | **The Key**: Only one person can enter the room at a time. | `sync.Mutex` |
| **Semaphore** | **The Parking Lot**: Limited spaces; wait if the lot is full. | `chan struct{}` |
| **Prod-Cons** | **The Phone Call**: Direct, synchronous data handover. | `chan T` |

---

## Visualizing Concurrency

### 1. Mutex (Mutual Exclusion)
> **The Concept**: Think of a single **Key** to a restricted room. To enter and perform work, you must hold the key. If someone else has it, you must wait in line.

*   **The Problem**: **Data Races**. When multiple threads (goroutines) try to write to the same memory location at once, the data becomes corrupted.
*   **The Solution**: The `sync.Mutex` ensures **Atomic Access**. It locks the data so only one worker can touch it at a time, making the operation "thread-safe."

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
count++ 
mu.Unlock()
```

*   **Goal**: Prevent data corruption in shared state.

### 2. Semaphore (Bounded Concurrency)
> **The Concept**: Think of a **Parking Lot** with 10 spaces. Cars can enter as long as there is a free spot. If the lot is full, new cars must wait until a space is vacated.

*   **The Problem**: **Resource Overload**. Running too many concurrent tasks (like 5,000 database queries) simultaneously can exhaust system memory or CPU.
*   **The Solution**: A **Buffered Channel** acts as a counter. By setting a capacity (N), you strictly limit how many goroutines are allowed to run the "heavy" code at once.

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
*   **The Solution**: An **Unbuffered Channel** (Capacity 0). It forces a "Rendezvous" where the Producer and Consumer meet to hand over the data packet directly.

![Producer-Consumer](./docs/img/prod-cons.png)

**Key Logic**:
*   **Send**: `ch <- data` — Wait for a receiver to pick up.
*   **Receive**: `data := <-ch` — Wait for a sender to dial in.
*   **Close**: `close(ch)` — Hang up the line; signals the end of communication.

**Example Code**:
```go
ch := make(chan string) // Zero capacity

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

## Usage

The project uses a `Makefile` to provide a consistent interface for development.

### Development & Execution
```bash
make run    # Build and run the concurrency simulation
make fmt    # Format all Go source files (idiomatic style)
make build  # Compile the binary without running
```

### Testing & Verification
```bash
make test   # Run all unit tests and race condition checks
make race   # Specifically run the code with the Race Detector
make vet    # Run 'go vet' to identify potential static errors
make vuln   # Check for known vulnerabilities (requires govulncheck)
```

### Maintenance
```bash
make clean  # Remove build artifacts and binary files
```
