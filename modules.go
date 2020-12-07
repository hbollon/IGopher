package instadm

import (
	"time"
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
	MaxDmHour int `yaml:"dm_per_hour"`
	// MaxDmDay: maximum dm quantity per day
	MaxDmDay int `yaml:"dm_per_day"`
	// Activated: quota manager activation boolean
	Activated bool `yaml:"activated"`
}

// InitializeQuotaManager initialize Quota manager with user settings
func (qm *QuotaManager) InitializeQuotaManager() {
	qm.HourTimestamp = time.Now()
	qm.DayTimestamp = time.Now()
}
