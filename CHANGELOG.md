# Changelog

All notable changes to **tm2hsl** will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial project structure and documentation
- GitHub Actions CI/CD pipeline
- Makefile with common development tasks
- Development setup script
- Comprehensive README and CONTRIBUTING guides
- Conventional commits configuration
- Go project restructuring following industry standards

### Changed
- Restructured codebase to follow Go best practices
- Updated import paths and package organization
- Improved error handling patterns

### Technical
- Added Git hooks and commit message templates
- Configured linting and formatting tools
- Set up automated testing and coverage reporting
- Created release automation workflow

## [0.1.0] - 2024-01-XX

### Added
- Basic TextMate grammar parsing (JSON/plist formats)
- Initial IR (Intermediate Representation) structure
- CLI framework with Kong
- Project skeleton with all major components
- Basic HSL bytecode format definition

### Technical
- Go 1.21+ requirement
- Module structure with proper internal/pkg separation
- Initial testing framework setup

---

## Types of Changes

- `Added` for new features
- `Changed` for changes in existing functionality
- `Deprecated` for soon-to-be removed features
- `Removed` for now removed features
- `Fixed` for any bug fixes
- `Security` in case of vulnerabilities
- `Technical` for internal changes, refactoring, tooling

## Versioning

This project follows [Semantic Versioning](https://semver.org/):

- **MAJOR** version for incompatible API changes
- **MINOR** version for backwards-compatible functionality additions
- **PATCH** version for backwards-compatible bug fixes

## Release Process

1. Update version in relevant files
2. Update CHANGELOG.md
3. Create git tag
4. GitHub Actions automatically creates release
5. Update documentation if needed