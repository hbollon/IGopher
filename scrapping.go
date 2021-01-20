package igopher

import (
	"github.com/ahmdrz/goinsta/v2"
	"github.com/sirupsen/logrus"
)

// Scrapper struct store goinsta.Instagram instance pointer and its corresponding configuration
type Scrapper struct {
	Instance *goinsta.Instagram
	Config   ScrapperConfig `yaml:"config"`
}

// ScrapperConfig store scrapper configuration for user fetching
// It also store fetched usernames
type ScrapperConfig struct {
	SrcAccounts     []string `yaml:"src_accounts"`
	FetchedAccounts []string
	Quantity        int `yaml:"fetch_quantity" validate:"numeric"`
}

func (sc *Scrapper) CreateScrapperInstance(a *Account) {
	sc.Instance = goinsta.New(
		a.Username,
		a.Password,
	)
}

func (sc *Scrapper) ScrapperLogin() error {
	if err := sc.Instance.Login(); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

// func (ig *Scrapper) FetchUsersFromUserFollowers() error {

// }
