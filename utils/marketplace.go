package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
)

var (
	duration = 20 * time.Microsecond
)

func ListItems(browser *rod.Browser, page *rod.Page, phones []item) error {
	defer browser.MustClose()

	itemsProccessed := 0
	start := time.Now()

	for _, p := range phones {
		page = page.MustNavigate("https://www.facebook.com/marketplace/create/item").MustWaitLoad()

		fmt.Println("Navigated to create item page")

		page.MustElement(`input[type="file"][accept="image/*,image/heif,image/heic"]`).MustSetFiles(p.Images...)

		fmt.Println("images inserted")

		title := p.Title

		inp := page.MustElement(`label[aria-label="Title"]`)

		for _, char := range title {
			inp.MustInput(string(char))
			time.Sleep(duration)
		}

		fmt.Println("title inserted")

		time.Sleep(3 * time.Second)

		price := p.Price
		inp = page.MustElement(`label[aria-label="Price"]`)

		for _, char := range price {
			inp.MustInput(string(char))
			time.Sleep(duration)
		}

		fmt.Println("price inserted")

		time.Sleep(3 * time.Second)

		page.MustElement(`label[aria-label="Category"]`).MustClick()

		cats := page.MustElements(`div[data-visualcompletion="ignore-dynamic"]`)

		for _, cat := range cats {
			if strings.ToLower(cat.MustText()) == p.Category {
				cat.MustClick()
				break
			}
		}

		fmt.Println("category selected")

		time.Sleep(3 * time.Second)

		page.MustElement(`label[aria-label="Condition"]`).MustClick()

		options := page.MustElements(`div[role="option"]`)

		fmt.Println("Conditions retrieved")

		for _, option := range options {
			if strings.ToLower(option.MustText()) == p.Condition {
				option.MustClick()
				break
			}
		}

		fmt.Println("condition selected")

		time.Sleep(3 * time.Second)

		inp = page.MustElement(`label[aria-label="Description"]`)
		desc := p.Description
		for _, char := range desc {
			inp.MustInput(string(char))
			time.Sleep(duration)
		}

		fmt.Println("description inserted")

		time.Sleep(3 * time.Second)

		tagsInput := page.MustElement(`label[aria-label="Product tags"]`)

		for _, tag := range p.Tags {
			tag = strings.TrimSpace(tag)

			for _, char := range tag {
				tagsInput.MustInput(string(char))
				time.Sleep(duration)
			}

			err := page.Keyboard.Press(input.Enter)
			if err != nil {
				return fmt.Errorf("error pressing enter key: %v", err)
			}

			time.Sleep(1 * time.Second)
		}

		fmt.Println("tags inserted")

		time.Sleep(3 * time.Second)

		page.MustElement(`div[aria-label="Next"]`).MustClick()

		fmt.Println("go to next page")

		time.Sleep(3 * time.Second)

		page.MustScreenshot("home.png")

		groups := page.MustElements(`div[role="checkbox"]`)

		fmt.Println("All Groups:", len(groups))

		var tickedGroups int = 0

		for _, group := range groups {
			group.MustClick().MustScreenshot("group.png")
			time.Sleep(100 * time.Millisecond)
			tickedGroups++
		}

		fmt.Println("Ticked Groups:", tickedGroups)

		page.MustElement(`div[aria-label="Publish"]`).MustClick()

		for i := 0; i < 30; i++ {
			if page.MustInfo().URL == "https://web.facebook.com/marketplace/you/selling" {
				fmt.Printf("%s published successfully\n", title)
				itemsProccessed++
				break
			}
			fmt.Printf("Publishing %s...\n", title)
			time.Sleep(2 * time.Second)
		}

		fmt.Printf("%d items proccessed\n", itemsProccessed)

		time.Sleep(10 * time.Second)
	}

	duration := time.Since(start).Minutes()

	fmt.Println("All Items Listed Successfully In %w", duration)

	return nil
}
