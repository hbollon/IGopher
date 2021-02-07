package igopher

import (
	"github.com/sirupsen/logrus"
	"github.com/tebeka/selenium"
)

// SimulateHandWriting simulate human writing by typing input string character by character with random interruptions
// between letters
func SimulateHandWriting(element selenium.WebElement, input string) bool {
	var err error
	if err = element.Click(); err == nil {
		for _, c := range input {
			if err = element.SendKeys(string(c)); err != nil {
				logrus.Debug("Unable to send key during message typing")
				logrus.Errorf("Error during message sending: %v", err)
				return false
			}
			randomSleepCustom(0.25, 1.0)
		}
		return true
	}
	logrus.Debug("Can't click on user searchbar")
	logrus.Errorf("Error during message sending: %v", err)
	return false
}
