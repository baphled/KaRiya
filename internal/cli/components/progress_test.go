package components

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestProgressBar tests the ProgressBar component
func TestProgressBar(t *testing.T) {
	t.Run("creates progress bar with initial values", func(t *testing.T) {
		pb := NewProgressBar("Exporting", 100)
		assert.NotNil(t, pb)
		assert.Equal(t, "Exporting", pb.label)
		assert.Equal(t, 0, pb.current)
		assert.Equal(t, 100, pb.total)
	})

	t.Run("set current progress", func(t *testing.T) {
		pb := NewProgressBar("Exporting", 100)
		pb.SetCurrent(50)
		current, total := pb.GetProgress()
		assert.Equal(t, 50, current)
		assert.Equal(t, 100, total)
	})

	t.Run("increment progress", func(t *testing.T) {
		pb := NewProgressBar("Exporting", 100)
		pb.Increment()
		current, _ := pb.GetProgress()
		assert.Equal(t, 1, current)
	})

	t.Run("increment by amount", func(t *testing.T) {
		pb := NewProgressBar("Exporting", 100)
		pb.IncrementBy(25)
		current, _ := pb.GetProgress()
		assert.Equal(t, 25, current)
	})

	t.Run("clamps current to bounds", func(t *testing.T) {
		pb := NewProgressBar("Exporting", 100)
		pb.SetCurrent(-10)
		current, _ := pb.GetProgress()
		assert.Equal(t, 0, current)

		pb.SetCurrent(150)
		current, _ = pb.GetProgress()
		assert.Equal(t, 100, current)
	})

	t.Run("calculates percentage correctly", func(t *testing.T) {
		pb := NewProgressBar("Exporting", 100)
		assert.Equal(t, 0, pb.GetPercentage())

		pb.SetCurrent(25)
		assert.Equal(t, 25, pb.GetPercentage())

		pb.SetCurrent(50)
		assert.Equal(t, 50, pb.GetPercentage())

		pb.SetCurrent(100)
		assert.Equal(t, 100, pb.GetPercentage())
	})

	t.Run("handles zero total", func(t *testing.T) {
		pb := NewProgressBar("Exporting", 0)
		assert.Equal(t, 0, pb.GetPercentage())
	})

	t.Run("detects completion", func(t *testing.T) {
		pb := NewProgressBar("Exporting", 100)
		assert.False(t, pb.IsComplete())

		pb.SetCurrent(99)
		assert.False(t, pb.IsComplete())

		pb.SetCurrent(100)
		assert.True(t, pb.IsComplete())
	})

	t.Run("renders progress bar", func(t *testing.T) {
		pb := NewProgressBar("Exporting", 100)
		pb.SetCurrent(50)
		rendered := pb.Render()

		// Should contain the label
		assert.Contains(t, rendered, "Exporting")
		// Should contain the percentage
		assert.Contains(t, rendered, "50%")
		// Should contain progress info
		assert.Contains(t, rendered, "50/100")
	})

	t.Run("renders without label", func(t *testing.T) {
		pb := NewProgressBar("Exporting", 100)
		pb.SetShowLabel(false)
		rendered := pb.Render()

		// Should not contain the label
		assert.NotContains(t, rendered, "Exporting")
		// Should still contain percentage
		assert.Contains(t, rendered, "%")
	})

	t.Run("renders without percentage", func(t *testing.T) {
		pb := NewProgressBar("Exporting", 100)
		pb.SetShowPercent(false)
		rendered := pb.Render()

		// Should contain the label
		assert.Contains(t, rendered, "Exporting")
		// Should not contain percentage sign
		assert.NotContains(t, rendered, "%")
	})

	t.Run("sets custom width", func(t *testing.T) {
		pb := NewProgressBar("Exporting", 100)
		pb.SetWidth(20)
		// Should not panic and should render
		rendered := pb.Render()
		assert.NotEmpty(t, rendered)
	})

	t.Run("minimum width is enforced", func(t *testing.T) {
		pb := NewProgressBar("Exporting", 100)
		pb.SetWidth(5)
		// Width should be set to minimum of 10
		assert.Equal(t, 10, pb.width)
	})
}

// TestProgressIndicator tests the ProgressIndicator component
func TestProgressIndicator(t *testing.T) {
	t.Run("creates indicator with label", func(t *testing.T) {
		pi := NewProgressIndicator("Loading")
		assert.NotNil(t, pi)
		assert.Equal(t, "Loading", pi.label)
		assert.Equal(t, 0, pi.index)
		assert.Greater(t, len(pi.frames), 0)
	})

	t.Run("advances to next frame", func(t *testing.T) {
		pi := NewProgressIndicator("Loading")
		initialFrame := pi.GetFrame()
		pi.Next()
		nextFrame := pi.GetFrame()
		assert.NotEqual(t, initialFrame, nextFrame)
	})

	t.Run("wraps around frames", func(t *testing.T) {
		pi := NewProgressIndicator("Loading")
		firstFrame := pi.GetFrame()

		// Advance through all frames
		for i := 0; i < len(pi.frames); i++ {
			pi.Next()
		}

		// Should be back at the first frame
		assert.Equal(t, firstFrame, pi.GetFrame())
	})

	t.Run("sets label", func(t *testing.T) {
		pi := NewProgressIndicator("Loading")
		pi.SetLabel("Processing")
		assert.Equal(t, "Processing", pi.label)
	})

	t.Run("renders indicator", func(t *testing.T) {
		pi := NewProgressIndicator("Loading")
		rendered := pi.Render()

		// Should contain the label
		assert.Contains(t, rendered, "Loading")
		// Should contain a frame character
		assert.NotEmpty(t, rendered)
	})

	t.Run("renders without label", func(t *testing.T) {
		pi := NewProgressIndicator("Loading")
		pi.SetShowLabel(false)
		rendered := pi.Render()

		// Should not contain the label
		assert.NotContains(t, rendered, "Loading")
		// Should still render a frame
		assert.NotEmpty(t, rendered)
	})

	t.Run("renders different frames", func(t *testing.T) {
		pi := NewProgressIndicator("Loading")
		frame1 := pi.Render()

		pi.Next()
		frame2 := pi.Render()

		// Rendered output should be different
		assert.NotEqual(t, frame1, frame2)
	})
}

// TestSimpleProgressBar tests the simple progress bar function
func TestSimpleProgressBar(t *testing.T) {
	t.Run("renders simple progress bar", func(t *testing.T) {
		result := SimpleProgressBar(50, 100, 40)
		assert.NotEmpty(t, result)
		// Should contain percentage
		assert.Contains(t, result, "%")
		// Should contain progress info
		assert.Contains(t, result, "50/100")
	})

	t.Run("handles zero total", func(t *testing.T) {
		result := SimpleProgressBar(0, 0, 40)
		assert.NotEmpty(t, result)
		// Should contain 0% progress
		assert.Contains(t, result, "0%")
	})

	t.Run("enforces minimum width", func(t *testing.T) {
		result := SimpleProgressBar(50, 100, 5)
		assert.NotEmpty(t, result)
	})

	t.Run("renders different progress levels", func(t *testing.T) {
		bar0 := SimpleProgressBar(0, 100, 40)
		bar50 := SimpleProgressBar(50, 100, 40)
		bar100 := SimpleProgressBar(100, 100, 40)

		// All should be different
		assert.NotEqual(t, bar0, bar50)
		assert.NotEqual(t, bar50, bar100)
	})
}

// TestSimpleProgressIndicator tests the simple progress indicator function
func TestSimpleProgressIndicator(t *testing.T) {
	t.Run("renders indicator at different indices", func(t *testing.T) {
		ind0 := SimpleProgressIndicator(0)
		ind1 := SimpleProgressIndicator(1)
		ind9 := SimpleProgressIndicator(9)

		assert.NotEmpty(t, ind0)
		assert.NotEmpty(t, ind1)
		assert.NotEmpty(t, ind9)

		// Different indices should produce different output
		assert.NotEqual(t, ind0, ind1)
	})

	t.Run("wraps around indicator frames", func(t *testing.T) {
		ind0 := SimpleProgressIndicator(0)
		ind10 := SimpleProgressIndicator(10)

		// Index 10 should wrap back to index 0
		assert.Equal(t, ind0, ind10)
	})

	t.Run("handles large indices", func(t *testing.T) {
		result := SimpleProgressIndicator(1000)
		assert.NotEmpty(t, result)
	})
}

// BenchmarkProgressBar benchmarks progress bar rendering
func BenchmarkProgressBar(b *testing.B) {
	pb := NewProgressBar("Exporting", 100)
	pb.SetCurrent(50)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pb.Render()
	}
}

// BenchmarkProgressIndicator benchmarks progress indicator rendering
func BenchmarkProgressIndicator(b *testing.B) {
	pi := NewProgressIndicator("Loading")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pi.Render()
		pi.Next()
	}
}

// BenchmarkSimpleProgressBar benchmarks simple progress bar
func BenchmarkSimpleProgressBar(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SimpleProgressBar(50, 100, 40)
	}
}

// BenchmarkSimpleProgressIndicator benchmarks simple progress indicator
func BenchmarkSimpleProgressIndicator(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SimpleProgressIndicator(i)
	}
}
