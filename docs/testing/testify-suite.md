# Zentests - Docs - Testing - testify/suite

`testify/suite` provides suite-level and per-test lifecycle hooks. Combine it with `zentests` when
several tests share an expensive resource, such as a migrated database or app, or when the project
already uses suites. `zentests` supplies the request helper and database helpers; `testify/suite`
continues to run the test methods.

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
	s.app = fiber.New()
	// routes ...
}

func (s *UserSuite) SetupTest() {
	s.zt = zentests.New(s.T()) // bind assertions to this test's *testing.T
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
- A fatal failure (`require`/`FailNow`) in `SetupSuite` exits before test methods run and before
  `TearDownSuite` is registered. A non-fatal `assert` failure marks the suite test failed but allows
  setup to return normally.
- `s.Require()` stops the current test; `s.Assert()` records the failure and continues.

## `s.T()` and cleanup — the one thing to get right

`s.T()` returns the `*testing.T` of the **currently running** test. In `SetupSuite` that is the
suite-level `T`, in `SetupTest` and in a test method it is that test's `T`.

Resource helpers such as `SetupTestDBWithModels` and `NewApp` register cleanup with `t.Cleanup`.
Cleanup runs at the scope of the `*testing.T` passed to the helper:

| Called in    | Cleanup runs                     |
| ------------ | -------------------------------- |
| `SetupSuite` | after the last test in the suite |
| `SetupTest`  | after that one test              |
| test method  | after that one test              |

So `SetupTestDBWithModels(s.T(), ...)` in `SetupSuite` keeps one database for the whole suite, and
`zentests.NewApp(s.T())` in `SetupSuite` keeps one app for the whole suite. Both register cleanup on
the suite-level test; neither needs a manual close in `TearDownSuite`. `zentests.New` itself owns no
resources; create it in `SetupTest` so assertions are attached to each individual test.

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
}

func (s *BaseSuite) SetupTest() {
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
		Timeout:       10 * time.Second,
		FailOnTimeout: true,
	}).OK()
}
```

## Pitfalls

- **Forgetting `suite.Run`.** A suite without it compiles and never runs.
- **`SetupSuite` instead of `SetupTest` for data.** Rows created in `SetupSuite` are visible to
  every test; reset them in `SetupTest` or create them per test.
- **Parallel test methods are unsupported.** Do not call `s.T().Parallel()` in a suite; suite state
  and its current `*testing.T` are shared across methods. Create a new `zentests.T` in `SetupTest`.
- **Relying on `TearDownSuite` after fatal `SetupSuite` failure.** It is not registered until setup
  returns. Use resource helpers that register `t.Cleanup` for resources that must be released even
  when setup stops early.

## Working code

[examples/03-testify-suite](../../examples/example-03.md) — suite with an in-memory database, three
tests, runnable as-is.

## Related

- [Testing styles overview](./README.md)
- [Database helpers](../database/README.md)
- [Request methods](../requests/methods.md)

---

Back to: [Testing styles](./README.md)
