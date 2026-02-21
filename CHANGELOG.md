# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [unreleased]

### 🚀 Features

- Add modernize & modernize-check commands
- Add versions command
- Add changelog command
- Add git-cliff configs (git-cliff --init)
- Add release command
- Add git-push (github-push, gitea-push) & docs, docs-md, clean commands

### 🚜 Refactor

- Reworked Makefile
- Update .PHONY list

### ⚙️ Miscellaneous Tasks

- Remove bad comments
- Updated README.md with Makefile commands
- Go mod tidy

## [0.1.0] - 2026-02-19

### Added
- New content type assertion methods: IsCSS(), IsJS(), IsXML(), IsXHR(), IsImage(), IsPNG(), IsJPEG(), IsGIF(), IsSVG(), IsWebP()
- Comprehensive documentation with Parameters and Returns sections for all functions
- New examples in `examples/` directory: flat_style.go, routes.go, suite.go
- Test coverage improved from 94.7% to 98.2%
- GitHub Actions CI configuration with testing, linting, and coverage artifacts
- Gitea Actions CI configuration with Docker-optimized workflows
- Makefile with release automation

### Changed
- Refactored error handling in HasInt() and HasFloat() to use assert.Fail() instead of t.Errorf()
- Made getNestedValue() directly testable for better coverage
- Added Implementation: headers before function definitions for better readability

### Fixed
- Fixed GitHub Actions artifact upload (Gitea-compatible tarball approach)
- Fixed Gitea CI Docker container configuration

### 💼 Other

- Improved Gitea workflow
- Added Gitea workflow configs
