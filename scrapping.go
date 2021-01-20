package igopher

import (
	"github.com/ahmdrz/goinsta/v2"
	"github.com/sirupsen/logrus"
)

type scrapper struct {
	instance *goinsta.Instagram
}

var insta scrapper

func (c *ClientConfig) CreateScrapperInstance() {
	insta.instance = goinsta.New(
		c.BotConfig.UserAccount.Username,
		c.BotConfig.UserAccount.Password,
	)
}

func (ig *scrapper) ScrapperLogin() error {
	if err := ig.instance.Login(); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
