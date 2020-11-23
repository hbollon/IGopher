package main

// WebDriver is the main Selenium instance used by this bot
var WebDriver Selenium

func main() {
	// Download dependencies
	lib.DownloadDependencies(true, true)

	// var webDriver Selenium
	WebDriver.InitializeSelenium()
	WebDriver.InitFirefoxWebDriver()
	defer WebDriver.CloseSelenium()
}
