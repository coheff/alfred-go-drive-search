package main

import (
	"fmt"
	"os"
	"os/exec"

	alco "github.com/coheff/al-co"
	aw "github.com/deanishe/awgo"
	"golang.org/x/oauth2"
)

// run contains the main logic for the application. It parses the input query
// and flags, retrieves/requests a token, and triggers a search.
func run() {
	arg := alco.ParseQuery(wf.Args()[0], &opts)
	handleOpts()
	search(arg, cachedToken())
	wf.WarnEmpty("No results", "Try another search?")
	wf.SendFeedback()
}

// handleOpts contains the flag/option logic.
func handleOpts() {
	if opts.Link != "" {
		err := exec.Command("open", opts.Link).Run()
		if err != nil {
			wf.Fatal(fmt.Sprintf("Error opening Google Drive link: %v", err))
		}
		os.Exit(0)
	}
	if opts.Token {
		token, err := newToken(config)
		if err != nil {
			wf.Fatal(fmt.Sprintf("Error retrieving new OAuth2 token: %v", err))
		}
		cacheToken(token)
		os.Exit(0)
	}
}

// search uses the results from the Google Drive API to populate a script filter.
// If there's an error while searching prompt the user to attempt token refresh
// (this solves the majority of issues).
func search(arg string, token *oauth2.Token) {
	results, err := searchDrive(arg, config, token)
	if err != nil {
		wf.NewItem(fmt.Sprintf("Error searching Google Drive: %v", err)).
			Subtitle("Hit return to request a new token").
			Icon(aw.IconWarning).
			Arg("--token").
			Valid(true)
		wf.SendFeedback()
		os.Exit(0)
	}
	for _, result := range results {
		wf.
			NewItem(result.name).
			Subtitle(fmt.Sprintf("%s - %s", result.owners, result.modified)).
			Icon(&aw.Icon{
				Value: result.icon,
			}).
			Arg("--link=" + result.link).
			Valid(true)
	}
}
