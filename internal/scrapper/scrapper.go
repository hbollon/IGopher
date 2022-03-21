package scrapper

import (
	"errors"
	"fmt"
	"time"

	"github.com/hbollon/igopher/internal/config/types"
	"github.com/hbollon/igopher/internal/utils"
	"github.com/hbollon/igopher/internal/xpath"
	"github.com/sirupsen/logrus"
	"github.com/tebeka/selenium"
	"github.com/vbauerster/mpb/v6"
	"github.com/vbauerster/mpb/v6/decor"
)

// FetchUsersFromUserFollowers scrap username list from users followers.
// Source accounts and quantity are set by the bot user.
func FetchUsersFromUserFollowers(bot *types.IGopher) ([]string, error) {
	logrus.Info("Fetching users from user's followers...")

	var igUsers []string
	// Valid configuration checking before fetching process
	if len(bot.ScrapperManager.SrcAccounts) == 0 || bot.ScrapperManager.SrcAccounts == nil {
		return nil, errors.New("No source users are set, please check your scrapper settings and retry")
	}
	if bot.ScrapperManager.Quantity <= 0 {
		return nil, errors.New("Scrapping quantity is null or negative, please check your scrapper settings and retry")
	}

	p := mpb.New(
		mpb.WithWidth(60),
		mpb.WithRefreshRate(180*time.Millisecond),
	)
	totalBar := p.Add(int64(len(bot.ScrapperManager.SrcAccounts)),
		mpb.NewBarFiller("[=>-|"),
		mpb.BarRemoveOnComplete(),
		mpb.PrependDecorators(
			decor.CountersNoUnit("%d / %d"),
		),
		mpb.AppendDecorators(
			decor.Percentage(),
		),
	)

	for _, srcUsername := range bot.ScrapperManager.SrcAccounts {
		logrus.Debugf("Fetch from '%s' user", srcUsername)
		finded, err := navigateUserFollowersList(bot, srcUsername)
		if !finded || err != nil {
			totalBar.IncrBy(1)
			continue
		}

		userBar := p.Add(int64(bot.ScrapperManager.Quantity),
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
		for len(scrappedUsers) < bot.ScrapperManager.Quantity {
			if len(scrappedUsers) != 0 {
				// Scroll to the end of the list to gather more followers from ig
				_, err = bot.SeleniumStruct.WebDriver.ExecuteScript("window.scrollTo(0, document.body.scrollHeight);", nil)
				if err != nil {
					logrus.Warnf(
						"Error during followers dialog box scroll for '%s' user. The user certainly did not have enough followers for the request",
						srcUsername,
					)
					userBar.Abort(true)
					break
				}
			}
			utils.RandomSleepCustom(3, 4)
			scrappedUsers, err = bot.SeleniumStruct.GetElements(xpath.XPathSelectors["profile_followers_list"], "xpath")
			if err != nil {
				logrus.Errorf(
					"Error during users scrapping from followers dialog box for '%s' user",
					srcUsername,
				)
				userBar.Abort(true)
				break
			}
			scrappedUsers = bot.Blacklist.FilterScrappedUsers(scrappedUsers)
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
func navigateUserFollowersList(bot *types.IGopher, srcUsername string) (bool, error) {
	// Navigate to Instagram user page
	if err := bot.SeleniumStruct.WebDriver.Get(fmt.Sprintf("https://www.instagram.com/%s/?hl=en", srcUsername)); err != nil {
		logrus.Warnf("Requested user '%s' doesn't exist, skip it", srcUsername)
		return false, errors.New("Error during access to requested user")
	}
	utils.RandomSleepCustom(1, 3)
	// Access to followers list view
	if find, err := bot.SeleniumStruct.WaitForElement(xpath.XPathSelectors["profile_followers_button"], "xpath", 10); err == nil && find {
		elem, _ := bot.SeleniumStruct.GetElement(xpath.XPathSelectors["profile_followers_button"], "xpath")
		elem.Click()
		logrus.Debug("Clicked on user followers list")
	} else {
		return true, errors.New("Error during access to user followers list")
	}

	return true, nil
}
