package instadm

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/tebeka/selenium"
)

// ConnectToInstagram get ig login webpage and connect user account
func (s *Selenium) ConnectToInstagram() {
	log.Info("Connecting to Instagram account...")
	// Access Instagram url
	if err := s.WebDriver.Get("https://instagram.com/?hl=en"); err != nil {
		s.fatal("Can't access to Instagram. ", err)
	}
	randomSleep()
	// Accept cookies if requested
	if find, err := s.WaitForElement("//button[text()='Accept']", "xpath", 10); err == nil && find {
		elem, _ := s.GetElement("//button[text()='Accept']", "xpath")
		elem.Click()
		log.Debug("Cookies validation done!")
	} else {
		log.Info("Cookies validation button not found, skipping.")
	}
	randomSleep()
	// Access to login screen if needed
	if find, err := s.WaitForElement("//button[text()='Log In']", "xpath", 10); err == nil && find {
		elem, _ := s.GetElement("//button[text()='Log In']", "xpath")
		elem.Click()
		log.Debug("Log in screen access done!")
	} else {
		log.Info("Login button not found, skipping.")
	}
	randomSleep()
	// Inject username and password to input fields and log in
	if find, err := s.WaitForElement("username", "name", 10); err == nil && find {
		elem, _ := s.GetElement("username", "name")
		elem.SendKeys(s.Config.BotConfig.UserAccount.Username)
		log.Debug("Username injection done!")
	} else {
		s.fatal("Exception during username inject: ", err)
	}
	if find, err := s.WaitForElement("password", "name", 10); err == nil && find {
		elem, _ := s.GetElement("password", "name")
		elem.SendKeys(s.Config.BotConfig.UserAccount.Password)
		log.Debug("Password injection done!")
	} else {
		s.fatal("Exception during password inject: ", err)
	}
	if find, err := s.WaitForElement("//button/*[text()='Log In']", "xpath", 10); err == nil && find {
		elem, _ := s.GetElement("//button/*[text()='Log In']", "xpath")
		elem.Click()
		log.Debug("Sent login request")
	} else {
		s.fatal("Log in button not found: ", err)
	}
	randomSleepCustom(10, 15)
	// Check if login was successful
	if s.IsElementPresent(selenium.ByXPATH, "//*[@aria-label='Home'] | //button[text()='Save Info'] | //button[text()='Not Now']") {
		log.Info("Login Successful!")
	} else {
		log.Warnf("Instagram does not ask for informations saving, the login process may have failed.")
	}
}

// SendMessage navigate to Instagram direct message interface and send one to specified user
// by simulating human typing
func (s *Selenium) SendMessage(user, message string) (bool, error) {
	log.Infof("Send message to %s...", user)
	// Navigate to Instagram new direct message page
	if err := s.WebDriver.Get("https://www.instagram.com/direct/new/?hl=en"); err != nil {
		s.fatal("Can't access to Instagram direct message redaction page! ", err)
	}
	randomSleepCustom(6, 10)

	// Type and select user to dm
	if find, err := s.WaitForElement("//*[@id=\"react-root\"]/section/div[2]/div/div[1]/div/div[2]/input", "xpath", 10); err == nil && find {
		elem, _ := s.GetElement("//*[@id=\"react-root\"]/section/div[2]/div/div[1]/div/div[2]/input", "xpath")
		log.Debug("Finded an retrieved user searchbar")
		if res := SimulateHandWriting(elem, user); res != true {
			return false, errors.New("Error during user searching")
		}
		randomSleep()
		if usernames, err := s.WebDriver.FindElements(selenium.ByXPATH, "//div[@aria-labelledby]/div/span//img[@data-testid='user-avatar']"); err != nil {
			return false, errors.New("Error during user selection")
		} else {
			usernames[0].Click()
			log.Debug("User to dm selected")
		}
	} else {
		return false, errors.New("Error during user selection")
	}

	// Type and send message by simulating human writing
	if err := s.typeMessage(message); err != nil {
		return false, errors.New("Error during message typing")
	}
	log.Debug("Message sended!")

	return true, nil
}

func (s *Selenium) typeMessage(message string) error {
	if find, err := s.WaitForElement("//button/*[text()='Next']", "xpath", 5); err == nil && find {
		elem, _ := s.GetElement("//button/*[text()='Next']", "xpath")
		elem.Click()
	} else {
		log.Errorf("Error during message sending: %v", err)
		return err
	}
	randomSleep()
	if find, err := s.WaitForElement("//textarea[@placeholder]", "xpath", 5); err == nil && find {
		elem, _ := s.GetElement("//textarea[@placeholder]", "xpath")
		if res := SimulateHandWriting(elem, message); res != true {
			return errors.New("Error during message typing")
		}
	} else {
		log.Errorf("Error during message sending: %v", err)
		return err
	}
	randomSleep()
	if find, err := s.WaitForElement("//button[text()='Send']", "xpath", 5); err == nil && find {
		elem, _ := s.GetElement("//button[text()='Send']", "xpath")
		elem.Click()
	} else {
		log.Errorf("Error during message sending: %v", err)
		return err
	}
	randomSleep()

	return nil
}
