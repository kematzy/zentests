# Zentests - Docs - Database - Transactions

> **Status:** DRAFT — structure and wording are open for discussion.

`DBTx` runs a function inside a transaction, commits when it returns `nil` and rolls back when it
returns an error — useful when a test needs to prove that a set of writes stands or falls together.

## Signature

```go
zentests.DBTx(t, db, func(tx *gorm.DB) error) error
```

The callback receives the transaction handle; use `tx`, not `db`, inside it. The returned error is
the one your function produced, so a test can assert on it.

## Committing

```go
err := zentests.DBTx(t, db, func(tx *gorm.DB) error {
	if err := tx.Create(&order).Error; err != nil {
		return err
	}
	return tx.Create(&payment).Error
})
require.NoError(t, err)

zentests.DBCount(t, db, &Order{})   // 1
zentests.DBCount(t, db, &Payment{}) // 1
```

## Rolling back

```go
err := zentests.DBTx(t, db, func(tx *gorm.DB) error {
	if err := tx.Create(&order).Error; err != nil {
		return err
	}
	return errors.New("payment declined")
})

require.Error(t, err)
zentests.DBCount(t, db, &Order{}) // 0 — the order was rolled back
```

That is the shape to use when testing the service code that owns the transaction: assert that the
error surfaces, then assert that nothing was written.

## When a transaction belongs in a test

- **Asserting atomicity** — a failure halfway through leaves no rows behind.
- **Proving a rollback path** — the handler returns 500 and the database is unchanged.
- **Building state that must be consistent** — an order and its line items, in one write.

For ordinary setup, a plain `DBCreate` is clearer; a transaction adds a failure mode without adding
value.

## Notes

- **A commit failure fails the test**, through `require`; a rollback does not, because rolling back
  is the expected outcome of an error.
- **`DBTx` does not nest automatically.** If your code under test opens its own transaction, use
  GORM savepoints or test the code path directly.
- **The callback's error is yours to assert on.** `DBTx` returns it unchanged after rolling back.

## Related

- [Records](records.md) — the writes a transaction wraps
- [Mutations](mutations.md) — single updates and deletes
- [testing/testify-suite](../testing/testify-suite.md) — where the database usually lives

---

Back to: [Database](./README.md)
