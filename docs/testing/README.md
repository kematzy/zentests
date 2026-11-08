# Zentests - Docs - Testing Styles

> **Status:** DRAFT — structure and wording are open for discussion.

`zentests` works with three styles of Go tests. The helper is the same in all of them
(`zt := zentests.New(t)`); what changes is where the app, the database and the assertions live.

Pick one style per file — mixing styles inside one package is fine and common, but a file should
read as one idea.

## Contents

- [Which style](#which-style)
- [What every style shares](#what-every-style-shares)
- [Flat tests](flat.md)
- [BDD tests](./bdd.md)
- [testify/suite tests](./testify-suite.md)

## Which style

| Style                               | Setup                              | Teardown                         | Best for                                                           |
| ----------------------------------- | ---------------------------------- | -------------------------------- | ------------------------------------------------------------------ |
| [Flat](flat.md)                     | `zentests.New(t)` inside each test | none — `t.Cleanup` handles it    | independent tests, small files, `t.Parallel()`                     |
| [BDD](./bdd.md)                     | `BeforeEach` on a `Describe` block | `AfterEach`                      | grouping by behaviour, readable subtest trees                      |
| [testify/suite](./testify-suite.md) | `SetupSuite` / `SetupTest`         | `TearDownTest` / `TearDownSuite` | one app or database shared by many tests, existing suite codebases |

Rules of thumb:

- Start **flat**. Move to BDD when the test names stop describing the behaviour.
- Move to **testify/suite** when setup is expensive (a database, a fixture server) and rebuilding it
  per test costs more than the isolation is worth.
- If the project already has `testify/suite` suites, keep them and add `zentests` inside — there is
  nothing to migrate.

## What every style shares

- `zentests.New(t)` returns the helper; every assertion hangs off a `*Response` returned by a
  request.
- Each test gets its own helper and its own `*Response`; nothing is global.
- `fiber.App.Test` drives requests in-process, so no listener, port or goroutine is involved and no
  shutdown is required.
- Database helpers register their own cleanup through `t.Cleanup`; only close a connection early
  when a test is about reconnecting.
- `ZENTESTS_TIMEOUT` and `fiber.TestConfig` behave identically in all three styles.

## Examples

Each style has a runnable example, with the same app tested three ways:

- [examples/01-flat-style](../../examples/example-01.md)
- [examples/02-bdd](../../examples/example-02.md)
- [examples/03-testify-suite](../../examples/example-03.md)

---

Back to: [README](../README.md)
