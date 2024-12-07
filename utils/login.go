package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/haron1996/fb-bot/viper"
)

func Login() (*rod.Browser, *rod.Page, error) {
	dir := "~/.config/google-chrome"

	u := launcher.New().UserDataDir(dir).Leakless(true).NoSandbox(true).Headless(false).MustLaunch()

	browser := rod.New().ControlURL(u).MustConnect().NoDefaultDevice()

	page := browser.MustPage("https://web.facebook.com/").MustWaitLoad().MustWindowMaximize()

	pageHasLoginButton := page.MustHas(`button[name="login"]`)

	c, err := viper.LoadConfig(".")
	if err != nil {
		return nil, nil, fmt.Errorf("error loading config file: %v", err)
	}

	switch {
	case pageHasLoginButton:
		cookies := []*proto.NetworkCookieParam{
			{
				Name:     "c_user",
				Value:    c.C_User,
				Domain:   ".facebook.com",
				Path:     "/",
				HTTPOnly: true,
				Secure:   true,
				SameSite: "None",
				Priority: "Medium",
			},
			{
				Name:     "datr",
				Value:    c.Datr,
				Domain:   ".facebook.com",
				Path:     "/",
				HTTPOnly: true,
				Secure:   true,
				SameSite: "None",
				Priority: "Medium",
			},
			{
				Name:     "fr",
				Value:    c.Fr,
				Domain:   ".facebook.com",
				Path:     "/",
				HTTPOnly: true,
				Secure:   true,
				SameSite: "None",
				Priority: "Medium",
			},
			{
				Name:     "presence",
				Value:    c.Presence,
				Domain:   ".facebook.com",
				Path:     "/",
				HTTPOnly: true,
				Secure:   true,
				SameSite: "None",
				Priority: "Medium",
			},
			{
				Name:     "sb",
				Value:    c.Sb,
				Domain:   ".facebook.com",
				Path:     "/",
				HTTPOnly: true,
				Secure:   true,
				SameSite: "None",
				Priority: "Medium",
			},
			{
				Name:     "wd",
				Value:    c.Wd,
				Domain:   ".facebook.com",
				Path:     "/",
				HTTPOnly: true,
				Secure:   true,
				SameSite: "None",
				Priority: "Medium",
			},
			{
				Name:     "xs",
				Value:    c.Xs,
				Domain:   ".facebook.com",
				Path:     "/",
				HTTPOnly: true,
				Secure:   true,
				SameSite: "None",
				Priority: "Medium",
			},
		}

		err := browser.SetCookies(cookies)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to set session cookie: %v", err)
		}

		page = page.MustNavigate("https://web.facebook.com/").MustWaitLoad()

		pageHasLoginForm, _, err := page.Has(`form[data-testid="royal_login_form"]`)
		if err != nil {
			log.Println("Error checking if page has login form:", err)
			return nil, nil, fmt.Errorf("error checking if page has login form: %v", err)
		}
		switch {
		case pageHasLoginForm:
			fmt.Println("Invalid or expired cookies 😞")
			os.Exit(1)
		default:
			fmt.Println("Log in complete 😊")
			page.MustScreenshot("home.png")
			return browser, page, nil
		}
	default:
		fmt.Println("User is logged in 😊")
		page.MustScreenshot("home.png")
	}

	return browser, page, nil
}
