package xpath

var (
	// XPathSelectors is a map regrouping all xpaths used by igopher
	// to find elements on the web page.
	// Also contains some elements names.
	XPathSelectors = map[string]string{
		// Login page
		"login_username":                 "username",
		"login_password":                 "password",
		"login_button":                   "//button[text()='Log In']",
		"login_alternate_button":         "//button/*[text()='Log In']",
		"login_accept_cookies":           "//button[text()='Accept All' or text()='Allow essential and optional cookies']",
		"login_alternate_accept_cookies": "//button[text()='Allow All Cookies']",
		"login_information_saving":       "//*[@aria-label='Home'] | //button[text()='Save Info'] | //button[text()='Not Now']",

		// DM related elements
		"dm_user_search": "//section/div[2]/div/div[1]/div/div[2]/input",
		"dm_next_button": "//button/*[text()='Next']",
		"dm_placeholder": "//textarea[@placeholder]",
		"dm_send_button": "//button[text()='Send']",

		// Profile related elements
		"profile_followers_button": "//section/main/div/ul/li[2]/a",
		"profile_followers_list":   "//*/li/div/div/div/div/a",
	}
)
