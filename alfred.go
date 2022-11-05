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
	search(arg, token())
	wf.WarnEmpty("No results", "Try another search?")
	wf.SendFeedback()
}

// handleOpts contains the flag/option logic.
func handleOpts() {
	if opts.Link != "" {
		exec.Command("open", opts.Link).Run()
		os.Exit(0)
	}
	if opts.Token {
		token, err := newToken(config)
		if err != nil {
			wf.Fatal(fmt.Sprintf("Error retrieving new OAuth2 token: %v", err))
		}
		cacheToken(kc, token)
		os.Exit(0)
	}
}

// token either retrieves a cached token or prompts the user to request a new one
// via 3-legged OAuth2 flow.
func token() *oauth2.Token {
	token, err := cachedToken(kc)
	if err != nil {
		wf.NewItem("Error retrieving cached OAuth2 token").
			Subtitle("Hit return to request a new one").
			Icon(aw.IconWarning).
			Arg("--token").
			Valid(true)
		wf.SendFeedback()
		os.Exit(0)
	}
	return token
}

// search uses the results from the Google Drive API to populate a script filter.
func search(arg string, token *oauth2.Token) {
	results, err := searchDrive(arg, config, token)
	if err != nil {
		wf.Fatal(fmt.Sprintf("Error searching Google Drive: %v", err))
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
