# TODOs

## Confirmed

Any items within this section have been confirmed as to be implemented.

### Report failures at the caller's line (`t.Helper()`)

Assertions and request helpers report the library's file:line in the failure headline today, e.g.
`response_http.go:24` for `.OK()`, because only a dozen of the ~105 exported methods call
`t.Helper()`. Add it to every exported method so `go test`, editors and CI annotations point at the
line in the test file.

Not a correctness problem — the test still fails with the right message and `testify`'s trace
contains the caller — but jump-to-failure should land on the test, not on the library.

| File | Methods without `t.Helper()` | Note |
| ---- | ----------------------------- | ---- |
| `response_http.go` | 41 (~32 missing) | `OK()` delegates to `Status()`; both need it |
| `request.go` | 25 | `Get()` delegates to `GetWithConfig()`; both need it |
| `response_json.go` | 18 (~15 missing) | |
| `response_body.go` | 5 (all missing) | uses `assert.*` directly |
| `response.go` | 5 (all missing) | includes `Debug()` and `Dump()` |
| `jar.go` | 8 (all missing) | `j.t.Helper()` — `testing.TB` has `Helper()` |
| `suite.go` | 5 | `Describe`/`It`/hooks |
| `fiber.go` | 2 | already uses `zt.Helper()` |
| `zentests.go` | 2 | `New()`, `Use()` |

Not applicable to `logger.go` — `DefaultLogger`/`MockLogger` have no `*testing.T`.

Acceptance: a failing chain (`zt.Get(app, "/nope").OK()`) prints the test file:line in the failure
headline, and the existing test suite still passes unchanged.

## Ideas

Any items within this section are ideas and thoughts kept for future consideration, and NOT be
implemented.

### Split DB functionality into `dbtest` package

**Status:** on hold — the current single-package setup stays as it is.

**Why consider it at all.** Two separate costs are currently paid by every user of the library:

1. `db_gorm.go` imports `gorm.io/driver/sqlite`, which is the CGO binding (`mattn/go-sqlite3`), so
   anyone running tests for a project that only uses the HTTP helpers still needs `CGO_ENABLED=1`
   and a C toolchain to compile a package they never call.
2. The GORM + SQLite dependencies sit in the root package, so they also land in the dependency
   graph of projects that only want assertions.

Go only compiles what is imported, so moving the helpers into their own package removes both costs
without touching the API surface:

```
github.com/kematzy/zentests          T, Response, Jar, Describe/It      deps: fiber, testify
github.com/kematzy/zentests/dbtest   Setup, Create, Find, Tx, Seed …    deps: gorm, sqlite driver
```

The `/dbtest` suffix follows the standard library convention (`httptest`, `fstest`, `iotest`).

**Symbols would lose the redundant `DB` prefix** — `dbtest.Setup`, `dbtest.Create`, `dbtest.Find`
read better than `dbtest.DBSetup`, and the import path already says "database".

| Today (`zentests.`) | After (`dbtest.`) |
| ------------------- | ----------------- |
| `SetupTestDB` / `SetupTestDBWithModels` | `Setup` / `SetupWithModels` |
| `DBMigrate` / `DBReset` / `CloseTestDB` | `Migrate` / `Reset` / `Close` |
| `DBCreate` / `DBCreateN` / `DBSeed` / `DBSeedSlice` | `Create` / `CreateN` / `Seed` / `SeedSlice` |
| `DBFind` / `DBFindBy` / `DBFirst` / `DBLast` | `Find` / `FindBy` / `First` / `Last` |
| `DBCount` / `DBExists` / `DBNotExists` | `Count` / `Exists` / `NotExists` |
| `DBUpdate` / `DBUpdateBy` / `DBDelete` / `DBDeleteBy` / `DBTx` | `Update` / `UpdateBy` / `Delete` / `DeleteBy` / `Tx` |

**Sub-package first, sub-module only if ever needed.** `github.com/kematzy/zentests/dbtest` is the
same import path either way, so it can start as a package inside this module (one `go.mod`, one
release, no CI split) and be extracted into its own module later with a nested tag (`dbtest/v0.x.0`)
and zero API change. Doing it the other way round costs the release tooling immediately and gains
users nothing visible.

**Migration outline.** Create `dbtest/` and `git mv db_gorm.go` → rename symbols → move `db_gorm_test.go`
with it → update README + `docs/database/*` → CHANGELOG entry marked breaking (v0.4.0) with the
table above. No compatibility aliases in the root package: an alias would re-add the GORM import
and undo the whole point.

**Caveat.** Keeping a default in-memory SQLite is what makes the helpers one-liners, so the driver
decision has to be made together with the split: either stay on `gorm.io/driver/sqlite` (CGO, but
"GORM's recommended" and the only driver this project actually tests against), or switch to a
pure-Go driver (`glebarez/sqlite` over `modernc.org/sqlite`) whose provenance and maintenance
activity need checking first.

### Verify the Markdown examples in CI

The runnable examples live in `examples/*.md` and no pipeline compiles them. `go vet` cannot see
Markdown, and the old `//go:build manual` Go examples were removed when the pages became the single
source, so `README.md` can keep claiming "the files above are the ones that run" while an API change
breaks them silently.

Proposed `make examples-check` target:

1. Extract every ```` ```go ```` block from `examples/*.md` into a temp directory, one package per
   file, keeping the title as the package name.
2. Write a minimal `go.mod` whose module `replace`s `github.com/kematzy/zentests` with the local
   checkout, and `go mod tidy` it.
3. Run `go vet ./...` and `go test ./...` in that directory.
4. Exit non-zero on any failure, so `make check` (and CI) catches a stale example.

Cost: ~30-40 lines of shell or Python, no new Go dependencies. Caveat: the suite example uses the
GORM SQLite driver, so the target needs `CGO_ENABLED=1` and a C compiler — both CI images already
install `libsqlite3-dev`.
