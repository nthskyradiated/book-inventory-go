package main

import (
	"os"
	"github.com/nthskyradiated/book-inventory-go/common"
	"github.com/nthskyradiated/book-inventory-go/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main () {
	err := run()

	if err != nil {
		panic(err)
	}
}

func run ()  error {
	//initialize env
	err := common.LoadEnv()
	if err != nil {
		return err
	}

	//initialize DB
	err = common.InitDB()
	if err != nil {
		return err
	}

	//defer closing DB
	defer common.CloseDB()

	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	//add routes
	router.AddBookGroup(app)

	//start server
	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "8080"
	}
	app.Listen(":" + port)
	return nil

}