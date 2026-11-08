# Zentests - Docs - Testing - BDD

> **Status:** DRAFT — structure and wording are open for discussion.

BDD style groups tests with `Describe`, runs setup with `BeforeEach`, and writes one `It` per
behaviour. It produces a readable sub-test tree in the output and needs no dependency beyond
`zentests`.

Use it when the _names_ of the tests are the documentation — a public API, a workflow, a set of
rules — rather than a long list of unrelated checks.

## The shape

```go
func TestUserAPI(t *testing.T) {
	zt := zentests.New(t)
	var app *fiber.App

	zt.Describe("User API", func(ctx *zentests.SuiteContext) {
		ctx.BeforeEach(func(zt *zentests.T) {
			app = newApp() // fresh app before every It
		})

		ctx.It("reports its health", func(zt *zentests.T) {
			zt.Get(app, "/api/v1/healthz").OK()
		})
	})
}
```

`Describe` creates a subtest, `It` creates a sub-test inside it, so the output is a tree:

```
--- PASS: TestUserAPI
    --- PASS: TestUserAPI/User_API
        --- PASS: TestUserAPI/User_API/reports_its_health
```

## API

| Method                   | Purpose                                                 |
| ------------------------ | ------------------------------------------------------- |
| `zt.Describe(name, fn)`  | Top-level group; receives a `*SuiteContext`             |
| `ctx.Describe(name, fn)` | Nested group; inherits the enclosing block's hooks      |
| `ctx.BeforeEach(fn)`     | Runs before each `It`; returns the context for chaining |
| `ctx.AfterEach(fn)`      | Runs after each `It`; returns the context for chaining  |
| `ctx.It(name, fn)`       | One behaviour; receives a fresh `*zentests.T`           |

## Nesting and hook order

Hooks of enclosing blocks apply to nested blocks:

- `BeforeEach` runs **outermost first**, so shared setup exists before specific setup.
- `AfterEach` runs **innermost first**, so specific teardown runs before shared teardown.

```go
zt.Describe("Users", func(ctx *zentests.SuiteContext) {
	ctx.BeforeEach(func(zt *zentests.T) { app = newApp() })

	ctx.Describe("when logged in", func(ctx *zentests.SuiteContext) {
		ctx.BeforeEach(func(zt *zentests.T) { zt.WithBearer("secret-token") })

		ctx.It("lists users", func(zt *zentests.T) {
			zt.Get(app, "/users").OK()
		})
	})
})
```

For the `It` above: `newApp()` runs, then the bearer header is applied, then the test runs, then the
inner `AfterEach`, then the outer one.

## Rules that surprise people

- **Hooks are captured when `It` is called.** Declare `BeforeEach`/`AfterEach` _before_ the `It`
  calls they should apply to; a hook added afterwards does not retroactively apply.
- **`AfterEach` runs even when the test stops early.** Hooks are registered with `t.Cleanup`, so a
  `require` failure inside `BeforeEach` or the test body still runs teardown.
- **Each `It` gets a fresh `*zentests.T`.** Headers, cookies and state set inside one test never
  leak into the next, even when the app is shared.
- **Sub-tests run sequentially.** There is no `t.Parallel()` hook; for parallel tests use
  [flat](flat.md) or [testify/suite](./testify-suite.md).

## Limits

- No suite-level setup: work that should happen once for many tests belongs in
  [testify/suite](./testify-suite.md), or in a plain helper function called from `BeforeEach`.
- No shared base blocks: `Describe` bodies are functions, so compose with ordinary Go helpers.

## Working code

[examples/02-bdd](../../examples/example-02.md) — nested blocks, both hooks, runnable as-is.

## Related

- [Testing styles overview](./README.md)
- [testify/suite](./testify-suite.md)

---

Back to: [Testing styles](./README.md) &mdash; Next: [testify/suite](./testify-suite.md)
