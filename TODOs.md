# TODOs

## Confirmed

Any items within this section have been confirmed as to be implemented.

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
update README + `docs/database/*` → record the breaking change in the changelog. If removing the
`DB` prefix, rename the symbols as shown in the table. Avoid root-package aliases: wrappers importing
`dbtest` would restore the SQLite compile dependency.

**Caveat.** A package split does not remove the CGO requirement for users of the SQLite helpers.
Keep the tested `gorm.io/driver/sqlite` unless database users have a separate need for a pure-Go
driver; evaluate alternatives such as `glebarez/sqlite` over `modernc.org/sqlite` for compatibility
and maintenance before switching.
