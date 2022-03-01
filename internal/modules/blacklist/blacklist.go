package blacklist

import (
	"encoding/csv"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/tebeka/selenium"
)

const (
	fileBlacklistPath = "data/blacklist.csv"
)

// BlacklistManager data
type BlacklistManager struct {
	// BlacklistedUsers: list of all blacklisted usernames
	BlacklistedUsers [][]string
	// Activated: quota manager activation boolean
	Activated bool `yaml:"activated"`
}

// InitializeBlacklist check existence of the blacklist csv file and initialize it if it doesn't exist.
func (bm *BlacklistManager) InitializeBlacklist() error {
	var err error
	// Check if blacklist csv exist
	_, err = os.Stat(fileBlacklistPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Create data folder if not exist
			if _, err = os.Stat("data/"); os.IsNotExist(err) {
				os.Mkdir("data/", os.ModePerm)
			}
			// Create and open csv blacklist
			var f *os.File
			f, err = os.OpenFile(fileBlacklistPath, os.O_RDWR|os.O_CREATE, 0755)
			if err != nil {
				return err
			}
			defer f.Close()
			// Write csv header
			writer := csv.NewWriter(f)
			err = writer.Write([]string{"Username"})
			defer writer.Flush()
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		// Open existing blacklist and recover blacklisted usernames
		f, err := os.OpenFile(fileBlacklistPath, os.O_RDONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		reader := csv.NewReader(f)
		bm.BlacklistedUsers, err = reader.ReadAll()
		if err != nil {
			return err
		}
	}

	return nil
}

// AddUser add argument username to the blacklist
func (bm *BlacklistManager) AddUser(user string) {
	bm.BlacklistedUsers = append(bm.BlacklistedUsers, []string{user})
	f, err := os.OpenFile(fileBlacklistPath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		logrus.Errorf("Failed to blacklist current user: %v", err)
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	err = writer.Write([]string{user})
	defer writer.Flush()
	if err != nil {
		logrus.Errorf("Failed to blacklist current user: %v", err)
	}
}

// IsBlacklisted check if the given user is already blacklisted
func (bm *BlacklistManager) IsBlacklisted(user string) bool {
	for _, username := range bm.BlacklistedUsers {
		if username[0] == user {
			return true
		}
	}
	return false
}

// FilterScrappedUsers remove blacklisted users from WebElement slice and return it
func (bm *BlacklistManager) FilterScrappedUsers(users []selenium.WebElement) []selenium.WebElement {
	var filteredUsers []selenium.WebElement
	for _, user := range users {
		username, err := user.Text()
		if !bm.IsBlacklisted(username) && err == nil {
			filteredUsers = append(filteredUsers, user)
		}
	}
	return filteredUsers
}
