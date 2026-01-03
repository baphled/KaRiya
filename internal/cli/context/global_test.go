package context

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewGlobalContext tests GlobalContext creation
func TestNewGlobalContext(t *testing.T) {
	t.Run("creates context with nil preferences and config", func(t *testing.T) {
		gc := NewGlobalContext(nil, nil)
		assert.NotNil(t, gc)
		assert.Equal(t, 0, len(gc.PreferenceKeys()))
		assert.Equal(t, 0, len(gc.ConfigKeys()))
		assert.Equal(t, 0, len(gc.TransientStateKeys()))
	})

	t.Run("creates context with initial preferences and config", func(t *testing.T) {
		prefs := map[string]interface{}{"theme": "dark"}
		cfg := map[string]interface{}{"version": "1.0"}

		gc := NewGlobalContext(prefs, cfg)
		assert.NotNil(t, gc)
		assert.Equal(t, 1, len(gc.PreferenceKeys()))
		assert.Equal(t, 1, len(gc.ConfigKeys()))

		val, exists := gc.GetPreference("theme")
		assert.True(t, exists)
		assert.Equal(t, "dark", val)

		val, exists = gc.GetConfig("version")
		assert.True(t, exists)
		assert.Equal(t, "1.0", val)
	})
}

// TestPreferences tests preference getter/setter
func TestPreferences(t *testing.T) {
	gc := NewGlobalContext(nil, nil)

	t.Run("set and get preference", func(t *testing.T) {
		gc.SetPreference("theme", "dark")

		val, exists := gc.GetPreference("theme")
		assert.True(t, exists)
		assert.Equal(t, "dark", val)
	})

	t.Run("get non-existent preference returns false", func(t *testing.T) {
		val, exists := gc.GetPreference("nonexistent")
		assert.False(t, exists)
		assert.Nil(t, val)
	})

	t.Run("get all preferences returns copy", func(t *testing.T) {
		gc.SetPreference("theme", "dark")
		gc.SetPreference("language", "en")

		prefs := gc.GetAllPreferences()
		assert.Equal(t, 2, len(prefs))
		assert.Equal(t, "dark", prefs["theme"])
		assert.Equal(t, "en", prefs["language"])

		// Verify it's a copy
		prefs["theme"] = "light"
		val, _ := gc.GetPreference("theme")
		assert.Equal(t, "dark", val)
	})

	t.Run("preference keys returns all keys", func(t *testing.T) {
		gc.SetPreference("theme", "dark")
		gc.SetPreference("language", "en")

		keys := gc.PreferenceKeys()
		assert.Equal(t, 2, len(keys))
		assert.Contains(t, keys, "theme")
		assert.Contains(t, keys, "language")
	})
}

// TestTransientState tests transient state getter/setter
func TestTransientState(t *testing.T) {
	gc := NewGlobalContext(nil, nil)

	t.Run("set and get transient state", func(t *testing.T) {
		gc.SetTransientState("selected_event", "evt123")

		val, exists := gc.GetTransientState("selected_event")
		assert.True(t, exists)
		assert.Equal(t, "evt123", val)
	})

	t.Run("get non-existent transient state returns false", func(t *testing.T) {
		val, exists := gc.GetTransientState("nonexistent")
		assert.False(t, exists)
		assert.Nil(t, val)
	})

	t.Run("clear transient state", func(t *testing.T) {
		gc.SetTransientState("selected_event", "evt123")
		gc.ClearTransientState("selected_event")

		val, exists := gc.GetTransientState("selected_event")
		assert.False(t, exists)
		assert.Nil(t, val)
	})

	t.Run("get all transient state returns copy", func(t *testing.T) {
		gc.SetTransientState("selected_event", "evt123")
		gc.SetTransientState("filter_tags", []string{"go", "rust"})

		state := gc.GetAllTransientState()
		assert.Equal(t, 2, len(state))
		assert.Equal(t, "evt123", state["selected_event"])

		// Verify it's a copy
		state["selected_event"] = "evt456"
		val, _ := gc.GetTransientState("selected_event")
		assert.Equal(t, "evt123", val)
	})

	t.Run("clear all transient state", func(t *testing.T) {
		gc.SetTransientState("key1", "val1")
		gc.SetTransientState("key2", "val2")
		gc.ClearAllTransientState()

		assert.Equal(t, 0, len(gc.TransientStateKeys()))
	})

	t.Run("transient state keys returns all keys", func(t *testing.T) {
		gc.SetTransientState("selected_event", "evt123")
		gc.SetTransientState("filter_tags", []string{"go"})

		keys := gc.TransientStateKeys()
		assert.Equal(t, 2, len(keys))
		assert.Contains(t, keys, "selected_event")
		assert.Contains(t, keys, "filter_tags")
	})
}

// TestConfig tests configuration getter/setter
func TestConfig(t *testing.T) {
	gc := NewGlobalContext(nil, nil)

	t.Run("set and get config", func(t *testing.T) {
		gc.SetConfig("version", "1.0.0")

		val, exists := gc.GetConfig("version")
		assert.True(t, exists)
		assert.Equal(t, "1.0.0", val)
	})

	t.Run("get non-existent config returns false", func(t *testing.T) {
		val, exists := gc.GetConfig("nonexistent")
		assert.False(t, exists)
		assert.Nil(t, val)
	})

	t.Run("get all config returns copy", func(t *testing.T) {
		gc.SetConfig("version", "1.0.0")
		gc.SetConfig("app_name", "KaRiya")

		cfg := gc.GetAllConfig()
		assert.Equal(t, 2, len(cfg))
		assert.Equal(t, "1.0.0", cfg["version"])
		assert.Equal(t, "KaRiya", cfg["app_name"])

		// Verify it's a copy
		cfg["version"] = "2.0.0"
		val, _ := gc.GetConfig("version")
		assert.Equal(t, "1.0.0", val)
	})

	t.Run("config keys returns all keys", func(t *testing.T) {
		gc.SetConfig("version", "1.0.0")
		gc.SetConfig("app_name", "KaRiya")

		keys := gc.ConfigKeys()
		assert.Equal(t, 2, len(keys))
		assert.Contains(t, keys, "version")
		assert.Contains(t, keys, "app_name")
	})
}

// TestThreadSafety tests concurrent access to GlobalContext
func TestThreadSafety(t *testing.T) {
	gc := NewGlobalContext(nil, nil)

	t.Run("concurrent preference access is safe", func(t *testing.T) {
		var wg sync.WaitGroup
		done := make(chan bool)

		// Writer goroutines
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for j := 0; j < 100; j++ {
					gc.SetPreference("key", id*100+j)
				}
			}(i)
		}

		// Reader goroutines
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 100; j++ {
					gc.GetPreference("key")
					gc.GetAllPreferences()
				}
			}()
		}

		go func() {
			wg.Wait()
			done <- true
		}()

		<-done
		// If we get here, no deadlock occurred
		assert.True(t, true)
	})

	t.Run("concurrent transient state access is safe", func(t *testing.T) {
		gc := NewGlobalContext(nil, nil)
		var wg sync.WaitGroup
		done := make(chan bool)

		// Writer goroutines
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for j := 0; j < 100; j++ {
					gc.SetTransientState("key", id*100+j)
				}
			}(i)
		}

		// Reader goroutines
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 100; j++ {
					gc.GetTransientState("key")
					gc.GetAllTransientState()
				}
			}()
		}

		go func() {
			wg.Wait()
			done <- true
		}()

		<-done
		// If we get here, no deadlock occurred
		assert.True(t, true)
	})

	t.Run("concurrent config access is safe", func(t *testing.T) {
		gc := NewGlobalContext(nil, nil)
		var wg sync.WaitGroup
		done := make(chan bool)

		// Writer goroutines
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for j := 0; j < 100; j++ {
					gc.SetConfig("key", id*100+j)
				}
			}(i)
		}

		// Reader goroutines
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 100; j++ {
					gc.GetConfig("key")
					gc.GetAllConfig()
				}
			}()
		}

		go func() {
			wg.Wait()
			done <- true
		}()

		<-done
		// If we get here, no deadlock occurred
		assert.True(t, true)
	})
}

// TestDataTypes tests that GlobalContext handles various data types
func TestDataTypes(t *testing.T) {
	gc := NewGlobalContext(nil, nil)

	t.Run("stores and retrieves strings", func(t *testing.T) {
		gc.SetPreference("name", "Alice")
		val, _ := gc.GetPreference("name")
		assert.Equal(t, "Alice", val.(string))
	})

	t.Run("stores and retrieves integers", func(t *testing.T) {
		gc.SetPreference("count", 42)
		val, _ := gc.GetPreference("count")
		assert.Equal(t, 42, val.(int))
	})

	t.Run("stores and retrieves booleans", func(t *testing.T) {
		gc.SetPreference("enabled", true)
		val, _ := gc.GetPreference("enabled")
		assert.Equal(t, true, val.(bool))
	})

	t.Run("stores and retrieves slices", func(t *testing.T) {
		tags := []string{"go", "rust", "python"}
		gc.SetTransientState("tags", tags)
		val, _ := gc.GetTransientState("tags")
		assert.Equal(t, tags, val.([]string))
	})

	t.Run("stores and retrieves maps", func(t *testing.T) {
		data := map[string]int{"a": 1, "b": 2}
		gc.SetTransientState("data", data)
		val, _ := gc.GetTransientState("data")
		assert.Equal(t, data, val.(map[string]int))
	})

	t.Run("stores and retrieves nil values", func(t *testing.T) {
		gc.SetTransientState("nullable", nil)
		val, exists := gc.GetTransientState("nullable")
		assert.True(t, exists)
		assert.Nil(t, val)
	})
}

// TestOverwriting tests that values can be overwritten
func TestOverwriting(t *testing.T) {
	gc := NewGlobalContext(nil, nil)

	t.Run("overwrite preference", func(t *testing.T) {
		gc.SetPreference("theme", "dark")
		gc.SetPreference("theme", "light")

		val, _ := gc.GetPreference("theme")
		assert.Equal(t, "light", val)
	})

	t.Run("overwrite transient state", func(t *testing.T) {
		gc.SetTransientState("selected", "item1")
		gc.SetTransientState("selected", "item2")

		val, _ := gc.GetTransientState("selected")
		assert.Equal(t, "item2", val)
	})

	t.Run("overwrite config", func(t *testing.T) {
		gc.SetConfig("version", "1.0.0")
		gc.SetConfig("version", "2.0.0")

		val, _ := gc.GetConfig("version")
		assert.Equal(t, "2.0.0", val)
	})
}

