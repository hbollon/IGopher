package igopher

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/tebeka/selenium"
)

// ConnectToInstagram get ig login webpage and connect user account
func (bot *IGopher) ConnectToInstagram() {
	bot.connectToInstagramWebDriver()
}

func (bot *IGopher) connectToInstagramWebDriver() {
	log.Info("Connecting to Instagram account...")
	// Access Instagram url
	if err := bot.SeleniumStruct.WebDriver.Get("https://www.instagram.com/?hl=en"); err != nil {
		bot.SeleniumStruct.Fatal("Can't access to Instagram. ", err)
	}
	randomSleep()
	// Accept cookies if requested
	if find, err := bot.SeleniumStruct.WaitForElement("//button[text()='Accept All']", "xpath", 10); err == nil && find {
		elem, _ := bot.SeleniumStruct.GetElement("//button[text()='Accept All']", "xpath")
		elem.Click()
		log.Debug("Cookies validation done!")
	} else {
		log.Info("Cookies validation button not found, skipping.")
	}
	randomSleep()
	// Access to login screen if needed
	if find, err := bot.SeleniumStruct.WaitForElement("//button[text()='Log In']", "xpath", 10); err == nil && find {
		elem, _ := bot.SeleniumStruct.GetElement("//button[text()='Log In']", "xpath")
		elem.Click()
		log.Debug("Log in screen access done!")
	} else {
		log.Info("Login button not found, skipping.")
	}
	randomSleep()
	// Inject username and password to input fields and log in
	if find, err := bot.SeleniumStruct.WaitForElement("username", "name", 10); err == nil && find {
		elem, _ := bot.SeleniumStruct.GetElement("username", "name")
		elem.SendKeys(bot.UserAccount.Username)
		log.Debug("Username injection done!")
	} else {
		bot.SeleniumStruct.Fatal("Exception during username inject: ", err)
	}
	if find, err := bot.SeleniumStruct.WaitForElement("password", "name", 10); err == nil && find {
		elem, _ := bot.SeleniumStruct.GetElement("password", "name")
		elem.SendKeys(bot.UserAccount.Password)
		log.Debug("Password injection done!")
	} else {
		bot.SeleniumStruct.Fatal("Exception during password inject: ", err)
	}
	if find, err := bot.SeleniumStruct.WaitForElement("//button/*[text()='Log In']", "xpath", 10); err == nil && find {
		elem, _ := bot.SeleniumStruct.GetElement("//button/*[text()='Log In']", "xpath")
		elem.Click()
		log.Debug("Sent login request")
	} else {
		bot.SeleniumStruct.Fatal("Log in button not found: ", err)
	}
	randomSleepCustom(10, 15)
	// Accept second cookies prompt if requested
	if find, err := bot.SeleniumStruct.WaitForElement("//button[text()='Allow All Cookies']", "xpath", 10); err == nil && find {
		elem, _ := bot.SeleniumStruct.GetElement("//button[text()='Allow All Cookies']", "xpath")
		elem.Click()
		log.Debug("Second cookies validation done!")
		randomSleep()
	} else {
		log.Info("Second cookies validation button not found, skipping.")
	}
	// Check if login was successful
	if bot.SeleniumStruct.IsElementPresent(selenium.ByXPATH,
		"//*[@aria-label='Home'] | //button[text()='Save Info'] | //button[text()='Not Now']") {
		log.Info("Login Successful!")
	} else {
		log.Warnf("Instagram does not ask for informations saving, the login process may have failed.")
	}
}

// SendMessage navigate to Instagram direct message interface and send one to specified user
// by simulating human typing
func (bot *IGopher) SendMessage(user, message string) (bool, error) {
	if bot.Scheduler.CheckTime() == nil && (!bot.Blacklist.Activated || !bot.Blacklist.IsBlacklisted(user)) {
		res, err := bot.sendMessageWebDriver(user, message)
		if res && err == nil {
			if bot.Quotas.Activated {
				bot.Quotas.AddDm()
			}
			if bot.Blacklist.Activated {
				bot.Blacklist.AddUser(user)
			}
			log.Info("Message successfully sent!")
		}

		return res, err
	}
	return false, nil
}

func (bot *IGopher) sendMessageWebDriver(user, message string) (bool, error) {
	log.Infof("Send message to %s...", user)
	// Navigate to Instagram new direct message page
	if err := bot.SeleniumStruct.WebDriver.Get("https://www.instagram.com/direct/new/?hl=en"); err != nil {
		bot.SeleniumStruct.Fatal("Can't access to Instagram direct message redaction page! ", err)
	}
	randomSleepCustom(6, 10)

	// Type and select user to dm
	if find, err := bot.SeleniumStruct.WaitForElement(
		"//*[@id=\"react-root\"]/div/div/section/div[2]/div/div[1]/div/div[2]/input", "xpath", 10); err == nil && find {
		elem, _ := bot.SeleniumStruct.GetElement("//*[@id=\"react-root\"]/div/div/section/div[2]/div/div[1]/div/div[2]/input", "xpath")
		log.Debug("Finded an retrieved user searchbar")
		if res := SimulateHandWriting(elem, user); !res {
			return false, errors.New("Error during user searching")
		}
		randomSleep()
		usernames, err := bot.SeleniumStruct.WebDriver.FindElements(selenium.ByXPATH,
			"//div[@aria-labelledby]/div/span//img[@data-testid='user-avatar']")
		if err != nil {
			return false, errors.New("Error during user selection")
		}
		usernames[0].Click()
		log.Debug("User to dm selected")
	} else {
		return false, errors.New("Error during user selection")
	}

	// Type and send message by simulating human writing
	if err := bot.typeMessage(message); err != nil {
		return false, errors.New("Error during message typing")
	}
	log.Debug("Message sended!")

	return true, nil
}

func (bot *IGopher) typeMessage(message string) error {
	if find, err := bot.SeleniumStruct.WaitForElement("//button/*[text()='Next']", "xpath", 5); err == nil && find {
		elem, _ := bot.SeleniumStruct.GetElement("//button/*[text()='Next']", "xpath")
		elem.Click()
	} else {
		log.Errorf("Error during message sending: %v", err)
		return err
	}
	randomSleep()
	if find, err := bot.SeleniumStruct.WaitForElement("//textarea[@placeholder]", "xpath", 5); err == nil && find {
		elem, _ := bot.SeleniumStruct.GetElement("//textarea[@placeholder]", "xpath")
		if res := SimulateHandWriting(elem, message); !res {
			return errors.New("Error during message typing")
		}
	} else {
		log.Errorf("Error during message sending: %v", err)
		return err
	}
	randomSleep()
	if find, err := bot.SeleniumStruct.WaitForElement("//button[text()='Send']", "xpath", 5); err == nil && find {
		elem, _ := bot.SeleniumStruct.GetElement("//button[text()='Send']", "xpath")
		elem.Click()
	} else {
		log.Errorf("Error during message sending: %v", err)
		return err
	}
	randomSleep()

	return nil
}
