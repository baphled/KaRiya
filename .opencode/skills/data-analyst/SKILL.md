---
name: data-analyst
description: Data exploration, statistical analysis, visualisation concepts, and deriving insights
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

# Data Analyst Skill

You are a data analyst expert skilled in data exploration, statistical analysis, visualization concepts, and deriving insights from data.

## Overview

Data analysis transforms raw data into actionable insights. Understand the data, ask the right questions, and communicate findings clearly.

---

## Data Analysis Process

```
1. Define the Question
   └─> What are we trying to learn?

2. Collect Data
   └─> What data do we need?

3. Clean Data
   └─> Handle missing values, outliers, errors

4. Explore Data
   └─> Descriptive statistics, visualizations

5. Analyze Data
   └─> Statistical tests, modeling

6. Communicate Results
   └─> Clear, actionable insights
```

---

## Data Exploration

### Descriptive Statistics

```go
import (
    "math"
    "sort"
)

type DataStats struct {
    Count    int
    Sum      float64
    Mean     float64
    Median   float64
    Mode     float64
    Min      float64
    Max      float64
    Range    float64
    Variance float64
    StdDev   float64
    Q1       float64
    Q3       float64
    IQR      float64
}

func Describe(data []float64) DataStats {
    n := len(data)
    sorted := make([]float64, n)
    copy(sorted, data)
    sort.Float64s(sorted)
    
    sum := 0.0
    for _, v := range data {
        sum += v
    }
    mean := sum / float64(n)
    
    variance := 0.0
    for _, v := range data {
        variance += (v - mean) * (v - mean)
    }
    variance /= float64(n)
    
    return DataStats{
        Count:    n,
        Sum:      sum,
        Mean:     mean,
        Median:   percentile(sorted, 50),
        Min:      sorted[0],
        Max:      sorted[n-1],
        Range:    sorted[n-1] - sorted[0],
        Variance: variance,
        StdDev:   math.Sqrt(variance),
        Q1:       percentile(sorted, 25),
        Q3:       percentile(sorted, 75),
        IQR:      percentile(sorted, 75) - percentile(sorted, 25),
    }
}
```

### Data Quality Checks

```go
type DataQuality struct {
    TotalRows     int
    NullCount     map[string]int
    UniqueCount   map[string]int
    DuplicateRows int
}

func CheckQuality(data []map[string]interface{}) DataQuality {
    quality := DataQuality{
        TotalRows:   len(data),
        NullCount:   make(map[string]int),
        UniqueCount: make(map[string]int),
    }
    
    // Check for nulls
    for _, row := range data {
        for key, value := range row {
            if value == nil {
                quality.NullCount[key]++
            }
        }
    }
    
    // Count unique values per column
    uniqueSets := make(map[string]map[interface{}]bool)
    for _, row := range data {
        for key, value := range row {
            if uniqueSets[key] == nil {
                uniqueSets[key] = make(map[interface{}]bool)
            }
            uniqueSets[key][value] = true
        }
    }
    for key, set := range uniqueSets {
        quality.UniqueCount[key] = len(set)
    }
    
    return quality
}
```

---

## Data Cleaning

### Handling Missing Values

```go
// Remove rows with missing values
func RemoveNulls(data []map[string]interface{}, column string) []map[string]interface{} {
    result := make([]map[string]interface{}, 0)
    for _, row := range data {
        if row[column] != nil {
            result = append(result, row)
        }
    }
    return result
}

// Fill with mean
func FillWithMean(data []map[string]float64, column string) {
    sum := 0.0
    count := 0
    for _, row := range data {
        if v, ok := row[column]; ok && !math.IsNaN(v) {
            sum += v
            count++
        }
    }
    mean := sum / float64(count)
    
    for i := range data {
        if math.IsNaN(data[i][column]) {
            data[i][column] = mean
        }
    }
}

// Fill with median (better for skewed data)
func FillWithMedian(data []map[string]float64, column string) {
    values := make([]float64, 0)
    for _, row := range data {
        if v, ok := row[column]; ok && !math.IsNaN(v) {
            values = append(values, v)
        }
    }
    sort.Float64s(values)
    median := values[len(values)/2]
    
    for i := range data {
        if math.IsNaN(data[i][column]) {
            data[i][column] = median
        }
    }
}
```

### Outlier Detection

```go
// IQR method
func DetectOutliers(data []float64) []int {
    sorted := make([]float64, len(data))
    copy(sorted, data)
    sort.Float64s(sorted)
    
    q1 := percentile(sorted, 25)
    q3 := percentile(sorted, 75)
    iqr := q3 - q1
    
    lower := q1 - 1.5*iqr
    upper := q3 + 1.5*iqr
    
    outliers := make([]int, 0)
    for i, v := range data {
        if v < lower || v > upper {
            outliers = append(outliers, i)
        }
    }
    return outliers
}

// Z-score method
func DetectOutliersZScore(data []float64, threshold float64) []int {
    mean := mean(data)
    stdDev := stdDev(data)
    
    outliers := make([]int, 0)
    for i, v := range data {
        zScore := (v - mean) / stdDev
        if math.Abs(zScore) > threshold {
            outliers = append(outliers, i)
        }
    }
    return outliers
}
```

---

## Aggregation & Grouping

```go
// Group by category
func GroupBy(data []map[string]interface{}, key string) map[interface{}][]map[string]interface{} {
    groups := make(map[interface{}][]map[string]interface{})
    for _, row := range data {
        k := row[key]
        groups[k] = append(groups[k], row)
    }
    return groups
}

// Aggregate functions
type AggFunc func([]float64) float64

func Sum(values []float64) float64 {
    sum := 0.0
    for _, v := range values {
        sum += v
    }
    return sum
}

func Count(values []float64) float64 {
    return float64(len(values))
}

func Avg(values []float64) float64 {
    return Sum(values) / Count(values)
}

// Group and aggregate
func GroupAggregate(data []map[string]interface{}, groupKey, valueKey string, agg AggFunc) map[interface{}]float64 {
    groups := GroupBy(data, groupKey)
    result := make(map[interface{}]float64)
    
    for key, rows := range groups {
        values := make([]float64, 0)
        for _, row := range rows {
            if v, ok := row[valueKey].(float64); ok {
                values = append(values, v)
            }
        }
        result[key] = agg(values)
    }
    return result
}
```

---

## Correlation Analysis

```go
// Pearson correlation coefficient
func Correlation(x, y []float64) float64 {
    if len(x) != len(y) {
        return 0
    }
    
    n := float64(len(x))
    sumX, sumY, sumXY := 0.0, 0.0, 0.0
    sumX2, sumY2 := 0.0, 0.0
    
    for i := range x {
        sumX += x[i]
        sumY += y[i]
        sumXY += x[i] * y[i]
        sumX2 += x[i] * x[i]
        sumY2 += y[i] * y[i]
    }
    
    numerator := n*sumXY - sumX*sumY
    denominator := math.Sqrt((n*sumX2 - sumX*sumX) * (n*sumY2 - sumY*sumY))
    
    if denominator == 0 {
        return 0
    }
    return numerator / denominator
}

// Interpretation
// |r| < 0.3: Weak
// 0.3 <= |r| < 0.7: Moderate
// |r| >= 0.7: Strong
```

---

## Time Series Analysis

```go
// Moving average
func MovingAverage(data []float64, window int) []float64 {
    result := make([]float64, len(data)-window+1)
    
    for i := 0; i <= len(data)-window; i++ {
        sum := 0.0
        for j := 0; j < window; j++ {
            sum += data[i+j]
        }
        result[i] = sum / float64(window)
    }
    return result
}

// Growth rate
func GrowthRate(current, previous float64) float64 {
    if previous == 0 {
        return 0
    }
    return (current - previous) / previous * 100
}

// Year-over-year comparison
func YoYGrowth(current, previousYear float64) float64 {
    return GrowthRate(current, previousYear)
}
```

---

## Data Validation

```go
// Validate data against rules
type ValidationRule struct {
    Column    string
    Validator func(interface{}) bool
    Message   string
}

func Validate(data []map[string]interface{}, rules []ValidationRule) []string {
    errors := make([]string, 0)
    
    for i, row := range data {
        for _, rule := range rules {
            value := row[rule.Column]
            if !rule.Validator(value) {
                errors = append(errors, 
                    fmt.Sprintf("Row %d, %s: %s", i, rule.Column, rule.Message))
            }
        }
    }
    return errors
}

// Common validators
func NotNull(v interface{}) bool {
    return v != nil
}

func InRange(min, max float64) func(interface{}) bool {
    return func(v interface{}) bool {
        f, ok := v.(float64)
        return ok && f >= min && f <= max
    }
}

func MatchPattern(pattern string) func(interface{}) bool {
    re := regexp.MustCompile(pattern)
    return func(v interface{}) bool {
        s, ok := v.(string)
        return ok && re.MatchString(s)
    }
}
```

---

## Report Generation

```go
type Report struct {
    Title       string
    GeneratedAt time.Time
    Summary     string
    Metrics     map[string]float64
    Findings    []string
}

func GenerateReport(data []map[string]interface{}, title string) Report {
    stats := Describe(extractColumn(data, "value"))
    
    findings := make([]string, 0)
    
    // Identify key findings
    if stats.StdDev > stats.Mean*0.5 {
        findings = append(findings, "High variability in data")
    }
    
    outliers := DetectOutliers(extractColumn(data, "value"))
    if len(outliers) > 0 {
        findings = append(findings, 
            fmt.Sprintf("%d outliers detected", len(outliers)))
    }
    
    return Report{
        Title:       title,
        GeneratedAt: time.Now(),
        Metrics: map[string]float64{
            "count":  float64(stats.Count),
            "mean":   stats.Mean,
            "median": stats.Median,
            "stddev": stats.StdDev,
        },
        Findings: findings,
    }
}
```

---

## Best Practices

### Data Analysis Checklist

- [ ] Understand the business question
- [ ] Verify data source and quality
- [ ] Check for missing values
- [ ] Check for outliers
- [ ] Understand distributions
- [ ] Look for correlations
- [ ] Validate assumptions
- [ ] Document methodology
- [ ] Communicate clearly

### Avoid Common Pitfalls

1. **Correlation ≠ Causation**
2. **Survivorship bias** - Missing data from failures
3. **Simpson's paradox** - Aggregation can reverse trends
4. **Overfitting** - Model fits noise, not signal
5. **Cherry-picking** - Selecting data to support conclusion

---

## Related Skills

- `math-expert` - Statistical foundations
- `computer-science` - Data structures
- `performance` - Efficient data processing
