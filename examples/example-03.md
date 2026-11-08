# Zentests - Examples - Example 03 - testify/suite` Style

A suite builds the app and the database once, resets data between tests, and keeps the `zentests`
helper on the suite. Use it when several tests share expensive setup, or when the project already
uses `testify/suite`.

## Run it

```bash
mkdir zentests-suite && cd zentests-suite
go mod init example.com/zentests-suite
go get github.com/gofiber/fiber/v3 github.com/kematzy/zentests github.com/stretchr/testify
# save the code below as api_test.go
go test -v ./...
```

The database helpers use GORM's SQLite driver, which links a C library: build with `CGO_ENABLED=1`
and a C compiler available. See [docs/database](../docs/database/README.md).

## The code

```go
package api_test

import (
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/kematzy/zentests"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// User is the model under test.
type User struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Name  string `json:"name"`
	Email string `gorm:"uniqueIndex" json:"email"`
}

type UserSuite struct {
	suite.Suite
	db  *gorm.DB
	app *fiber.App
	zt  *zentests.T
}

// SetupSuite runs once, before the first test in the suite.
func (s *UserSuite) SetupSuite() {
	// In-memory SQLite with the User table migrated; closed automatically after the suite.
	s.db = zentests.SetupTestDBWithModels(s.T(), &User{})
	s.zt = zentests.New(s.T())

	s.app = fiber.New()
	s.app.Get("/api/v1/users", func(c fiber.Ctx) error {
		var users []User
		if err := s.db.Find(&users).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"users": users})
	})
}

// SetupTest runs before every test: same schema, empty tables.
func (s *UserSuite) SetupTest() {
	zentests.DBReset(s.T(), s.db)
}

func (s *UserSuite) TestListUsers_Empty() {
	s.zt.Get(s.app, "/api/v1/users").
		OK().
		IsJSON().
		ArrayLength("users", 0)
}

func (s *UserSuite) TestListUsers_Seeded() {
	zentests.DBCreateN(s.T(), s.db, 3, func(i int) User {
		return User{Name: fmt.Sprintf("User %d", i), Email: fmt.Sprintf("user%d@example.com", i)}
	})

	s.zt.Get(s.app, "/api/v1/users").
		OK().
		IsJSON().
		ArrayLength("users", 3).
		HasString("users.0.name", "User 1")
}

func (s *UserSuite) TestCreateAndFind() {
	user := zentests.DBCreate(s.T(), s.db, &User{Name: "Jane", Email: "jane@example.com"})
	s.NotZero(user.ID)

	zentests.DBExists(s.T(), s.db, &User{}, "email = ?", "jane@example.com")
}

// TestUserSuite is the single entry point that runs every method above.
func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}
```

## What it shows

- `SetupSuite` for work that is expensive and reusable: one database, one app, one `*T`.
- `SetupTest` for isolation: `DBReset` keeps the schema and empties the tables.
- `DBCreateN` with a 1-based factory to seed rows without repeating literals.
- `DBCreate` + `DBExists` for a record assertion that never goes through HTTP.
- Test methods named `TestX` on the suite struct; `suite.Run` is the only entry point.

## Lifecycle

```txt
SetupSuite()            once, before any test
  ├── SetupTest() → TestA() → TearDownTest()
  ├── SetupTest() → TestB() → TearDownTest()
  └── SetupTest() → TestC() → TearDownTest()
TearDownSuite()         once, after every test
```

## Related

- Suite guide: [docs/testing/testify-suite](../docs/testing/testify-suite.md)
- Database helpers: [docs/database](../docs/database/README.md)
- Other examples: [Example 01](example-01.md), [Example 02](example-02.md)
