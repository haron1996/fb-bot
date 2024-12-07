package utils

import (
	"fmt"
	"time"

	"math/rand"

	"github.com/go-rod/rod"
)

func PostToGroups(browser *rod.Browser, page *rod.Page, items []item) {
	defer browser.MustClose()

	page = page.MustNavigate("https://web.facebook.com/groups/joins").MustWaitLoad()

	fmt.Println("Navigated to my groups page")

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	seen := make(map[string]struct{})

	scrollAttempts := 0
	maxAttempts := 5

	for i := 0; i < 100; i++ {
		page.Mouse.MustScroll(0, 1000)

		time.Sleep(3 * time.Second)

		classSelector := ".x9f619.x1gryazu.xkrivgy.x1ikqzku.x1h0ha7o.xg83lxy.xh8yej3"
		parents := page.MustElements(classSelector)

		var parent *rod.Element
		parentsLength := len(parents)

		if parentsLength == 1 {
			parent = parents[0]
		} else if parentsLength == 2 {
			parent = parents[1]
		} else {
			fmt.Println("no parents found")
			return
		}

		selector := ".x1i10hfl.x1qjc9v5.xjbqb8w.xjqpnuy.xa49m3k.xqeqjp1.x2hbi6w.x13fuv20.xu3j5b3.x1q0q8m5.x26u7qi.x972fbf.xcfux6l.x1qhh985.xm0m39n.x9f619.x1ypdohk.xdl72j9.x2lah0s.xe8uvvx.xdj266r.x11i5rnm.xat24cr.x1mh8g0r.x2lwn1j.xeuugli.xexx8yu.x4uap5.x18d9i69.xkhd6sd.x1n2onr6.x16tdsg8.x1hl2dhg.xggy1nq.x1ja2u2z.x1t137rt.x1o1ewxj.x3x9cwd.x1e5q0jg.x13rtm0m.x1q0g3np.x87ps6o.x1lku1pv.x1rg5ohu.x1a2a7pz.xh8yej3"
		links := parent.MustElements(selector)

		initialCount := len(seen)

		for _, link := range links {
			anchorIDPtr := link.MustAttribute("href")
			if anchorIDPtr != nil {
				anchorID := *anchorIDPtr
				if _, exists := seen[anchorID]; !exists {
					seen[anchorID] = struct{}{}
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

	}

	fmt.Println("Total Unique Anchors:", len(seen))

	hrefs := make([]string, 0, len(seen))
	for href := range seen {
		hrefs = append(hrefs, href)
	}

	r.Shuffle(len(hrefs), func(i, j int) {
		hrefs[i], hrefs[j] = hrefs[j], hrefs[i]
	})

	r.Shuffle(len(items), func(i, j int) {
		items[i], items[j] = items[j], items[i]
	})

	randomPostContent := items[0]

	for _, href := range hrefs {
		fmt.Println(href)

		page = page.MustNavigate(href).MustWaitLoad()
		fmt.Println("Navigated to random group")

		if page.MustHas(`div.x1yztbdb.x1j9u4d2`) {
			fmt.Println("Blocked from accessing this group")
			hrefs = removeHref(hrefs, href)
			fmt.Println("Removed random href from selected groups")
			fmt.Println("Remaining Groups:", len(hrefs))
			fmt.Println("Continue to the next iteration")
			continue
		}

		pageHasSellBtn, sellBtn, _ := page.Has(`div[aria-label="Sell Something"]`)

		switch {
		case pageHasSellBtn:
			fmt.Println("This is a buy and sell group:", href)
			fmt.Println("Handle a case where page has sell button")
			sellBtn.MustScreenshot("sellButton.png")
			continue
		default:
			fmt.Println("This is a normal group")

			pageHasWriteSomethingDiv := page.MustHas(`div.x1i10hfl.x1ejq31n.xd10rxx.x1sy0etr.x17r0tee.x972fbf.xcfux6l.x1qhh985.xm0m39n.x9f619.x1ypdohk.xe8uvvx.xdj266r.x11i5rnm.xat24cr.x1mh8g0r.x16tdsg8.x1hl2dhg.xggy1nq.x87ps6o.x1lku1pv.x1a2a7pz.x6s0dn4.xmjcpbm.x107yiy2.xv8uw2v.x1tfwpuw.x2g32xy.x78zum5.x1q0g3np.x1iyjqo2.x1nhvcw1.x1n2onr6.xt7dq6l.x1ba4aug.x1y1aw1k.xn6708d.xwib8y2.x1ye3gou > div.xi81zsa.x1lkfr7t.xkjl1po.x1mzt3pk.xh8yej3.x13faqbe > span.x1lliihq.x6ikm8r.x10wlt62.x1n2onr6`)

			fmt.Println("Page Has Write Something Element:", pageHasWriteSomethingDiv)

			if !pageHasWriteSomethingDiv {
				fmt.Println("Page does not have write something div")
				page.MustScreenshot("home.png")
				continue
			}

			writeSomethingDiv := page.MustElement(`div.x1i10hfl.x1ejq31n.xd10rxx.x1sy0etr.x17r0tee.x972fbf.xcfux6l.x1qhh985.xm0m39n.x9f619.x1ypdohk.xe8uvvx.xdj266r.x11i5rnm.xat24cr.x1mh8g0r.x16tdsg8.x1hl2dhg.xggy1nq.x87ps6o.x1lku1pv.x1a2a7pz.x6s0dn4.xmjcpbm.x107yiy2.xv8uw2v.x1tfwpuw.x2g32xy.x78zum5.x1q0g3np.x1iyjqo2.x1nhvcw1.x1n2onr6.xt7dq6l.x1ba4aug.x1y1aw1k.xn6708d.xwib8y2.x1ye3gou > div.xi81zsa.x1lkfr7t.xkjl1po.x1mzt3pk.xh8yej3.x13faqbe > span.x1lliihq.x6ikm8r.x10wlt62.x1n2onr6`)

			if writeSomethingDiv == nil {
				fmt.Println("writeSomethingDiv is nil, cannot click")
				continue
			}

			writeSomethingDiv.MustClick()

			time.Sleep(5 * time.Second)

			fmt.Println("Clicked on write something button")

			dialog := page.MustElement(`div.x1n2onr6.x1ja2u2z.x1afcbsf.x78zum5.xdt5ytf.x1a2a7pz.x6ikm8r.x10wlt62.x71s49j.x1jx94hy.x1qpq9i9.xdney7k.xu5ydu1.xt3gfkd.x104qc98.x1g2kw80.x16n5opg.xl7ujzl.xhkep3z.x193iq5w[role="dialog"]`).MustWaitLoad()

			if !dialog.MustInteractable() {
				fmt.Println("Dialog is not interactable")
				time.Sleep(30 * time.Second)
				if !dialog.MustInteractable() {
					fmt.Println("Dialog is not interactable after waiting for 30 seconds")
					continue
				}
			}

			fmt.Println("Past dialog")

			for i := 0; i < 10; i++ {
				element := page.MustElement(`div[aria-label="Photo/video"]`)
				if element != nil {
					fmt.Println("Found photo/video element")
					element.MustClick()
					break
				}
				time.Sleep(2 * time.Second)
			}

			fmt.Println("clicked photo/video icon")

			images := randomPostContent.Images

			for i := 0; i < 10; i++ {
				element := dialog.MustElement(`input[type="file"]`)
				if element != nil {
					fmt.Println("Found photo/video element")
					element.MustSetFiles(images...)
					break
				}
				time.Sleep(2 * time.Second)
			}

			fmt.Println("Images inserted")

			dialogHasContentEditableDiv, contentEditableDiv, _ := dialog.Has(`div[aria-label="Create a public post…"]`)

			desc := randomPostContent.Description

			if dialogHasContentEditableDiv {
				fmt.Println("Dialog has a create a public post div")
				contentEditableDiv.MustInput(desc)
			} else {
				fmt.Println("Dialog does not have create a public post div")
				fmt.Println(page.Has(`div[aria-label="Write something..."]`))
				page.MustElement(`div[aria-label="Write something..."]`).MustInput(desc)
			}

			fmt.Println("Description inserted")

			time.Sleep(5 * time.Second)

			btn := dialog.MustElement(`div[aria-label="Post"]`)

			btn.MustClick()

			fmt.Println("Post button clicked")

			for i := 0; i < 60; i++ {
				if dialog.MustVisible() {
					fmt.Println("Processing...")
					time.Sleep(2 * time.Second)
					continue
				}

				fmt.Println("Ad posted successfully 🥳🥳🥳🥳🥳🥳🥳🥳🥳")
				break
			}

			title := randomPostContent.Title

			fmt.Printf("%s posted successfully 🥳🥳🥳🥳🥳🥳🥳🥳🥳\n", title)

			hrefs = removeHref(hrefs, href)

			fmt.Println("Removed random href from selected groups")
			fmt.Println("Remaining Groups:", len(hrefs))

			page.MustScreenshot("home.png")

			fmt.Println(time.Now().Local())

			time.Sleep(30 * time.Second)

			fmt.Println(btn.Visible())

			fmt.Println("Continue to the next iteration")

			continue
		}
	}

	fmt.Println("Finished posting to all selected groups")
}

func removeHref(slice []string, item string) []string {
	newSlice := []string{}

	for _, v := range slice {
		if v != item {
			newSlice = append(newSlice, v)
		}
	}
	return newSlice
}
