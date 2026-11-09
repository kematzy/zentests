# Zentests - Docs - Database

Use a real GORM database in tests without repeating setup and record boilerplate. The built-in
constructors open in-memory SQLite, optionally migrate models, and register cleanup; the record,
mutation, and transaction helpers operate on a supplied `*gorm.DB`.

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

- **GORM.** Database helpers use GORM models and connections. Setup functions create the
  `*gorm.DB`; record, mutation, and transaction helpers operate on the connection they receive.
- **SQLite in memory.** `SetupTestDB` and `SetupTestDBWithModels` use `gorm.io/driver/sqlite`
  with `":memory:"`. Other drivers and DSNs require opening your own `*gorm.DB`.
- **CGO.** That driver links `mattn/go-sqlite3`, so the database tests need `CGO_ENABLED=1` and a C
  compiler. The HTTP half of `zentests` does not.
- **A `*testing.T`.** Helpers take `t` first and report database errors as test failures. Most do
  not return errors; `DBTx` returns the callback error after rollback and fails the test on begin or
  commit errors.

## Cleanup is automatic

`SetupTestDB` and `SetupTestDBWithModels` register `t.Cleanup`, so the connection closes when the
test — or the whole suite, when created in `SetupSuite` — finishes.

Use `CloseTestDB` only when a test is _about_ closing the connection.

## Related

- [Testing styles](../testing/testify-suite.md) — one DB shared across a suite
- [examples/database test](../../examples/example-03.md) — a suite with a seeded database

---

Back to: [Docs](../README.md) &mdash; Next: [Setup](setup.md)
