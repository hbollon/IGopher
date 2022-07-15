package actions

import (
	"errors"

	"github.com/hbollon/igopher/internal/config/types"
	"github.com/hbollon/igopher/internal/simulation"
	"github.com/hbollon/igopher/internal/utils"
	"github.com/hbollon/igopher/internal/xpath"
	log "github.com/sirupsen/logrus"
	"github.com/tebeka/selenium"
)

// SendMessage navigate to Instagram direct message interface and send one to specified user
// by simulating human typing
func SendMessage(bot *types.IGopher, user, message string) (bool, error) {
	if bot.Scheduler.CheckTime() == nil && (!bot.Blacklist.Activated || !bot.Blacklist.IsBlacklisted(user)) {
		res, err := sendMessageWebDriver(bot, user, message)
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

func sendMessageWebDriver(bot *types.IGopher, user, message string) (bool, error) {
	log.Infof("Send message to %s...", user)
	// Navigate to Instagram new direct message page
	if err := bot.SeleniumStruct.WebDriver.Get("https://www.instagram.com/direct/new/?hl=en"); err != nil {
		bot.SeleniumStruct.Fatal("Can't access to Instagram direct message redaction page! ", err)
	}
	utils.RandomSleepCustom(6, 10)

	// Type and select user to dm
	if find, err := bot.SeleniumStruct.WaitForElement(
		xpath.XPathSelectors["dm_user_search"], "xpath", 10); err == nil && find {
		elem, _ := bot.SeleniumStruct.GetElement(xpath.XPathSelectors["dm_user_search"], "xpath")
		log.Debug("Finded an retrieved user searchbar")
		if res := simulation.SimulateHandWriting(elem, user); !res {
			return false, errors.New("Error during user searching")
		}
		utils.RandomSleep()
		usernames, err := bot.SeleniumStruct.WebDriver.FindElements(selenium.ByXPATH,
			xpath.XPathSelectors["dm_profile_pictures_links"])
		if err != nil {
			return false, errors.New("Error during user selection")
		}
		usernames[0].Click()
		log.Debug("User to dm selected")
	} else {
		return false, errors.New("Error during user selection")
	}

	// Type and send message by simulating human writing
	if err := typeMessage(bot, message); err != nil {
		return false, errors.New("Error during message typing")
	}
	log.Debug("Message sended!")

	return true, nil
}

func typeMessage(bot *types.IGopher, message string) error {
	if find, err := bot.SeleniumStruct.WaitForElement(xpath.XPathSelectors["dm_next_button"], "xpath", 5); err == nil && find {
		elem, _ := bot.SeleniumStruct.GetElement(xpath.XPathSelectors["dm_next_button"], "xpath")
		elem.Click()
	} else {
		log.Errorf("Error during message sending: %v", err)
		return err
	}
	utils.RandomSleep()
	if find, err := bot.SeleniumStruct.WaitForElement(xpath.XPathSelectors["dm_placeholder"], "xpath", 5); err == nil && find {
		elem, _ := bot.SeleniumStruct.GetElement(xpath.XPathSelectors["dm_placeholder"], "xpath")
		if res := simulation.SimulateHandWriting(elem, message); !res {
			return errors.New("Error during message typing")
		}
	} else {
		log.Errorf("Error during message sending: %v", err)
		return err
	}
	utils.RandomSleep()
	if find, err := bot.SeleniumStruct.WaitForElement(xpath.XPathSelectors["dm_send_button"], "xpath", 5); err == nil && find {
		elem, _ := bot.SeleniumStruct.GetElement(xpath.XPathSelectors["dm_send_button"], "xpath")
		elem.Click()
	} else {
		log.Errorf("Error during message sending: %v", err)
		return err
	}
	utils.RandomSleep()

	return nil
}
