package controllers

// info struct is used to display informations to
// the user via info.html template
type info struct {
	BrowserTitle string // title displayed in browser
	Title        string // title displayed in card
	Msg          string // message displayed below the title card
	LoggedInUser User // currently logged in user
}
