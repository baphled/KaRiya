package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CLIContext", func() {
	Describe("NewCLIContext", func() {
		Context("in-memory mode", func() {
			It("should create context with in-memory flag", func() {
				ctx := NewCLIContext("", true)
				Expect(ctx).NotTo(BeNil())
				Expect(ctx.inMemory).To(BeTrue())
				Expect(ctx.dbPath).To(Equal(""))
				Expect(ctx.svc).To(BeNil())
			})
		})

		Context("SQLite with custom path", func() {
			It("should create context with custom database path", func() {
				ctx := NewCLIContext("/tmp/test.db", false)
				Expect(ctx).NotTo(BeNil())
				Expect(ctx.inMemory).To(BeFalse())
				Expect(ctx.dbPath).To(Equal("/tmp/test.db"))
				Expect(ctx.svc).To(BeNil())
			})
		})

		Context("SQLite with empty path", func() {
			It("should create context with empty path", func() {
				ctx := NewCLIContext("", false)
				Expect(ctx).NotTo(BeNil())
				Expect(ctx.inMemory).To(BeFalse())
				Expect(ctx.dbPath).To(Equal(""))
				Expect(ctx.svc).To(BeNil())
			})
		})
	})

	Describe("InitService", func() {
		Context("with in-memory database", func() {
			It("should initialize service successfully", func() {
				ctx := NewCLIContext("", true)
				errBuf := new(bytes.Buffer)

				err := ctx.InitService(errBuf)
				Expect(err).NotTo(HaveOccurred())
				Expect(ctx.svc).NotTo(BeNil())
				Expect(errBuf.Len()).To(Equal(0))
			})

			It("should initialize all repositories", func() {
				ctx := NewCLIContext("", true)
				errBuf := new(bytes.Buffer)

				err := ctx.InitService(errBuf)
				Expect(err).NotTo(HaveOccurred())

				Expect(ctx.svc.GetEventRepository()).NotTo(BeNil())
				Expect(ctx.svc.GetFactRepository()).NotTo(BeNil())
				Expect(ctx.svc.GetBurstRepository()).NotTo(BeNil())
				Expect(ctx.svc.GetSkillRepository()).NotTo(BeNil())
			})
		})

		Context("with SQLite and custom path", func() {
			It("should initialize service successfully", func() {
				tmpDir := GinkgoT().TempDir()
				dbPath := filepath.Join(tmpDir, "test.db")

				ctx := NewCLIContext(dbPath, false)
				errBuf := new(bytes.Buffer)

				err := ctx.InitService(errBuf)
				Expect(err).NotTo(HaveOccurred())
				Expect(ctx.svc).NotTo(BeNil())
				Expect(errBuf.Len()).To(Equal(0))
				DeferCleanup(func() {
					ctx.Close()
				})
			})

			It("should create database file", func() {
				tmpDir := GinkgoT().TempDir()
				dbPath := filepath.Join(tmpDir, "test.db")

				ctx := NewCLIContext(dbPath, false)
				errBuf := new(bytes.Buffer)

				err := ctx.InitService(errBuf)
				Expect(err).NotTo(HaveOccurred())
				DeferCleanup(func() {
					ctx.Close()
				})

				_, err = os.Stat(dbPath)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should initialize all repositories", func() {
				tmpDir := GinkgoT().TempDir()
				dbPath := filepath.Join(tmpDir, "test.db")

				ctx := NewCLIContext(dbPath, false)
				errBuf := new(bytes.Buffer)

				err := ctx.InitService(errBuf)
				Expect(err).NotTo(HaveOccurred())
				DeferCleanup(func() {
					ctx.Close()
				})

				Expect(ctx.svc.GetEventRepository()).NotTo(BeNil())
				Expect(ctx.svc.GetFactRepository()).NotTo(BeNil())
				Expect(ctx.svc.GetBurstRepository()).NotTo(BeNil())
				Expect(ctx.svc.GetSkillRepository()).NotTo(BeNil())
			})
		})

		Context("with SQLite and default path", func() {
			It("should initialize service successfully", func() {
				homeDir := GinkgoT().TempDir()
				originalHome := os.Getenv("HOME")
				originalUserProfile := os.Getenv("USERPROFILE")
				DeferCleanup(func() {
					os.Setenv("HOME", originalHome)
					os.Setenv("USERPROFILE", originalUserProfile)
				})
				os.Setenv("HOME", homeDir)
				os.Setenv("USERPROFILE", homeDir)

				ctx := NewCLIContext("", false)
				errBuf := new(bytes.Buffer)

				err := ctx.InitService(errBuf)
				Expect(err).NotTo(HaveOccurred())
				Expect(ctx.svc).NotTo(BeNil())
				Expect(errBuf.Len()).To(Equal(0))
				DeferCleanup(func() {
					ctx.Close()
				})
			})

			It("should create database at default path", func() {
				homeDir := GinkgoT().TempDir()
				originalHome := os.Getenv("HOME")
				originalUserProfile := os.Getenv("USERPROFILE")
				DeferCleanup(func() {
					os.Setenv("HOME", originalHome)
					os.Setenv("USERPROFILE", originalUserProfile)
				})
				os.Setenv("HOME", homeDir)
				os.Setenv("USERPROFILE", homeDir)

				ctx := NewCLIContext("", false)
				errBuf := new(bytes.Buffer)

				err := ctx.InitService(errBuf)
				Expect(err).NotTo(HaveOccurred())
				DeferCleanup(func() {
					ctx.Close()
				})

				expectedPath := filepath.Join(homeDir, ".kariya", "events.db")
				_, err = os.Stat(expectedPath)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should create .kariya directory with correct permissions", func() {
				if runtime.GOOS == "windows" {
					Skip("Windows does not enforce Unix permissions")
				}
				homeDir := GinkgoT().TempDir()
				originalHome := os.Getenv("HOME")
				originalUserProfile := os.Getenv("USERPROFILE")
				DeferCleanup(func() {
					os.Setenv("HOME", originalHome)
					os.Setenv("USERPROFILE", originalUserProfile)
				})
				os.Setenv("HOME", homeDir)
				os.Setenv("USERPROFILE", homeDir)

				ctx := NewCLIContext("", false)
				errBuf := new(bytes.Buffer)

				err := ctx.InitService(errBuf)
				Expect(err).NotTo(HaveOccurred())
				DeferCleanup(func() {
					ctx.Close()
				})

				kariyaDir := filepath.Join(homeDir, ".kariya")
				info, err := os.Stat(kariyaDir)
				Expect(err).NotTo(HaveOccurred())
				Expect(info.IsDir()).To(BeTrue())
				Expect(info.Mode().Perm()).To(Equal(os.FileMode(0o750)))
			})

			It("should initialize all repositories", func() {
				homeDir := GinkgoT().TempDir()
				originalHome := os.Getenv("HOME")
				originalUserProfile := os.Getenv("USERPROFILE")
				DeferCleanup(func() {
					os.Setenv("HOME", originalHome)
					os.Setenv("USERPROFILE", originalUserProfile)
				})
				os.Setenv("HOME", homeDir)
				os.Setenv("USERPROFILE", homeDir)

				ctx := NewCLIContext("", false)
				errBuf := new(bytes.Buffer)

				err := ctx.InitService(errBuf)
				Expect(err).NotTo(HaveOccurred())
				DeferCleanup(func() {
					ctx.Close()
				})

				Expect(ctx.svc.GetEventRepository()).NotTo(BeNil())
				Expect(ctx.svc.GetFactRepository()).NotTo(BeNil())
				Expect(ctx.svc.GetBurstRepository()).NotTo(BeNil())
				Expect(ctx.svc.GetSkillRepository()).NotTo(BeNil())
			})
		})

		Context("with invalid database path", func() {
			It("should return error", func() {
				ctx := NewCLIContext("/invalid/path/that/does/not/exist/test.db", false)
				errBuf := new(bytes.Buffer)

				err := ctx.InitService(errBuf)
				Expect(err).To(HaveOccurred())
				Expect(ctx.svc).To(BeNil())
				Expect(errBuf.Len()).To(BeNumerically(">", 0))
			})
		})

		Context("with directory creation error", func() {
			It("should return error when directory cannot be created", func() {
				if runtime.GOOS == "windows" {
					Skip("Windows does not enforce Unix permissions")
				}
				homeDir := GinkgoT().TempDir()
				originalHome := os.Getenv("HOME")
				originalUserProfile := os.Getenv("USERPROFILE")
				DeferCleanup(func() {
					os.Setenv("HOME", originalHome)
					os.Setenv("USERPROFILE", originalUserProfile)
					os.Chmod(filepath.Join(homeDir, "readonly"), 0o755)
				})
				os.Setenv("HOME", homeDir)
				os.Setenv("USERPROFILE", homeDir)

				readOnlyDir := filepath.Join(homeDir, "readonly")
				Expect(os.Mkdir(readOnlyDir, 0o555)).To(Succeed())

				ctx := NewCLIContext(filepath.Join(readOnlyDir, ".kariya", "events.db"), false)
				errBuf := new(bytes.Buffer)

				err := ctx.InitService(errBuf)
				Expect(err).To(HaveOccurred())
				Expect(ctx.svc).To(BeNil())
				Expect(errBuf.Len()).To(BeNumerically(">", 0))
				DeferCleanup(func() {
					ctx.Close()
				})
			})
		})

		Context("with bad HOME directory", func() {
			It("should return error for nonexistent HOME", func() {
				if runtime.GOOS == "windows" {
					Skip("Windows resolves Unix-style paths differently")
				}
				originalHome := os.Getenv("HOME")
				originalUserProfile := os.Getenv("USERPROFILE")
				DeferCleanup(func() {
					os.Setenv("HOME", originalHome)
					os.Setenv("USERPROFILE", originalUserProfile)
				})
				os.Setenv("HOME", "/nonexistent/path/that/cannot/be/created")
				os.Setenv("USERPROFILE", "/nonexistent/path/that/cannot/be/created")

				ctx := NewCLIContext("", false)
				errBuf := new(bytes.Buffer)

				err := ctx.InitService(errBuf)
				Expect(err).To(HaveOccurred())
				Expect(ctx.svc).To(BeNil())
				Expect(errBuf.Len()).To(BeNumerically(">", 0))
			})
		})

		Context("with empty HOME directory", func() {
			It("should return error for empty HOME", func() {
				originalHome := os.Getenv("HOME")
				originalUserProfile := os.Getenv("USERPROFILE")
				DeferCleanup(func() {
					os.Setenv("HOME", originalHome)
					os.Setenv("USERPROFILE", originalUserProfile)
				})
				os.Setenv("HOME", "")
				os.Setenv("USERPROFILE", "")

				ctx := NewCLIContext("", false)
				errBuf := new(bytes.Buffer)

				err := ctx.InitService(errBuf)
				Expect(err).To(HaveOccurred())
				Expect(ctx.svc).To(BeNil())
				Expect(errBuf.Len()).To(BeNumerically(">", 0))
			})
		})

		Context("with corrupted database", func() {
			It("should return error for corrupted DB file", func() {
				tmpDir := GinkgoT().TempDir()
				dbPath := filepath.Join(tmpDir, "corrupted.db")

				Expect(os.WriteFile(dbPath, []byte("not a valid sqlite database"), 0o600)).To(Succeed())

				ctx := NewCLIContext(dbPath, false)
				errBuf := new(bytes.Buffer)

				err := ctx.InitService(errBuf)
				Expect(err).To(HaveOccurred())
				Expect(ctx.svc).To(BeNil())
				Expect(errBuf.Len()).To(BeNumerically(">", 0))
			})
		})

		Context("with directory as database path", func() {
			It("should return error when DB path is a directory", func() {
				tmpDir := GinkgoT().TempDir()
				dbPath := filepath.Join(tmpDir, "db_dir")

				Expect(os.Mkdir(dbPath, 0o755)).To(Succeed())

				ctx := NewCLIContext(dbPath, false)
				errBuf := new(bytes.Buffer)

				err := ctx.InitService(errBuf)
				Expect(err).To(HaveOccurred())
				Expect(ctx.svc).To(BeNil())
				Expect(errBuf.Len()).To(BeNumerically(">", 0))
			})
		})

		Context("with readonly database path", func() {
			It("should return error for readonly DB path", func() {
				if runtime.GOOS == "windows" {
					Skip("Windows does not enforce Unix permissions")
				}
				tmpDir := GinkgoT().TempDir()
				readOnlyDir := filepath.Join(tmpDir, "readonly")
				DeferCleanup(func() {
					os.Chmod(readOnlyDir, 0o755)
				})

				Expect(os.Mkdir(readOnlyDir, 0o555)).To(Succeed())

				dbPath := filepath.Join(readOnlyDir, "test.db")

				ctx := NewCLIContext(dbPath, false)
				errBuf := new(bytes.Buffer)

				err := ctx.InitService(errBuf)
				Expect(err).To(HaveOccurred())
				Expect(ctx.svc).To(BeNil())
				Expect(errBuf.Len()).To(BeNumerically(">", 0))
			})
		})

		Context("with readonly home directory", func() {
			It("should return error for readonly HOME", func() {
				if runtime.GOOS == "windows" {
					Skip("Windows does not enforce Unix permissions")
				}
				tmpDir := GinkgoT().TempDir()
				readOnlyHome := filepath.Join(tmpDir, "readonly_home")
				originalHome := os.Getenv("HOME")
				originalUserProfile := os.Getenv("USERPROFILE")
				DeferCleanup(func() {
					os.Chmod(readOnlyHome, 0o755)
					os.Setenv("HOME", originalHome)
					os.Setenv("USERPROFILE", originalUserProfile)
				})

				Expect(os.Mkdir(readOnlyHome, 0o555)).To(Succeed())
				os.Setenv("HOME", readOnlyHome)
				os.Setenv("USERPROFILE", readOnlyHome)

				ctx := NewCLIContext("", false)
				errBuf := new(bytes.Buffer)

				err := ctx.InitService(errBuf)
				Expect(err).To(HaveOccurred())
				Expect(ctx.svc).To(BeNil())
				Expect(errBuf.Len()).To(BeNumerically(">", 0))
			})
		})

		Context("with invalid database path (dev/null)", func() {
			It("should return error for /dev/null path", func() {
				ctx := NewCLIContext("/dev/null/invalid/path/test.db", false)
				errBuf := new(bytes.Buffer)

				err := ctx.InitService(errBuf)
				Expect(err).To(HaveOccurred())
				Expect(ctx.svc).To(BeNil())
				Expect(errBuf.Len()).To(BeNumerically(">", 0))
			})
		})
	})

	Describe("Service", func() {
		Context("before initialization", func() {
			It("should return nil", func() {
				ctx := NewCLIContext("", true)
				Expect(ctx.Service()).To(BeNil())
			})
		})

		Context("after initialization", func() {
			It("should return the initialized service", func() {
				ctx := NewCLIContext("", true)
				errBuf := new(bytes.Buffer)

				err := ctx.InitService(errBuf)
				Expect(err).NotTo(HaveOccurred())

				svc := ctx.Service()
				Expect(svc).NotTo(BeNil())
				Expect(svc).To(Equal(ctx.svc))
			})
		})
	})

	Describe("Multiple InitService calls", func() {
		It("should create a new service on each call", func() {
			ctx := NewCLIContext("", true)
			errBuf := new(bytes.Buffer)

			err1 := ctx.InitService(errBuf)
			Expect(err1).NotTo(HaveOccurred())
			svc1 := ctx.svc

			err2 := ctx.InitService(errBuf)
			Expect(err2).NotTo(HaveOccurred())
			svc2 := ctx.svc

			Expect(svc1).NotTo(Equal(svc2))
		})
	})
})
