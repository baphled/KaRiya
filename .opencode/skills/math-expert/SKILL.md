---
name: math-expert
description: Mathematical concepts, problem solving, and mathematical reasoning for software development
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

# Math Expert Skill

You are a mathematics expert who can explain mathematical concepts, solve problems, and apply mathematical reasoning to software development.

## Overview

Mathematics is the foundation of computer science. Apply rigorous mathematical thinking to algorithm design, complexity analysis, and problem-solving.

---

## Big-O Complexity

### Time Complexity

| Notation | Name | Example |
|----------|------|---------|
| O(1) | Constant | Map lookup, array index |
| O(log n) | Logarithmic | Binary search |
| O(n) | Linear | Linear search, single loop |
| O(n log n) | Linearithmic | Merge sort, heap sort |
| O(n²) | Quadratic | Nested loops, bubble sort |
| O(n³) | Cubic | Matrix multiplication |
| O(2ⁿ) | Exponential | Recursive fibonacci |
| O(n!) | Factorial | Permutations |

### Analyzing Code

```go
// O(1) - Constant
func getFirst(items []int) int {
    return items[0]
}

// O(n) - Linear
func findMax(items []int) int {
    max := items[0]
    for _, item := range items {  // n iterations
        if item > max {
            max = item
        }
    }
    return max
}

// O(n²) - Quadratic
func bubbleSort(items []int) {
    for i := 0; i < len(items); i++ {      // n iterations
        for j := 0; j < len(items)-1; j++ { // n iterations
            // ...
        }
    }
}

// O(log n) - Logarithmic
func binarySearch(items []int, target int) int {
    left, right := 0, len(items)-1
    for left <= right {
        mid := (left + right) / 2  // Halving each time
        if items[mid] == target {
            return mid
        } else if items[mid] < target {
            left = mid + 1
        } else {
            right = mid - 1
        }
    }
    return -1
}
```

### Space Complexity

```go
// O(1) space - In-place
func reverse(items []int) {
    for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
        items[i], items[j] = items[j], items[i]
    }
}

// O(n) space - Creates new slice
func duplicate(items []int) []int {
    result := make([]int, len(items))
    copy(result, items)
    return result
}

// O(n) space - Recursive call stack
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n-1)  // n stack frames
}
```

---

## Common Algorithms

### Sorting

```go
// Quick Sort - O(n log n) average, O(n²) worst
func quickSort(arr []int, low, high int) {
    if low < high {
        pivot := partition(arr, low, high)
        quickSort(arr, low, pivot-1)
        quickSort(arr, pivot+1, high)
    }
}

// Merge Sort - O(n log n) guaranteed
func mergeSort(arr []int) []int {
    if len(arr) <= 1 {
        return arr
    }
    mid := len(arr) / 2
    left := mergeSort(arr[:mid])
    right := mergeSort(arr[mid:])
    return merge(left, right)
}
```

### Searching

```go
// Binary Search - O(log n)
// Requires sorted array
func binarySearch(arr []int, target int) int {
    left, right := 0, len(arr)-1
    for left <= right {
        mid := left + (right-left)/2  // Avoid overflow
        if arr[mid] == target {
            return mid
        }
        if arr[mid] < target {
            left = mid + 1
        } else {
            right = mid - 1
        }
    }
    return -1
}
```

### Graph Algorithms

```go
// BFS - O(V + E)
func bfs(graph map[int][]int, start int) []int {
    visited := make(map[int]bool)
    queue := []int{start}
    result := []int{}
    
    for len(queue) > 0 {
        node := queue[0]
        queue = queue[1:]
        
        if visited[node] {
            continue
        }
        visited[node] = true
        result = append(result, node)
        
        for _, neighbor := range graph[node] {
            if !visited[neighbor] {
                queue = append(queue, neighbor)
            }
        }
    }
    return result
}

// DFS - O(V + E)
func dfs(graph map[int][]int, node int, visited map[int]bool) {
    if visited[node] {
        return
    }
    visited[node] = true
    
    for _, neighbor := range graph[node] {
        dfs(graph, neighbor, visited)
    }
}
```

---

## Number Theory

### Common Operations

```go
// GCD - Euclidean algorithm
func gcd(a, b int) int {
    for b != 0 {
        a, b = b, a%b
    }
    return a
}

// LCM
func lcm(a, b int) int {
    return a * b / gcd(a, b)
}

// Prime check
func isPrime(n int) bool {
    if n < 2 {
        return false
    }
    for i := 2; i*i <= n; i++ {
        if n%i == 0 {
            return false
        }
    }
    return true
}

// Modular exponentiation - (base^exp) % mod
func modPow(base, exp, mod int) int {
    result := 1
    base = base % mod
    for exp > 0 {
        if exp%2 == 1 {
            result = (result * base) % mod
        }
        exp = exp / 2
        base = (base * base) % mod
    }
    return result
}
```

---

## Statistics

### Basic Statistics

```go
import "math"

// Mean
func mean(data []float64) float64 {
    sum := 0.0
    for _, v := range data {
        sum += v
    }
    return sum / float64(len(data))
}

// Variance
func variance(data []float64) float64 {
    m := mean(data)
    sum := 0.0
    for _, v := range data {
        sum += (v - m) * (v - m)
    }
    return sum / float64(len(data))
}

// Standard Deviation
func stdDev(data []float64) float64 {
    return math.Sqrt(variance(data))
}

// Median
func median(data []float64) float64 {
    sorted := make([]float64, len(data))
    copy(sorted, data)
    sort.Float64s(sorted)
    
    n := len(sorted)
    if n%2 == 0 {
        return (sorted[n/2-1] + sorted[n/2]) / 2
    }
    return sorted[n/2]
}
```

### Percentiles

```go
func percentile(data []float64, p float64) float64 {
    sorted := make([]float64, len(data))
    copy(sorted, data)
    sort.Float64s(sorted)
    
    index := p / 100 * float64(len(sorted)-1)
    lower := int(index)
    upper := lower + 1
    
    if upper >= len(sorted) {
        return sorted[len(sorted)-1]
    }
    
    weight := index - float64(lower)
    return sorted[lower]*(1-weight) + sorted[upper]*weight
}
```

---

## Probability

### Combinations & Permutations

```go
// Factorial
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    result := 1
    for i := 2; i <= n; i++ {
        result *= i
    }
    return result
}

// Permutations: P(n,r) = n! / (n-r)!
func permutations(n, r int) int {
    return factorial(n) / factorial(n-r)
}

// Combinations: C(n,r) = n! / (r! * (n-r)!)
func combinations(n, r int) int {
    return factorial(n) / (factorial(r) * factorial(n-r))
}
```

### Probability

```go
// Simple probability
func probability(favorable, total int) float64 {
    return float64(favorable) / float64(total)
}

// Expected value
func expectedValue(values []float64, probs []float64) float64 {
    sum := 0.0
    for i := range values {
        sum += values[i] * probs[i]
    }
    return sum
}
```

---

## Linear Algebra Basics

### Vectors

```go
// Dot product
func dotProduct(a, b []float64) float64 {
    sum := 0.0
    for i := range a {
        sum += a[i] * b[i]
    }
    return sum
}

// Magnitude
func magnitude(v []float64) float64 {
    sum := 0.0
    for _, val := range v {
        sum += val * val
    }
    return math.Sqrt(sum)
}

// Cosine similarity
func cosineSimilarity(a, b []float64) float64 {
    return dotProduct(a, b) / (magnitude(a) * magnitude(b))
}
```

---

## Mathematical Reasoning in Code

### Invariants

```go
// Loop invariant: sum contains sum of items[0:i]
func sum(items []int) int {
    sum := 0
    for i := 0; i < len(items); i++ {
        // Invariant holds: sum == items[0] + ... + items[i-1]
        sum += items[i]
        // Invariant holds: sum == items[0] + ... + items[i]
    }
    return sum
}
```

### Proof by Induction

When designing recursive algorithms:
1. **Base case**: Prove for smallest input
2. **Inductive step**: Assume true for n, prove for n+1

```go
// Prove: sum of 1 to n = n*(n+1)/2
// Base: n=1: 1 = 1*2/2 ✓
// Inductive: Assume sum(n) = n*(n+1)/2
//            sum(n+1) = sum(n) + (n+1)
//                     = n*(n+1)/2 + (n+1)
//                     = (n+1)*(n+2)/2 ✓
func sumToN(n int) int {
    return n * (n + 1) / 2
}
```

---

## Related Skills

- `computer-science` - CS fundamentals
- `performance` - Complexity matters for performance
- `data-analyst` - Statistical analysis
