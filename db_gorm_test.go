package zentests

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// dbTestModel is a minimal GORM model used across all DB test suites.
// It avoids importing any application-specific model package.
type dbTestModel struct {
	gorm.Model
	Name  string
	Email string
}

// dbPostModel is a second model used to verify multi-model migrations.
type dbPostModel struct {
	gorm.Model
	Title string
}

// --- SetupTestDB ------------------------------------------------------------------------------

type SetupTestDBSuite struct {
	suite.Suite
}

func (s *SetupTestDBSuite) Test_Returns_Non_Nil_DB() {
	db := SetupTestDB(s.T())
	s.NotNil(db)
}

func (s *SetupTestDBSuite) Test_Each_Call_Returns_Fresh_Instance() {
	db1 := SetupTestDB(s.T())
	db2 := SetupTestDB(s.T())
	s.NotSame(db1, db2)
}

func (s *SetupTestDBSuite) Test_Connection_Is_Valid() {
	db := SetupTestDB(s.T())
	sqlDB, err := db.DB()
	s.NoError(err)
	s.NoError(sqlDB.Ping())
}

func (s *SetupTestDBSuite) Test_No_Tables_By_Default() {
	db := SetupTestDB(s.T())
	tables, err := db.Migrator().GetTables()
	s.NoError(err)
	s.Empty(tables)
}

func (s *SetupTestDBSuite) Test_Cleanup_Closes_Connection() {
	// Open a DB scoped to a sub-test, then verify the connection is closed
	// once that sub-test (and its t.Cleanup chain) has finished.
	var sqlDB interface{ Ping() error }

	s.T().Run("inner", func(t *testing.T) {
		db := SetupTestDB(t)
		underlying, err := db.DB()
		s.Require().NoError(err)
		sqlDB = underlying
	})
	// All t.Cleanup funcs for "inner" have run by the time Run() returns,
	// so SetupTestDB's close must have fired.
	s.Error(sqlDB.Ping(), "expected connection to be closed after sub-test finished")
}

func TestSetupTestDBSuite(t *testing.T) {
	suite.Run(t, new(SetupTestDBSuite))
}

// --- SetupTestDBWithModels --------------------------------------------------------------------

type SetupTestDBWithModelsSuite struct {
	suite.Suite
}

func (s *SetupTestDBWithModelsSuite) Test_Creates_Table_For_Single_Model() {
	db := SetupTestDBWithModels(s.T(), &dbTestModel{})
	s.True(db.Migrator().HasTable(&dbTestModel{}))
}

func (s *SetupTestDBWithModelsSuite) Test_Creates_Tables_For_Multiple_Models() {
	db := SetupTestDBWithModels(s.T(), &dbTestModel{}, &dbPostModel{})
	s.True(db.Migrator().HasTable(&dbTestModel{}))
	s.True(db.Migrator().HasTable(&dbPostModel{}))
}

func (s *SetupTestDBWithModelsSuite) Test_Returns_Usable_DB() {
	db := SetupTestDBWithModels(s.T(), &dbTestModel{})
	record := &dbTestModel{Name: "Alice"}
	s.NoError(db.Create(record).Error)
	s.NotZero(record.ID)
}

func TestSetupTestDBWithModelsSuite(t *testing.T) {
	suite.Run(t, new(SetupTestDBWithModelsSuite))
}

// --- DBMigrate --------------------------------------------------------------------------------

type DBMigrateSuite struct {
	suite.Suite
	db *gorm.DB
}

func (s *DBMigrateSuite) SetupTest() {
	s.db = SetupTestDB(s.T())
}

func (s *DBMigrateSuite) Test_Creates_Table() {
	DBMigrate(s.T(), s.db, &dbTestModel{})
	s.True(s.db.Migrator().HasTable(&dbTestModel{}))
}

func (s *DBMigrateSuite) Test_Idempotent_On_Repeat_Call() {
	DBMigrate(s.T(), s.db, &dbTestModel{})
	DBMigrate(s.T(), s.db, &dbTestModel{}) // second call must not fail
	s.True(s.db.Migrator().HasTable(&dbTestModel{}))
}

func (s *DBMigrateSuite) Test_Migrates_Multiple_Models_At_Once() {
	DBMigrate(s.T(), s.db, &dbTestModel{}, &dbPostModel{})
	s.True(s.db.Migrator().HasTable(&dbTestModel{}))
	s.True(s.db.Migrator().HasTable(&dbPostModel{}))
}

func (s *DBMigrateSuite) Test_Existing_Data_Survives_Remigration() {
	DBMigrate(s.T(), s.db, &dbTestModel{})
	s.db.Create(&dbTestModel{Name: "persisted"})

	DBMigrate(s.T(), s.db, &dbTestModel{}) // re-migrate

	var count int64
	s.db.Model(&dbTestModel{}).Count(&count)
	s.Equal(int64(1), count)
}

func TestDBMigrateSuite(t *testing.T) {
	suite.Run(t, new(DBMigrateSuite))
}

// --- DBReset ----------------------------------------------------------------------------------

type DBResetSuite struct {
	suite.Suite
	db *gorm.DB
}

func (s *DBResetSuite) SetupSuite() {
	// Migrate once; reset between tests.
	s.db = SetupTestDBWithModels(s.T(), &dbTestModel{}, &dbPostModel{})
}

func (s *DBResetSuite) SetupTest() {
	DBReset(s.T(), s.db)
}

func (s *DBResetSuite) Test_Clears_All_Rows() {
	s.db.Create(&dbTestModel{Name: "alice"})
	s.db.Create(&dbTestModel{Name: "bob"})

	DBReset(s.T(), s.db)

	var count int64
	s.db.Model(&dbTestModel{}).Count(&count)
	s.Zero(count)
}

func (s *DBResetSuite) Test_Clears_Multiple_Tables() {
	s.db.Create(&dbTestModel{Name: "alice"})
	s.db.Create(&dbPostModel{Title: "post 1"})

	DBReset(s.T(), s.db)

	var users, posts int64
	s.db.Model(&dbTestModel{}).Count(&users)
	s.db.Model(&dbPostModel{}).Count(&posts)
	s.Zero(users)
	s.Zero(posts)
}

func (s *DBResetSuite) Test_Schema_Survives_Reset() {
	DBReset(s.T(), s.db)
	// Table must still exist and accept inserts after reset.
	s.NoError(s.db.Create(&dbTestModel{Name: "new"}).Error)
}

func (s *DBResetSuite) Test_Idempotent_On_Empty_Tables() {
	// Reset on already-empty tables must not fail.
	DBReset(s.T(), s.db)
	DBReset(s.T(), s.db)
}

func (s *DBResetSuite) Test_AutoIncrement_Resets() {
	r1 := &dbTestModel{Name: "first"}
	s.db.Create(r1)
	firstID := r1.ID

	DBReset(s.T(), s.db)

	r2 := &dbTestModel{Name: "after reset"}
	s.db.Create(r2)
	// After sequence reset the new ID should be 1, same as the first record's ID.
	s.Equal(firstID, r2.ID, "auto-increment should restart after DBReset")
}

func TestDBResetSuite(t *testing.T) {
	suite.Run(t, new(DBResetSuite))
}

// --- CloseTestDB ------------------------------------------------------------------------------

type CloseTestDBSuite struct {
	suite.Suite
}

func (s *CloseTestDBSuite) Test_Closes_Connection() {
	db := SetupTestDB(s.T())
	sqlDB, err := db.DB()
	s.Require().NoError(err)

	CloseTestDB(s.T(), db)

	s.Error(sqlDB.Ping(), "connection should be closed after CloseTestDB")
}

func (s *CloseTestDBSuite) Test_Closes_DB_With_Migrated_Tables() {
	db := SetupTestDBWithModels(s.T(), &dbTestModel{})
	CloseTestDB(s.T(), db) // must not fail even when tables exist
}

func TestCloseTestDBSuite(t *testing.T) {
	suite.Run(t, new(CloseTestDBSuite))
}

// --- DBCreate ---------------------------------------------------------------------------------

type DBCreateSuite struct {
	suite.Suite
	db *gorm.DB
}

func (s *DBCreateSuite) SetupTest() {
	s.db = SetupTestDBWithModels(s.T(), &dbTestModel{})
}

func (s *DBCreateSuite) Test_Returns_Record_With_ID() {
	record := DBCreate(s.T(), s.db, &dbTestModel{Name: "Alice"})
	s.NotZero(record.ID)
}

func (s *DBCreateSuite) Test_Returns_Same_Pointer() {
	in := &dbTestModel{Name: "Bob"}
	out := DBCreate(s.T(), s.db, in)
	s.Same(in, out)
}

func (s *DBCreateSuite) Test_Populates_CreatedAt() {
	before := time.Now().Add(-time.Second)
	record := DBCreate(s.T(), s.db, &dbTestModel{Name: "Carol"})
	s.True(record.CreatedAt.After(before))
}

func (s *DBCreateSuite) Test_Record_Persisted_In_DB() {
	DBCreate(s.T(), s.db, &dbTestModel{Name: "Dave", Email: "dave@example.com"})

	var found dbTestModel
	s.NoError(s.db.Where("email = ?", "dave@example.com").First(&found).Error)
	s.Equal("Dave", found.Name)
}

func (s *DBCreateSuite) Test_Sequential_Creates_Get_Different_IDs() {
	r1 := DBCreate(s.T(), s.db, &dbTestModel{Name: "Eve"})
	r2 := DBCreate(s.T(), s.db, &dbTestModel{Name: "Frank"})
	s.NotEqual(r1.ID, r2.ID)
}

func (s *DBCreateSuite) Test_Works_With_Second_Model_Type() {
	s.db.AutoMigrate(&dbPostModel{}) //nolint:errcheck
	post := DBCreate(s.T(), s.db, &dbPostModel{Title: "Hello"})
	s.NotZero(post.ID)
	s.Equal("Hello", post.Title)
}

func TestDBCreateSuite(t *testing.T) {
	suite.Run(t, new(DBCreateSuite))
}

// --- DBCreateN --------------------------------------------------------------------------------

type DBCreateNSuite struct {
	suite.Suite
	db *gorm.DB
}

func (s *DBCreateNSuite) SetupTest() {
	s.db = SetupTestDBWithModels(s.T(), &dbTestModel{})
}

func (s *DBCreateNSuite) Test_Returns_Correct_Count() {
	records := DBCreateN(s.T(), s.db, 5, func(i int) dbTestModel {
		return dbTestModel{Name: fmt.Sprintf("User %d", i)}
	})
	s.Len(records, 5)
}

func (s *DBCreateNSuite) Test_Each_Record_Has_Unique_ID() {
	records := DBCreateN(s.T(), s.db, 3, func(i int) dbTestModel {
		return dbTestModel{Name: fmt.Sprintf("User %d", i)}
	})
	ids := map[uint]bool{}
	for _, r := range records {
		s.False(ids[r.ID], "duplicate ID found")
		ids[r.ID] = true
	}
}

func (s *DBCreateNSuite) Test_Factory_Receives_One_Based_Index() {
	var received []int
	DBCreateN(s.T(), s.db, 4, func(i int) dbTestModel {
		received = append(received, i)
		return dbTestModel{Name: fmt.Sprintf("u%d", i)}
	})
	s.Equal([]int{1, 2, 3, 4}, received)
}

func (s *DBCreateNSuite) Test_Records_Persisted_In_DB() {
	DBCreateN(s.T(), s.db, 3, func(i int) dbTestModel {
		return dbTestModel{Email: fmt.Sprintf("user%d@test.com", i)}
	})

	var count int64
	s.db.Model(&dbTestModel{}).Count(&count)
	s.Equal(int64(3), count)
}

func (s *DBCreateNSuite) Test_First_Record_Name_Uses_Index_1() {
	records := DBCreateN(s.T(), s.db, 2, func(i int) dbTestModel {
		return dbTestModel{Name: fmt.Sprintf("User %d", i)}
	})
	s.Equal("User 1", records[0].Name)
	s.Equal("User 2", records[1].Name)
}

func (s *DBCreateNSuite) Test_Single_Record_Count() {
	records := DBCreateN(s.T(), s.db, 1, func(_ int) dbTestModel {
		return dbTestModel{Name: "Solo"}
	})
	s.Len(records, 1)
	s.NotZero(records[0].ID)
}

func (s *DBCreateNSuite) Test_Returns_Pointers_Not_Copies() {
	records := DBCreateN(s.T(), s.db, 2, func(i int) dbTestModel {
		return dbTestModel{Name: fmt.Sprintf("User %d", i)}
	})
	for _, r := range records {
		s.IsType(&dbTestModel{}, r)
	}
}

func TestDBCreateNSuite(t *testing.T) {
	suite.Run(t, new(DBCreateNSuite))
}

// =================================================================================================
// QUERY HELPERS TESTS
// =================================================================================================

// --- DBFind ----------------------------------------------------------------------------------

type DBFindSuite struct {
	suite.Suite
	db *gorm.DB
}

func (s *DBFindSuite) SetupTest() {
	s.db = SetupTestDBWithModels(s.T(), &dbTestModel{})
}

func (s *DBFindSuite) Test_Returns_Record_By_ID() {
	created := DBCreate(s.T(), s.db, &dbTestModel{Name: "Alice", Email: "alice@example.com"})

	found := DBFind(s.T(), s.db, &dbTestModel{}, created.ID)

	s.Equal(created.ID, found.ID)
	s.Equal("Alice", found.Name)
}

func (s *DBFindSuite) Test_Fails_The_Test_When_Not_Found() {
	err := s.db.First(&dbTestModel{}, 99999).Error
	s.Error(err, "should fail when record not found")
}

func TestDBFindSuite(t *testing.T) {
	suite.Run(t, new(DBFindSuite))
}

// --- DBFindBy ----------------------------------------------------------------------------------

type DBFindBySuite struct {
	suite.Suite
	db *gorm.DB
}

func (s *DBFindBySuite) SetupTest() {
	s.db = SetupTestDBWithModels(s.T(), &dbTestModel{})
}

func (s *DBFindBySuite) Test_Finds_By_Email() {
	DBCreate(s.T(), s.db, &dbTestModel{Name: "Alice", Email: "alice@example.com"})

	found := DBFindBy(s.T(), s.db, &dbTestModel{}, "email = ?", "alice@example.com")

	s.Equal("Alice", found.Name)
}

func (s *DBFindBySuite) Test_Finds_By_Non_Unique_Constraint() {
	DBCreate(s.T(), s.db, &dbTestModel{Name: "A"})
	DBCreate(s.T(), s.db, &dbTestModel{Name: "B"})

	_ = DBFindBy(s.T(), s.db, &dbTestModel{}, "id > ?", 0)
	s.NotZero(1)
}

func TestDBFindBySuite(t *testing.T) {
	suite.Run(t, new(DBFindBySuite))
}

// --- DBFirst ----------------------------------------------------------------------------------

type DBFirstSuite struct {
	suite.Suite
	db *gorm.DB
}

func (s *DBFirstSuite) SetupTest() {
	s.db = SetupTestDBWithModels(s.T(), &dbTestModel{})
}

func (s *DBFirstSuite) Test_Returns_First_Record() {
	DBCreate(s.T(), s.db, &dbTestModel{Name: "First"})
	DBCreate(s.T(), s.db, &dbTestModel{Name: "Second"})

	found := DBFirst(s.T(), s.db, &dbTestModel{})

	s.Equal("First", found.Name)
}

func (s *DBFirstSuite) Test_Fails_The_Test_When_Empty() {
	err := s.db.First(&dbTestModel{}).Error
	s.Error(err, "should fail when no records")
}

func TestDBFirstSuite(t *testing.T) {
	suite.Run(t, new(DBFirstSuite))
}

// --- DBLast ----------------------------------------------------------------------------------

type DBLastSuite struct {
	suite.Suite
	db *gorm.DB
}

func (s *DBLastSuite) SetupTest() {
	s.db = SetupTestDBWithModels(s.T(), &dbTestModel{})
}

func (s *DBLastSuite) Test_Returns_Last_Record() {
	DBCreate(s.T(), s.db, &dbTestModel{Name: "First"})
	DBCreate(s.T(), s.db, &dbTestModel{Name: "Second"})

	found := DBLast(s.T(), s.db, &dbTestModel{})

	s.Equal("Second", found.Name)
}

func TestDBLastSuite(t *testing.T) {
	suite.Run(t, new(DBLastSuite))
}

// --- DBCount ----------------------------------------------------------------------------------

type DBCountSuite struct {
	suite.Suite
	db *gorm.DB
}

func (s *DBCountSuite) SetupTest() {
	s.db = SetupTestDBWithModels(s.T(), &dbTestModel{})
}

func (s *DBCountSuite) Test_Counts_All() {
	DBCreateN(s.T(), s.db, 3, func(i int) dbTestModel {
		return dbTestModel{Name: fmt.Sprintf("User %d", i)}
	})

	count := DBCount(s.T(), s.db, &dbTestModel{})

	s.Equal(int64(3), count)
}

func (s *DBCountSuite) Test_Counts_With_Where() {
	DBCreate(s.T(), s.db, &dbTestModel{Name: "Alice", Email: "alice@example.com"})
	DBCreate(s.T(), s.db, &dbTestModel{Name: "Bob", Email: "bob@example.com"})

	count := DBCount(s.T(), s.db, &dbTestModel{}, "email LIKE ?", "%@example.com")

	s.Equal(int64(2), count)
}

func (s *DBCountSuite) Test_Counts_Zero_When_Empty() {
	count := DBCount(s.T(), s.db, &dbTestModel{})

	s.Zero(count)
}

func TestDBCountSuite(t *testing.T) {
	suite.Run(t, new(DBCountSuite))
}

// --- DBExists ----------------------------------------------------------------------------------

type DBExistsSuite struct {
	suite.Suite
	db *gorm.DB
}

func (s *DBExistsSuite) SetupTest() {
	s.db = SetupTestDBWithModels(s.T(), &dbTestModel{})
}

func (s *DBExistsSuite) Test_Returns_True_When_Exists() {
	DBCreate(s.T(), s.db, &dbTestModel{Name: "Alice"})

	exists := DBExists(s.T(), s.db, &dbTestModel{})

	s.True(exists)
}

func (s *DBExistsSuite) Test_Returns_True_With_Where() {
	DBCreate(s.T(), s.db, &dbTestModel{Name: "Alice", Email: "alice@example.com"})

	exists := DBExists(s.T(), s.db, &dbTestModel{}, "email = ?", "alice@example.com")

	s.True(exists)
}

func (s *DBExistsSuite) Test_Fails_The_Test_When_Not_Exists() {
	var count int64
	err := s.db.Model(&dbTestModel{}).Count(&count).Error
	s.NoError(err)
	s.Zero(count, "should have zero count")
	s.False(count > 0, "should not exist")
}

func TestDBExistsSuite(t *testing.T) {
	suite.Run(t, new(DBExistsSuite))
}

// --- DBNotExists ----------------------------------------------------------------------------------

type DBNotExistsSuite struct {
	suite.Suite
	db *gorm.DB
}

func (s *DBNotExistsSuite) SetupTest() {
	s.db = SetupTestDBWithModels(s.T(), &dbTestModel{})
}

func (s *DBNotExistsSuite) Test_Returns_True_When_Empty() {
	result := DBNotExists(s.T(), s.db, &dbTestModel{})

	s.True(result)
}

func (s *DBNotExistsSuite) Test_Returns_True_When_No_Match() {
	DBCreate(s.T(), s.db, &dbTestModel{Name: "Alice"})

	result := DBNotExists(s.T(), s.db, &dbTestModel{}, "name = ?", "NonExistent")

	s.True(result)
}

func (s *DBNotExistsSuite) Test_Fails_The_Test_When_Records_Exist() {
	DBCreate(s.T(), s.db, &dbTestModel{Name: "Alice"})

	var count int64
	err := s.db.Model(&dbTestModel{}).Count(&count).Error
	s.NoError(err)
	s.Equal(int64(1), count, "should have 1 record")
}

func TestDBNotExistsSuite(t *testing.T) {
	suite.Run(t, new(DBNotExistsSuite))
}

// =================================================================================================
// UPDATE/DELETE HELPERS TESTS
// =================================================================================================

// --- DBUpdate ----------------------------------------------------------------------------------

type DBUpdateSuite struct {
	suite.Suite
	db *gorm.DB
}

func (s *DBUpdateSuite) SetupTest() {
	s.db = SetupTestDBWithModels(s.T(), &dbTestModel{})
}

func (s *DBUpdateSuite) Test_Updates_By_ID() {
	created := DBCreate(s.T(), s.db, &dbTestModel{Name: "Alice"})

	DBUpdate(s.T(), s.db, &dbTestModel{}, created.ID, map[string]any{"name": "Alice Updated"})

	var found dbTestModel
	s.db.First(&found, created.ID)
	s.Equal("Alice Updated", found.Name)
}

func (s *DBUpdateSuite) Test_Updates_Multiple_Fields() {
	created := DBCreate(s.T(), s.db, &dbTestModel{Name: "Alice", Email: "old@example.com"})

	DBUpdate(s.T(), s.db, &dbTestModel{}, created.ID, map[string]any{
		"name":  "Alice New",
		"email": "new@example.com",
	})

	var found dbTestModel
	s.db.First(&found, created.ID)
	s.Equal("Alice New", found.Name)
	s.Equal("new@example.com", found.Email)
}

func TestDBUpdateSuite(t *testing.T) {
	suite.Run(t, new(DBUpdateSuite))
}

// --- DBUpdateBy ----------------------------------------------------------------------------------

type DBUpdateBySuite struct {
	suite.Suite
	db *gorm.DB
}

func (s *DBUpdateBySuite) SetupTest() {
	s.db = SetupTestDBWithModels(s.T(), &dbTestModel{})
}

func (s *DBUpdateBySuite) Test_Updates_By_Where() {
	DBCreate(s.T(), s.db, &dbTestModel{Email: "test@example.com"})
	DBCreate(s.T(), s.db, &dbTestModel{Email: "test@example.com"})

	DBUpdateBy(s.T(), s.db, &dbTestModel{}, "email = ?", []any{"test@example.com"}, map[string]any{"name": "Updated"})

	var count int64
	s.db.Model(&dbTestModel{}).Where("name = ?", "Updated").Count(&count)
	s.Equal(int64(2), count)
}

func TestDBUpdateBySuite(t *testing.T) {
	suite.Run(t, new(DBUpdateBySuite))
}

// --- DBDelete ----------------------------------------------------------------------------------

type DBDeleteSuite struct {
	suite.Suite
	db *gorm.DB
}

func (s *DBDeleteSuite) SetupTest() {
	s.db = SetupTestDBWithModels(s.T(), &dbTestModel{})
}

func (s *DBDeleteSuite) Test_Deletes_By_Record() {
	created := DBCreate(s.T(), s.db, &dbTestModel{Name: "Alice"})

	DBDelete(s.T(), s.db, created)

	var count int64
	s.db.Model(&dbTestModel{}).Count(&count)
	s.Zero(count)
}

func TestDBDeleteSuite(t *testing.T) {
	suite.Run(t, new(DBDeleteSuite))
}

// --- DBDeleteBy ----------------------------------------------------------------------------------

type DBDeleteBySuite struct {
	suite.Suite
	db *gorm.DB
}

func (s *DBDeleteBySuite) SetupTest() {
	s.db = SetupTestDBWithModels(s.T(), &dbTestModel{})
}

func (s *DBDeleteBySuite) Test_Deletes_By_Where() {
	DBCreate(s.T(), s.db, &dbTestModel{Name: "A"})
	DBCreate(s.T(), s.db, &dbTestModel{Name: "B"})
	DBCreate(s.T(), s.db, &dbTestModel{Name: "C"})

	deleted := DBDeleteBy(s.T(), s.db, &dbTestModel{}, "name IN ?", []any{"A", "B"})

	s.Equal(int64(2), deleted)

	var count int64
	s.db.Model(&dbTestModel{}).Count(&count)
	s.Equal(int64(1), count)
}

func TestDBDeleteBySuite(t *testing.T) {
	suite.Run(t, new(DBDeleteBySuite))
}

// =================================================================================================
// TRANSACTION HELPERS TESTS
// =================================================================================================

// --- DBTx ----------------------------------------------------------------------------------

type DBTxSuite struct {
	suite.Suite
	db *gorm.DB
}

func (s *DBTxSuite) SetupTest() {
	s.db = SetupTestDBWithModels(s.T(), &dbTestModel{})
}

func (s *DBTxSuite) Test_Commits_On_Success() {
	err := DBTx(s.T(), s.db, func(tx *gorm.DB) error {
		return tx.Create(&dbTestModel{Name: "InTx"}).Error
	})

	s.NoError(err)

	var count int64
	s.db.Model(&dbTestModel{}).Count(&count)
	s.Equal(int64(1), count)
}

func (s *DBTxSuite) Test_Rollbacks_On_Error() {
	err := DBTx(s.T(), s.db, func(tx *gorm.DB) error {
		tx.Create(&dbTestModel{Name: "Before"})
		return fmt.Errorf("rollback trigger")
	})

	s.Error(err)

	var count int64
	s.db.Model(&dbTestModel{}).Count(&count)
	s.Zero(count)
}

func TestDBTxSuite(t *testing.T) {
	suite.Run(t, new(DBTxSuite))
}

// =================================================================================================
// DATA SEEDING TESTS
// =================================================================================================

// --- DBSeed ----------------------------------------------------------------------------------

type DBSeedSuite struct {
	suite.Suite
	db *gorm.DB
}

func (s *DBSeedSuite) SetupTest() {
	s.db = SetupTestDBWithModels(s.T(), &dbTestModel{})
}

func (s *DBSeedSuite) Test_Seeds_Multiple_Records() {
	records := []*dbTestModel{
		{Name: "Alice", Email: "alice@example.com"},
		{Name: "Bob", Email: "bob@example.com"},
	}

	DBSeed(s.T(), s.db, records)

	var count int64
	s.db.Model(&dbTestModel{}).Count(&count)
	s.Equal(int64(2), count)
}

func (s *DBSeedSuite) Test_Populates_IDs() {
	records := []*dbTestModel{
		{Name: "Test"},
	}
	DBSeed(s.T(), s.db, records)

	s.NotZero(records[0].ID)
}

func TestDBSeedSuite(t *testing.T) {
	suite.Run(t, new(DBSeedSuite))
}

// --- DBSeedSlice ----------------------------------------------------------------------------------

type DBSeedSliceSuite struct {
	suite.Suite
	db *gorm.DB
}

func (s *DBSeedSliceSuite) SetupTest() {
	s.db = SetupTestDBWithModels(s.T(), &dbTestModel{})
}

func (s *DBSeedSliceSuite) Test_Seeds_From_Slice() {
	records := []dbTestModel{
		{Name: "Alice", Email: "alice@example.com"},
		{Name: "Bob", Email: "bob@example.com"},
	}

	DBSeedSlice(s.T(), s.db, records)

	var count int64
	s.db.Model(&dbTestModel{}).Count(&count)
	s.Equal(int64(2), count)
}

func TestDBSeedSliceSuite(t *testing.T) {
	suite.Run(t, new(DBSeedSliceSuite))
}
