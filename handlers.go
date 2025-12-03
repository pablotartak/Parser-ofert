package main
import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/pablotartak/parser-ofert/parser/services"
	"strings"
)
func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("index.html")
	})

	app.Post("/parser", func(c *fiber.Ctx) error {
		raw := string(c.Body())
		if raw == "" {
			return fmt.Errorf("Brak danych w body POST")
		}

		lines := strings.Split(raw, "\n")

		var paramLines []string
		var mediaLine, infoLine string
		var opisLines []string
		mode := "parametry"
//przełączanie sekcji według nagłówków
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
// Wywołanie parserów dla każdej sekcji
		parametry := parseParametry(paramLines)
		media := parseMedia(mediaLine)
		info := parseInfo(infoLine)
		opis, endParams := parseOpis(opisLines)

		final := map[string]interface{}{
			"Parametry oferty": parametry,
			"Media": media,
			"Informacje dodatkowe": info,
			"Opis": opis,
		}
		for k, v := range endParams {
			final[k] = v
		}

		jsonData, err := json.MarshalIndent(final, "", "  ")
		if err != nil {
			return c.Status(500).SendString(fmt.Sprintf("Błąd konwersji do JSON: %v", err))
		}

		return c.SendString(string(jsonData))
	})

	app.Listen(":8080")
}