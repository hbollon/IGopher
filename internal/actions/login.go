package actions

import (
	"github.com/hbollon/igopher/internal/config/types"
	"github.com/hbollon/igopher/internal/utils"
	"github.com/hbollon/igopher/internal/xpath"
	log "github.com/sirupsen/logrus"
	"github.com/tebeka/selenium"
)

// ConnectToInstagram get ig login webpage and connect user account
func ConnectToInstagram(bot *types.IGopher) {
	connectToInstagramWebDriver(bot)
}

func connectToInstagramWebDriver(bot *types.IGopher) {
	log.Info("Connecting to Instagram account...")
	// Access Instagram url
	if err := bot.SeleniumStruct.WebDriver.Get("https://www.instagram.com/?hl=en"); err != nil {
		bot.SeleniumStruct.Fatal("Can't access to Instagram. ", err)
	}
	utils.RandomSleep()
	// Accept cookies if requested
	if find, err := bot.SeleniumStruct.WaitForElement(xpath.XPathSelectors["login_accept_cookies"],
		"xpath", 10); err == nil && find {
		elem, _ := bot.SeleniumStruct.GetElement(xpath.XPathSelectors["login_accept_cookies"], "xpath")
		elem.Click()
		log.Debug("Cookies validation done!")
	} else {
		log.Info("Cookies validation button not found, skipping.")
	}
	utils.RandomSleep()
	// Access to login screen if needed
	if find, err := bot.SeleniumStruct.WaitForElement(xpath.XPathSelectors["login_button"], "xpath", 10); err == nil && find {
		elem, _ := bot.SeleniumStruct.GetElement(xpath.XPathSelectors["login_button"], "xpath")
		elem.Click()
		log.Debug("Log in screen access done!")
	} else {
		log.Info("Login button not found, skipping.")
	}
	utils.RandomSleep()
	// Inject username and password to input fields and log in
	if find, err := bot.SeleniumStruct.WaitForElement(xpath.XPathSelectors["login_username"], "name", 10); err == nil && find {
		elem, _ := bot.SeleniumStruct.GetElement(xpath.XPathSelectors["login_username"], "name")
		elem.SendKeys(bot.UserAccount.Username)
		log.Debug("Username injection done!")
	} else {
		bot.SeleniumStruct.Fatal("Exception during username inject: ", err)
	}
	if find, err := bot.SeleniumStruct.WaitForElement(xpath.XPathSelectors["login_password"], "name", 10); err == nil && find {
		elem, _ := bot.SeleniumStruct.GetElement(xpath.XPathSelectors["login_password"], "name")
		elem.SendKeys(bot.UserAccount.Password)
		log.Debug("Password injection done!")
	} else {
		bot.SeleniumStruct.Fatal("Exception during password inject: ", err)
	}
	if find, err := bot.SeleniumStruct.WaitForElement(xpath.XPathSelectors["login_alternate_button"], "xpath", 10); err == nil && find {
		elem, _ := bot.SeleniumStruct.GetElement(xpath.XPathSelectors["login_alternate_button"], "xpath")
		elem.Click()
		log.Debug("Sent login request")
	} else {
		bot.SeleniumStruct.Fatal("Log in button not found: ", err)
	}
	utils.RandomSleepCustom(10, 15)
	// Accept second cookies prompt if requested
	if find, err := bot.SeleniumStruct.WaitForElement(xpath.XPathSelectors["login_alternate_button"], "xpath", 10); err == nil && find {
		elem, _ := bot.SeleniumStruct.GetElement(xpath.XPathSelectors["login_alternate_button"], "xpath")
		elem.Click()
		log.Debug("Second cookies validation done!")
		utils.RandomSleep()
	} else {
		log.Info("Second cookies validation button not found, skipping.")
	}
	// Check if login was successful
	if bot.SeleniumStruct.IsElementPresent(selenium.ByXPATH,
		xpath.XPathSelectors["login_information_saving"]) {
		log.Info("Login Successful!")
	} else {
		if err := bot.SeleniumStruct.WebDriver.Refresh(); err != nil {
			bot.SeleniumStruct.Fatal("Can't refresh page: ", err)
		}
		if find, err := bot.SeleniumStruct.WaitForElement(xpath.XPathSelectors["login_information_saving"],
			"xpath", 10); err != nil || !find {
			log.Warnf("Instagram does not ask for informations saving or app download, the login process may have failed.")
		}
	}
}
