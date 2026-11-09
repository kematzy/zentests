# Zentests - Docs - Testing Styles

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

| Style                               | Setup                                   | Teardown / isolation                           | Best for                                               |
| ----------------------------------- | --------------------------------------- | ---------------------------------------------- | ------------------------------------------------------ |
| [Flat](flat.md)                     | `zentests.New(t)` and fixtures per test | No listener; resource helpers register cleanup | Independent tests; parallel when fixtures are isolated |
| [BDD](./bdd.md)                     | `BeforeEach` per `Describe` block       | `AfterEach` per `It`                           | Named behavior groups and readable subtest trees       |
| [testify/suite](./testify-suite.md) | `SetupSuite` / `SetupTest`              | `TearDownTest` / `TearDownSuite`               | Shared expensive fixtures; existing suite codebases    |

Rules of thumb:

- Start **flat**. Move to BDD when the test names stop describing the behaviour.
- Move to **testify/suite** when setup is expensive (a database, a fixture server) and rebuilding it
  per test costs more than the isolation is worth.
- If the project already has `testify/suite` suites, keep them and add `zentests` inside — there is
  nothing to migrate.

## What every style shares

- `zentests.New(t)` creates a helper for the current test; response assertions chain from the
  `*Response` returned by each request.
- Create a helper per test. A new helper does not isolate state in a shared app, database, or cookie
  jar; reset or recreate those fixtures as needed.
- `fiber.App.Test` handles requests in-process. No network listener or port is opened.
- Resource helpers register cleanup through `t.Cleanup`; close a connection early only when a test
  needs to release or reopen it explicitly.
- `ZENTESTS_TIMEOUT` and `fiber.TestConfig` behave identically in all three styles.

## Examples

Each style has a runnable example; the examples show how setup and request tests differ by style:

- [examples/01-flat-style](../../examples/example-01.md)
- [examples/02-bdd](../../examples/example-02.md)
- [examples/03-testify-suite](../../examples/example-03.md)

---

Back to: [README](../README.md)
