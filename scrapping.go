package igopher

import (
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tebeka/selenium"
	"github.com/vbauerster/mpb/v6"
	"github.com/vbauerster/mpb/v6/decor"
)

// ScrapperConfig store scrapper configuration for user fetching
// It also store fetched usernames
type ScrapperConfig struct {
	SrcAccounts     []string `yaml:"src_accounts"`
	FetchedAccounts []string
	Quantity        int `yaml:"fetch_quantity" validate:"numeric"`
}

// FetchUsersFromUserFollowers scrap username list from users followers.
// Source accounts and quantity are set by the bot user.
func (sc *IGopher) FetchUsersFromUserFollowers() ([]string, error) {
	logrus.Info("Fetching users from users followers...")

	var igUsers []string
	// Valid configuration checking before fetching process
	if len(sc.ScrapperManager.SrcAccounts) == 0 || sc.ScrapperManager.SrcAccounts == nil {
		return nil, errors.New("No source users are set, please check your scrapper settings and retry")
	}
	if sc.ScrapperManager.Quantity <= 0 {
		return nil, errors.New("Scrapping quantity is null or negative, please check your scrapper settings and retry")
	}

	p := mpb.New(
		mpb.WithWidth(60),
		mpb.WithRefreshRate(180*time.Millisecond),
	)
	totalBar := p.Add(int64(len(sc.ScrapperManager.SrcAccounts)),
		mpb.NewBarFiller("[=>-|"),
		mpb.BarRemoveOnComplete(),
		mpb.PrependDecorators(
			decor.CountersNoUnit("%d / %d"),
		),
		mpb.AppendDecorators(
			decor.Percentage(),
		),
	)

	for _, srcUsername := range sc.ScrapperManager.SrcAccounts {
		logrus.Debugf("Fetch from '%s' user", srcUsername)
		finded, err := sc.navigateUserFollowersList(srcUsername)
		if !finded || err != nil {
			totalBar.IncrBy(1)
			continue
		}

		userBar := p.Add(int64(sc.ScrapperManager.Quantity),
			mpb.NewBarFiller("[=>-|"),
			mpb.BarRemoveOnComplete(),
			mpb.PrependDecorators(
				decor.Name(fmt.Sprintf("Scrapping users from %s account: ", srcUsername)),
				decor.CountersNoUnit("%d / %d"),
			),
			mpb.AppendDecorators(
				decor.Percentage(),
			),
		)

		// Scrap users until it has the right amount defined in ScrapperManager.Quantity by the user
		var scrappedUsers []selenium.WebElement
		for len(scrappedUsers) < sc.ScrapperManager.Quantity {
			if len(scrappedUsers) != 0 {
				// Scroll to the end of the list to gather more followers from ig
				_, err = sc.SeleniumStruct.WebDriver.ExecuteScript("window.scrollTo(0, document.body.scrollHeight);", nil)
				if err != nil {
					logrus.Warnf(
						"Error during followers dialog box scroll for '%s' user. The user certainly did not have enough followers for the request",
						srcUsername,
					)
					userBar.Abort(true)
					break
				}
			}
			randomSleepCustom(3, 4)
			scrappedUsers, err = sc.SeleniumStruct.GetElements("//*/li/div/div/div/div/a", "xpath")
			if err != nil {
				logrus.Errorf(
					"Error during users scrapping from followers dialog box for '%s' user",
					srcUsername,
				)
				userBar.Abort(true)
				break
			}
			scrappedUsers = sc.Blacklist.FilterScrappedUsers(scrappedUsers)
			userBar.SetCurrent(int64(len(scrappedUsers)))
			logrus.Debugf("Users count finded: %d", len(scrappedUsers))
		}

		if len(scrappedUsers) != 0 {
			for _, user := range scrappedUsers {
				username, err := user.Text()
				if err == nil {
					igUsers = append(igUsers, username)
				}
			}
		}

		logrus.Debugf("Scrapped users: %v\n", igUsers)
		if !userBar.Completed() {
			userBar.Abort(true)
		}
		totalBar.IncrBy(1)
	}
	p.Wait()
	if len(igUsers) == 0 {
		return nil, errors.New("Empty users result")
	}
	return igUsers, nil
}

// Go to user followers list with webdriver
func (sc *IGopher) navigateUserFollowersList(srcUsername string) (bool, error) {
	// Navigate to Instagram user page
	if err := sc.SeleniumStruct.WebDriver.Get(fmt.Sprintf("https://www.instagram.com/%s/?hl=en", srcUsername)); err != nil {
		logrus.Warnf("Requested user '%s' doesn't exist, skip it", srcUsername)
		return false, errors.New("Error during access to requested user")
	}
	randomSleepCustom(1, 3)
	// Access to followers list view
	if find, err := sc.SeleniumStruct.WaitForElement("//*[@id=\"react-root\"]/div/div/section/main/div/ul/li[2]/a", "xpath", 10); err == nil && find {
		elem, _ := sc.SeleniumStruct.GetElement("//*[@id=\"react-root\"]/div/div/section/main/div/ul/li[2]/a", "xpath")
		elem.Click()
		logrus.Debug("Clicked on user followers list")
	} else {
		return true, errors.New("Error during access to user followers list")
	}

	return true, nil
}
