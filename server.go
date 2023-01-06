package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/frngrit/assessment/db"
	"github.com/frngrit/assessment/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {

	//connect to db
	db.StartDB()

	//initial server
	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	h := handler.NewApplication(db.DB)

	e.POST("/expenses", h.CreateExpenseHandler)
	e.GET("/expenses/:id", h.GetExpenseById)
	e.PUT("/expenses/:id", h.UpdateExpenseById)
	e.GET("/expenses", h.GetAllExpenses)

	go func() {
		if err := e.Start(os.Getenv("PORT")); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
