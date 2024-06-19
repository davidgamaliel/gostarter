package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bitzero/gostarter/pkg/filestructure"
	"github.com/bitzero/gostarter/pkg/gomodule"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/google/uuid"
)

func main() {
	// GenerateBoilerPlate()

	engine := html.New("./web", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/healthz", healthz)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})
	app.Get("/generate", generate)

	if err := app.Listen(":5555"); err != nil {
		log.Fatal("Failed to run service")
	}
	log.Print("Service is running on port 5555")
}

func healthz(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "ok"})
}

func generate(c *fiber.Ctx) error {
	filepath := GenerateBoilerPlate()
	return c.Status(http.StatusOK).SendFile(filepath)
	// return c.SendString("Boilerplate generated successfully")

}

func GenerateBoilerPlate() string {
	tempFolderPath := "./tmp/gostarter"
	requestID := uuid.New().String()
	fullpath := tempFolderPath + "/" + requestID
	fmt.Println("Fullpath:", fullpath)

	err := filestructure.CreateFolderStructure(fullpath)
	if err != nil {
		fmt.Println("Failed to create folder structure:", err)
		return ""
	}

	err = gomodule.GenerateGoMod(fullpath)
	if err != nil {
		fmt.Println("Failed to generate go.mod file:", err)
		return ""
	}

	err = filestructure.WrapToZip(fullpath)
	if err != nil {
		fmt.Println("Failed to wrap files and folders to zip:", err)
		return ""
	}

	fmt.Println("Successfully generated go.mod file and wrapped to zip.")
	go func() {
		time.Sleep(10 * time.Second)
		filestructure.CleanUp(fullpath)
	}()

	return fullpath + "/generated.zip"
}
