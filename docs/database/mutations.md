# Zentests - Docs - Database - Mutations

> **Status:** DRAFT — structure and wording are open for discussion.

Updating and deleting records mid-test, to set up the state a case needs or to trigger a code path
that depends on missing data.

## Updating

| Function                                           | Target                         |
| -------------------------------------------------- | ------------------------------ |
| `DBUpdate(t, db, &User{}, id, updates)`            | One row, by primary key        |
| `DBUpdateBy(t, db, &User{}, where, args, updates)` | Every row matching a condition |

```go
zentests.DBUpdate(t, db, &User{}, 1, map[string]any{"status": "active"})

zentests.DBUpdateBy(t, db, &User{}, "status = ?", []any{"pending"}, map[string]any{
	"processed": true,
})
```

- `updates` is a map, so only the named columns change. Zero values are written, which is why the
  update uses `Updates` rather than `Save`.
- Passing `nil` as the ID updates every row: `DBUpdate(t, db, &User{}, nil, updates)`.
- `DBUpdateBy` takes its arguments as a `[]any`, not variadic, because `updates` comes after them.

## Deleting

| Function                                     | Target                                            |
| -------------------------------------------- | ------------------------------------------------- |
| `DBDelete(t, db, record)`                    | The row for that specific instance                |
| `DBDeleteBy(t, db, &User{}, where, args...)` | Every row matching a condition, returns the count |

```go
user := zentests.DBCreate(t, db, &User{Name: "Temp"})
zentests.DBDelete(t, db, user)

deleted := zentests.DBDeleteBy(t, db, &User{}, "status = ?", "deleted")
// deleted is an int64 row count, useful for an assertion
```

`DBDelete` uses the record's primary key, so the instance must have been created or loaded first.
`DBDeleteBy` returns the affected row count instead of failing when nothing matches.

## Soft deletes

If your model embeds `gorm.DeletedAt`, both delete functions respect GORM's soft-delete rules: the
row stays in the table with `deleted_at` set, and `DBFind`/`DBCount` no longer see it. Use
`db.Unscoped()` directly when a test must assert on the soft-deleted row.

## Notes

- **Updates are not transactions.** Two calls in a row can leave the database half-changed if the
  second fails; use [DBTx](transactions.md) when both must succeed or neither.
- **No hooks bypass.** These helpers call GORM's `Updates`/`Delete`, so model hooks and `UpdatedAt`
  behave as they do in application code.
- **Conditions are GORM syntax**, not raw SQL strings concatenated by the helper.

## Related

- [Records](records.md) — creating and reading
- [Transactions](transactions.md) — several mutations, all or nothing

---

Back to: [Database](./README.md) &mdash; Next: [Transactions](transactions.md)
