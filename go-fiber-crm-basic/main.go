package main

import (
	"fmt"

	"github.com/carissaayo/go-starter-projects/go-fiber-basic/database"
	"github.com/carissaayo/go-starter-projects/go-fiber-basic/lead"
	"github.com/gofiber/fiber"
	"github.com/jinzhu/gorm"
)

func setupRoutes(app *fiber.App) {
	app.Get(lead.GetLeads)
	app.Get(lead.GetLead)
	app.Post(lead.NewLead)
	app.Delete(lead.DeleteLead)
}

func initDatabase() {
	var err error
	database.DBConn, err = gorm.Open("sqlite3", "leads.db")

	if err != nil {
		panic("failed to connect to database")
	}

	fmt.Println("Connection opened to database")

	database.DBConn.AutoMigrate((&lead.Lead{}))
}

func main() {
	app := fiber.New()

	initDatabase()
	setupRoutes(app)

	app.Listen(3000)

	defer database.DBConn.Close()

}
