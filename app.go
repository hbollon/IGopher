package main

import "github.com/hbollon/go-instadm/lib"

// WebDriver is the main Selenium instance used by this bot
var WebDriver Selenium

func main() {
	// Download dependencies
	lib.DownloadDependencies(true, true)

	// var webDriver Selenium
	WebDriver.InitializeWebDriver()
}
