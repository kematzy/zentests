# Zentests - Docs - Database - Records

Creating and reading records without repeating `db.Create(...)` plus an error check in every test.

## Creating

| Helper                                  | Behavior                                                   | Returns                        |
| --------------------------------------- | ---------------------------------------------------------- | ------------------------------ |
| `DBCreate(t, db, record)`               | Inserts one record                                         | Same pointer with GORM fields  |
| `DBCreateN(t, db, count, factory)`      | Inserts `count` values; factory indices run from 1 to count | Slice of inserted pointers     |

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

| Helper                         | Input   | Behavior                    |
| ------------------------------ | ------- | --------------------------- |
| `DBSeed(t, db, records)`       | `[]*T`  | Inserts each pointer in turn |
| `DBSeedSlice(t, db, records)`  | `[]T`   | Inserts each value in turn   |

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

`DBSeed` takes pointers; `DBSeedSlice` takes values and creates pointers internally. Both insert
one record at a time. A constraint violation fails the test on that insert; `DBSeedSlice` includes
the record index in its failure message.

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
pointer to the populated record otherwise. `DBCount` returns a count and fails only on a database
error. `DBExists` and `DBNotExists` are assertions: each returns `true` when its condition passes
and fails the test when it does not.

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
- **Ordering is by primary key.** `DBFirst` and `DBLast` use primary-key order, not insertion time.
  Use a GORM query with an explicit `Order` clause when you need custom ordering.

## Related

- [Setup](setup.md) — opening and resetting the database
- [Mutations](mutations.md) — updating and deleting
- [Transactions](transactions.md) — several changes, all or nothing

---

Back to: [Database](./README.md) &mdash; Next: [Mutations](mutations.md)
