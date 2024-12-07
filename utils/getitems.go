package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type item struct {
	Images      []string
	Title       string
	Price       string
	Category    string
	Condition   string
	Description string
	Tags        []string
}

func GetItems(root string) ([]item, error) {
	var items []item

	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("error reading root directory: %w", err)
	}

	for _, entry := range entries {
		subDir := filepath.Join(root, entry.Name())
		detailsFile := filepath.Join(subDir, "details.txt")

		subEntries, err := os.ReadDir(subDir)
		if err != nil {
			continue
		}

		var imageFiles []string
		for _, subEntry := range subEntries {
			if !subEntry.IsDir() && filepath.Ext(subEntry.Name()) != ".txt" {
				filePath := filepath.Join(subDir, subEntry.Name())
				imageFiles = append(imageFiles, filePath)
			}
		}

		file, err := os.Open(detailsFile)
		if err != nil {
			continue
		}
		defer file.Close()

		var title, price, category, condition, description, tagsString string

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()

			switch {
			case strings.HasPrefix(line, "title:"):
				title = strings.ToUpper(strings.TrimSpace(line[len("title:"):]))
			case strings.HasPrefix(line, "price:"):
				price = strings.TrimSpace(line[len("price:"):])
			case strings.HasPrefix(line, "category"):
				category = strings.ToLower(strings.TrimSpace(line[len("category:"):]))
			case strings.HasPrefix(line, "condition"):
				condition = strings.ToLower(strings.TrimSpace(line[len("condition:"):]))
			case strings.HasPrefix(line, "description:"):
				description = strings.ToUpper(strings.TrimSpace(line[len("description:"):]))
			case strings.HasPrefix(line, "tags:"):
				tagsString = strings.ToUpper(strings.TrimSpace(line[len("tags:"):]))
			}
		}

		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("error reading file: %v", err)
		}

		parts := strings.Split(description, "...")

		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}

		description = strings.Join(parts, "\n\n")

		tags := strings.Split(tagsString, ",")

		for i := len(tags) - 1; i >= 0; i-- {
			if tags[i] == "" {
				tags = tags[:i]
			} else {
				break
			}
		}

		items = append(items, item{
			Images:      imageFiles,
			Title:       title,
			Price:       price,
			Category:    category,
			Condition:   condition,
			Description: description,
			Tags:        tags,
		})
	}

	return items, nil
}
