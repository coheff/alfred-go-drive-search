package main

import (
	aw "github.com/deanishe/awgo"
	awkc "github.com/deanishe/awgo/keychain"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
)

var (
	// OAuth Config
	config *oauth2.Config

	// Option flags
	opts struct {
		FullText string `short:"f" long:"full" description:"Search full text"`
		Link     string `long:"link" description:"Open Google Drive link"`
		Token    bool   `long:"token" description:"Trigger OAuth2 flow to retrieve new token"`
	}

	// Workflow
	kc *awkc.Keychain
	wf *aw.Workflow
)

func init() {
	wf = aw.New()
	kc = awkc.New(wf.BundleID())

	config = &oauth2.Config{
		ClientID:     wf.Config.Get("client_id"),
		ClientSecret: wf.Config.Get("client_secret"),
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:1337/callback",
		Scopes:       []string{drive.DriveMetadataReadonlyScope},
	}
}

func main() {
	wf.Run(run)
}
