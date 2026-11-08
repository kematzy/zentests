# Zentests - Docs - Database - Records

> **Status:** DRAFT — structure and wording are open for discussion.

Creating and reading records without repeating `db.Create(...)` plus an error check in every test.

## Creating

```go
user := zentests.DBCreate(t, db, &User{Name: "Jane", Email: "jane@example.com"})
// user is the same pointer, with ID and timestamps filled in

users := zentests.DBCreateN(t, db, 3, func(i int) User {
	return User{Email: fmt.Sprintf("user%d@example.com", i)} // i runs 1..3, not 0..2
})
```

`DBCreateN`'s factory is **1-based**, so the sequence reads `user1`, `user2`, `user3`. Returning a
value (not a pointer) keeps the factory free of loop-variable aliasing bugs.

## Seeding

For a fixed fixture, seed a slice instead of calling create in a loop:

```go
zentests.DBSeed(t, db, []*User{
	{Name: "Alice", Email: "alice@example.com"},
	{Name: "Bob", Email: "bob@example.com"},
})

zentests.DBSeedSlice(t, db, []User{
	{Name: "Carol", Email: "carol@example.com"},
	{Name: "Dave", Email: "dave@example.com"},
})
```

`DBSeed` takes pointers; `DBSeedSlice` takes values and creates a pointer internally. Both insert
row by row, so a constraint violation fails at the offending index.

## Reading

| Function                                     | Returns                          |
| -------------------------------------------- | -------------------------------- |
| `DBFind(t, db, &User{}, id)`                 | The record with that primary key |
| `DBFindBy(t, db, &User{}, where, args...)`   | The first match for a condition  |
| `DBFirst(t, db, &User{})`                    | The first record by primary key  |
| `DBLast(t, db, &User{})`                     | The last record by primary key   |
| `DBCount(t, db, &User{}, conditions...)`     | The number of matching rows      |
| `DBExists(t, db, &User{}, conditions...)`    | Whether at least one row matches |
| `DBNotExists(t, db, &User{}, conditions...)` | Whether no row matches           |

```go
user := zentests.DBFind(t, db, &User{}, 1)
jane := zentests.DBFindBy(t, db, &User{}, "email = ?", "jane@example.com")

count := zentests.DBCount(t, db, &User{})
active := zentests.DBCount(t, db, &User{}, "status = ?", "active")

zentests.DBExists(t, db, &User{}, "email = ?", "jane@example.com")
zentests.DBNotExists(t, db, &User{}, "status = ?", "deleted")
```

`DBFind`, `DBFindBy`, `DBFirst` and `DBLast` fail the test when nothing matches, and return a
pointer to the populated record otherwise. `DBCount`, `DBExists` and `DBNotExists` return values you
can assert on, and fail only on a database error.

## Conditions

`conditions...` is an optional `where` string followed by its arguments — the same shape as GORM's
`Where`. The first argument must be a string; anything else fails the test with a message naming the
offending type rather than panicking.

```go
zentests.DBCount(t, db, &User{}, "created_at > ? AND status = ?", cutoff, "active")
```

## Notes

- **The model argument is for the table.** `&User{}` needs no fields set; only its type matters.
- **`DBFind` allocates nothing for you** — pass a pointer to a zero value, as shown.
- **Ordering is the database's.** `DBFirst`/`DBLast` order by primary key; use `DBFindBy` with your
  own ordering when insertion order matters.

## Related

- [Setup](setup.md) — opening and resetting the database
- [Mutations](mutations.md) — updating and deleting
- [Transactions](transactions.md) — several changes, all or nothing

---

Back to: [Database](./README.md) &mdash; Next: [Mutations](mutations.md)
