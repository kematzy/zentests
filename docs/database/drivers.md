# Zentests - Docs - Database - Drivers

The built-in database constructors use GORM's SQLite driver with an in-memory database. Record,
mutation, and transaction helpers accept a `*gorm.DB`, so they can also be used with a connection
opened by your application.

## SQLite, via GORM

```go
db := zentests.SetupTestDBWithModels(t, &User{}, &Post{})
// equivalent to:
// gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: silent, NowFunc: local time})
```

- **Driver:** `gorm.io/driver/sqlite` — GORM's recommended SQLite driver.
- **DSN:** `":memory:"`, so the database does not touch the filesystem. Each constructor call opens
  a fresh database; reusing a handle preserves its data until reset.
- **Logger:** GORM logging is silent; database helpers report errors through `testing.T`.
- **Time:** `time.Now().Local()`, so timestamp assertions are not surprised by UTC conversion.

## CGO is required

`gorm.io/driver/sqlite` uses `mattn/go-sqlite3`, which requires CGO. Running tests that use the
SQLite constructors requires:

- `CGO_ENABLED=1` (the default when a C compiler is present), and
- a C toolchain — `gcc`/`clang` on Linux and macOS (Xcode command line tools count), MSYS2 or MinGW
  on Windows.

If your CI image is `golang:alpine` or uses `CGO_ENABLED=0`, install a C toolchain and enable CGO,
or omit database-helper builds. The HTTP assertions compile with CGO disabled.

## Built-in constructor limits

- **PostgreSQL, MySQL, or another driver in the constructors.** `SetupTestDB` and
  `SetupTestDBWithModels` always use in-memory SQLite and accept no dialector or DSN. Open another
  GORM connection yourself, then pass it to the record, mutation, migration, or transaction helpers.
  `DBReset` is intended for SQLite; its sequence reset is SQLite-specific.
- **File-backed SQLite through the constructors.** They hardwire `":memory:"`; open a file-backed
  connection yourself if you need one.
- **Pure-Go SQLite in the built-in constructors** (`modernc.org/sqlite`, `glebarez/sqlite`). These
  drivers are not wired in; the constructors use `mattn/go-sqlite3` and require CGO.

## Using your own connection

Record, mutation, migration, and transaction helpers take `db *gorm.DB`; the constructors are
convenience functions for in-memory SQLite, not a requirement. Open and migrate your own connection:

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
