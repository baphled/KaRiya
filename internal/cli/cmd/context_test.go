package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestNewCLIContext(t *testing.T) {
	tests := []struct {
		name     string
		dbPath   string
		inMemory bool
		wantPath string
		wantMem  bool
	}{
		{
			name:     "in-memory mode",
			dbPath:   "",
			inMemory: true,
			wantPath: "",
			wantMem:  true,
		},
		{
			name:     "sqlite with custom path",
			dbPath:   "/tmp/test.db",
			inMemory: false,
			wantPath: "/tmp/test.db",
			wantMem:  false,
		},
		{
			name:     "sqlite with empty path",
			dbPath:   "",
			inMemory: false,
			wantPath: "",
			wantMem:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewCLIContext(tt.dbPath, tt.inMemory)

			if ctx == nil {
				t.Fatal("NewCLIContext returned nil")
			}
			if ctx.inMemory != tt.wantMem {
				t.Errorf("inMemory = %v, want %v", ctx.inMemory, tt.wantMem)
			}
			if ctx.dbPath != tt.wantPath {
				t.Errorf("dbPath = %q, want %q", ctx.dbPath, tt.wantPath)
			}
			if ctx.svc != nil {
				t.Error("service should not be initialized in constructor")
			}
		})
	}
}

func TestCLIContextInitServiceInMemory(t *testing.T) {
	ctx := NewCLIContext("", true)
	errBuf := new(bytes.Buffer)

	err := ctx.InitService(errBuf)

	if err != nil {
		t.Fatalf("InitService failed: %v", err)
	}
	if ctx.svc == nil {
		t.Fatal("service not initialized")
	}
	if errBuf.Len() > 0 {
		t.Errorf("unexpected error output: %s", errBuf.String())
	}

	if ctx.svc.GetEventRepository() == nil {
		t.Error("event repository not set")
	}
	if ctx.svc.GetFactRepository() == nil {
		t.Error("fact repository not set")
	}
	if ctx.svc.GetBurstRepository() == nil {
		t.Error("burst repository not set")
	}
	if ctx.svc.GetSkillRepository() == nil {
		t.Error("skill repository not set")
	}
}

func TestCLIContextInitServiceSQLiteCustomPath(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	ctx := NewCLIContext(dbPath, false)
	errBuf := new(bytes.Buffer)

	err := ctx.InitService(errBuf)

	if err != nil {
		t.Fatalf("InitService failed: %v", err)
	}
	if ctx.svc == nil {
		t.Fatal("service not initialized")
	}
	if errBuf.Len() > 0 {
		t.Errorf("unexpected error output: %s", errBuf.String())
	}

	if _, err := os.Stat(dbPath); err != nil {
		t.Errorf("database file not created: %v", err)
	}

	if ctx.svc.GetEventRepository() == nil {
		t.Error("event repository not set")
	}
	if ctx.svc.GetFactRepository() == nil {
		t.Error("fact repository not set")
	}
	if ctx.svc.GetBurstRepository() == nil {
		t.Error("burst repository not set")
	}
	if ctx.svc.GetSkillRepository() == nil {
		t.Error("skill repository not set")
	}
}

func TestCLIContextInitServiceSQLiteDefaultPath(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	ctx := NewCLIContext("", false)
	errBuf := new(bytes.Buffer)

	err := ctx.InitService(errBuf)

	if err != nil {
		t.Fatalf("InitService failed: %v", err)
	}
	if ctx.svc == nil {
		t.Fatal("service not initialized")
	}
	if errBuf.Len() > 0 {
		t.Errorf("unexpected error output: %s", errBuf.String())
	}

	expectedPath := filepath.Join(homeDir, ".kariya", "events.db")
	if _, err := os.Stat(expectedPath); err != nil {
		t.Errorf("database file not created at default path: %v", err)
	}

	kariyaDir := filepath.Join(homeDir, ".kariya")
	info, err := os.Stat(kariyaDir)
	if err != nil {
		t.Errorf("kariya directory not created: %v", err)
	}
	if !info.IsDir() {
		t.Error("kariya path is not a directory")
	}
	if info.Mode().Perm() != 0o750 {
		t.Errorf("kariya directory permissions = %o, want 0o750", info.Mode().Perm())
	}

	if ctx.svc.GetEventRepository() == nil {
		t.Error("event repository not set")
	}
	if ctx.svc.GetFactRepository() == nil {
		t.Error("fact repository not set")
	}
	if ctx.svc.GetBurstRepository() == nil {
		t.Error("burst repository not set")
	}
	if ctx.svc.GetSkillRepository() == nil {
		t.Error("skill repository not set")
	}
}

func TestCLIContextInitServiceInvalidDBPath(t *testing.T) {
	ctx := NewCLIContext("/invalid/path/that/does/not/exist/test.db", false)
	errBuf := new(bytes.Buffer)

	err := ctx.InitService(errBuf)

	if err == nil {
		t.Fatal("expected error for invalid DB path")
	}
	if ctx.svc != nil {
		t.Error("service should not be initialized on error")
	}
	if errBuf.Len() == 0 {
		t.Error("expected error output to errOut")
	}
}

func TestCLIContextService(t *testing.T) {
	ctx := NewCLIContext("", true)
	errBuf := new(bytes.Buffer)

	if ctx.Service() != nil {
		t.Error("Service should return nil before initialization")
	}

	err := ctx.InitService(errBuf)
	if err != nil {
		t.Fatalf("InitService failed: %v", err)
	}

	svc := ctx.Service()
	if svc == nil {
		t.Fatal("Service returned nil after initialization")
	}
	if svc != ctx.svc {
		t.Error("Service should return the initialized service")
	}
}

func TestCLIContextMultipleInitCalls(t *testing.T) {
	ctx := NewCLIContext("", true)
	errBuf := new(bytes.Buffer)

	err1 := ctx.InitService(errBuf)
	if err1 != nil {
		t.Fatalf("first InitService failed: %v", err1)
	}
	svc1 := ctx.svc

	err2 := ctx.InitService(errBuf)
	if err2 != nil {
		t.Fatalf("second InitService failed: %v", err2)
	}
	svc2 := ctx.svc

	if svc1 == svc2 {
		t.Error("multiple InitService calls should create a new service")
	}
}

func TestCLIContextServiceType(t *testing.T) {
	ctx := NewCLIContext("", true)
	errBuf := new(bytes.Buffer)

	err := ctx.InitService(errBuf)
	if err != nil {
		t.Fatalf("InitService failed: %v", err)
	}

	if ctx.svc == nil {
		t.Fatal("service is nil")
	}
}

func TestCLIContextInitServiceSQLiteCreateDirError(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	readOnlyDir := filepath.Join(homeDir, "readonly")
	if err := os.Mkdir(readOnlyDir, 0o555); err != nil {
		t.Fatalf("failed to create readonly dir: %v", err)
	}
	defer os.Chmod(readOnlyDir, 0o755)

	ctx := NewCLIContext(filepath.Join(readOnlyDir, ".kariya", "events.db"), false)
	errBuf := new(bytes.Buffer)

	err := ctx.InitService(errBuf)

	if err == nil {
		t.Fatal("expected error for directory creation failure")
	}
	if ctx.svc != nil {
		t.Error("service should not be initialized on error")
	}
	if errBuf.Len() == 0 {
		t.Error("expected error output to errOut")
	}
}

func TestCLIContextInitServiceSQLiteMigrationError(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	ctx := NewCLIContext(dbPath, false)
	errBuf := new(bytes.Buffer)

	err := ctx.InitService(errBuf)
	if err != nil {
		t.Fatalf("first InitService failed: %v", err)
	}

	ctx2 := NewCLIContext(dbPath, false)
	errBuf2 := new(bytes.Buffer)

	err2 := ctx2.InitService(errBuf2)
	if err2 != nil {
		t.Fatalf("second InitService failed: %v", err2)
	}
}

func TestCLIContextInitServiceSQLiteDefaultPathWithBadHome(t *testing.T) {
	t.Setenv("HOME", "/nonexistent/path/that/cannot/be/created")

	ctx := NewCLIContext("", false)
	errBuf := new(bytes.Buffer)

	err := ctx.InitService(errBuf)

	if err == nil {
		t.Fatal("expected error for bad home directory")
	}
	if ctx.svc != nil {
		t.Error("service should not be initialized on error")
	}
	if errBuf.Len() == 0 {
		t.Error("expected error output to errOut")
	}
}

func TestCLIContextInitServiceSQLiteRepositoriesError(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	ctx := NewCLIContext(dbPath, false)
	errBuf := new(bytes.Buffer)

	err := ctx.InitService(errBuf)
	if err != nil {
		t.Fatalf("InitService failed: %v", err)
	}

	if ctx.svc == nil {
		t.Fatal("service should be initialized")
	}
	if ctx.svc.GetEventRepository() == nil {
		t.Error("event repository not set")
	}
}

func TestCLIContextInitServiceSQLiteOpenDBError(t *testing.T) {
	ctx := NewCLIContext("/dev/null/invalid/path/test.db", false)
	errBuf := new(bytes.Buffer)

	err := ctx.InitService(errBuf)

	if err == nil {
		t.Fatal("expected error for invalid DB path")
	}
	if ctx.svc != nil {
		t.Error("service should not be initialized on error")
	}
	if errBuf.Len() == 0 {
		t.Error("expected error output to errOut")
	}
}

func TestCLIContextInitServiceSQLiteFullSuccess(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	ctx := NewCLIContext(dbPath, false)
	errBuf := new(bytes.Buffer)

	err := ctx.InitService(errBuf)

	if err != nil {
		t.Fatalf("InitService failed: %v", err)
	}
	if ctx.svc == nil {
		t.Fatal("service not initialized")
	}
	if errBuf.Len() > 0 {
		t.Errorf("unexpected error output: %s", errBuf.String())
	}

	eventRepo := ctx.svc.GetEventRepository()
	if eventRepo == nil {
		t.Fatal("event repository is nil")
	}

	factRepo := ctx.svc.GetFactRepository()
	if factRepo == nil {
		t.Fatal("fact repository is nil")
	}

	burstRepo := ctx.svc.GetBurstRepository()
	if burstRepo == nil {
		t.Fatal("burst repository is nil")
	}

	skillRepo := ctx.svc.GetSkillRepository()
	if skillRepo == nil {
		t.Fatal("skill repository is nil")
	}

	if _, err := os.Stat(dbPath); err != nil {
		t.Errorf("database file not created: %v", err)
	}
}

func TestCLIContextInitServiceInMemoryFullSuccess(t *testing.T) {
	ctx := NewCLIContext("", true)
	errBuf := new(bytes.Buffer)

	err := ctx.InitService(errBuf)

	if err != nil {
		t.Fatalf("InitService failed: %v", err)
	}
	if ctx.svc == nil {
		t.Fatal("service not initialized")
	}
	if errBuf.Len() > 0 {
		t.Errorf("unexpected error output: %s", errBuf.String())
	}

	eventRepo := ctx.svc.GetEventRepository()
	if eventRepo == nil {
		t.Fatal("event repository is nil")
	}

	factRepo := ctx.svc.GetFactRepository()
	if factRepo == nil {
		t.Fatal("fact repository is nil")
	}

	burstRepo := ctx.svc.GetBurstRepository()
	if burstRepo == nil {
		t.Fatal("burst repository is nil")
	}

	skillRepo := ctx.svc.GetSkillRepository()
	if skillRepo == nil {
		t.Fatal("skill repository is nil")
	}
}

func TestCLIContextInitServiceSQLiteDefaultPathSuccess(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	ctx := NewCLIContext("", false)
	errBuf := new(bytes.Buffer)

	err := ctx.InitService(errBuf)

	if err != nil {
		t.Fatalf("InitService failed: %v", err)
	}
	if ctx.svc == nil {
		t.Fatal("service not initialized")
	}
	if errBuf.Len() > 0 {
		t.Errorf("unexpected error output: %s", errBuf.String())
	}

	expectedPath := filepath.Join(homeDir, ".kariya", "events.db")
	if _, err := os.Stat(expectedPath); err != nil {
		t.Errorf("database file not created at default path: %v", err)
	}

	eventRepo := ctx.svc.GetEventRepository()
	if eventRepo == nil {
		t.Fatal("event repository is nil")
	}

	factRepo := ctx.svc.GetFactRepository()
	if factRepo == nil {
		t.Fatal("fact repository is nil")
	}

	burstRepo := ctx.svc.GetBurstRepository()
	if burstRepo == nil {
		t.Fatal("burst repository is nil")
	}

	skillRepo := ctx.svc.GetSkillRepository()
	if skillRepo == nil {
		t.Fatal("skill repository is nil")
	}
}

func TestCLIContextInitServiceSQLiteReadOnlyDBPath(t *testing.T) {
	tmpDir := t.TempDir()
	readOnlyDir := filepath.Join(tmpDir, "readonly")
	if err := os.Mkdir(readOnlyDir, 0o555); err != nil {
		t.Fatalf("failed to create readonly dir: %v", err)
	}
	defer os.Chmod(readOnlyDir, 0o755)

	dbPath := filepath.Join(readOnlyDir, "test.db")

	ctx := NewCLIContext(dbPath, false)
	errBuf := new(bytes.Buffer)

	err := ctx.InitService(errBuf)

	if err == nil {
		t.Fatal("expected error for readonly DB path")
	}
	if ctx.svc != nil {
		t.Error("service should not be initialized on error")
	}
	if errBuf.Len() == 0 {
		t.Error("expected error output to errOut")
	}
}

func TestCLIContextInitServiceSQLiteReadOnlyHomeDir(t *testing.T) {
	tmpDir := t.TempDir()
	readOnlyHome := filepath.Join(tmpDir, "readonly_home")
	if err := os.Mkdir(readOnlyHome, 0o555); err != nil {
		t.Fatalf("failed to create readonly home dir: %v", err)
	}
	defer os.Chmod(readOnlyHome, 0o755)

	t.Setenv("HOME", readOnlyHome)

	ctx := NewCLIContext("", false)
	errBuf := new(bytes.Buffer)

	err := ctx.InitService(errBuf)

	if err == nil {
		t.Fatal("expected error for readonly home directory")
	}
	if ctx.svc != nil {
		t.Error("service should not be initialized on error")
	}
	if errBuf.Len() == 0 {
		t.Error("expected error output to errOut")
	}
}

func TestCLIContextInitServiceSQLiteEmptyHomeDir(t *testing.T) {
	t.Setenv("HOME", "")

	ctx := NewCLIContext("", false)
	errBuf := new(bytes.Buffer)

	err := ctx.InitService(errBuf)

	if err == nil {
		t.Fatal("expected error for empty HOME")
	}
	if ctx.svc != nil {
		t.Error("service should not be initialized on error")
	}
	if errBuf.Len() == 0 {
		t.Error("expected error output to errOut")
	}
}

func TestCLIContextInitServiceSQLiteCorruptedDB(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "corrupted.db")

	if err := os.WriteFile(dbPath, []byte("not a valid sqlite database"), 0o600); err != nil {
		t.Fatalf("failed to create corrupted DB: %v", err)
	}

	ctx := NewCLIContext(dbPath, false)
	errBuf := new(bytes.Buffer)

	err := ctx.InitService(errBuf)

	if err == nil {
		t.Fatal("expected error for corrupted DB")
	}
	if ctx.svc != nil {
		t.Error("service should not be initialized on error")
	}
	if errBuf.Len() == 0 {
		t.Error("expected error output to errOut")
	}
}

func TestCLIContextInitServiceSQLiteDirectoryAsDB(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "db_dir")

	if err := os.Mkdir(dbPath, 0o755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	ctx := NewCLIContext(dbPath, false)
	errBuf := new(bytes.Buffer)

	err := ctx.InitService(errBuf)

	if err == nil {
		t.Fatal("expected error for directory as DB path")
	}
	if ctx.svc != nil {
		t.Error("service should not be initialized on error")
	}
	if errBuf.Len() == 0 {
		t.Error("expected error output to errOut")
	}
}
