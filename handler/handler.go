package handler

import (
	"fmt"
	"net/http"

	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

type Err struct {
	Message string `json:"message"`
}

type handler struct {
	DB *sql.DB
}

func NewApplication(db *sql.DB) *handler {
	return &handler{db}
}

type Expense struct {
	Id     int      `json:"id"`
	Title  string   `json:"title"`
	Amount float64  `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

func (h handler) CreateExpenseHandler(c echo.Context) error {
	expense := Expense{}
	c.Bind(&expense)

	//check if user provide body parameter correctly
	if expense.Title == "" || expense.Note == "" || expense.Amount == 0 || expense.Tags == nil {
		return c.JSON(http.StatusBadRequest, Err{"bad body request"})
	}

	row := h.DB.QueryRow("INSERT INTO expenses (title, amount, note, tags) VALUES ($1, $2, $3, $4) RETURNING id", expense.Title, expense.Amount, expense.Note, pq.Array(expense.Tags))
	var id int
	if err := row.Scan(&id); err != nil {
		fmt.Printf("\n\n%v\n\n", err.Error())
		return c.JSON(http.StatusInternalServerError, Err{err.Error()})
	}
	expense.Id = id
	return c.JSON(http.StatusCreated, expense)
}

func (h handler) GetExpenseById(c echo.Context) error {
	expense := Expense{}

	stmt, err := h.DB.Prepare("SELECT * FROM expenses WHERE id=$1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{err.Error()})
	}

	rowId := c.Param("id")
	rows := stmt.QueryRow(rowId)

	err = rows.Scan(&expense.Id, &expense.Title, &expense.Amount, &expense.Note, pq.Array(&expense.Tags))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{err.Error()})
	}

	return c.JSON(http.StatusOK, expense)
}

func (h handler) UpdateExpenseById(c echo.Context) error {
	expense := Expense{}
	c.Bind(&expense)
	rowId := c.Param("id")

	//check if user provide body parameter correctly
	if expense.Title == "" || expense.Note == "" || expense.Amount == 0 || expense.Tags == nil {
		return c.JSON(http.StatusBadRequest, Err{"bad body request"})
	}

	stmt, err := h.DB.Prepare("UPDATE expenses SET title=$1, amount=$2, note=$3, tags=$4 WHERE id=$5;")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{err.Error()})
	}

	if _, err := stmt.Exec(&expense.Title, &expense.Amount, &expense.Note, pq.Array(&expense.Tags), rowId); err != nil {
		return c.JSON(http.StatusInternalServerError, Err{err.Error()})
	}

	return c.JSON(http.StatusOK, expense)
}
