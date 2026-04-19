# go-ipc-sync-lab

Practical Go implementations of IPC and synchronization primitives.

This project originated from project work for the **Fondamenti di Informatica** module (Bachelor's Degree), where I contributed to improving the theoretical descriptions of IPC mechanisms on the Italian Wikipedia page for [Comunicazione tra processi](https://it.wikipedia.org/wiki/Comunicazione_tra_processi).

---

## Patterns

| Pattern | Concept (Analogy) | Go Primitive |
|---|---|---|
| **Mutex** | **The Single Key**: Only one person can enter the room at a time. | `sync.Mutex` |
| **Semaphore** | **The Admission Ticket**: Limit the total number of guests in the pool. | `chan struct{}` |
| **Prod-Cons** | **The Handshake**: A direct, synchronous handover of data. | `chan T` |

---

## Visualizing Concurrency

### 1. Mutex (Mutual Exclusion)
> **The Concept**: Think of a single key to a bathroom. If you have the key, you can enter; if not, you must wait in line until the current occupant comes out and hands you the key.

*   **The Problem**: **Race Conditions**. Multiple goroutines trying to modify the same data simultaneously, leading to unpredictable results and corrupted state.
*   **The Solution**: Use a `sync.Mutex` to protect the **Critical Section**. It ensures that only one goroutine can execute that block of code at any given time.

![Mutex](./docs/img/mutex.png)

**Key Logic**:
*   `mu.Lock()`: Acquire the key. If someone else has it, you block (wait).
*   `mu.Unlock()`: Release the key. The next goroutine in line can now take it.

**Example Code**:
```go
var (
    mu    sync.Mutex
    count int
)

// Thread-safe increment
mu.Lock()
count++ 
mu.Unlock()
```

### 2. Semaphore (Bounded Concurrency)
> **The Concept**: Think of a spot with a maximum capacity. spot (the semaphore) only let a new guest in (goroutine) when someone else leaves.

*   **The Problem**: **Resource Exhaustion**. Launching too many concurrent tasks (e.g., thousands of API calls) can crash your system or get you rate-limited.
*   **The Solution**: Use a **Buffered Channel** as a counting semaphore. The channel's capacity defines the "limit" of allowed concurrent workers.

![Semaphore](./docs/img/semaphore.png)

**Key Logic**:
*   **Acquire**: `sem <- struct{}{}` (Fill a slot). Blocks if the "club" is full.
*   **Release**: `<-sem` (Free a slot). Someone left the club.

**Example Code**:
```go
// Limit to 3 concurrent workers
sem := make(chan struct{}, 3)

for i := 0; i < 10; i++ {
    go func() {
        sem <- struct{}{} // Acquire
        defer func() { <-sem }() // Release
        doWork()
    }()
}
```

### 3. Producer-Consumer (The Synchronous Pipeline)
> **The Concept**: Think of a physical handshake. You can't complete the handshake unless both people are present at the same time and reach out.

*   **The Problem**: **Tight Coupling**. One part of your system generates data while another processes it. You need a way to pass data safely without them knowing too much about each other.
*   **The Solution**: Use an **Unbuffered Channel** as a synchronous "Conveyor Belt" or "Pipe". It enforces a perfect rendezvous between the Producer and the Consumer.

![Producer-Consumer](./docs/img/prod-cons.png)

**Key Logic**:
*   **Send**: `ch <- data` — The Producer offers an item and waits for a receiver.
*   **Receive**: `data := <-ch` — The Consumer waits for an item to arrive.
*   **Close**: `close(ch)` — The Producer signals that the "shift is over."

**Example Code**:
```go
ch := make(chan string) // Unbuffered

// Producer
go func() {
    ch <- "Data Packet" 
    close(ch)
}()

// Consumer
for item := range ch {
    fmt.Println("Received:", item)
}
```

---

## Usage

The project includes a `Makefile` to simplify common tasks.

```bash
# Run the simulation (builds and executes)
make run

# Run all tests (includes race detection)
make test

# Check for race conditions specifically
make race

# Format the codebase
make fmt

# Check for vulnerabilities (requires govulncheck)
make vuln

# Clean build artifacts
make clean
```
