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
	req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(`{
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

func TestGetExpenseById(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses", strings.NewReader(``))
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

	newsMockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow(expenseMock.Id, expenseMock.Title, expenseMock.Amount, expenseMock.Note, pq.Array(expenseMock.Tags))

	mock.ExpectPrepare("SELECT")
	mock.ExpectQuery("SELECT (.+) FROM expenses").WithArgs("1").WillReturnRows(newsMockRows)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	h := handler{db}
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")
	expected := "{\"id\":1,\"title\":\"strawberry smoothie\",\"amount\":79,\"note\":\"night market promotion discount 10 bath\",\"tags\":[\"food\",\"beverage\"]}"

	// Act
	err = h.GetExpenseById(c)
	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}

func TestUpdateExpenseById(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/expenses/1", strings.NewReader(`
	{
		"id": 1,
		"title": "strawberry smoothie",
		"amount": 89,
		"note": "no discount",
		"tags": ["beverage"]
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
		89,
		"no discount",
		[]string{"beverage"},
	}

	mock.ExpectPrepare("UPDATE")
	mock.ExpectExec("UPDATE expenses").WithArgs(expenseMock.Title, expenseMock.Amount, expenseMock.Note, pq.Array(expenseMock.Tags), "1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	h := handler{db}
	c := e.NewContext(req, rec)
	c.SetParamNames("id") //set path params
	c.SetParamValues("1") //set path params
	expected := "{\"id\":1,\"title\":\"strawberry smoothie\",\"amount\":89,\"note\":\"no discount\",\"tags\":[\"beverage\"]}"

	// Act
	err = h.UpdateExpenseById(c)
	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}

func TestGetAllExpenses(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses", strings.NewReader(``))
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

	newsMockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow(expenseMock.Id, expenseMock.Title, expenseMock.Amount, expenseMock.Note, pq.Array(expenseMock.Tags))

	mock.ExpectPrepare("SELECT")
	mock.ExpectQuery("SELECT (.+) FROM expenses").WillReturnRows(newsMockRows)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	h := handler{db}
	c := e.NewContext(req, rec)
	expected := "[{\"id\":1,\"title\":\"strawberry smoothie\",\"amount\":79,\"note\":\"night market promotion discount 10 bath\",\"tags\":[\"food\",\"beverage\"]}]"

	// Act
	err = h.GetAllExpenses(c)
	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}
