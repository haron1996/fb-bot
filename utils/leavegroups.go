package utils

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
)

func LeaveGroups(browser *rod.Browser, page *rod.Page) {
	defer browser.MustClose()

	page = page.MustNavigate("https://web.facebook.com/groups/joins").MustWaitLoad().MustWaitDOMStable().MustWaitIdle()

	fmt.Println("Navigated to groups page")

	seen := make(map[string]struct{})

	scrollAttempts := 0
	maxAttempts := 5

	for i := 0; i < 100; i++ {
		page.Mouse.MustScroll(0, 1000)
		time.Sleep(5 * time.Second)

		if !page.MustHas(`div[role="listitem"]`) {
			fmt.Println("Page has no list item cards")
			return
		}

		cards := page.MustElements(`div[role="listitem"]`)

		initialCount := len(seen)

		for _, card := range cards {

			link, err := card.ElementR(`a[role="link"]`, `.*`)
			if err != nil {
				fmt.Println("Error getting link:", err)
				continue
			}

			if link != nil {
				href := link.MustAttribute("href")
				if href != nil {
					if _, exists := seen[*href]; !exists {
						seen[*href] = struct{}{}
					}
				}
			}
		}

		if len(seen) == initialCount {
			scrollAttempts++
			if scrollAttempts >= maxAttempts {
				fmt.Println("End of page reached - no new elements detected")
				break
			}
		} else {
			scrollAttempts = 0
		}

		page.MustScreenshot("home.png")
	}

	fmt.Printf("Total unique hrefs founs: %d\n", len(seen))
	for href := range seen {
		fmt.Println(href)
		page := page.MustNavigate(href).MustWaitLoad().MustWaitIdle().MustWaitDOMStable()
		if !page.MustHas(`div[aria-label="Joined"]`) {
			fmt.Println("Page has no joined button")
			delete(seen, href)
			fmt.Println("Remaining hrefs:", len(seen))
			continue
		}
		joined := page.MustElement(`div[aria-label="Joined"]`)
		joined.MustClick()
		time.Sleep(3 * time.Second)
		btns := page.MustElements(`div[role="menuitem"]`)
		for _, btn := range btns {
			btnText := btn.MustText()
			if btnText == "Leave group" {
				btn.MustClick()
				time.Sleep(3 * time.Second)
				btn = page.MustElement(`div[aria-label="Leave Group"]`)
				btn.MustClick()
				time.Sleep(5 * time.Second)
				if btn.MustVisible() {
					fmt.Println("Btn still visible")
					time.Sleep(5 * time.Second)
				}
				fmt.Println("Left group")
				page.MustScreenshot("home.png")
			}
		}
		delete(seen, href)
		fmt.Println("Remaining hrefs:", len(seen))
	}

	fmt.Println("Finished quitting selected groups")
}
