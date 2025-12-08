package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/pablotartak/Parser-Ofert/services"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("index.html")
	})

	app.Post("/parser", func(c *fiber.Ctx) error {
		raw := strings.TrimSpace(string(c.Body()))
		if raw == "" {
			return c.Status(400).SendString("Brak danych w body POST")
		}

		lines := strings.Split(raw, "\n")

		var paramLines []string
		var mediaLine, infoLine string
		var opisLines []string
		mode := "parametry"

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			switch line {
			case "Media":
				mode = "media"
				continue
			case "Informacje dodatkowe":
				mode = "info"
				continue
			case "Opis":
				mode = "opis"
				continue
			}

			switch mode {
			case "parametry":
				paramLines = append(paramLines, line)
			case "media":
				if mediaLine != "" {
					mediaLine += " "
				}
				mediaLine += line
			case "info":
				if infoLine != "" {
					infoLine += " "
				}
				infoLine += line
			case "opis":
				opisLines = append(opisLines, line)
			}
		}

		parametry, err := services.ParseParametry(paramLines)
		if err != nil {
			return c.Status(400).SendString(err.Error())
		}

		media, err := services.ParseMedia(mediaLine)
		if err != nil {
			return c.Status(400).SendString(err.Error())
		}

		info, err := services.ParseInfo(infoLine)
		if err != nil {
			return c.Status(400).SendString(err.Error())
		}

		opis, endParams, err := services.ParseOpis(opisLines)
		if err != nil {
			return c.Status(400).SendString(err.Error())
		}

		final := map[string]interface{}{
			"Parametry oferty":     parametry,
			"Media":                media,
			"Informacje dodatkowe": info,
			"Opis":                 opis,
		}
		for k, v := range endParams {
			final[k] = v
		}

		jsonData, err := json.MarshalIndent(final, "", "  ")
		if err != nil {
			return c.Status(500).SendString(fmt.Sprintf("Błąd JSON: %v", err))
		}

		return c.Send(jsonData)
	})

	app.Listen(":8080")
}
