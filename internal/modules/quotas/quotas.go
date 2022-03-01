package quotas

import (
	"time"

	"github.com/sirupsen/logrus"
)

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
