package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os/exec"

	"golang.org/x/oauth2"
)

// newToken starts the OAuth2 flow in order to get a token.
func newToken(config *oauth2.Config) (*oauth2.Token, error) {
	codeCh, err := startWebServer()
	if err != nil {
		return nil, err
	}

	// TODO: fix state to protect against CSRF - https://github.com/douglasmakey/oauth2-example
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	if exec.Command("open", authURL).Start() != nil {
		return nil, err
	}

	fmt.Println("Your browser has been opened to an authorization URL.",
		" This program will resume once authorization has been provided.")
	fmt.Println(authURL)

	// Wait for the web server to get the code.
	code := <-codeCh
	return exchangeToken(config, code)
}

// startWebServer starts a web server and listens for OAuth2 code
// returned as part of the three-legged auth flow.
func startWebServer() (chan string, error) {
	listener, err := net.Listen("tcp", "localhost:1337")
	if err != nil {
		return nil, err
	}

	codeCh := make(chan string)
	go http.Serve(listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code := r.FormValue("code")
		codeCh <- code // send code to OAuth flow
		listener.Close()
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Received code: %v\nYou can now safely close this browser window and start using the workflow.", code)
	}))
	return codeCh, nil
}

// exchangeToken swaps the authorization code for an access token.
func exchangeToken(config *oauth2.Config, code string) (*oauth2.Token, error) {
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}
	return token, nil
}
