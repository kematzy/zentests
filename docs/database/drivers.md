# Zentests - Docs - Database - Drivers

> **Status:** DRAFT — structure and wording are open for discussion.

The database helpers are built around one driver, on purpose: GORM's SQLite driver, in memory. That
is what makes them one-liners, and it is the only database configuration this project tests.

## SQLite, via GORM

```go
db := zentests.SetupTestDBWithModels(t, &User{}, &Post{})
// equivalent to:
// gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: silent, NowFunc: local time})
```

- **Driver:** `gorm.io/driver/sqlite` — GORM's recommended SQLite driver.
- **DSN:** `":memory:"`, so nothing touches the filesystem and every test starts empty.
- **Logger:** silent, so a failing query reports through `testing` rather than the console.
- **Time:** `time.Now().Local()`, so timestamp assertions are not surprised by UTC conversion.

## CGO is required

`gorm.io/driver/sqlite` links `mattn/go-sqlite3`, which is a C library. Building the database
helpers therefore needs:

- `CGO_ENABLED=1` (the default when a C compiler is present), and
- a C toolchain — `gcc`/`clang` on Linux and macOS (Xcode command line tools count), MSYS2 or MinGW
  on Windows.

If your CI image is `golang:alpine` or `CGO_ENABLED=0`, either install a toolchain and `gcc`, or
keep the database helpers out of that build — the HTTP assertions compile without any of it.

## What is not supported

- **PostgreSQL, MySQL, MariaDB.** No constructor accepts a foreign dialector, and no DSN parameter
  exists. You can use those databases in your own tests, but not through these helpers.
- **File-backed SQLite.** The helpers hardwire `":memory:"`. A `DBFind` on a `*gorm.DB` you opened
  yourself works fine, since every helper takes the `*gorm.DB` as an argument.
- **Pure-Go SQLite** (`modernc.org/sqlite`, `glebarez/sqlite`). Not wired in; the CGO caveat above
  applies.

## Using your own connection

Every helper takes `db *gorm.DB`, so the constructors are a convenience, not a requirement:

```go
db, err := gorm.Open(sqlite.Open("file:test.db"), &gorm.Config{})
require.NoError(t, err)

zentests.DBMigrate(t, db, &User{})
zentests.DBCreate(t, db, &User{Name: "Jane"})
```

That is the seam for anything the defaults do not cover: a different DSN, a different driver your
project already uses, or a database shared across packages.

## Related

- [Setup](setup.md) — what the constructors do
- [Requirements](../../README.md#requirements) — the CGO note in the README

---

Back to: [Database](./README.md) &mdash; Next: [Records](records.md)
