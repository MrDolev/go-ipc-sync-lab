# go-ipc-sync-lab

Practical Go implementations of IPC and synchronization primitives.

This project originated from project work for the **Fondamenti di Informatica** module (Bachelor's Degree), where I contributed to improving the theoretical descriptions of IPC mechanisms on the Italian Wikipedia page for [Comunicazione tra processi](https://it.wikipedia.org/wiki/Comunicazione_tra_processi).

---

## 🟢 Phase 1: Foundations

### Beginner On-Ramp: What is a Goroutine?
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

## 🔵 Phase 2: The Toolkit & Language Bridge

### 1. The Go Concurrency Toolkit
Before looking at the patterns, here are the four main tools we use in Go:

1.  **`sync.Mutex` (The Lock)**: Used to protect a single piece of data. Only one goroutine can hold the lock at a time.
2.  **`sync.WaitGroup` (The Counter)**: Used to wait for a collection of goroutines to finish. It's like a guest list; the program won't close until everyone has "checked out."
3.  **Unbuffered Channel (`chan T`)**: A pipe with **zero capacity**. It forces a "Direct Handover" where both the sender and receiver must be ready at the same instant.
4.  **Buffered Channel (`chan T, n`)**: A pipe with a **fixed capacity**. It acts like a queue, allowing the sender to keep working until the pipe is full.

### 2. The Signal Value: `struct{}{}`
In the semaphore pattern, you'll see `make(chan struct{}, n)`.

`struct{}` is a **zero-byte type**. It carries no data and uses no memory. Think of it as a **Signal Light**: we don't care *what* color the light is, only that it is **ON** (a value is in the channel) or **OFF** (the channel is empty). We use it to notify the system that a slot is taken without wasting memory on unnecessary data.

### 3. Conceptual Map for Polyglot Developers
If you already know threads and locks from Python, Java, or C++, use this translation layer.

<details>
<summary><b>View Language Mappings</b></summary>

| Feature | Thread (Python/Java/C++) | Goroutine (Go) |
|---|---|---|
| **Spawn Cost** | ~1–8 MB stack | ~2 KB stack |
| **Scheduling** | OS Kernel (Expensive) | Go Runtime (Cheap) |
| **Creation** | `new Thread(function).start()` | `go function()` |

#### Python Developers
| Python | Go equivalent | Key Difference |
|---|---|---|
| `threading.Lock()` | `sync.Mutex` | Use `defer mu.Unlock()` for the same safety as `with lock:`. |
| `threading.Semaphore(n)` | `make(chan struct{}, n)` | No special type; a buffered channel *is* the semaphore. |

#### Java Developers
| Java | Go equivalent | Key Difference |
|---|---|---|
| `synchronized(this)` | `mu.Lock()` / `mu.Unlock()` | Always explicit; no method-level locking sugar. |
| `BlockingQueue` | `chan T` | Channels are first-class citizens in the language. |

#### C++ Developers
| C++ | Go equivalent | Key Difference |
|---|---|---|
| `std::lock_guard` | `defer mu.Unlock()` | `defer` is Go's RAII (runs at function exit). |
| `std::thread::join()` | `sync.WaitGroup` | Synchronization is handled via WaitGroups or channels. |

</details>

---

## 🟡 Phase 3: Visualizing Concurrency Patterns

| Pattern | Concept (Analogy) | Go Primitive |
|---|---|---|
| **Mutex** | **The Key**: Only one person can enter the room at a time. | `sync.Mutex` |
| **Semaphore** | **The Parking Lot**: Limited spaces; wait if the lot is full. | `chan struct{}` |
| **Prod-Cons** | **The Phone Call**: Direct, synchronous data handover. | `chan T` |

> **Pro-Tip**: Use **Channels** to distribute work or orchestrate data flow (Messaging). Use **Mutexes** to protect a single piece of shared data (Memory Sharing).

### 1. Mutex (Mutual Exclusion)
> **The Concept**: Think of a single **Key** to a restricted room. To enter and perform work, you must hold the key. If someone else has it, you must wait in line.

*   **The Problem**: **Data Races**. When multiple threads (goroutines) try to write to the same memory location at once, the data becomes corrupted.
*   **The Solution**: The `sync.Mutex` ensures **Atomic Access**. It locks the data so only one worker can touch it at a time, making the operation "thread-safe."
*   **Avoiding Deadlocks**: Always use `defer mu.Unlock()` immediately after `mu.Lock()`. This ensures the key is returned even if the function crashes or returns early.
*   **Real-World Scenario**: Imagine 5,000 concurrent users trying to buy the last seat on a flight; a Mutex ensures only one user succeeds and the seat isn't sold twice.

![Mutex](./docs/img/mutex.png)

**Key Logic**:
```go
// Safe update across all goroutines
mu.Lock()
defer mu.Unlock()
count++ 
```

### 2. Semaphore (Bounded Concurrency)
> **The Concept**: Think of a **Parking Lot** with 10 spaces. Cars can enter as long as there is a free spot. If the lot is full, new cars must wait until a space is vacated.

*   **The Solution**: A **Buffered Channel** acts as a counter. By setting a capacity (N), you strictly limit how many goroutines are allowed to run the "heavy" code at once.
*   **Avoiding Starvation**: Use `defer func() { <-sem }()` immediately after acquiring a slot. This guarantees that even if the worker panics, the slot is returned to the pool.
*   **Real-World Scenario**: Imagine 10,000 workers needing to call an external API that only allows 50 concurrent requests.

![Semaphore](./docs/img/semaphore.png)

**Key Logic**:
```go
sem := make(chan struct{}, 3) 

for i := 0; i < 10; i++ {
    go func() {
        sem <- struct{}{}        // Acquire
        defer func() { <-sem }() // Release
        doWork()
    }()
}
```

### 3. Producer-Consumer (Synchronous Handover)
> **The Concept**: Think of a **Phone Call**. A conversation can only happen if both the caller and the receiver are on the line at the same time.

*   **The Solution**: An **Unbuffered Channel** (Capacity 0). It forces a **Direct Handover** (where both sides must meet at the same time) for the data packet.
*   **Avoiding Deadlocks**: The Producer *must* `close(ch)` when finished. If the channel is never closed, the Consumer will wait forever (deadlock).
*   **Real-World Scenario**: Imagine a log-processing pipeline where one service generates logs and another pushes them to a database.

![Producer-Consumer](./docs/img/prod-cons.png)

**Key Logic**:
```go
ch := make(chan string) // Zero capacity (Direct Handover)

go func() {
    ch <- "Signal" 
    close(ch)
}()

for msg := range ch {
    fmt.Println("Received:", msg)
}
```

---

## 🔴 Phase 4: Defensive Engineering

### Common Pitfalls: How Deadlocks Happen
A **Deadlock** occurs when goroutines are stuck waiting for each other, and none can proceed.

| Pitfall | Consequence | How to Fix |
|---|---|---|
| **Forgetting to Unlock** | Other goroutines wait forever for the Mutex. | Use `defer mu.Unlock()` immediately after locking. |
| **Forgetting to Close** | The Consumer `range` loop waits forever for more data. | Always `close(ch)` in the Producer when done. |
| **Circular Waiting** | G1 waits for G2, and G2 waits for G1. | Always acquire multiple locks in the same consistent order (A then B). |
| **Sending to Full Channel** | If no one is receiving, the Producer blocks forever. | Ensure a Consumer is active or use a large enough buffer. |

---

## 🛠️ Phase 5: Usage & Environment

### Prerequisites
*   **Go 1.24+**, **Make**, and **Git**.

### Execution
```bash
make check         # Verify your environment
make install-tools # Install govulncheck
make run           # Build and run (or: go run ./cmd/main.go)
```

### Verification
```bash
make test   # Run unit tests (or: go test -v ./...)
make race   # Run with Race Detector (or: go run -race ./cmd/main.go)
make vuln   # Check for vulnerabilities (or: govulncheck ./...)
```

---

## 📚 Phase 6: References & Further Learning

*   **🎓 Tutorials**: [Official Go Tour](https://go.dev/tour/concurrency/1) | [Go by Example](https://gobyexample.com/mutexes) | [Ardan Labs: Channels](https://www.ardanlabs.com/blog/2014/02/the-nature-of-channels-in-go.html)
*   **🛡️ Safety**: [Race Detector Guide](https://yourbasic.org/golang/detect-data-races/) | [Common Mistakes](https://go101.org/article/concurrent-common-mistakes.html)
*   **🛠️ Tooling**: [Golang Makefile Guide](https://earthly.dev/blog/golang-makefile/) | [Time-Saving Makefiles](https://www.alexedwards.net/blog/a-time-saving-makefile-for-your-go-projects)

---
*This lab is a living document. See the `pkg/` directory for implementation details of each pattern.*
