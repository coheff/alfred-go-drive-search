package main

import (
	"log"
	"net/url"
	"os/exec"

	ac "github.com/coheff/al-co"
	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/keychain"
)

var (
	Kc   = keychain.New(Wf.BundleID())
	opts struct {
		FullText string `short:"f" long:"full" description:"Search full text"`
	}
)

type Result struct {
	Title    string
	Subtitle string
	Arg      string
	Icon     string
}

func Run() {
	q := ac.ParseQuery(Wf.Args()[0], &opts)

	// if query is a link open it
	if _, err := url.ParseRequestURI(q); err == nil {
		err := exec.Command("open", q).Run()
		if err != nil {
			log.Fatalf("Error opening Google Drive link: %v", err)
		}
	}

	for _, result := range SearchDrive(q) {
		Wf.
			NewItem(result.Title).
			Subtitle(result.Subtitle).
			Arg(result.Arg).
			Icon(&aw.Icon{
				Value: result.Icon,
			}).
			Valid(true)
	}

	Wf.WarnEmpty("No results", "Try another search?")
	Wf.SendFeedback()
}
