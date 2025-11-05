# QuickTemplate Wrapper (qtcwrap)

[![BSD-3-Clause](https://img.shields.io/badge/BSD--3--Clause-green?style=flat)](https://github.com/valksor/php-bundle/blob/master/LICENSE)
[![Coverage Status](https://coveralls.io/repos/github/valksor/go-qtcwrap/badge.svg?branch=master)](https://coveralls.io/github/valksor/go-qtcwrap?branch=master)

A Go package that provides utilities for working with the QuickTemplate compiler (qtc), simplifying the process of generating Go code from QuickTemplate (.qtpl) files.

## Features

- **Simple API**: Easy-to-use Go API wrapper around the qtc command-line tool
- **Flexible Configuration**: Support for both directory-based and single-file compilation
- **Error Handling**: Intelligent error filtering and proper warning suppression
- **Validation**: Configuration validation to ensure proper setup
- **Convenience Functions**: Multiple helper functions for common use cases
- **Template Discovery**: Built-in template file discovery utilities

## Installation

```bash
go get github.com/valksor/go-qtcwrap
```

## Prerequisites

This package requires the `qtc` (QuickTemplate compiler) tool to be installed and available in your system PATH.

Install qtc:
```bash
go install github.com/valyala/quicktemplate/qtc@latest
```

## Quick Start

### Basic Usage

```go
package main

import "github.com/valksor/go-qtcwrap"

func main() {
    // Compile all .qtpl files in current directory
    qtcwrap.QtcWrap()
    
    // Or specify a directory
    qtcwrap.CompileDirectory("templates")
    
    // Or compile a single file
    qtcwrap.CompileFile("templates/home.qtpl")
}
```

### Configuration-Based Usage

```go
package main

import "github.com/valksor/go-qtcwrap"

func main() {
    config := qtcwrap.Config{
        Dir:              "templates",
        SkipLineComments: true,
        Ext:              ".qtpl",
    }
    
    qtcwrap.WithConfig(config)
}
```

### Validation and Error Handling

```go
package main

import (
    "fmt"
    "github.com/valksor/go-qtcwrap"
)

func main() {
    config := qtcwrap.Config{
        Dir:              "templates",
        SkipLineComments: true,
        Ext:              ".qtpl",
    }
    
    // Validate configuration before compilation
    if err := qtcwrap.CompileWithValidation(config); err != nil {
        fmt.Printf("Compilation failed: %v\n", err)
        return
    }
    
    fmt.Println("Templates compiled successfully!")
}
```

## Configuration

### Config Structure

```go
type Config struct {
    // Directory containing .qtpl files (ignored if File is set)
    Dir string
    
    // Skip line comments in generated code for cleaner output
    SkipLineComments bool
    
    // File extension for template files (ignored if File is set)
    Ext string
    
    // Single file to compile (takes precedence over Dir/Ext)
    File string
}
```

### Configuration Options

- **Dir**: Directory containing template files. Defaults to current directory if empty.
- **SkipLineComments**: When `true`, generates cleaner code without line comments. Recommended for production.
- **Ext**: File extension filter for template files. Defaults to `.qtpl` if empty.
- **File**: Single file to compile. When specified, `Dir` and `Ext` are ignored.

## API Reference

### Core Functions

#### `QtcWrap()`
Compiles templates with default configuration (current directory, skip line comments).

#### `WithConfig(config Config)`
Compiles templates with custom configuration.

#### `CompileWithValidation(config Config) error`
Compiles templates with configuration validation. Returns error if validation fails.

### Convenience Functions

#### `CompileDirectory(dir string)`
Compiles all .qtpl files in the specified directory.

#### `CompileFile(file string)`
Compiles a single template file.

#### `CompileWithExtension(dir, ext string)`
Compiles template files with a specific extension in the given directory.

### Utility Functions

#### `GetDefaultConfig() Config`
Returns a Config struct with sensible defaults.

#### `ValidateConfig(config Config) error`
Validates configuration without compilation.

#### `IsQtcAvailable() bool`
Checks if qtc tool is available in the system.

#### `GetQtcVersion() (string, error)`
Returns the version of the qtc tool.

#### `FindTemplateFiles(dir, ext string) ([]string, error)`
Discovers template files in a directory (useful for preprocessing).

## Usage Examples

### Example 1: Basic Template Compilation

```go
package main

import (
    "fmt"
    "github.com/valksor/go-qtcwrap"
)

func main() {
    // Check if qtc is available
    if !qtcwrap.IsQtcAvailable() {
        fmt.Println("qtc tool not found. Please install it first.")
        return
    }
    
    // Compile templates in the templates directory
    qtcwrap.CompileDirectory("templates")
    
    fmt.Println("Template compilation completed!")
}
```

### Example 2: Advanced Configuration

```go
package main

import (
    "fmt"
    "github.com/valksor/go-qtcwrap"
)

func main() {
    config := qtcwrap.Config{
        Dir:              "src/templates",
        SkipLineComments: true,
        Ext:              ".qtpl",
    }
    
    // Validate configuration
    if err := qtcwrap.ValidateConfig(config); err != nil {
        fmt.Printf("Invalid configuration: %v\n", err)
        return
    }
    
    // Compile with validation
    if err := qtcwrap.CompileWithValidation(config); err != nil {
        fmt.Printf("Compilation failed: %v\n", err)
        return
    }
    
    fmt.Println("Templates compiled successfully!")
}
```

### Example 3: Template Discovery and Preprocessing

```go
package main

import (
    "fmt"
    "github.com/valksor/go-qtcwrap"
)

func main() {
    // Find all template files
    files, err := qtcwrap.FindTemplateFiles("templates", ".qtpl")
    if err != nil {
        fmt.Printf("Error finding templates: %v\n", err)
        return
    }
    
    fmt.Printf("Found %d template files:\n", len(files))
    for _, file := range files {
        fmt.Printf("  - %s\n", file)
    }
    
    // Compile each file individually
    for _, file := range files {
        fmt.Printf("Compiling %s...\n", file)
        qtcwrap.CompileFile(file)
    }
}
```

### Example 4: Integration with Build Systems

```go
package main

import (
    "fmt"
    "os"
    "github.com/valksor/go-qtcwrap"
)

func main() {
    // Get template directory from environment or use default
    templateDir := os.Getenv("TEMPLATE_DIR")
    if templateDir == "" {
        templateDir = "templates"
    }
    
    // Get qtc version for logging
    version, err := qtcwrap.GetQtcVersion()
    if err != nil {
        fmt.Printf("Warning: Cannot determine qtc version: %v\n", err)
    } else {
        fmt.Printf("Using qtc version: %s\n", version)
    }
    
    // Create build configuration
    config := qtcwrap.Config{
        Dir:              templateDir,
        SkipLineComments: true,  // Cleaner output for production
        Ext:              ".qtpl",
    }
    
    // Compile with error handling
    if err := qtcwrap.CompileWithValidation(config); err != nil {
        fmt.Printf("Build failed: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Println("Build completed successfully!")
}
```

## Error Handling

The package provides intelligent error handling:

- **Validation Errors**: Configuration validation errors are returned with descriptive messages
- **Tool Errors**: Missing qtc tool errors are handled gracefully
- **Compilation Errors**: Actual template compilation errors are displayed
- **Warning Suppression**: Common temporary file warnings are suppressed with informative messages

### Common Error Scenarios

1. **qtc tool not found**: Install qtc using `go install github.com/valyala/quicktemplate/qtc@latest`
2. **Invalid directory**: Ensure the specified directory exists and is accessible
3. **Invalid file**: Ensure the specified file exists and is readable
4. **Permission errors**: Ensure proper file/directory permissions
5. **Template syntax errors**: Fix syntax errors in your .qtpl files

## Integration with Build Systems

### Makefile Integration

Create a simple build script `scripts/compile-templates.go`:
```go
package main

import "github.com/valksor/go-qtcwrap"

func main() {
    qtcwrap.CompileDirectory("templates")
}
```

Then use it in your Makefile:
```makefile
.PHONY: templates
templates:
	@echo "Compiling templates..."
	@go run scripts/compile-templates.go

.PHONY: build
build: templates
	@echo "Building application..."
	@go build -o bin/myapp ./cmd/myapp
```

### Go Generate Integration

If you have a qtcwrap.go file that calls qtcwrap functions, you can run it directly:
```go
//go:generate go run ../qtcwrap/qtcwrap.go
package main
```

Or create a template compilation script (e.g., `scripts/gen-templates.go`):
```go
package main

import "github.com/valksor/go-qtcwrap"

func main() {
    qtcwrap.CompileDirectory("templates")
}
```

Then use it with go:generate:
```go
//go:generate go run scripts/gen-templates.go
package main
```

### CI/CD Integration

Create a template compilation script in your repository and use it in CI/CD:

`scripts/compile-templates.go`:
```go
package main

import "github.com/valksor/go-qtcwrap"

func main() {
    qtcwrap.CompileDirectory("templates")
}
```

GitHub Actions example:
```yaml
- name: Compile Templates
  run: |
    go install github.com/valyala/quicktemplate/qtc@latest
    go run scripts/compile-templates.go
```

## Best Practices

1. **Use SkipLineComments for Production**: Set `SkipLineComments: true` for cleaner generated code
2. **Validate Configuration**: Always validate configuration in production builds
3. **Check Tool Availability**: Verify qtc is available before compilation
4. **Error Handling**: Always handle errors from compilation functions
5. **Directory Organization**: Keep templates in a dedicated directory structure
6. **Version Pinning**: Pin qtc version in your CI/CD for reproducible builds

## Performance Considerations

- **Batch Compilation**: Use directory-based compilation for better performance with many files
- **Parallel Processing**: The package handles qtc execution efficiently
- **Memory Usage**: Large template sets may require adequate memory allocation
- **File System**: SSD storage recommended for large template compilations

## Compatibility

- **Go Version**: Requires Go 1.24 or later
- **qtc Version**: Compatible with all recent versions of quicktemplate/qtc
- **Operating Systems**: Works on all platforms supported by Go and qtc
- **Architectures**: Supports all Go-supported architectures

## Contributing

This package is part of the Valksor project. Contributions are welcome through:

1. Issue reporting for bugs and feature requests
2. Pull requests with improvements and fixes
3. Documentation improvements
4. Test coverage enhancements

## License

BSD 3-Clause License - see [LICENSE](LICENSE) file for details.

## Related Projects

- **QuickTemplate**: The underlying template engine - https://github.com/valyala/quicktemplate
- **qtc**: The QuickTemplate compiler - https://github.com/valyala/quicktemplate/tree/master/qtc

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for version history and changes.
