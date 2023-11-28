package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Product struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Category string `json:"category"`
}

func main() {
	dsn := "root:root@tcp(localhost:3306)/belajar-go?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database")
	}

	db.AutoMigrate(&Product{})

	app := fiber.New()

	// Create a product
	app.Post("/products", func(c *fiber.Ctx) error {
		var product Product
		if err := c.BodyParser(&product); err != nil {
			return err
		}
		db.Create(&product)
		return c.JSON(product)
	})

	// Read all products
	app.Get("/products", func(c *fiber.Ctx) error {
		var products []Product
		if err := db.Find(&products).Error; err != nil {
			return err
		}
		return c.JSON(products)
	})

	// Read a single product
	app.Get("/products/:id", func(c *fiber.Ctx) error {
		var product Product
		id := c.Params("id")
		if err := db.First(&product, id).Error; err != nil {
			return err
		}
		return c.JSON(product)
	})

	// Update a product
	app.Put("/products/:id", func(c *fiber.Ctx) error {
		var product Product
		id := c.Params("id")
		if err := db.First(&product, id).Error; err != nil {
			return err
		}
		if err := c.BodyParser(&product); err != nil {
			return err
		}
		db.Save(&product)
		return c.JSON(product)
	})

	// Delete a product
	app.Delete("/products/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if err := db.Delete(&Product{}, id).Error; err != nil {
			return err
		}
		return c.SendStatus(fiber.StatusNoContent)
	})

	port := "3000"
	app.Listen(":" + port)
}
