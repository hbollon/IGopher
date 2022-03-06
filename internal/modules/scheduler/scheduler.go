package scheduler

import (
	"errors"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Manager data
type Manager struct {
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
func (s *Manager) InitializeScheduler() error {
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
func (s *Manager) CheckTime() error {
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
			// if engine.BotStruct.ExitCh != nil {
			// 	select {
			// 	case <-engine.BotStruct.HotReloadCallback:
			// 		if err = engine.BotStruct.HotReload(); err != nil {
			// 			logrus.Errorf("Bot hot reload failed: %v", err)
			// 			engine.BotStruct.HotReloadCallback <- false
			// 		} else {
			// 			logrus.Info("Bot hot reload successfully.")
			// 			engine.BotStruct.HotReloadCallback <- true
			// 		}
			// 		break
			// 	case <-engine.BotStruct.ExitCh:
			// 		logrus.Info("Bot process successfully stopped.")
			// 		return bot.ErrStopBot
			// 	default:
			// 		break
			// 	}
			// }
			time.Sleep(10 * time.Second)
		}
		logrus.Info("Back to work!")
	}
	return nil
}

// Check if current time is between scheduler working interval
func (s *Manager) isWorkingTime() (bool, error) {
	if s.BeginAtTimestamp.Equal(s.EndAtTimestamp) {
		return false, errors.New("Bad scheduler configuration")
	}
	currentTime := time.Date(0, time.January, 1, time.Now().Hour(), time.Now().Minute(), 0, 0, time.Local)
	if s.BeginAtTimestamp.Before(s.EndAtTimestamp) {
		return !currentTime.Before(s.BeginAtTimestamp) && !currentTime.After(s.EndAtTimestamp), nil
	}
	return !s.BeginAtTimestamp.After(currentTime) || !s.EndAtTimestamp.Before(currentTime), nil
}
