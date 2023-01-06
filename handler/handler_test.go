//go:build unit
// +build unit

package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreateExpense(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/expense", strings.NewReader(`{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	expenseMock := Expense{
		1,
		"strawberry smoothie",
		79,
		"night market promotion discount 10 bath",
		[]string{"food", "beverage"},
	}

	newsMockRows := sqlmock.NewRows([]string{"id"}).
		AddRow("123")

	mock.ExpectQuery("INSERT INTO expenses (.+) RETURNING id").
		WithArgs(expenseMock.Title, expenseMock.Amount, expenseMock.Note, pq.Array(expenseMock.Tags)).
		WillReturnRows(newsMockRows)
	mock.ExpectCommit()

	h := handler{db}
	c := e.NewContext(req, rec)
	expected := "{\"id\":123,\"title\":\"strawberry smoothie\",\"amount\":79,\"note\":\"night market promotion discount 10 bath\",\"tags\":[\"food\",\"beverage\"]}"

	// Act
	err = h.CreateExpenseHandler(c)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}
