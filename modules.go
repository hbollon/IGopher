package igopher

import (
	"encoding/csv"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tebeka/selenium"
)

const (
	fileBlacklistPath = "data/blacklist.csv"
)

/* Quota manager */

// QuotaManager data
type QuotaManager struct {
	// HourTimestamp: hourly timestamp used to handle hour limitations
	HourTimestamp time.Time
	// DayTimestamp: daily timestamp used to handle day limitations
	DayTimestamp time.Time
	// DmSent: quantity of dm sent in the last hour
	DmSent int
	// DmSentDay: quantity of dm sent in the last day
	DmSentDay int
	// MaxDmHour: maximum dm quantity per hour
	MaxDmHour int `yaml:"dm_per_hour" validate:"numeric"`
	// MaxDmDay: maximum dm quantity per day
	MaxDmDay int `yaml:"dm_per_day" validate:"numeric"`
	// Activated: quota manager activation boolean
	Activated bool `yaml:"activated"`
}

// InitializeQuotaManager initialize Quota manager with user settings
func (qm *QuotaManager) InitializeQuotaManager() {
	qm.HourTimestamp = time.Now()
	qm.DayTimestamp = time.Now()
}

// ResetDailyQuotas reset daily dm counter and update timestamp
func (qm *QuotaManager) ResetDailyQuotas() {
	qm.DmSentDay = 0
	qm.DayTimestamp = time.Now()
}

// ResetHourlyQuotas reset hourly dm counter and update timestamp
func (qm *QuotaManager) ResetHourlyQuotas() {
	qm.DmSent = 0
	qm.HourTimestamp = time.Now()
}

// AddDm report to the manager a message sending. It increment dm counter and check if quotas are still valid.
func (qm *QuotaManager) AddDm() {
	qm.DmSent++
	qm.DmSentDay++
	qm.CheckQuotas()
}

// CheckQuotas check if quotas have not been exceeded and pauses the program otherwise.
func (qm *QuotaManager) CheckQuotas() {
	// Hourly quota checking
	if qm.DmSent >= qm.MaxDmHour && qm.Activated {
		if time.Since(qm.HourTimestamp).Seconds() < 3600 {
			sleepDur := 3600 - time.Since(qm.HourTimestamp).Seconds()
			logrus.Infof("Hourly quota reached, sleeping %f seconds...", sleepDur)
			time.Sleep(time.Duration(sleepDur) * time.Second)
		} else {
			qm.ResetHourlyQuotas()
			logrus.Info("Hourly quotas resetted.")
		}
	}
	// Daily quota checking
	if qm.DmSentDay >= qm.MaxDmDay && qm.Activated {
		if time.Since(qm.DayTimestamp).Seconds() < 86400 {
			sleepDur := 86400 - time.Since(qm.DayTimestamp).Seconds()
			logrus.Infof("Daily quota reached, sleeping %f seconds...", sleepDur)
			time.Sleep(time.Duration(sleepDur) * time.Second)
		} else {
			qm.ResetDailyQuotas()
			logrus.Info("Daily quotas resetted.")
		}
	}
}

/* Schedule manager */

// SchedulerManager data
type SchedulerManager struct {
	// BeginAt: Begin time setting
	BeginAt string `yaml:"begin_at" validate:"contains=:"`
	// EndAt: End time setting
	EndAt string `yaml:"end_at" validate:"contains=:"`
	// BeginAtTimestamp: begin timestamp
	BeginAtTimestamp time.Time
	// EndAtTimestamp: end timestamp
	EndAtTimestamp time.Time
	// Activated: quota manager activation boolean
	Activated bool `yaml:"activated"`
}

// InitializeScheduler convert string time from config to time.Time instances
func (s *SchedulerManager) InitializeScheduler() error {
	ttBegin, err := time.Parse("15:04", strings.TrimSpace(s.BeginAt))
	if err != nil {
		return err
	}
	s.BeginAtTimestamp = ttBegin
	ttEnd, err := time.Parse("15:04", strings.TrimSpace(s.EndAt))
	if err != nil {
		return err
	}
	s.EndAtTimestamp = ttEnd
	return nil
}

// CheckTime check scheduler and pause the bot if it's not working time
func (s *SchedulerManager) CheckTime() error {
	if !s.Activated {
		return nil
	}
	res, err := s.isWorkingTime()
	if err == nil {
		if res {
			return nil
		}
		logrus.Info("Reached end of service. Sleeping...")
		for {
			if res, _ = s.isWorkingTime(); res {
				break
			}
			if BotStruct.exitCh != nil {
				select {
				case <-BotStruct.hotReloadCallback:
					if err = BotStruct.HotReload(); err != nil {
						logrus.Errorf("Bot hot reload failed: %v", err)
						BotStruct.hotReloadCallback <- false
					} else {
						logrus.Info("Bot hot reload successfully.")
						BotStruct.hotReloadCallback <- true
					}
					break
				case <-BotStruct.exitCh:
					logrus.Info("Bot process successfully stopped.")
					return errStopBot
				default:
					break
				}
			}
			time.Sleep(10 * time.Second)
		}
		logrus.Info("Back to work!")
	}
	return nil
}

// Check if current time is between scheduler working interval
func (s *SchedulerManager) isWorkingTime() (bool, error) {
	if s.BeginAtTimestamp.Equal(s.EndAtTimestamp) {
		return false, errors.New("Bad scheduler configuration")
	}
	currentTime := time.Date(0, time.January, 1, time.Now().Hour(), time.Now().Minute(), 0, 0, time.Local)
	if s.BeginAtTimestamp.Before(s.EndAtTimestamp) {
		return !currentTime.Before(s.BeginAtTimestamp) && !currentTime.After(s.EndAtTimestamp), nil
	}
	return !s.BeginAtTimestamp.After(currentTime) || !s.EndAtTimestamp.Before(currentTime), nil
}

/* Blacklist manager */

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
