package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
)

func parse(rawText string) (string, error) {
	fields := strings.Split(rawText, ",")
	result := make(map[string]interface{})

	for i := 0; i < len(fields); i++ {
		parts := strings.Split(fields[i], ":")
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch{
		case value == "tak":
			result[key] = true
		case value == "nie":
			result[key] = false
		case key == "iloscpokoi" :
			val, err := strconv.Atoi(value)
			if err != nil {
				return "", fmt.Errorf("Błąd konwersji iloscpokoi: %v", err)
			}
			result[key] = val
		case key == "metraz" :
			val, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return "", fmt.Errorf("Błąd konwersji metraz: %v", err)
			}
			result[key] = val 
		default: 
			result[key] = value
		
	}
	}
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("Błąd zamiany na JSON: %v", err)
	}
	return string(jsonData), nil
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

		jsonString, err := parse(raw)
		if err != nil {
			return c.Status(400).SendString(err.Error())
		}

		return c.SendString(jsonString)
	})

	app.Listen(":8080")
}
