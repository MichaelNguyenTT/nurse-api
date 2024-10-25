package service_test

import (
	"nms/api/service"
	"nms/mock"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func TestRepositoryRequests(t *testing.T) {
	t.Parallel()

	db, mock, err := mock.LoadMockDB()
	assertNoError(t, err)

	repo := service.NewRepository(db)

	rows := sqlmock.NewRows([]string{"id", "Name", "Category"}).AddRow(uuid.New(), "Name1", "Category1").AddRow(uuid.New(), "Brain", "Neuroscience")

	mock.ExpectQuery("^SELECT (.+) FROM \"service\"").
		WillReturnRows(rows)

	res, err := repo.List()
	assertNoError(t, err)
	assertEqual(t, len(res), 2)
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
