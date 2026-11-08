# Zentests - Docs - Testing - testify/suite

> **Status:** DRAFT — structure and wording are open for discussion.

`testify/suite` builds setup once for a group of tests. Combined with `zentests` it is the right
choice when several tests share an expensive resource — a migrated database, a fixture app, a seeded
table — or when the project already has suite-based tests.

Everything in the suite is ordinary `testify`: `zentests` only supplies the helper and the database
helpers.

## The shape

```go
type UserSuite struct {
	suite.Suite
	db  *gorm.DB
	app *fiber.App
	zt  *zentests.T
}

func (s *UserSuite) SetupSuite() {
	s.db = zentests.SetupTestDBWithModels(s.T(), &User{}) // closed after the suite
	s.zt = zentests.New(s.T())
	s.app = fiber.New()
	// routes ...
}

func (s *UserSuite) SetupTest() {
	zentests.DBReset(s.T(), s.db) // same schema, empty tables
}

func (s *UserSuite) TestListUsers() {
	s.zt.Get(s.app, "/api/v1/users").OK().IsJSON().ArrayLength("users", 0)
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}
```

Only methods named `TestX` on the suite struct run; `suite.Run` is the single entry point.

## Lifecycle

```txt
SetupSuite()            once, before any test method
  ├── SetupTest() → TestA() → TearDownTest()
  ├── SetupTest() → TestB() → TearDownTest()
  └── SetupTest() → TestC() → TearDownTest()
TearDownSuite()         once, after every test method
```

- `SetupTest`/`TearDownTest` run around **each** test method. `TearDownTest` is deferred before
  `SetupTest` runs, so it still runs when setup fails.
- A failure in `SetupSuite` skips every test **and** `TearDownSuite` — its `defer` is only
  registered after `SetupSuite` returns.
- `s.Require()` stops the current test; `s.Assert()` records the failure and continues.

## `s.T()` and cleanup — the one thing to get right

`s.T()` returns the `*testing.T` of the **currently running** test. In `SetupSuite` that is the
suite-level `T`, in `SetupTest` and in a test method it is that test's `T`.

This matters because every `zentests` helper registers cleanup with `t.Cleanup`:

| Called in    | Cleanup runs                     |
| ------------ | -------------------------------- |
| `SetupSuite` | after the last test in the suite |
| `SetupTest`  | after that one test              |
| test method  | after that one test              |

So `SetupTestDBWithModels(s.T(), ...)` in `SetupSuite` keeps one database for the whole suite, and
`zentests.NewApp(s.T())` in `SetupSuite` keeps one app for the whole suite. Neither needs a manual
close in `TearDownSuite`.

## Database isolation between tests

`DBReset` deletes every row from every migrated table and resets SQLite's auto-increment counter.
The schema stays, so it is cheap enough to call in `SetupTest` for every test:

```go
func (s *UserSuite) SetupTest() {
	zentests.DBReset(s.T(), s.db)
}
```

Use `DBCreate`, `DBCreateN` or `DBSeed` inside a test to build only the rows that test needs.

## Shared base suites

Embed a base suite for setup that several suites repeat. If the base defines a lifecycle method and
the outer suite defines the same one, the outer method **shadows** it — call the base explicitly:

```go
type BaseSuite struct {
	suite.Suite
	app *fiber.App
	zt  *zentests.T
}

func (s *BaseSuite) SetupSuite() {
	s.app = zentests.NewApp(s.T())
	s.zt = zentests.New(s.T())
}

type UserSuite struct{ BaseSuite }

func (s *UserSuite) SetupSuite() {
	s.BaseSuite.SetupSuite() // required, otherwise the base setup never runs
	s.app.Get("/api/v1/users", listUsers)
}
```

## Timeouts inside a suite

`TestConfig` works the same way; use the `*WithConfig` methods for slow handlers:

```go
func (s *SlowSuite) TestReport() {
	s.zt.GetWithConfig(s.app, "/api/v1/report", fiber.TestConfig{
		Timeout: 10 * time.Second,
	}).OK()
}
```

## Pitfalls

- **Forgetting `suite.Run`.** A suite without it compiles and never runs.
- **`SetupSuite` instead of `SetupTest` for data.** Rows created in `SetupSuite` are visible to
  every test; reset them in `SetupTest` or create them per test.
- **Sharing one `*T` across parallel tests.** `zentests.T` wraps a specific `*testing.T`; do not
  reuse it inside `t.Parallel()` subtests — create one per test.
- **Relying on `TearDownSuite` after a failed `SetupSuite`.** It does not run, so keep essential
  cleanup in `t.Cleanup` (which `zentests` helpers already use) rather than in `TearDownSuite`.

## Working code

[examples/03-testify-suite](../../examples/example-03.md) — suite with an in-memory database, three
tests, runnable as-is.

## Related

- [Testing styles overview](./README.md)
- [Database helpers](../database/README.md)
- [Request methods](../requests/methods.md)

---

Back to: [Testing styles](./README.md)
