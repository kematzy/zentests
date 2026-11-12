## [0.4.0] - 2026-11-12

### 🚀 Features

- *(makefile)* Add test-dox and test-watch targets
- *(jar)* Add cookie-keeping Jar for multi-request tests
- *(response)* Add redirect and status-class assertions
- *(response)* Add HasCookie and CookieNotPresent assertions
- *(response)* Add ArrayContains, ArrayEmpty and Bind
- *(request,jar)* Add default headers on T and Jar
- *(request,jar)* Add PostMultipart for file uploads
- *(request)* Make the default timeout configurable with ZENTESTS_TIMEOUT
- *(request)* Add GetQuery for encoded query parameters
- Add DB/GORM test helpers

### 🐛 Bug Fixes

- *(request)* Stop test on transport error instead of panicking
- *(request)* Url-encode form bodies in PostForm
- *(suite)* Run AfterEach via t.Cleanup and snapshot hooks at It
- *(json)* Tighten JSON assertions and document subset semantics
- *(cookies)* Report cookies that were never set
- *(debug)* Redact sensitive headers and share formatting with Dump
- *(zentests)* Fail fast when Use is called with a nil app
- *(body)* Stop BodyMatches after an invalid regex
- Satisfy godoclint and noctx so make lint passes
- *(response)* Redact X-Api-Key and X-Auth-Token in Dump and Debug
- *(suite)* Run AfterEach hooks even when a BeforeEach hook skips
- *(response)* Redact more token headers in Debug/Dump
- *(testing)* Report assertion failures at caller lines
- *(release)* Generate changelog from untagged commits

### 💼 Other

- Fix coverage target and align CI with supported Go versions
- *(makefile)* Document target flags and fix the ignored check target
- *(makefile)* Generate coverage HTML with go-better-html-coverage
- *(makefile)* Run the suite once in check and restore coverage-profile phony
- *(makefile)* Fail fmt-check on unformatted files and harden check
- Add AGENTS.md
- *(makefile)* Write the coverage report into .code-status
- *(makefile)* Drop the personal git-push helpers
- *(makefile)* Commit the changelog and push on release

### 🚜 Refactor

- *(response)* Accept testing.TB for internal assertion target
- *(request)* Route all requests through execute and executeJSON
- *(request)* Split request execution out of execute
- Fix make release task
- *(db_gorm)* Centralise condition handling and validate the where string

### 📚 Documentation

- Update README.md with 'testify/suite' examples
- Regenerate API reference and modernize examples
- *(readme)* Document numeric Has, zero-value TestConfig, hooks and JSONHasValues
- *(api)* Regenerate API reference with gomarkdoc
- Correct Status and JSONHasValues examples and coverage figure
- *(response)* Deprecate JSONMatches in favour of JSONHasValues
- Refresh the coverage figure to 97.8%
- *(api)* Regenerate API reference with the new APIs
- Move the testify/suite tutorial to docs/SUITE.md
- Trim the Features section to a cheat sheet
- Move the Fiber v2 migration notes to docs/MIGRATION-FIBER-V2.md
- Move branch policy and dev workflow to CONTRIBUTING.md
- Add contents bar, trim Examples to a list, make coverage claims drift-proof
- Drop the Fiber v2 migration guide
- Replace the license text with a link and shorten credits
- *(api)* Deprecate Use and correct the request examples
- Update README.md
- *(api)* Regenerate the API reference
- *(changelog)* Fix the DefaultLogger typo
- *(license)* Use 2026 as the copyright year
- *(examples)* Replace the Go examples with runnable Markdown pages
- *(assertions)* Add the assertion reference
- *(requests)* Add the request reference
- *(database)* Add the database reference
- *(helpers)* Add the helper reference
- *(testing)* Add the testing style guides
- Add the docs index and the comparison page
- *(contributing)* Move the contributing guide into docs
- *(readme)* Rewrite the README as a front door
- *(readme)* Clarify package scope and value
- *(assertions)* Clarify failure and matching behavior
- *(requests)* Document config and cookie semantics
- *(testing)* Correct lifecycle and isolation guidance
- *(database)* Clarify driver and transaction behavior
- *(helpers)* Document app cleanup and logger behavior
- *(guide)* Clarify navigation and comparison

### 🎨 Styling

- Use the American spelling of canonicalize

### 🧪 Testing

- Cover assertion failure paths
- *(json)* Cover numeric normalization and regex error branches
- Additional DB/GORM helpers tests
- Document the test suite types
- *(jar)* Create the jar in SetupTest
- *(examples)* Verify Markdown examples in CI

### ⚙️ Miscellaneous Tasks

- Vet examples under the manual build tag
- Test the go.mod minimum and lint with a Go 1.27-compatible golangci-lint
- *(cspell)* Skip generated files and register tool names
- Align editor and prettier config for json and markdown
- *(gitignore)* Ignore the local TMP scratch directory
- *(editor)* Add Zed settings for JSON and YAML formatting
- *(cspell)* Register chmouel as a known word
- Regenerate the coverage HTML report
- Regenerate the coverage report
- *(spellcheck)* Extend the dictionary for the new docs
- *(tooling)* Pin node for cspell and format Markdown on save
- Record the deferred work in TODOs.md
- Updated TODOs
- *(cspell)* Allow historical changelog spelling
## [0.3.0] - 2026-04-05

### 🚀 Features

- Add Logger, MockLogger & DefaultLogger
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
