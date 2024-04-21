package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nthskyradiated/book-inventory-go/common"
	"github.com/nthskyradiated/book-inventory-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddBookGroup(app *fiber.App) {
	bookGroup := app.Group("/books")
	bookGroup.Get("/", getBooks)
	bookGroup.Get("/:id", getBook)
	bookGroup.Post("/", addBook)
	bookGroup.Put("/:id", updateBook)
	bookGroup.Delete("/:id", deleteBook)
}

func getBooks(c *fiber.Ctx) error {
	col := common.GetDB("books")

	//get all books
	books := make([]models.Book, 0)
	cursor, err := col.Find(c.Context(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	for cursor.Next(c.Context()) {
		book := models.Book{}
		err := cursor.Decode(&book)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		books = append(books, book)
	}
	return c.Status(200).JSON(fiber.Map{"data": books})
}

func getBook(c *fiber.Ctx) error {
	col := common.GetDB("books")
	
	//find single book
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "book id not provided",
		})
	}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid id",
		})
	}
	book := models.Book{}
	err = col.FindOne(c.Context(), bson.M{"_id": objectId}).Decode(&book)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{"data" : book})
}

type createDTO struct {
	Title  string `json:"title" bson:"title"`
	Author string `json:"author" bson:"author"`
	Year   int `json:"year" bson:"year"`
}

func addBook(c *fiber.Ctx) error {
	//validate body
	b := new(createDTO)
	if err := c.BodyParser(b); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid body",
		})
	}
	//add the book
	col := common.GetDB("books")
	result, err := col.InsertOne(c.Context(), b)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to add book",
			"message": err.Error(),
		})
	}
	return c.Status(201).JSON(fiber.Map{
		"result": result,
	})
}

type updateDTO struct {
	Title  string `json:"title,omitempty" bson:"title,omitempty"`
	Author string `json:"author,omitempty" bson:"author,omitempty"`
	Year   int `json:"year,omitempty" bson:"year,omitempty"`
}

func updateBook(c *fiber.Ctx) error {
	//validate body
	b := new(updateDTO)
	if err := c.BodyParser(b); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid body",
		})
	}

	//get id
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "book id required",
		})
	}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid id",
		})
	}
	//update book
	col := common.GetDB("books")
	result, err := col.UpdateOne(c.Context(), bson.M{"_id": objectId}, bson.M{"$set": b})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to update book",
			"message": err.Error(),
		})
	}
	//return book
	return c.Status(200).JSON(fiber.Map{
		"result": result,
	})
}

func deleteBook(c *fiber.Ctx) error {
	//get id
	id := c.Params("id")
	if id == ""  {
		return c.Status(400).JSON(fiber.Map{
			"error": "book id required",
		})
	}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	//delete book
	col := common.GetDB("books")
	result, err := col.DeleteOne(c.Context(), bson.M{"_id": objectId})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to delete book",
			"message": err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"result": result,
	})
}