## [0.3.0] - 2026-04-05

### 🚀 Features

- Add Logger, MockLogger & DefaultLoggger
- Add Fiber helpers HasRoute(), NotHasRoute(), NewApp()
- Add DB GORM test helpers

### 🐛 Bug Fixes

- Updated 'make release' command
- Fix typos & add test-watch
- Removed TODO
- Typo

### 🚜 Refactor

- Security improvements

### 📚 Documentation

- Additions to README.md file

### 🧪 Testing

- Logger, MockLogger & DefaultLogger tests
- Fiber helpers HasRoute(), NotHasRoute() & NewApp() tests
- DB GORM helpers tests

### ⚙️ Miscellaneous Tasks

- Go mod tidy (gorm packages)
- Updated code-coverage.html
- Spellchecker configs
## [0.2.0] - 2026-03-19

### 🚀 Features

- Add modernize & modernize-check commands
- Add versions command
- Add changelog command
- Add git-cliff configs (git-cliff --init)
- Add release command
- Add git-push (github-push, gitea-push) & docs, docs-md, clean commands
- Add Cspell configs
- Add *WithConfig()` functions

### 🐛 Bug Fixes

- Golangci configs
- Github CI failure on Lint job
- Reactivating Gitea Lint job
- GitHub CI failure on Lint job (try #2)

### 💼 Other

- Add required installs

### 🚜 Refactor

- Reworked Makefile
- Update .PHONY list
- Make modernize (code updates)
- Update to Fiber v3
- Update fiber imports
- Update handler signatures (from pointer to interface)
- Fiber v3 - BodyParser() -> Bind().Body()

### 📚 Documentation

- Add API.md (gomarkdoc trial)
- Updated README.md about Fiber v2 -> v3 migration & functionality
- Update README.md with 'testify/suite' examples

### 🧪 Testing

- Updated HTTP response tests
- Add code-coverage.html (temporary test)
- Fiber v3 JSON adds `charset=utf-8`
- Add tests for *WithConfig() functions
- Convert to testify/suite structure

### ⚙️ Miscellaneous Tasks

- Remove bad comments
- Updated README.md with Makefile commands
- Go mod tidy
- Add CHANGELOG.md
- Docs update
- Updated make test-verbose output
- Updated .gitignore
## [0.1.0] - 2026-02-19

### 💼 Other

- Improved Gitea workflow Maybe???
- Gitea workflow configs
- Gitea workflow again - ensure NodeJS is available
