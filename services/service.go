package services

import (
	"errors"
	"fmt"
	"strings"
)

// Funkcja dostaje liste linijek i dzieli je na key i value
func ParseParametry(lines []string) (map[string]string, error) {
	result := make(map[string]string)

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			return nil, fmt.Errorf("linia %d nie zawiera ':' -> %q", i+1, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "" || value == "" {
			return nil, fmt.Errorf("linia %d ma pusty key lub value -> %q", i+1, line)
		}

		result[key] = value
	}

	return result, nil
}

// Rozdziela tekst na elementy, brak ':' = true
func ParseMedia(line string) (map[string]interface{}, error) {
	if strings.TrimSpace(line) == "" {
		return nil, errors.New("pusta linia media")
	}

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

			if key == "" || value == "" {
				return nil, fmt.Errorf("niepoprawny media item: %q", item)
			}

			result[key] = value
		} else {
			result[item] = true
		}
	}

	return result, nil
}

// Dzieli tekst na listę informacji
func ParseInfo(line string) ([]string, error) {
	if strings.TrimSpace(line) == "" {
		return nil, errors.New("pusta linia info")
	}

	var result []string
	for _, p := range strings.Split(line, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}

	if len(result) == 0 {
		return nil, errors.New("brak poprawnych informacji")
	}

	return result, nil
}

// Oddziela opis od parametrów
func ParseOpis(lines []string) (string, map[string]string, error) {
	opisBuilder := strings.Builder{}
	endParams := make(map[string]string)

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			if key == "" {
				return "", nil, fmt.Errorf("pusty key w linii %d", i+1)
			}

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

	if opisBuilder.Len() == 0 && len(endParams) == 0 {
		return "", nil, errors.New("brak opisu i parametrów")
	}

	return opisBuilder.String(), endParams, nil
}
