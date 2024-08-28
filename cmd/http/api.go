package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/bitzero/gostarter/internal/program"
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
	projectName := "github.com/builder"

	tempFolderPath := "/mnt/d/source/go/src/github.com/learn_stuff/go-starter/tmp/gostarter"
	defaultDirStructPath := "/mnt/d/source/go/src/github.com/learn_stuff/go-starter/assets/default_structure.yaml"
	requestID := uuid.New().String()
	fullpath := tempFolderPath + "/" + requestID
	fmt.Println("Fullpath:", fullpath)

	project, err := program.CreateProject(projectName, fullpath, defaultDirStructPath)
	if err != nil {
		slog.Error("Failed to create project", err)
	}

	project.SetDBDriver("mysql")

	err = project.CreateMainFile()
	if err != nil {
		slog.Error("Failed to create main file", err)
		return ""
	}

	return fullpath + "/generated.zip"
}
