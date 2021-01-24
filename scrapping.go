package igopher

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tebeka/selenium"
)

// ScrapperConfig store scrapper configuration for user fetching
// It also store fetched usernames
type ScrapperConfig struct {
	SrcAccounts     []string `yaml:"src_accounts"`
	FetchedAccounts []string
	Quantity        int `yaml:"fetch_quantity" validate:"numeric"`
}

func (sc *IGopher) FetchUsersFromUserFollowers() ([]string, error) {
	logrus.Info("Fetching users from users followers...")
	var igUsers []string
	for _, srcUsername := range sc.ScrapperManager.SrcAccounts {
		// Navigate to Instagram user page
		if err := sc.SeleniumStruct.WebDriver.Get(fmt.Sprintf("https://www.instagram.com/%s/?hl=en", srcUsername)); err != nil {
			logrus.Warnf("Requested user '%s' doesn't exist, skip it", srcUsername)
		}
		randomSleepCustom(1, 3)
		if find, err := sc.SeleniumStruct.WaitForElement("//*[@id=\"react-root\"]/section/main/div/ul/li[2]/a", "xpath", 10); err == nil && find {
			elem, _ := sc.SeleniumStruct.GetElement("//*[@id=\"react-root\"]/section/main/div/ul/li[2]/a", "xpath")
			elem.Click()
			logrus.Debug("Clicked on user followers list")
		} else {
			return nil, errors.New("Error during access to user followers list")
		}
		randomSleepCustom(1, 3)
		var dialog selenium.WebElement
		if find, err := sc.SeleniumStruct.WaitForElement("//*[@id=\"react-root\"]/section/main/div", "xpath", 10); err == nil && find {
			dialog, _ = sc.SeleniumStruct.GetElement("//*[@id=\"react-root\"]/section/main/div", "xpath")
			dialog.Click()
			logrus.Debug("Clicked on user followers dialog box")
		} else {
			return nil, errors.New("Error during focus user followers list dialog")
		}

		var scrappedUsers []selenium.WebElement
		for len(scrappedUsers) < sc.ScrapperManager.Quantity {
			if len(scrappedUsers) != 0 {
				_, err = sc.SeleniumStruct.WebDriver.ExecuteScript("arguments[0].scrollIntoView();", []interface{}{dialog})
				if err != nil {
					return nil, errors.New("Error during followers dialog box scroll")
				}
			}
			randomSleepCustom(1, 2)
			scrappedUsers, err = sc.SeleniumStruct.WebDriver.FindElements(selenium.ByXPATH, "//*/li/div/div/div/div/a")
			if err != nil {
				logrus.Error(err)
				return nil, errors.New("Error during users scrapping from followers dialog box")
			}
			fmt.Println(len(scrappedUsers))
		}

		for _, user := range scrappedUsers {
			username, err := user.Text()
			if err == nil {
				igUsers = append(igUsers, username)
			}
		}

		logrus.Debugf("Scrapped users: %v", igUsers)
	}
	if igUsers == nil || len(igUsers) == 0 {
		return nil, errors.New("Empty users result")
	}
	return igUsers, nil
}
