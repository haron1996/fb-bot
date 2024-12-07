package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/haron1996/fb-bot/utils"
)

var mu sync.Mutex
var isRunning bool

func runEvery30Seconds() {
	location, err := time.LoadLocation("Africa/Nairobi")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return
	}

	fmt.Println("Location:", location)

	currentTime := time.Now().In(location)
	hour := currentTime.Hour()

	if hour == 7 || hour == 11 || hour == 16 || hour == 19 {
		mu.Lock()

		if isRunning {
			fmt.Println("Function is already running, skipping this cycle:", currentTime)
			mu.Unlock()
			return
		}

		isRunning = true
		mu.Unlock()

		fmt.Println("Function started at", currentTime)

		browser, page, err := utils.Login()
		if err != nil {
			fmt.Println(err)
			return
		}
		if browser == nil || page == nil {
			if browser == nil {
				fmt.Println("Browser is nill")
			}

			if page == nil {
				fmt.Println("Page is nill")
			}

			return
		}

		os.Exit(1)

		root := "/home/kwandapchumba/Pictures/SIMU"

		items, err := utils.GetItems(root)
		if err != nil {
			fmt.Println(err)
			return
		}

		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		r.Shuffle(len(items), func(i, j int) {
			items[i], items[j] = items[j], items[i]
		})

		if err := utils.ListItemsInMarketplace(browser, page, items); err != nil {
			fmt.Println(err)
			return
		}

		// OTHER OPTIONS

		// utils.ListInMorePlaces(browser, page)
		// utils.PostToGroups(browser, page, items)
		// utils.LeaveGroups(browser, page)

		fmt.Println("All Phones Have Been Listed")

		mu.Lock()
		isRunning = false
		mu.Unlock()
	} else {
		fmt.Println("Current time is outside the specified window:", currentTime)
	}
}

func main() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			runEvery30Seconds()
		}
	}()

	select {}
}
