package main

import (
	"fmt"

	"github.com/carissaayo/go-starter-projects/go-fiber-basic/database"
	"github.com/carissaayo/go-starter-projects/go-fiber-basic/lead"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupRoutes(app *fiber.App) {
	app.Get("/api/v1/lead", lead.GetLeads)
	app.Get("/api/v1/lead/:id", lead.GetLead)
	app.Post("/api/v1/lead", lead.NewLead)
	app.Delete("/api/v1/lead/:id", lead.DeleteLead)
}

func initDatabase() {
	var err error

	database.DBConn, err = gorm.Open(
		sqlite.Open("leads.db"),
		&gorm.Config{},
	)

	if err != nil {
		panic(err)
	}

	fmt.Println("Database connected")

	database.DBConn.AutoMigrate(&lead.Lead{})
}

func main() {
	app := fiber.New()

	initDatabase()
	setupRoutes(app)

	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}
}
