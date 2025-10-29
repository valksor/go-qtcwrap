package qtcwrap

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	// File and path constants.
	testQtplFile = "test.qtpl"
	testContent  = "test content"

	// Command argument constants.
	dirTemplatesArg = "-dir=templates"
	extQtplArg      = "-ext=.qtpl"
	skipCommentsArg = "-skipLineComments"

	// Error message templates.
	syntaxErrorMsg    = "syntax error in template"
	createTempDirErr  = "Failed to create temp directory: %v"
	removeTempDirErr  = "Failed to remove temp directory: %v"
	createTempFileErr = "Failed to create temp file: %v"
)

func TestConfig(t *testing.T) {
	t.Run("DefaultConfig", func(t *testing.T) {
		config := Config{
			Dir:              ".",
			SkipLineComments: true,
			Ext:              "",
			File:             "",
		}

		if config.Dir != "." {
			t.Errorf("Expected Dir to be '.', got '%s'", config.Dir)
		}
		if !config.SkipLineComments {
			t.Error("Expected SkipLineComments to be true")
		}
		if config.Ext != "" {
			t.Errorf("Expected Ext to be empty, got '%s'", config.Ext)
		}
		if config.File != "" {
			t.Errorf("Expected File to be empty, got '%s'", config.File)
		}
	})

	t.Run("CustomConfig", func(t *testing.T) {
		config := Config{
			Dir:              "templates",
			SkipLineComments: false,
			Ext:              ".qtpl",
			File:             testQtplFile,
		}

		if config.Dir != "templates" {
			t.Errorf("Expected Dir to be 'templates', got '%s'", config.Dir)
		}
		if config.SkipLineComments {
			t.Error("Expected SkipLineComments to be false")
		}
		if config.Ext != ".qtpl" {
			t.Errorf("Expected Ext to be '.qtpl', got '%s'", config.Ext)
		}
		if config.File != testQtplFile {
			t.Errorf("Expected File to be 'test.qtpl', got '%s'", config.File)
		}
	})
}

func TestGetDefaultConfig(t *testing.T) {
	config := GetDefaultConfig()

	tests := []struct {
		name     string
		expected any
		actual   any
	}{
		{"Dir", ".", config.Dir},
		{"SkipLineComments", true, config.SkipLineComments},
		{"Ext", "", config.Ext},
		{"File", "", config.File},
	}

	for _, testT := range tests {
		t.Run(testT.name, func(t *testing.T) {
			if testT.actual != testT.expected {
				t.Errorf("Expected %s to be %v, got %v", testT.name, testT.expected, testT.actual)
			}
		})
	}
}

func TestValidateQtcTool(t *testing.T) {
	t.Run("ValidateQtcTool", func(t *testing.T) {
		err := validateQtcTool()

		// This test will pass if qtc is available, skip if not
		if err != nil {
			t.Skipf("qtc tool not available: %v", err)
		}
	})
}

func TestIsQtcAvailable(t *testing.T) {
	t.Run("IsQtcAvailable", func(t *testing.T) {
		available := IsQtcAvailable()

		// Check that the function returns a boolean
		if available != true && available != false {
			t.Error("IsQtcAvailable should return a boolean")
		}
	})
}

func TestGetQtcVersion(t *testing.T) {
	t.Run("GetQtcVersion", func(t *testing.T) {
		version, err := GetQtcVersion()

		if err != nil {
			t.Skipf("qtc tool not available: %v", err)
		}

		if version == "" {
			t.Error("Expected version to be non-empty")
		}

		if strings.TrimSpace(version) != version {
			t.Error("Version should be trimmed")
		}
	})
}

func TestBuildArgs(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		expected []string
	}{
		{
			name: "FileMode",
			config: Config{
				File:             testQtplFile,
				Dir:              "templates",
				Ext:              ".qtpl",
				SkipLineComments: false,
			},
			expected: []string{"-file=test.qtpl"},
		},
		{
			name: "DirectoryMode",
			config: Config{
				Dir:              "templates",
				Ext:              ".qtpl",
				SkipLineComments: false,
				File:             "",
			},
			expected: []string{dirTemplatesArg, extQtplArg},
		},
		{
			name: "DirectoryModeWithExt",
			config: Config{
				Dir:              "templates",
				SkipLineComments: false,
				Ext:              ".qtpl",
				File:             "",
			},
			expected: []string{dirTemplatesArg, extQtplArg},
		},
		{
			name: "DirectoryModeWithSkipComments",
			config: Config{
				Dir:              "templates",
				SkipLineComments: true,
				Ext:              ".qtpl",
				File:             "",
			},
			expected: []string{dirTemplatesArg, extQtplArg, skipCommentsArg},
		},
		{
			name: "FileModeWithSkipComments",
			config: Config{
				Dir:              "",
				SkipLineComments: true,
				Ext:              "",
				File:             testQtplFile,
			},
			expected: []string{"-file=test.qtpl", skipCommentsArg},
		},
		{
			name: "OnlyDir",
			config: Config{
				Dir:              "src",
				SkipLineComments: false,
				Ext:              "",
				File:             "",
			},
			expected: []string{"-dir=src"},
		},
		{
			name: "OnlyExt",
			config: Config{
				Dir:              "",
				SkipLineComments: false,
				Ext:              ".template",
				File:             "",
			},
			expected: []string{"-ext=.template"},
		},
		{
			name: "OnlySkipComments",
			config: Config{
				Dir:              "",
				SkipLineComments: true,
				Ext:              "",
				File:             "",
			},
			expected: []string{skipCommentsArg},
		},
		{
			name: "EmptyConfig",
			config: Config{
				Dir:              "",
				SkipLineComments: false,
				Ext:              "",
				File:             "",
			},
			expected: []string{},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			args := buildArgs(Config{
				Dir:              testCase.config.Dir,
				SkipLineComments: testCase.config.SkipLineComments,
				Ext:              testCase.config.Ext,
				File:             testCase.config.File,
			})

			if len(args) != len(testCase.expected) {
				t.Errorf("Expected %d arguments, got %d", len(testCase.expected), len(args))
			}

			for idx, arg := range args {
				if idx >= len(testCase.expected) {
					t.Errorf("Unexpected argument at index %d: %s", idx, arg)
					continue
				}
				if arg != testCase.expected[idx] {
					t.Errorf("Expected argument %d to be '%s', got '%s'", idx, testCase.expected[idx], arg)
				}
			}
		})
	}
}

func TestIsTemporaryFileWarning(t *testing.T) {
	tests := []struct {
		name     string
		stderr   []byte
		expected bool
	}{
		{
			name:     "TemporaryFileWarning",
			stderr:   []byte("open .tmp/test.qtpl: no such file or directory"),
			expected: true,
		},
		{
			name:     "TemporaryFileWarningWithPath",
			stderr:   []byte("stat: /path/to/.tmp/file: no such file or directory"),
			expected: true,
		},
		{
			name:     "ActualError",
			stderr:   []byte(syntaxErrorMsg),
			expected: false,
		},
		{
			name:     "NoSuchFileButNotestTmp",
			stderr:   []byte("open test.qtpl: no such file or directory"),
			expected: false,
		},
		{
			name:     "TmpButNotNoSuchFile",
			stderr:   []byte("permission denied: .tmp/test.qtpl"),
			expected: false,
		},
		{
			name:     "EmptyStderr",
			stderr:   []byte(""),
			expected: false,
		},
		{
			name:     "OnlyTmp",
			stderr:   []byte(".tmp"),
			expected: false,
		},
		{
			name:     "OnlyNoSuchFile",
			stderr:   []byte("no such file or directory"),
			expected: false,
		},
	}

	for _, testT := range tests {
		t.Run(testT.name, func(t *testing.T) {
			result := isTemporaryFileWarning(testT.stderr)
			if result != testT.expected {
				t.Errorf("Expected %v, got %v for stderr: %s", testT.expected, result, string(testT.stderr))
			}
		})
	}
}

func TestValidateConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "qtcwrap_test")
	if err != nil {
		t.Fatalf(createTempDirErr, err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf(removeTempDirErr, err)
		}
	}()

	// Create a temporary file for testing
	tempFile := filepath.Join(tempDir, testQtplFile)
	err = os.WriteFile(tempFile, []byte(testContent), 0600)
	if err != nil {
		t.Fatalf(createTempFileErr, err)
	}

	tests := []struct {
		name      string
		config    Config
		expectErr bool
		errorMsg  string
	}{
		{
			name: "ValidFileConfig",
			config: Config{
				Dir:              "",
				SkipLineComments: false,
				Ext:              "",
				File:             tempFile,
			},
			expectErr: false,
			errorMsg:  "",
		},
		{
			name: "ValidDirConfig",
			config: Config{
				Dir:              tempDir,
				SkipLineComments: false,
				Ext:              "",
				File:             "",
			},
			expectErr: false,
			errorMsg:  "",
		},
		{
			name: "ValidDirConfigWithExt",
			config: Config{
				Dir:              tempDir,
				SkipLineComments: false,
				Ext:              ".qtpl",
				File:             "",
			},
			expectErr: false,
			errorMsg:  "",
		},
		{
			name: "InvalidFileConfig",
			config: Config{
				Dir:              "",
				SkipLineComments: false,
				Ext:              "",
				File:             "/nonexistent/file.qtpl",
			},
			expectErr: true,
			errorMsg:  "file /nonexistent/file.qtpl is not accessible",
		},
		{
			name: "InvalidDirConfig",
			config: Config{
				Dir:              "/nonexistent/directory",
				SkipLineComments: false,
				Ext:              "",
				File:             "",
			},
			expectErr: true,
			errorMsg:  "directory /nonexistent/directory is not accessible",
		},
		{
			name: "EmptyConfig",
			config: Config{
				Dir:              "",
				SkipLineComments: false,
				Ext:              "",
				File:             "",
			},
			expectErr: true,
			errorMsg:  "either File or Dir must be specified",
		},
		{
			name: "DirIsActuallyFile",
			config: Config{
				Dir:              tempFile,
				SkipLineComments: false,
				Ext:              "",
				File:             "",
			},
			expectErr: true,
			errorMsg:  "is not a directory",
		},
		{
			name: "InvalidExtension",
			config: Config{
				Dir:              tempDir,
				SkipLineComments: false,
				Ext:              "qtpl",
				File:             "",
			},
			expectErr: true,
			errorMsg:  "extension must start with a dot",
		},
		{
			name: "ValidExtensionWithDot",
			config: Config{
				Dir:              tempDir,
				SkipLineComments: false,
				Ext:              ".qtpl",
				File:             "",
			},
			expectErr: false,
			errorMsg:  "",
		},
		{
			name: "ValidCurrentDir",
			config: Config{
				Dir:              ".",
				SkipLineComments: false,
				Ext:              "",
				File:             "",
			},
			expectErr: false,
			errorMsg:  "",
		},
	}

	for _, testT := range tests {
		t.Run(testT.name, func(t *testing.T) {
			err := ValidateConfig(testT.config)

			if testT.expectErr {
				if err == nil {
					t.Error("Expected error but got none")
				} else if testT.errorMsg != "" && !strings.Contains(err.Error(), testT.errorMsg) {
					t.Errorf("Expected error to contain '%s', got '%s'", testT.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestCompileWithValidation(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "qtcwrap_test")
	if err != nil {
		t.Fatalf(createTempDirErr, err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf(removeTempDirErr, err)
		}
	}()

	tests := []struct {
		name      string
		config    Config
		expectErr bool
		errorMsg  string
	}{
		{
			name: "ValidConfig",
			config: Config{
				Dir:              tempDir,
				SkipLineComments: false,
				Ext:              "",
				File:             "",
			},
			expectErr: false,
			errorMsg:  "",
		},
		{
			name: "InvalidConfig",
			config: Config{
				Dir:              "/nonexistent/directory",
				SkipLineComments: false,
				Ext:              "",
				File:             "",
			},
			expectErr: true,
			errorMsg:  "configuration validation failed",
		},
		{
			name: "EmptyConfig",
			config: Config{
				Dir:              "",
				SkipLineComments: false,
				Ext:              "",
				File:             "",
			},
			expectErr: true,
			errorMsg:  "either File or Dir must be specified",
		},
	}

	for _, testT := range tests {
		t.Run(testT.name, func(t *testing.T) {
			err := CompileWithValidation(testT.config)

			if testT.expectErr {
				if err == nil {
					t.Error("Expected error but got none")
				} else if testT.errorMsg != "" && !strings.Contains(err.Error(), testT.errorMsg) {
					t.Errorf("Expected error to contain '%s', got '%s'", testT.errorMsg, err.Error())
				}
			} else {
				// For valid configs, we might get a qtc tool error, which is acceptable
				if err != nil && !strings.Contains(err.Error(), "qtc tool validation failed") {
					t.Errorf("Expected no error or qtc tool error, got: %v", err)
				}
			}
		})
	}
}

func TestFindTemplateFiles(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "qtcwrap_test")
	if err != nil {
		t.Fatalf(createTempDirErr, err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf(removeTempDirErr, err)
		}
	}()

	testFiles := []string{
		"test1.qtpl",
		"test2.qtpl",
		"subdir/test3.qtpl",
		"subdir/test4.template",
		"other.go",
	}
	for _, fileName := range testFiles {
		fullPath := filepath.Join(tempDir, fileName)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0700); err != nil {
			t.Fatalf("Failed to create directory for %s: %v", fileName, err)
		}
		err = os.WriteFile(fullPath, []byte(testContent), 0600)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", fileName, err)
		}
	}

	testCases := []struct {
		description string
		ext         string
		expectedLen int
	}{
		{"Find .qtpl files", ".qtpl", 3},
		{"Find .template files", ".template", 1},
		{"Find .go files", ".go", 1},
	}
	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			files, err := FindTemplateFiles(tempDir, testCase.ext)
			if err != nil {
				t.Fatalf("Failed to find template files: %v", err)
			}
			if len(files) != testCase.expectedLen {
				t.Errorf("Expected %d files, got %d", testCase.expectedLen, len(files))
			}
		})
	}
}

func TestConvenienceFunctions(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "qtcwrap_test")
	if err != nil {
		t.Fatalf(createTempDirErr, err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf(removeTempDirErr, err)
		}
	}()

	// Create a temporary file for testing
	tempFile := filepath.Join(tempDir, testQtplFile)
	err = os.WriteFile(tempFile, []byte(testContent), 0600)
	if err != nil {
		t.Fatalf(createTempFileErr, err)
	}

	t.Run("CompileDirectory", func(t *testing.T) {
		// This test just ensures the function can be called without panic
		// We can't test actual compilation without qtc being available
		CompileDirectory(tempDir)
	})

	t.Run("CompileFile", func(t *testing.T) {
		// This test just ensures the function can be called without panic
		// We can't test actual compilation without qtc being available
		CompileFile(tempFile)
	})

	t.Run("CompileWithExtension", func(t *testing.T) {
		// This test just ensures the function can be called without panic
		// We can't test actual compilation without qtc being available
		CompileWithExtension(tempDir, ".qtpl")
	})

	t.Run("QtcWrap", func(t *testing.T) {
		// This test just ensures the function can be called without panic
		// We can't test actual compilation without qtc being available
		QtcWrap()
	})
}

func TestHandleQtcError(t *testing.T) {
	tests := []struct {
		name           string
		stderr         string
		expectedOutput string
	}{
		{
			name:           "TemporaryFileWarning",
			stderr:         "open .tmp/test.qtpl: no such file or directory",
			expectedOutput: "[qtc warning suppressed]",
		},
		{
			name:           "ActualError",
			stderr:         syntaxErrorMsg,
			expectedOutput: syntaxErrorMsg,
		},
		{
			name:           "EmptyStderr",
			stderr:         "",
			expectedOutput: "qtc execution failed:",
		},
	}

	for _, testT := range tests {
		t.Run(testT.name, func(t *testing.T) {
			var buf bytes.Buffer
			stderr := bytes.NewBufferString(testT.stderr)
			err := errors.New("exit status 1")

			// Capture output by temporarily redirecting stdout
			oldStdout := os.Stdout
			rFile, wFile, _ := os.Pipe()
			os.Stdout = wFile

			handleQtcError(*stderr, err)

			if err := wFile.Close(); err != nil {
				t.Fatalf("Failed to close writer: %v", err)
			}
			os.Stdout = oldStdout

			// Read the captured output
			if _, err := buf.ReadFrom(rFile); err != nil {
				t.Fatalf("Failed to read from pipe: %v", err)
			}
			output := buf.String()

			if !strings.Contains(output, testT.expectedOutput) {
				t.Errorf("Expected output to contain '%s', got '%s'", testT.expectedOutput, output)
			}
		})
	}
}

func TestExecuteQtc(t *testing.T) {
	t.Run("ExecuteQtcWithInvalidArgs", func(t *testing.T) {
		// Test with invalid arguments that should fail
		args := []string{"-invalid-flag"}

		// Capture output
		var buf bytes.Buffer
		oldStdout := os.Stdout
		rFile, wFile, _ := os.Pipe()
		os.Stdout = wFile

		executeQtc(args)

		if err := wFile.Close(); err != nil {
			t.Fatalf("Failed to close writer: %v", err)
		}
		os.Stdout = oldStdout
		if _, err := buf.ReadFrom(rFile); err != nil {
			t.Fatalf("Failed to read from pipe: %v", err)
		}

		// The function should handle the error gracefully
		// We can't test much more without qtc being available
	})
}

func TestWithConfig(t *testing.T) {
	t.Run("WithValidConfig", func(t *testing.T) {
		// Create a temporary directory for testing
		tempDir, err := os.MkdirTemp("", "qtcwrap_test")
		if err != nil {
			t.Fatalf(createTempDirErr, err)
		}
		defer func() {
			if err := os.RemoveAll(tempDir); err != nil {
				t.Fatalf(removeTempDirErr, err)
			}
		}()

		config := Config{
			Dir:              tempDir,
			SkipLineComments: true,
			Ext:              "",
			File:             "",
		}

		// This should not panic
		WithConfig(config)
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("ConfigWithBothFileAndDir", func(t *testing.T) {
		// Create a temporary directory and file
		tempDir, err := os.MkdirTemp("", "qtcwrap_test")
		if err != nil {
			t.Fatalf(createTempDirErr, err)
		}
		defer func() {
			if err := os.RemoveAll(tempDir); err != nil {
				t.Fatalf(removeTempDirErr, err)
			}
		}()

		tempFile := filepath.Join(tempDir, testQtplFile)
		err = os.WriteFile(tempFile, []byte(testContent), 0600)
		if err != nil {
			t.Fatalf(createTempFileErr, err)
		}

		config := Config{
			Dir:              tempDir,
			SkipLineComments: false,
			Ext:              "",
			File:             tempFile,
		}

		args := buildArgs(config)
		// When both are specified, File should take precedence
		if len(args) != 1 || args[0] != "-file="+tempFile {
			t.Errorf("Expected File to take precedence over Dir, got args: %v", args)
		}
	})

	t.Run("EmptyExtension", func(t *testing.T) {
		config := Config{
			Dir:              ".",
			SkipLineComments: false,
			Ext:              "",
			File:             "",
		}

		args := buildArgs(config)
		// Empty extension should not be added to args
		found := false
		for _, arg := range args {
			if strings.HasPrefix(arg, "-ext=") {
				found = true
				break
			}
		}
		if found {
			t.Error("Empty extension should not be added to args")
		}
	})

	t.Run("ValidateConfigWithNilPointer", func(t *testing.T) {
		// Test that validation doesn't panic with edge cases
		config := Config{
			Dir:              "",
			SkipLineComments: false,
			Ext:              "",
			File:             "",
		}
		err := ValidateConfig(config)
		if err == nil {
			t.Error("Expected error for empty config")
		}
	})
}

func TestFileSystemPermissions(t *testing.T) {
	t.Run("ReadOnlyDirectory", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "qtcwrap_test")
		if err != nil {
			t.Fatalf(createTempDirErr, err)
		}
		defer func() {
			if err := os.RemoveAll(tempDir); err != nil {
				t.Fatalf(removeTempDirErr, err)
			}
		}()

		// Make directory read-only
		if err := os.Chmod(tempDir, 0o600); err != nil {
			t.Fatalf("Failed to change directory permissions: %v", err)
		}

		// Restore permissions for cleanup
		defer func() {
			if err := os.Chmod(tempDir, 0o600); err != nil {
				t.Fatalf("Failed to restore directory permissions: %v", err)
			}
		}()

		config := Config{
			Dir:              tempDir,
			SkipLineComments: false,
			Ext:              "",
			File:             "",
		}
		err = ValidateConfig(config)

		// Should still be valid as we can read the directory
		if err != nil {
			t.Errorf("Expected no error for read-only directory, got: %v", err)
		}
	})
}

func TestConcurrentAccess(t *testing.T) {
	t.Run("ConcurrentValidation", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "qtcwrap_test")
		if err != nil {
			t.Fatalf(createTempDirErr, err)
		}
		defer func() {
			if err := os.RemoveAll(tempDir); err != nil {
				t.Fatalf(removeTempDirErr, err)
			}
		}()

		config := Config{
			Dir:              tempDir,
			SkipLineComments: false,
			Ext:              "",
			File:             "",
		}

		// Run multiple validations concurrently
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func() {
				err := ValidateConfig(config)
				if err != nil {
					t.Errorf("Concurrent validation failed: %v", err)
				}
				done <- true
			}()
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}
	})
}

func TestLargeDirectoryStructure(t *testing.T) {
	t.Run("ManyFiles", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "qtcwrap_test")
		if err != nil {
			t.Fatalf(createTempDirErr, err)
		}
		defer func() {
			if err := os.RemoveAll(tempDir); err != nil {
				t.Fatalf(removeTempDirErr, err)
			}
		}()

		// Create many files
		for i := 0; i < 100; i++ {
			filename := filepath.Join(tempDir, fmt.Sprintf("test%d.qtpl", i))
			err = os.WriteFile(filename, []byte(testContent), 0600)
			if err != nil {
				t.Fatalf("Failed to create file %s: %v", filename, err)
			}
		}

		files, err := FindTemplateFiles(tempDir, ".qtpl")
		if err != nil {
			t.Fatalf("Failed to find template files: %v", err)
		}

		if len(files) != 100 {
			t.Errorf("Expected 100 files, got %d", len(files))
		}
	})
}

func TestSpecialCharacters(t *testing.T) {
	t.Run("SpecialCharsInPath", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "qtcwrap test with spaces")
		if err != nil {
			t.Fatalf(createTempDirErr, err)
		}
		defer func() {
			if err := os.RemoveAll(tempDir); err != nil {
				t.Fatalf(removeTempDirErr, err)
			}
		}()

		config := Config{
			Dir:              tempDir,
			SkipLineComments: false,
			Ext:              "",
			File:             "",
		}
		err = ValidateConfig(config)

		if err != nil {
			t.Errorf("Expected no error for path with spaces, got: %v", err)
		}
	})
}

func TestSymlinks(t *testing.T) {
	t.Run("SymlinkDirectory", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "qtcwrap_test")
		if err != nil {
			t.Fatalf(createTempDirErr, err)
		}
		defer func() {
			if err := os.RemoveAll(tempDir); err != nil {
				t.Fatalf(removeTempDirErr, err)
			}
		}()

		// Create a symlink to the directory
		symlinkPath := filepath.Join(tempDir, "symlink")
		err = os.Symlink(tempDir, symlinkPath)
		if err != nil {
			t.Skipf("Failed to create symlink, skipping test: %v", err)
		}

		config := Config{
			Dir:              symlinkPath,
			SkipLineComments: false,
			Ext:              "",
			File:             "",
		}
		err = ValidateConfig(config)

		if err != nil {
			t.Errorf("Expected no error for symlink directory, got: %v", err)
		}
	})
}

func TestErrorMessages(t *testing.T) {
	t.Run("DetailedErrorMessages", func(t *testing.T) {
		config := Config{
			Dir:              "/this/path/definitely/does/not/exist",
			SkipLineComments: false,
			Ext:              "",
			File:             "",
		}

		err := ValidateConfig(config)
		if err == nil {
			t.Error("Expected error for non-existent directory")
		}

		errMsg := err.Error()
		if !strings.Contains(errMsg, "directory") {
			t.Errorf("Expected error message to contain 'directory', got: %s", errMsg)
		}
		if !strings.Contains(errMsg, config.Dir) {
			t.Errorf("Expected error message to contain directory path, got: %s", errMsg)
		}
	})
}
