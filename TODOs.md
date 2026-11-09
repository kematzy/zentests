# TODOs

## Confirmed

Any items within this section have been confirmed as to be implemented.

### Report failures at the caller's line (`t.Helper()`)

A failed assertion currently points at a library file:line in the failure report, such as
`response_http.go:24` for `.OK()`. Calling `t.Helper()` (or `r.t.Helper()` / `j.t.Helper()`) in the
failing method and each wrapper in its call chain should make the report point to the test call site.

This is not a correctness problem—the test still fails—but it makes failures harder to jump to in
editors and CI. Add helpers to exported methods that report test failures and to wrappers that lead
to those methods. Do not add them indiscriminately to every exported method: pure setters, accessors,
and constructors that cannot report a test failure do not need them. Logger methods have no testing
context.

**Audit before editing:** the old per-file counts were stale. For example, several `Jar` methods and
the route assertions already call `Helper()`, while `zentests.New` does not itself report failures.
Recheck the current call chains rather than relying on the previous count table.

**Acceptance:** a deliberately failing chain such as `zt.Get(app, "/nope").OK()` reports the test
file and the line containing the chain, not a `zentests` source line. Verify at least one direct
assertion and one delegated chain, then run the existing test suite unchanged.

---

### Verify the Markdown examples in CI

The runnable examples live in `examples/*.md`, but CI does not compile or run them. `go vet` cannot
see Markdown, so an API change can silently break examples that the README presents as runnable.

Proposed `make examples-check` target:

1. Extract the designated runnable `go` block from each example into a temporary module. Each
   current example has one complete Go block; fail clearly if that convention changes. Do not treat
   arbitrary future Go snippets as runnable examples.
2. Derive temporary directory and module names from the example filenames, and keep valid Go package
   names from the source. Do not use Markdown titles as package names.
3. Write a temporary `go.mod` requiring `github.com/kematzy/zentests` and replacing it with the local
   checkout; run `go mod tidy` there.
4. Run `go vet ./...` and `go test ./...`; return non-zero on any failure.
5. Run the target from `make check` and CI.

Cost: roughly 30–40 lines of shell or Python, with no new Go dependencies. Example 03 uses the
SQLite helpers, so CI must have `CGO_ENABLED=1` and a C compiler. Verify the compiler in both jobs;
the bundled `go-sqlite3` driver does not require a system SQLite development package by itself.

---

## Ideas

Any items within this section are ideas and thoughts kept for future consideration, and should NOT
be implemented.

---

## Archived

Any items within this section are old ideas kept for future reference ONLY. They should NOT be
implemented.

### Split DB functionality into `dbtest` package

**Status:** Archived. Revisit only if users report significant CGO/toolchain friction for
HTTP-only tests.

**Why it was considered.** The root package includes `db_gorm.go`, which imports GORM's SQLite
driver (`mattn/go-sqlite3`). With CGO enabled, an HTTP-only consumer on a machine without a C
compiler may still fail to build the root package; setting `CGO_ENABLED=0` is the current workaround
for HTTP-only use.

A `dbtest` subpackage would stop compiling the SQLite code when users import only
`github.com/kematzy/zentests`:

```txt
github.com/kematzy/zentests          T, Response, Jar, Describe/It      imports: fiber, testify
github.com/kematzy/zentests/dbtest   Setup, Create, Find, Tx, Seed …    imports: gorm, sqlite driver, testify
```

However, a subpackage in the same module does **not** remove GORM and SQLite from the root
`go.mod`; it may not remove those module requirements from a consumer's module graph. Verify that
benefit with a small downstream module. Full dependency separation may require a nested module.

The `/dbtest` suffix follows the standard library convention (`httptest`, `fstest`, `iotest`).

Removing the redundant `DB` prefix (`dbtest.Create` instead of `dbtest.DBCreate`) is optional and
would be an additional API rename, not a requirement of the package split.

| Today (`zentests.`)                                            | After (`dbtest.`)                                    |
| -------------------------------------------------------------- | ---------------------------------------------------- |
| `SetupTestDB` / `SetupTestDBWithModels`                        | `Setup` / `SetupWithModels`                          |
| `DBMigrate` / `DBReset` / `CloseTestDB`                        | `Migrate` / `Reset` / `Close`                        |
| `DBCreate` / `DBCreateN` / `DBSeed` / `DBSeedSlice`            | `Create` / `CreateN` / `Seed` / `SeedSlice`          |
| `DBFind` / `DBFindBy` / `DBFirst` / `DBLast`                   | `Find` / `FindBy` / `First` / `Last`                 |
| `DBCount` / `DBExists` / `DBNotExists`                         | `Count` / `Exists` / `NotExists`                     |
| `DBUpdate` / `DBUpdateBy` / `DBDelete` / `DBDeleteBy` / `DBTx` | `Update` / `UpdateBy` / `Delete` / `DeleteBy` / `Tx` |

**Package vs. module.** A subpackage in this module is the smaller API/code change, but keeps the
GORM/SQLite requirements in the root `go.mod`. A nested module can isolate those requirements but
adds release and CI complexity. Choose based on whether the goal is only to avoid compiling SQLite
for root-package users or also to remove it from their module graph.

**Migration outline.** Create `dbtest/` and `git mv db_gorm.go` → move `db_gorm_test.go` with it →
update README + `docs/database/*` → CHANGELOG entry marked breaking (v0.4.0). If removing the `DB`
prefix, rename the symbols as shown in the table. Avoid root-package aliases: wrappers importing
`dbtest` would restore the SQLite compile dependency.

**Caveat.** A package split does not remove the CGO requirement for users of the SQLite helpers.
Keep the tested `gorm.io/driver/sqlite` unless database users have a separate need for a pure-Go
driver; evaluate alternatives such as `glebarez/sqlite` over `modernc.org/sqlite` for compatibility
and maintenance before switching.
