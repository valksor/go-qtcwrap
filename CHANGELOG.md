# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of qtcwrap package
- Core functionality for QuickTemplate compilation via `QtcWrap()` and `WithConfig()`
- Comprehensive configuration system with `Config` struct
- Support for directory-based template compilation
- Support for single-file template compilation
- Intelligent error handling and warning suppression for temporary file issues
- Configuration validation with `ValidateConfig()` and `CompileWithValidation()`
- Template file discovery with `FindTemplateFiles()`
- Tool availability checking with `IsQtcAvailable()` and `GetQtcVersion()`
- Convenience functions: `CompileDirectory()`, `CompileFile()`, `CompileWithExtension()`
- Default configuration helper with `GetDefaultConfig()`
- Flexible file extension filtering support
- Skip line comments option for cleaner generated code
- Proper command-line argument building for qtc tool
- Security-conscious command execution with validated arguments
- Complete package documentation with usage examples
- Extensive test coverage (25+ test functions, 100+ sub-tests)

### Configuration Features
- `Dir`: Directory-based template compilation
- `SkipLineComments`: Toggle for cleaner generated code output
- `Ext`: Custom file extension filtering (defaults to .qtpl)
- `File`: Single file compilation mode (overrides Dir/Ext)

### Error Handling
- Graceful handling of missing qtc tool
- Intelligent temporary file warning suppression
- Detailed validation error messages
- Proper error propagation for compilation failures

### Convenience Functions
- `QtcWrap()`: Default compilation with sensible defaults
- `CompileDirectory(dir)`: Directory-focused compilation
- `CompileFile(file)`: Single file compilation
- `CompileWithExtension(dir, ext)`: Extension-filtered compilation
- `CompileWithValidation(config)`: Validated compilation

### Utility Functions
- `GetDefaultConfig()`: Sensible default configuration
- `ValidateConfig(config)`: Configuration validation
- `IsQtcAvailable()`: Tool availability check
- `GetQtcVersion()`: Version information retrieval
- `FindTemplateFiles(dir, ext)`: Template discovery

### Technical Details
- Module name: `github.com/valksor/go-qtcwrap`
- Go version: 1.24+
- Package name: `qtcwrap`
- Zero external dependencies (except qtc tool requirement)
- Cross-platform compatibility
- Proper file system permission handling
- Concurrent operation support
- Symlink support for directories and files
- Special character handling in file paths

### Documentation
- Complete API documentation with examples
- README with comprehensive usage guide and best practices
- Integration examples for build systems (Makefile, Go Generate, CI/CD)
- Performance considerations and optimization tips
- Troubleshooting guide for common issues
- Extensive test suite demonstrating all functionality

### Security
- Command injection prevention with validated arguments
- Proper error message sanitization
- Safe file system operations
- Secure temporary file handling
