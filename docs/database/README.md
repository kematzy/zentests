# Zentests - Docs - Database

> **Status:** DRAFT — structure and wording are open for discussion.

A real database in the test, gone when the test ends: the helpers open in-memory SQLite, migrate
your GORM models and close the connection themselves, then give you short functions for the create,
read, update and delete that every test repeats.

## Contents

| Page                            | What it answers                                        |
| ------------------------------- | ------------------------------------------------------ |
| [Setup](setup.md)               | Opening a test DB, migrating, resetting, closing early |
| [Drivers](drivers.md)           | SQLite, CGO, and what is not supported                 |
| [Records](records.md)           | Create, seed, find, count, exists                      |
| [Mutations](mutations.md)       | Update and delete by ID or condition                   |
| [Transactions](transactions.md) | `DBTx`, and when a transaction belongs in a test       |

## The shape

```go
db := zentests.SetupTestDBWithModels(t, &User{}, &Post{}) // :memory:, AutoMigrate, auto-close

user := zentests.DBCreate(t, db, &User{Name: "Jane", Email: "jane@example.com"})

found := zentests.DBFind(t, db, &User{}, user.ID)
zentests.DBUpdate(t, db, found, found.ID, map[string]any{"name": "Janet"})
zentests.DBCount(t, db, &User{}) // 1
```

## What the helpers assume

- **GORM.** The helpers take a `*gorm.DB`, so they work with GORM models and GORM's query builder.
- **SQLite in memory.** `SetupTestDB` and `SetupTestDBWithModels` hardwire `gorm.io/driver/sqlite`
  with `":memory:"`. Other databases are not supported by these constructors.
- **CGO.** That driver links `mattn/go-sqlite3`, so the database tests need `CGO_ENABLED=1` and a C
  compiler. The HTTP half of `zentests` does not.
- **A `*testing.T`.** Every function takes `t` as its first argument, calls `t.Helper()` and fails
  the test on error — there are no error returns to check.

## Cleanup is automatic

`SetupTestDB` and `SetupTestDBWithModels` register `t.Cleanup`, so the connection closes when the
test — or the whole suite, when created in `SetupSuite` — finishes.

Use `CloseTestDB` only when a test is _about_ closing the connection.

## Related

- [Testing styles](../testing/testify-suite.md) — one DB shared across a suite
- [examples/database test](../../examples/example-03.md) — a suite with a seeded database

---

Back to: [Docs](../README.md) &mdash; Next: [Setup](setup.md)
