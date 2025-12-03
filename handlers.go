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

	// Serwowanie index.html
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("index.html")
	})

	// Endpoint POST /parser
	app.Post("/parser", func(c *fiber.Ctx) error {
		raw := string(c.Body())
		if raw == "" {
			return fmt.Errorf("Brak danych w body POST")
		}

		lines := strings.Split(raw, "\n")

		// zmienne lokalne do przechowywania linii dla każdej sekcji
		var paramLines []string
		var mediaLine, infoLine string
		var opisLines []string
		mode := "parametry"

		// Przełączanie sekcji według nagłówków
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

		// Wywołanie funkcji parsera z pakietu services
		parametry := services.ParseParametry(paramLines)
		media := services.ParseMedia(mediaLine)
		info := services.ParseInfo(infoLine)
		opis, endParams := services.ParseOpis(opisLines)

		// Łączenie wyników w jedną mapę
		final := map[string]interface{}{
			"Parametry oferty":     parametry,
			"Media":                media,
			"Informacje dodatkowe": info,
			"Opis":                 opis,
		}
		for k, v := range endParams {
			final[k] = v
		}

		// Konwersja do JSON
		jsonData, err := json.MarshalIndent(final, "", "  ")
		if err != nil {
			return c.Status(500).SendString(fmt.Sprintf("Błąd konwersji do JSON: %v", err))
		}

		return c.SendString(string(jsonData))
	})

	// Uruchomienie serwera na porcie 8080
	app.Listen(":8080")
}
