package zentests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupTestDB opens a fresh in-memory SQLite database configured for testing.
// The database uses silent logging to suppress GORM output during tests, and
// local time for all timestamp operations.
//
// The connection is automatically closed when the test finishes via [testing.T.Cleanup],
// so no explicit teardown is required in the common case. Use [CloseTestDB] only when
// you need to close the connection before the test ends (e.g. to test reconnect behavior).
//
// Parameters:
//   - t: The testing.T instance; marks this as a test helper and registers cleanup.
//
// Returns:
//   - *gorm.DB: A connected, empty database. No tables are created — use
//     [SetupTestDBWithModels] or [DBMigrate] to add schema.
//
// Example:
//
//	db := zentests.SetupTestDB(t)
//	// db is connected; no tables yet
//	zentests.DBMigrate(t, db, &Post{})
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		NowFunc: func() time.Time { return time.Now().Local() },
		Logger:  logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err, "SetupTestDB: failed to open in-memory SQLite")

	t.Cleanup(func() {
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	})

	return db
}

// SetupTestDBWithModels opens an in-memory SQLite database and runs AutoMigrate
// for every provided model in a single call. This is the most common entry point
// for database tests — it replaces the three-line open/migrate/check sequence with
// one line.
//
// Parameters:
//   - t: The testing.T instance.
//   - models: One or more pointers to GORM model structs (e.g. &User{}, &Post{}).
//
// Returns:
//   - *gorm.DB: A connected database with the requested tables created.
//
// Example — suite-level setup, data reset between tests:
//
//	func (s *MySuite) SetupSuite() {
//	    s.db = zentests.SetupTestDBWithModels(s.T(), &User{}, &Post{})
//	}
//
//	func (s *MySuite) SetupTest() {
//	    zentests.DBReset(s.T(), s.db)
//	}
func SetupTestDBWithModels(t *testing.T, models ...any) *gorm.DB {
	t.Helper()

	db := SetupTestDB(t)
	DBMigrate(t, db, models...)

	return db
}

// DBMigrate runs GORM AutoMigrate for the provided models on an existing database.
// Use this when you have a DB from [SetupTestDB] and want to add or update tables
// without recreating the connection.
//
// Parameters:
//   - t: The testing.T instance.
//   - db: An open *gorm.DB connection.
//   - models: One or more pointers to GORM model structs.
//
// Example:
//
//	db := zentests.SetupTestDB(t)
//	zentests.DBMigrate(t, db, &User{}, &Post{})
func DBMigrate(t *testing.T, db *gorm.DB, models ...any) {
	t.Helper()

	require.NoError(t, db.AutoMigrate(models...), "DBMigrate: AutoMigrate failed")
}

// DBReset deletes all rows from every table in the database and resets SQLite auto-increment
// sequences, leaving the schema intact. This is faster than recreating the database and is
// the right tool for [testing.T.SetupTest] hooks in suites that share a single
// database connection across tests.
//
// NOTE: This function is SQLite-specific. The sequence reset step (sqlite_sequence) is a
// no-op on other databases but the DELETE statements will still clear data correctly on
// any GORM-supported driver.
//
// Parameters:
//   - t: The testing.T instance.
//   - db: An open *gorm.DB connection.
//
// Example:
//
//	func (s *MySuite) SetupTest() {
//	    zentests.DBReset(s.T(), s.db)
//	}
func DBReset(t *testing.T, db *gorm.DB) {
	t.Helper()

	tables, err := db.Migrator().GetTables()
	require.NoError(t, err, "DBReset: could not list tables")

	for _, table := range tables {
		require.NoError(t,
			db.Exec("DELETE FROM "+table).Error, //nolint:gosec // table names come from the migrator, not user input
			"DBReset: failed to delete from table %s", table,
		)
	}

	// Reset SQLite auto-increment counters. This is a no-op if the
	// sqlite_sequence table does not exist (i.e. no AUTOINCREMENT columns).
	db.Exec("DELETE FROM sqlite_sequence") //nolint:errcheck // intentionally best-effort
}

// CloseTestDB explicitly closes the underlying database connection.
// In most tests this is not needed because [SetupTestDB] registers an automatic close via
// [testing.T.Cleanup]. Use this only when you need to close the connection before the test ends
// — for example, to verify reconnect behavior or to release a file-based SQLite lock mid-test.
//
// Parameters:
//   - t: The testing.T instance.
//   - db: An open *gorm.DB connection.
//
// Example:
//
//	db := zentests.SetupTestDB(t)
//	// ... tests ...
//	zentests.CloseTestDB(t, db) // close early; t.Cleanup will be a no-op
func CloseTestDB(t *testing.T, db *gorm.DB) {
	t.Helper()

	sqlDB, err := db.DB()
	require.NoError(t, err, "CloseTestDB: could not retrieve underlying sql.DB")
	require.NoError(t, sqlDB.Close(), "CloseTestDB: close failed")
}

// DBCreate saves a single record to the database and returns it with any auto-populated fields
// (ID, CreatedAt, etc.) filled in. The test is failed immediately if the insert fails.
//
// Type parameter T must be a GORM model struct; pass a pointer to an initialized value.
//
// Parameters:
//   - t: The testing.T instance.
//   - db: An open *gorm.DB with the model's table already migrated.
//   - record: Pointer to the record to insert.
//
// Returns:
//   - *T: The same pointer, with auto-generated fields populated by GORM.
//
// Example:
//
//	user := zentests.DBCreate(t, db, &User{Name: "Alice", Email: "alice@example.com"})
//
//	s.NotZero(user.ID)      // ID populated after save
//	s.NotZero(user.CreatedAt)
func DBCreate[T any](t *testing.T, db *gorm.DB, record *T) *T {
	t.Helper()

	require.NoError(t, db.Create(record).Error, "DBCreate: insert failed")

	return record
}

// DBCreateN creates count records by calling factory once per record with a 1-based index
// (1, 2, … count). Each record is inserted and returned with auto-populated fields.
// The test is failed immediately if any insert fails.
//
// The 1-based index means you can use i directly in format strings without the i+1 offset:
//
//	zentests.DBCreateN(t, db, 5, func(i int) User {
//	    return User{Email: fmt.Sprintf("user%d@example.com", i)} // user1, user2, ...
//	})
//
// Type parameter T must be a GORM model struct (not a pointer — the factory returns values;
// DBCreateN handles taking their address internally).
//
// Parameters:
//   - t: The testing.T instance.
//   - db: An open *gorm.DB with the model's table already migrated.
//   - count: Number of records to create (must be > 0).
//   - factory: Function called with i = 1..count; returns an initialized T value.
//
// Returns:
//   - []*T: Slice of inserted record pointers, each with auto-generated fields.
//
// Example:
//
//	posts := zentests.DBCreateN(t, db, 3, func(i int) Post {
//	    return Post{Title: fmt.Sprintf("Post %d", i), Body: "content"}
//	})
//
//	s.Len(posts, 3)
//	s.Equal("Post 1", posts[0].Title)
func DBCreateN[T any](t *testing.T, db *gorm.DB, count int, factory func(i int) T) []*T {
	t.Helper()

	records := make([]*T, count)

	for i := range count {
		v := factory(i + 1) // 1-based index
		records[i] = DBCreate(t, db, &v)
	}

	return records
}
