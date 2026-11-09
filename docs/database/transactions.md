# Zentests - Docs - Database - Transactions

`DBTx` runs a callback inside a transaction: it commits when the callback returns `nil` and rolls
back when it returns an error. Use it to test writes that must succeed or fail together.

## Signature

```go
zentests.DBTx(t, db, func(tx *gorm.DB) error) error
```

The callback receives the transaction handle; use `tx`, not `db`, for writes that must participate
in the transaction. Its returned error is surfaced after rollback; `nil` commits the transaction.

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

- **Begin and commit failures fail the test** through `require`; rollback errors are not returned by
  `DBTx`.
- **`DBTx` does not nest automatically.** If your code under test opens its own transaction, use
  GORM savepoints or test the code path directly.
- **Return errors from the callback to request rollback.** Do not call `require`/`FailNow` inside it:
  those stop the test before `DBTx` can roll back. Assert on the returned callback error afterward.

## Related

- [Records](records.md) — the writes a transaction wraps
- [Mutations](mutations.md) — single updates and deletes
- [testing/testify-suite](../testing/testify-suite.md) — where the database usually lives

---

Back to: [Database](./README.md)
