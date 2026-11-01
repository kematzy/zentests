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

// =================================================================================================
// QUERY HELPERS
// =================================================================================================

// DBFind finds a record by ID and populates the result into the provided pointer.
// Fails the test if the record is not found or on database error.
// This is the idiomatic way to fetch a record by primary key in tests.
//
// Parameters:
//   - t: The testing.T instance.
//   - db: An open *gorm.DB with the model's table migrated.
//   - record: Pointer to model struct to populate (must be allocated).
//   - id: The primary key value to find.
//
// Returns:
//   - *T: The same pointer, populated with the found record.
//
// Example:
//
//	user := &User{}
//	zentests.DBFind(t, db, user, 1)
//	s.Equal("alice", user.Name)
func DBFind[T any](t *testing.T, db *gorm.DB, record *T, id any) *T {
	t.Helper()

	err := db.First(record, id).Error
	require.NoError(t, err, "DBFind: record not found with id=%v", id)

	return record
}

// DBFindBy finds a record by a specific field condition.
// Fails the test if exactly one record is not found or on error.
// Use this for non-ID lookups (e.g., by email, unique constraint).
//
// Parameters:
//   - t: The testing.T instance.
//   - db: An open *gorm.DB with the model's table migrated.
//   - record: Pointer to model struct to populate.
//   - where: Where clause (e.g., "email = ?" or "name LIKE ?").
//   - args: Arguments matching the where clause placeholders.
//
// Returns:
//   - *T: The same pointer, populated with the found record.
//
// Example:
//
//	user := &User{}
//	zentests.DBFindBy(t, db, user, "email = ?", "alice@example.com")
//	s.Equal("Alice", user.Name)
func DBFindBy[T any](t *testing.T, db *gorm.DB, record *T, where string, args ...any) *T {
	t.Helper()

	err := db.Where(where, args...).First(record).Error
	require.NoError(t, err, "DBFindBy: record not found where %s", where)

	return record
}

// DBFirst retrieves the first record ordered by primary key.
// Fails the test if no records exist or on error.
// Use this when you need the oldest/first record by natural order.
//
// Parameters:
//   - t: The testing.T instance.
//   - db: An open *gorm.DB with the model's table migrated.
//   - record: Pointer to model struct to populate.
//
// Returns:
//   - *T: The same pointer, populated with the first record.
//
// Example:
//
//	user := &User{}
//	zentests.DBFirst(t, db, user)
//	s.Equal("Alice", user.Name) // first created user
func DBFirst[T any](t *testing.T, db *gorm.DB, record *T) *T {
	t.Helper()

	err := db.First(record).Error
	require.NoError(t, err, "DBFirst: no records found")

	return record
}

// DBLast retrieves the last record ordered by primary key (descending).
// Fails the test if no records exist or on error.
// Use this when you need the most recent record.
//
// Parameters:
//   - t: The testing.T instance.
//   - db: An open *gorm.DB with the model's table migrated.
//   - record: Pointer to model struct to populate.
//
// Returns:
//   - *T: The same pointer, populated with the last record.
//
// Example:
//
//	user := &User{}
//	zentests.DBLast(t, db, user)
//	s.Equal("Zoe", user.Name) // most recently created
func DBLast[T any](t *testing.T, db *gorm.DB, record *T) *T {
	t.Helper()

	err := db.Last(record).Error
	require.NoError(t, err, "DBLast: no records found")

	return record
}

// DBCount counts records matching the given conditions.
// Fails the test on database error.
//
// Parameters:
//   - t: The testing.T instance.
//   - db: An open *gorm.DB with the model's table migrated.
//   - model: Pointer to model struct (used only for type, can be &Model{}).
//   - conditions: Optional - pass a where string, then args. e.g., "status = ?", "active"
//
// Returns:
//   - int64: The count of matching records.
//
// Example:
//
//	count := zentests.DBCount(t, db, &User{})
//	s.Equal(int64(5), count)
//
//	count := zentests.DBCount(t, db, &User{}, "status = ?", "active")
func DBCount[T any](t *testing.T, db *gorm.DB, model *T, conditions ...any) int64 {
	t.Helper()

	var count int64
	query := db.Model(model)
	if len(conditions) > 0 {
		query = query.Where(conditions[0].(string), conditions[1:]...)
	}
	err := query.Count(&count).Error
	require.NoError(t, err, "DBCount: count failed")

	return count
}

// DBExists asserts that at least one record matches the given conditions.
// Fails the test if no records are found.
//
// Parameters:
//   - t: The testing.T instance.
//   - db: An open *gorm.DB with the model's table migrated.
//   - model: Pointer to model struct (type used for table).
//   - conditions: Optional - pass a where string, then args. e.g., "status = ?", "active"
//
// Returns:
//   - bool: Always true if check passes.
//
// Example:
//
//	s.True(zentests.DBExists(t, db, &User{}))
//	s.True(zentests.DBExists(t, db, &User{}, "status = ?", "active"))
func DBExists[T any](t *testing.T, db *gorm.DB, model *T, conditions ...any) bool {
	t.Helper()

	var count int64
	query := db.Model(model)
	if len(conditions) > 0 {
		query = query.Where(conditions[0].(string), conditions[1:]...)
	}
	err := query.Count(&count).Error
	require.NoError(t, err, "DBExists: count failed")

	require.Positive(t, count, "DBExists: no records found matching condition")

	return true
}

// DBNotExists asserts that no records match the given conditions.
// Fails the test if any records are found.
//
// Parameters:
//   - t: The testing.T instance.
//   - db: An open *gorm.DB with the model's table migrated.
//   - model: Pointer to model struct (type used for table).
//   - conditions: Optional - pass a where string, then args. e.g., "email = ?", "deleted@example.com"
//
// Returns:
//
//	bool: Always true if check passes.
//
// Example:
//
//	s.True(zentests.DBNotExists(t, db, &User{}))
//	s.True(zentests.DBNotExists(t, db, &User{}, "email = ?", "deleted@example.com"))
func DBNotExists[T any](t *testing.T, db *gorm.DB, model *T, conditions ...any) bool {
	t.Helper()

	var count int64
	query := db.Model(model)
	if len(conditions) > 0 {
		query = query.Where(conditions[0].(string), conditions[1:]...)
	}
	err := query.Count(&count).Error
	require.NoError(t, err, "DBNotExists: count failed")

	require.Zero(t, count, "DBNotExists: expected no records but found %d", count)

	return true
}

// =================================================================================================
// UPDATE/DELETE HELPERS
// =================================================================================================

// DBUpdate performs an update query and fails the test on error.
// Use this for simple field updates without loading the record first.
//
// Parameters:
//   - t: The testing.T instance.
//   - db: An open *gorm.DB with the model's table migrated.
//   - model: Pointer to model struct (used for table, must have ID if using byID).
//   - byID: Primary key value to identify record; use nil for where clause.
//   - updates: Map of column to new value (e.g., map[string]any{"status": "inactive"}).
//
// Example:
//
//	zentests.DBUpdate(t, db, &User{}, 1, map[string]any{"status": "inactive"})
func DBUpdate[T any](t *testing.T, db *gorm.DB, model *T, byID any, updates map[string]any) {
	t.Helper()

	query := db.Model(model)
	if byID != nil {
		query = query.Where("id = ?", byID)
	}
	err := query.Updates(updates).Error
	require.NoError(t, err, "DBUpdate: update failed")
}

// DBUpdateBy performs an update using a where clause.
// Use this when you need to update records matching a non-ID condition.
//
// Parameters:
//   - t: The testing.T instance.
//   - db: An open *gorm.DB with the model's table migrated.
//   - model: Pointer to model struct.
//   - where: Where clause (e.g., "status = ?" or "id > ?").
//   - args: Arguments for where clause placeholders.
//   - updates: Map of column to new value.
//
// Example:
//
//	zentests.DBUpdateBy(t, db, &User{}, "status = ?", "pending", map[string]any{"processed": true})
func DBUpdateBy[T any](t *testing.T, db *gorm.DB, model *T, where string, args []any, updates map[string]any) {
	t.Helper()

	err := db.Model(model).Where(where, args...).Updates(updates).Error
	require.NoError(t, err, "DBUpdateBy: update failed")
}

// DBDelete deletes a record by ID and fails the test on error.
//
// Parameters:
//   - t: The testing.T instance.
//   - db: An open *gorm.DB with the model's table migrated.
//   - record: Pointer to model struct with ID set.
//
// Example:
//
//	user := &User{ID: 1}
//	zentests.DBDelete(t, db, user)
func DBDelete[T any](t *testing.T, db *gorm.DB, record *T) {
	t.Helper()

	err := db.Delete(record).Error
	require.NoError(t, err, "DBDelete: delete failed")
}

// DBDeleteBy deletes records matching a where clause and fails the test on error.
//
// Parameters:
//   - t: The testing.T instance.
//   - db: An open *gorm.DB with the model's table migrated.
//   - model: Pointer to model struct (type used for table).
//   - where: Where clause.
//   - args: Arguments for where clause.
//
// Returns:
//   - int64: Number of records deleted.
//
// Example:
//
//	deleted := zentests.DBDeleteBy(t, db, &User{}, "status = ?", "deleted")
func DBDeleteBy[T any](t *testing.T, db *gorm.DB, model *T, where string, args ...any) int64 {
	t.Helper()

	result := db.Where(where, args...).Delete(model)
	err := result.Error
	require.NoError(t, err, "DBDeleteBy: delete failed")

	return result.RowsAffected
}

// =================================================================================================
// TRANSACTION HELPERS
// =================================================================================================

// DBTx executes a function within a transaction.
// The transaction is automatically rolled back if the function returns an error,
// otherwise it is committed when the function completes successfully.
// This is useful for testing multiple operations that should succeed or fail together.
//
// Parameters:
//   - t: The testing.T instance.
//   - db: An open *gorm.DB connection.
//   - fn: Function that executes operations on a transacted *gorm.DB.
//
// Returns:
//   - error: Any error returned by fn, which causes rollback.
//
// Example:
//
//	err := zentests.DBTx(t, db, func(tx *gorm.DB) error {
//	    if err := tx.Create(&user).Error; err != nil {
//	        return err
//	    }
//	    return tx.Create(&balance).Error
//	})
//	s.NoError(err)
func DBTx(t *testing.T, db *gorm.DB, fn func(tx *gorm.DB) error) error {
	t.Helper()

	tx := db.Begin()
	if tx.Error != nil {
		require.NoError(t, tx.Error, "DBTx: begin failed")
	}

	err := fn(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	require.NoError(t, tx.Commit().Error, "DBTx: commit failed")

	return nil
}

// =================================================================================================
// DATA SEEDING
// =================================================================================================

// DBSeed loads test data from a map of slices (table -> records).
// This is useful for populating a database with test fixtures defined inline.
// Each key in data should be a pointer to a model struct; the value should be a slice of
// pointers to already-populated model instances.
//
// Parameters:
//   - t: The testing.T instance.
//   - db: An open *gorm.DB.
//   - data: Map of model pointers to slices of populated records.
//
// Example:
//
//	users := []*User{
//	    {Name: "Alice", Email: "alice@example.com"},
//	    {Name: "Bob", Email: "bob@example.com"},
//	}
//	posts := []*Post{
//	    {Title: "First Post", Body: "Hello world"},
//	}
//	zentests.DBSeed(t, db, users, posts)
func DBSeed[T any](t *testing.T, db *gorm.DB, records []*T) {
	t.Helper()

	for _, record := range records {
		err := db.Create(record).Error
		require.NoError(t, err, "DBSeed: failed to create record")
	}
}

// DBSeedSlice loads test data from a slice of records.
// This is a convenience wrapper around DBSeed for single-table seeding.
//
// Parameters:
//   - t: The testing.T instance.
//   - db: An open *gorm.DB.
//   - records: Slice of pointers to already-populated model instances.
//
// Example:
//
//	users := []*User{
//	    {Name: "Alice", Email: "alice@example.com"},
//	    {Name: "Bob", Email: "bob@example.com"},
//	}
//	zentests.DBSeedSlice(t, db, users)
func DBSeedSlice[T any](t *testing.T, db *gorm.DB, records []T) {
	t.Helper()

	for i := range records {
		err := db.Create(&records[i]).Error
		require.NoError(t, err, "DBSeedSlice: failed to create record at index %d", i)
	}
}
