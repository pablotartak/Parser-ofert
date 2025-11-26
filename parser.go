package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strings"
)
//Funkcja dostaje liste linijek pocietych w main i dzieli je na key i value
func parseParametry(lines []string) map[string]string {
	result := make(map[string]string)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		result[key] = value
	}
	return result
}
//Rozdziela tekst na osobne elementy na podstawie przecinka a poźniej dzieli to na key i value na podstawie dwukropka jeśli nie ma dwukropka traktuje to jako true
func parseMedia(line string) map[string]interface{} {
	result := make(map[string]interface{})
	items := strings.Split(line, ",")
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		parts := strings.SplitN(item, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			result[key] = value
		} else {
			result[item] = true
		}
	}
	return result
}
//dzieli tekst na wartości biorąc "informacje dodatkowe " jako key
func parseInfo(line string) []string {
	var result []string
	for _, p := range strings.Split(line, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
//  Oddziela opis od parametrów , opis łączy w jeden string i przypisuje do klucza "opis"
// Linie w formacie key: value trafiają do mapy endParams
func parseOpis(lines []string) (string, map[string]string) {
	opisBuilder := strings.Builder{}
	endParams := make(map[string]string)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			if value != "" {
				endParams[key] = value
				continue
			}
		}
		if opisBuilder.Len() > 0 {
			opisBuilder.WriteString(" ")
		}
		opisBuilder.WriteString(line)
	}

	return opisBuilder.String(), endParams
}

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