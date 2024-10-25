package service_test

import (
	"nms/api/service"
	"nms/mock"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func TestRepositoryList(t *testing.T) {
	t.Parallel()

	db, mock, err := mock.LoadMockDB()
	assertNoError(t, err)

	repo := service.NewRepository(db)

	id1 := uuid.New()
	id2 := uuid.New()
	id3 := uuid.New()

	rows := sqlmock.NewRows(setMockRow()).
		AddRow(id1, "Name1", "Category1", "notes1").
		AddRow(id2, "Brain", "Neuroscience", "notes2").
		AddRow(id3, "Brain", "Neuroscience", "notes3")

	mock.ExpectQuery(`(?i)SELECT.*FROM "services"`).WillReturnRows(rows)

	res, err := repo.List()
	assertNoError(t, err)
	assertEqual(t, len(res), 3)

	t.Logf("retrieved %d services", len(res))
	for i, s := range res {
		t.Logf("Service %d: id=%v, name=%s", i, s.ID, s.Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectation error: %s", err)
	}
}
func TestRepository_Create(t *testing.T) {
	t.Parallel()

	db, mock, err := mock.LoadMockDB()
	assertNoError(t, err)

	repo := service.NewRepository(db)

	now := time.Now()
	testID := uuid.New()

	testService := &service.Service{
		ID:        testID,
		Name:      "TestName",
		Category:  "TestCategory",
		Notes:     "TestNotes",
		CreatedAt: now,
		UpdatedAt: now,
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "services"`).
		WithArgs(testID, "TestName", "TestCategory", "TestNotes", sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	res, err := repo.Create(testService)
	assertNoError(t, err)

	assertEqual(t, res.ID, testID)
	assertEqual(t, res.Name, "TestName")
	assertEqual(t, res.Category, "TestCategory")
	assertEqual(t, res.Notes, "TestNotes")
}

func TestRepositoryRead(t *testing.T) {
	t.Parallel()

	db, mock, err := mock.LoadMockDB()
	assertNoError(t, err)

	repo := service.NewRepository(db)

	id := uuid.New()

	rows := sqlmock.NewRows(setMockRow()).
		AddRow(id, "TestName", "TestCategory", "TestNotes")

	mock.ExpectQuery(`(?i)SELECT.*FROM "services" WHERE`).
		WithArgs(id, 1).
		WillReturnRows(rows)

	service, err := repo.Read(id)
	assertNoError(t, err)
	assertEqual(t, "TestName", service.Name)
}

func TestRepositoryUpdate(t *testing.T) {
	t.Parallel()

	db, mock, err := mock.LoadMockDB()
	assertNoError(t, err)

	repo := service.NewRepository(db)

	t.Run("should update service successfully", func(t *testing.T) {
		// Create test data
		testID := uuid.MustParse("616c9ed9-5936-476b-9972-24a02ae92b6f")
		testService := &service.Service{
			ID:       testID,
			Name:     "TestName",
			Category: "TestCategory",
			Notes:    "TestNotes",
		}

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "services"`).
			WithArgs(
				"TestName",
				"TestCategory",
				"TestNotes",
				sqlmock.AnyArg(),
				testID,
			).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		rowsAffected, err := repo.Update(testService)
		assertNoError(t, err)
		assertEqual(t, rowsAffected, int64(1))

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("should handle not found case", func(t *testing.T) {
		testID := uuid.New()
		testService := &service.Service{
			ID:       testID,
			Name:     "TestName",
			Category: "TestCategory",
			Notes:    "TestNotes",
		}

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "services"`).
			WithArgs(
				"TestName",
				"TestCategory",
				"TestNotes",
				sqlmock.AnyArg(),
				testID,
			).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		rows, err := repo.Update(testService)
		assertNoError(t, err)
		assertEqual(t, rows, int64(0))

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
func TestRepositoryDelete(t *testing.T) {
	t.Parallel()

	db, mock, err := mock.LoadMockDB()
	assertNoError(t, err)

	repo := service.NewRepository(db)
	id := uuid.New()

	_ = sqlmock.NewRows(setMockRow()).
		AddRow(id, "TestName", "TestCategory", "TestNotes")

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE \"services\" SET \"deleted_at\"").
		WithArgs(sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	rows, err := repo.Delete(id)
	assertNoError(t, err)
	assertEqual(t, int64(1), rows)
}

func setMockRow() []string {
	return []string{
		"id",
		"Name",
		"Category",
		"Notes",
	}
}

func assertExpectations(t *testing.T, msg string, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: %e", msg, err)
	}
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("err: %e", err)
	}
}

func assertEqual[T comparable](t *testing.T, x, y T) {
	if x != y {
		t.Fatalf("not equal: %v, %v", x, y)
	}
}
