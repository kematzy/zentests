# Using zentests with testify/suite

`zentests` integrates seamlessly with [`testify/suite`](https://github.com/stretchr/testify#suite-package)
for more complex test scenarios requiring suite-level setup/teardown.

See the [README](../README.md) for the zentests basics, and [API.md](API.md) for the full reference.

## Basic Structure

```go
package myapp_test

import (
    "testing"
    "github.com/gofiber/fiber/v3"
    "github.com/kematzy/zentests"
    "github.com/stretchr/testify/suite"
)

// Define your suite struct
type UserAPISuite struct {
    suite.Suite
    app *fiber.App
    zt  *zentests.T
}

// SetupSuite runs ONCE before all tests in the suite
func (s *UserAPISuite) SetupSuite() {
    s.app = fiber.New()
    s.zt = zentests.New(s.T())  // Initialize zentests.T

    // Setup routes
    s.app.Get("/users", func(c fiber.Ctx) error {
        return c.JSON(fiber.Map{"users": []string{"John", "Jane"}})
    })

    s.app.Post("/users", func(c fiber.Ctx) error {
        return c.Status(201).JSON(fiber.Map{"id": 1, "name": "John"})
    })
}

// TearDownSuite runs ONCE after all tests in the suite
func (s *UserAPISuite) TearDownSuite() {
    // Cleanup resources (close database, etc.)
}

// SetupTest runs BEFORE EACH test
func (s *UserAPISuite) SetupTest() {
    // Reset state before each test
}

// TearDownTest runs AFTER EACH test
func (s *UserAPISuite) TearDownTest() {
    // Cleanup after each test
}

// Tests use the zentests.T instance
func (s *UserAPISuite) TestListUsers() {
    s.zt.Get(s.app, "/users").
        OK().
        IsJSON().
        ArrayLength("users", 2)
}

func (s *UserAPISuite) TestCreateUser() {
    s.zt.PostJSON(s.app, "/users", map[string]any{
        "name": "John",
    }).
        Created().
        Has("id", float64(1))
}

// Run the suite
func TestUserAPISuite(t *testing.T) {
    suite.Run(t, new(UserAPISuite))
}
```

## Lifecycle Hooks Order

```
SetupSuite()           ← Runs once before all tests
    ├── TestA()
    │       ├── SetupTest()
    │       ├── [test code]
    │       └── TearDownTest()
    │
    └── TestB()
            ├── SetupTest()
            ├── [test code]
            └── TearDownTest()
TearDownSuite()        ← Runs once after all tests
```

## Database Integration

```go
type DatabaseSuite struct {
    suite.Suite
    app *fiber.App
    zt  *zentests.T
    db  *sql.DB  // or *gorm.DB
}

func (s *DatabaseSuite) SetupSuite() {
    // Connect to test database
    s.db, _ = sql.Open("sqlite3", ":memory:")

    s.app = fiber.New()
    s.zt = zentests.New(s.T())

    // Setup routes that use database
    s.app.Get("/users", func(c fiber.Ctx) error {
        var users []User
        s.db.Find(&users)
        return c.JSON(users)
    })
}

func (s *DatabaseSuite) SetupTest() {
    // Clean database before each test
    s.db.Exec("DELETE FROM users")
}

func (s *DatabaseSuite) TearDownSuite() {
    s.db.Close()
}
```

## Multiple Suites with Shared Setup

```go
// Base suite with common setup
type BaseSuite struct {
    suite.Suite
    App *fiber.App
    ZT  *zentests.T
}

func (s *BaseSuite) SetupSuite() {
    s.App = fiber.New()
    s.ZT = zentests.New(s.T())
}

// UserAPISuite inherits BaseSuite
type UserAPISuite struct {
    BaseSuite
}

func (s *UserAPISuite) SetupSuite() {
    s.BaseSuite.SetupSuite()
    // Add user-specific routes
    s.App.Get("/users", s.listUsers)
}

func TestUserAPISuite(t *testing.T) {
    suite.Run(t, new(UserAPISuite))
}
```

## Using TestConfig with Suite

```go
func (s *SlowAPISuite) TestSlowEndpoint() {
    // Use WithConfig for slow endpoints
    s.zt.GetWithConfig(s.app, "/slow-endpoint", fiber.TestConfig{
        Timeout: 10 * time.Second,
    }).OK()
}
```
