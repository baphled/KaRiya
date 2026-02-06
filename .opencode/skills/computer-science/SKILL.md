# Computer Science Skill

You are a computer science expert with deep knowledge of data structures, algorithms, system design, and computational theory.

## Overview

Computer science fundamentals are the bedrock of software engineering. Master these concepts to write better, more efficient code.

---

## Data Structures

### Arrays & Slices

```go
// Fixed array
var arr [5]int

// Dynamic slice
slice := make([]int, 0, 10)  // len=0, cap=10
slice = append(slice, 1, 2, 3)

// Time complexity
arr[i]           // O(1) - access
slice = append() // O(1) amortized, O(n) worst (reallocation)
```

### Linked Lists

```go
type Node struct {
    Value int
    Next  *Node
}

type LinkedList struct {
    Head *Node
}

// Operations
// Insert at head: O(1)
// Insert at tail: O(n) or O(1) with tail pointer
// Search: O(n)
// Delete: O(n)
```

### Stacks (LIFO)

```go
type Stack struct {
    items []int
}

func (s *Stack) Push(item int) {
    s.items = append(s.items, item)
}

func (s *Stack) Pop() (int, bool) {
    if len(s.items) == 0 {
        return 0, false
    }
    item := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return item, true
}

// All operations: O(1)
```

### Queues (FIFO)

```go
type Queue struct {
    items []int
}

func (q *Queue) Enqueue(item int) {
    q.items = append(q.items, item)
}

func (q *Queue) Dequeue() (int, bool) {
    if len(q.items) == 0 {
        return 0, false
    }
    item := q.items[0]
    q.items = q.items[1:]
    return item, true
}

// Enqueue: O(1) amortized
// Dequeue: O(n) - consider ring buffer for O(1)
```

### Hash Maps

```go
// Go's built-in map
m := make(map[string]int)
m["key"] = value

// Time complexity (average case)
m[key]          // O(1) - lookup
m[key] = value  // O(1) - insert
delete(m, key)  // O(1) - delete

// Worst case (hash collisions): O(n)
```

### Trees

```go
type TreeNode struct {
    Value int
    Left  *TreeNode
    Right *TreeNode
}

// Binary Search Tree
// Search, Insert, Delete: O(log n) average, O(n) worst
// Balanced BST (AVL, Red-Black): O(log n) guaranteed

// Traversals
func inOrder(node *TreeNode, result *[]int) {
    if node == nil {
        return
    }
    inOrder(node.Left, result)
    *result = append(*result, node.Value)
    inOrder(node.Right, result)
}

func preOrder(node *TreeNode, result *[]int) {
    if node == nil {
        return
    }
    *result = append(*result, node.Value)
    preOrder(node.Left, result)
    preOrder(node.Right, result)
}

func postOrder(node *TreeNode, result *[]int) {
    if node == nil {
        return
    }
    postOrder(node.Left, result)
    postOrder(node.Right, result)
    *result = append(*result, node.Value)
}
```

### Heaps

```go
import "container/heap"

// Min-heap implementation
type MinHeap []int

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x interface{}) {
    *h = append(*h, x.(int))
}

func (h *MinHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}

// Operations
// Insert: O(log n)
// Extract min/max: O(log n)
// Peek: O(1)
```

### Graphs

```go
// Adjacency list
type Graph struct {
    vertices map[int][]int
}

func NewGraph() *Graph {
    return &Graph{vertices: make(map[int][]int)}
}

func (g *Graph) AddEdge(from, to int) {
    g.vertices[from] = append(g.vertices[from], to)
}

// Adjacency matrix (for dense graphs)
type MatrixGraph struct {
    matrix [][]bool
    size   int
}
```

---

## Design Patterns

### Creational

```go
// Singleton
var instance *Config
var once sync.Once

func GetConfig() *Config {
    once.Do(func() {
        instance = &Config{}
    })
    return instance
}

// Factory
type Vehicle interface {
    Drive()
}

func CreateVehicle(vType string) Vehicle {
    switch vType {
    case "car":
        return &Car{}
    case "truck":
        return &Truck{}
    default:
        return nil
    }
}

// Builder
type ServerBuilder struct {
    server *Server
}

func (b *ServerBuilder) WithPort(port int) *ServerBuilder {
    b.server.Port = port
    return b
}

func (b *ServerBuilder) WithHost(host string) *ServerBuilder {
    b.server.Host = host
    return b
}

func (b *ServerBuilder) Build() *Server {
    return b.server
}
```

### Structural

```go
// Adapter
type OldSystem struct{}
func (o *OldSystem) OldMethod() string { return "old" }

type Adapter struct {
    old *OldSystem
}

func (a *Adapter) NewMethod() string {
    return a.old.OldMethod()
}

// Decorator
type Coffee interface {
    Cost() float64
}

type SimpleCoffee struct{}
func (c *SimpleCoffee) Cost() float64 { return 1.0 }

type MilkDecorator struct {
    coffee Coffee
}
func (m *MilkDecorator) Cost() float64 {
    return m.coffee.Cost() + 0.5
}
```

### Behavioral

```go
// Observer
type Observer interface {
    Update(data interface{})
}

type Subject struct {
    observers []Observer
}

func (s *Subject) Subscribe(o Observer) {
    s.observers = append(s.observers, o)
}

func (s *Subject) Notify(data interface{}) {
    for _, o := range s.observers {
        o.Update(data)
    }
}

// Strategy
type SortStrategy interface {
    Sort([]int) []int
}

type QuickSort struct{}
func (q *QuickSort) Sort(arr []int) []int { /* ... */ }

type MergeSort struct{}
func (m *MergeSort) Sort(arr []int) []int { /* ... */ }

type Sorter struct {
    strategy SortStrategy
}

func (s *Sorter) Sort(arr []int) []int {
    return s.strategy.Sort(arr)
}
```

---

## System Design Concepts

### Scalability

- **Vertical scaling**: Add more power (CPU, RAM)
- **Horizontal scaling**: Add more machines
- **Load balancing**: Distribute traffic
- **Caching**: Reduce database load
- **Database sharding**: Split data across databases

### CAP Theorem

Can only have 2 of 3:
- **Consistency**: All nodes see same data
- **Availability**: System always responds
- **Partition tolerance**: System works despite network failures

### ACID vs BASE

**ACID** (Traditional databases):
- Atomicity, Consistency, Isolation, Durability

**BASE** (NoSQL):
- Basically Available, Soft state, Eventually consistent

---

## Concurrency

### Concepts

```go
// Race condition - avoid with synchronization
var counter int
var mu sync.Mutex

func increment() {
    mu.Lock()
    counter++
    mu.Unlock()
}

// Deadlock - avoid by consistent lock ordering
// A locks X, waits for Y
// B locks Y, waits for X -> Deadlock!

// Channels - Go's preferred concurrency primitive
ch := make(chan int)
go func() { ch <- 42 }()
value := <-ch
```

### Patterns

```go
// Worker pool
func workerPool(jobs <-chan int, results chan<- int, workers int) {
    var wg sync.WaitGroup
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {
                results <- process(job)
            }
        }()
    }
    wg.Wait()
    close(results)
}

// Fan-out, fan-in
func fanOut(input <-chan int, workers int) []<-chan int {
    channels := make([]<-chan int, workers)
    for i := 0; i < workers; i++ {
        channels[i] = worker(input)
    }
    return channels
}

func fanIn(channels ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    out := make(chan int)
    
    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for v := range c {
                out <- v
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

---

## Computational Complexity

### P vs NP

- **P**: Problems solvable in polynomial time
- **NP**: Problems verifiable in polynomial time
- **NP-Complete**: Hardest problems in NP
- **NP-Hard**: At least as hard as NP-Complete

### Common NP-Complete Problems

- Traveling Salesman
- Knapsack Problem
- Graph Coloring
- Boolean Satisfiability (SAT)

---

## Memory Management

### Stack vs Heap

```go
// Stack allocation (fast, automatic)
func stackAlloc() {
    x := 42  // On stack
    // Automatically freed when function returns
}

// Heap allocation (slower, manual/GC)
func heapAlloc() *int {
    x := new(int)  // On heap
    *x = 42
    return x  // Escapes to heap
}
```

### Go's Garbage Collection

- Concurrent, tri-color mark-and-sweep
- Write barriers for concurrent marking
- Tunable via GOGC environment variable

---

## Related Skills

- `math-expert` - Mathematical foundations
- `performance` - Complexity and optimization
- `go-expert` - Go-specific implementations
