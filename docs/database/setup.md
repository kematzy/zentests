# Zentests - Docs - Database - Setup

> **Status:** DRAFT — structure and wording are open for discussion.

Five functions cover the lifecycle of a test database: open, migrate, reset, close.

| Function                              | Purpose                                                         |
| ------------------------------------- | --------------------------------------------------------------- |
| `SetupTestDB(t)`                      | Open an empty in-memory SQLite database, closed via `t.Cleanup` |
| `SetupTestDBWithModels(t, models...)` | The same, plus `AutoMigrate` for each model                     |
| `DBMigrate(t, db, models...)`         | Migrate more models onto an existing connection                 |
| `DBReset(t, db)`                      | Delete every row, keep the schema, reset SQLite's counter       |
| `CloseTestDB(t, db)`                  | Close early, when the test is about closing                     |

## Opening

```go
// no schema yet
db := zentests.SetupTestDB(t)
zentests.DBMigrate(t, db, &User{}, &Post{})

// the common case: open and migrate in one line
db := zentests.SetupTestDBWithModels(t, &User{}, &Post{})
```

Both use `":memory:"` with GORM's silent logger, so a failing query produces a test failure rather
than log noise, and `time.Now().Local()` for timestamps.

## Resetting between tests

`DBReset` deletes all rows from every migrated table and clears `sqlite_sequence`, so auto-increment
IDs start at 1 again. Schema and indexes stay, which makes it cheap enough to run before every test:

```go
func (s *UserSuite) SetupTest() {
	zentests.DBReset(s.T(), s.db)
}
```

In a flat test file, the equivalent is a fresh `SetupTestDBWithModels` per test — the in-memory
database is gone when the test ends anyway.

## Closing early

`SetupTestDB` registers the close through `t.Cleanup`, so the usual test needs nothing.
`CloseTestDB` exists for tests that must observe a closed connection:

```go
zentests.CloseTestDB(t, db)

zt.Get(app, "/api/v1/users").ServerError() // the handler now hits a dead connection
```

Calling it twice is harmless: the second close returns an error that the helper ignores.

## Notes

- **`DBReset` needs at least one migration.** With no tables it is a no-op.
- **In-memory means private.** Two `SetupTestDB` calls give two unrelated databases, even in the
  same test.
- **Shared cache is not enabled** — no `?cache=shared` — so parallel tests each get their own
  database instead of fighting over one.

## Related

- [Drivers](drivers.md) — why SQLite, and the CGO requirement
- [Records](records.md) — creating data once the schema exists
- [Testing styles](../testing/testify-suite.md) — suite-level setup and teardown

---

Back to: [Database](./README.md) &mdash; Next: [Drivers](drivers.md)
