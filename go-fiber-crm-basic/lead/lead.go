package lead

import (
	"github.com/carissaayo/go-starter-projects/go-fiber-basic/database"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Lead struct {
	gorm.Model
	Name    string `json:"name"`
	Company string `json:"company"`
	Email   string `json:"email"`
	Phone   int    `json:"phone"`
}

func GetLeads(c *fiber.Ctx) error {
	var leads []Lead
	database.DBConn.Find(&leads)
	return c.JSON(leads)
}

func GetLead(c *fiber.Ctx) error {
	id := c.Params("id")
	var lead Lead
	database.DBConn.First(&lead, id)
	return c.JSON(lead)
}

func NewLead(c *fiber.Ctx) error {
	var lead Lead

	if err := c.BodyParser(&lead); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	database.DBConn.Create(&lead)
	return c.JSON(lead)
}

func DeleteLead(c *fiber.Ctx) error {
	id := c.Params("id")

	var lead Lead
	result := database.DBConn.First(&lead, id)

	if result.Error != nil {
		return c.Status(404).SendString("No lead found with that ID")
	}

	database.DBConn.Delete(&lead)
	return c.SendString("Lead successfully deleted")
}
